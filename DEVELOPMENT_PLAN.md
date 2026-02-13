# Plano de Desenvolvimento - Microserviços para Estudo de Observabilidade

## 📋 Visão Geral do Projeto

**Projeto de Portfólio** para demonstrar competências em **Observability Engineering** e **SRE practices**, alinhado com as tecnologias e práticas da Wolt/DoorDash.

Este projeto implementa uma arquitetura de microserviços com **stack de observabilidade completa** para processar telemetria em escala (métricas, traces, logs), utilizando **OpenTelemetry**, **ClickHouse**, **Prometheus**, **Grafana** e práticas de **Site Reliability Engineering**.

### 🎯 Objetivo do Portfólio

Demonstrar habilidades práticas em:
- ✅ Desenvolvimento de sistemas distribuídos escaláveis
- ✅ Instrumentação com OpenTelemetry (traces, metrics, logs)
- ✅ Arquitetura e manutenção de stack de observabilidade
- ✅ Go e Python para automação e tooling
- ✅ ClickHouse para análise de telemetria em escala
- ✅ Prometheus + Grafana para métricas e visualização
- ✅ SRE practices (SLIs/SLOs, incident response)
- ✅ Kubernetes deployment e cloud-native patterns
- ✅ Kafka para event streaming

### Arquitetura Proposta

```
                                    Observability Stack
                    ┌─────────────────────────────────────────────┐
                    │                                             │
                    │  ┌──────────────┐      ┌──────────────┐   │
                    │  │ OpenTelemetry│─────▶│  ClickHouse  │   │
                    │  │  Collector   │      │   (Traces,   │   │
                    │  │              │      │   Logs)      │   │
                    │  └──────────────┘      └──────────────┘   │
                    │         │                      │           │
                    │         │                      │           │
                    │         ▼                      ▼           │
                    │  ┌──────────────┐      ┌──────────────┐   │
                    │  │  Prometheus  │      │   Grafana    │   │
                    │  │  (Metrics)   │      │ (Dashboards) │   │
                    │  └──────────────┘      └──────────────┘   │
                    └─────────────────────────────────────────────┘
                                    ▲
                                    │ (OTLP)
                                    │
┌─────────────────┐      ┌─────────────┐      ┌──────────────────┐      ┌────────────┐
│   Product API   │─────▶│    Kafka    │─────▶│ Product Consumer │─────▶│ PostgreSQL │
│    (Python)     │      │             │      │     (Golang)     │      │            │
│   + OTel SDK    │      │             │      │   + OTel SDK     │      │            │
└─────────────────┘      └─────────────┘      └──────────────────┘      └────────────┘
   FastAPI +               Topic: product          Segmentio +              DB
   Instrumentation         -registrations          kafka-go +
                                                   Instrumentation
```

### Tecnologias Utilizadas

#### Microserviços
- **Python**: FastAPI, Kafka-Python, Pydantic, OpenTelemetry SDK
- **Golang**: Segmentio Kafka-Go, pgx (PostgreSQL driver), OpenTelemetry SDK
- **Infraestrutura**: Docker Compose, Kafka, Zookeeper, PostgreSQL

#### Observability Stack (Alinhado com Wolt)
- **OpenTelemetry**: Instrumentação automática e manual, OTLP exporters
- **ClickHouse**: Storage de traces e logs em escala
- **Prometheus**: Métricas de aplicação e infraestrutura
- **Grafana**: Dashboards, visualização e alerting
- **Jaeger** (opcional): UI para visualização de traces

#### SRE & DevOps
- **Kubernetes**: Deployment e orquestração (opcional, fase avançada)
- **Docker**: Containerização
- **SLIs/SLOs**: Service Level Indicators e Objectives
- **Incident Response**: Alerting, runbooks, postmortems

### Estrutura de Diretórios

```
microservices-observability-study/
├── services/
│   ├── product-api/              # API Python
│   │   ├── app/
│   │   │   ├── __init__.py
│   │   │   ├── main.py
│   │   │   ├── api/
│   │   │   │   ├── __init__.py
│   │   │   │   └── routes.py
│   │   │   ├── models/
│   │   │   │   ├── __init__.py
│   │   │   │   └── product.py
│   │   │   ├── services/
│   │   │   │   ├── __init__.py
│   │   │   │   └── kafka_producer.py
│   │   │   ├── middleware/
│   │   │   │   ├── __init__.py
│   │   │   │   ├── error_injection.py
│   │   │   │   └── latency_injection.py
│   │   │   └── config.py
│   │   ├── tests/
│   │   ├── requirements.txt
│   │   ├── Dockerfile
│   │   └── README.md
│   └── product-consumer/         # Consumer Golang
│       ├── cmd/
│       │   └── main.go
│       ├── internal/
│       │   ├── consumer/
│       │   │   └── message_handler.go
│       │   ├── repository/
│       │   │   └── product_repository.go
│       │   ├── models/
│       │   │   └── product.go
│       │   └── config/
│       │       └── config.go
│       ├── tests/
│       ├── go.mod
│       ├── go.sum
│       ├── Dockerfile
│       └── README.md
├── infrastructure/
│   ├── postgres/
│   │   └── init.sql
│   ├── kafka/
│   │   └── topics.txt
│   └── observability/
│       ├── otel-collector/
│       │   └── otel-config.yaml
│       ├── clickhouse/
│       │   ├── init-db.sql
│       │   └── config.xml
│       ├── prometheus/
│       │   └── prometheus.yml
│       └── grafana/
│           ├── dashboards/
│           │   ├── microservices-overview.json
│           │   ├── traces-analysis.json
│           │   └── slo-dashboard.json
│           └── provisioning/
│               ├── datasources.yml
│               └── dashboards.yml
├── k8s/                          # Kubernetes manifests (opcional)
│   ├── deployments/
│   ├── services/
│   └── configmaps/
├── sre/
│   ├── slos/
│   │   └── product-service-slo.yaml
│   ├── runbooks/
│   │   └── incident-response.md
│   └── postmortems/
│       └── template.md
├── docker-compose.yml
├── docker-compose.observability.yml
├── .env.example
├── .gitignore
└── README.md
```

---

## 🎯 Fase 1: Setup Inicial da Infraestrutura

### Objetivo
Configurar Docker Compose com todos os serviços de infraestrutura necessários (Kafka, Zookeeper, PostgreSQL).

### Descrição
Criar a base de infraestrutura que suportará os microserviços. Todos os serviços devem ter health checks, redes isoladas e configurações adequadas para desenvolvimento.

### Prompt para Implementação

```
Preciso criar a infraestrutura base para um projeto de microserviços usando Docker Compose.

Criar os seguintes arquivos:

1. **docker-compose.yml** com:
   - Zookeeper (Confluent)
   - Kafka (Confluent) com configurações:
     * KAFKA_ADVERTISED_LISTENERS para localhost e rede interna
     * Auto-criação de tópicos habilitada
     * Replication factor adequado para desenvolvimento
   - PostgreSQL 15:
     * Database: products_db
     * User/Password configuráveis via .env
     * Volume persistente
     * Health check
   - Redes:
     * kafka-network (Kafka e serviços)
     * db-network (PostgreSQL e consumer)

2. **infrastructure/postgres/init.sql**:
   - CREATE TABLE products com campos:
     * id UUID PRIMARY KEY
     * name VARCHAR(255) NOT NULL
     * price DECIMAL(10,2) NOT NULL
     * description TEXT
     * created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   - Índices apropriados

3. **.env.example**:
   - Todas as variáveis de ambiente necessárias
   - Documentação inline

4. **.gitignore**:
   - Python (__pycache__, .pytest_cache, venv)
   - Go (vendor/, *.exe)
   - Docker (.env, volumes/)
   - IDE (.vscode/, .idea/)

Requisitos:
- Usar versões estáveis e recentes
- Health checks em todos os serviços
- Restart policies adequados
- Volumes nomeados para persistência
- Configurações otimizadas para desenvolvimento local
```

### Critérios de Aceite

- [ ] Docker Compose sobe todos os serviços sem erros
- [ ] Kafka está acessível na porta 9092
- [ ] PostgreSQL está acessível na porta 5432
- [ ] Health checks estão funcionando
- [ ] Tabela products é criada automaticamente
- [ ] .env.example documenta todas as variáveis

### Comandos de Teste

```bash
# Subir infraestrutura
docker-compose up -d zookeeper kafka postgres

# Verificar status
docker-compose ps

# Verificar logs
docker-compose logs kafka
docker-compose logs postgres

# Testar conexão PostgreSQL
docker-compose exec postgres psql -U postgres -d products_db -c "\dt"

# Testar Kafka
docker-compose exec kafka kafka-topics --list --bootstrap-server localhost:9092
```

### Entregáveis
- `docker-compose.yml`
- `infrastructure/postgres/init.sql`
- `.env.example`
- `.gitignore`

---

## 🐍 Fase 2: Product API (Python) - Estrutura Básica

### Objetivo
Criar API REST com FastAPI para cadastro de produtos, seguindo clean architecture e boas práticas Python.

### Descrição
Implementar endpoint POST /products com validação de dados usando Pydantic, estrutura modular e preparação para integração com Kafka.

### Prompt para Implementação

```
Criar a API REST em Python usando FastAPI para cadastro de produtos.

Estrutura a implementar:

1. **app/models/product.py**:
   - ProductCreate (Pydantic model) com validações:
     * name: string, 3-255 caracteres, não vazio
     * price: decimal, maior que 0, máximo 2 casas decimais
     * description: string opcional, máximo 1000 caracteres
   - ProductResponse com id UUID e created_at

2. **app/config.py**:
   - Classe Settings usando pydantic BaseSettings
   - Variáveis de ambiente:
     * APP_NAME, APP_VERSION
     * KAFKA_BOOTSTRAP_SERVERS
     * LOG_LEVEL
   - Singleton pattern para config

3. **app/api/routes.py**:
   - Router FastAPI
   - POST /products endpoint:
     * Recebe ProductCreate
     * Gera UUID
     * Valida dados
     * Retorna ProductResponse (mock por enquanto)
   - GET /health endpoint para health check

4. **app/main.py**:
   - Criar aplicação FastAPI
   - Configurar CORS
   - Incluir routers
   - Configurar logging estruturado (JSON)
   - Lifespan events (startup/shutdown)

5. **requirements.txt**:
   - fastapi
   - uvicorn[standard]
   - pydantic
   - pydantic-settings
   - python-json-logger
   - kafka-python (para próxima fase)

6. **Dockerfile**:
   - Multi-stage build
   - Python 3.11-slim
   - Non-root user
   - Otimizado para layer caching
   - Health check

7. **README.md**:
   - Descrição do serviço
   - Tecnologias usadas
   - Como rodar localmente
   - Endpoints disponíveis
   - Exemplos de requisições

Requisitos:
- Seguir PEP 8 e type hints
- Documentação automática (OpenAPI/Swagger)
- Tratamento de erros com status codes adequados
- Logs estruturados em JSON
- Código testável e desacoplado
```

