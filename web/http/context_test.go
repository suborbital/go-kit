package http_test

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

func TestRID(t *testing.T) {
	const expectedRequestID = "rid-testing-in-progress"

	requestMiddleware := middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return expectedRequestID
		},
	})

	tests := []struct {
		name        string
		middlewares []echo.MiddlewareFunc
		want        string
	}{
		{
			name:        "with request ID middleware",
			middlewares: []echo.MiddlewareFunc{requestMiddleware},
			want:        expectedRequestID,
		},
		{
			name:        "without request ID middleware",
			middlewares: nil,
			want:        "unknown-request",
		},
		{
			name: "with request ID and custom context middlewares",
			middlewares: []echo.MiddlewareFunc{
				requestMiddleware,
				mid.CustomContext(),
			},
			want: expectedRequestID,
		},
		{
			name: "no request ID middleware, yes custom context middleware",
			middlewares: []echo.MiddlewareFunc{
				mid.CustomContext(),
			},
			want: "unknown-request",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Use(tt.middlewares...)
			e.GET("/", func(c echo.Context) error {
				return c.String(http.StatusOK, kitHttp.RID(c))
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			e.ServeHTTP(w, req)

			b, err := io.ReadAll(w.Result().Body)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, w.Result().StatusCode)
			assert.Equal(t, tt.want, string(b))
		})
	}
}
