# ClickHouse para Observabilidade em Escala

## Visão Geral

ClickHouse configurado como backend de observabilidade para armazenar e consultar bilhões de traces e logs com performance otimizada.

## Arquitetura

```
┌─────────────┐         ┌──────────────────┐         ┌──────────────┐
│  Services   │ ─OTLP─> │ OTel Collector   │ ────>   │ ClickHouse   │
│ (API/Cons.) │         │  (Processor)     │         │  (Storage)   │
└─────────────┘         └──────────────────┘         └──────────────┘
                                │                             │
                                └─────────────────────────────┘
                                        Queries
```

## Schema Otimizado

### Tabela de Traces

```sql
CREATE TABLE observability.traces (
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
    events Array(Tuple(...)) CODEC(ZSTD),
    resource_attributes Map(String, String) CODEC(ZSTD),
    ...
) ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (service_name, operation_name, timestamp)
TTL timestamp + INTERVAL 30 DAY
```

**Otimizações:**
- **Compression**: ZSTD para campos grandes, Delta+ZSTD para timestamps
- **LowCardinality**: Para colunas com poucos valores distintos
- **Particionamento**: Por dia para facilitar drops de dados antigos
- **TTL**: Retenção automática de 30 dias
- **Índices**: Bloom filter para trace_id, set indexes para service e operation

### Materialized View para Métricas

```sql
CREATE MATERIALIZED VIEW observability.traces_metrics_mv
ENGINE = SummingMergeTree()
AS SELECT
    toStartOfMinute(timestamp) AS timestamp_minute,
    service_name,
    operation_name,
    status_code,
    count() AS request_count,
    avg(duration_ns) AS avg_duration_ns,
    quantile(0.95)(duration_ns) AS p95_duration_ns,
    quantile(0.99)(duration_ns) AS p99_duration_ns
FROM observability.traces
GROUP BY timestamp_minute, service_name, operation_name, status_code
```

**Benefícios:**
- Agregações pré-calculadas
- Queries de métricas instantâneas
- Redução de processamento em consultas

## Queries Comuns

### 1. RED Metrics (Rate, Errors, Duration)

```sql
SELECT 
    service_name,
    count() AS request_rate,
    countIf(status_code = 'ERROR') AS error_count,
    (error_count / request_rate) * 100 AS error_rate_pct,
    avg(duration_ns) / 1000000 AS avg_duration_ms,
    quantile(0.95)(duration_ns) / 1000000 AS p95_ms,
    quantile(0.99)(duration_ns) / 1000000 AS p99_ms
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY service_name
```

### 2. Top Operações Mais Lentas

```sql
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
LIMIT 10
```

### 3. Análise de Distributed Tracing

```sql
SELECT 
    trace_id,
    countDistinct(service_name) AS services_involved,
    groupArray(DISTINCT service_name) AS services,
    count() AS total_spans,
    dateDiff('millisecond', min(timestamp), max(timestamp)) AS duration_ms
FROM observability.traces
WHERE timestamp >= now() - INTERVAL 1 HOUR
GROUP BY trace_id
HAVING services_involved > 1
ORDER BY duration_ms DESC
```

### 4. Trace Completo por ID

```sql
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
ORDER BY timestamp
```

## Scripts de Análise

### Executar Análise Completa

```bash
./scripts/analyze-traces.sh
```

**Saída inclui:**
- Database size e row count
- Top 10 operações mais lentas
- Taxa de erro por serviço
- RED metrics
- Throughput por minuto
- Operações mais frequentes
- Distribuição de severidade de logs
- Traces distribuídos
- Volume de dados por serviço

### Executar Query Customizada

```bash
docker compose exec clickhouse clickhouse-client --query="
SELECT * FROM observability.traces LIMIT 10
"
```

## Performance

### Capacidade de Escala

**Configuração atual suporta:**
- Bilhões de spans
- Milhões de spans por segundo (ingestão)
- Queries sub-segundo em datasets grandes
- Compression ratio ~10:1

### Otimizações Aplicadas

1. **Compression Codecs**:
   - `Delta, ZSTD` para timestamps (alta compressão)
   - `T64, ZSTD` para durations (números)
   - `ZSTD` para strings e maps

2. **Data Types**:
   - `LowCardinality(String)` para campos com poucos valores
   - `Map(String, String)` para attributes flexíveis
   - `DateTime64(9)` para precisão de nanossegundos

3. **Indexes**:
   - Bloom filter para `trace_id` (busca exata)
   - Set indexes para `service_name` e `operation_name`
   - Primary key otimizada: `(service_name, operation_name, timestamp)`