### Critérios de Aceite

- [ ] API inicia sem erros
- [ ] Endpoint POST /products valida dados corretamente
- [ ] Endpoint GET /health retorna 200
- [ ] Swagger UI acessível em /docs
- [ ] Validações Pydantic funcionam
- [ ] Logs estruturados em JSON
- [ ] Dockerfile build com sucesso

### Comandos de Teste

```bash
# Build da imagem
cd services/product-api
docker build -t product-api:latest .

# Rodar localmente (sem Docker Compose)
pip install -r requirements.txt
uvicorn app.main:app --reload --host 0.0.0.0 --port 8000

# Testar endpoint
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Notebook Dell",
    "price": 3500.00,
    "description": "Notebook para desenvolvimento"
  }'

# Testar validação (preço negativo)
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Teste",
    "price": -100,
    "description": "Deve falhar"
  }'

# Health check
curl http://localhost:8000/health

# Acessar documentação
open http://localhost:8000/docs
```

### Entregáveis
- `services/product-api/app/` (todos os módulos)
- `services/product-api/requirements.txt`
- `services/product-api/Dockerfile`
- `services/product-api/README.md`

---

## 📨 Fase 3: Product API - Integração Kafka Producer

### Objetivo
Integrar producer Kafka na API Python para publicar mensagens no tópico "product-registrations".

### Descrição
Implementar serviço de producer Kafka desacoplado, com tratamento de erros, retry e logs estruturados.

### Prompt para Implementação

```
Integrar Kafka Producer na API Python existente.

Implementar:

1. **app/services/kafka_producer.py**:
   - Classe KafkaProducerService:
     * Inicialização com configurações do Kafka
     * Método send_product_event(product: dict):
       - Serializa para JSON
       - Envia para tópico "product-registrations"
       - Retry com backoff exponencial (3 tentativas)
       - Log de sucesso/falha
     * Método close() para shutdown gracioso
     * Connection pooling
   - Tratamento de exceções específicas do Kafka
   - Callback de confirmação de entrega

2. **app/api/routes.py** (modificar):
   - Injetar dependency do KafkaProducerService
   - No POST /products:
     * Enviar mensagem para Kafka após validação
     * Tratar erros de envio
     * Retornar 202 Accepted em caso de sucesso
     * Retornar 500 se Kafka falhar após retries
   - Adicionar correlation_id nos logs

3. **app/main.py** (modificar):
   - Inicializar KafkaProducerService no startup
   - Fechar producer no shutdown
   - Singleton pattern para producer

4. **app/config.py** (modificar):
   - Adicionar configurações:
     * KAFKA_TOPIC_PRODUCTS
     * KAFKA_RETRIES
     * KAFKA_RETRY_BACKOFF_MS
     * KAFKA_ACKS (all, 1, 0)

5. **Atualizar docker-compose.yml**:
   - Adicionar serviço product-api
   - Depends_on: kafka
   - Environment variables
   - Conectar às redes kafka-network
   - Port mapping 8000:8000

Requisitos:
- Producer deve ser thread-safe
- Mensagens em formato JSON estruturado:
  ```json
  {
    "id": "uuid",
    "name": "string",
    "price": float,
    "description": "string",
    "timestamp": "ISO8601"
  }
  ```
- Logs devem incluir correlation_id
- Tratamento de erros sem perder dados
- Configuração de acks=all para garantias de entrega
```

### Critérios de Aceite

- [ ] API envia mensagens para Kafka com sucesso
- [ ] Mensagens aparecem no tópico "product-registrations"
- [ ] Retry funciona em caso de falha temporária
- [ ] Logs incluem correlation_id
- [ ] Erros são tratados adequadamente
- [ ] Producer é fechado graciosamente no shutdown
- [ ] Container sobe via docker-compose

### Comandos de Teste

```bash
# Subir serviços
docker-compose up -d

# Criar tópico manualmente (se necessário)
docker-compose exec kafka kafka-topics --create \
  --topic product-registrations \
  --bootstrap-server localhost:9092 \
  --partitions 3 \
  --replication-factor 1

# Enviar produto via API
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mouse Gamer",
    "price": 150.00,
    "description": "RGB e DPI ajustável"
  }'

# Consumir mensagens do tópico para validar
docker-compose exec kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic product-registrations \
  --from-beginning

# Verificar logs da API
docker-compose logs -f product-api

# Testar com Kafka offline (simular falha)
docker-compose stop kafka
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Teste", "price": 10.00}'
# Deve retornar erro após retries
```

### Entregáveis
- `services/product-api/app/services/kafka_producer.py`
- `services/product-api/app/api/routes.py` (atualizado)
- `services/product-api/app/main.py` (atualizado)
- `docker-compose.yml` (atualizado com product-api)

---

## 🐹 Fase 4: Product Consumer (Golang) - Estrutura Básica

### Objetivo
Criar estrutura base do projeto Go seguindo Go project layout e clean architecture.

### Descrição
Implementar estrutura de pacotes, models, configuração e preparação para repository pattern.

### Prompt para Implementação

```
Criar a estrutura base do consumer Golang com organização clean e idiomática.

Implementar:

1. **go.mod**:
   - Module name: github.com/seu-usuario/product-consumer
   - Go version: 1.21+
   - Dependências:
     * github.com/segmentio/kafka-go
     * github.com/jackc/pgx/v5
     * github.com/joho/godotenv
     * github.com/sirupsen/logrus (structured logging)

2. **internal/models/product.go**:
   - Struct Product com tags JSON e DB:
     ```go
     type Product struct {
         ID          string    `json:"id" db:"id"`
         Name        string    `json:"name" db:"name"`
         Price       float64   `json:"price" db:"price"`
         Description string    `json:"description" db:"description"`
         CreatedAt   time.Time `json:"created_at" db:"created_at"`
     }
     ```
   - Método Validate() error

3. **internal/config/config.go**:
   - Struct Config com todas as configurações:
     * Kafka (brokers, topic, group_id)
     * PostgreSQL (host, port, user, password, database)
     * Application (log_level, shutdown_timeout)
   - Função LoadConfig() que lê de variáveis de ambiente
   - Validação de configurações obrigatórias
   - Uso de godotenv para .env

4. **cmd/main.go**:
   - Estrutura básica com:
     * Load config
     * Setup logger
     * Context com cancel
     * Signal handling (SIGTERM, SIGINT)
     * Graceful shutdown
     * Placeholder para consumer

5. **Dockerfile**:
   - Multi-stage build:
     * Builder: golang:1.21-alpine
     * Runtime: alpine:latest
   - CGO_ENABLED=0 para static binary
   - Non-root user
   - Health check (preparar endpoint HTTP)

6. **README.md**:
   - Descrição do serviço
   - Arquitetura do código
   - Variáveis de ambiente
   - Como buildar e rodar
   - Comandos úteis

Requisitos:
- Seguir Go best practices e idioms
- Error handling adequado
- Structured logging desde o início
- Código preparado para testes
- Documentação em GoDoc format
```

### Critérios de Aceite

- [ ] Projeto Go compila sem erros
- [ ] Configurações são carregadas corretamente
- [ ] Logger estruturado funciona
- [ ] Graceful shutdown está implementado
- [ ] Dockerfile gera imagem otimizada
- [ ] go.mod tem todas as dependências
- [ ] Código segue Go conventions

### Comandos de Teste

```bash
cd services/product-consumer

# Inicializar módulo
go mod init github.com/seu-usuario/product-consumer

# Adicionar dependências
go get github.com/segmentio/kafka-go
go get github.com/jackc/pgx/v5
go get github.com/joho/godotenv
go get github.com/sirupsen/logrus

# Tidy dependencies
go mod tidy

# Build local
go build -o bin/consumer cmd/main.go

# Rodar localmente
cp ../../.env.example .env
# Editar .env com configurações
./bin/consumer

# Testar com variáveis de ambiente
KAFKA_BROKERS=localhost:9092 \
KAFKA_TOPIC=product-registrations \
DB_HOST=localhost \
DB_PORT=5432 \
./bin/consumer

# Build Docker
docker build -t product-consumer:latest .

# Lint (opcional)
go install golang.org/x/lint/golint@latest
golint ./...

# Vet
go vet ./...
```

### Entregáveis
- `services/product-consumer/go.mod`
- `services/product-consumer/go.sum`
- `services/product-consumer/internal/models/product.go`
- `services/product-consumer/internal/config/config.go`
- `services/product-consumer/cmd/main.go`
- `services/product-consumer/Dockerfile`
- `services/product-consumer/README.md`

---

## 🗄️ Fase 5: Product Consumer - Repository PostgreSQL

### Objetivo
Implementar repository pattern para persistência de produtos no PostgreSQL usando pgx.

### Descrição
Criar camada de repository com connection pool, prepared statements, transações e tratamento de duplicatas.

### Prompt para Implementação

