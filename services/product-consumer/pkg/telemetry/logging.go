package telemetry

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/bridges/otellogrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func grpcEndpointForLog(endpoint string) string {
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

type LoggingConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
}

func SetupLogging(config LoggingConfig) (func(context.Context) error, error) {
	ctx := context.Background()
	endpoint := grpcEndpointForLog(config.OTLPEndpoint)

	// Create gRPC connection to OTLP collector
	conn, err := grpc.DialContext(ctx, endpoint,
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

	// Set global logger provider and keep reference for AddLogrusHook
	global.SetLoggerProvider(loggerProvider)
	otelLoggerProvider = loggerProvider

	// Return shutdown function
	return func(ctx context.Context) error {
		if err := loggerProvider.Shutdown(ctx); err != nil {
			return err
		}
		return conn.Close()
	}, nil
}

// otelLoggerProvider holds the provider set by SetupLogging so AddLogrusHook can use it.
var otelLoggerProvider *log.LoggerProvider

// AddLogrusHook attaches an OTLP hook to the given logrus logger so that all
// logrus entries are also exported as OTLP log records (e.g. to ClickStack).
// Call after SetupLogging. Returns the hook or nil if SetupLogging was not run.
func AddLogrusHook(logger *logrus.Logger, serviceName string) *otellogrus.Hook {
	if otelLoggerProvider == nil {
		return nil
	}
	hook := otellogrus.NewHook(serviceName, otellogrus.WithLoggerProvider(otelLoggerProvider))
	logger.AddHook(hook)
	return hook
}
