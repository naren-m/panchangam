#!/bin/bash

set -e

echo "🧪 Running Gateway Integration Tests..."
echo "====================================="

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Start servers
echo -e "${YELLOW}Starting servers...${NC}"
./scripts/start-servers.sh &
SERVER_PID=$!

# Wait for servers to start
echo "Waiting for servers to be ready..."
sleep 5

# Check if servers are running
if ! curl -s http://localhost:8080/api/v1/health > /dev/null; then
    echo -e "${RED}❌ Gateway server failed to start${NC}"
    kill $SERVER_PID 2>/dev/null || true
    exit 1
fi

echo -e "${GREEN}✅ Servers started successfully${NC}"
echo ""

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Function to run a test
run_test() {
    local test_name=$1
    local command=$2
    local expected_status=$3
    local check_function=$4
    
    echo -n "Testing: $test_name... "
    
    response=$(eval "$command" 2>&1) || true
    status_code=$(echo "$response" | grep -E "< HTTP/1.1" | awk '{print $3}' || echo "0")
    
    if [ "$status_code" = "$expected_status" ] && eval "$check_function \"$response\""; then
        echo -e "${GREEN}✅ PASSED${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}❌ FAILED${NC}"
        echo "Expected status: $expected_status, Got: $status_code"
        echo "Response: $response"
        ((TESTS_FAILED++))
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
    "curl -s -v http://localhost:8080/api/v1/health 2>&1" \
    "200" \
    check_health

# Test 2: Valid panchangam request
run_test "Valid panchangam request" \
    "curl -s -v 'http://localhost:8080/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946&tz=Asia/Kolkata' 2>&1" \
    "200" \
    check_panchangam_data

# Test 3: Missing parameter error
run_test "Missing date parameter" \
    "curl -s -v 'http://localhost:8080/api/v1/panchangam?lat=12.9716&lng=77.5946' 2>&1" \
    "400" \
    check_error_format

# Test 4: Invalid parameter error
run_test "Invalid latitude format" \
    "curl -s -v 'http://localhost:8080/api/v1/panchangam?date=2024-01-15&lat=invalid&lng=77.5946' 2>&1" \
    "400" \
    check_error_format

# Test 5: CORS headers
run_test "CORS headers for allowed origin" \
    "curl -s -v -H 'Origin: http://localhost:5173' 'http://localhost:8080/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946' 2>&1" \
    "200" \
    check_cors_headers

# Test 6: Request ID tracking
echo -n "Testing: Request ID tracking... "
response=$(curl -s -v -H "X-Request-Id: test-integration-123" http://localhost:8080/api/v1/health 2>&1)
if echo "$response" | grep -q "X-Request-Id: test-integration-123"; then
    echo -e "${GREEN}✅ PASSED${NC}"
    ((TESTS_PASSED++))
else
    echo -e "${RED}❌ FAILED${NC}"
    ((TESTS_FAILED++))
fi

# Test 7: Performance test
echo -n "Testing: Response time (<100ms)... "
start_time=$(date +%s%N)
curl -s "http://localhost:8080/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946" > /dev/null
end_time=$(date +%s%N)
duration=$(( (end_time - start_time) / 1000000 ))

if [ $duration -lt 100 ]; then
    echo -e "${GREEN}✅ PASSED (${duration}ms)${NC}"
    ((TESTS_PASSED++))
else
    echo -e "${YELLOW}⚠️  WARNING (${duration}ms)${NC}"
fi

# Cleanup
echo ""
echo "Cleaning up..."
pkill -f grpc-server || true
pkill -f gateway-server || true
kill $SERVER_PID 2>/dev/null || true

# Summary
echo ""
echo "====================================="
echo -e "${GREEN}Tests Passed: $TESTS_PASSED${NC}"
if [ $TESTS_FAILED -gt 0 ]; then
    echo -e "${RED}Tests Failed: $TESTS_FAILED${NC}"
else
    echo -e "${GREEN}Tests Failed: $TESTS_FAILED${NC}"
fi
echo "====================================="

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ All integration tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed${NC}"
    exit 1
fi