```
Implementar repository pattern para PostgreSQL no consumer Golang.

Implementar:

1. **internal/repository/product_repository.go**:
   - Interface ProductRepository:
     ```go
     type ProductRepository interface {
         CreateProduct(ctx context.Context, product *models.Product) error
         GetProductByID(ctx context.Context, id string) (*models.Product, error)
         Close() error
     }
     ```
   
   - Struct productRepository com pgxpool.Pool
   
   - Função NewProductRepository(config *config.Config) (ProductRepository, error):
     * Criar connection pool
     * Configurar pool size, timeouts
     * Testar conexão com Ping
   
   - Implementar CreateProduct:
     * INSERT com ON CONFLICT DO NOTHING (idempotência)
     * Usar prepared statements
     * Tratamento de constraint violations
     * Context timeout
     * Logs estruturados
   
   - Implementar GetProductByID (para testes)
   
   - Método Close para fechar pool

2. **internal/repository/repository_test.go**:
   - Testes unitários básicos
   - Mock do database (opcional)
   - Integration tests com testcontainers (opcional)

3. **cmd/main.go** (atualizar):
   - Inicializar repository no startup
   - Passar para consumer (preparar)
   - Fechar repository no shutdown
   - Tratamento de erros de conexão

4. **infrastructure/postgres/init.sql** (atualizar se necessário):
   - Adicionar índices:
     * CREATE INDEX idx_products_created_at ON products(created_at)
   - Garantir UUID como PK

Requisitos:
- Connection pool otimizado para produção
- Context-aware em todas as queries
- Prepared statements para performance
- Idempotência (duplicatas não causam erro)
- Transações quando necessário
- Error wrapping com contexto
- Logs de todas as operações
```

### Critérios de Aceite

- [ ] Repository conecta no PostgreSQL com sucesso
- [ ] CreateProduct insere dados corretamente
- [ ] Duplicatas são tratadas (idempotência)
- [ ] Connection pool funciona
- [ ] Context timeout é respeitado
- [ ] Logs incluem operações do DB
- [ ] Repository é fechado graciosamente

### Comandos de Teste

```bash
# Rodar consumer com repository
cd services/product-consumer
go build -o bin/consumer cmd/main.go

# Configurar .env
cat > .env << EOF
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=products_db
DB_MAX_CONNS=10
DB_MIN_CONNS=2
EOF

# Garantir que PostgreSQL está rodando
docker-compose up -d postgres

# Rodar consumer
./bin/consumer

# Em outro terminal, testar inserção direta
docker-compose exec postgres psql -U postgres -d products_db

# Inserir produto manualmente para testar
INSERT INTO products (id, name, price, description)
VALUES ('123e4567-e89b-12d3-a456-426614174000', 'Produto Teste', 99.90, 'Descrição');

# Verificar logs do consumer

# Testar idempotência (inserir duplicata)
INSERT INTO products (id, name, price, description)
VALUES ('123e4567-e89b-12d3-a456-426614174000', 'Produto Teste', 99.90, 'Descrição');
# Não deve dar erro

# Verificar pool de conexões
SELECT * FROM pg_stat_activity WHERE datname = 'products_db';
```

### Entregáveis
- `services/product-consumer/internal/repository/product_repository.go`
- `services/product-consumer/cmd/main.go` (atualizado)
- `infrastructure/postgres/init.sql` (atualizado)

---

## 🔄 Fase 6: Product Consumer - Kafka Consumer Integration

### Objetivo
Integrar consumer Kafka completo com repository, implementando processamento robusto de mensagens.

### Descrição
Consumir mensagens do Kafka, deserializar, persistir no PostgreSQL com retry, DLQ e garantias de entrega.

### Prompt para Implementação

```
Implementar consumer Kafka completo integrado com repository PostgreSQL.

Implementar:

1. **internal/consumer/message_handler.go**:
   - Struct MessageHandler:
     * Dependencies: ProductRepository, Logger, Config
   
   - Método ProcessMessage(ctx context.Context, msg []byte) error:
     * Deserializar JSON para Product
     * Validar produto
     * Chamar repository.CreateProduct
     * Retry com backoff exponencial (3 tentativas)
     * Retornar erro se falhar após retries
   
   - Método handleDeserializationError(msg []byte):
     * Log de mensagens inválidas
     * Enviar para DLQ (tópico: product-registrations-dlq)
   
   - Função NewMessageHandler(repo, config) *MessageHandler

2. **internal/consumer/kafka_consumer.go**:
   - Struct KafkaConsumer:
     * kafka.Reader
     * MessageHandler
     * Context
   
   - Função NewKafkaConsumer(config, handler) *KafkaConsumer:
     * Configurar Reader com:
       - Brokers, Topic, GroupID
       - MinBytes, MaxBytes
       - CommitInterval (manual commit)
       - StartOffset = kafka.LastOffset
   
   - Método Start():
     * Loop de leitura de mensagens
     * Para cada mensagem:
       - Processar com handler.ProcessMessage
       - Se sucesso: CommitMessages
       - Se erro: Log e continue (não commita)
     * Graceful shutdown quando context cancelado
   
   - Método Close():
     * Fechar reader
     * Aguardar processamento atual

3. **cmd/main.go** (atualizar):
   - Inicializar todos os componentes:
     * Config
     * Logger
     * Repository
     * MessageHandler
     * KafkaConsumer
   
   - Iniciar consumer em goroutine
   
   - Signal handling:
     * Capturar SIGTERM/SIGINT
     * Cancel context
     * Aguardar shutdown gracioso (timeout 30s)
     * Fechar resources
   
   - Health check HTTP endpoint (opcional):
     * GET /health
     * GET /ready (verifica DB e Kafka)

4. **docker-compose.yml** (atualizar):
   - Adicionar serviço product-consumer:
     * depends_on: kafka, postgres
     * Environment variables
     * Redes: kafka-network, db-network
     * Restart: unless-stopped

Requisitos:
- At-least-once delivery garantido
- Commit manual apenas após sucesso no DB
- Idempotência no processamento
- Graceful shutdown sem perda de mensagens
- Logs estruturados com correlation_id
- Métricas básicas (contador de mensagens)
- Dead Letter Queue para mensagens inválidas
- Context propagation em toda a stack
```

### Critérios de Aceite

- [ ] Consumer conecta e lê mensagens do Kafka
- [ ] Mensagens são deserializadas corretamente
- [ ] Produtos são salvos no PostgreSQL
- [ ] Commit ocorre apenas após sucesso
- [ ] Retry funciona em caso de erro temporário
- [ ] DLQ recebe mensagens inválidas
- [ ] Graceful shutdown funciona
- [ ] Logs são estruturados e informativos
- [ ] Fluxo end-to-end funciona completamente

### Comandos de Teste

```bash
# Subir toda a stack
docker-compose up -d

# Verificar todos os serviços
docker-compose ps

# Acompanhar logs do consumer
docker-compose logs -f product-consumer

# Enviar produto via API
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Teclado Mecânico",
    "price": 450.00,
    "description": "Switch blue, RGB"
  }'

# Verificar logs da API
docker-compose logs -f product-api

# Verificar no PostgreSQL
docker-compose exec postgres psql -U postgres -d products_db \
  -c "SELECT * FROM products ORDER BY created_at DESC LIMIT 5;"

# Testar mensagem inválida
docker-compose exec kafka kafka-console-producer \
  --bootstrap-server localhost:9092 \
  --topic product-registrations
# Enviar JSON inválido:
{"invalid": "data"}
# Deve ir para DLQ

# Verificar DLQ
docker-compose exec kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic product-registrations-dlq \
  --from-beginning

# Testar graceful shutdown
docker-compose stop product-consumer
# Verificar logs - deve processar mensagem atual antes de parar

# Testar idempotência (enviar mesmo produto 2x)
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Produto Duplicado",
    "price": 100.00,
    "description": "Teste de idempotência"
  }'
# Enviar novamente
# Verificar que só existe 1 registro no DB

# Stress test (enviar múltiplos produtos)
for i in {1..100}; do
  curl -X POST http://localhost:8000/products \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"Produto $i\",\"price\":$i.00,\"description\":\"Teste $i\"}"
  sleep 0.1
done

# Verificar contagem
docker-compose exec postgres psql -U postgres -d products_db \
  -c "SELECT COUNT(*) FROM products;"
```

### Entregáveis
- `services/product-consumer/internal/consumer/message_handler.go`
- `services/product-consumer/internal/consumer/kafka_consumer.go`
- `services/product-consumer/cmd/main.go` (atualizado)
- `docker-compose.yml` (atualizado com product-consumer)

---

## 💥 Fase 7: Injeção de Erros e Latência (Chaos Engineering)

### Objetivo
Adicionar capacidade de simular erros e latência para testar observabilidade e resiliência do sistema.

### Descrição
Implementar middlewares e configurações que permitam injetar falhas controladas para validar métricas, logs e comportamento sob stress.

### Prompt para Implementação

```
Implementar chaos engineering com injeção controlada de erros e latência.

Implementar:

1. **Product API (Python) - Middleware de Chaos**:

   **app/middleware/error_injection.py**:
   - Middleware FastAPI que:
     * Lê configuração ERROR_RATE (0.0 a 1.0)
     * Aleatoriamente injeta erros HTTP:
       - 500 Internal Server Error (50%)
       - 503 Service Unavailable (30%)
       - 429 Too Many Requests (20%)
     * Aplica apenas em endpoints específicos (não em /health)
     * Log de erros injetados
   
   **app/middleware/latency_injection.py**:
   - Middleware que:
     * Lê configuração LATENCY_MS_MIN e LATENCY_MS_MAX
     * Adiciona delay aleatório antes de processar request
     * Log de latência injetada
     * Aplica em todos os endpoints exceto /health
   
   **app/main.py** (atualizar):
   - Adicionar middlewares condicionalmente:
     * Apenas se ENABLE_CHAOS=true
     * Ordem: latency → error → application
   
   **app/config.py** (atualizar):
   - Adicionar:
     * ENABLE_CHAOS: bool = False
     * ERROR_RATE: float = 0.1
     * LATENCY_MS_MIN: int = 100
     * LATENCY_MS_MAX: int = 2000

2. **Product Consumer (Golang) - Chaos Injection**:

   **internal/consumer/chaos.go**:
   - Funções de chaos:
     * InjectRandomError(rate float64) error
     * InjectLatency(minMs, maxMs int)
     * ShouldFail(rate float64) bool
   
   **internal/consumer/message_handler.go** (atualizar):
   - No ProcessMessage, adicionar:
     * Injeção de latência antes de processar
     * Injeção de erro aleatório
     * Controlado por env vars:
       - ENABLE_CHAOS
       - CHAOS_ERROR_RATE
       - CHAOS_LATENCY_MIN_MS
       - CHAOS_LATENCY_MAX_MS
   
   **internal/config/config.go** (atualizar):
   - Adicionar configurações de chaos

3. **Health Checks**:

   **Product API**:
   - GET /health: sempre retorna 200
   - GET /ready: verifica conectividade Kafka
     * Retorna 200 se OK
     * Retorna 503 se Kafka offline
   
   **Product Consumer**:
   - HTTP server em porta 8081:
     * GET /health: sempre retorna 200
     * GET /ready: verifica Kafka e PostgreSQL
   
   **docker-compose.yml** (atualizar):
   - Adicionar health checks em todos os serviços
   - Usar endpoints /health

4. **Documentação**:
   - README.md principal atualizar com:
     * Como habilitar chaos engineering
     * Variáveis de ambiente
     * Exemplos de uso
     * Casos de teste

Requisitos:
- Chaos desabilitado por padrão (produção-safe)
- Configuração via variáveis de ambiente
- Logs detalhados de chaos events
- Health checks não afetados por chaos
- Fácil ligar/desligar para testes
```

