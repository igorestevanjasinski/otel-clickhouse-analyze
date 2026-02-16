package metrics

import (
	"context"
	"log"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	otlpmetricgrpc "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func grpcEndpoint(endpoint string) string {
	s := strings.TrimSpace(endpoint)
	s = strings.TrimPrefix(s, "http://")
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimSuffix(s, "/")
	if idx := strings.Index(s, "/"); idx > 0 {
		s = s[:idx]
	}
	if s == "" {
		return "localhost:4317"
	}
	return s
}

var (
	messagesConsumedCounter   metric.Int64Counter
	processingDurationHist    metric.Float64Histogram
	inFlightUpDown           metric.Int64UpDownCounter
	databaseOperationsCounter metric.Int64Counter
	databaseOpDurationHist    metric.Float64Histogram
)

type MetricsConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
}

// SetupMetrics initializes OpenTelemetry metrics and exports them via OTLP.
// Returns a shutdown function to be called on application exit.
func SetupMetrics(ctx context.Context, config MetricsConfig) (func(context.Context) error, error) {
	endpoint := grpcEndpoint(config.OTLPEndpoint)
	conn, err := grpc.DialContext(ctx, endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		conn.Close()
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
		conn.Close()
		return nil, err
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
	)

	meter := provider.Meter(config.ServiceName)

	messagesConsumedCounter, err = meter.Int64Counter(
		"messages_consumed_total",
		metric.WithDescription("Total messages consumed from Kafka"),
	)
	if err != nil {
		provider.Shutdown(ctx)
		return nil, err
	}

	processingDurationHist, err = meter.Float64Histogram(
		"message_processing_duration_seconds",
		metric.WithDescription("Message processing duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		provider.Shutdown(ctx)
		return nil, err
	}

	inFlightUpDown, err = meter.Int64UpDownCounter(
		"messages_in_flight",
		metric.WithDescription("Messages currently being processed"),
	)
	if err != nil {
		provider.Shutdown(ctx)
		return nil, err
	}

	databaseOperationsCounter, err = meter.Int64Counter(
		"database_operations_total",
		metric.WithDescription("Total database operations"),
	)
	if err != nil {
		provider.Shutdown(ctx)
		return nil, err
	}

	databaseOpDurationHist, err = meter.Float64Histogram(
		"database_operation_duration_seconds",
		metric.WithDescription("Database operation duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		provider.Shutdown(ctx)
		return nil, err
	}

	log.Printf("OpenTelemetry metrics initialized: service=%s endpoint=%s",
		config.ServiceName, endpoint)

	return provider.Shutdown, nil
}

// RecordMessagesConsumed records a consumed message (topic, status).
func RecordMessagesConsumed(topic, status string, value int64) {
	if messagesConsumedCounter != nil {
		messagesConsumedCounter.Add(context.Background(), value,
			metric.WithAttributes(
				attribute.String("topic", topic),
				attribute.String("status", status),
			),
		)
	}
}

// RecordProcessingDuration records message processing duration in seconds.
func RecordProcessingDuration(topic string, durationSeconds float64) {
	if processingDurationHist != nil {
		processingDurationHist.Record(context.Background(), durationSeconds,
			metric.WithAttributes(attribute.String("topic", topic)),
		)
	}
}

// InFlightAdd adds delta to the in-flight messages gauge (e.g. +1 on start, -1 on end).
func InFlightAdd(delta int64) {
	if inFlightUpDown != nil {
		inFlightUpDown.Add(context.Background(), delta)
	}
}

// RecordDatabaseOperations records a database operation (operation, status).
func RecordDatabaseOperations(operation, status string, value int64) {
	if databaseOperationsCounter != nil {
		databaseOperationsCounter.Add(context.Background(), value,
			metric.WithAttributes(
				attribute.String("operation", operation),
				attribute.String("status", status),
			),
		)
	}
}

// RecordDatabaseOperationDuration records database operation duration in seconds.
func RecordDatabaseOperationDuration(operation string, durationSeconds float64) {
	if databaseOpDurationHist != nil {
		databaseOpDurationHist.Record(context.Background(), durationSeconds,
			metric.WithAttributes(attribute.String("operation", operation)),
		)
	}
}
