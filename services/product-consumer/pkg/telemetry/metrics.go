package telemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Metrics struct {
	MessagesConsumed   metric.Int64Counter
	ProcessingDuration metric.Float64Histogram
	MessagesInFlight   metric.Int64UpDownCounter
	ProcessingErrors   metric.Int64Counter
	DatabaseOperations metric.Int64Counter
}

var AppMetrics *Metrics

type MetricsConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
}

func SetupMetrics(config MetricsConfig) (func(context.Context) error, error) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, config.OTLPEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(config.ServiceName),
			semconv.ServiceVersion(config.ServiceVersion),
			semconv.DeploymentEnvironment(config.Environment),
			attribute.String("service.namespace", "product-portfolio"),
		),
	)
	if err != nil {
		return nil, err
	}

	reader := sdkmetric.NewPeriodicReader(exporter)
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(reader),
	)

	otel.SetMeterProvider(meterProvider)

	meter := meterProvider.Meter(config.ServiceName)

	AppMetrics = &Metrics{}

	AppMetrics.MessagesConsumed, err = meter.Int64Counter(
		"messages.consumed.total",
		metric.WithDescription("Total number of messages consumed from Kafka"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	AppMetrics.ProcessingDuration, err = meter.Float64Histogram(
		"messages.processing.duration",
		metric.WithDescription("Duration of message processing in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	AppMetrics.MessagesInFlight, err = meter.Int64UpDownCounter(
		"messages.in_flight",
		metric.WithDescription("Number of messages currently being processed"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	AppMetrics.ProcessingErrors, err = meter.Int64Counter(
		"messages.processing.errors.total",
		metric.WithDescription("Total number of message processing errors"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	AppMetrics.DatabaseOperations, err = meter.Int64Counter(
		"database.operations.total",
		metric.WithDescription("Total number of database operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	log.Printf("OpenTelemetry metrics initialized")

	return meterProvider.Shutdown, nil
}