### Critérios de Aceite

- [ ] Chaos pode ser habilitado via env var
- [ ] Erros são injetados na taxa configurada
- [ ] Latência é adicionada nos ranges configurados
- [ ] Health checks não são afetados
- [ ] Logs identificam eventos de chaos
- [ ] Sistema se recupera após erros injetados
- [ ] Ready checks validam dependências

### Comandos de Teste

```bash
# Testar sem chaos (padrão)
docker-compose up -d
curl http://localhost:8000/products
# Deve funcionar normalmente

# Habilitar chaos na API
docker-compose stop product-api
# Editar .env ou docker-compose.yml:
# ENABLE_CHAOS=true
# ERROR_RATE=0.3
# LATENCY_MS_MIN=500
# LATENCY_MS_MAX=3000
docker-compose up -d product-api

# Enviar múltiplas requests e observar comportamento
for i in {1..20}; do
  echo "Request $i:"
  time curl -X POST http://localhost:8000/products \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"Produto $i\",\"price\":100.00}"
  echo ""
done
# Observar: algumas falham, outras têm latência alta

# Verificar logs
docker-compose logs product-api | grep -i "chaos\|error\|latency"

# Testar health checks
curl http://localhost:8000/health
# Sempre deve retornar 200

curl http://localhost:8000/ready
# Retorna status das dependências

# Habilitar chaos no consumer
# Editar docker-compose.yml:
# ENABLE_CHAOS=true
# CHAOS_ERROR_RATE=0.2
# CHAOS_LATENCY_MIN_MS=1000
# CHAOS_LATENCY_MAX_MS=5000
docker-compose restart product-consumer

# Enviar produtos e observar processamento lento/falhas
docker-compose logs -f product-consumer

# Verificar que mensagens são reprocessadas após falha
# (não commitadas)

# Testar ready check do consumer
docker exec product-consumer curl localhost:8081/ready

# Simular Kafka offline
docker-compose stop kafka
curl http://localhost:8000/ready
# Deve retornar 503

curl http://localhost:8081/ready
# Consumer deve retornar 503

# Religar Kafka
docker-compose start kafka

# Aguardar e verificar ready novamente
sleep 10
curl http://localhost:8000/ready
# Deve voltar a 200
```

### Entregáveis
- `services/product-api/app/middleware/error_injection.py`
- `services/product-api/app/middleware/latency_injection.py`
- `services/product-consumer/internal/consumer/chaos.go`
- Arquivos config atualizados
- Health check endpoints
- Documentação atualizada

---

## ✅ Fase 8: Testes, Validação e Documentação Final

### Objetivo
Adicionar testes automatizados, scripts de validação e documentação completa do projeto.

### Descrição
Garantir qualidade através de testes unitários e integração, além de documentar todo o projeto para facilitar manutenção e aprendizado.

### Prompt para Implementação

```
Adicionar testes, scripts de automação e documentação final do projeto.

Implementar:

1. **Testes Python (Product API)**:

   **tests/test_routes.py**:
   - Testes do endpoint POST /products:
     * Produto válido retorna 202
     * Validação de campos obrigatórios
     * Validação de tipos (price decimal)
     * Validação de ranges (price > 0)
   - Testes do /health e /ready
   - Mock do KafkaProducer
   - Usar pytest fixtures
   
   **tests/test_kafka_producer.py**:
   - Testar envio de mensagens
   - Testar retry em falhas
   - Mock do KafkaProducer real
   
   **tests/conftest.py**:
   - Fixtures compartilhadas
   - TestClient do FastAPI
   
   **requirements-dev.txt**:
   - pytest
   - pytest-cov
   - pytest-asyncio
   - httpx
   - faker

2. **Testes Go (Product Consumer)**:

   **tests/message_handler_test.go**:
   - Testes de ProcessMessage:
     * JSON válido
     * JSON inválido
     * Erro de repository (mock)
     * Retry behavior
   
   **tests/repository_test.go** (opcional):
   - Integration tests com testcontainers
   - Testar CreateProduct
   - Testar idempotência
   
   - Usar testify/assert e testify/mock

3. **Scripts de Automação**:

   **scripts/setup.sh**:
   ```bash
   #!/bin/bash
   # Setup completo do projeto:
   # - Copiar .env.example para .env
   # - Build de todas as imagens
   # - Criar redes Docker
   # - Inicializar banco de dados
   # - Criar tópicos Kafka
   # - Aguardar health checks
   ```
   
   **scripts/test-e2e.sh**:
   ```bash
   #!/bin/bash
   # Teste end-to-end:
   # - Subir stack
   # - Aguardar ready
   # - Enviar produtos via API
   # - Validar no PostgreSQL
   # - Cleanup
   # - Reportar resultados
   ```
   
   **scripts/run-tests.sh**:
   ```bash
   #!/bin/bash
   # Rodar todos os testes:
   # - Testes unitários Python
   # - Testes unitários Go
   # - Coverage report
   ```

4. **Documentação**:

   **README.md** (raiz do projeto):
   - Descrição e objetivos
   - Arquitetura (diagrama)
   - Tecnologias utilizadas
   - Pré-requisitos
   - Quick start
   - Estrutura de diretórios
   - Variáveis de ambiente
   - Como rodar testes
   - Chaos engineering
   - Troubleshooting
   - Próximos passos (observabilidade)
   
   **services/product-api/README.md**:
   - Detalhes do serviço Python
   - API endpoints (OpenAPI)
   - Como desenvolver localmente
   - Como rodar testes
   - Variáveis de ambiente
   
   **services/product-consumer/README.md**:
   - Detalhes do serviço Go
   - Arquitetura interna
   - Como desenvolver localmente
   - Como rodar testes
   - Variáveis de ambiente
   
   **CONTRIBUTING.md**:
   - Guidelines de contribuição
   - Code style
   - Commit conventions
   - Pull request process

5. **Collections e Exemplos**:

   **postman/microservices-collection.json**:
   - Collection Postman com:
     * Criar produto
     * Health checks
     * Ready checks
     * Testes com dados inválidos
     * Variáveis de ambiente
   
   **examples/curl-examples.sh**:
   - Scripts curl para todas as operações
   - Exemplos de sucesso e erro
   - Comentários explicativos

6. **CI/CD Preparação** (opcional):

   **.github/workflows/tests.yml**:
   - GitHub Actions para rodar testes
   - Lint e format check
   - Build de imagens Docker
   - Integration tests

Requisitos:
- Coverage mínimo de 70%
- Todos os testes devem passar
- Documentação clara e objetiva
- Scripts devem ser idempotentes
- Exemplos funcionais
```

### Critérios de Aceite

- [ ] Testes unitários Python executam com sucesso
- [ ] Testes unitários Go executam com sucesso
- [ ] Coverage >= 70% em ambos os serviços
- [ ] Scripts de setup e teste funcionam
- [ ] README está completo e atualizado
- [ ] Postman collection funciona
- [ ] Documentação de cada serviço está clara
- [ ] Projeto pode ser clonado e executado facilmente

### Comandos de Teste

```bash
# Testes Python
cd services/product-api
pip install -r requirements-dev.txt
pytest tests/ -v --cov=app --cov-report=html
# Abrir htmlcov/index.html

# Testes Go
cd services/product-consumer
go test ./... -v -cover
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Rodar setup completo
chmod +x scripts/*.sh
./scripts/setup.sh

# Rodar testes E2E
./scripts/test-e2e.sh

# Validar com Postman
newman run postman/microservices-collection.json \
  --environment postman/environment.json

# Validar curl examples
chmod +x examples/curl-examples.sh
./examples/curl-examples.sh

# Lint Python
cd services/product-api
pip install black flake8 mypy
black app/ --check
flake8 app/
mypy app/

# Lint Go
cd services/product-consumer
go fmt ./...
go vet ./...
golangci-lint run

# Build todas as imagens
docker-compose build

# Subir stack completa
docker-compose up -d

# Validar saúde
curl http://localhost:8000/health
curl http://localhost:8000/ready
curl http://localhost:8081/health
curl http://localhost:8081/ready

# Enviar produto teste
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Produto Final",
    "price": 999.99,
    "description": "Teste final do sistema"
  }'

# Validar no DB
docker-compose exec postgres psql -U postgres -d products_db \
  -c "SELECT * FROM products WHERE name = 'Produto Final';"

# Verificar logs
docker-compose logs --tail=50 product-api
docker-compose logs --tail=50 product-consumer

# Cleanup
docker-compose down -v
```

### Entregáveis
- Testes unitários completos (Python e Go)
- Scripts de automação (setup, test-e2e, run-tests)
- README.md completo e detalhado
- Documentação de cada serviço
- Postman collection
- Curl examples
- CONTRIBUTING.md (opcional)
- CI/CD configs (opcional)

---

## � Fase 9: OpenTelemetry Collector + Instrumentação Completa

### Objetivo
Instrumentar ambos microserviços com OpenTelemetry SDK e configurar OTLP Collector para receber e processar telemetria.