4. **Particionamento**:
   - Por dia (`toYYYYMMDD(timestamp)`)
   - Permite drop eficiente de dados antigos
   - Queries filtradas por data usam apenas partições relevantes

## Monitoramento do ClickHouse

### Health Check

```bash
docker compose exec clickhouse clickhouse-client --query "SELECT version()"
```

### Verificar Tamanho das Tabelas

```bash
docker compose exec clickhouse clickhouse-client --query="
SELECT 
    database,
    table,
    formatReadableSize(sum(bytes)) AS size,
    sum(rows) AS rows,
    count() AS parts
FROM system.parts
WHERE database = 'observability'
GROUP BY database, table
FORMAT PrettyCompact
"
```

### Verificar Compression Ratio

```bash
docker compose exec clickhouse clickhouse-client --query="
SELECT 
    table,
    formatReadableSize(sum(bytes_on_disk)) AS compressed,
    formatReadableSize(sum(data_uncompressed_bytes)) AS uncompressed,
    round(sum(data_uncompressed_bytes) / sum(bytes_on_disk), 2) AS ratio
FROM system.parts
WHERE database = 'observability'
GROUP BY table
FORMAT PrettyCompact
"
```

### Queries Ativas

```bash
docker compose exec clickhouse clickhouse-client --query="
SELECT 
    query_id,
    user,
    elapsed,
    read_rows,
    formatReadableSize(memory_usage) AS memory,
    query
FROM system.processes
FORMAT PrettyCompact
"
```

## Manutenção

### Limpeza Manual de Dados Antigos

```sql
-- Deletar partições antigas (mais de 30 dias)
ALTER TABLE observability.traces DROP PARTITION '20240101'
```

### Optimize Merge

```sql
-- Forçar merge de parts (background automation normalmente faz isso)
OPTIMIZE TABLE observability.traces FINAL
```

### Backup

```bash
# Backup via clickhouse-backup (instalar separadamente)
docker compose exec clickhouse clickhouse-backup create backup_name

# Export para arquivo
docker compose exec clickhouse clickhouse-client --query="
SELECT * FROM observability.traces
WHERE timestamp >= '2024-01-01'
FORMAT Parquet
" > traces_backup.parquet
```

## Integração com Grafana

1. Adicionar datasource ClickHouse no Grafana
2. Endpoint: `http://clickhouse:8123`
3. Database: `observability`
4. User: `default`

**Exemplo de query para dashboard:**

```sql
SELECT 
    $__timeInterval(timestamp) as t,
    service_name,
    avg(duration_ns) / 1000000 AS avg_latency_ms
FROM observability.traces
WHERE $__timeFilter(timestamp)
GROUP BY t, service_name
ORDER BY t
```

## Troubleshooting

### ClickHouse não inicia

```bash
# Ver logs
docker compose logs clickhouse

# Verificar permissões de arquivo
ls -la infrastructure/observability/clickhouse/

# Verificar ulimits
docker compose exec clickhouse sh -c 'ulimit -n'
```

### Queries lentas

```bash
# Verificar query log
docker compose exec clickhouse clickhouse-client --query="
SELECT 
    query,
    query_duration_ms,
    read_rows,
    formatReadableSize(memory_usage) AS memory
FROM system.query_log
WHERE type = 'QueryFinish'
ORDER BY query_duration_ms DESC
LIMIT 10
FORMAT PrettyCompact
"

# Analisar query plan
EXPLAIN SELECT ... FROM observability.traces ...
```

### Collector não exporta para ClickHouse

```bash
# Verificar logs do collector
docker compose logs otel-collector | grep clickhouse

# Verificar conectividade
docker compose exec otel-collector wget -O- http://clickhouse:8123/ping

# Verificar tabelas existem
docker compose exec clickhouse clickhouse-client --query "SHOW TABLES FROM observability"
```

## Recursos Adicionais

- [ClickHouse Documentation](https://clickhouse.com/docs)
- [ClickHouse Performance](https://clickhouse.com/docs/en/operations/performance)
- [OpenTelemetry ClickHouse Exporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/clickhouseexporter)
- [ClickHouse Query Optimization](https://clickhouse.com/docs/en/guides/improving-query-performance)

## Próximos Passos

1. Configurar Grafana com datasource ClickHouse
2. Criar dashboards de observabilidade
3. Configurar alerting baseado em queries ClickHouse
4. Implementar retenção customizada por tipo de span
5. Adicionar tabelas de métricas (além de traces)
