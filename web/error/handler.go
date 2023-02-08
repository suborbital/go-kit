package error

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	kitHttp "github.com/suborbital/go-kit/web/http"
)

func Handler(logger zerolog.Logger) echo.HTTPErrorHandler {
	ll := logger.With().Str("middleware", "errorHandler").Logger()
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		ll.Err(err).
			Str("requestID", kitHttp.RID(c)).
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

		// Issue #1426
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
