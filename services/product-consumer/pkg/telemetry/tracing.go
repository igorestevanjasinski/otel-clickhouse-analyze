package telemetry

import (
	"context"
	"log"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// tracer: inicia com no-op; após SetupTracing com sucesso passa a exportar para OTLP
var tracer trace.Tracer = otel.Tracer("product-consumer")

// grpcEndpoint normaliza o endpoint para host:port (gRPC não usa scheme).
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

type TracingConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
}

func SetupTracing(config TracingConfig) (func(context.Context) error, error) {
	ctx := context.Background()
	endpoint := grpcEndpoint(config.OTLPEndpoint)

	conn, err := grpc.DialContext(ctx, endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
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

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	tracer = otel.Tracer(config.ServiceName)

	log.Printf("OpenTelemetry tracing initialized: service=%s endpoint=%s",
		config.ServiceName, endpoint)

	return tracerProvider.Shutdown, nil
}

func GetTracer() trace.Tracer {
	return tracer
}
