package utils

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"

	"github.com/shlembo598/text-lexicon-go/pkg/logger/sl"
)

// ReqIDCtxKey is a key used for the Request ID in context
type ReqIDCtxKey struct{}

// Get request id from echo context
func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// Get user ip address
func GetIPAddress(c echo.Context) string {
	return c.Request().RemoteAddr
}

// Get context  with request id
func GetRequestCtx(c echo.Context) context.Context {
	return context.WithValue(c.Request().Context(), ReqIDCtxKey{}, GetRequestID(c))
}

// Read request body and validate
func ReadRequest(ctx echo.Context, request interface{}) error {
	if err := ctx.Bind(request); err != nil {
		return err
	}
	return validate.StructCtx(ctx.Request().Context(), request)
}

// Error response with logging error for echo context
func LogResponseError(ctx echo.Context, err error) {
	slog.Error(
		"response error",
		slog.String("RequestID", GetRequestID(ctx)),
		slog.String("IPAddress", GetIPAddress(ctx)),
		sl.Err(err),
	)
}
