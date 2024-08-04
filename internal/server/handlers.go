package server

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	authHttp "github.com/shlembo598/text-lexicon-go/internal/auth/delivery/http"
	authRepository "github.com/shlembo598/text-lexicon-go/internal/auth/repository"
	authUseCase "github.com/shlembo598/text-lexicon-go/internal/auth/usecase"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	// Init repositories
	authRepo := authRepository.NewAuthRepository(s.db)

	// Init useCases
	authUC := authUseCase.NewAuthUserCase(s.cfg, authRepo)

	// Init handlers
	authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC)
	// Init middleware

	v1 := e.Group("/api/v1")

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")

	authHttp.MapAuthRoutes(authGroup, authHandlers)

	health.GET(
		"", func(c echo.Context) error {
			slog.Info("health check ")
			return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
		},
	)

	return nil
}
