# Guia de Testes - Product Portfolio

Este guia fornece instruções completas para testar todos os componentes do projeto, desde a infraestrutura até os testes end-to-end.

## 📋 Índice

1. [Pré-requisitos](#pré-requisitos)
2. [Setup Inicial](#setup-inicial)
3. [Testes de Infraestrutura](#testes-de-infraestrutura)
4. [Testes da API (Python)](#testes-da-api-python)
5. [Testes do Consumer (Go)](#testes-do-consumer-go)
6. [Testes End-to-End](#testes-end-to-end)
7. [Testes de Chaos Engineering](#testes-de-chaos-engineering)
8. [Validação Completa](#validação-completa)
9. [Troubleshooting](#troubleshooting)

---

## Pré-requisitos

### Software Necessário

```bash
# Verificar instalações
docker --version          # >= 20.10
docker-compose --version  # >= 1.29
python3 --version         # >= 3.11
go version               # >= 1.21
curl --version
jq --version
```

### Portas Utilizadas

Certifique-se de que estas portas estão livres:

- **5432**: PostgreSQL
- **2181**: Zookeeper
- **9092**: Kafka
- **8000**: Product API
- **8081**: Consumer Health Server

```bash
# Verificar portas em uso (Linux/Mac)
netstat -tuln | grep -E '5432|2181|9092|8000|8081'

# Matar processos se necessário
sudo lsof -ti:8000 | xargs kill -9
sudo lsof -ti:8081 | xargs kill -9
```

---

## Setup Inicial

### 1. Script Automatizado (Recomendado)

```bash
# Executar setup completo
./scripts/setup.sh
```

**O que o script faz:**
- ✅ Verifica pré-requisitos (Docker, Python, Go)
- ✅ Inicia infraestrutura (PostgreSQL, Kafka, Zookeeper)
- ✅ Cria Python virtual environment e instala dependências
- ✅ Cria schema do PostgreSQL
- ✅ Compila o consumer Go
- ✅ Cria arquivo `.env` a partir do template

**Saída esperada:**
```
==================================================
  Setup Complete!
==================================================

Next steps:
  - Run tests: ./scripts/run-tests.sh
  - Start API: cd services/product-api && source venv/bin/activate && uvicorn app.main:app --reload
  - Start Consumer: cd services/product-consumer && ./bin/consumer
  - Run E2E: ./scripts/test-e2e.sh
```

### 2. Setup Manual

Se preferir fazer manualmente:

```bash
# 1. Infraestrutura
docker-compose up -d
docker-compose ps  # Verificar se todos estão UP

# 2. Python API
cd services/product-api
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
pip install -r requirements-dev.txt
cd ../..

# 3. Go Consumer
cd services/product-consumer
go mod download
go build -o bin/consumer cmd/main.go
cd ../..

# 4. Database Schema
docker exec -i postgres psql -U postgres -d products -c "
CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL
);"
```

---

## Testes de Infraestrutura

### Validar Docker Compose

```bash
# Verificar status de todos os serviços
docker-compose ps

# Saída esperada:
# NAME        IMAGE               STATUS          PORTS
# kafka       confluentinc/...    Up (healthy)    0.0.0.0:9092->9092/tcp
# postgres    postgres:15-alpine  Up (healthy)    0.0.0.0:5432->5432/tcp
# zookeeper   confluentinc/...    Up              2181/tcp
```

### Testar PostgreSQL

```bash
# Conectar ao PostgreSQL
docker exec -it postgres psql -U postgres -d products

# Dentro do psql:
\dt                          # Listar tabelas
\d products                  # Descrever tabela products
SELECT COUNT(*) FROM products; # Deve retornar 0 ou mais
\q                           # Sair

# Ou via comando único:
docker exec -i postgres psql -U postgres -d products -c "\dt"
docker exec -i postgres psql -U postgres -d products -c "SELECT * FROM products LIMIT 5;"
```

### Testar Kafka

```bash
# Listar tópicos
docker exec -it kafka kafka-topics --bootstrap-server localhost:9092 --list

# Criar tópico de teste (se não existir)
docker exec -it kafka kafka-topics --bootstrap-server localhost:9092 \
  --create --topic products.events --partitions 3 --replication-factor 1

# Descrever tópico
docker exec -it kafka kafka-topics --bootstrap-server localhost:9092 \
  --describe --topic products.events

# Produzir mensagem de teste
echo '{"test": "message"}' | docker exec -i kafka kafka-console-producer \
  --bootstrap-server localhost:9092 --topic products.events

# Consumir mensagens
docker exec -it kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 --topic products.events \
  --from-beginning --max-messages 1 --timeout-ms 3000
```

**Saída esperada:**
```
{"test": "message"}
Processed a total of 1 messages
```

---

## Testes da API (Python)

### 1. Testes Automatizados

```bash
cd services/product-api

# Ativar ambiente virtual
source venv/bin/activate

# Executar todos os testes
pytest tests/ -v

# Executar com cobertura
pytest tests/ -v --cov=app --cov-report=term

# Executar com cobertura HTML (para visualizar detalhes)
pytest tests/ --cov=app --cov-report=html
open htmlcov/index.html  # No Mac/Linux
# ou xdg-open htmlcov/index.html

# Executar apenas um arquivo de teste
pytest tests/test_routes.py -v

# Executar apenas um teste específico
pytest tests/test_routes.py::TestProductRoutes::test_create_product_missing_name -v
```

**Cobertura esperada:** ≥ 73%

### 2. Iniciar API Manualmente

```bash
cd services/product-api
source venv/bin/activate

# Modo desenvolvimento (com reload)
uvicorn app.main:app --reload --host 0.0.0.0 --port 8000

# Ou em background
uvicorn app.main:app --host 0.0.0.0 --port 8000 > api.log 2>&1 &
```

### 3. Testes Manuais da API

#### Health Checks

```bash
# Health check básico
curl http://localhost:8000/health | jq

# Saída esperada:
# {
#   "status": "healthy"
# }

# Ready check (valida Kafka)
curl http://localhost:8000/ready | jq

# Saída esperada:
# {
#   "status": "ready",
#   "kafka": "connected"
# }
```

#### Criar Produto - Sucesso

```bash
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -H "X-Correlation-ID: test-$(date +%s)" \
  -d '{
    "name": "Laptop Dell XPS 15",
    "price": 1299.99,
    "description": "High-performance laptop"
  }' | jq

# Saída esperada:
# {
#   "id": "550e8400-e29b-41d4-a716-446655440000",
#   "name": "Laptop Dell XPS 15",
#   "price": "1299.99",
#   "description": "High-performance laptop",
#   "created_at": "2026-02-13T10:30:00"
# }
```

#### Validação de Erros

```bash
# Nome muito curto (< 3 caracteres)
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -d '{"name": "AB", "price": 99.99}' | jq

# Preço negativo
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Product", "price": -10}' | jq

# Campo obrigatório faltando
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Product"}' | jq

# Todos devem retornar HTTP 422 com detalhes do erro
```

#### Testar Documentação OpenAPI

```bash
# Abrir no navegador
open http://localhost:8000/docs

# Ou via curl
curl http://localhost:8000/openapi.json | jq
```

### 4. Executar Exemplos Automatizados

```bash
# Script com múltiplos exemplos de requisições
./examples/curl-examples.sh
```

### 5. Usar Postman Collection

```bash
# Importar no Postman
# Arquivo: examples/Product-Portfolio.postman_collection.json

# Ou via Newman (CLI do Postman)
npm install -g newman
newman run examples/Product-Portfolio.postman_collection.json
```

---

## Testes do Consumer (Go)

### 1. Testes Automatizados

```bash
cd services/product-consumer

# Executar todos os testes
go test ./... -v -cover

# Apenas testes unitários (rápido)
go test ./... -v -short -cover

# Teste com detalhamento de cobertura
go test ./internal/consumer/... -v -coverprofile=coverage.out
go tool cover -html=coverage.out  # Abrir no navegador

# Executar apenas consumer tests
go test ./internal/consumer/... -v
```

**Cobertura esperada:** 
- Consumer: ≥ 37%
- Models: ≥ 45%

### 2. Iniciar Consumer Manualmente

```bash
cd services/product-consumer

# Modo normal
./bin/consumer

# Ou recompilar e executar
go build -o bin/consumer cmd/main.go
./bin/consumer

# Em background com log
./bin/consumer > consumer.log 2>&1 &

# Ver logs em tempo real
tail -f consumer.log
```

**Logs esperados:**
```json
{"level":"info","msg":"Starting Product Consumer","app":"product-consumer","version":"1.0.0"}
{"level":"info","msg":"Connected to PostgreSQL","max_conns":25}
{"level":"info","msg":"Kafka consumer started","brokers":"localhost:9092","topic":"products.events"}
{"level":"info","msg":"Health server started","port":8081}
```

### 3. Testar Health Endpoints do Consumer

```bash
# Health check básico
curl http://localhost:8081/health | jq

# Saída esperada:
# {
#   "status": "healthy"
# }

# Ready check (valida PostgreSQL e Kafka)
curl http://localhost:8081/ready | jq

# Saída esperada quando tudo OK:
# {
#   "status": "ready",
#   "postgres": "connected",
#   "kafka": "connected"
# }

# Testar falha do Kafka
docker-compose stop kafka
sleep 2
curl http://localhost:8081/ready | jq
# Deve retornar status "not ready"

# Restartar Kafka
docker-compose start kafka
sleep 10  # Aguardar inicialização
curl http://localhost:8081/ready | jq
# Deve voltar a "ready"
```

### 4. Monitorar Processamento

```bash
# Ver logs de processamento
tail -f consumer.log | grep "product_id"

# Em outro terminal, criar produto via API
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Monitor LG 27",
    "price": 399.99,
    "description": "4K Monitor"
  }' | jq

# Você deve ver no log do consumer:
# {"level":"info","msg":"Processing product message","product_id":"...","product_name":"Monitor LG 27"}
# {"level":"info","msg":"Product persisted successfully","product_id":"..."}
```

---

## Testes End-to-End

### 1. Script Automatizado Completo

```bash
# Executar teste E2E automatizado
./scripts/test-e2e.sh
```

**O que o script faz:**
1. Verifica infraestrutura
2. Inicia API e Consumer
3. Faz health checks
4. Cria produto via API
5. Aguarda processamento
6. Valida persistência no PostgreSQL
7. Mostra resultados

**Saída esperada:**
```
==================================================
  End-to-End Test
==================================================

[1/7] Checking if infrastructure is running...
✓ Infrastructure already running

[2/7] Starting Product API...
✓ API started (PID: 12345)

[3/7] Starting Product Consumer...
✓ Consumer started (PID: 12346)

[4/7] Health checks...
✓ API health: healthy
✓ Consumer ready: ready

[5/7] Creating test product...
Response: {"id":"abc-123","name":"E2E Test Product",...}
✓ Created product with ID: abc-123

[6/7] Waiting for consumer to process message...

[7/7] Verifying in database...
✓ Product found in database
 id  |       name        | price  
-----+-------------------+--------
 abc-123 | E2E Test Product | 149.99

==================================================
  E2E Test Passed!
==================================================
```

### 2. Teste Manual Passo a Passo

```bash
# Passo 1: Garantir que tudo está rodando
docker-compose ps
curl http://localhost:8000/health
curl http://localhost:8081/health

# Passo 2: Criar produto via API
RESPONSE=$(curl -s -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -H "X-Correlation-ID: manual-test-$(date +%s)" \
  -d '{
    "name": "Manual Test Product",
    "price": 99.99,
    "description": "Testing E2E flow manually"
  }')

echo $RESPONSE | jq

# Passo 3: Extrair ID do produto
PRODUCT_ID=$(echo $RESPONSE | jq -r '.id')
echo "Product ID: $PRODUCT_ID"

# Passo 4: Aguardar processamento (3-5 segundos)
sleep 5

# Passo 5: Verificar no banco de dados
docker exec -i postgres psql -U postgres -d products -c \
  "SELECT id, name, price, created_at FROM products WHERE id = '$PRODUCT_ID';"

# Passo 6: Verificar logs
echo "=== API Logs ==="
docker-compose logs --tail=20 product-api 2>/dev/null || echo "API não está no compose"

echo "=== Consumer Logs ==="
tail -20 consumer.log 2>/dev/null || docker-compose logs --tail=20 product-consumer
```

### 3. Teste de Múltiplos Produtos

```bash
# Criar 10 produtos em sequência
for i in {1..10}; do
  echo "Creating product $i..."
  curl -s -X POST http://localhost:8000/products \
    -H "Content-Type: application/json" \
    -H "X-Correlation-ID: batch-$i-$(date +%s)" \
    -d "{
      \"name\": \"Batch Product $i\",
      \"price\": $(echo "scale=2; $i * 10.99" | bc),
      \"description\": \"Product number $i in batch\"
    }" | jq -c '{id: .id, name: .name, price: .price}'
  
  sleep 0.5
done

# Aguardar processamento
echo "Waiting for processing..."
sleep 5

# Verificar quantos foram persistidos
docker exec -i postgres psql -U postgres -d products -c \
  "SELECT COUNT(*) as total FROM products WHERE name LIKE 'Batch Product%';"

# Deve retornar total = 10
```

---

## Testes de Chaos Engineering

### 1. Script Automatizado

```bash
# Executar testes de chaos
./scripts/test-chaos.sh
```

**O que o script faz:**
- Envia 20 requisições com chaos habilitado
- Coleta estatísticas de sucesso/falha
- Analisa latências
- Calcula taxa de erro

**Saída esperada:**
```
Testing Chaos Engineering
Rate: 40% errors, 30% latency

Sending 20 requests with chaos enabled...
Request 1: 202 (394ms)
Request 2: 500 (267ms) FAILED
Request 3: 202 (611ms)
...

Results:
- Success: 12/20 (60%)
- Failures: 8/20 (40%)
- Average latency: 445ms
```

### 2. Teste Manual de Chaos

#### Erro Injection

```bash
# Iniciar API com chaos habilitado
cd services/product-api
source venv/bin/activate

ENABLE_CHAOS=true \
ERROR_RATE=0.5 \
LATENCY_MS_MIN=200 \
LATENCY_MS_MAX=1000 \
uvicorn app.main:app --host 0.0.0.0 --port 8000

# Em outro terminal, fazer requisições
for i in {1..10}; do
  echo "Request $i:"
  curl -s -o /dev/null -w "HTTP %{http_code} - %{time_total}s\n" \
    -X POST http://localhost:8000/products \
    -H "Content-Type: application/json" \
    -d '{"name": "Chaos Test", "price": 99.99}'
  sleep 1
done

# Você deve ver um mix de:
# HTTP 202 - 0.234s  (sucesso)
# HTTP 500 - 0.156s  (erro injetado)
# HTTP 503 - 0.789s  (serviço indisponível)
# HTTP 429 - 0.567s  (rate limit)
```

#### Latency Injection

```bash
# Chaos apenas com latência (sem erros)
ENABLE_CHAOS=true \
ERROR_RATE=0.0 \
LATENCY_MS_MIN=500 \
LATENCY_MS_MAX=2000 \
uvicorn app.main:app --host 0.0.0.0 --port 8000

# Testar latência
time curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Latency Test", "price": 99.99}' | jq

# Deve demorar entre 0.5s e 2s
```

#### Consumer Chaos

```bash
# Consumer com chaos
cd services/product-consumer

ENABLE_CHAOS=true \
ERROR_RATE=0.3 \
LATENCY_MS_MIN=100 \
LATENCY_MS_MAX=500 \
./bin/consumer

# Enviar mensagens e observar falhas/latências nos logs
```

---

## Validação Completa

### Script de Validação Total

```bash
# Executar todos os testes em sequência
./scripts/run-tests.sh
```

**O que valida:**
- ✅ Testes Python (12 testes)
- ✅ Cobertura Python ≥ 73%
- ✅ Testes Go consumer (3 testes)
- ✅ Cobertura Go ≥ 37%

### Checklist Manual de Validação

```bash
# 1. Infraestrutura
[ ] Docker Compose com todos os serviços UP
[ ] PostgreSQL acessível na porta 5432
[ ] Kafka acessível na porta 9092
[ ] Tabela products criada

# 2. API Python
[ ] Testes unitários passando (12/12)
[ ] Cobertura ≥ 73%
[ ] Health check retorna 200
[ ] POST /products cria produto
[ ] Validação de campos funciona
[ ] Kafka producer envia mensagem

# 3. Consumer Go
[ ] Testes unitários passando (3/3)
[ ] Compila sem erros
[ ] Health check retorna 200 (porta 8081)
[ ] Ready check valida PostgreSQL e Kafka
[ ] Processa mensagens do Kafka
[ ] Persiste no PostgreSQL

# 4. End-to-End
[ ] Fluxo API → Kafka → Consumer → DB funciona
[ ] Produto aparece no banco após criação
[ ] Logs mostram processamento

# 5. Chaos Engineering
[ ] Erros são injetados quando habilitado
[ ] Latência é adicionada quando configurado
[ ] Health checks continuam funcionando
[ ] Aplicação se recupera após erros

# 6. Documentação
[ ] README.md completo e atualizado
[ ] Postman collection importa corretamente
[ ] curl-examples.sh executa sem erros
[ ] Service READMEs estão completos
```

---

## Troubleshooting

### Problema: Docker Compose não sobe

```bash
# Verificar logs
docker-compose logs

# Verificar portas em uso
netstat -tuln | grep -E '5432|9092|2181'

# Limpar tudo e recomeçar
docker-compose down -v
docker-compose up -d
```

### Problema: Kafka não conecta

```bash
# Verificar se Kafka está pronto
docker-compose logs kafka | grep "started (kafka.server.KafkaServer)"

# Aguardar mais tempo (Kafka pode demorar 30-60s)
sleep 30

# Testar conectividade
docker exec -it kafka kafka-broker-api-versions \
  --bootstrap-server localhost:9092
```

### Problema: PostgreSQL não conecta

```bash
# Verificar se está rodando
docker-compose ps postgres

# Verificar logs
docker-compose logs postgres

# Testar conexão
docker exec -it postgres psql -U postgres -c "SELECT version();"

# Recriar database
docker-compose down
docker volume rm otel-clickhouse-analyze_postgres_data
docker-compose up -d postgres
```

### Problema: Testes Python falhando

```bash
# Reinstalar dependências
cd services/product-api
rm -rf venv
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
pip install -r requirements-dev.txt

# Executar testes com mais detalhes
pytest tests/ -vv -s
```

### Problema: Consumer não processa mensagens

```bash
# Verificar consumer group
docker exec -it kafka kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --describe --group product-consumer-group

# Verificar tópico tem mensagens
docker exec -it kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic products.events \
  --from-beginning --max-messages 5

# Verificar logs do consumer
tail -100 consumer.log | grep -E "(ERROR|WARN|Processing)"
```

### Problema: Cobertura de testes baixa

```bash
# Ver arquivos não cobertos
pytest tests/ --cov=app --cov-report=term-missing

# Gerar relatório HTML para análise
pytest tests/ --cov=app --cov-report=html
open htmlcov/index.html
```

---

## Próximos Passos

Após validar todos os testes:

1. **Commit das mudanças**
   ```bash
   git add .
   git commit -m "test: complete Phase 8 - testing and documentation"
   git push
   ```

2. **Preparar para Fase 9**
   - Revisar DEVELOPMENT_PLAN.md - Fase 9
   - OpenTelemetry Collector
   - Instrumentação de traces
   - Métricas customizadas

3. **Melhorias opcionais**
   - Aumentar cobertura de testes
   - Adicionar mais testes de integração
   - Implementar testes de carga (K6, Locust)
   - Configurar CI/CD pipeline

---

## 📊 Métricas de Qualidade

### Objetivos Atingidos

- ✅ Cobertura Python: **73%** (meta: ≥70%)
- ✅ Testes Go passando: **100%**
- ✅ Infraestrutura estável
- ✅ Documentação completa
- ✅ Scripts de automação funcionais
- ✅ Chaos engineering validado

### Comandos Rápidos de Validação

```bash
# Validação completa em 1 comando
./scripts/run-tests.sh && ./scripts/test-e2e.sh && echo "✅ ALL TESTS PASSED"

# Status rápido
docker-compose ps && \
curl -s http://localhost:8000/health | jq && \
curl -s http://localhost:8081/health | jq
```

---

**🎉 Projeto testado e validado! Pronto para produção (ou próximas fases de observabilidade).**
