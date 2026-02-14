CREATE DATABASE IF NOT EXISTS observability;

CREATE TABLE IF NOT EXISTS observability.traces (
    timestamp DateTime64(9) CODEC(Delta, ZSTD),
    trace_id String CODEC(ZSTD),
    span_id String CODEC(ZSTD),
    parent_span_id String CODEC(ZSTD),
    service_name LowCardinality(String),
    operation_name LowCardinality(String),
    duration_ns UInt64 CODEC(T64, ZSTD),
    status_code LowCardinality(String),
    span_kind LowCardinality(String),
    attributes Map(String, String) CODEC(ZSTD),
    events Array(Tuple(timestamp DateTime64(9), name String, attributes Map(String, String))) CODEC(ZSTD),
    resource_attributes Map(String, String) CODEC(ZSTD),
    INDEX idx_trace_id trace_id TYPE bloom_filter GRANULARITY 1,
    INDEX idx_service service_name TYPE set(100) GRANULARITY 1,
    INDEX idx_operation operation_name TYPE set(100) GRANULARITY 1
) ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (service_name, operation_name, timestamp)
TTL toDateTime(timestamp) + INTERVAL 30 DAY
SETTINGS index_granularity = 8192;

CREATE TABLE IF NOT EXISTS observability.logs (
    timestamp DateTime64(9) CODEC(Delta, ZSTD),
    trace_id String CODEC(ZSTD),
    span_id String CODEC(ZSTD),
    severity LowCardinality(String),
    service_name LowCardinality(String),
    message String CODEC(ZSTD),
    attributes Map(String, String) CODEC(ZSTD),
    resource_attributes Map(String, String) CODEC(ZSTD),
    INDEX idx_trace_id trace_id TYPE bloom_filter GRANULARITY 1,
    INDEX idx_severity severity TYPE set(10) GRANULARITY 1,
    INDEX idx_service service_name TYPE set(100) GRANULARITY 1
) ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (service_name, severity, timestamp)
TTL toDateTime(timestamp) + INTERVAL 30 DAY
SETTINGS index_granularity = 8192;

CREATE MATERIALIZED VIEW IF NOT EXISTS observability.traces_metrics_mv
ENGINE = SummingMergeTree()
PARTITION BY toYYYYMMDD(timestamp_minute)
ORDER BY (service_name, operation_name, status_code, timestamp_minute)
AS SELECT
    toStartOfMinute(timestamp) AS timestamp_minute,
    service_name,
    operation_name,
    status_code,
    count() AS request_count,
    avg(duration_ns) AS avg_duration_ns,
    quantile(0.50)(duration_ns) AS p50_duration_ns,
    quantile(0.95)(duration_ns) AS p95_duration_ns,
    quantile(0.99)(duration_ns) AS p99_duration_ns,
    max(duration_ns) AS max_duration_ns
FROM observability.traces
GROUP BY timestamp_minute, service_name, operation_name, status_code;

CREATE VIEW IF NOT EXISTS observability.error_analysis AS
SELECT
    toStartOfHour(timestamp) AS hour,
    service_name,
    operation_name,
    status_code,
    count() AS error_count,
    avg(duration_ns) / 1000000 AS avg_duration_ms,
    groupArray(5)(trace_id) AS sample_traces
FROM observability.traces
WHERE status_code = 'ERROR'
GROUP BY hour, service_name, operation_name, status_code
ORDER BY hour DESC, error_count DESC;

CREATE VIEW IF NOT EXISTS observability.latency_analysis AS
SELECT
    toStartOfMinute(timestamp) AS minute,
    service_name,
    operation_name,
    count() AS request_count,
    avg(duration_ns) / 1000000 AS avg_ms,
    quantile(0.50)(duration_ns) / 1000000 AS p50_ms,
    quantile(0.95)(duration_ns) / 1000000 AS p95_ms,
    quantile(0.99)(duration_ns) / 1000000 AS p99_ms
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY minute, service_name, operation_name
ORDER BY minute DESC;
