from prometheus_client import Counter, Histogram, Gauge, CollectorRegistry, make_asgi_app

prometheus_registry = CollectorRegistry()

REQUEST_COUNT = Counter(
    'products_created_total',
    'Total products created',
    ['status', 'endpoint'],
    registry=prometheus_registry
)

REQUEST_LATENCY = Histogram(
    'product_creation_duration_seconds',
    'Product creation latency in seconds',
    ['endpoint'],
    buckets=[0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1.0, 2.5, 5.0, 7.5, 10.0],
    registry=prometheus_registry
)

ACTIVE_REQUESTS = Gauge(
    'active_requests',
    'Active requests',
    ['endpoint'],
    registry=prometheus_registry
)

KAFKA_MESSAGES_PUBLISHED = Counter(
    'kafka_messages_published_total',
    'Total Kafka messages published',
    ['topic', 'status'],
    registry=prometheus_registry
)

metrics_app = make_asgi_app(registry=prometheus_registry)
