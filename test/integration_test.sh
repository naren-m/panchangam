#!/bin/bash

set -e

echo "Running Gateway Integration Tests..."
echo "====================================="

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

HTTP_PORT=${HTTP_PORT:-8080}
HTTP_BASE_URL="http://localhost:${HTTP_PORT}"

cleanup() {
    echo ""
    echo "Cleaning up..."
    ./scripts/stop-servers.sh || true
}

trap cleanup EXIT

# Start servers
echo -e "${YELLOW}Starting servers...${NC}"
./scripts/start-servers.sh

# Wait for servers to start
echo "Waiting for servers to be ready..."
sleep 5

# Check if servers are running
if ! curl -s "${HTTP_BASE_URL}/api/v1/health" > /dev/null; then
    echo -e "${RED}FAIL: gateway server failed to start${NC}"
    exit 1
fi

echo -e "${GREEN}PASS: servers started successfully${NC}"
echo ""

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Function to run a test
run_test() {
    local test_name=$1
    local expected_status=$2
    local check_function=$3
    shift 3
    local response
    local status_code

    echo -n "Testing: $test_name... "

    response=$("$@" 2>&1) || true
    status_code=$(printf '%s\n' "$response" | awk '/< HTTP\/1\.1/ { print $3; exit }')
    status_code=${status_code:-0}

    if [ "$status_code" = "$expected_status" ] && "$check_function" "$response"; then
        echo -e "${GREEN}PASS${NC}"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${RED}FAIL${NC}"
        echo "Expected status: $expected_status, Got: $status_code"
        echo "Response: $response"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
}

# Check functions
check_health() {
    echo "$1" | grep -q '"status":"healthy"'
}

check_panchangam_data() {
    echo "$1" | grep -q '"date":"2024-01-15"' && \
    echo "$1" | grep -q '"tithi":'
}

check_error_format() {
    echo "$1" | grep -q '"error":' && \
    echo "$1" | grep -q '"code":' && \
    echo "$1" | grep -q '"message":'
}

check_cors_headers() {
    echo "$1" | grep -q "Access-Control-Allow-Origin: http://localhost:5173"
}

echo "Running integration tests..."
echo "============================"

# Test 1: Health check
run_test "Health check endpoint" \
    "200" \
    check_health \
    curl -s -v "${HTTP_BASE_URL}/api/v1/health"

# Test 2: Valid panchangam request
run_test "Valid panchangam request" \
    "200" \
    check_panchangam_data \
    curl -s -v "${HTTP_BASE_URL}/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946&tz=Asia/Kolkata"

# Test 3: Missing parameter error
run_test "Missing date parameter" \
    "400" \
    check_error_format \
    curl -s -v "${HTTP_BASE_URL}/api/v1/panchangam?lat=12.9716&lng=77.5946"

# Test 4: Invalid parameter error
run_test "Invalid latitude format" \
    "400" \
    check_error_format \
    curl -s -v "${HTTP_BASE_URL}/api/v1/panchangam?date=2024-01-15&lat=invalid&lng=77.5946"

# Test 5: CORS headers
run_test "CORS headers for allowed origin" \
    "200" \
    check_cors_headers \
    curl -s -v -H "Origin: http://localhost:5173" "${HTTP_BASE_URL}/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946"

# Test 6: Request ID tracking
echo -n "Testing: Request ID tracking... "
response=$(curl -s -v -H "X-Request-Id: test-integration-123" "${HTTP_BASE_URL}/api/v1/health" 2>&1) || true
if echo "$response" | grep -q "X-Request-Id: test-integration-123"; then
    echo -e "${GREEN}PASS${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${RED}FAIL${NC}"
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi

# Test 7: Performance test
echo -n "Testing: Response time (<100ms)... "
start_time=$(date +%s%N)
curl -s "${HTTP_BASE_URL}/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946" > /dev/null
end_time=$(date +%s%N)
duration=$(( (end_time - start_time) / 1000000 ))

if [ "$duration" -lt 100 ]; then
    echo -e "${GREEN}PASS (${duration}ms)${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${YELLOW}WARN (${duration}ms)${NC}"
fi

# Cleanup
cleanup
trap - EXIT

# Summary
echo ""
echo "====================================="
echo -e "${GREEN}Tests Passed: $TESTS_PASSED${NC}"
if [ "$TESTS_FAILED" -gt 0 ]; then
    echo -e "${RED}Tests Failed: $TESTS_FAILED${NC}"
else
    echo -e "${GREEN}Tests Failed: $TESTS_FAILED${NC}"
fi
echo "====================================="

if [ "$TESTS_FAILED" -eq 0 ]; then
    echo -e "${GREEN}PASS: all integration tests passed${NC}"
    exit 0
else
    echo -e "${RED}FAIL: some tests failed${NC}"
    exit 1
fi
