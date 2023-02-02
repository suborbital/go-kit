package mid

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CORSOptions represents configuration options for the CORS middleware to further customize the built-in echo CORS
// middleware to suit our needs.
type CORSOptions struct {
	domains           []string
	additionalHeaders []string
	skipper           func(echo.Context) bool
}

// OptionModifier is a type of function that changes values on a CORSOptions struct in place.
type OptionModifier func(c *CORSOptions)

// CORS configures echo's CORS middleware with the following default allowed headers for the specified domain:
// Origin, Content-Type, Accept, Accept-Encoding, Content-Length, Authorization, Cache-Control.
//
// You can use the two modifier functions to add additional domains, or additional headers, for example:
//   - WithDomains("example.net", "jquery.com")
//   - WithHeaders("X-Suborbital-Something", "X-Marks-The-Spot")
//   - WithSkipper(func(c echo.Context) bool { return c.Path() != "/home" })
func CORS(domain string, options ...OptionModifier) echo.MiddlewareFunc {
	corsOptions := CORSOptions{
		domains:           []string{domain},
		additionalHeaders: make([]string, 0),
		skipper:           middleware.DefaultSkipper,
	}

	for _, o := range options {
		o(&corsOptions)
	}

	corsHeaders := []string{
		echo.HeaderOrigin,
		echo.HeaderContentType,
		echo.HeaderAccept,
		echo.HeaderContentLength,
		echo.HeaderAuthorization,
		echo.HeaderAcceptEncoding,
		echo.HeaderCacheControl,
	}

	corsHeaders = append(corsHeaders, corsOptions.domains...)

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: corsOptions.domains,
		AllowHeaders: corsHeaders,
		Skipper:      corsOptions.skipper,
	})
}

// WithDomains adds additional domains to the CORS middleware as permitted origins besides the one already passed to the
// middleware constructor.
func WithDomains(domains ...string) OptionModifier {
	return func(c *CORSOptions) {
		c.domains = append(c.domains, domains...)
	}
}

// WithHeaders adds additional header keys that are returned in cross-origin responses besides the ones the CORS
// middleware already allows.
func WithHeaders(headers ...string) OptionModifier {
	return func(c *CORSOptions) {
		c.additionalHeaders = append(c.additionalHeaders, headers...)
	}
}

// WithSkipper configures a skipper function for the CORS header. If not set, it will use middleware.DefaultSkipper,
// which enables the middleware to be used on all routes it's attached to.
func WithSkipper(skipper func(echo.Context) bool) OptionModifier {
	return func(c *CORSOptions) {
		c.skipper = skipper
	}
}
