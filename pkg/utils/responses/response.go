package responses

import (
	"github.com/labstack/echo/v4"
)

const (
	statusSuccess = "success"
	statusError   = "error"
)

const (
	ErrBadRequest         = "Bad request"
	ErrEmailAlreadyExists = "User with given email already exists"
	ErrNoSuchUser         = "User not found"
	ErrWrongCredentials   = "Wrong Credentials"
	ErrNotFound           = "Not Found"
	ErrUnauthorized       = "Unauthorized"
	ErrForbidden          = "Forbidden"
	ErrBadQueryParams     = "Invalid query params"
)

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Status int         `json:"code"`
	Error  interface{} `json:"error"`
}

func SuccessResponse(c echo.Context, code int, data interface{}) error {
	response := &Response{
		Status: statusSuccess,
		Data:   data,
		Error:  nil,
	}
	return c.JSON(code, response)
}

func ErrorResponse(c echo.Context, code int, error interface{}) error {
	response := &Response{
		Status: statusError,
		Data:   nil,
		Error: &APIError{
			Status: code,
			Error:  error,
		},
	}
	return c.JSON(code, response)
}
