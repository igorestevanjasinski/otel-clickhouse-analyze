import logging
from opentelemetry._logs import set_logger_provider
from opentelemetry.sdk._logs import LoggerProvider, LoggingHandler
from opentelemetry.sdk._logs.export import BatchLogRecordProcessor
from opentelemetry.exporter.otlp.proto.grpc._log_exporter import OTLPLogExporter
from opentelemetry.sdk.resources import Resource, SERVICE_NAME, SERVICE_VERSION, DEPLOYMENT_ENVIRONMENT


def setup_logging(service_name: str, service_version: str, environment: str, otlp_endpoint: str):
    """
    Configure OpenTelemetry logging with OTLP exporter.
    
    Args:
        service_name: Name of the service
        service_version: Version of the service
        environment: Deployment environment (dev, staging, prod)
        otlp_endpoint: OTLP collector endpoint (e.g., http://localhost:4317)
    """
    resource = Resource(attributes={
        SERVICE_NAME: service_name,
        SERVICE_VERSION: service_version,
        DEPLOYMENT_ENVIRONMENT: environment,
        "service.namespace": "product-portfolio"
    })
    
    # Create logger provider
    logger_provider = LoggerProvider(resource=resource)
    
    # Create OTLP log exporter
    otlp_exporter = OTLPLogExporter(
        endpoint=otlp_endpoint,
        insecure=True
    )
    
    # Add batch processor
    logger_provider.add_log_record_processor(
        BatchLogRecordProcessor(otlp_exporter)
    )
    
    # Set global logger provider
    set_logger_provider(logger_provider)
    
    # Create and configure OpenTelemetry logging handler
    handler = LoggingHandler(
        level=logging.NOTSET,
        logger_provider=logger_provider
    )
    
    # Attach handler to root logger
    logging.getLogger().addHandler(handler)
    
    logging.info(
        "OpenTelemetry logging initialized",
        extra={
            "service_name": service_name,
            "otlp_endpoint": otlp_endpoint
        }
    )
