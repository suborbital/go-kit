package observability

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"google.golang.org/grpc"
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
func OtelMeter(ctx context.Context, conn *grpc.ClientConn, meterConfig MeterConfig) error {
	if meterConfig.CollectPeriod < time.Second {
		return errors.New("collect period is shorter than a second, please choose a longer period to avoid" +
			" overloading the collector")
	}

	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithGRPCConn(conn),
	)
	if err != nil {
		return errors.Wrap(err, "otlpmetricgrpc.New")
	}

	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(meterConfig.ServiceName),
		semconv.TelemetrySDKLanguageGo,
		semconv.ServiceNamespaceKey.String(meterConfig.ServiceNamespace),
		semconv.ServiceVersionKey.String(meterConfig.ServiceVersion),
		semconv.DBSystemPostgreSQL,
		attribute.String("exporter", "grpc"),
	)

	reader := metric.NewPeriodicReader(exporter, metric.WithInterval(meterConfig.CollectPeriod))
	provider := metric.NewMeterProvider(metric.WithResource(r), metric.WithReader(reader))

	global.SetMeterProvider(provider)
	return nil
}
