-- ========================================
-- QUERIES BÁSICAS
-- ========================================

-- Listar todos os trace_ids recentes
SELECT DISTINCT
    trace_id,
    min(timestamp) AS first_seen,
    max(timestamp) AS last_seen,
    countDistinct(service_name) AS num_services
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY trace_id
ORDER BY first_seen DESC
LIMIT 100;

-- Listar todos os traces com detalhes
SELECT
    timestamp,
    trace_id,
    span_id,
    parent_span_id,
    service_name,
    operation_name,
    duration_ns / 1000000 AS duration_ms,
    status_code,
    span_kind
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
ORDER BY timestamp DESC
LIMIT 100;

-- Listar todos os logs recentes
SELECT
    timestamp,
    trace_id,
    span_id,
    severity,
    service_name,
    message
FROM observability.logs
WHERE timestamp >= now() - INTERVAL 1 HOUR
ORDER BY timestamp DESC
LIMIT 100;

-- Correlacionar traces e logs por trace_id
SELECT
    t.trace_id,
    t.service_name AS trace_service,
    t.operation_name,
    t.duration_ns / 1000000 AS duration_ms,
    t.status_code,
    l.severity AS log_severity,
    l.message AS log_message,
    l.timestamp AS log_timestamp
FROM observability.traces t
INNER JOIN observability.logs l ON t.trace_id = l.trace_id
WHERE t.timestamp >= now() - INTERVAL 1 HOUR
ORDER BY t.timestamp DESC
LIMIT 50;

-- Buscar trace específico por ID (substitua <TRACE_ID>)
SELECT
    timestamp,
    span_id,
    parent_span_id,
    service_name,
    operation_name,
    duration_ns / 1000000 AS duration_ms,
    status_code,
    attributes
FROM observability.traces
WHERE trace_id = '<TRACE_ID>'
ORDER BY timestamp;

-- Buscar logs de um trace específico (substitua <TRACE_ID>)
SELECT
    timestamp,
    severity,
    service_name,
    message,
    attributes
FROM observability.logs
WHERE trace_id = '<TRACE_ID>'
ORDER BY timestamp;

-- Contar total de traces e logs
SELECT
    'traces' AS type,
    count() AS total,
    min(timestamp) AS oldest,
    max(timestamp) AS newest
FROM observability.traces
UNION ALL
SELECT
    'logs' AS type,
    count() AS total,
    min(timestamp) AS oldest,
    max(timestamp) AS newest
FROM observability.logs;

-- ========================================
-- QUERIES AVANÇADAS
-- ========================================

-- Top 10 operações mais lentas
SELECT
    service_name,
    operation_name,
    count() AS count,
    avg(duration_ns) / 1000000 AS avg_ms,
    quantile(0.99)(duration_ns) / 1000000 AS p99_ms
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY service_name, operation_name
ORDER BY p99_ms DESC
LIMIT 10;

-- Taxa de erro por serviço
SELECT
    service_name,
    countIf(status_code = 'ERROR') AS errors,
    count() AS total,
    (errors / total) * 100 AS error_rate_pct
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY service_name
ORDER BY error_rate_pct DESC;

-- Trace completo (distributed tracing)
-- Substitua <TRACE_ID> pelo trace_id desejado
SELECT
    span_id,
    parent_span_id,
    service_name,
    operation_name,
    duration_ns / 1000000 AS duration_ms,
    status_code,
    attributes
FROM observability.traces
WHERE trace_id = '<TRACE_ID>'
ORDER BY timestamp;

-- Análise de throughput por minuto
SELECT
    toStartOfMinute(timestamp) AS minute,
    service_name,
    count() AS requests_per_minute
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY minute, service_name
ORDER BY minute DESC;

-- Logs correlacionados com traces
-- Substitua <TRACE_ID> pelo trace_id desejado
SELECT
    l.timestamp,
    l.service_name,
    l.severity,
    l.message,
    t.operation_name,
    t.duration_ns / 1000000 AS duration_ms
FROM observability.logs l
LEFT JOIN observability.traces t ON l.trace_id = t.trace_id
WHERE l.trace_id = '<TRACE_ID>'
ORDER BY l.timestamp;

-- RED Metrics (Rate, Errors, Duration) por serviço
SELECT
    service_name,
    count() AS request_rate,
    countIf(status_code = 'ERROR') AS error_count,
    (error_count / request_rate) * 100 AS error_rate_pct,
    avg(duration_ns) / 1000000 AS avg_duration_ms,
    quantile(0.50)(duration_ns) / 1000000 AS p50_ms,
    quantile(0.95)(duration_ns) / 1000000 AS p95_ms,
    quantile(0.99)(duration_ns) / 1000000 AS p99_ms
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY service_name
ORDER BY request_rate DESC;

-- Análise de distribuição de duração
SELECT
    service_name,
    operation_name,
    quantile(0.25)(duration_ns) / 1000000 AS p25_ms,
    quantile(0.50)(duration_ns) / 1000000 AS p50_ms,
    quantile(0.75)(duration_ns) / 1000000 AS p75_ms,
    quantile(0.90)(duration_ns) / 1000000 AS p90_ms,
    quantile(0.95)(duration_ns) / 1000000 AS p95_ms,
    quantile(0.99)(duration_ns) / 1000000 AS p99_ms,
    max(duration_ns) / 1000000 AS max_ms
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY service_name, operation_name
ORDER BY p99_ms DESC;

-- Tendências de erro ao longo do tempo
SELECT
    toStartOfMinute(timestamp) AS minute,
    service_name,
    countIf(status_code = 'ERROR') AS error_count,
    count() AS total_requests,
    (error_count / total_requests) * 100 AS error_rate_pct
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 6 HOUR
GROUP BY minute, service_name
ORDER BY minute DESC, error_rate_pct DESC;

-- Análise de volume de dados
SELECT
    toStartOfHour(timestamp) AS hour,
    service_name,
    count() AS span_count,
    uniq(trace_id) AS unique_traces,
    avg(length(toString(attributes))) AS avg_attributes_size,
    sum(length(toString(attributes))) AS total_attributes_size
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 24 HOUR
GROUP BY hour, service_name
ORDER BY hour DESC;

-- Operações mais frequentes
SELECT
    service_name,
    operation_name,
    count() AS execution_count,
    avg(duration_ns) / 1000000 AS avg_duration_ms
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY service_name, operation_name
ORDER BY execution_count DESC
LIMIT 20;

-- Análise de severidade de logs
SELECT
    severity,
    service_name,
    count() AS log_count,
    uniq(trace_id) AS unique_traces
FROM observability.logs
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY severity, service_name
ORDER BY log_count DESC;

-- Traces com maior duração total (mais caro)
SELECT
    trace_id,
    service_name,
    count() AS span_count,
    sum(duration_ns) / 1000000 AS total_duration_ms,
    max(duration_ns) / 1000000 AS max_span_duration_ms
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY trace_id, service_name
ORDER BY total_duration_ms DESC
LIMIT 10;

-- Análise de propagação de contexto
SELECT
    trace_id,
    countDistinct(service_name) AS services_involved,
    groupArray(DISTINCT service_name) AS services,
    count() AS total_spans,
    min(timestamp) AS trace_start,
    max(timestamp) AS trace_end,
    dateDiff('millisecond', trace_start, trace_end) AS trace_duration_ms
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY trace_id
HAVING services_involved > 1
ORDER BY trace_duration_ms DESC
LIMIT 10;
