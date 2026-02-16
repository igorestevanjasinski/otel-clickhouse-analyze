# Prometheus + Grafana - Métricas e Dashboards

## Visão Geral

Esta documentação descreve a implementação completa de Prometheus e Grafana para monitoramento de métricas no projeto de observabilidade.

## Arquitetura

```
┌─────────────┐     ┌──────────────┐     ┌────────────┐
│ Product API │────▶│  Prometheus  │────▶│  Grafana   │
│ /metrics    │     │   (scrape)   │     │ Dashboards │
└─────────────┘     └──────────────┘     └────────────┘
                           │                    │
┌─────────────┐           │                    │
│  Consumer   │───────────┘                    │
│ /metrics    │                                │
└─────────────┘                                │
                                               │
┌─────────────┐                                │
│ ClickHouse  │───────────────────────────────┘
│  (traces)   │
└─────────────┘
```

## Componentes

### 1. Prometheus

**Configuração**: `infrastructure/observability/prometheus/prometheus.yml`

#### Scrape Configs

- **product-api**: `http://product-api:8000/metrics` (10s interval)
- **product-consumer**: `http://product-consumer:8081/metrics` (10s interval)
- **otel-collector**: `http://otel-collector:8888/metrics` (15s interval)
- **prometheus**: `http://localhost:9090/metrics` (self-monitoring)

#### Recording Rules

Arquivo: `infrastructure/observability/prometheus/rules/app-rules.yml`

**Product API Rules**:
- `job:product_api:requests_per_second`: Taxa de requisições por segundo
- `job:product_api:error_rate`: Taxa de erro (%)
- `job:product_api:latency_seconds:avg`: Latência média

**Product Consumer Rules**:
- `job:product_consumer:messages_per_second`: Mensagens por segundo
- `job:product_consumer:error_rate`: Taxa de erro de processamento
- `job:product_consumer:processing_seconds:avg`: Tempo médio de processamento

### 2. Métricas Expostas

#### Product API (Python)

```python
# Counter - Total de produtos criados
products_created_total{status="success|error", endpoint="/products"}

# Histogram - Latência de criação de produtos
product_creation_duration_seconds{endpoint="/products"}
# Buckets: 5ms, 10ms, 25ms, 50ms, 75ms, 100ms, 250ms, 500ms, 750ms, 1s, 2.5s, 5s, 7.5s, 10s

# Gauge - Requisições ativas
active_requests{endpoint="/products"}

# Counter - Mensagens publicadas no Kafka
kafka_messages_published_total{topic="products.events", status="success|error"}
```

**Implementação**: 
- `app/telemetry/prometheus_metrics.py`: Definição de métricas
- `app/api/routes.py`: Instrumentação nas rotas
- Endpoint: `http://localhost:8000/metrics`

#### Product Consumer (Go)

```go
// Counter - Total de mensagens consumidas
messages_consumed_total{topic="products.events", status="success|error"}

// Histogram - Duração de processamento
message_processing_duration_seconds{topic="products.events"}
// Buckets: 5ms, 10ms, 25ms, 50ms, 75ms, 100ms, 250ms, 500ms, 750ms, 1s, 2.5s, 5s, 7.5s, 10s

// Gauge - Mensagens em processamento
messages_in_flight

// Counter - Operações de banco de dados
database_operations_total{operation="insert", status="success|error|duplicate"}

// Histogram - Latência de operações de banco
database_operation_duration_seconds{operation="insert"}
// Buckets: 1ms, 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s
```

**Implementação**:
- `pkg/metrics/prometheus.go`: Definição e servidor HTTP
- `internal/consumer/message_handler.go`: Instrumentação do handler
- `internal/repository/product_repository.go`: Métricas de database
- Endpoint: `http://localhost:8081/metrics`

### 3. Grafana

**Datasources**: 
- **Prometheus**: `http://prometheus:9090` (default)
- **ClickHouse**: `http://clickhouse:8123` (database: observability)

**Provisionamento**: `infrastructure/observability/grafana/provisioning/`

#### Dashboard 1: Microservices Overview - RED Metrics

Arquivo: `infrastructure/observability/grafana/dashboards/microservices-overview.json`

**Painéis**:

1. **Request Rate (RPS) - Product API**
   - Query: `sum(rate(products_created_total[5m])) by (status)`
   - Visualização: Graph
   - Labels: success, error

2. **Error Rate (%) - Product API**
   - Query: `100 * sum(rate(products_created_total{status="error"}[5m])) / sum(rate(products_created_total[5m]))`
   - Visualização: Graph
   - Alert: > 5%

3. **P95 Latency - Product API**
   - Queries:
     - P50: `histogram_quantile(0.50, sum(rate(product_creation_duration_seconds_bucket[5m])) by (le))`
     - P95: `histogram_quantile(0.95, ...)`
     - P99: `histogram_quantile(0.99, ...)`
   - Visualização: Graph

4. **Active Requests - Product API**
   - Query: `sum(active_requests)`
   - Visualização: Stat

