#!/bin/bash

set -e

echo "========================================="
echo "Testing Phase 11: Prometheus + Grafana"
echo "========================================="
echo ""

echo "Step 1: Starting Prometheus and Grafana..."
docker compose up -d prometheus grafana

echo ""
echo "Step 2: Waiting for services to be healthy..."
sleep 15

echo ""
echo "Step 3: Checking Prometheus health..."
if curl -s http://localhost:9090/-/healthy | grep -q "Prometheus is Healthy"; then
    echo "✓ Prometheus is healthy"
else
    echo "✗ Prometheus health check failed"
    exit 1
fi

echo ""
echo "Step 4: Checking Grafana health..."
if curl -s http://localhost:3000/api/health | grep -q "ok"; then
    echo "✓ Grafana is healthy"
else
    echo "✗ Grafana health check failed"
    exit 1
fi

echo ""
echo "Step 5: Verifying Prometheus targets..."
TARGETS=$(curl -s http://localhost:9090/api/v1/targets | jq -r '.data.activeTargets[] | select(.health == "up") | .labels.job')
echo "Active targets:"
echo "$TARGETS"

if echo "$TARGETS" | grep -q "product-api"; then
    echo "✓ product-api target is up"
else
    echo "✗ product-api target is down"
fi

if echo "$TARGETS" | grep -q "product-consumer"; then
    echo "✓ product-consumer target is up"
else
    echo "✗ product-consumer target is down"
fi

if echo "$TARGETS" | grep -q "otel-collector"; then
    echo "✓ otel-collector target is up"
else
    echo "✗ otel-collector target is down"
fi

echo ""
echo "Step 6: Checking Product API metrics endpoint..."
if curl -s http://localhost:8000/metrics | grep -q "products_created_total"; then
    echo "✓ Product API /metrics endpoint is working"
    echo "  Sample metrics:"
    curl -s http://localhost:8000/metrics | grep "products_created_total" | head -3
else
    echo "✗ Product API /metrics endpoint failed"
fi

echo ""
echo "Step 7: Checking Product Consumer metrics endpoint..."
if curl -s http://localhost:8081/metrics | grep -q "messages_consumed_total"; then
    echo "✓ Product Consumer /metrics endpoint is working"
    echo "  Sample metrics:"
    curl -s http://localhost:8081/metrics | grep "messages_consumed_total" | head -3
else
    echo "✗ Product Consumer /metrics endpoint failed"
fi

echo ""
echo "Step 8: Generating load for metrics..."
echo "Creating 100 products..."
for i in {1..100}; do
    curl -s -X POST http://localhost:8000/products \
        -H "Content-Type: application/json" \
        -d "{\"name\":\"Grafana Test Product $i\",\"price\":$(( ( RANDOM % 100 ) + 1 )).99,\"description\":\"Testing Grafana dashboards\"}" \
        > /dev/null &
    
    if [ $(( i % 20 )) -eq 0 ]; then
        echo "  Created $i products..."
        wait
    fi
done
wait

echo "✓ Load generation complete"

echo ""
echo "Step 9: Waiting for metrics to be scraped..."
sleep 20

echo ""
echo "Step 10: Verifying metrics in Prometheus..."
PRODUCT_COUNT=$(curl -s 'http://localhost:9090/api/v1/query?query=products_created_total' | jq -r '.data.result[0].value[1]')
if [ ! -z "$PRODUCT_COUNT" ] && [ "$PRODUCT_COUNT" != "null" ]; then
    echo "✓ products_created_total: $PRODUCT_COUNT"
else
    echo "✗ Failed to query products_created_total"
fi

MESSAGE_COUNT=$(curl -s 'http://localhost:9090/api/v1/query?query=messages_consumed_total' | jq -r '.data.result[0].value[1]')
if [ ! -z "$MESSAGE_COUNT" ] && [ "$MESSAGE_COUNT" != "null" ]; then
    echo "✓ messages_consumed_total: $MESSAGE_COUNT"
else
    echo "✗ Failed to query messages_consumed_total"
fi

echo ""
echo "Step 11: Testing recording rules..."
RULES=$(curl -s http://localhost:9090/api/v1/rules | jq -r '.data.groups[].rules[].name')
if [ ! -z "$RULES" ]; then
    echo "✓ Recording rules loaded:"
    echo "$RULES"
else
    echo "✗ No recording rules found"
fi

echo ""
echo "Step 12: Verifying Grafana datasources..."
DATASOURCES=$(curl -s -u admin:admin http://localhost:3000/api/datasources | jq -r '.[].name')
if echo "$DATASOURCES" | grep -q "Prometheus"; then
    echo "✓ Prometheus datasource configured"
else
    echo "✗ Prometheus datasource missing"
fi

if echo "$DATASOURCES" | grep -q "ClickHouse"; then
    echo "✓ ClickHouse datasource configured"
else
    echo "✗ ClickHouse datasource missing"
fi

echo ""
echo "Step 13: Verifying Grafana dashboards..."
DASHBOARDS=$(curl -s -u admin:admin http://localhost:3000/api/search?type=dash-db | jq -r '.[].title')
if echo "$DASHBOARDS" | grep -q "Microservices Overview"; then
    echo "✓ Microservices Overview dashboard loaded"
else
    echo "✗ Microservices Overview dashboard missing"
fi

if echo "$DASHBOARDS" | grep -q "Distributed Tracing"; then
    echo "✓ Distributed Tracing dashboard loaded"
else
    echo "✗ Distributed Tracing dashboard missing"
fi

echo ""
echo "========================================="
echo "Phase 11 Testing Complete!"
echo "========================================="
echo ""
echo "Access the dashboards:"
echo "  Prometheus: http://localhost:9090"
echo "  Grafana:    http://localhost:3000 (admin/admin)"
echo ""
echo "Dashboards available:"
echo "  - Microservices Overview - RED Metrics"
echo "  - Distributed Tracing - ClickHouse"
echo ""
