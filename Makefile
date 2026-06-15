# Panchangam Project Makefile

# Variables
GO_VERSION := 1.23
NODE_VERSION := 18
DOCKER_REGISTRY := ghcr.io
PROJECT_NAME := panchangam
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT_SHA := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Docker image names
BACKEND_IMAGE := $(DOCKER_REGISTRY)/$(PROJECT_NAME)-backend
FRONTEND_IMAGE := $(DOCKER_REGISTRY)/$(PROJECT_NAME)-frontend

# Build flags
LDFLAGS := -ldflags="-w -s -X main.Version=$(VERSION) -X main.CommitSHA=$(COMMIT_SHA) -X main.BuildTime=$(BUILD_TIME)"

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

.PHONY: build run test clean demo proto deps help all

# Default target
all: build

# =============================================================================
# Building
# =============================================================================

# Build all components
build: build-backend build-frontend

# Build backend binaries
build-backend:
	@echo "$(BLUE)Building backend binaries...$(NC)"
	@mkdir -p bin
	@CGO_ENABLED=0 go build $(LDFLAGS) -o bin/panchangam-gateway ./cmd/gateway
	@CGO_ENABLED=0 go build $(LDFLAGS) -o bin/panchangam-grpc ./cmd/server

# Build frontend
build-frontend:
	@echo "$(BLUE)Building frontend...$(NC)"
	@cd ui && npm run build

# Build the demo client
build-demo:
	@echo "Building sunrise demo client..."
	go build -o bin/sunrise-demo ./cmd/sunrise-demo

# Build the CLI client
build-cli:
	@echo "Building panchangam CLI client..."
	go build -o bin/panchangam-cli ./cmd/panchangam-cli

# Build the simple sunrise client
build-simple:
	@echo "Building simple sunrise client..."
	go build -o bin/sunrise-simple ./cmd/sunrise-simple

# Run the server
run:
	@echo "Starting panchangam server..."
	go run ./cmd/server

# Legacy server run (compatibility)
run_server:
	@echo "Starting panchangam server..."
	go run ./cmd/server

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -cover ./...

# Run only astronomy tests
test-astronomy:
	@echo "Running astronomy package tests..."
	go test -v ./astronomy

# Run historical validation tests
test-validation:
	@echo "Running historical validation tests..."
	go test ./astronomy -run TestHistoricalValidation -v

# Run the demo client with default settings
demo:
	@echo "Running sunrise demo (default: New York)..."
	go run ./cmd/sunrise-demo

# Run the CLI client (validate connection)
cli:
	@echo "Running CLI client validation..."
	go run ./cmd/panchangam-cli validate

# Run the simple sunrise client
simple:
	@echo "Running simple sunrise client (London)..."
	go run ./cmd/sunrise-simple -location london

# Run demo with London
demo-london:
	@echo "Running sunrise demo for London..."
	go run ./cmd/sunrise-demo -location london

# Run demo with Tokyo
demo-tokyo:
	@echo "Running sunrise demo for Tokyo..."
	go run ./cmd/sunrise-demo -location tokyo

# Run interactive demo examples
demo-interactive:
	@echo "Running interactive demo examples..."
	./scripts/demo-examples.sh

# Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/panchangam.proto

# Legacy proto generation (compatibility)
gen: proto

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# =============================================================================
# Testing
# =============================================================================

# Run frontend tests
test-frontend:
	@echo "$(BLUE)Running frontend tests...$(NC)"
	@cd ui && npm run test

# Run integration tests
test-integration:
	@echo "$(BLUE)Running integration tests...$(NC)"
	@go test -tags=integration ./...

# Run end-to-end tests
test-e2e:
	@echo "$(BLUE)Running end-to-end tests...$(NC)"
	@cd ui && npm run test:e2e

# =============================================================================
# Code Quality
# =============================================================================

# Clean build artifacts
clean:
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	rm -rf bin/
	rm -rf ui/dist/
	rm -f panchangam sunrise-demo
	rm -rf proto/panchangam/
	rm -f coverage.out coverage.html

# Format code
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	go fmt ./...
	@cd ui && npm run format

# Legacy format (compatibility)
format: fmt

# Run linter (if available)
lint:
	@echo "$(BLUE)Running linter...$(NC)"
	@go vet ./...
	@gofmt -s -l . | grep -v vendor | tee /dev/stderr | test -z "$$(cat)"
	@cd ui && npm run lint
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "WARN: golangci-lint not found, skipping additional lint checks"; \
	fi

# Run security scans
security:
	@echo "$(BLUE)Running security scans...$(NC)"
	@command -v gosec >/dev/null 2>&1 || go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@gosec ./...
	@cd ui && npm audit --audit-level=high

# Run all checks (format, lint, test)
check: fmt lint test

# =============================================================================
# Development
# =============================================================================

# Development server with auto-reload (requires air)
dev:
	@echo "Starting development server with auto-reload..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "WARN: air not found, falling back to regular run"; \
		make run; \
	fi

