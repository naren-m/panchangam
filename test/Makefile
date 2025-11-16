# Makefile for Panchangam API Testing Framework
# Provides convenient commands for Docker-based testing

.PHONY: help build test test-smoke test-integration test-performance test-security
.PHONY: test-all test-verbose test-coverage test-html shell logs clean
.PHONY: install-deps check-deps lint format security-scan
.PHONY: test-data-validation test-cache test-e2e test-issue-81 test-load

# Default target
.DEFAULT_GOAL := help

# Variables
DOCKER_COMPOSE = docker-compose
CONTAINER_NAME = panchangam-tests
PROJECT_NAME = panchangam

# Colors for output
BLUE = \033[0;34m
GREEN = \033[0;32m
YELLOW = \033[1;33m
RED = \033[0;31m
NC = \033[0m # No Color

## Help
help: ## Show this help message
	@echo "$(BLUE)Panchangam API Testing Framework$(NC)"
	@echo "================================="
	@echo ""
	@echo "$(GREEN)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Examples:$(NC)"
	@echo "  make test                  # Run all tests"
	@echo "  make test-smoke            # Run smoke tests only"
	@echo "  make test-coverage         # Run tests with coverage"
	@echo "  make shell                 # Open interactive shell"
	@echo "  make clean                 # Clean up containers"

## Docker Operations
build: ## Build Docker images
	@echo "$(BLUE)Building Docker images...$(NC)"
	$(DOCKER_COMPOSE) build --parallel

rebuild: ## Force rebuild Docker images
	@echo "$(BLUE)Force rebuilding Docker images...$(NC)"
	$(DOCKER_COMPOSE) build --no-cache --parallel

up: ## Start services in background
	@echo "$(BLUE)Starting services...$(NC)"
	$(DOCKER_COMPOSE) up -d

down: ## Stop and remove containers
	@echo "$(BLUE)Stopping services...$(NC)"
	$(DOCKER_COMPOSE) down --remove-orphans

## Testing Commands
test: build ## Run all tests
	@echo "$(BLUE)Running all tests...$(NC)"
	./run-docker-tests.sh all

test-smoke: build ## Run smoke tests only
	@echo "$(BLUE)Running smoke tests...$(NC)"
	./run-docker-tests.sh smoke

test-integration: build ## Run integration tests only
	@echo "$(BLUE)Running integration tests...$(NC)"
	./run-docker-tests.sh integration

test-performance: build ## Run performance tests only
	@echo "$(BLUE)Running performance tests...$(NC)"
	./run-docker-tests.sh performance

test-security: build ## Run security tests only
	@echo "$(BLUE)Running security tests...$(NC)"
	./run-docker-tests.sh security

test-data-validation: build ## Run data validation tests (Issue #81)
	@echo "$(BLUE)Running data validation tests...$(NC)"
	$(DOCKER_COMPOSE) run --rm $(CONTAINER_NAME) pytest -v -m data_validation

test-cache: build ## Run cache integration tests (Issue #81)
	@echo "$(BLUE)Running cache integration tests...$(NC)"
	$(DOCKER_COMPOSE) run --rm $(CONTAINER_NAME) pytest -v -m cache test_cache_integration.py

test-e2e: build ## Run end-to-end integration tests (Issue #81)
	@echo "$(BLUE)Running end-to-end integration tests...$(NC)"
	$(DOCKER_COMPOSE) run --rm $(CONTAINER_NAME) pytest -v -m e2e test_e2e_integration.py

test-issue-81: build ## Run all Issue #81 specific tests
	@echo "$(BLUE)Running Issue #81 comprehensive test suite...$(NC)"
	@echo "$(YELLOW)1. Data Validation Tests$(NC)"
	$(DOCKER_COMPOSE) run --rm $(CONTAINER_NAME) pytest -v test_data_validation.py
	@echo "$(YELLOW)2. Performance Tests$(NC)"
	$(DOCKER_COMPOSE) run --rm $(CONTAINER_NAME) pytest -v test_performance.py
	@echo "$(YELLOW)3. Cache Integration Tests$(NC)"
	$(DOCKER_COMPOSE) run --rm $(CONTAINER_NAME) pytest -v test_cache_integration.py
	@echo "$(YELLOW)4. End-to-End Tests$(NC)"
	$(DOCKER_COMPOSE) run --rm $(CONTAINER_NAME) pytest -v test_e2e_integration.py
	@echo "$(GREEN)Issue #81 test suite completed!$(NC)"

test-load: build ## Run load tests with Locust (Issue #81)
	@echo "$(BLUE)Running load tests with Locust...$(NC)"
	@echo "$(YELLOW)Starting Locust web interface on http://localhost:8089$(NC)"
	@echo "$(YELLOW)Use Ctrl+C to stop$(NC)"
	$(DOCKER_COMPOSE) run --rm -p 8089:8089 $(CONTAINER_NAME) \
		locust -f locustfile.py --host=http://localhost:8080 --web-host=0.0.0.0

test-load-headless: build ## Run load tests headless with 50 users (Issue #81 requirement)
	@echo "$(BLUE)Running headless load test (50 concurrent users)...$(NC)"
	$(DOCKER_COMPOSE) run --rm $(CONTAINER_NAME) \
		locust -f locustfile.py --host=http://localhost:8080 \
		--headless --users 50 --spawn-rate 10 --run-time 60s

test-all: build ## Run all tests with verbose output
	@echo "$(BLUE)Running all tests with verbose output...$(NC)"
	./run-docker-tests.sh all --verbose

test-verbose: build ## Run tests with verbose output
	@echo "$(BLUE)Running tests with verbose output...$(NC)"
	./run-docker-tests.sh all --verbose

test-coverage: build ## Run tests with coverage report
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	./run-docker-tests.sh all --coverage

test-html: build ## Run tests with HTML reports
	@echo "$(BLUE)Running tests with HTML reports...$(NC)"
	./run-docker-tests.sh all --coverage --html-report

test-full: build ## Run complete test suite with all reports
	@echo "$(BLUE)Running complete test suite...$(NC)"
	./run-docker-tests.sh all --verbose --coverage --html-report

## Development Commands
shell: build ## Open interactive shell in test container
	@echo "$(BLUE)Opening shell in test container...$(NC)"
	./run-docker-tests.sh --shell

logs: ## Show container logs
	@echo "$(BLUE)Showing container logs...$(NC)"
	./run-docker-tests.sh --logs

watch: ## Watch test files and run smoke tests on changes
	@echo "$(BLUE)Watching for file changes...$(NC)"
	@while true; do \
		inotifywait -e modify,create,delete -r . --exclude '\.git|__pycache__|\.pyc' 2>/dev/null || true; \
		echo "$(YELLOW)Files changed, running smoke tests...$(NC)"; \
		make test-smoke; \
		sleep 2; \
	done

## Code Quality
lint: build ## Run code linting
	@echo "$(BLUE)Running code linting...$(NC)"
	$(DOCKER_COMPOSE) run --rm $(CONTAINER_NAME) bash -c "\
		flake8 . --count --select=E9,F63,F7,F82 --show-source --statistics && \
		flake8 . --count --exit-zero --max-complexity=10 --max-line-length=127 --statistics"

format: build ## Format code with black and isort
	@echo "$(BLUE)Formatting code...$(NC)"
	$(DOCKER_COMPOSE) run --rm $(CONTAINER_NAME) bash -c "\
		black . && \
		isort ."

type-check: build ## Run type checking with mypy
	@echo "$(BLUE)Running type checking...$(NC)"
	$(DOCKER_COMPOSE) run --rm $(CONTAINER_NAME) mypy .

security-scan: build ## Run security scanning
	@echo "$(BLUE)Running security scan...$(NC)"
	$(DOCKER_COMPOSE) run --rm $(CONTAINER_NAME) bash -c "\
		safety check && \
		bandit -r ."

check-deps: ## Check for dependency vulnerabilities
	@echo "$(BLUE)Checking dependencies...$(NC)"
	$(DOCKER_COMPOSE) run --rm $(CONTAINER_NAME) safety check

## Reporting
report: test-html ## Generate and open test reports
	@echo "$(BLUE)Generating reports...$(NC)"
	@if [ -f "./reports/report.html" ]; then \
		echo "$(GREEN)Test report: ./reports/report.html$(NC)"; \
		command -v open >/dev/null 2>&1 && open ./reports/report.html || true; \
	fi
	@if [ -f "./htmlcov/index.html" ]; then \
		echo "$(GREEN)Coverage report: ./htmlcov/index.html$(NC)"; \
		command -v open >/dev/null 2>&1 && open ./htmlcov/index.html || true; \
	fi

coverage-report: test-coverage ## Generate coverage report only
	@echo "$(BLUE)Generating coverage report...$(NC)"
	@if [ -f "./htmlcov/index.html" ]; then \
		echo "$(GREEN)Coverage report: ./htmlcov/index.html$(NC)"; \
		command -v open >/dev/null 2>&1 && open ./htmlcov/index.html || true; \
	fi

## Cleanup
clean: ## Clean up containers and images
	@echo "$(BLUE)Cleaning up...$(NC)"
	$(DOCKER_COMPOSE) down --remove-orphans --volumes
	docker system prune -f
	@echo "$(GREEN)Cleanup completed$(NC)"

clean-all: ## Clean up everything including images
	@echo "$(BLUE)Cleaning up everything...$(NC)"
	$(DOCKER_COMPOSE) down --remove-orphans --volumes --rmi all
	docker system prune -a -f
	@echo "$(GREEN)Complete cleanup finished$(NC)"

clean-reports: ## Clean test reports and coverage files
	@echo "$(BLUE)Cleaning reports...$(NC)"
	rm -rf ./reports ./htmlcov ./logs
	rm -f .coverage coverage.xml pytest_cache
	@echo "$(GREEN)Reports cleaned$(NC)"

## Installation
install-deps: ## Install dependencies (local development)
	@echo "$(BLUE)Installing Python dependencies...$(NC)"
	pip install -r requirements.txt
	pip install -r requirements-dev.txt

## CI/CD
ci-test: build ## Run tests in CI mode
	@echo "$(BLUE)Running CI tests...$(NC)"
	$(DOCKER_COMPOSE) run --rm $(CONTAINER_NAME) python run_tests.py --verbose --coverage --html-report

ci-full: ## Run complete CI pipeline
	@echo "$(BLUE)Running complete CI pipeline...$(NC)"
	make lint
	make type-check
	make security-scan
	make ci-test

## Status
status: ## Show Docker container status
	@echo "$(BLUE)Container status:$(NC)"
	$(DOCKER_COMPOSE) ps

ps: status ## Alias for status

## Quick Actions
quick-test: ## Quick smoke test
	@echo "$(BLUE)Quick smoke test...$(NC)"
	./run-docker-tests.sh smoke

quick-check: ## Quick health check
	@echo "$(BLUE)Quick health check...$(NC)"
	$(DOCKER_COMPOSE) run --rm $(CONTAINER_NAME) python -c "import requests; print('✅ Modules available')"

## Development Workflow
dev-setup: build install-deps ## Setup development environment
	@echo "$(GREEN)Development environment ready!$(NC)"
	@echo "Run 'make test' to execute tests"
	@echo "Run 'make shell' for interactive development"

dev-test: ## Development test cycle (format, lint, test)
	@echo "$(BLUE)Development test cycle...$(NC)"
	make format
	make lint
	make test-smoke

## Documentation
docs: ## Generate documentation
	@echo "$(BLUE)Generating documentation...$(NC)"
	$(DOCKER_COMPOSE) run --rm $(CONTAINER_NAME) bash -c "\
		sphinx-build -b html docs docs/_build/html"

## Environment Info
info: ## Show environment information
	@echo "$(BLUE)Environment Information:$(NC)"
	@echo "Docker version: $$(docker --version)"
	@echo "Docker Compose version: $$(docker-compose --version)"
	@echo "Project directory: $$(pwd)"
	@echo "Container status:"
	@$(DOCKER_COMPOSE) ps || echo "No containers running"