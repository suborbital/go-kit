package mid

import (
	"github.com/labstack/echo/v4"
	"github.com/suborbital/go-kit/web/http"
)

// CustomContext wraps the default echo.Context into our own decorated struct, http.Context. The only addition in ours
// is a .RequestID() convenience method that returns the request ID stored in the echo.HeaderXRequestID header.
func CustomContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &http.Context{Context: c}
			return next(cc)
		}
	}
}
