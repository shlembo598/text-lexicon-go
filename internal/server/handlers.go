package server

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/shlembo598/text-lexicon-go/docs"
	authHttp "github.com/shlembo598/text-lexicon-go/internal/auth/delivery/http"
	authRepository "github.com/shlembo598/text-lexicon-go/internal/auth/repository"
	authUseCase "github.com/shlembo598/text-lexicon-go/internal/auth/usecase"
	apiMiddlewares "github.com/shlembo598/text-lexicon-go/internal/middleware"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	// Init repositories
	authRepo := authRepository.NewAuthRepository(s.db)

	// Init useCases
	authUC := authUseCase.NewAuthUserCase(s.cfg, authRepo)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC)

	// Init middleware
	mw := apiMiddlewares.NewMiddlewareManager(authUC, s.cfg, []string{"*"})

	e.Use(mw.RequestLoggerMiddleware)

	docs.SwaggerInfo.Title = "Text lexicon REST API"
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Use(
		middleware.CORSWithConfig(
			middleware.CORSConfig{
				AllowOrigins: []string{"*"},
				AllowHeaders: []string{
					echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID,
				},
			},
		),
	)
	e.Use(
		middleware.RecoverWithConfig(
			middleware.RecoverConfig{
				StackSize:         1 << 10, // 1 KB
				DisablePrintStack: true,
				DisableStackAll:   true,
			},
		),
	)
	e.Use(middleware.RequestID())
	e.Use(
		middleware.GzipWithConfig(
			middleware.GzipConfig{
				Level: 5,
				Skipper: func(c echo.Context) bool {
					return strings.Contains(c.Request().URL.Path, "swagger")
				},
			},
		),
	)
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))
	if s.cfg.Server.Debug {
		e.Use(mw.DebugMiddleware)
	}

	v1 := e.Group("/api/v1")

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")

	authHttp.MapAuthRoutes(authGroup, authHandlers, mw, authUC, s.cfg)

	health.GET(
		"", func(c echo.Context) error {
			slog.Info("health check ")
			return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
		},
	)

	return nil
}
