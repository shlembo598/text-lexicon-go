package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/shlembo598/text-lexicon-go/internal/auth"
	"github.com/shlembo598/text-lexicon-go/internal/config"
	"github.com/shlembo598/text-lexicon-go/pkg/utils"
	"github.com/shlembo598/text-lexicon-go/pkg/utils/httpErrors"
)

// JWT way of auth using Authorization header
func (mw *MiddlewareManager) AuthJWTMiddleware(authUC auth.UseCase, cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerHeader := c.Request().Header.Get("Authorization")

			slog.Info("auth middleware", slog.String("bearerHeader", bearerHeader))

			if bearerHeader != "" {
				headerParts := strings.Split(bearerHeader, " ")
				if len(headerParts) != 2 {
					slog.Error("auth middleware", slog.String("headerParts", "len(headerParts) != 2"))

					return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
				}

				tokenString := headerParts[1]

				if err := mw.validateJWTToken(tokenString, authUC, c, cfg); err != nil {
					slog.Error("middleware validateJWTToken", slog.String("headerJWT", err.Error()))

					return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
				}

				return next(c)
			}

			slog.Error("authorization header is missing")

			return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}
	}
}

func (mw *MiddlewareManager) validateJWTToken(
	tokenString string, authUC auth.UseCase, c echo.Context, cfg *config.Config,
) error {
	if tokenString == "" {
		return httpErrors.InvalidJWTToken
	}

	token, err := jwt.Parse(
		tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
			}
			secret := []byte(cfg.Server.JwtSecretKey)
			return secret, nil
		},
	)
	if err != nil {
		return err
	}

	if !token.Valid {
		return httpErrors.InvalidJWTToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["id"].(string)
		if !ok {
			return httpErrors.InvalidJWTClaims
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return err
		}

		u, err := authUC.GetByID(c.Request().Context(), userUUID)
		if err != nil {
			return err
		}

		c.Set("user", u)

		ctx := context.WithValue(c.Request().Context(), utils.UserCtxKey{}, u)

		c.SetRequest(c.Request().WithContext(ctx))
	}
	return nil
}
