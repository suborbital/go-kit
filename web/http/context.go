package http

import (
	"github.com/labstack/echo/v4"
)

const noRequestID = "unknown-request"

type Context struct {
	echo.Context
}

// RequestID will return the request ID stored in the echo.HeaderXRequestID header, or "unknown-request" if the value of
// the header was an empty string.
func (c *Context) RequestID() string {
	rid := c.Request().Header.Get(echo.HeaderXRequestID)
	if rid == "" {
		return noRequestID
	}

	return rid
}
