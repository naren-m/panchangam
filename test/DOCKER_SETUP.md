# Docker-Based Testing Setup

Complete Docker containerization for the Panchangam API testing framework, providing isolated, reproducible testing environment.

## 🚀 Quick Start

### Prerequisites
- Docker Engine 20.10+
- Docker Compose 2.0+
- Make (optional, for convenience commands)

### Instant Setup
```bash
# Navigate to test directory
cd test/

# Run all tests (builds automatically)
make test

# Or using the shell script directly
./run-docker-tests.sh
```

## 🐳 Docker Configuration

### Container Architecture
```
panchangam-test-runner/
├── Python 3.11 runtime
├── Go development tools
├── Testing dependencies
├── Project source code (mounted)
└── Test output directories
```

### Key Features
- **Isolated Environment**: No dependency conflicts with host system
- **Reproducible Results**: Consistent testing across different machines
- **Automatic Server Management**: Built-in gRPC and Gateway server lifecycle
- **Volume Mounting**: Real-time code changes reflected in container
- **Report Persistence**: Test reports saved to host filesystem

## 📁 File Structure

```
test/
├── Dockerfile              # Main container definition
├── docker-compose.yml      # Service orchestration
├── .dockerignore           # Docker build context exclusions
├── run-docker-tests.sh     # Enhanced test runner script
├── Makefile                # Convenient command shortcuts
├── requirements.txt        # Core Python dependencies
├── requirements-dev.txt    # Development dependencies
└── DOCKER_SETUP.md        # This documentation
```

## 🔧 Usage Methods

### Method 1: Make Commands (Recommended)
```bash
# Quick tests
make test-smoke              # Fast validation
make test-integration        # Full API testing
make test-performance        # Performance benchmarks
make test-security           # Security validation

# Development
make shell                   # Interactive container shell
make test-coverage           # Tests with coverage
make test-html              # Tests with HTML reports

# Maintenance
make clean                   # Remove containers
make rebuild                # Force image rebuild
```

### Method 2: Shell Script
```bash
# Basic usage
./run-docker-tests.sh                    # All tests
./run-docker-tests.sh smoke              # Smoke tests
./run-docker-tests.sh integration -v     # Verbose integration tests

# Advanced options
./run-docker-tests.sh all --coverage --html-report --build --cleanup
./run-docker-tests.sh --shell            # Interactive mode
./run-docker-tests.sh --logs             # View logs
```

### Method 3: Direct Docker Compose
```bash
# Build and start
docker-compose build
docker-compose up -d panchangam-tests

# Execute tests
docker-compose exec panchangam-tests python run_tests.py --type smoke

# Interactive shell
docker-compose exec panchangam-tests /bin/bash

# Cleanup
docker-compose down
```

## 🧪 Test Execution Examples

### Quick Validation
```bash
# 30-second smoke test
make quick-test

# Health check
make quick-check
```

### Comprehensive Testing
```bash
# Full test suite with reports
make test-full

# CI/CD pipeline simulation
make ci-full
```

### Development Workflow
```bash
# Setup development environment
make dev-setup

# Development test cycle
make dev-test

# Code quality checks
make lint format type-check security-scan
```

### Performance Testing
```bash
# Performance benchmarks
make test-performance

# Load testing (if available)
docker-compose run --rm panchangam-tests locust --headless -u 10 -r 2 -t 30s
```

## 📊 Reports and Output

### Generated Reports
- **Test Report**: `./reports/report.html`
- **Coverage Report**: `./htmlcov/index.html`
- **Logs**: `./logs/`

### Viewing Reports
```bash
# Generate and open reports
make report

# Coverage only
make coverage-report

# Clean old reports
make clean-reports
```

## 🔧 Configuration

### Environment Variables
```bash
# API configuration
PANCHANGAM_API_URL=http://localhost:8080

# Test configuration
SKIP_SERVER_START=false
PYTHONPATH=/workspace
PYTEST_DISABLE_PLUGIN_AUTOLOAD=1

# Development settings
PYTEST_DEBUG=1
```

### Docker Compose Override
Create `docker-compose.override.yml` for local customizations:
```yaml
version: '3.8'
services:
  panchangam-tests:
    environment:
      - DEBUG=true
      - LOG_LEVEL=debug
    ports:
      - "9090:8080"  # Alternative port mapping
```

## 🚀 Advanced Features

### Multi-Stage Testing
```bash
# Stage 1: Smoke tests
./run-docker-tests.sh smoke

# Stage 2: Integration tests (if smoke passes)
./run-docker-tests.sh integration

# Stage 3: Performance validation
./run-docker-tests.sh performance
```

