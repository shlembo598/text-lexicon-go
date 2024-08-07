package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"

	"github.com/labstack/echo/v4"
)

// Debug dump request middleware
func (mw *MiddlewareManager) DebugMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if mw.cfg.Server.Debug {
			dump, err := httputil.DumpRequest(c.Request(), true)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}
			slog.Info(
				"Request dump",
				fmt.Sprintf("\nbegin :--------------\n\n%s\n\nend :--------------", dump),
			)
		}
		return next(c)
	}
}
