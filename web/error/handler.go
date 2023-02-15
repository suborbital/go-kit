package error

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	kitHttp "github.com/suborbital/go-kit/web/http"
)

// Handler is a modified version of echo's own DefaultHTTPErrorHandler function. It allows us to do the following:
// - log a committed response, both that return an error, and ones that don't
// - log all internal errors without exposing them to the client
// - modify the response json to also include the status code in the response body
func Handler(logger zerolog.Logger) echo.HTTPErrorHandler {
	ll := logger.With().Str("middleware", "errorHandler").Logger()
	return func(err error, c echo.Context) {
		rid := kitHttp.RID(c)

		if c.Response().Committed {
			ll.Err(err).Str("requestID", rid).Msg("response already committed")
			return
		}

		ll.Err(err).
			Str("requestID", rid).
			Msg("request returned an error")

		he, ok := err.(*echo.HTTPError)
		if ok {
			if he.Internal != nil {
				if herr, ok := he.Internal.(*echo.HTTPError); ok {
					he = herr
				}
			}
		} else {
			he = &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			}
		}

		code := he.Code
		message := he.Message
		if m, ok := he.Message.(string); ok {
			if c.Echo().Debug {
				message = echo.Map{"status": code, "message": m, "error": err.Error()}
			} else {
				message = echo.Map{"status": code, "message": m}
			}
		}

		// Send response
		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(he.Code)
		} else {
			err = c.JSON(code, message)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
