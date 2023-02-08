package mid

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// UUIDRequestID configures echo's built in request id middleware so that the ID generated is an UUIDv4, and the
// generated request ID is added to the following three parts of the request:
// - the echo.HeaderXRequestID header, this is by default
// - echo.Context's own Set method with the RequestIDKey key
// - request context with key RequestIDKey
//
// Value of the RequestIDKey is "requestID", however for stability, use the exported constant when referring to it.
func UUIDRequestID() echo.MiddlewareFunc {
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.New().String()
		},
	})
}
