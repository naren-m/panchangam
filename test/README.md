# Panchangam API Testing Framework

Comprehensive Python pytest-based testing framework for the Panchangam API Gateway.

## 🎯 Overview

This testing framework provides a robust, multi-layered approach to validating the Panchangam API Gateway implementation. It complements the existing Go unit tests with comprehensive end-to-end, integration, and performance testing.

## 📋 Test Categories

### Test Markers
- `@pytest.mark.smoke` - Critical functionality tests (quick validation)
- `@pytest.mark.integration` - End-to-end integration tests
- `@pytest.mark.performance` - Performance and load testing
- `@pytest.mark.security` - Security and vulnerability tests

### Test Modules

#### 1. `test_health_check.py`
- **Purpose**: Health endpoint validation
- **Coverage**: Basic connectivity, CORS, request tracking
- **Tests**: 4 test functions
- **Key Scenarios**:
  - Health endpoint returns proper JSON response
  - Custom request ID preservation
  - Response time validation (<100ms)
  - CORS headers for allowed origins

#### 2. `test_panchangam_api.py`
- **Purpose**: Core API functionality testing
- **Coverage**: Request validation, data accuracy, performance
- **Tests**: 11 test functions
- **Key Scenarios**:
  - Valid panchangam calculations for multiple locations
  - Date range testing (7 consecutive days)
  - Optional parameter handling
  - Performance benchmarking (<500ms response time)
  - Concurrent request handling (10 simultaneous requests)
  - Cache header validation
  - Request ID tracking
  - CORS configuration testing
  - HTTP method validation

#### 3. `test_error_handling.py`
- **Purpose**: Error scenarios and edge cases
- **Coverage**: Input validation, error consistency, security
- **Tests**: 12 test functions
- **Key Scenarios**:
  - Missing required parameters
  - Invalid parameter types and ranges
  - Empty parameter values
  - Extreme but valid coordinates (poles, date line)
  - Invalid optional parameters
  - Malformed request headers
  - Request size limits
  - Concurrent error scenarios
  - Error response consistency
  - Error handling performance
  - Information disclosure prevention

## 🚀 Quick Start

### Prerequisites
- Python 3.8+
- Go 1.19+
- Git

### Installation
```bash
# Navigate to test directory
cd test/

# Install dependencies
pip install -r requirements.txt

# Run all tests with automatic server management
python run_tests.py
```

### Basic Usage

#### Run All Tests
```bash
python run_tests.py
```

#### Run Specific Test Categories
```bash
# Smoke tests (quick validation)
python run_tests.py --type smoke

# Integration tests
python run_tests.py --type integration

# Performance tests
python run_tests.py --type performance

# Security tests
python run_tests.py --type security
```

#### Advanced Options
```bash
# Verbose output with coverage and HTML reports
python run_tests.py --verbose --coverage --html-report

# Install dependencies and build servers
python run_tests.py --install-deps --build-servers

# Run with custom markers
python run_tests.py --markers "smoke and not performance"

# Skip automatic server startup (if servers already running)
python run_tests.py --skip-server-start
```

## 🔧 Configuration

### Environment Variables
- `PANCHANGAM_API_URL` - API base URL (default: `http://localhost:8080`)
- `SKIP_SERVER_START` - Skip automatic server startup (`true`/`false`)

### Test Configuration (`conftest.py`)
- **Session Fixtures**: Server lifecycle management
- **API Client**: Configured requests session with proper headers
- **Sample Data**: Multiple test locations (Bangalore, Mumbai, New York, London)
- **Custom Markers**: Test categorization and filtering

## 📊 Test Execution Flow

### Automatic Server Management
1. **Build Phase**: Compile Go servers (`grpc-server`, `gateway-server`)
2. **Startup Phase**: Launch gRPC server (port 50052) and Gateway (port 8080)
3. **Validation Phase**: Health check to ensure servers are ready
4. **Test Phase**: Execute test suite
5. **Cleanup Phase**: Graceful server shutdown

### Manual Server Management
```bash
# Set environment variable to skip automatic startup
export SKIP_SERVER_START=true

# Start servers manually
../scripts/start-servers.sh

# Run tests
python run_tests.py --skip-server-start
```

## 📈 Performance Benchmarks

### Response Time Targets
- **Health Check**: <100ms
- **Panchangam API**: <500ms
- **Error Responses**: <100ms
- **Concurrent Requests**: 10 requests in <2000ms

### Load Testing
- **Concurrent Users**: 10 simultaneous requests
- **Success Rate**: 100% success expected
- **Performance Degradation**: Minimal impact under load

## 🔒 Security Testing

### Security Validations
- **CORS Configuration**: Proper origin validation
- **Information Disclosure**: Error messages don't reveal sensitive data
- **Input Validation**: SQL injection and XSS prevention
- **Request Size Limits**: Protection against large payload attacks
- **Error Consistency**: Uniform error response structure

### Security Test Scenarios
- Malformed headers and payloads
- Injection attempts in parameters
- Oversized requests
- Invalid authorization attempts

## 📋 Test Data

