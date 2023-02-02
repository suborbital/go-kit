package mid

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

// Logger middleware configures echo's built in middleware.RequestLoggerWithConfig middleware. The passed in
// zerolog.Logger is used to output the logs. For each request there are two log entries:
// - one for the incoming, before it gets passed onto the next handler, and
// - another one that logs info about the response before the client receives it.
//
// A non-empty slice of paths will be used to entirely skip logging for those routes. Do be careful to pass in the
// routes as you added them to echo. For example, if it's e.Get("/path/:something", ...), then the entry should look
// like []string{"/path/:something"}. If the path is inside a group with a path prefix, you need to add the full path,
// which includes the group prefix.
//
// Important to note that the Logger is configuring the echo middleware to handle errors first, and then return to the
// middleware. This has a side effect that middlewares further up the chain won't be able to change the response body
// or status code or headers. The upside is that any error handling happens between the "request started" and "request
// finished" log entries.
//
// The following fields are printed into the structured output for the incoming request along with a message of "request
// started":
//   - path - c.Path() - this is the one you pass to the request functions. Has the placeholders in.
//   - URI - c.Request().RequestURI - this is the actual request path that gets served. The placeholders are substituted
//   - requestID - c.Request().Header().Get(echo.HeaderXRequestID) - the RequestID middleware puts that in the header.
//     That middleware needs to wrap this one within it.
//   - method - c.Request().Method - the http verb used, e.g. GET, POST, PUT, etc...
//
// When the request is finished, the log will have a "request finished" message, and the following structured info:
//   - URI - request URI, same as above
//   - status - response status code
//   - requestID - same as above
//   - method - same as above
//   - latency - how long the request took. The time is between the logger middleware seeing the request go in, and the
//     logger middleware seeing the same request return, so it's a total time of every handler and middleware
//     inside of this middleware
func Logger(l zerolog.Logger, skipPaths []string) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		BeforeNextFunc: func(c echo.Context) {
			l.Info().
				Str("path", c.Path()).
				Str("URI", c.Request().RequestURI).
				Str("requestID", c.Response().Header().Get(echo.HeaderXRequestID)).
				Str("method", c.Request().Method).
				Msg("request started")
		},
		Skipper: func(c echo.Context) bool {
			path := c.Path()

			for _, sp := range skipPaths {
				if sp == path {
					return true
				}
			}

			return false
		},
		HandleError:  true,
		LogURI:       true,
		LogStatus:    true,
		LogRequestID: true,
		LogMethod:    true,
		LogLatency:   true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			l.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Str("requestID", v.RequestID).
				Str("method", v.Method).
				Dur("latency", v.Latency).
				Msg("request finished")

			return nil
		},
	})
}
