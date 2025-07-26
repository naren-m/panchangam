#!/bin/bash

# Deployment validation script for Panchangam application
set -euo pipefail

# Default values
ENVIRONMENT="staging"
BASE_URL=""
TIMEOUT=300
VERBOSE=false

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Usage function
usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Validate Panchangam application deployment.

OPTIONS:
    -e, --environment ENV    Environment to validate (staging/production) [default: staging]
    -u, --url URL           Base URL to test [auto-detected if not provided]
    -t, --timeout SECONDS   Timeout for tests in seconds [default: 300]
    -v, --verbose           Verbose output
    -h, --help              Show this help message

EXAMPLES:
    $0 -e staging
    $0 -e production -u https://api.panchangam.example.com
    $0 --verbose --timeout 600

EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -e|--environment)
            ENVIRONMENT="$2"
            shift 2
            ;;
        -u|--url)
            BASE_URL="$2"
            shift 2
            ;;
        -t|--timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            usage
            exit 1
            ;;
    esac
done

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_verbose() {
    if [ "$VERBOSE" = true ]; then
        echo -e "${BLUE}[VERBOSE]${NC} $1"
    fi
}

# Auto-detect base URL if not provided
detect_base_url() {
    if [ -n "$BASE_URL" ]; then
        return
    fi
    
    if [ "$ENVIRONMENT" = "production" ]; then
        BASE_URL="https://api.panchangam.example.com"
    else
        # Try to get from kubectl
        if command -v kubectl &> /dev/null; then
            local service_ip
            service_ip=$(kubectl get svc panchangam-backend -n panchangam -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "")
            if [ -n "$service_ip" ]; then
                BASE_URL="http://$service_ip"
            else
                BASE_URL="http://localhost:8080"
            fi
        else
            BASE_URL="http://localhost:8080"
        fi
    fi
    
    log_info "Auto-detected base URL: $BASE_URL"
}

# Health check validation
validate_health_check() {
    log_info "Validating health check endpoint..."
    
    local health_url="$BASE_URL/api/v1/health"
    local start_time=$(date +%s)
    local max_time=$((start_time + TIMEOUT))
    
    while [ $(date +%s) -lt $max_time ]; do
        log_verbose "Testing health endpoint: $health_url"
        
        local response
        response=$(curl -s -w "\n%{http_code}" "$health_url" 2>/dev/null || echo -e "\n000")
        local body=$(echo "$response" | head -n -1)
        local status_code=$(echo "$response" | tail -n 1)
        
        log_verbose "Response status: $status_code"
        log_verbose "Response body: $body"
        
        if [ "$status_code" = "200" ]; then
            if echo "$body" | grep -q "healthy"; then
                log_success "Health check passed"
                return 0
            else
                log_warning "Health endpoint returned 200 but body doesn't contain 'healthy'"
            fi
        fi
        
        log_verbose "Health check failed, retrying in 5 seconds..."
        sleep 5
    done
    
    log_error "Health check failed after $TIMEOUT seconds"
    return 1
}