# Start development servers
dev-start:
	@echo "$(BLUE)Starting development servers...$(NC)"
	@trap 'kill 0' SIGINT; \
	(go run ./cmd/server) & \
	(sleep 3 && go run ./cmd/gateway) & \
	(cd ui && npm run dev) & \
	wait

# Set up development environment
dev-setup:
	@echo "$(BLUE)Setting up development environment...$(NC)"
	@go mod download
	@cd ui && npm install

# =============================================================================
# Docker
# =============================================================================

# Build all Docker images
docker-build: docker-build-backend docker-build-frontend

# Build backend Docker image
docker-build-backend:
	@echo "$(BLUE)Building backend Docker image...$(NC)"
	@docker build -f docker/Dockerfile.backend -t $(BACKEND_IMAGE):$(VERSION) -t $(BACKEND_IMAGE):latest .

# Build frontend Docker image
docker-build-frontend:
	@echo "$(BLUE)Building frontend Docker image...$(NC)"
	@docker build -f ui/Dockerfile -t $(FRONTEND_IMAGE):$(VERSION) -t $(FRONTEND_IMAGE):latest ui/

# Push Docker images to registry
docker-push:
	@echo "$(BLUE)Pushing Docker images...$(NC)"
	@docker push $(BACKEND_IMAGE):$(VERSION)
	@docker push $(BACKEND_IMAGE):latest
	@docker push $(FRONTEND_IMAGE):$(VERSION)
	@docker push $(FRONTEND_IMAGE):latest

# Run application with Docker Compose
docker-run:
	@echo "$(BLUE)Starting application with Docker Compose...$(NC)"
	@docker-compose -f docker-compose.prod.yml up --build

# Stop Docker Compose
docker-stop:
	@echo "$(BLUE)Stopping Docker Compose...$(NC)"
	@docker-compose -f docker-compose.prod.yml down

# Clean Docker images and containers
docker-clean:
	@echo "$(BLUE)Cleaning Docker images and containers...$(NC)"
	@docker-compose -f docker-compose.prod.yml down --volumes --remove-orphans
	@docker image prune -f
	@docker volume prune -f

# =============================================================================
# Deployment
# =============================================================================

# Deploy to development using Docker Compose
deploy-dev:
	@echo "$(BLUE)Deploying to development...$(NC)"
	@./deployments/scripts/deploy.sh development docker-compose

# Deploy to staging using Docker Compose
deploy-staging:
	@echo "$(BLUE)Deploying to staging...$(NC)"
	@./deployments/scripts/deploy.sh staging docker-compose

# Deploy to production using Kubernetes
deploy-production:
	@echo "$(BLUE)Deploying to production...$(NC)"
	@./deployments/scripts/deploy.sh production kubernetes

# Deploy to Kubernetes (any environment)
deploy-k8s:
	@echo "$(BLUE)Deploying to Kubernetes ($(ENVIRONMENT))...$(NC)"
	@cd deployments/k8s/overlays/$(ENVIRONMENT) && kustomize build . | kubectl apply -f -

# Deploy to Docker Compose (production)
deploy-compose:
	@echo "$(BLUE)Deploying with Docker Compose...$(NC)"
	@docker-compose -f docker-compose.prod.yml up -d --remove-orphans

# Monitoring & Health Checks
# =============================================================================

# Check service health
health-check:
	@echo "$(BLUE)Checking service health...$(NC)"
	@curl -f http://localhost:8080/api/v1/health || echo "$(RED)API Gateway is down$(NC)"
	@curl -f http://localhost:80/health || echo "$(RED)Frontend is down$(NC)"

# View logs
logs:
	@echo "$(BLUE)Viewing service logs...$(NC)"
	@docker-compose -f docker-compose.prod.yml logs -f --tail=100

# View Kubernetes logs
logs-k8s:
	@echo "$(BLUE)Viewing Kubernetes logs...$(NC)"
	@kubectl logs -f deployment/panchangam-gateway -n panchangam

# =============================================================================
# Kubernetes Operations
# =============================================================================

# Apply Kubernetes manifests
k8s-apply:
	@echo "$(BLUE)Applying Kubernetes manifests...$(NC)"
	@kubectl apply -f deployments/k8s/base/

# Delete Kubernetes resources
k8s-delete:
	@echo "$(RED)Deleting Kubernetes resources...$(NC)"
	@kubectl delete -f deployments/k8s/base/

# Get Kubernetes pod status
k8s-status:
	@echo "$(BLUE)Kubernetes pod status:$(NC)"
	@kubectl get pods -n panchangam

# Describe Kubernetes deployment
k8s-describe:
	@echo "$(BLUE)Kubernetes deployment details:$(NC)"
	@kubectl describe deployment -n panchangam

# Scale Kubernetes deployment
k8s-scale:
	@echo "$(BLUE)Scaling deployment to $(replicas) replicas...$(NC)"
	@kubectl scale deployment/panchangam-gateway --replicas=$(replicas) -n panchangam
	@kubectl scale deployment/panchangam-grpc --replicas=$(replicas) -n panchangam

