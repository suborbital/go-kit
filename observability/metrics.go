package observability

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc"
)

const (
	defaultReaderTimeout = 5 * time.Second
)

type MeterConfig struct {
	CollectPeriod    time.Duration
	ServiceName      string
	ServiceNamespace string
	ServiceVersion   string
}

// OtelMeter takes a grpc connection to an otel collector, a MeterConfig that holds important data like collection
// period, service and namespace names, and the version of the application we're attaching the meter to. It then
// configures the meter provider, and returns a shutdown function and nil error if successful.
//
// This function merely sets up the scaffolding to ship collected metered data to the opentelemetry collector. It does
// not set up the specific meters for the applications.
func OtelMeter(ctx context.Context, conn *grpc.ClientConn, meterConfig MeterConfig) (func(context.Context) error, error) {
	if meterConfig.CollectPeriod < time.Second {
		return nil, errors.New("collect period is shorter than a second, please choose a longer period to avoid" +
			" overloading the collector")
	}

	// exporter is the thing that will send the data from the app to wherever else it needs to go.
	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithGRPCConn(conn),
	)
	if err != nil {
		return nil, errors.Wrap(err, "otlpmetricgrpc.New")
	}

	// periodicReader is the thing that collects the data every cycle, and then uses the exporter above to send it
	// someplace.
	periodicReader := metric.NewPeriodicReader(exporter,
		metric.WithTimeout(defaultReaderTimeout),
		metric.WithInterval(meterConfig.CollectPeriod),
	)

	// resource configures the very basic attributes of every measurement taken.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(meterConfig.ServiceName),
			semconv.ServiceNamespaceKey.String(meterConfig.ServiceNamespace),
			semconv.ServiceVersionKey.String(meterConfig.ServiceVersion),
			semconv.DBSystemPostgreSQL,
			attribute.String("exporter", "grpc"),
		))

	// meterProvider takes the resource, and the reader, to provide a thing that we can create actual instruments out
	// of, so we can start measuring things.
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(r),
		metric.WithReader(periodicReader),
	)

	otel.SetMeterProvider(meterProvider)
	return meterProvider.Shutdown, nil
}