### Descrição
Implementar observabilidade nativa com OpenTelemetry, incluindo distributed tracing, context propagation e exportação de telemetria para múltiplos backends. Esta é uma das **competências principais da vaga Wolt**.

### Prompt para Implementação

```
Implementar instrumentação OpenTelemetry completa nos microserviços.

Implementar:

1. **Product API (Python) - OpenTelemetry SDK**:

   **app/telemetry/tracing.py**:
   - Configurar OpenTelemetry SDK para Python
   - Auto-instrumentation do FastAPI
   - Manual instrumentation para Kafka producer
   - Resource attributes (service.name, service.version, deployment.environment)
   - OTLP exporter para Collector
   - Sampling strategy (AlwaysOn para dev)
   
   **app/telemetry/metrics.py**:
   - Configurar métricas OpenTelemetry:
     * Counter: products_created_total
     * Histogram: product_creation_duration_seconds
     * Gauge: active_requests
   - OTLP metrics exporter
   
   **app/main.py** (atualizar):
   - Inicializar telemetry no startup
   - Adicionar middleware de tracing
   - Context propagation nos headers HTTP
   - Trace ID logging
   
   **requirements.txt** (adicionar):
   - opentelemetry-api
   - opentelemetry-sdk
   - opentelemetry-instrumentation-fastapi
   - opentelemetry-instrumentation-kafka-python
   - opentelemetry-exporter-otlp-proto-grpc

2. **Product Consumer (Golang) - OpenTelemetry SDK**:

   **pkg/telemetry/tracing.go**:
   - Configurar OpenTelemetry SDK para Go
   - TracerProvider com OTLP exporter
   - Resource attributes
   - Context propagation
   - Span attributes customizados
   
   **pkg/telemetry/metrics.go**:
   - MeterProvider
   - Métricas:
     * Int64Counter: messages_consumed_total
     * Float64Histogram: message_processing_duration
     * Int64UpDownCounter: messages_in_flight
   
   **internal/consumer/message_handler.go** (atualizar):
   - Extrair trace context da mensagem Kafka
   - Criar span para processamento
   - Adicionar span attributes (product_id, etc)
   - Registrar events em caso de erro
   - Propagate context para repository
   
   **internal/repository/product_repository.go** (atualizar):
   - Aceitar context com trace
   - Criar span para operações DB
   - Span attributes (query, table)
   
   **go.mod** (adicionar):
   - go.opentelemetry.io/otel
   - go.opentelemetry.io/otel/sdk
   - go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
   - go.opentelemetry.io/contrib/instrumentation/github.com/segmentio/kafka-go/otelsarama

3. **OpenTelemetry Collector**:

   **infrastructure/observability/otel-collector/otel-config.yaml**:
   ```yaml
   receivers:
     otlp:
       protocols:
         grpc:
           endpoint: 0.0.0.0:4317
         http:
           endpoint: 0.0.0.0:4318
   
   processors:
     batch:
       timeout: 10s
       send_batch_size: 1024
     
   exporters:
     otlp/clickhouse:
       endpoint: clickhouse:9000
       tls:
         insecure: true
     debug:
        verbosity: detailed
     
     prometheus:
       endpoint: "0.0.0.0:8889"

   
   service:
     pipelines:
       traces:
         receivers: [otlp]
         processors: [resource, tail_sampling, batch]
         exporters: [debug,otlp/clickhouse]
       
       metrics:
         receivers: [otlp]
         processors: [resource, batch]
         exporters: [prometheus]
   ```

4. **docker-compose.observability.yml**:
   - Serviço otel-collector (otel/opentelemetry-collector-contrib:latest)
   - Port mappings: 4317 (gRPC), 4318 (HTTP), 8889 (Prometheus)
   - Volume mount para config
   - Serviço Jaeger UI (opcional)
   - Depends_on: clickhouse

5. **Kafka Message Context Propagation**:
   - No producer Python: injetar trace context nos headers Kafka
   - No consumer Go: extrair trace context dos headers
   - Manter trace_id consistente entre serviços

Requisitos:
- Trace context propagation funcionando end-to-end
- Spans detalhados em cada operação
- Resource attributes identificando serviços
- OTLP exportando para Collector
- Métricas sendo coletadas
- Sampling configurável
- Logs correlacionados com trace_id
```

### Critérios de Aceite

- [ ] Traces são gerados em Python e Go
- [ ] Context propagation funciona via Kafka
- [ ] Collector recebe telemetria OTLP
- [ ] Spans têm attributes relevantes
- [ ] Métricas são exportadas
- [ ] Sampling está funcionando
- [ ] Jaeger UI mostra traces (opcional)
- [ ] Logs incluem trace_id

### Comandos de Teste

```bash
# Subir stack com observabilidade
docker-compose -f docker-compose.yml -f docker-compose.observability.yml up -d

# Verificar Collector
docker-compose logs otel-collector | grep "Everything is ready"

# Enviar produto e gerar trace
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Produto Traced",
    "price": 299.99,
    "description": "Com OpenTelemetry"
  }'

# Ver logs do Collector (traces sendo recebidos)
docker-compose logs -f otel-collector

# Acessar Jaeger UI (se habilitado)
open http://localhost:16686
# Procurar por traces do serviço "product-api"

# Verificar métricas no endpoint Prometheus
curl http://localhost:8889/metrics | grep product

# Verificar que trace_id está nos logs
docker-compose logs product-api | grep trace_id
docker-compose logs product-consumer | grep trace_id

# Validar context propagation
# Os logs devem mostrar o mesmo trace_id em API e Consumer
```

### Entregáveis
- OpenTelemetry SDK configurado (Python e Go)
- OTLP Collector funcionando
- Context propagation via Kafka
- Métricas instrumentadas
- docker-compose.observability.yml

---

## 📊 Fase 10: ClickHouse para Traces e Logs em Escala

### Objetivo
Configurar ClickHouse como backend de observabilidade para armazenar e consultar bilhões de traces e logs, demonstrando capacidade de trabalhar com **large-scale distributed databases**.

### Descrição
Implementar schema otimizado no ClickHouse para telemetria, criar queries de análise e demonstrar capacidade de processar dados em escala, alinhado com a experiência que a Wolt valoriza.

### Prompt para Implementação

```
Configurar ClickHouse como backend de observabilidade escalável.

Implementar:

1. **infrastructure/observability/clickhouse/init-db.sql**:
   
   ```sql
   -- Database para observabilidade
   CREATE DATABASE IF NOT EXISTS observability;
   
   -- Tabela de traces otimizada para escala
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
   TTL timestamp + INTERVAL 30 DAY
   SETTINGS index_granularity = 8192;
   
   -- Tabela de logs estruturados
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
   TTL timestamp + INTERVAL 30 DAY
   SETTINGS index_granularity = 8192;
   
   -- Materialized View para métricas agregadas (RED metrics)
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
   
   -- View para análise de erros
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
   
   -- View para análise de latência
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
   ```

2. **Queries de Análise (queries/clickhouse-analysis.sql)**:
   
   ```sql
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
   ```

3. **OpenTelemetry Collector ClickHouse Exporter** (atualizar config):
   
   Usar `otel/opentelemetry-collector-contrib` com ClickHouse exporter:
   ```yaml
   exporters:
     clickhouse:
       endpoint: tcp://clickhouse:9000?database=observability
       traces_table_name: traces
       logs_table_name: logs
       ttl_days: 30
       compression: lz4
   ```

4. **docker-compose.observability.yml** (atualizar):
   - ClickHouse Server (clickhouse/clickhouse-server:latest)
   - Persistent volumes
   - Port 8123 (HTTP), 9000 (Native)
   - Health check
   - Init script montado

5. **Scripts de Análise (scripts/analyze-traces.sh)**:
   ```bash
   #!/bin/bash
   # Script para análises comuns no ClickHouse
   
   CLICKHOUSE_CLIENT="docker-compose exec clickhouse clickhouse-client"
   
   echo "=== Top 10 Slowest Operations ==="
   $CLICKHOUSE_CLIENT --query="
   SELECT service_name, operation_name, 
          avg(duration_ns)/1000000 as avg_ms
   FROM observability.traces
   WHERE timestamp >= now() - INTERVAL 1 HOUR
   GROUP BY service_name, operation_name
   ORDER BY avg_ms DESC LIMIT 10
   FORMAT PrettyCompact
   "
   
   echo "=== Error Rate by Service ==="
   $CLICKHOUSE_CLIENT --query="
   SELECT service_name,
          countIf(status_code='ERROR') as errors,
          count() as total,
          (errors/total)*100 as error_rate
   FROM observability.traces
   WHERE timestamp >= now() - INTERVAL 1 HOUR
   GROUP BY service_name
   FORMAT PrettyCompact
   "
   ```

Requisitos:
- Schema otimizado para bilhões de registros
- Particionamento por data
- Compression codecs apropriados
- TTL para retenção de dados
- Índices para queries rápidas
- Materialized views para agregações
- Queries de análise prontas para uso
```

### Critérios de Aceite

- [ ] ClickHouse está rodando e acessível
- [ ] Tabelas de traces e logs criadas
- [ ] Collector exporta para ClickHouse
- [ ] Queries de análise executam rapidamente
- [ ] Materialized views agregam dados
- [ ] TTL está configurado
- [ ] Compression está funcionando
- [ ] Scripts de análise funcionam

### Comandos de Teste

```bash
# Verificar ClickHouse
docker-compose exec clickhouse clickhouse-client --query "SELECT version()"

# Verificar databases
docker-compose exec clickhouse clickhouse-client --query "SHOW DATABASES"

# Verificar tabelas
docker-compose exec clickhouse clickhouse-client --query "SHOW TABLES FROM observability"

# Gerar traces
for i in {1..100}; do
  curl -X POST http://localhost:8000/products \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"Product $i\",\"price\":$i.00}"
  sleep 0.1
done

# Aguardar alguns segundos para exportação

# Verificar dados no ClickHouse
docker-compose exec clickhouse clickhouse-client --query="
SELECT count() FROM observability.traces
"

# Análise de latência
docker-compose exec clickhouse clickhouse-client --query="
SELECT
    service_name,
    operation_name,
    count() as count,
    avg(duration_ns)/1000000 as avg_ms,
    quantile(0.99)(duration_ns)/1000000 as p99_ms
FROM observability.traces
GROUP BY service_name, operation_name
FORMAT PrettyCompact
"

# Ver traces recentes
docker-compose exec clickhouse clickhouse-client --query="
SELECT
    timestamp,
    service_name,
    operation_name,
    duration_ns/1000000 as duration_ms,
    status_code
FROM observability.traces
ORDER BY timestamp DESC
LIMIT 10
FORMAT PrettyCompact
"

# Executar script de análise
chmod +x scripts/analyze-traces.sh
./scripts/analyze-traces.sh
```

