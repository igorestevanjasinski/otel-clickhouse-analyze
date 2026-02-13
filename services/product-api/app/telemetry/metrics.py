from opentelemetry import metrics
from opentelemetry.sdk.metrics import MeterProvider
from opentelemetry.sdk.metrics.export import PeriodicExportingMetricReader
from opentelemetry.sdk.resources import Resource, SERVICE_NAME, SERVICE_VERSION, DEPLOYMENT_ENVIRONMENT
from opentelemetry.exporter.otlp.proto.grpc.metric_exporter import OTLPMetricExporter
import logging

logger = logging.getLogger(__name__)


class Metrics:
    """OpenTelemetry metrics for Product API."""
    
    def __init__(self):
        self.meter = None
        self.products_created_counter = None
        self.product_creation_duration = None
        self.active_requests_gauge = None
        self.kafka_publish_counter = None
        self.kafka_publish_errors = None
    
    def setup(self, service_name: str, service_version: str, environment: str, otlp_endpoint: str):
        """Initialize metrics with OTLP exporter."""
        resource = Resource(attributes={
            SERVICE_NAME: service_name,
            SERVICE_VERSION: service_version,
            DEPLOYMENT_ENVIRONMENT: environment,
            "service.namespace": "product-portfolio"
        })
        
        exporter = OTLPMetricExporter(
            endpoint=otlp_endpoint,
            insecure=True
        )
        
        reader = PeriodicExportingMetricReader(exporter, export_interval_millis=5000)
        
        provider = MeterProvider(resource=resource, metric_readers=[reader])
        metrics.set_meter_provider(provider)
        
        self.meter = metrics.get_meter(__name__)
        
        self.products_created_counter = self.meter.create_counter(
            name="products.created.total",
            description="Total number of products created",
            unit="1"
        )
        
        self.product_creation_duration = self.meter.create_histogram(
            name="products.creation.duration",
            description="Duration of product creation in seconds",
            unit="s"
        )
        
        self.active_requests_gauge = self.meter.create_up_down_counter(
            name="http.server.active_requests",
            description="Number of active HTTP requests",
            unit="1"
        )
        
        self.kafka_publish_counter = self.meter.create_counter(
            name="kafka.messages.published.total",
            description="Total number of messages published to Kafka",
            unit="1"
        )
        
        self.kafka_publish_errors = self.meter.create_counter(
            name="kafka.publish.errors.total",
            description="Total number of Kafka publish errors",
            unit="1"
        )
        
        logger.info("OpenTelemetry metrics initialized")


app_metrics = Metrics()
