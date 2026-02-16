"""
OpenTelemetry metrics for Product API (OTLP export).
Replaces Prometheus metrics for ClickStack compatibility.
"""
import logging
from opentelemetry import metrics
from opentelemetry.sdk.metrics import MeterProvider
from opentelemetry.sdk.metrics.export import PeriodicExportingMetricReader
from opentelemetry.sdk.resources import (
    Resource,
    SERVICE_NAME,
    SERVICE_VERSION,
    DEPLOYMENT_ENVIRONMENT,
)
from opentelemetry.exporter.otlp.proto.grpc.metric_exporter import OTLPMetricExporter

logger = logging.getLogger(__name__)

# Set in setup_metrics(); instruments created once there
_request_count_counter = None
_request_latency_histogram = None
_active_requests_counter = None
_kafka_messages_published_counter = None


def setup_metrics(
    service_name: str,
    service_version: str,
    environment: str,
    otlp_endpoint: str,
    export_interval_millis: int = 10000,
) -> None:
    """
    Configure OpenTelemetry metrics with OTLP exporter (same endpoint as traces).
    """
    global _request_count_counter, _request_latency_histogram
    global _active_requests_counter, _kafka_messages_published_counter

    # gRPC endpoint: host:port without scheme
    endpoint = otlp_endpoint.replace("http://", "").replace("https://", "").rstrip("/")
    if "/" in endpoint:
        endpoint = endpoint.split("/")[0]

    resource = Resource(
        attributes={
            SERVICE_NAME: service_name,
            SERVICE_VERSION: service_version,
            DEPLOYMENT_ENVIRONMENT: environment,
            "service.namespace": "product-portfolio",
        }
    )
    exporter = OTLPMetricExporter(endpoint=endpoint, insecure=True)
    reader = PeriodicExportingMetricReader(
        exporter,
        export_interval_millis=export_interval_millis,
    )
    provider = MeterProvider(resource=resource, metric_readers=[reader])
    metrics.set_meter_provider(provider)
    meter = metrics.get_meter(service_name, service_version)

    _request_count_counter = meter.create_counter(
        "products_created_total",
        description="Total products created",
        unit="1",
    )
    _request_latency_histogram = meter.create_histogram(
        "product_creation_duration_seconds",
        description="Product creation latency in seconds",
        unit="s",
    )
    _active_requests_counter = meter.create_up_down_counter(
        "active_requests",
        description="Active requests",
        unit="1",
    )
    _kafka_messages_published_counter = meter.create_counter(
        "kafka_messages_published_total",
        description="Total Kafka messages published",
        unit="1",
    )

    logger.info(
        "OpenTelemetry metrics initialized",
        extra={"service_name": service_name, "otlp_endpoint": endpoint},
    )


# Public API: call these from routes (attributes match former Prometheus labels)
def record_request_count(status: str, endpoint: str, value: int = 1) -> None:
    try:
        if _request_count_counter is not None:
            _request_count_counter.add(value, {"status": status, "endpoint": endpoint})
    except Exception as e:
        logger.debug("Metrics record_request_count skipped: %s", e)


def record_request_latency(endpoint: str, duration_seconds: float) -> None:
    try:
        if _request_latency_histogram is not None:
            _request_latency_histogram.record(duration_seconds, {"endpoint": endpoint})
    except Exception as e:
        logger.debug("Metrics record_request_latency skipped: %s", e)


def active_requests_inc(endpoint: str, delta: int = 1) -> None:
    try:
        if _active_requests_counter is not None:
            _active_requests_counter.add(delta, {"endpoint": endpoint})
    except Exception as e:
        logger.debug("Metrics active_requests_inc skipped: %s", e)


def active_requests_dec(endpoint: str, delta: int = 1) -> None:
    try:
        if _active_requests_counter is not None:
            _active_requests_counter.add(-delta, {"endpoint": endpoint})
    except Exception as e:
        logger.debug("Metrics active_requests_dec skipped: %s", e)


def record_kafka_messages_published(topic: str, status: str, value: int = 1) -> None:
    try:
        if _kafka_messages_published_counter is not None:
            _kafka_messages_published_counter.add(
                value, {"topic": topic, "status": status}
            )
    except Exception as e:
        logger.debug("Metrics record_kafka_messages_published skipped: %s", e)