### Entregáveis
- ClickHouse configurado e otimizado
- Schema de traces e logs
- Materialized views para agregações
- Queries de análise prontas
- Scripts de automação
- Documentação de queries

---

## 📈 Fase 11: Prometheus + Grafana Dashboards

### Objetivo
Implementar coleta de métricas com Prometheus e criar dashboards profissionais no Grafana, demonstrando expertise em **Prometheus, Grafana** conforme requisitos da vaga.

### Descrição
Configurar stack de métricas completa com RED metrics (Rate, Errors, Duration), criar dashboards customizados e integrar com ClickHouse para visualização unificada.

### Prompt para Implementação

```
Implementar stack completa de métricas com Prometheus e Grafana.

Implementar:

1. **Prometheus Configuration**:

   **infrastructure/observability/prometheus/prometheus.yml**:
   ```yaml
   global:
     scrape_interval: 15s
     evaluation_interval: 15s
     external_labels:
       cluster: 'dev'
       environment: 'development'
   
   scrape_configs:
     # OpenTelemetry Collector metrics
     - job_name: 'otel-collector'
       static_configs:
         - targets: ['otel-collector:8889']
       metric_relabel_configs:
         - source_labels: [__name__]
           regex: 'otelcol_.*'
           action: keep
     
     # Product API metrics
     - job_name: 'product-api'
       static_configs:
         - targets: ['product-api:8000']
       metrics_path: '/metrics'
       scrape_interval: 10s
     
     # Product Consumer metrics
     - job_name: 'product-consumer'
       static_configs:
         - targets: ['product-consumer:8081']
       metrics_path: '/metrics'
       scrape_interval: 10s
     
     # Infrastructure metrics
     - job_name: 'kafka'
       static_configs:
         - targets: ['kafka:9999']  # JMX exporter
     
     - job_name: 'postgres'
       static_configs:
         - targets: ['postgres-exporter:9187']
     
     # Prometheus self-monitoring
     - job_name: 'prometheus'
       static_configs:
         - targets: ['localhost:9090']
   
   # Recording rules para otimização
   rule_files:
     - '/etc/prometheus/rules/*.yml'
   ```

   **infrastructure/observability/prometheus/rules/app-rules.yml**:
   ```yaml
   groups:
     - name: product_api_rules
       interval: 30s
       rules:
         # Request rate
         - record: job:product_api:requests_per_second
           expr: rate(products_created_total[1m])
         
         # Error rate
         - record: job:product_api:error_rate
           expr: |
             rate(products_created_total{status="error"}[5m]) 
             / 
             rate(products_created_total[5m])
         
         # Average latency
         - record: job:product_api:latency_seconds:avg
           expr: |
             rate(product_creation_duration_seconds_sum[5m])
             /
             rate(product_creation_duration_seconds_count[5m])
   ```

2. **Expor Métricas nos Microserviços**:

   **Product API (Python) - app/api/metrics_endpoint.py**:
   ```python
   from prometheus_client import make_asgi_app
   
   # Endpoint /metrics para Prometheus
   metrics_app = make_asgi_app()
   ```
   
   Integrar no FastAPI main.py:
   ```python
   from prometheus_client import Counter, Histogram, Gauge
   
   # RED Metrics
   REQUEST_COUNT = Counter(
       'products_created_total',
       'Total products created',
       ['status', 'endpoint']
   )
   
   REQUEST_LATENCY = Histogram(
       'product_creation_duration_seconds',
       'Product creation latency',
       ['endpoint']
   )
   
   ACTIVE_REQUESTS = Gauge(
       'active_requests',
       'Active requests'
   )
   
   # Mount metrics endpoint
   app.mount("/metrics", metrics_app)
   ```

   **Product Consumer (Go) - pkg/metrics/prometheus.go**:
   ```go
   package metrics
   
   import (
       "github.com/prometheus/client_golang/prometheus"
       "github.com/prometheus/client_golang/prometheus/promhttp"
       "net/http"
   )
   
   var (
       MessagesConsumed = prometheus.NewCounterVec(
           prometheus.CounterOpts{
               Name: "messages_consumed_total",
               Help: "Total messages consumed",
           },
           []string{"topic", "status"},
       )
       
       ProcessingDuration = prometheus.NewHistogramVec(
           prometheus.HistogramOpts{
               Name:    "message_processing_duration_seconds",
               Help:    "Message processing duration",
               Buckets: prometheus.DefBuckets,
           },
           []string{"topic"},
       )
       
       InFlightMessages = prometheus.NewGauge(
           prometheus.GaugeOpts{
               Name: "messages_in_flight",
               Help: "Messages currently being processed",
           },
       )
   )
   
   func InitMetrics() {
       prometheus.MustRegister(MessagesConsumed)
       prometheus.MustRegister(ProcessingDuration)
       prometheus.MustRegister(InFlightMessages)
   }
   
   func StartMetricsServer(port string) {
       http.Handle("/metrics", promhttp.Handler())
       go http.ListenAndServe(":"+port, nil)
   }
   ```

3. **Grafana Dashboards**:

   **infrastructure/observability/grafana/provisioning/datasources.yml**:
   ```yaml
   apiVersion: 1
   
   datasources:
     - name: Prometheus
       type: prometheus
       access: proxy
       url: http://prometheus:9090
       isDefault: true
       editable: true
     
     - name: ClickHouse
       type: vertamedia-clickhouse-datasource
       access: proxy
       url: http://clickhouse:8123
       jsonData:
         defaultDatabase: observability
       editable: true
   ```

   **infrastructure/observability/grafana/dashboards/microservices-overview.json**:
   ```json
   {
     "dashboard": {
       "title": "Microservices Overview",
       "panels": [
         {
           "title": "Request Rate (RPS)",
           "targets": [{
             "expr": "sum(rate(products_created_total[5m])) by (service_name)"
           }],
           "type": "graph"
         },
         {
           "title": "Error Rate (%)",
           "targets": [{
             "expr": "100 * sum(rate(products_created_total{status=\"error\"}[5m])) / sum(rate(products_created_total[5m]))"
           }],
           "type": "graph"
         },
         {
           "title": "P95 Latency",
           "targets": [{
             "expr": "histogram_quantile(0.95, rate(product_creation_duration_seconds_bucket[5m]))"
           }],
           "type": "graph"
         },
         {
           "title": "Active Requests",
           "targets": [{
             "expr": "sum(active_requests) by (service_name)"
           }],
           "type": "stat"
         }
       ]
     }
   }
   ```

   **Dashboard para ClickHouse Traces**:
   - Panel com query ClickHouse para latência por operação
   - Heatmap de distribuição de latências
   - Table de top slow queries
   - Trace search interface

4. **docker-compose.observability.yml** (atualizar):
   ```yaml
   prometheus:
     image: prom/prometheus:latest
     volumes:
       - ./infrastructure/observability/prometheus:/etc/prometheus
       - prometheus_data:/prometheus
     command:
       - '--config.file=/etc/prometheus/prometheus.yml'
       - '--storage.tsdb.path=/prometheus'
       - '--storage.tsdb.retention.time=30d'
     ports:
       - "9090:9090"
     networks:
       - observability
   
   grafana:
     image: grafana/grafana:latest
     volumes:
       - ./infrastructure/observability/grafana/provisioning:/etc/grafana/provisioning
       - ./infrastructure/observability/grafana/dashboards:/var/lib/grafana/dashboards
       - grafana_data:/var/lib/grafana
     environment:
       - GF_SECURITY_ADMIN_PASSWORD=admin
       - GF_USERS_ALLOW_SIGN_UP=false
     ports:
       - "3000:3000"
     networks:
       - observability
     depends_on:
       - prometheus
       - clickhouse
   
   volumes:
     prometheus_data:
     grafana_data:
   ```

Requisitos:
- Prometheus scraping todos os targets
- Métricas RED implementadas
- Recording rules otimizadas
- Grafana com datasources configurados
- Dashboards provisionados automaticamente
- Queries ClickHouse integradas
- Alerting preparado (próxima fase)
```

### Critérios de Aceite

- [ ] Prometheus coleta métricas de todos os serviços
- [ ] Endpoints /metrics funcionando
- [ ] Recording rules executando
- [ ] Grafana acessível e configurado
- [ ] Dashboards carregam automaticamente
- [ ] ClickHouse datasource conectado
- [ ] Visualizações mostram dados reais
- [ ] RED metrics estão visíveis

### Comandos de Teste

```bash
# Verificar Prometheus targets
open http://localhost:9090/targets
# Todos devem estar "UP"

# Testar query PromQL
curl 'http://localhost:9090/api/v1/query?query=up'

# Ver métricas da API
curl http://localhost:8000/metrics | grep products_created_total

# Ver métricas do Consumer
curl http://localhost:8081/metrics | grep messages_consumed_total

# Gerar carga para métricas
for i in {1..1000}; do
  curl -X POST http://localhost:8000/products \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"Load test $i\",\"price\":99.99}" &
done
wait

# Acessar Grafana
open http://localhost:3000
# Login: admin/admin

# Verificar dashboards
# Deve mostrar:
# - Request rate aumentando
# - Latência P95
# - Error rate
# - Active requests

# Query manual no Prometheus
# http://localhost:9090/graph
# Testar: rate(products_created_total[5m])

# Verificar recording rules
curl http://localhost:9090/api/v1/rules | jq
```

### Entregáveis
- Prometheus configurado e scraping
- Recording rules otimizadas
- Grafana com datasources
- Dashboards profissionais
- Integração ClickHouse
- Métricas RED implementadas

