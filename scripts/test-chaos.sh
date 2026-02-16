#!/bin/bash

set -e

echo "========================================="
echo "Chaos Engineering Test"
echo "========================================="

echo ""
echo "1. Testing API health checks..."
curl -s http://localhost:8000/health | jq .
curl -s http://localhost:8000/ready | jq .

echo ""
echo "2. Sending requests WITHOUT chaos (normal behavior)..."
for i in {1..5}; do
  echo -n "Request $i: "
  START=$(date +%s%3N)
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8000/products \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"Normal Product $i\",\"price\":100.00,\"description\":\"Test without chaos\"}")
  END=$(date +%s%3N)
  DURATION=$((END - START))
  echo "Status: $STATUS, Duration: ${DURATION}ms"
done

echo ""
echo "3. Enabling chaos (restart services with chaos enabled)..."
echo "Note: You need to manually set ENABLE_CHAOS=true in docker-compose.yml"
echo "or use environment variables when starting services."
echo ""
read -p "Press Enter when chaos is enabled..."

echo ""
echo "4. Sending requests WITH chaos (expect some failures and high latency)..."
SUCCESSES=0
FAILURES=0
TOTAL_LATENCY=0

for i in {1..100}; do
  echo -n "Request $i: "
  START=$(date +%s%3N)
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8000/products \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"Chaos Product $i\",\"price\":100.00,\"description\":\"Test with chaos\"}")
  END=$(date +%s%3N)
  DURATION=$((END - START))
  TOTAL_LATENCY=$((TOTAL_LATENCY + DURATION))
  
  if [ "$STATUS" = "202" ]; then
    echo "Status: $STATUS ✓, Duration: ${DURATION}ms"
    SUCCESSES=$((SUCCESSES + 1))
  else
    echo "Status: $STATUS ✗, Duration: ${DURATION}ms (CHAOS ERROR)"
    FAILURES=$((FAILURES + 1))
  fi
done

echo ""
echo "========================================="
echo "Results Summary:"
echo "========================================="
echo "Total Requests:    20"
echo "Successes:         $SUCCESSES"
echo "Failures:          $FAILURES"
echo "Success Rate:      $(awk "BEGIN {printf \"%.1f\", ($SUCCESSES/20)*100}")%"
echo "Average Latency:   $((TOTAL_LATENCY / 20))ms"
echo "========================================="

echo ""
echo "5. Checking logs for chaos events..."
echo "API chaos logs (last 10):"
docker compose logs product-api 2>/dev/null | grep -i "chaos" | tail -10 || echo "Run with docker compose to see logs"

echo ""
echo "Consumer chaos logs (last 10):"
docker compose logs product-consumer 2>/dev/null | grep -i "chaos" | tail -10 || echo "Run with docker compose to see logs"

echo ""
echo "========================================="
echo "Chaos Engineering Test Complete!"
echo "========================================="
