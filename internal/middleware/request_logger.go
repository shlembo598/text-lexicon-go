package middleware

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/shlembo598/text-lexicon-go/pkg/utils"
)

func (mw *MiddlewareManager) RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		err := next(ctx)

		req := ctx.Request()
		res := ctx.Response()
		status := res.Status
		size := res.Size
		s := time.Since(start).String()
		requestID := utils.GetRequestID(ctx)

		slog.Info(
			"logger middleware", slog.String("RequestID", requestID), slog.String("Method", req.Method),
			slog.String("URI", req.URL.String()), slog.String("Status", fmt.Sprint(status)),
			slog.String("Size", fmt.Sprint(size)), slog.String("Time", s),
		)

		return err
	}
}