---

## 🎯 Fase 12: SRE Practices - SLIs/SLOs + Alerting

### Objetivo
Implementar práticas SRE com definição de SLIs/SLOs, error budgets, alerting inteligente e runbooks de incident response, demonstrando **SRE principles** conforme requisitos da vaga.

### Descrição
Aplicar metodologia SRE para garantir confiabilidade, definir Service Level Objectives, configurar alerting e documentar processos de resposta a incidentes.

### Prompt para Implementação

```
Implementar SRE practices com SLIs/SLOs, alerting e incident response.

Implementar:

1. **Definição de SLIs/SLOs**:

   **sre/slos/product-service-slo.yaml**:
   ```yaml
   apiVersion: v1
   kind: SLO
   metadata:
     name: product-service-availability
     description: "Product service availability and performance SLOs"
   
   spec:
     service: product-api
     
     slis:
       # Availability SLI
       - name: availability
         description: "Percentage of successful requests"
         query: |
           sum(rate(products_created_total{status!~"5.."}[5m]))
           /
           sum(rate(products_created_total[5m]))
         unit: percent
       
       # Latency SLI
       - name: latency_p95
         description: "95th percentile latency"
         query: |
           histogram_quantile(0.95, 
             rate(product_creation_duration_seconds_bucket[5m])
           )
         unit: seconds
       
       # Throughput SLI
       - name: throughput
         description: "Requests per second"
         query: |
           sum(rate(products_created_total[1m]))
         unit: rps
     
     slos:
       # 99.9% availability target
       - name: availability_target
         sli: availability
         target: 0.999
         window: 30d
         
       # P95 latency < 500ms
       - name: latency_target
         sli: latency_p95
         target: 0.5
         window: 30d
       
       # Minimum throughput
       - name: throughput_target
         sli: throughput
         target: 10
         window: 1h
     
     error_budget:
       window: 30d
       calculation: |
         # 99.9% uptime = 0.1% error budget
         # In 30 days = 43200 minutes
         # Error budget = 43.2 minutes of downtime
         (1 - slo_target) * window_minutes
   ```

2. **Alerting Rules**:

   **infrastructure/observability/prometheus/rules/alerts.yml**:
   ```yaml
   groups:
     - name: slo_alerts
       interval: 30s
       rules:
         # Error budget burn rate alerts (multi-window)
         - alert: ErrorBudgetBurnRateCritical
           expr: |
             (
               sum(rate(products_created_total{status=~"5.."}[1h]))
               /
               sum(rate(products_created_total[1h]))
             ) > (14.4 * (1 - 0.999))
           for: 2m
           labels:
             severity: critical
             slo: availability
           annotations:
             summary: "Critical error budget burn rate"
             description: "Burning through error budget 14.4x faster than target (1h window)"
             runbook_url: "https://runbooks.example.com/error-budget-burn"
         
         - alert: ErrorBudgetBurnRateWarning
           expr: |
             (
               sum(rate(products_created_total{status=~"5.."}[6h]))
               /
               sum(rate(products_created_total[6h]))
             ) > (6 * (1 - 0.999))
           for: 15m
           labels:
             severity: warning
             slo: availability
           annotations:
             summary: "High error budget burn rate"
             description: "Burning through error budget 6x faster than target (6h window)"
         
         # Latency SLO violation
         - alert: LatencySLOViolation
           expr: |
             histogram_quantile(0.95,
               rate(product_creation_duration_seconds_bucket[5m])
             ) > 0.5
           for: 10m
           labels:
             severity: warning
             slo: latency
           annotations:
             summary: "P95 latency exceeds SLO"
             description: "P95 latency is {{ $value }}s, target is 0.5s"
             runbook_url: "https://runbooks.example.com/high-latency"
         
         # Service down
         - alert: ServiceDown
           expr: up{job="product-api"} == 0
           for: 1m
           labels:
             severity: critical
           annotations:
             summary: "Service {{ $labels.job }} is down"
             description: "{{ $labels.instance }} has been down for more than 1 minute"
             runbook_url: "https://runbooks.example.com/service-down"
         
         # High error rate
         - alert: HighErrorRate
           expr: |
             (
               sum(rate(products_created_total{status=~"5.."}[5m]))
               /
               sum(rate(products_created_total[5m]))
             ) > 0.05
           for: 5m
           labels:
             severity: warning
           annotations:
             summary: "High error rate detected"
             description: "Error rate is {{ $value | humanizePercentage }} (threshold: 5%)"
         
         # Consumer lag
         - alert: KafkaConsumerLag
           expr: |
             sum(kafka_consumer_lag) by (topic, consumer_group) > 1000
           for: 10m
           labels:
             severity: warning
           annotations:
             summary: "High Kafka consumer lag"
             description: "Consumer lag is {{ $value }} messages"
   
     - name: infrastructure_alerts
       interval: 30s
       rules:
         - alert: HighMemoryUsage
           expr: |
             (node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes)
             / node_memory_MemTotal_bytes > 0.9
           for: 5m
           labels:
             severity: warning
           annotations:
             summary: "High memory usage"
         
         - alert: DiskSpaceRunningOut
           expr: |
             (node_filesystem_avail_bytes / node_filesystem_size_bytes) < 0.1
           for: 10m
           labels:
             severity: critical
           annotations:
             summary: "Disk space running out"
   ```

3. **Alertmanager Configuration**:

   **infrastructure/observability/alertmanager/alertmanager.yml**:
   ```yaml
   global:
     resolve_timeout: 5m
   
   route:
     group_by: ['alertname', 'severity']
     group_wait: 10s
     group_interval: 10s
     repeat_interval: 12h
     receiver: 'default'
     
     routes:
       # Critical alerts -> PagerDuty (simulado)
       - match:
           severity: critical
         receiver: 'pagerduty'
         continue: true
       
       # Warnings -> Slack
       - match:
           severity: warning
         receiver: 'slack'
       
       # SLO violations -> Dedicated channel
       - match_re:
           alertname: '.*SLO.*|.*ErrorBudget.*'
         receiver: 'slo-alerts'
   
   receivers:
     - name: 'default'
       webhook_configs:
         - url: 'http://webhook-logger:8080/webhook'
     
     - name: 'pagerduty'
       # pagerduty_configs:
       #   - service_key: '<KEY>'
       webhook_configs:
         - url: 'http://webhook-logger:8080/critical'
     
     - name: 'slack'
       # slack_configs:
       #   - api_url: '<WEBHOOK_URL>'
       #     channel: '#alerts'
       webhook_configs:
         - url: 'http://webhook-logger:8080/slack'
     
     - name: 'slo-alerts'
       webhook_configs:
         - url: 'http://webhook-logger:8080/slo'
   
   inhibit_rules:
     # Silence warnings if critical alert is firing
     - source_match:
         severity: 'critical'
       target_match:
         severity: 'warning'
       equal: ['alertname', 'service']
   ```

4. **Runbooks**:

   **sre/runbooks/high-latency.md**:
   ```markdown
   # Runbook: High Latency
   
   ## Alert
   - **Name**: LatencySLOViolation
   - **Severity**: Warning
   - **SLO**: P95 latency < 500ms
   
   ## Impact
   - User experience degraded
   - Potential SLO violation
   - Error budget consumption
   
   ## Diagnosis
   
   1. **Check Grafana Dashboard**: 
      - URL: http://localhost:3000/d/microservices-overview
      - Look for latency spikes
   
   2. **Query ClickHouse for slow traces**:
      ```sql
      SELECT trace_id, service_name, operation_name, duration_ns/1000000 as duration_ms
      FROM observability.traces
      WHERE timestamp >= now() - INTERVAL 15 MINUTE
        AND duration_ns > 500000000
      ORDER BY duration_ns DESC
      LIMIT 20;
      ```
   
   3. **Check database performance**:
      ```bash
      docker-compose exec postgres pg_stat_statements
      ```
   
   4. **Review recent deployments**: Check if latency correlates with deployment
   
   ## Mitigation
   
   1. **Immediate**:
      - Scale horizontally if CPU/memory constrained
      - Enable caching if applicable
      - Reduce traffic with rate limiting
   
   2. **Short-term**:
      - Optimize slow database queries
      - Review and optimize code hot paths
      - Check external dependencies
   
   3. **Long-term**:
      - Implement query optimization
      - Add database indexes
      - Consider async processing
   
   ## Resolution
   
   - Monitor until P95 latency < 500ms for 15 minutes
   - Document root cause in postmortem
   - Update error budget tracking
   ```

   **sre/postmortems/template.md**:
   ```markdown
   # Postmortem Template
   
   **Incident Date**: YYYY-MM-DD  
   **Duration**: X hours  
   **Severity**: Critical/High/Medium/Low  
   **Services Affected**: [list]  
   **Error Budget Impact**: X.XX%
   
   ## Summary
   Brief description of the incident.
   
   ## Timeline (UTC)
   - **HH:MM** - Incident detected
   - **HH:MM** - Alert fired
   - **HH:MM** - Investigation started
   - **HH:MM** - Root cause identified
   - **HH:MM** - Mitigation applied
   - **HH:MM** - Service restored
   
   ## Root Cause
   Detailed explanation of what caused the incident.
   
   ## Impact
   - Requests affected: X
   - Users impacted: X
   - Revenue impact: $X
   - SLO violation: Yes/No
   - Error budget consumed: X.XX%
   
   ## Detection
   How was the incident detected?
   
   ## Response
   What was done to resolve the incident?
   
   ## Lessons Learned
   
   ### What Went Well
   - 
   
   ### What Went Wrong
   - 
   
   ### Where We Got Lucky
   - 
   
   ## Action Items
   - [ ] Item 1 (Owner: @person, Due: YYYY-MM-DD)
   - [ ] Item 2 (Owner: @person, Due: YYYY-MM-DD)
   ```

5. **Error Budget Tracking Dashboard**:

   **Grafana dashboard** com:
   - Error budget remaining (gauge)
   - Burn rate graph (30d, 7d, 24h, 1h)
   - Time until budget exhausted (projection)
   - Historical SLO compliance
   - Incident timeline

Requisitos:
- SLIs medindo availability, latency, throughput
- SLOs definidos com targets claros
- Error budget calculado e tracked
- Multi-window burn rate alerting
- Runbooks documentados
- Postmortem templates
- Alertmanager configurado
```