5. **Total Requests - Product API**
   - Query: `sum(products_created_total)`
   - Visualização: Stat

6. **Message Processing Rate - Consumer**
   - Query: `sum(rate(messages_consumed_total[5m])) by (status)`
   - Visualização: Graph

7. **Consumer Processing Latency**
   - Query: `histogram_quantile(0.95, sum(rate(message_processing_duration_seconds_bucket[5m])) by (le))`
   - Visualização: Graph

8. **Database Operations**
   - Query: `sum(rate(database_operations_total[5m])) by (operation, status)`
   - Visualização: Graph

9. **Database Operation Latency**
   - Query: `histogram_quantile(0.95, sum(rate(database_operation_duration_seconds_bucket[5m])) by (le, operation))`
   - Visualização: Graph

10. **Messages In Flight**
    - Query: `messages_in_flight`
    - Visualização: Stat

11. **Kafka Messages Published**
    - Query: `sum(kafka_messages_published_total)`
    - Visualização: Stat

#### Dashboard 2: Distributed Tracing - ClickHouse

Arquivo: `infrastructure/observability/grafana/dashboards/clickhouse-traces.json`

**Painéis**:

1. **Trace Count by Service**
   - Query: `SELECT ServiceName, count() as count FROM observability.traces WHERE Timestamp >= now() - INTERVAL 1 HOUR GROUP BY ServiceName`
   - Datasource: ClickHouse
   - Visualização: Pie Chart

2. **Top 10 Slowest Operations**
   - Query: `SELECT ServiceName, SpanName, quantile(0.99)(Duration) / 1000000 as p99_ms FROM observability.traces WHERE Timestamp >= now() - INTERVAL 1 HOUR GROUP BY ServiceName, SpanName ORDER BY p99_ms DESC LIMIT 10`
   - Datasource: ClickHouse
   - Visualização: Table

3. **Request Rate per Service**
   - Query: `SELECT toStartOfMinute(Timestamp) as time, ServiceName, count() / 60 as rps FROM observability.traces WHERE Timestamp >= now() - INTERVAL 1 HOUR GROUP BY time, ServiceName ORDER BY time`
   - Datasource: ClickHouse
   - Visualização: Graph

4. **Error Rate by Service**
   - Query: `SELECT toStartOfMinute(Timestamp) as time, ServiceName, countIf(StatusCode = 'STATUS_CODE_ERROR') * 100.0 / count() as error_rate FROM observability.traces WHERE Timestamp >= now() - INTERVAL 1 HOUR GROUP BY time, ServiceName ORDER BY time`
   - Datasource: ClickHouse
   - Visualização: Graph

5. **Latency Heatmap**
   - Query: `SELECT toStartOfMinute(Timestamp) as time, Duration / 1000000 as duration_ms FROM observability.traces WHERE Timestamp >= now() - INTERVAL 1 HOUR ORDER BY time`
   - Datasource: ClickHouse
   - Visualização: Heatmap

6. **Distributed Traces (Multi-Service)**
   - Query: `SELECT TraceId, groupArray(ServiceName) as services, count() as span_count, max(Duration) / 1000000 as total_duration_ms, countIf(StatusCode = 'STATUS_CODE_ERROR') as errors FROM observability.traces WHERE Timestamp >= now() - INTERVAL 1 HOUR GROUP BY TraceId HAVING length(services) > 1 ORDER BY total_duration_ms DESC LIMIT 20`
   - Datasource: ClickHouse
   - Visualização: Table

7. **P50/P95/P99 Latency by Service**
   - Query: `SELECT toStartOfMinute(Timestamp) as time, ServiceName, quantile(0.50)(Duration) / 1000000 as p50, quantile(0.95)(Duration) / 1000000 as p95, quantile(0.99)(Duration) / 1000000 as p99 FROM observability.traces WHERE Timestamp >= now() - INTERVAL 1 HOUR GROUP BY time, ServiceName ORDER BY time`
   - Datasource: ClickHouse
   - Visualização: Graph

8. **Recent Errors**
   - Query: `SELECT Timestamp, TraceId, ServiceName, SpanName, StatusCode FROM observability.traces WHERE StatusCode = 'STATUS_CODE_ERROR' AND Timestamp >= now() - INTERVAL 1 HOUR ORDER BY Timestamp DESC LIMIT 50`
   - Datasource: ClickHouse
   - Visualização: Table

## Uso

### Iniciar Serviços

```bash
# Subir Prometheus e Grafana
docker compose up -d prometheus grafana

# Verificar status
docker compose ps prometheus grafana

# Ver logs
docker compose logs -f prometheus grafana
```

### Acessar Interfaces

**Prometheus**:
- URL: http://localhost:9090
- Targets: http://localhost:9090/targets
- Rules: http://localhost:9090/rules
- Graph: http://localhost:9090/graph

**Grafana**:
- URL: http://localhost:3000
- Login: admin / admin (configurável via .env)
- Dashboards: http://localhost:3000/dashboards

### Testar Métricas

