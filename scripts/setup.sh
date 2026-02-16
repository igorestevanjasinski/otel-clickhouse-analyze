#!/bin/bash
set -e

echo "=================================================="
echo "  Environment Setup - Product Portfolio"
echo "=================================================="

cd "$(dirname "$0")/.."

echo ""
echo "[1/6] Checking prerequisites..."
command -v docker >/dev/null 2>&1 || { echo "Error: Docker not found"; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "Error: docker-compose not found"; exit 1; }
command -v python3 >/dev/null 2>&1 || { echo "Error: Python3 not found"; exit 1; }
command -v go >/dev/null 2>&1 || { echo "Error: Go not found"; exit 1; }
echo "✓ All prerequisites met"

echo ""
echo "[2/6] Starting infrastructure..."
if docker ps | grep -q postgres; then
    echo "✓ Infrastructure already running"
else
    docker-compose up -d
    echo "⏳ Waiting for services to be ready..."
    sleep 10
    echo "✓ Infrastructure started"
fi

echo ""
echo "[3/6] Setting up Python virtual environment..."
if [ -d "services/product-api/venv" ]; then
    echo "✓ Virtual environment already exists"
else
    cd services/product-api
    python3 -m venv venv
    source venv/bin/activate
    pip install -r requirements.txt
    pip install -r requirements-dev.txt
    deactivate
    cd ../..
    echo "✓ Python environment created"
fi

echo ""
echo "[4/6] Creating PostgreSQL schema..."
docker exec -i postgres psql -U postgres -d products -c "
CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL
);
" 2>/dev/null || echo "Schema already exists"
echo "✓ Database schema ready"

echo ""
echo "[5/6] Compiling Go consumer..."
cd services/product-consumer
if [ ! -f "bin/consumer" ]; then
    go mod download
    go build -o bin/consumer cmd/main.go
    echo "✓ Consumer compiled"
else
    echo "✓ Consumer already compiled"
fi
cd ../..

echo ""
echo "[6/6] Creating .env files if missing..."
if [ ! -f ".env" ]; then
    cp .env.example .env
    echo "✓ .env created from template"
else
    echo "✓ .env already exists"
fi

echo ""
echo "=================================================="
echo "  Setup Complete!"
echo "=================================================="
echo ""
echo "Next steps:"
echo "  - Run tests: ./scripts/run-tests.sh"
echo "  - Start API: cd services/product-api && source venv/bin/activate && uvicorn app.main:app --reload"
echo "  - Start Consumer: cd services/product-consumer && ./bin/consumer"
echo "  - Run E2E: ./scripts/test-e2e.sh"
echo ""
