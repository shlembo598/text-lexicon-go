package auth

import (
	"github.com/labstack/echo/v4"
)

type Handlers interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	GetUserByID() echo.HandlerFunc
	GetMe() echo.HandlerFunc
}