### Parallel Test Execution
```bash
# Parallel test execution (if supported)
docker-compose run --rm panchangam-tests pytest -n auto

# Custom worker count
docker-compose run --rm panchangam-tests pytest -n 4
```

### Custom Test Selection
```bash
# Run specific test files
docker-compose run --rm panchangam-tests pytest test_health_check.py -v

# Run tests matching pattern
docker-compose run --rm panchangam-tests pytest -k "test_error" -v

# Run tests with custom markers
docker-compose run --rm panchangam-tests pytest -m "smoke and not slow" -v
```

## 🔍 Troubleshooting

### Common Issues

#### Port Conflicts
```bash
# Check port usage
lsof -i :8080
lsof -i :50052

# Use alternative ports
PANCHANGAM_API_URL=http://localhost:8081 ./run-docker-tests.sh
```

#### Build Failures
```bash
# Clean rebuild
make clean-all
make rebuild

# Check Docker system
docker system df
docker system prune
```

#### Container Issues
```bash
# Debug container
docker-compose run --rm panchangam-tests /bin/bash

# Check logs
make logs

# Container status
make status
```

#### Test Failures
```bash
# Verbose debugging
./run-docker-tests.sh all -v --logs

# Interactive debugging
make shell
# Inside container:
python -m pytest test_health_check.py::test_health_check_endpoint -v -s
```

### Debug Commands
```bash
# System information
make info

# Container health check
docker-compose exec panchangam-tests python --version
docker-compose exec panchangam-tests go version

# Network connectivity
docker-compose exec panchangam-tests curl http://localhost:8080/api/v1/health
```

## 🔐 Security Considerations

### Container Security
- **Non-root User**: Container runs with limited privileges
- **Read-only Filesystem**: Core system files are read-only
- **Network Isolation**: Containers run in isolated network
- **Minimal Attack Surface**: Only necessary packages installed

### Secret Management
```bash
# Environment file for secrets (not committed)
echo "API_KEY=your-secret-key" > .env

# Docker Compose will automatically load .env
docker-compose up
```

## 📈 Performance Optimization

### Build Optimization
- **Multi-stage Build**: Optimized Docker layer caching
- **Minimal Base Image**: Python slim image reduces size
- **Dependency Caching**: Requirements cached separately

### Runtime Optimization
- **Volume Mounting**: Avoid copying large codebases
- **Parallel Execution**: Multiple test workers when supported
- **Resource Limits**: Configurable memory and CPU limits

### Resource Configuration
```yaml
# docker-compose.override.yml
services:
  panchangam-tests:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 1G
        reservations:
          cpus: '1.0'
          memory: 512M
```

## 🚀 CI/CD Integration

### GitHub Actions Example
```yaml
- name: Run Docker Tests
  run: |
    cd test/
    make ci-test
```

### GitLab CI Example
```yaml
test:
  script:
    - cd test/
    - make ci-full
  artifacts:
    reports:
      junit: test/reports/junit.xml
      coverage_report:
        coverage_format: cobertura
        path: test/coverage.xml
```

### Jenkins Pipeline
```groovy
stage('Docker Tests') {
    steps {
        dir('test') {
            sh 'make ci-full'
        }
    }
    post {
        always {
            publishHTML([
                allowMissing: false,
                alwaysLinkToLastBuild: true,
                keepAll: true,
                reportDir: 'test/reports',
                reportFiles: 'report.html',
                reportName: 'Test Report'
            ])
        }
    }
}
```

## 🎯 Benefits

### Development Benefits
- **Consistent Environment**: Same results across all machines
- **Easy Onboarding**: New developers get running instantly
- **Dependency Isolation**: No conflicts with host system
- **Version Control**: Docker configuration is versioned

### Testing Benefits
- **Reproducible Results**: Eliminates "works on my machine" issues
- **Parallel Testing**: Run multiple test suites simultaneously
- **Clean State**: Fresh environment for each test run
- **Comprehensive Reporting**: Detailed test and coverage reports

### Operations Benefits
- **CI/CD Ready**: Drop-in solution for automated pipelines
- **Scalable Testing**: Easy to scale testing across infrastructure
- **Environment Parity**: Development matches production testing
- **Resource Management**: Controlled resource usage

## 📚 Additional Resources

### Documentation Links
- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Reference](https://docs.docker.com/compose/)
- [Pytest Documentation](https://docs.pytest.org/)
- [Python Testing Best Practices](https://docs.python-guide.org/writing/tests/)

### Project-Specific Docs
- [Test Framework README](./README.md)
- [Go Unit Tests](../gateway/README.md)
- [API Documentation](../docs/api.md)

This Docker setup provides a complete, production-ready testing environment that scales from local development to enterprise CI/CD pipelines.