package http

import (
	"github.com/labstack/echo/v4"

	"github.com/shlembo598/text-lexicon-go/internal/auth"
)

func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers) {
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.PUT("/:user_id", h.Update())
	authGroup.DELETE("/:user_id", h.Delete())
	authGroup.GET("/:user_id", h.GetUserByID())
	authGroup.GET("/me", h.GetMe())
}
