# Microservices Observability Portfolio

Platform de observabilidade demonstrando práticas modernas com OpenTelemetry, ClickHouse, Prometheus, Grafana e SRE practices.

## Arquitetura

```
┌─────────────┐      ┌──────────┐      ┌──────────────┐      ┌────────────┐
│ Product API │ ---> │  Kafka   │ ---> │   Consumer   │ ---> │ PostgreSQL │
│  (Python)   │      │          │      │     (Go)     │      │            │
└─────────────┘      └──────────┘      └──────────────┘      └────────────┘
                           │
                           v
                     ┌──────────┐
                     │   DLQ    │
                     └──────────┘
```

## Stack Tecnológica

- **API**: Python FastAPI + Pydantic
- **Consumer**: Golang + kafka-go + pgx
- **Message Broker**: Apache Kafka (Confluent)
- **Database**: PostgreSQL 15
- **Observability**: OpenTelemetry, Prometheus, Grafana
- **Storage**: ClickHouse (traces/logs)

## Quick Start

```bash
# Clone repository
git clone <repo-url>
cd otel-clickhouse-analyze

# Start infrastructure
docker compose up -d

# Check services
docker compose ps

# Create product via API
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Product", "price": 99.99, "description": "Test"}'

# Verify in database
docker compose exec postgres psql -U postgres -d products_db \
  -c "SELECT * FROM products ORDER BY created_at DESC LIMIT 1;"
```

## Health Checks

### API Health
```bash
curl http://localhost:8000/health
# Returns: {"status": "healthy"}

curl http://localhost:8000/ready
# Returns: {"status": "ready", "kafka": "connected"}
```

### Consumer Health
```bash
curl http://localhost:8081/health
# Returns: {"status": "healthy"}

curl http://localhost:8081/ready
# Returns: {"status": "ready", "postgres": "connected", "kafka": "connected"}
```

## Chaos Engineering

### Habilitar Chaos na API

Edite `docker-compose.yml` ou defina variáveis de ambiente:

```yaml
environment:
  ENABLE_CHAOS: "true"
  ERROR_RATE: "0.3"           # 30% de erros
  LATENCY_MS_MIN: "500"       # Latência mínima 500ms
  LATENCY_MS_MAX: "3000"      # Latência máxima 3s
```

### Habilitar Chaos no Consumer

```yaml
environment:
  ENABLE_CHAOS: "true"
  CHAOS_ERROR_RATE: "0.2"         # 20% de erros
  CHAOS_LATENCY_MIN_MS: "1000"    # Latência mínima 1s
  CHAOS_LATENCY_MAX_MS: "5000"    # Latência máxima 5s
```

### Tipos de Erros Injetados (API)

- **500 Internal Server Error** (50% dos erros)
- **503 Service Unavailable** (30% dos erros)
- **429 Too Many Requests** (20% dos erros)

### Testar Chaos

```bash
# Enviar múltiplas requests com chaos habilitado
for i in {1..20}; do
  echo "Request $i:"
  time curl -X POST http://localhost:8000/products \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"Product $i\",\"price\":100.00}"
  echo ""
done

# Observar logs com eventos de chaos
docker compose logs product-api | grep -i "chaos"
docker compose logs product-consumer | grep -i "chaos"
```

### Importante

- Chaos está **desabilitado por padrão** (production-safe)
- Health checks (`/health`) **não são afetados** por chaos
- Mensagens com erro no consumer **não são commitadas** (serão reprocessadas)
- Idempotência no DB garante que reprocessamento não cria duplicatas

## Variáveis de Ambiente

Veja [.env.example](.env.example) para todas as configurações disponíveis.

## Desenvolvimento Local

### Pré-requisitos

- Docker & Docker Compose
- Python 3.11+ (para desenvolvimento da API)
- Go 1.21+ (para desenvolvimento do consumer)

### API (Python)

```bash
cd services/product-api
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
uvicorn app.main:app --reload
```

### Consumer (Go)

```bash
cd services/product-consumer
go mod download
go build -o bin/consumer cmd/main.go
./bin/consumer
```

## Testes

### End-to-End

```bash
./scripts/test-e2e.sh
```

### Unit Tests

```bash
# API
cd services/product-api
pytest

# Consumer
cd services/product-consumer
go test ./... -v
```

## Estrutura do Projeto

```
├── docker-compose.yml
├── infrastructure/
│   ├── kafka/
│   ├── postgres/
│   │   └── init.sql
│   └── observability/
├── scripts/
│   └── test-e2e.sh
├── services/
│   ├── product-api/
│   │   ├── app/
│   │   │   ├── api/
│   │   │   ├── middleware/
│   │   │   ├── models/
│   │   │   ├── services/
│   │   │   ├── config.py
│   │   │   └── main.py
│   │   ├── Dockerfile
│   │   └── requirements.txt
│   └── product-consumer/
│       ├── cmd/
│       ├── internal/
│       │   ├── config/
│       │   ├── consumer/
│       │   ├── health/
│       │   ├── models/
│       │   └── repository/
│       ├── Dockerfile
│       └── go.mod
└── sre/
    ├── slos/
    ├── runbooks/
    └── postmortems/
```

## Fases Implementadas

- ✅ Fase 1: Infrastructure Setup (Docker Compose)
- ✅ Fase 2: Product API (FastAPI + Kafka Producer)
- ✅ Fase 3: Kafka Integration
- ✅ Fase 4: Go Consumer Base Structure
- ✅ Fase 5: PostgreSQL Repository Pattern
- ✅ Fase 6: Complete Kafka Consumer
- ✅ Fase 7: Chaos Engineering
- ⏳ Fase 8: Testing & Documentation
- ⏳ Fase 9-12: OpenTelemetry, ClickHouse, Prometheus, Grafana

## License

MIT
