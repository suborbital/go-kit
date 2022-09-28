package observability

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/view"
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

// OtelMeter takes a grpc connection to an otel collector, a collectPeriod time.Duration (more than a second), and
// sets the meter collector up and assigns it to a global meter. If the function returns nil, assume that the global
// meter is now set and ready to be used.
//
// This function merely sets up the scaffolding to ship collected metered data to the opentelemetry collector. It does
// not set up the specific meters for the applications.
func OtelMeter(ctx context.Context, conn *grpc.ClientConn, meterConfig MeterConfig) (error, func(context.Context) error) {
	if meterConfig.CollectPeriod < time.Second {
		return errors.New("collect period is shorter than a second, please choose a longer period to avoid" +
			" overloading the collector"), nil
	}

	// exporter is the thing that will send the data from the app to wherever else it needs to go.
	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithGRPCConn(conn),
	)
	if err != nil {
		return errors.Wrap(err, "otlpmetricgrpc.New"), nil
	}

	// periodicReader is the thing that collects the data every cycle, and then uses the exporter above to send it
	// someplace.
	periodicReader := metric.NewPeriodicReader(exporter,
		metric.WithTimeout(defaultReaderTimeout),
		metric.WithInterval(meterConfig.CollectPeriod),
		metric.WithTemporalitySelector(func(view.InstrumentKind) metricdata.Temporality {
			return metricdata.DeltaTemporality
		}),
	)

	// resource configures the very basic attributes of every measurement taken.
	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(meterConfig.ServiceName),
		semconv.TelemetrySDKLanguageGo,
		semconv.ServiceNamespaceKey.String(meterConfig.ServiceNamespace),
		semconv.ServiceVersionKey.String(meterConfig.ServiceVersion),
		semconv.DBSystemPostgreSQL,
		attribute.String("exporter", "grpc"),
	)

	// meterProvider takes the resource, and the reader, to provide a thing that we can create actual instruments out
	// of, so we can start measuring things.
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(r),
		metric.WithReader(periodicReader),
	)

	global.SetMeterProvider(meterProvider)
	return nil, meterProvider.Shutdown
}