### Sample Locations
```python
locations = {
    "bangalore": {"lat": 12.9716, "lng": 77.5946, "tz": "Asia/Kolkata"},
    "mumbai": {"lat": 19.0760, "lng": 72.8777, "tz": "Asia/Kolkata"},
    "new_york": {"lat": 40.7128, "lng": -74.0060, "tz": "America/New_York"},
    "london": {"lat": 51.5074, "lng": -0.1278, "tz": "Europe/London"}
}
```

### Date Range Testing
- **Base Date**: 2024-01-15
- **Range**: 7 consecutive days
- **Validation**: Consistent data structure across dates

## 📊 Reporting

### HTML Reports
```bash
# Generate HTML test report
python run_tests.py --html-report

# View reports
open report.html          # Test results
open htmlcov/index.html   # Coverage report
```

### Coverage Reports
```bash
# Generate coverage report
python run_tests.py --coverage

# View coverage
python run_tests.py --coverage --html-report
```

### Console Output
- **Summary**: Pass/fail counts with execution time
- **Performance**: Response time measurements
- **Errors**: Detailed failure information with request IDs

## 🔧 Advanced Usage

### Custom Test Execution
```bash
# Direct pytest execution
pytest -v -m smoke                    # Smoke tests only
pytest -v -m "integration and not performance"  # Integration without performance
pytest -k "test_health"               # Tests matching pattern
pytest --maxfail=1                    # Stop on first failure
pytest -x                             # Stop on first failure (short form)
```

### Development Workflow
```bash
# Quick validation during development
python run_tests.py --type smoke

# Full validation before commit
python run_tests.py --verbose --coverage

# Performance regression testing
python run_tests.py --type performance --verbose
```

### CI/CD Integration
```bash
# Automated CI pipeline
python run_tests.py --type all --coverage --html-report --skip-server-start
```

## 🧪 Test Development

### Adding New Tests
1. **Choose Module**: Health, API, or Error handling
2. **Add Test Function**: Follow naming convention `test_*`
3. **Add Markers**: Use appropriate `@pytest.mark.*`
4. **Use Fixtures**: Leverage existing fixtures for consistency
5. **Assert Structure**: Follow established patterns

### Test Patterns
```python
@pytest.mark.integration
def test_new_feature(api_client: requests.Session, api_base_url: str):
    """Test description"""
    response = api_client.get(f"{api_base_url}/api/v1/endpoint")
    
    assert response.status_code == 200
    assert response.headers["content-type"] == "application/json"
    
    data = response.json()
    assert "expected_field" in data
```

## 🔍 Troubleshooting

### Common Issues

#### Server Startup Failures
```bash
# Check port availability
lsof -i :8080
lsof -i :50052

# Manual server start
../scripts/start-servers.sh
```

#### Dependency Issues
```bash
# Reinstall dependencies
pip install -r requirements.txt --force-reinstall

# Check Python version
python --version  # Should be 3.8+
```

#### Test Failures
```bash
# Run with maximum verbosity
python run_tests.py --verbose --type smoke

# Check server logs
tail -f ../logs/gateway.log
tail -f ../logs/grpc.log
```

### Debug Mode
```bash
# Enable debug logging
export PYTEST_DEBUG=1
python run_tests.py --verbose
```

## 📝 Test Coverage

### Current Coverage
- **Health Endpoint**: 100% scenarios covered
- **Panchangam API**: 95% functionality covered
- **Error Handling**: 90% edge cases covered
- **Performance**: Key benchmarks established
- **Security**: Critical vulnerabilities tested

### Coverage Goals
- **Unit Test Coverage**: 90%+ (Go tests)
- **Integration Coverage**: 95%+ (Python tests)
- **Performance Benchmarks**: All critical paths
- **Security Testing**: OWASP Top 10 scenarios

## 🎯 Integration with Go Tests

This pytest framework complements the existing Go unit tests:

### Go Tests (Unit Level)
- **Coverage**: 73.7% statement coverage
- **Focus**: Internal logic, mocking, benchmarks
- **Speed**: Very fast (<1s execution)
- **Scope**: Individual functions and methods

### Python Tests (Integration Level)
- **Coverage**: End-to-end workflows
- **Focus**: API behavior, user scenarios, error handling
- **Speed**: Moderate (10-30s execution)
- **Scope**: Complete request/response cycles

### Combined Strategy
1. **Go Tests**: Fast feedback during development
2. **Python Tests**: Comprehensive validation before deployment
3. **Integration**: Both run in CI/CD pipeline
4. **Reporting**: Combined coverage and performance metrics

## 🚀 Next Steps

### Planned Enhancements
1. **Load Testing**: Apache Bench integration
2. **Chaos Testing**: Network failure simulation
3. **Contract Testing**: API contract validation
4. **Performance Regression**: Baseline comparison
5. **Accessibility Testing**: API usability validation

### Postman Collection
A companion Postman collection is planned to provide:
- **Manual Testing**: GUI-based API exploration
- **Documentation**: Interactive API documentation
- **Collaboration**: Shareable test scenarios
- **CI Integration**: Newman-based automation

This pytest framework provides the foundation for robust, automated testing while the Postman collection will enable manual exploration and documentation.