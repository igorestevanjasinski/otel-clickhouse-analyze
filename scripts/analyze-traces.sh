#!/bin/bash

set -e

CLICKHOUSE_CLIENT="docker compose exec clickhouse clickhouse-client"

echo "========================================="
echo "ClickHouse Trace Analysis Dashboard"
echo "========================================="
echo ""

echo "=== Database Info ==="
$CLICKHOUSE_CLIENT --query="
SELECT 
    database,
    formatReadableSize(sum(bytes)) AS size,
    sum(rows) AS rows
FROM system.parts
WHERE database = 'observability'
GROUP BY database
FORMAT PrettyCompact
"
echo ""

echo "=== Top 10 Slowest Operations ==="
$CLICKHOUSE_CLIENT --query="
SELECT 
    service_name,
    operation_name,
    count() AS count,
    round(avg(duration_ns) / 1000000, 2) AS avg_ms,
    round(quantile(0.99)(duration_ns) / 1000000, 2) AS p99_ms
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY service_name, operation_name
ORDER BY p99_ms DESC
LIMIT 10
FORMAT PrettyCompact
"
echo ""

echo "=== Error Rate by Service ==="
$CLICKHOUSE_CLIENT --query="
SELECT 
    service_name,
    countIf(status_code = 'ERROR') AS errors,
    count() AS total,
    round((errors / total) * 100, 2) AS error_rate_pct
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY service_name
ORDER BY error_rate_pct DESC
FORMAT PrettyCompact
"
echo ""

echo "=== RED Metrics (Rate, Errors, Duration) ==="
$CLICKHOUSE_CLIENT --query="
SELECT 
    service_name,
    count() AS request_rate,
    countIf(status_code = 'ERROR') AS error_count,
    round((error_count / request_rate) * 100, 2) AS error_rate_pct,
    round(avg(duration_ns) / 1000000, 2) AS avg_duration_ms,
    round(quantile(0.50)(duration_ns) / 1000000, 2) AS p50_ms,
    round(quantile(0.95)(duration_ns) / 1000000, 2) AS p95_ms,
    round(quantile(0.99)(duration_ns) / 1000000, 2) AS p99_ms
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY service_name
ORDER BY request_rate DESC
FORMAT PrettyCompact
"
echo ""

echo "=== Throughput by Minute ==="
$CLICKHOUSE_CLIENT --query="
SELECT 
    toStartOfMinute(timestamp) AS minute,
    service_name,
    count() AS requests_per_minute
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 15 MINUTE
GROUP BY minute, service_name
ORDER BY minute DESC
LIMIT 20
FORMAT PrettyCompact
"
echo ""

echo "=== Most Frequent Operations ==="
$CLICKHOUSE_CLIENT --query="
SELECT 
    service_name,
    operation_name,
    count() AS execution_count,
    round(avg(duration_ns) / 1000000, 2) AS avg_duration_ms
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY service_name, operation_name
ORDER BY execution_count DESC
LIMIT 10
FORMAT PrettyCompact
"
echo ""

echo "=== Log Severity Distribution ==="
$CLICKHOUSE_CLIENT --query="
SELECT 
    severity,
    service_name,
    count() AS log_count
FROM observability.logs
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY severity, service_name
ORDER BY log_count DESC
FORMAT PrettyCompact
"
echo ""

echo "=== Distributed Traces (Multi-Service) ==="
$CLICKHOUSE_CLIENT --query="
SELECT 
    trace_id,
    countDistinct(service_name) AS services_involved,
    groupArray(DISTINCT service_name) AS services,
    count() AS total_spans,
    round(dateDiff('millisecond', min(timestamp), max(timestamp)), 2) AS trace_duration_ms
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY trace_id
HAVING services_involved > 1
ORDER BY trace_duration_ms DESC
LIMIT 10
FORMAT PrettyCompact
"
echo ""

echo "=== Data Volume by Service ==="
$CLICKHOUSE_CLIENT --query="
SELECT 
    service_name,
    count() AS span_count,
    uniq(trace_id) AS unique_traces,
    formatReadableSize(sum(length(toString(attributes)))) AS attributes_size
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY service_name
ORDER BY span_count DESC
FORMAT PrettyCompact
"
echo ""

echo "========================================="
echo "Analysis completed successfully!"
echo "========================================="
