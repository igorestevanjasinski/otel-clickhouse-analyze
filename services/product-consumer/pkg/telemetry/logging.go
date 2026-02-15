package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoggingConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
}

func SetupLogging(config LoggingConfig) (func(context.Context) error, error) {
	ctx := context.Background()

	// Create gRPC connection to OTLP collector
	conn, err := grpc.DialContext(ctx, config.OTLPEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	// Create OTLP log exporter
	exporter, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	// Create resource with service information
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

	// Create logger provider with batch processor
	processor := log.NewBatchProcessor(exporter)
	loggerProvider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(processor),
	)

	// Set global logger provider
	global.SetLoggerProvider(loggerProvider)

	// Return shutdown function
	return func(ctx context.Context) error {
		if err := loggerProvider.Shutdown(ctx); err != nil {
			return err
		}
		return conn.Close()
	}, nil
}
