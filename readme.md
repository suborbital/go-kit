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
