#!/bin/bash

set -e

echo "=================================================="
echo "  OpenTelemetry Instrumentation Test"
echo "=================================================="

echo ""
echo "[1/6] Starting observability stack..."
docker compose -f docker-compose.observability.yml up -d

echo "Waiting for OpenTelemetry Collector to be ready..."
sleep 5

COLLECTOR_READY=false
for i in {1..30}; do
    if curl -s http://localhost:13133/ > /dev/null 2>&1; then
        COLLECTOR_READY=true
        break
    fi
    echo -n "."
    sleep 1
done

if [ "$COLLECTOR_READY" = false ]; then
    echo ""
    echo "ERROR: OpenTelemetry Collector failed to start"
    docker compose -f docker-compose.observability.yml logs otel-collector
    exit 1
fi

echo ""
echo "✓ OpenTelemetry Collector is ready"

echo ""
echo "[2/6] Installing Python dependencies..."
cd services/product-api
if [ ! -d "venv" ]; then
    python3 -m venv venv
fi
source venv/bin/activate
pip install -q -r requirements.txt
cd ../..

echo "✓ Python dependencies installed"

echo ""
echo "[3/6] Building Go consumer..."
cd services/product-consumer
go mod download
go mod tidy
go build -o bin/consumer cmd/main.go
cd ../..

echo "✓ Go consumer built"

echo ""
echo "[4/6] Starting Product API..."
cd services/product-api
source venv/bin/activate
OTLP_ENDPOINT=http://localhost:4317 uvicorn app.main:app --host 0.0.0.0 --port 8000 > ../../api-telemetry.log 2>&1 &
API_PID=$!
cd ../..

echo "API started (PID: $API_PID)"
sleep 3

if ! ps -p $API_PID > /dev/null; then
    echo "ERROR: API failed to start"
    cat api-telemetry.log
    exit 1
fi

echo "✓ API is running"

echo ""
echo "[5/6] Starting Product Consumer..."
cd services/product-consumer
OTLP_ENDPOINT=localhost:4317 ./bin/consumer > ../../consumer-telemetry.log 2>&1 &
CONSUMER_PID=$!
cd ../..

echo "Consumer started (PID: $CONSUMER_PID)"
sleep 3

if ! ps -p $CONSUMER_PID > /dev/null; then
    echo "ERROR: Consumer failed to start"
    cat consumer-telemetry.log
    kill $API_PID 2>/dev/null || true
    exit 1
fi

echo "✓ Consumer is running"

echo ""
echo "[6/6] Testing telemetry generation..."

echo ""
echo "Creating test product..."
RESPONSE=$(curl -s -X POST http://localhost:8000/products \
  -H "Content-Type: application/json" \
  -H "X-Correlation-ID: test-telemetry-$(date +%s)" \
  -d '{
    "name": "Telemetry Test Product",
    "price": 199.99,
    "description": "Testing OpenTelemetry instrumentation"
  }')

PRODUCT_ID=$(echo $RESPONSE | jq -r '.id')
echo "Created product: $PRODUCT_ID"

echo ""
echo "Waiting for message processing..."
sleep 5

echo ""
echo "=================================================="
echo "  Checking OpenTelemetry Collector Logs"
echo "=================================================="
echo ""
echo "Last 50 lines of collector logs:"
docker compose -f docker-compose.observability.yml logs --tail=50 otel-collector

echo ""
echo "=================================================="
echo "  Validation"
echo "=================================================="

echo ""
echo "Checking for traces in collector logs..."
if docker compose -f docker-compose.observability.yml logs otel-collector | grep -q "Span #"; then
    echo "✓ Traces found in collector logs"
    TRACES_FOUND=true
else
    echo "✗ No traces found in collector logs"
    TRACES_FOUND=false
fi

echo ""
echo "Checking for metrics in collector logs..."
if docker compose -f docker-compose.observability.yml logs otel-collector | grep -q "Metric #"; then
    echo "✓ Metrics found in collector logs"
    METRICS_FOUND=true
else
    echo "✗ No metrics found in collector logs"
    METRICS_FOUND=false
fi

echo ""
echo "Checking for resource attributes..."
if docker compose -f docker-compose.observability.yml logs otel-collector | grep -q "service.name"; then
    echo "✓ Resource attributes found"
    RESOURCES_FOUND=true
else
    echo "✗ No resource attributes found"
    RESOURCES_FOUND=false
fi

echo ""
echo "=================================================="
echo "  Summary"
echo "=================================================="
echo ""
echo "Services:"
echo "  - API PID: $API_PID (http://localhost:8000)"
echo "  - Consumer PID: $CONSUMER_PID"
echo "  - Collector: http://localhost:4317"
echo "  - Jaeger UI: http://localhost:16686"
echo ""
echo "Telemetry Status:"
echo "  - Traces: $([ "$TRACES_FOUND" = true ] && echo '✓' || echo '✗')"
echo "  - Metrics: $([ "$METRICS_FOUND" = true ] && echo '✓' || echo '✗')"
echo "  - Resources: $([ "$RESOURCES_FOUND" = true ] && echo '✓' || echo '✗')"
echo ""
echo "Logs:"
echo "  - API: tail -f api-telemetry.log"
echo "  - Consumer: tail -f consumer-telemetry.log"
echo "  - Collector: docker compose -f docker-compose.observability.yml logs -f otel-collector"
echo ""
echo "Commands:"
echo "  - View traces: docker compose -f docker-compose.observability.yml logs otel-collector | grep 'Span #'"
echo "  - View metrics: docker compose -f docker-compose.observability.yml logs otel-collector | grep 'Metric #'"
echo "  - Open Jaeger: open http://localhost:16686"
echo ""
echo "Cleanup:"
echo "  kill $API_PID $CONSUMER_PID"
echo "  docker compose -f docker-compose.observability.yml down"
echo ""

if [ "$TRACES_FOUND" = true ] && [ "$RESOURCES_FOUND" = true ]; then
    echo "=================================================="
    echo "  ✓ Telemetry Validation Successful!"
    echo "=================================================="
    exit 0
else
    echo "=================================================="
    echo "  ✗ Telemetry Validation Failed"
    echo "=================================================="
    echo ""
    echo "Check logs for errors:"
    echo "  - cat api-telemetry.log"
    echo "  - cat consumer-telemetry.log"
    exit 1
fi