```bash
# Verificar endpoint da API
curl http://localhost:8000/metrics | grep products_created_total

# Verificar endpoint do Consumer
curl http://localhost:8081/metrics | grep messages_consumed_total

# Query manual no Prometheus
curl 'http://localhost:9090/api/v1/query?query=up'
curl 'http://localhost:9090/api/v1/query?query=products_created_total'

# Gerar carga para visualizar nos dashboards
for i in {1..100}; do
  curl -X POST http://localhost:8000/products \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"Test $i\",\"price\":99.99}" &
done
wait
```

### Script de Teste Automatizado

```bash
# Executar suite completa de testes
./scripts/test-grafana.sh

# O script irá:
# 1. Verificar saúde de Prometheus e Grafana
# 2. Validar targets do Prometheus
# 3. Testar endpoints /metrics
# 4. Gerar carga de 100 requisições
# 5. Verificar métricas no Prometheus
# 6. Validar recording rules
# 7. Verificar datasources do Grafana
# 8. Validar dashboards provisionados
```

## Queries Úteis

### Prometheus (PromQL)

#### RED Metrics

```promql
# Rate - Taxa de requisições
sum(rate(products_created_total[5m])) by (status)

# Errors - Taxa de erro
100 * sum(rate(products_created_total{status="error"}[5m])) / sum(rate(products_created_total[5m]))

# Duration - Percentis de latência
histogram_quantile(0.95, sum(rate(product_creation_duration_seconds_bucket[5m])) by (le))
histogram_quantile(0.99, sum(rate(product_creation_duration_seconds_bucket[5m])) by (le))
```

#### Agregações

```promql
# Taxa de requisições por minuto
sum(increase(products_created_total[1m]))

# Requisições ativas médias
avg_over_time(active_requests[5m])

# Total de erros nas últimas 24h
sum(increase(products_created_total{status="error"}[24h]))
```

#### Consumer Metrics

```promql
# Mensagens processadas por segundo
sum(rate(messages_consumed_total[5m]))

# Lag de processamento (inferido por in-flight)
messages_in_flight

# Operações de DB por segundo
sum(rate(database_operations_total[5m])) by (operation, status)
```

## Alertas

Os alertas serão implementados na Fase 12 (SLIs/SLOs + Alerting).

Exemplos de alertas planejados:
- Error rate > 5% por 5 minutos
- P95 latency > 1s por 5 minutos
- Database operation errors > 10 em 1 minuto
- Messages in flight > 100 (possível backup)

## Troubleshooting

### Prometheus Não Scraping

```bash
# Verificar targets
curl http://localhost:9090/api/v1/targets

# Verificar logs
docker compose logs prometheus

# Validar configuração
docker compose exec prometheus promtool check config /etc/prometheus/prometheus.yml
```

### Grafana Datasource Não Conecta

```bash
# Verificar rede
docker compose exec grafana ping prometheus
docker compose exec grafana ping clickhouse

# Verificar logs
docker compose logs grafana

# Testar datasource manualmente
curl -u admin:admin http://localhost:3000/api/datasources
```

### Métricas Não Aparecem

```bash
# Verificar se endpoints estão respondendo
curl http://localhost:8000/metrics
curl http://localhost:8081/metrics

# Verificar se Prometheus está scraping
curl 'http://localhost:9090/api/v1/query?query=up{job="product-api"}'

# Forçar reload do Prometheus
curl -X POST http://localhost:9090/-/reload
```

### Dashboard Vazio

1. Verificar time range (último 1h por padrão)
2. Gerar carga de teste
3. Verificar query no Explore
4. Verificar se datasource está correto

## Performance

### Prometheus Retenção

- **Padrão**: 30 dias (`--storage.tsdb.retention.time=30d`)
- **Espaço**: ~1-2GB para 30 dias de métricas (depende da cardinalidade)

### Recording Rules

As recording rules são calculadas a cada 30s e armazenam resultados pré-calculados:
- Reduz tempo de query em dashboards
- Útil para cálculos complexos (percentis, taxas)
- Armazenadas como séries temporais normais

### Scrape Intervals

- **product-api**: 10s (alta frequência por ser user-facing)
- **product-consumer**: 10s (processamento contínuo)
- **otel-collector**: 15s (métricas internas)
- **prometheus**: 15s (self-monitoring)

## Próximos Passos

Fase 12 implementará:
1. **SLIs/SLOs**: Definição de Service Level Indicators e Objectives
2. **Alerting**: Alertmanager configurado com regras baseadas em SLOs
3. **Runbooks**: Documentação de resposta a incidentes
4. **Error Budgets**: Cálculo e tracking de budgets de erro

## Referências

- [Prometheus Documentation](https://prometheus.io/docs/)
- [PromQL Guide](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Grafana Provisioning](https://grafana.com/docs/grafana/latest/administration/provisioning/)
- [RED Method](https://www.weave.works/blog/the-red-method-key-metrics-for-microservices-architecture/)
- [Histograms and Summaries](https://prometheus.io/docs/practices/histograms/)
