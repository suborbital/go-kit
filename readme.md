# Suborbital go-kit

go-kit package is a reusable set of modules that all services share. Most set up codes that all would need are delegated here to reduce code duplication and make maintenance easier to follow.

## Metrics

Metrics returns a configured MeterProvider with a shutdown function. Here's how to use it from your service's `main` function:

```go
package main

func main() {
	grpcConn, err := metrics.GrpcConnection(ctx, endpoint)
	if err != nil {
		log.Fatal("Failed to get grpc connection")
	}
	
	mc := metrics.MeterConfig{
		CollectPeriod:     5 * time.Second,
		ServiceName:      "my-service",
		ServiceNamespace: "production",
		ServiceVersion:   "v0.0.1",
    }
	
	shutdownFunc, err := metrics.OtelMeter(ctx, grpcConn, mc)
	if err != nil {
		log.Fatal("failed to get otel meter provider")
    }
	
	defer shutdownFunc()
	
	// do the thing that blocks here, like accept incoming connections, etc
}
```

This will set up the meter provider and put it in a metrics global context. Then to actually create the meters, you will need to grab the meter provider from the global context, and spawn the meters from then:

```go
package someplace

import (
	"go.opentelemetry.io/otel/metric/global"
)

func Meters() {
	mp := global.Meter("instrumentationName")
}
```

## Tracing

## Web

There are three configurable echo middlewares included in the kit.

### CORS

Provides good enough defaults with a simple call signature for ease of use:
```go
func main() {
    e := echo.New()
	e.Use(
		mid.CORS("*"),
    )
}
```
In case it's needed, you can configure additional domains, additional allowed headers, and a skipper function in case there's a route you don't want the middleware to be applied to.

```go
func main() {
	e := echo.New()
	e.Use(
		mid.CORS(
			"domainone.com",
			mid.WithDomains("domaintwo.org", "domainthree.exe"),
			mid.WithHeaders("X-Suborbital-State"),
			mid.WithSkipper(func(c echo.Context) bool {
				return c.Path() != "/dont/cors/this"
			}),
        ),
    )
}
```

### Logger
Provides a middleware that will log when a request comes in and when the same request goes out. Error handling happens before the response is logged, which means neither the logger, nor any other middleware up the chain can further modify the response status code / body.

Important! The requestID logger needs to be outside of this middleware. In practical terms, it needs to be in the list passed to `e.Use` earlier.

It uses rs/zerolog.
```go
func main() {
	logger := zerolog.New(os.Stderr).With().Str("service", "myservice").Logger()
	
	e := echo.New()
	e.Use(
		// requestID middleware goes here somewhere. As long as it's above the logger.
		mid.Logger(logger),
    )
}
```

### RequestID

The `UUIDRequestID` middleware configures echo's built in request ID middleware to use UUIDv4s instead of a random twenty-something character string.

This middleware needs to wrap the logger middleware. In practical terms, the logger needs to come after this one when passed to the `e.Use` method.

It also saves the request ID in the echo context itself, as well as the standard request context. All three should be present if the middleware is used.

The key to retrieve the request ID value is `mid.RequestIDKey`.

```go
func main() {
	e := echo.New()
	
	e.Use(
		mid.UUIDRequestID(),
		// logger middleware should go somewhere here
    )
}

func Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		// to get the request ID from the echo context
		rid := c.Get(mid.RequestIDKey)
		
		// to get it from the request context
		rid, ok := c.Request().Context().Value(mid.RequestIDKey).(string)
		
		// to get it from the request header
		rid := c.Request().Header.Get(echo.HeaderXRequestID)
	}
}
```

### Tracing

OpenTelemetry contrib already has an echo tracing middleware, best to use that one. You still need to configure it beforehand.

The example that's in their repository is a minimally working implementation that's around 60 lines of code including the main function: https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/labstack/echo/otelecho/example/server.go#L46.
