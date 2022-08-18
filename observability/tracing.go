package observability

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"google.golang.org/grpc"
)

const (
	collectorExporterName = "collector"
	honeycombExporterName = "honeycomb"
)

// TracingConfig is the minimum configuration needed to configure any of the tracing solutions that aren't the no-op
// tracer.
type TracingConfig struct {
	Probability float64
	ServiceName string
}

// HoneycombTracingConfig embeds the TracingConfig struct, and adds other, specifically Honeycomb related fields.
type HoneycombTracingConfig struct {
	TracingConfig
	APIKey  string
	Dataset string
}

// OtelTracer sets up a trace provider that sends data to an opentelemetry collector.
func OtelTracer(ctx context.Context, conn *grpc.ClientConn, config TracingConfig) (*trace.TracerProvider, error) {
	exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithGRPCConn(conn),
	))
	if err != nil {
		return nil, errors.Wrap(err, "oltptrace.New with exporter as collector")
	}

	traceOpts := tracerOpts(exporter, config.ServiceName, collectorExporterName, config.Probability)
	traceProvider := trace.NewTracerProvider(traceOpts...)
	otel.SetTracerProvider(traceProvider)

	return traceProvider, nil
}

// HoneycombTracer returns a tracer provider configured to send traces to your Honeycomb account.
func HoneycombTracer(ctx context.Context, conn *grpc.ClientConn, config HoneycombTracingConfig) (*trace.TracerProvider, error) {
	exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithGRPCConn(conn),
		otlptracegrpc.WithHeaders(map[string]string{
			"x-honeycomb-team":    config.APIKey,
			"x-honeycomb-dataset": config.Dataset,
		}),
	))
	if err != nil {
		return nil, errors.Wrap(err, "oltptrace.New with exporter as honeycomb")
	}

	traceOpts := tracerOpts(exporter, config.ServiceName, honeycombExporterName, config.Probability)
	traceProvider := trace.NewTracerProvider(traceOpts...)

	otel.SetTracerProvider(traceProvider)

	return traceProvider, nil
}

// NoopTracer returns a non-configured empty trace provider that won't do anything.
func NoopTracer() (*trace.TracerProvider, error) {
	// Create the most default trace provider and escape early.
	traceProvider := trace.NewTracerProvider()
	otel.SetTracerProvider(traceProvider)

	return traceProvider, nil
}

// tracerOpts is a utility function to cut down on code duplication, as the tracer provider options overlap between the
// collector and honeycomb tracer implementations.
func tracerOpts(exporter trace.SpanExporter, serviceName, exporterName string, probability float64) []trace.TracerProviderOption {
	return []trace.TracerProviderOption{
		trace.WithSampler(trace.TraceIDRatioBased(probability)),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
				attribute.String("exporter", exporterName),
			),
		),
		trace.WithBatcher(exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
	}
}
