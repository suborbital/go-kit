package http

import (
	"github.com/labstack/echo/v4"
)

const noRequestID = "unknown-request"

type Context struct {
	echo.Context
}

// RequestID will return the request ID stored in the echo.HeaderXRequestID header on the http.Response, or
// "unknown-request" if the value of the header was an empty string.
//
// Importantly this is on the response header because the request headers are expected to come from outside the service.
func (c *Context) RequestID() string {
	rid := c.Response().Header().Get(echo.HeaderXRequestID)
	if rid == "" {
		return noRequestID
	}

	return rid
}

// RID is a short function that the echo context can be passed into, which returns the request ID as it grabs it from
// the response header. The upside of this over the custom context is that a custom context is not required.
func RID(c echo.Context) string {
	rid := c.Response().Header().Get(echo.HeaderXRequestID)
	if rid == "" {
		return noRequestID
	}

	return rid
}
