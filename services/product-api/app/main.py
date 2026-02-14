from contextlib import asynccontextmanager
import logging
import sys
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pythonjsonlogger import jsonlogger

from app.config import get_settings
from app.api.routes import router as products_router, set_kafka_producer
from app.services import KafkaProducerService
from app.middleware.error_injection import ErrorInjectionMiddleware
from app.middleware.latency_injection import LatencyInjectionMiddleware
from app.telemetry.tracing import setup_tracing, get_current_trace_id
from app.telemetry.metrics import app_metrics
from app.telemetry.prometheus_metrics import metrics_app


def setup_logging() -> None:
    settings = get_settings()
    
    class TraceIdFilter(logging.Filter):
        def filter(self, record):
            record.trace_id = get_current_trace_id()
            return True
    
    log_handler = logging.StreamHandler(sys.stdout)
    formatter = jsonlogger.JsonFormatter(
        fmt="%(asctime)s %(name)s %(levelname)s %(message)s %(trace_id)s",
        rename_fields={"asctime": "timestamp", "levelname": "level", "name": "logger"}
    )
    log_handler.setFormatter(formatter)
    log_handler.addFilter(TraceIdFilter())
    
    root_logger = logging.getLogger()
    root_logger.addHandler(log_handler)
    root_logger.setLevel(getattr(logging, settings.log_level.upper()))


kafka_producer_service: KafkaProducerService | None = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    global kafka_producer_service
    settings = get_settings()
    logger = logging.getLogger(__name__)
    
    logger.info(
        "Application starting",
        extra={
            "app_name": settings.app_name,
            "version": settings.app_version,
            "environment": settings.environment
        }
    )
    
    kafka_producer_service = KafkaProducerService(
        bootstrap_servers=settings.kafka_bootstrap_servers,
        topic=settings.kafka_topic_products,
        retries=settings.kafka_retries,
        retry_backoff_ms=settings.kafka_retry_backoff_ms,
        acks=settings.kafka_acks
    )
    set_kafka_producer(kafka_producer_service)
    
    logger.info("Kafka producer initialized")
    
    yield
    
    if kafka_producer_service:
        kafka_producer_service.close()
    
    logger.info("Application shutting down")


def create_application() -> FastAPI:
    setup_logging()
    settings = get_settings()
    
    app = FastAPI(
        title=settings.app_name,
        version=settings.app_version,
        description="Product API for microservices observability project",
        docs_url="/docs",
        redoc_url="/redoc",
        lifespan=lifespan
    )
    
    setup_tracing(
        app=app,
        service_name=settings.app_name,
        service_version=settings.app_version,
        environment=settings.environment,
        otlp_endpoint=settings.otlp_endpoint
    )
    
    app_metrics.setup(
        service_name=settings.app_name,
        service_version=settings.app_version,
        environment=settings.environment,
        otlp_endpoint=settings.otlp_endpoint
    )
    
    app.add_middleware(
        CORSMiddleware,
        allow_origins=settings.cors_origins,
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )
    
    if settings.enable_chaos:
        logger = logging.getLogger(__name__)
        logger.warning(
            "Chaos Engineering enabled",
            extra={
                "error_rate": settings.error_rate,
                "latency_min_ms": settings.latency_ms_min,
                "latency_max_ms": settings.latency_ms_max
            }
        )
        app.add_middleware(
            ErrorInjectionMiddleware,
            error_rate=settings.error_rate
        )
        app.add_middleware(
            LatencyInjectionMiddleware,
            min_ms=settings.latency_ms_min,
            max_ms=settings.latency_ms_max
        )
    
    app.include_router(products_router, tags=["products"])
    
    app.mount("/metrics", metrics_app)
    
    return app


app = create_application()