# API functionality validation
validate_api_functionality() {
    log_info "Validating API functionality..."
    
    local test_cases=(
        # Test case format: "description|endpoint|expected_field"
        "Basic panchangam calculation|/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946|tithi"
        "Different location|/api/v1/panchangam?date=2024-06-21&lat=40.7128&lng=-74.0060&tz=America/New_York|nakshatra" 
        "UK location|/api/v1/panchangam?date=2024-12-21&lat=51.5074&lng=-0.1278&tz=Europe/London|yoga"
    )
    
    local passed=0
    local total=${#test_cases[@]}
    
    for test_case in "${test_cases[@]}"; do
        IFS='|' read -r description endpoint expected_field <<< "$test_case"
        
        log_verbose "Testing: $description"
        log_verbose "Endpoint: $BASE_URL$endpoint"
        
        local response
        response=$(curl -s -w "\n%{http_code}" "$BASE_URL$endpoint" 2>/dev/null || echo -e "\n000")
        local body=$(echo "$response" | head -n -1)
        local status_code=$(echo "$response" | tail -n 1)
        
        log_verbose "Response status: $status_code"
        
        if [ "$status_code" = "200" ]; then
            if echo "$body" | grep -q "\"$expected_field\""; then
                log_success "✅ $description"
                ((passed++))
            else
                log_error "❌ $description - Expected field '$expected_field' not found"
                log_verbose "Response body: $body"
            fi
        else
            log_error "❌ $description - HTTP $status_code"
            log_verbose "Response body: $body"
        fi
    done
    
    log_info "API functionality tests: $passed/$total passed"
    
    if [ $passed -eq $total ]; then
        return 0
    else
        return 1
    fi
}

# Performance validation
validate_performance() {
    log_info "Validating API performance..."
    
    local endpoint="$BASE_URL/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946"
    local iterations=10
    local total_time=0
    local successful_requests=0
    
    for i in $(seq 1 $iterations); do
        log_verbose "Performance test iteration $i/$iterations"
        
        local start_time=$(date +%s.%3N)
        local response
        response=$(curl -s -w "\n%{http_code}" "$endpoint" 2>/dev/null || echo -e "\n000")
        local end_time=$(date +%s.%3N)
        local status_code=$(echo "$response" | tail -n 1)
        
        if [ "$status_code" = "200" ]; then
            local request_time=$(echo "$end_time - $start_time" | bc)
            total_time=$(echo "$total_time + $request_time" | bc)
            ((successful_requests++))
            log_verbose "Request $i: ${request_time}s"
        else
            log_verbose "Request $i failed with status $status_code"
        fi
    done
    
    if [ $successful_requests -gt 0 ]; then
        local avg_time=$(echo "scale=3; $total_time / $successful_requests" | bc)
        log_info "Performance results:"
        log_info "  Successful requests: $successful_requests/$iterations"
        log_info "  Average response time: ${avg_time}s"
        
        # Check if average time is acceptable (< 5 seconds)
        if (( $(echo "$avg_time < 5.0" | bc -l) )); then
            log_success "Performance validation passed"
            return 0
        else
            log_warning "Performance validation failed - average time ${avg_time}s > 5.0s"
            return 1
        fi
    else
        log_error "Performance validation failed - no successful requests"
        return 1
    fi
}

# Error handling validation
validate_error_handling() {
    log_info "Validating error handling..."
    
    local error_test_cases=(
        # Test case format: "description|endpoint|expected_status"
        "Missing date parameter|/api/v1/panchangam?lat=12.9716&lng=77.5946|400"
        "Invalid latitude|/api/v1/panchangam?date=2024-01-15&lat=999&lng=77.5946|400"
        "Invalid date format|/api/v1/panchangam?date=invalid&lat=12.9716&lng=77.5946|400"
    )
    
    local passed=0
    local total=${#error_test_cases[@]}
    
    for test_case in "${error_test_cases[@]}"; do
        IFS='|' read -r description endpoint expected_status <<< "$test_case"
        
        log_verbose "Testing error case: $description"
        
        local response
        response=$(curl -s -w "\n%{http_code}" "$BASE_URL$endpoint" 2>/dev/null || echo -e "\n000")
        local status_code=$(echo "$response" | tail -n 1)
        
        if [ "$status_code" = "$expected_status" ]; then
            log_success "✅ $description"
            ((passed++))
        else
            log_error "❌ $description - Expected $expected_status, got $status_code"
        fi
    done
    
    log_info "Error handling tests: $passed/$total passed"
    
    if [ $passed -eq $total ]; then
        return 0
    else
        return 1
    fi
}

# Frontend validation (if accessible)
validate_frontend() {
    log_info "Validating frontend accessibility..."
    
    local frontend_url
    if [ "$ENVIRONMENT" = "production" ]; then
        frontend_url="https://panchangam.example.com"
    else
        # Try to detect frontend URL
        if command -v kubectl &> /dev/null; then
            local service_ip
            service_ip=$(kubectl get svc panchangam-frontend -n panchangam -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "")
            if [ -n "$service_ip" ]; then
                frontend_url="http://$service_ip"
            else
                frontend_url="http://localhost:80"
            fi
        else
            frontend_url="http://localhost:80"
        fi
    fi
    
    log_verbose "Testing frontend URL: $frontend_url"
    
    local response
    response=$(curl -s -w "\n%{http_code}" "$frontend_url" 2>/dev/null || echo -e "\n000")
    local status_code=$(echo "$response" | tail -n 1)
    
    if [ "$status_code" = "200" ]; then
        local body=$(echo "$response" | head -n -1)
        if echo "$body" | grep -q -i "panchangam"; then
            log_success "Frontend validation passed"
            return 0
        else
            log_warning "Frontend returned 200 but doesn't contain 'panchangam'"
            return 1
        fi
    else
        log_warning "Frontend validation failed - HTTP $status_code"
        return 1
    fi
}

# Generate validation report
generate_report() {
    local health_status=$1
    local api_status=$2
    local performance_status=$3
    local error_status=$4
    local frontend_status=$5
    
    echo
    log_info "=== DEPLOYMENT VALIDATION REPORT ==="
    log_info "Environment: $ENVIRONMENT"
    log_info "Base URL: $BASE_URL"
    log_info "Timestamp: $(date)"
    echo
    
    # Status indicators
    local health_indicator=$( [ $health_status -eq 0 ] && echo "✅ PASS" || echo "❌ FAIL" )
    local api_indicator=$( [ $api_status -eq 0 ] && echo "✅ PASS" || echo "❌ FAIL" )
    local performance_indicator=$( [ $performance_status -eq 0 ] && echo "✅ PASS" || echo "❌ FAIL" )
    local error_indicator=$( [ $error_status -eq 0 ] && echo "✅ PASS" || echo "❌ FAIL" )
    local frontend_indicator=$( [ $frontend_status -eq 0 ] && echo "✅ PASS" || echo "⚠️  WARN" )
    
    echo "| Test Category        | Status       |"
    echo "|---------------------|--------------|"
    echo "| Health Check        | $health_indicator   |"
    echo "| API Functionality   | $api_indicator   |"
    echo "| Performance         | $performance_indicator   |"
    echo "| Error Handling      | $error_indicator   |"
    echo "| Frontend Access     | $frontend_indicator   |"
    echo
    
    # Overall status
    if [ $health_status -eq 0 ] && [ $api_status -eq 0 ] && [ $performance_status -eq 0 ] && [ $error_status -eq 0 ]; then
        log_success "🎉 DEPLOYMENT VALIDATION SUCCESSFUL"
        return 0
    else
        log_error "💥 DEPLOYMENT VALIDATION FAILED"
        return 1
    fi
}

# Main validation flow
main() {
    log_info "Starting deployment validation for $ENVIRONMENT environment"
    
    detect_base_url
    
    # Run all validation tests
    local health_status=1
    local api_status=1
    local performance_status=1
    local error_status=1
    local frontend_status=1
    
    validate_health_check && health_status=0
    validate_api_functionality && api_status=0
    validate_performance && performance_status=0
    validate_error_handling && error_status=0
    validate_frontend && frontend_status=0
    
    # Generate final report
    generate_report $health_status $api_status $performance_status $error_status $frontend_status
}

# Check if bc is available for calculations
if ! command -v bc &> /dev/null; then
    log_error "bc (calculator) is required but not installed"
    exit 1
fi

# Run main function
main "$@"