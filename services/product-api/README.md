# Product API

FastAPI microservice for product creation with Kafka event publishing and comprehensive observability features.

## Features

- **FastAPI Framework**: High-performance async web framework
- **Pydantic Validation**: Strict input validation with custom validators
- **Kafka Integration**: Event-driven architecture with confluent-kafka
- **Structured Logging**: JSON logging with correlation IDs
- **Health Checks**: `/health` and `/ready` endpoints
- **Chaos Engineering**: Optional error and latency injection for testing

## API Endpoints

### Health Checks

- `GET /health` - Service health status
- `GET /ready` - Readiness check (validates Kafka connectivity)

### Products

- `POST /products` - Create a new product

**Request Body:**
```json
{
  "name": "Product Name",
  "price": 99.99,
  "description": "Optional description"
}
```bash
cd services/product-api

python -m venv venv
source venv/bin/activate

pip install -r requirements.txt
```

### Environment Variables

Copy `.env.example` to `.env` in the project root:

```bash
cp ../../.env.example ../../.env
```

### Run Locally

```bash
uvicorn app.main:app --reload --host 0.0.0.0 --port 8000
```

The API will be available at http://localhost:8000

## API Endpoints

### POST /products

Create a new product.

**Request:**
```json
{
  "name": "Laptop Dell XPS 13",
  "price": 1299.99,
  "description": "High-performance ultrabook"
}
```

**Response (202 Accepted):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Laptop Dell XPS 13",
  "price": 1299.99,
  "description": "High-performance ultrabook",
  "created_at": "2024-01-01T12:00:00Z"
}
```

**Validations:**
- `name`: 3-255 characters, non-empty
- `price`: Greater than 0, max 2 decimal places
- `description`: Optional, max 1000 characters

### GET /health

Health check endpoint.

**Response (200 OK):**
```json
{
  "status": "healthy"
}
```

## Documentation

Interactive API documentation is available at:

- Swagger UI: http://localhost:8000/docs
- ReDoc: http://localhost:8000/redoc

## Docker

### Build

```bash
docker build -t product-api:latest .
```

### Run

```bash
docker run -p 8000:8000 \
  -e KAFKA_BOOTSTRAP_SERVERS=kafka:9092 \
  -e LOG_LEVEL=INFO \
  product-api:latest
```

## Testing

### Manual Testing

```bash
# Health check
curl http://localhost:8000/health

# Create product
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Wireless Mouse",
    "price": 49.99,
    "description": "Ergonomic wireless mouse"
  }'

# Test validation (invalid price)
curl -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Invalid Product",
    "price": -10,
    "description": "Should fail"
  }'
```

## Logging

The application uses structured JSON logging. All logs include:

- `timestamp`: ISO 8601 format
- `level`: Log level (INFO, ERROR, etc)
- `logger`: Logger name
- `message`: Log message
- Additional context fields

Example log:
```json
{
  "timestamp": "2024-01-01T12:00:00.000Z",
  "level": "INFO",
  "logger": "app.api.routes",
  "message": "Product creation initiated",
  "product_id": "123e4567-e89b-12d3-a456-426614174000",
  "product_name": "Laptop Dell XPS 13",
  "product_price": 1299.99
}
```

## Next Steps (Phase 3)

- Integrate Kafka producer
- Publish product events to Kafka topic
- Add retry logic for message delivery
- Implement correlation IDs for distributed tracing
