#!/bin/bash
set -e

echo "=================================================="
echo "  End-to-End Test"
echo "=================================================="

cd "$(dirname "$0")/.."

API_PORT=8000
CONSUMER_HEALTH_PORT=8081
API_URL="http://localhost:${API_PORT}"
CONSUMER_URL="http://localhost:${CONSUMER_HEALTH_PORT}"

cleanup() {
    echo ""
    echo "Cleaning up test processes..."
    pkill -f "uvicorn app.main:app" 2>/dev/null || true
    pkill -f "bin/consumer" 2>/dev/null || true
}

trap cleanup EXIT

echo ""
echo "[1/7] Checking if infrastructure is running..."
if ! docker ps | grep -q postgres; then
    echo "Starting infrastructure..."
    docker-compose up -d
    echo "⏳ Waiting for services to be ready..."
    sleep 10
else
    echo "✓ Infrastructure already running"
fi

echo ""
echo "[2/7] Starting Product API..."
cd services/product-api
if [ -d "venv" ]; then
    source venv/bin/activate
    ENABLE_CHAOS=false uvicorn app.main:app --host 0.0.0.0 --port 8000 > /tmp/api.log 2>&1 &
    API_PID=$!
    deactivate
    sleep 3
    echo "✓ API started (PID: $API_PID)"
else
    echo "✗ Python venv not found. Run ./scripts/setup.sh first"
    exit 1
fi
cd ../..

echo ""
echo "[3/7] Starting Product Consumer..."
cd services/product-consumer
if [ -f "bin/consumer" ]; then
    ENABLE_CHAOS=false ./bin/consumer > /tmp/consumer.log 2>&1 &
    CONSUMER_PID=$!
    sleep 3
    echo "✓ Consumer started (PID: $CONSUMER_PID)"
else
    echo "✗ Consumer binary not found. Run ./scripts/setup.sh first"
    exit 1
fi
cd ../..

echo ""
echo "[4/7] Health checks..."
HEALTH=$(curl -s "${API_URL}/health" | jq -r '.status')
if [ "$HEALTH" = "healthy" ]; then
    echo "✓ API health: ${HEALTH}"
else
    echo "✗ API health check failed"
    exit 1
fi

CONSUMER_READY=$(curl -s "${CONSUMER_URL}/ready" | jq -r '.status')
if [ "$CONSUMER_READY" = "ready" ]; then
    echo "✓ Consumer ready: ${CONSUMER_READY}"
else
    echo "✗ Consumer ready check failed"
    exit 1
fi

echo ""
echo "[5/7] Creating test product..."
RESPONSE=$(curl -s -X POST "${API_URL}/products" \
  -H "Content-Type: application/json" \
  -H "X-Correlation-ID: e2e-test-$(date +%s)" \
  -d '{
    "name": "E2E Test Product",
    "price": 149.99,
    "description": "Product created during E2E test"
  }')

echo "Response: $RESPONSE"

PRODUCT_ID=$(echo $RESPONSE | jq -r '.id')
if [ -z "$PRODUCT_ID" ] || [ "$PRODUCT_ID" = "null" ]; then
    echo "✗ Failed to create product"
    exit 1
fi
echo "✓ Created product with ID: $PRODUCT_ID"

echo ""
echo "[6/7] Waiting for consumer to process message..."
sleep 5

echo ""
echo "[7/7] Verifying in database..."
DB_RESULT=$(docker exec -i postgres psql -U postgres -d products -t -c "SELECT COUNT(*) FROM products WHERE id = '$PRODUCT_ID';")
DB_COUNT=$(echo $DB_RESULT | tr -d ' ')

if [ "$DB_COUNT" = "1" ]; then
    echo "✓ Product found in database"
    docker exec -i postgres psql -U postgres -d products -c "SELECT * FROM products WHERE id = '$PRODUCT_ID';"
else
    echo "✗ Product not found in database (count: $DB_COUNT)"
    echo ""
    echo "API logs:"
    tail -n 20 /tmp/api.log
    echo ""
    echo "Consumer logs:"
    tail -n 20 /tmp/consumer.log
    exit 1
fi

echo ""
echo "=================================================="
echo "  E2E Test Passed!"
echo "=================================================="
echo ""
echo "Logs available at:"
echo "  - API: /tmp/api.log"
echo "  - Consumer: /tmp/consumer.log"
echo ""

echo ""
echo "7. Total products in database..."
docker compose exec -T postgres psql -U postgres -d products_db \
  -c "SELECT COUNT(*) as total FROM products;"

echo ""
echo "========================================="
echo "All tests completed successfully!"
echo "========================================="