# Rollback Kubernetes deployment
k8s-rollback:
	@echo "$(YELLOW)Rolling back Kubernetes deployment...$(NC)"
	@kubectl rollout undo deployment/panchangam-gateway -n panchangam
	@kubectl rollout undo deployment/panchangam-grpc -n panchangam

# Restart Kubernetes deployment
k8s-restart:
	@echo "$(BLUE)Restarting Kubernetes deployment...$(NC)"
	@kubectl rollout restart deployment/panchangam-gateway -n panchangam
	@kubectl rollout restart deployment/panchangam-grpc -n panchangam

# =============================================================================
# Infrastructure Setup
# =============================================================================

# Set up production infrastructure
infra-setup:
	@echo "$(BLUE)Setting up production infrastructure...$(NC)"
	@mkdir -p deployments/{prometheus/alerts,grafana/{provisioning,dashboards},alertmanager,loki,nginx}
	@echo "$(GREEN)PASS: Infrastructure directories created$(NC)"

# Validate deployment configurations
infra-validate:
	@echo "$(BLUE)Validating deployment configurations...$(NC)"
	@docker-compose -f docker-compose.prod.yml config
	@echo "$(GREEN)PASS: Docker Compose configuration is valid$(NC)"

# Generate SSL certificates (self-signed for development)
ssl-generate:
	@echo "$(BLUE)Generating self-signed SSL certificates...$(NC)"
	@mkdir -p deployments/nginx/ssl
	@openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
		-keyout deployments/nginx/ssl/privkey.pem \
		-out deployments/nginx/ssl/fullchain.pem \
		-subj "/C=US/ST=State/L=City/O=Organization/CN=panchangam.local"
	@echo "$(GREEN)PASS: SSL certificates generated$(NC)"

# =============================================================================
# Local Verification Commands
# =============================================================================

# Local lint gate
ci-lint:
	@echo "$(BLUE)Running local lint gate...$(NC)"
	@$(MAKE) lint
	@$(MAKE) security

# Local test gate
ci-test:
	@echo "$(BLUE)Running local test gate...$(NC)"
	@$(MAKE) test
	@$(MAKE) test-frontend
	@$(MAKE) test-integration

# Local build gate
ci-build:
	@echo "$(BLUE)Running local build gate...$(NC)"
	@$(MAKE) build
	@$(MAKE) docker-build

# =============================================================================
# Utilities
# =============================================================================

# Show version information
version:
	@echo "$(BLUE)Version Information:$(NC)"
	@echo "  Version: $(VERSION)"
	@echo "  Commit: $(COMMIT_SHA)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  Go Version: $$(go version | cut -d' ' -f3)"
	@echo "  Node Version: $$(cd ui && node --version 2>/dev/null || echo 'Not available')"

# Legacy docker start (compatibility)
start:
	@echo "Starting services..."
	docker-compose -f docker-compose.prod.yml up --force-recreate --remove-orphans --detach
	@echo ""
	@echo "Panchangam services are running."
	@echo "Frontend: http://localhost"
	@echo "API health: http://localhost:8080/api/v1/health"
	@echo "Jaeger traces: http://localhost:16686/jaeger/ui"

# Show help
help:
	@echo "Panchangam Development Commands:"
	@echo ""
	@echo "Build Commands:"
	@echo "  make build          - Build the main server"
	@echo "  make build-demo     - Build the demo client"
	@echo "  make build-cli      - Build the CLI client"
	@echo "  make build-simple   - Build the simple sunrise client"
	@echo ""
	@echo "Run Commands:"
	@echo "  make run            - Start the panchangam server"
	@echo "  make demo           - Run demo client (New York)"
	@echo "  make cli            - Run CLI client validation"
	@echo "  make simple         - Run simple sunrise client (London)"
	@echo "  make demo-london    - Run demo client (London)"
	@echo "  make demo-tokyo     - Run demo client (Tokyo)"
	@echo "  make demo-interactive - Run interactive demo examples"
	@echo ""
	@echo "Test Commands:"
	@echo "  make test           - Run all tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make test-astronomy - Run astronomy package tests"
	@echo "  make test-validation - Run historical validation tests"
	@echo ""
	@echo "Development Commands:"
	@echo "  make proto          - Generate protobuf files"
	@echo "  make deps           - Install dependencies"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter"
	@echo "  make check          - Run all checks (fmt, lint, test)"
	@echo "  make dev            - Start development server with auto-reload"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run Docker container"
	@echo ""
	@echo "Utility Commands:"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make help           - Show this help"
	@echo ""
	@echo "Quick Start:"
	@echo "  1. make run          # Start server"
	@echo "  2. make cli          # Test CLI client"
	@echo "  3. make simple       # Test simple client"
	@echo "  4. make test-validation # Validate accuracy"
