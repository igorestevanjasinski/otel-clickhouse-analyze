from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.sdk.resources import Resource, SERVICE_NAME, SERVICE_VERSION, DEPLOYMENT_ENVIRONMENT
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.instrumentation.logging import LoggingInstrumentor
from fastapi import FastAPI
import logging

logger = logging.getLogger(__name__)


def _grpc_endpoint(otlp_endpoint: str) -> str:
    """Normalize OTLP endpoint to host:port for gRPC (strip scheme and path)."""
    endpoint = (otlp_endpoint or "").strip().replace("http://", "").replace("https://", "").rstrip("/")
    if "/" in endpoint:
        endpoint = endpoint.split("/")[0]
    return endpoint or "localhost:4317"


def setup_tracing(app: FastAPI, service_name: str, service_version: str, environment: str, otlp_endpoint: str):
    """
    Configure OpenTelemetry tracing with OTLP exporter.
    
    Args:
        app: FastAPI application instance
        service_name: Name of the service
        service_version: Version of the service
        environment: Deployment environment (dev, staging, prod)
        otlp_endpoint: OTLP collector endpoint (e.g., http://localhost:4317 or clickstack:4317)
    """
    resource = Resource(attributes={
        SERVICE_NAME: service_name,
        SERVICE_VERSION: service_version,
        DEPLOYMENT_ENVIRONMENT: environment,
        "service.namespace": "product-portfolio"
    })
    
    provider = TracerProvider(resource=resource)
    
    otlp_exporter = OTLPSpanExporter(
        endpoint=_grpc_endpoint(otlp_endpoint),
        insecure=True
    )
    
    processor = BatchSpanProcessor(otlp_exporter)
    provider.add_span_processor(processor)
    
    trace.set_tracer_provider(provider)
    
    FastAPIInstrumentor.instrument_app(app)
    
    LoggingInstrumentor().instrument(set_logging_format=True)
    
    logger.info("OpenTelemetry tracing initialized", extra={
        "service_name": service_name,
        "otlp_endpoint": otlp_endpoint
    })


def get_tracer(name: str) -> trace.Tracer:
    """Get a tracer instance."""
    return trace.get_tracer(name)


def get_current_trace_id() -> str:
    """Get current trace ID as hex string."""
    span = trace.get_current_span()
    if span and span.get_span_context().is_valid:
        return format(span.get_span_context().trace_id, '032x')
    return ""


def get_current_span_id() -> str:
    """Get current span ID as hex string."""
    span = trace.get_current_span()
    if span and span.get_span_context().is_valid:
        return format(span.get_span_context().span_id, '016x')
    return ""
