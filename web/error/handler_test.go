package error

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/suborbital/go-kit/web/mid"
)

func TestHandler(t *testing.T) {
	mockRequestID := "1f550f47-0086-4f92-8d6a-1d5805b2e20e"

	tests := []struct {
		name             string
		handler          echo.HandlerFunc
		wantEmpty        bool
		wantFragments    []string
		wantStatusCode   int
		wantResponseBody string
	}{
		{
			name: "uncommitted response, no error, should stay silent",
			handler: func(c echo.Context) error {
				return c.String(http.StatusOK, "all good")
			},
			wantEmpty:        true,
			wantStatusCode:   http.StatusOK,
			wantResponseBody: "all good",
		},
		{
			name: "uncommitted response, error returned, should appear",
			handler: func(c echo.Context) error {
				return echo.NewHTTPError(http.StatusBadRequest, "ohno").SetInternal(errors.Wrap(errors.New("wrapped origin error"), "wrapping thing"))
			},
			wantEmpty: false,
			wantFragments: []string{
				"wrapped origin error",
				"wrapping thing",
				`"message":"request returned an error"`,
				mockRequestID,
			},
			wantStatusCode: http.StatusBadRequest,
			wantResponseBody: `{"message":"ohno","status":400}
`,
		},
		{
			name: "committed response, no error returned, should not appear",
			handler: func(c echo.Context) error {
				c.Response().WriteHeader(http.StatusOK)
				_, _ = c.Response().Write([]byte(`all good`))

				return nil
			},
			wantEmpty:        true,
			wantStatusCode:   http.StatusOK,
			wantResponseBody: "all good",
		},
		{
			name: "committed response, error returned, should appear",
			handler: func(c echo.Context) error {
				c.Response().WriteHeader(http.StatusUnauthorized)
				_, _ = c.Response().Write([]byte(`go away`))

				return errors.Wrap(errors.New("something went wrong"), "oh no")
			},
			wantEmpty: false,
			wantFragments: []string{
				"something went wrong",
				"oh no",
				mockRequestID,
			},
			wantStatusCode:   http.StatusUnauthorized,
			wantResponseBody: "go away",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// buffer is our mock log writer
			b := bytes.NewBuffer(nil)
			l := zerolog.New(b)

			// set up our echo instance with the request ID middleware, the error handler we're testing, and the handler
			// that returns things for the error handler to handle.
			e := echo.New()
			e.Use(mid.UUIDRequestID())
			e.HTTPErrorHandler = Handler(l)
			e.GET("/", tt.handler)

			// set up the http test request and response recorder
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Add(echo.HeaderXRequestID, mockRequestID)
			w := httptest.NewRecorder()

			// do the test
			e.ServeHTTP(w, req)

			// extract log content and response body
			logContent := b.String()

			body, err := io.ReadAll(w.Result().Body)
			require.NoError(t, err, "reading response body")

			// run the assertions on the results
			assert.Equal(t, tt.wantStatusCode, w.Result().StatusCode)
			assert.Equal(t, tt.wantResponseBody, string(body))

			if tt.wantEmpty {
				assert.Empty(t, logContent, "mock writer should have been empty")
				assert.Empty(t, tt.wantFragments, "if log is expected to be empty, want fragments should also be nil")
			} else {
				assert.NotEmpty(t, logContent, "mock writer should not have been empty")
			}

			for _, l := range tt.wantFragments {
				assert.Containsf(t, logContent, l, "log does not contain '%s'", l)
			}
		})
	}
}
