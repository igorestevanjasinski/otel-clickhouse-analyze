#!/bin/bash

echo "=================================================="
echo "  Product API - Example Requests"
echo "=================================================="

API_URL="http://localhost:8000"

echo ""
echo "[1] Health Check"
echo "--------------------------------------------------"
echo "$ curl ${API_URL}/health"
curl -s ${API_URL}/health | jq .

echo ""
echo "[2] Ready Check"
echo "--------------------------------------------------"
echo "$ curl ${API_URL}/ready"
curl -s ${API_URL}/ready | jq .

echo ""
echo "[3] Create Product - Success"
echo "--------------------------------------------------"
CORRELATION_ID="example-$(date +%s)"
echo "$ curl -X POST ${API_URL}/products \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -H 'X-Correlation-ID: ${CORRELATION_ID}' \\"
echo "  -d '{\"name\": \"Laptop\", \"price\": 1299.99, \"description\": \"High-performance laptop\"}'"
curl -s -X POST ${API_URL}/products \
  -H "Content-Type: application/json" \
  -H "X-Correlation-ID: ${CORRELATION_ID}" \
  -d '{
    "name": "Laptop",
    "price": 1299.99,
    "description": "High-performance laptop"
  }' | jq .

echo ""
echo "[4] Create Product - Invalid (Missing Name)"
echo "--------------------------------------------------"
echo "$ curl -X POST ${API_URL}/products \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{\"price\": 99.99, \"description\": \"Test\"}'"
curl -s -X POST ${API_URL}/products \
  -H "Content-Type: application/json" \
  -d '{
    "price": 99.99,
    "description": "Test"
  }' | jq .

echo ""
echo "[5] Create Product - Invalid (Negative Price)"
echo "--------------------------------------------------"
echo "$ curl -X POST ${API_URL}/products \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{\"name\": \"Test\", \"price\": -10, \"description\": \"Test\"}'"
curl -s -X POST ${API_URL}/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test",
    "price": -10,
    "description": "Test"
  }' | jq .

echo ""
echo "[6] Create Product - Bulk (Multiple Products)"
echo "--------------------------------------------------"
for i in {1..3}; do
    echo "Creating product $i..."
    curl -s -X POST ${API_URL}/products \
      -H "Content-Type: application/json" \
      -H "X-Correlation-ID: bulk-${i}-$(date +%s)" \
      -d "{
        \"name\": \"Product $i\",
        \"price\": $(echo "scale=2; $i * 10.99" | bc),
        \"description\": \"Bulk product number $i\"
      }" | jq -c '{id: .id, name: .name, price: .price}'
    sleep 0.5
done

echo ""
echo "=================================================="
echo "  Examples completed!"
echo "=================================================="
