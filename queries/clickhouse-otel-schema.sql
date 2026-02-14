-- ========================================
-- QUERIES BÁSICAS (Schema do OTel Collector)
-- ========================================

-- Listar todos os trace_ids recentes
SELECT DISTINCT
    TraceId,
    min(Timestamp) AS first_seen,
    max(Timestamp) AS last_seen,
    countDistinct(ServiceName) AS num_services
FROM observability.traces
WHERE Timestamp >= now() - INTERVAL 1 HOUR
GROUP BY TraceId
ORDER BY first_seen DESC
LIMIT 100;

-- Listar todos os traces com detalhes
SELECT
    Timestamp,
    TraceId,
    SpanId,
    ParentSpanId,
    ServiceName,
    SpanName AS operation_name,
    Duration / 1000000 AS duration_ms,
    StatusCode,
    SpanKind
FROM observability.traces
WHERE Timestamp >= now() - INTERVAL 1 HOUR
ORDER BY Timestamp DESC
LIMIT 100;

-- Listar todos os logs recentes
SELECT
    Timestamp,
    TraceId,
    SeverityText AS severity,
    ServiceName,
    Body AS message
FROM observability.logs
WHERE Timestamp >= now() - INTERVAL 1 HOUR
ORDER BY Timestamp DESC
LIMIT 100;

-- Correlacionar traces e logs por trace_id
SELECT
    t.TraceId,
    t.ServiceName AS trace_service,
    t.SpanName AS operation_name,
    t.Duration / 1000000 AS duration_ms,
    t.StatusCode,
    l.SeverityText AS log_severity,
    l.Body AS log_message,
    l.Timestamp AS log_timestamp
FROM observability.traces t
INNER JOIN observability.logs l ON t.TraceId = l.TraceId
WHERE t.Timestamp >= now() - INTERVAL 1 HOUR
ORDER BY t.Timestamp DESC
LIMIT 50;

-- Buscar trace específico por ID (substitua <TRACE_ID>)
SELECT
    Timestamp,
    SpanId,
    ParentSpanId,
    ServiceName,
    SpanName AS operation_name,
    Duration / 1000000 AS duration_ms,
    StatusCode,
    SpanAttributes
FROM observability.traces
WHERE TraceId = '<TRACE_ID>'
ORDER BY Timestamp;

-- Buscar logs de um trace específico (substitua <TRACE_ID>)
SELECT
    Timestamp,
    SeverityText AS severity,
    ServiceName,
    Body AS message,
    LogAttributes
FROM observability.logs
WHERE TraceId = '<TRACE_ID>'
ORDER BY Timestamp;

-- Contar total de traces e logs
SELECT
    'traces' AS type,
    count() AS total,
    min(Timestamp) AS oldest,
    max(Timestamp) AS newest
FROM observability.traces
UNION ALL
SELECT
    'logs' AS type,
    count() AS total,
    min(Timestamp) AS oldest,
    max(Timestamp) AS newest
FROM observability.logs;

-- ========================================
-- QUERIES AVANÇADAS
-- ========================================

-- Top 10 operações mais lentas
SELECT
    ServiceName,
    SpanName AS operation_name,
    count() AS count,
    avg(Duration) / 1000000 AS avg_ms,
    quantile(0.99)(Duration) / 1000000 AS p99_ms
FROM observability.traces
WHERE Timestamp >= now() - INTERVAL 1 HOUR
GROUP BY ServiceName, SpanName
ORDER BY p99_ms DESC
LIMIT 10;

-- Taxa de erro por serviço
SELECT
    ServiceName,
    countIf(StatusCode = 'STATUS_CODE_ERROR') AS errors,
    count() AS total,
    (errors / total) * 100 AS error_rate_pct
FROM observability.traces
WHERE Timestamp >= now() - INTERVAL 1 HOUR
GROUP BY ServiceName
ORDER BY error_rate_pct DESC;

-- Análise de throughput por minuto
SELECT
    toStartOfMinute(Timestamp) AS minute,
    ServiceName,
    count() AS requests_per_minute
FROM observability.traces
WHERE Timestamp >= now() - INTERVAL 1 HOUR
GROUP BY minute, ServiceName
ORDER BY minute DESC;

-- RED Metrics (Rate, Errors, Duration) por serviço
SELECT
    ServiceName,
    count() AS request_rate,
    countIf(StatusCode = 'STATUS_CODE_ERROR') AS error_count,
    round((error_count / request_rate) * 100, 2) AS error_rate_pct,
    round(avg(Duration) / 1000000, 2) AS avg_duration_ms,
    round(quantile(0.50)(Duration) / 1000000, 2) AS p50_ms,
    round(quantile(0.95)(Duration) / 1000000, 2) AS p95_ms,
    round(quantile(0.99)(Duration) / 1000000, 2) AS p99_ms
FROM observability.traces
WHERE Timestamp >= now() - INTERVAL 1 HOUR
GROUP BY ServiceName
ORDER BY request_rate DESC;

-- Operações mais frequentes
SELECT
    ServiceName,
    SpanName AS operation_name,
    count() AS execution_count,
    round(avg(Duration) / 1000000, 2) AS avg_duration_ms
FROM observability.traces
WHERE Timestamp >= now() - INTERVAL 1 HOUR
GROUP BY ServiceName, SpanName
ORDER BY execution_count DESC
LIMIT 20;

-- Análise de severidade de logs
SELECT
    SeverityText AS severity,
    ServiceName,
    count() AS log_count,
    uniq(TraceId) AS unique_traces
FROM observability.logs
WHERE Timestamp >= now() - INTERVAL 1 HOUR
GROUP BY SeverityText, ServiceName
ORDER BY log_count DESC;

-- Traces com maior duração total (mais caro)
SELECT
    TraceId,
    ServiceName,
    count() AS span_count,
    sum(Duration) / 1000000 AS total_duration_ms,
    max(Duration) / 1000000 AS max_span_duration_ms
FROM observability.traces
WHERE Timestamp >= now() - INTERVAL 1 HOUR
GROUP BY TraceId, ServiceName
ORDER BY total_duration_ms DESC
LIMIT 10;

-- Análise de propagação de contexto (distributed traces)
SELECT
    TraceId,
    countDistinct(ServiceName) AS services_involved,
    groupArray(DISTINCT ServiceName) AS services,
    count() AS total_spans,
    min(Timestamp) AS trace_start,
    max(Timestamp) AS trace_end,
    dateDiff('millisecond', trace_start, trace_end) AS trace_duration_ms
FROM observability.traces
WHERE Timestamp >= now() - INTERVAL 1 HOUR
GROUP BY TraceId
HAVING services_involved > 1
ORDER BY trace_duration_ms DESC
LIMIT 10;

-- Análise de distribuição de duração
SELECT
    ServiceName,
    SpanName AS operation_name,
    round(quantile(0.25)(Duration) / 1000000, 2) AS p25_ms,
    round(quantile(0.50)(Duration) / 1000000, 2) AS p50_ms,
    round(quantile(0.75)(Duration) / 1000000, 2) AS p75_ms,
    round(quantile(0.90)(Duration) / 1000000, 2) AS p90_ms,
    round(quantile(0.95)(Duration) / 1000000, 2) AS p95_ms,
    round(quantile(0.99)(Duration) / 1000000, 2) AS p99_ms,
    round(max(Duration) / 1000000, 2) AS max_ms
FROM observability.traces
WHERE Timestamp >= now() - INTERVAL 1 HOUR
GROUP BY ServiceName, SpanName
ORDER BY p99_ms DESC;
