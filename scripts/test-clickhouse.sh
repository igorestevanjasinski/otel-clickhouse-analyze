#!/bin/bash

set -e

echo "========================================="
echo "Fase 10: ClickHouse Testing Guide"
echo "========================================="
echo ""

echo "Step 1: Starting infrastructure..."
docker compose up -d clickhouse
echo "Waiting for ClickHouse to be healthy..."
sleep 10

echo ""
echo "Step 2: Verifying ClickHouse..."
docker compose exec clickhouse clickhouse-client --query "SELECT version()"

echo ""
echo "Step 3: Checking databases..."
docker compose exec clickhouse clickhouse-client --query "SHOW DATABASES"

echo ""
echo "Step 4: Checking tables..."
docker compose exec clickhouse clickhouse-client --query "SHOW TABLES FROM observability"

echo ""
echo "Step 5: Starting all services..."
docker compose up -d

echo ""
echo "Waiting for services to be ready..."
sleep 15

echo ""
echo "Step 6: Generating test data (100 requests)..."
for i in {1..100}; do
  curl -s -X POST http://localhost:8000/products \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"Product $i\",\"price\":$i.00,\"description\":\"Test product for ClickHouse\"}" > /dev/null
  
  if [ $((i % 10)) -eq 0 ]; then
    echo "Generated $i requests..."
  fi
  
  sleep 0.1
done

echo ""
echo "Step 7: Waiting for data export to ClickHouse..."
sleep 10

echo ""
echo "Step 8: Verifying data in ClickHouse..."
echo ""
echo "=== Total Traces ==="
docker compose exec clickhouse clickhouse-client --query="
SELECT count() as total_traces FROM observability.traces
"

echo ""
echo "=== Total Logs ==="
docker compose exec clickhouse clickhouse-client --query="
SELECT count() as total_logs FROM observability.logs
"

echo ""
echo "=== Traces by Service ==="
docker compose exec clickhouse clickhouse-client --query="
SELECT 
    service_name,
    count() as trace_count
FROM observability.traces
GROUP BY service_name
FORMAT PrettyCompact
"

echo ""
echo "Step 9: Running analysis script..."
./scripts/analyze-traces.sh

echo ""
echo "========================================="
echo "Testing completed successfully!"
echo "========================================="
echo ""
echo "Next steps:"
echo "1. Check Grafana dashboards (Phase 11)"
echo "2. Review ClickHouse documentation: docs/CLICKHOUSE.md"
echo "3. Explore custom queries in: queries/clickhouse-analysis.sql"
echo ""
echo "Useful commands:"
echo "  - View traces: docker compose exec clickhouse clickhouse-client"
echo "  - Analysis: ./scripts/analyze-traces.sh"
echo "  - Logs: docker compose logs clickhouse"
