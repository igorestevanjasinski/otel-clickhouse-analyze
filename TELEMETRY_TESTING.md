# Fase 9: OpenTelemetry Instrumentation - Testing Guide

## Objetivo

Validar que ambos os serviços (Python API e Go Consumer) estão gerando telemetria corretamente usando OpenTelemetry antes de configurar exportação para ClickHouse.

## Arquitetura de Telemetria

```
┌─────────────────┐         ┌──────────────────┐
│  Product API    │────────▶│  OTLP Collector  │
│  (Python)       │  gRPC   │                  │
│  Port: 8000     │  :4317  │  Debug Exporter  │
└─────────────────┘         │                  │
                            │  Logs: verbose   │
┌─────────────────┐         └──────────────────┘
│  Consumer       │────────▶         │
│  (Go)           │  gRPC            │
│  Port: 8081     │  :4317           ▼
└─────────────────┘         ┌──────────────────┐
                            │   Jaeger UI      │
         Kafka              │   Port: 16686    │
    (trace context)         └──────────────────┘
```

## 1. Preparação

### Verificar infraestrutura base

```bash
# Deve ter Kafka, Zookeeper e PostgreSQL rodando
docker compose ps

# Se não estiver, inicie:
docker compose up -d
```

### Verificar dependências Python

```bash
cd services/product-api

# Instalar/atualizar dependências OpenTelemetry
pip install -r requirements.txt

# Verificar instalação
python -c "import opentelemetry; print(opentelemetry.__version__)"
```

### Verificar dependências Go

```bash
cd services/product-consumer

# Baixar dependências
go mod download
go mod tidy

# Compilar
go build -o bin/consumer cmd/main.go
```

## 2. Iniciar Observability Stack

```bash
# Subir OpenTelemetry Collector e Jaeger
docker compose -f docker-compose.observability.yml up -d

# Verificar status
docker compose -f docker-compose.observability.yml ps

# Aguardar Collector estar pronto
curl http://localhost:13133/
```

**Portas expostas:**
- `4317`: OTLP gRPC receiver
- `4318`: OTLP HTTP receiver
- `13133`: Health check
- `16686`: Jaeger UI

## 3. Teste Automatizado (Recomendado)

```bash
# Executar script de teste completo
./scripts/test-telemetry.sh
```

O script irá:
1. ✅ Iniciar observability stack
2. ✅ Instalar dependências Python
3. ✅ Compilar Go consumer
4. ✅ Iniciar API e Consumer
5. ✅ Criar produto de teste
6. ✅ Validar telemetria nos logs do Collector

**Saída esperada:**
```
==================================================
  ✓ Telemetry Validation Successful!
==================================================
```

## 4. Teste Manual

### Passo 1: Iniciar Product API

```bash
cd services/product-api
source venv/bin/activate

# Com telemetria habilitada
OTLP_ENDPOINT=http://localhost:4317 \
uvicorn app.main:app --host 0.0.0.0 --port 8000
```

**Logs esperados:**
```json
{
  "timestamp": "2026-02-13T...",
  "level": "INFO",
  "message": "OpenTelemetry tracing initialized",
  "service_name": "Product API",
  "otlp_endpoint": "http://localhost:4317"
}
```

### Passo 2: Iniciar Product Consumer

```bash
cd services/product-consumer

OTLP_ENDPOINT=localhost:4317 \
./bin/consumer
```

**Logs esperados:**
```
OpenTelemetry tracing initialized: service=product-consumer endpoint=localhost:4317
OpenTelemetry metrics initialized
```

### Passo 3: Gerar Telemetria

```bash
# Criar produto
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -H "X-Correlation-ID: manual-test-$(date +%s)" \
  -d '{
    "name": "Telemetry Test Product",
    "price": 299.99,
    "description": "Testing OpenTelemetry"
  }' | jq
```

### Passo 4: Verificar Telemetria no Collector

```bash
# Ver logs do collector em tempo real
docker compose -f docker-compose.observability.yml logs -f otel-collector

# Ou filtrar apenas traces
docker compose -f docker-compose.observability.yml logs otel-collector | grep "Span #"

# Ou filtrar apenas métricas
docker compose -f docker-compose.observability.yml logs otel-collector | grep "Metric #"
```

## 5. Validações

### ✅ Traces Gerados

**No Collector logs, você deve ver:**

```
Span #0
    Trace ID       : abc123def456...
    Parent ID      : 
    ID             : 789012345678...
    Name           : POST /products
    Kind           : Server
    Start time     : 2026-02-13 10:30:00
    End time       : 2026-02-13 10:30:00.250
    Status code    : Ok
Resource attributes:
     -> service.name: Str(Product API)
     -> service.version: Str(1.0.0)
     -> deployment.environment: Str(development)
     -> service.namespace: Str(product-portfolio)
```

