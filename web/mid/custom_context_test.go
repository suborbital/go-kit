package mid_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	kitHttp "github.com/suborbital/go-kit/web/http"
	"github.com/suborbital/go-kit/web/mid"
)

func TestCustomContext(t *testing.T) {
	const (
		expectedRid = "hello-from-the-requestid"
	)

	uuidMW := middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return expectedRid
		},
	})

	tests := []struct {
		name               string
		middlewares        []echo.MiddlewareFunc
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "no request id mw, only custom context",
			middlewares: []echo.MiddlewareFunc{
				mid.CustomContext(),
			},
			expectedBody:       "unknown-request",
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "request id before custom context",
			middlewares: []echo.MiddlewareFunc{
				uuidMW,
				mid.CustomContext(),
			},
			expectedBody:       expectedRid,
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "request id after custom context",
			middlewares: []echo.MiddlewareFunc{
				mid.CustomContext(),
				uuidMW,
			},
			expectedBody:       expectedRid,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "no custom context mw",
			middlewares: nil,
			expectedBody: `{"message":"` + http.StatusText(http.StatusInternalServerError) + `"}
`,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Use(tt.middlewares...)
			e.GET("/", func(c echo.Context) error {
				cc, ok := c.(*kitHttp.Context)
				if !ok {
					return echo.NewHTTPError(http.StatusInternalServerError)
				}

				rid := cc.RequestID()

				return c.String(http.StatusOK, rid)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			e.ServeHTTP(w, req)
			b, err := io.ReadAll(w.Result().Body)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatusCode, w.Result().StatusCode)
			assert.Equalf(t, tt.expectedBody, string(b), "got body\n%s\n", b)
		})
	}
}
