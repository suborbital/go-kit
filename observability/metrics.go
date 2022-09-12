package observability

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"google.golang.org/grpc"
)

// OtelMeter takes a grpc connection to an otel collector, a collectPeriod time.Duration (more than a second), and
// sets the meter collector up and assigns it to a global meter. If the function returns nil, assume that the global
// meter is now set and ready to be used.
//
// This function merely sets up the scaffolding to ship collected metered data to the opentelemetry collector. It does
// not set up the specific meters for the applications.
func OtelMeter(ctx context.Context, conn *grpc.ClientConn, collectPeriod time.Duration) error {
	if collectPeriod < time.Second {
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
		semconv.ServiceNameKey.String("scn-thingy"),
		semconv.TelemetrySDKLanguageGo,
		semconv.ServiceInstanceIDKey.String("4523"),
		semconv.ServiceNamespaceKey.String("production"),
		semconv.ServiceVersionKey.String("v1.2.3"),
		semconv.DBSystemPostgreSQL,
		attribute.String("exporter", "grpc"),
	)

	cont := controller.New(
		processor.NewFactory(
			simple.NewWithInexpensiveDistribution(),
			exporter,
		),
		controller.WithResource(r),
		controller.WithExporter(exporter),
		controller.WithCollectPeriod(collectPeriod),
	)

	if err := cont.Start(context.Background()); err != nil {
		return errors.Wrap(err, "metric controller Start")
	}

	global.SetMeterProvider(cont)
	return nil
}