### Critérios de Aceite

- [ ] SLIs estão sendo medidos
- [ ] SLOs definidos e documentados
- [ ] Error budget tracking funciona
- [ ] Alerting rules configuradas
- [ ] Alerts testados e funcionando
- [ ] Runbooks escritos e acessíveis
- [ ] Alertmanager roteia corretamente
- [ ] Dashboard SLO criado
- [ ] Postmortem template pronto

### Comandos de Teste

```bash
# Verificar alerting rules no Prometheus
curl http://localhost:9090/api/v1/rules | jq '.data.groups[].rules[] | select(.type=="alerting")'

# Ver alerts ativos
open http://localhost:9090/alerts

# Simular violação de SLO (gerar erros)
docker-compose exec product-api bash -c '
export ENABLE_CHAOS=true
export ERROR_RATE=0.2
python -c "import uvicorn; uvicorn.run()"
'

# Gerar carga
for i in {1..1000}; do
  curl -X POST http://localhost:8000/products \
    -H "Content-Type: application/json" \
    -d '{"name":"Test","price":10.00}' || true
done

# Aguardar alert firing
sleep 120

# Verificar Alertmanager
open http://localhost:9093

# Calcular error budget atual
docker-compose exec prometheus promtool query instant \
  'http://localhost:9090' \
  '1 - (sum(rate(products_created_total{status!~"5.."}[30d])) / sum(rate(products_created_total[30d])))'

# Ver dashboard SLO no Grafana
open http://localhost:3000/d/slo-dashboard

# Testar webhook de alerting
curl -X POST http://localhost:8080/webhook \
  -H 'Content-Type: application/json' \
  -d '{
    "alerts": [{
      "status": "firing",
      "labels": {"alertname": "TestAlert", "severity": "warning"},
      "annotations": {"summary": "Test alert"}
    }]
  }'
```

### Entregáveis
- SLO definitions (YAML)
- Prometheus alerting rules
- Alertmanager configuration
- Runbooks completos
- Postmortem templates
- Error budget dashboard
- Documentação SRE

---

## �📊 Resumo das Fases

| Fase | Foco | Duração Estimada | Complexidade | Relevância Wolt |
|------|------|------------------|--------------|------------------|
| 1 | Infraestrutura (Docker Compose) | 2-3 horas | ⭐⭐ | Distributed systems |
| 2 | API Python (FastAPI + OTel) | 4-5 horas | ⭐⭐⭐ | Python, instrumentation |
| 3 | Kafka Producer (Python) | 2-3 horas | ⭐⭐⭐ | Event streaming |
| 4 | Estrutura Go (Base + OTel) | 3-4 horas | ⭐⭐⭐ | **Go (preferred)** |
| 5 | Repository PostgreSQL (Go) | 3-4 horas | ⭐⭐⭐ | Distributed systems |
| 6 | Kafka Consumer (Go completo) | 4-5 horas | ⭐⭐⭐⭐ | **Go + Kafka** |
| 7 | Chaos Engineering | 2-3 horas | ⭐⭐⭐ | SRE practices |
| 8 | Testes e Documentação | 4-5 horas | ⭐⭐⭐ | Documentation |
| **9** | **OpenTelemetry + Collector** | **4-6 horas** | **⭐⭐⭐⭐** | **OpenTelemetry** |
| **10** | **ClickHouse + Traces/Logs** | **5-7 horas** | **⭐⭐⭐⭐⭐** | **ClickHouse at scale** |
| **11** | **Prometheus + Grafana** | **4-5 horas** | **⭐⭐⭐⭐** | **Metrics + Viz** |
| **12** | **SLIs/SLOs + Alerting** | **3-4 horas** | **⭐⭐⭐⭐** | **SRE principles** |

**Total estimado**: 40-54 horas (projeto completo de portfólio)

---

## 🎓 Conceitos Aprendidos por Fase

### Fase 1
- Docker Compose multi-container
- Networking entre containers
- Health checks
- Volumes e persistência
- Variáveis de ambiente

### Fase 2
- FastAPI e async Python
- Pydantic validation
- REST API best practices
- Structured logging
- Docker multi-stage builds

### Fase 3
- Kafka producer patterns
- Event-driven architecture
- Retry e backoff strategies
- Error handling assíncrono
- Idempotência

### Fase 4 & 5
- Go project layout
- Clean architecture em Go
- Repository pattern
- Connection pooling
- Context em Go
- Error wrapping

### Fase 6
- Kafka consumer groups
- At-least-once delivery
- Manual commit strategies
- Graceful shutdown
- Dead letter queues
- Idempotência em consumers

### Fase 7
- Chaos Engineering
- Fault injection
- Resiliency patterns
- Health vs Ready checks
- Circuit breakers (conceito)

### Fase 8
- Test-driven development
- Integration testing
- Mocking strategies
- Coverage analysis
- Documentation as code
- CI/CD concepts

---

## 🚀 Fases Avançadas de Observabilidade (9-12)

As próximas 4 fases implementam a **stack de observabilidade completa**, demonstrando as competências principais da vaga:

### **Fase 9: OpenTelemetry Collector + Instrumentação**
- Instrumentar microserviços com OpenTelemetry SDK
- Configurar OTLP Collector para receber telemetria
- Trace context propagation entre serviços
- Exporters para múltiplos backends
- Sampling strategies

### **Fase 10: ClickHouse para Traces e Logs em Escala**
- Deploy ClickHouse como backend de observabilidade
- Schema otimizado para traces e logs
- Queries de análise de latência e erros
- Agregações para insights de performance
- Retention policies

### **Fase 11: Prometheus + Grafana**
- Exportar métricas de aplicação (RED metrics)
- Prometheus scraping e PromQL
- Grafana dashboards customizados
- Visualização unificada (metrics + traces + logs)
- Correlação de dados

### **Fase 12: SRE - SLIs/SLOs + Alerting**
- Definir SLIs (latência, erro rate, disponibilidade)
- Configurar SLOs com error budgets
- Alerting baseado em SLOs
- Runbooks de incident response
- Postmortem templates

---

## 📚 Recursos e Referências

### Documentação Oficial
- [FastAPI](https://fastapi.tiangolo.com/)
- [Kafka](https://kafka.apache.org/documentation/)
- [Go](https://go.dev/doc/)
- [PostgreSQL](https://www.postgresql.org/docs/)
- [Docker Compose](https://docs.docker.com/compose/)

### Bibliotecas
- [kafka-python](https://kafka-python.readthedocs.io/)
- [segmentio/kafka-go](https://github.com/segmentio/kafka-go)
- [pgx](https://github.com/jackc/pgx)
- [Pydantic](https://docs.pydantic.dev/)

### Padrões e Práticas
- [12 Factor App](https://12factor.net/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Microservices Patterns](https://microservices.io/patterns/)

---

## 💡 Dicas Importantes

1. **Commit frequentemente**: Faça commits ao final de cada fase
2. **Teste antes de avançar**: Valide cada fase completamente
3. **Documente enquanto desenvolve**: Não deixe para depois
4. **Use logs estruturados desde o início**: Facilita debugging
5. **Pense em observabilidade**: Prepare código para métricas/traces
6. **Erro é aprendizado**: Use os erros injetados para aprender
7. **Consulte documentação**: Não tenha medo de ler docs oficiais
8. **Peça ajuda**: Use este plano como guia, não como restrição

---

## 🎯 Checklist Final - Competências Demonstradas

Ao concluir todas as 12 fases, você terá um **portfólio completo** demonstrando:

### Requisitos Principais da Vaga ✅
- [ ] **Go proficiency**: Consumer completo em Go com best practices
- [ ] **Python**: API FastAPI com instrumentação
- [ ] **OpenTelemetry**: Instrumentação manual e automática
- [ ] **ClickHouse**: Backend de observabilidade em escala
- [ ] **Kafka**: Event streaming e processamento
- [ ] **Prometheus + Grafana**: Métricas e visualização
- [ ] **SRE principles**: SLIs/SLOs, incident response
- [ ] **Distributed systems**: Arquitetura multi-serviço

### Competências Técnicas ✅
- [ ] Observability platform architecture
- [ ] Telemetry data collection e processamento
- [ ] Scalable software solutions
- [ ] Automation and developer tooling
- [ ] Container orchestration (Docker)
- [ ] Troubleshooting complex systems
- [ ] Documentation e knowledge sharing

### Nice to Haves ✅
- [ ] OpenTelemetry at scale
- [ ] ClickHouse experience
- [ ] Kafka event streaming
- [ ] Chaos engineering
- [ ] Infrastructure as Code

### Entregáveis do Portfólio 📦
- [ ] Repositório GitHub público com código
- [ ] README detalhado com arquitetura
- [ ] Dashboards Grafana exportáveis
- [ ] Queries ClickHouse para análise
- [ ] Runbooks e documentação SRE
- [ ] Exemplos de traces, metrics, logs
- [ ] Screenshots de dashboards
- [ ] Vídeo demo (opcional, mas recomendado)

---

## 💼 Destacando o Projeto no Currículo/GitHub

### Título Sugerido
**"Scalable Observability Platform with OpenTelemetry, ClickHouse & Kafka"**

### Descrição
```
Full-stack observability platform demonstrating SRE practices and large-scale 
telemetry processing. Implements distributed tracing, metrics collection, and 
log aggregation using OpenTelemetry, ClickHouse, Prometheus, and Grafana.

Highlights:
• Go + Python microservices with full OpenTelemetry instrumentation
• ClickHouse backend processing billions of telemetry data points
• Kafka event streaming with exactly-once semantics
• Prometheus metrics + Grafana dashboards with SLO tracking
• Chaos engineering for reliability testing
• Complete SRE practices: SLIs/SLOs, runbooks, incident response
```

### Tags/Keywords
`observability` `opentelemetry` `clickhouse` `prometheus` `grafana` `go` `python` 
`kafka` `sre` `distributed-systems` `microservices` `devops` `monitoring`

---

**Boa sorte com o processo seletivo! 🚀**
