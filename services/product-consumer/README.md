# Product Consumer

Go-based Kafka consumer that processes product events and persists them to PostgreSQL with comprehensive observability.

## Features

- **Clean Architecture**: Repository pattern with interface-based design
- **PostgreSQL Integration**: Connection pooling with pgx/v5
- **Kafka Consumer**: kafka-go with consumer groups
- **Graceful Shutdown**: Proper cleanup of resources
- **Health Checks**: Dedicated HTTP server for health monitoring
- **Idempotent Processing**: Handles duplicate messages
- **DLQ Support**: Dead Letter Queue for failed messages
- **Chaos Engineering**: Optional error and latency injection

## Health Endpoints

Dedicated health server on port **8081**:

- `GET /health` - Service health (always returns 200 when running)
- `GET /ready` - Readiness check (validates PostgreSQL and Kafka connectivity)

**Ready Check Response:**
```json
{
  "status": "ready",
  "postgres": "connected",
  "kafka": "connected"
}
```
- Static binary compilation for minimal Docker image

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `KAFKA_BOOTSTRAP_SERVERS` | Kafka broker addresses | `localhost:9092` |
| `KAFKA_TOPIC_PRODUCTS` | Kafka topic to consume | `products.events` |
| `KAFKA_GROUP_ID` | Kafka consumer group ID | `product-consumer-group` |
| `POSTGRES_HOST` | PostgreSQL host | `localhost` |
| `POSTGRES_PORT` | PostgreSQL port | `5432` |
| `POSTGRES_USER` | PostgreSQL user | `postgres` |
| `POSTGRES_PASSWORD` | PostgreSQL password | `postgres` |
| `POSTGRES_DB` | PostgreSQL database | `products_db` |
| `POSTGRES_SSLMODE` | PostgreSQL SSL mode | `disable` |
| `LOG_LEVEL` | Logging level (DEBUG, INFO, WARN, ERROR) | `INFO` |
| `SHUTDOWN_TIMEOUT` | Graceful shutdown timeout | `30s` |

## Build and Run

### Local Development

```bash
# Install dependencies
go mod download

# Run locally
cp .env.example .env
go run cmd/main.go

# Build binary
go build -o bin/consumer cmd/main.go
./bin/consumer

# Run with custom env vars
KAFKA_BOOTSTRAP_SERVERS=kafka:9093 \
POSTGRES_HOST=postgres \
LOG_LEVEL=DEBUG \
./bin/consumer
```

### Docker

```bash
# Build image
docker build -t product-consumer:latest .

# Run container
docker run --rm \
  -e KAFKA_BOOTSTRAP_SERVERS=kafka:9093 \
  -e POSTGRES_HOST=postgres \
  product-consumer:latest
```

### Docker Compose

The service is configured in the main `docker-compose.yml`:

```bash
docker compose up -d product-consumer
```

## Development

### Code Quality

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run linter (requires golangci-lint)
golangci-lint run

# Tidy dependencies
go mod tidy
```

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Dependencies

- **kafka-go**: Kafka client for Go
- **pgx/v5**: PostgreSQL driver and toolkit
- **godotenv**: Load environment variables from .env
- **logrus**: Structured logging

## Logging Format

All logs are output in JSON format:

```json
{
  "timestamp": "2026-02-13T12:00:00Z",
  "level": "info",
  "message": "Application starting",
  "kafka_brokers": ["localhost:9092"],
  "kafka_topic": "products.events"
}
```

## Graceful Shutdown

The application handles SIGTERM and SIGINT signals for graceful shutdown:

1. Stop accepting new messages
2. Process in-flight messages
3. Close database connections
4. Close Kafka consumer
5. Exit with appropriate status code

Shutdown timeout is configurable via `SHUTDOWN_TIMEOUT` environment variable.

## Next Steps

- [ ] Implement Kafka consumer
- [ ] Implement PostgreSQL repository
- [ ] Add retry mechanism for failed messages
- [ ] Add metrics collection
- [ ] Add distributed tracing
- [ ] Implement dead letter queue
