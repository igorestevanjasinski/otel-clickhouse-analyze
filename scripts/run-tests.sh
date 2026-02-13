#!/bin/bash
set -e

echo "=================================================="
echo "  Test Suite Execution"
echo "=================================================="

cd "$(dirname "$0")/.."

FAILED=0
TOTAL_TESTS=0
PASSED_TESTS=0

echo ""
echo "[1/2] Running Python tests..."
echo "--------------------------------------------------"
cd services/product-api

if [ -d "venv" ]; then
    source venv/bin/activate
    
    pytest tests/ -v --cov=app --cov-report=term --cov-report=html:htmlcov
    PYTHON_EXIT=$?
    
    if [ $PYTHON_EXIT -eq 0 ]; then
        echo "✓ Python tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo "✗ Python tests failed"
        FAILED=$((FAILED + 1))
    fi
    
    deactivate
else
    echo "✗ Python venv not found. Run ./scripts/setup.sh first"
    FAILED=$((FAILED + 1))
fi

TOTAL_TESTS=$((TOTAL_TESTS + 1))
cd ../..

echo ""
echo "[2/2] Running Go tests..."
echo "--------------------------------------------------"
cd services/product-consumer

go test ./internal/consumer/... -v -cover
GO_CONSUMER_EXIT=$?

if [ $GO_CONSUMER_EXIT -eq 0 ]; then
    echo "✓ Go consumer tests passed"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    echo "✗ Go consumer tests failed"
    FAILED=$((FAILED + 1))
fi

TOTAL_TESTS=$((TOTAL_TESTS + 1))

echo ""
echo "Running Go integration tests (with testcontainers)..."
go test ./internal/repository/... -v -cover
GO_REPO_EXIT=$?

if [ $GO_REPO_EXIT -eq 0 ]; then
    echo "✓ Go repository tests passed"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    echo "✗ Go repository tests failed"
    FAILED=$((FAILED + 1))
fi

TOTAL_TESTS=$((TOTAL_TESTS + 1))
cd ../..

echo ""
echo "=================================================="
echo "  Test Results Summary"
echo "=================================================="
echo "Total test suites: $TOTAL_TESTS"
echo "Passed: $PASSED_TESTS"
echo "Failed: $FAILED"
echo ""

if [ $FAILED -eq 0 ]; then
    echo "✓ All tests passed!"
    echo ""
    echo "Coverage reports:"
    echo "  - Python: services/product-api/htmlcov/index.html"
    echo "  - Go: Run 'go test -coverprofile=coverage.out' for detailed coverage"
    exit 0
else
    echo "✗ Some tests failed"
    exit 1
fi