### ✅ Context Propagation via Kafka

**Você deve ver dois spans com MESMO Trace ID:**

1. **Span do Producer (Python API):**
```
Name: kafka.publish
Attributes:
  - messaging.system: kafka
  - messaging.destination: products.events
  - product.id: abc-123
```

2. **Span do Consumer (Go):**
```
Name: kafka.consume
Parent ID: <ID do span do producer>
Attributes:
  - messaging.system: kafka
  - messaging.destination: products.events
  - product.id: abc-123
```

### ✅ Database Operations

**Span de inserção no banco:**
```
Name: database.insert
Attributes:
  - db.system: postgresql
  - db.operation: INSERT
  - db.table: products
  - product.id: abc-123
```

### ✅ Métricas Geradas

```
Metric #0
Descriptor:
     -> Name: products.created.total
     -> Description: Total number of products created
     -> Unit: 1
     -> DataType: Sum
NumberDataPoints:
  Value: 1

Metric #1
Descriptor:
     -> Name: messages.consumed.total
     -> Description: Total number of messages consumed
```

## 6. Visualizar no Jaeger UI

```bash
# Abrir Jaeger
open http://localhost:16686
# ou
xdg-open http://localhost:16686
```

**No Jaeger:**
1. Selecione serviço: `Product API`
2. Clique em "Find Traces"
3. Você deve ver traces com múltiplos spans
4. Clique em um trace para ver detalhes
5. Valide que trace_id é propagado para `product-consumer`

## 7. Troubleshooting

### Problema: Collector não recebe telemetria

```bash
# 1. Verificar se Collector está rodando
curl http://localhost:13133/

# 2. Verificar logs do Collector
docker compose -f docker-compose.observability.yml logs otel-collector

# 3. Verificar se serviços estão apontando para endpoint correto
env | grep OTLP
```

### Problema: API ou Consumer não iniciam

```bash
# Ver erro completo
cat api-telemetry.log
cat consumer-telemetry.log

# Verificar dependências
cd services/product-api && pip list | grep opentelemetry
cd services/product-consumer && go list -m all | grep opentelemetry
```

### Problema: Context não propaga via Kafka

**Verificar headers Kafka:**
```bash
# Consumir mensagem e ver headers
docker exec -it kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic products.events \
  --from-beginning \
  --property print.headers=true \
  --max-messages 1
```

Deve mostrar header `traceparent`:
```
traceparent:00-abc123...
{"id":"...","name":"..."}
```

### Problema: Go build errors

```bash
cd services/product-consumer

# Limpar e reinstalar
go clean -modcache
go mod download
go mod tidy
go build -o bin/consumer cmd/main.go
```

## 8. Comandos Úteis

```bash
# Ver últimos 100 spans no Collector
docker compose -f docker-compose.observability.yml logs otel-collector | \
  grep "Span #" | tail -100

# Contar quantos traces foram gerados
docker compose -f docker-compose.observability.yml logs otel-collector | \
  grep "Trace ID" | wc -l

# Ver todas as métricas
docker compose -f docker-compose.observability.yml logs otel-collector | \
  grep "Metric #" -A 10

# Verificar resource attributes
docker compose -f docker-compose.observability.yml logs otel-collector | \
  grep "service.name"

# Limpar tudo
docker compose -f docker-compose.observability.yml down
pkill -f "uvicorn app.main:app"
pkill -f "./bin/consumer"
```

## 9. Checklist de Validação

- [ ] OpenTelemetry Collector iniciou sem erros
- [ ] Product API conectou ao Collector
- [ ] Product Consumer conectou ao Collector
- [ ] Traces aparecem nos logs do Collector
- [ ] Métricas aparecem nos logs do Collector
- [ ] Resource attributes estão corretos (service.name, version, environment)
- [ ] Context propagation funciona (mesmo trace_id entre API e Consumer)
- [ ] Spans têm attributes corretos (product.id, db.table, etc)
- [ ] Jaeger UI mostra traces
- [ ] Não há erros de conexão nos logs

## 10. Próximos Passos

Após validar telemetria com debug exporter:

1. ✅ Configurar ClickHouse
2. ✅ Criar tabelas `otel_traces` e `otel_logs`
3. ✅ Atualizar Collector config com ClickHouse exporter
4. ✅ Validar dados no ClickHouse
5. ✅ Criar queries de análise

## Referências

- [OpenTelemetry Python](https://opentelemetry-python.readthedocs.io/)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
- [OTLP Specification](https://opentelemetry.io/docs/specs/otlp/)
- [Collector Configuration](https://opentelemetry.io/docs/collector/configuration/)
