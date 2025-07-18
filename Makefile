# Panchangam Build and Development Makefile

.PHONY: build run test clean demo proto deps help all

# Default target
all: build

# Build the main server
build:
	@echo "🔨 Building panchangam server..."
	go build -o bin/panchangam .

# Build the demo client
build-demo:
	@echo "🔨 Building sunrise demo client..."
	go build -o bin/sunrise-demo cmd/sunrise-demo/main.go

# Build the CLI client
build-cli:
	@echo "🔨 Building panchangam CLI client..."
	go build -o bin/panchangam-cli cmd/panchangam-cli/main.go

# Build the simple sunrise client
build-simple:
	@echo "🔨 Building simple sunrise client..."
	go build -o bin/sunrise-simple cmd/sunrise-simple/main.go

# Run the server
run:
	@echo "🚀 Starting panchangam server..."
	go run main.go

# Legacy server run (compatibility)
run_server:
	@echo "🚀 Starting panchangam server..."
	go run main.go

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	go test -v -cover ./...

# Run only astronomy tests
test-astronomy:
	@echo "🌅 Running astronomy package tests..."
	go test -v ./astronomy

# Run historical validation tests
test-validation:
	@echo "📊 Running historical validation tests..."
	go test ./astronomy -run TestHistoricalValidation -v

# Run the demo client with default settings
demo:
	@echo "🌅 Running sunrise demo (default: New York)..."
	go run cmd/sunrise-demo/main.go

# Run the CLI client (validate connection)
cli:
	@echo "🖥️  Running CLI client validation..."
	go run cmd/panchangam-cli/main.go validate

# Run the simple sunrise client
simple:
	@echo "🌅 Running simple sunrise client (London)..."
	go run cmd/sunrise-simple/main.go -location london

# Run demo with London
demo-london:
	@echo "🌅 Running sunrise demo for London..."
	go run cmd/sunrise-demo/main.go -location london

# Run demo with Tokyo
demo-tokyo:
	@echo "🌅 Running sunrise demo for Tokyo..."
	go run cmd/sunrise-demo/main.go -location tokyo

# Run interactive demo examples
demo-interactive:
	@echo "🎬 Running interactive demo examples..."
	./scripts/demo-examples.sh

# Generate protobuf files
proto:
	@echo "🔧 Generating protobuf files..."
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/panchangam.proto

# Legacy proto generation (compatibility)
gen: proto

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -f panchangam sunrise-demo
	rm -rf proto/panchangam/

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Legacy format (compatibility)
format: fmt

# Run linter (if available)
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found, skipping lint check"; \
	fi

# Run all checks (format, lint, test)
check: fmt lint test

# Development server with auto-reload (requires air)
dev:
	@echo "🔄 Starting development server with auto-reload..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "⚠️  air not found, falling back to regular run"; \
		make run; \
	fi

# Docker build
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t panchangam .

# Docker run
docker-run:
	@echo "🐳 Running Docker container..."
	docker run -p 8080:8080 panchangam

# Legacy docker start (compatibility)
start:
	@echo "🐳 Starting services..."
	docker compose up --force-recreate --remove-orphans --detach
	@echo ""
	@echo "OpenTelemetry Demo is running."
	@echo "Go to http://192.168.68.73:16686/ for the demo UI."
	@echo "Go to http://localhost:16686/jaeger/ui for the Jaeger UI."
	@echo "Go to http://localhost:8080/grafana/ for the Grafana UI."

# Show help
help:
	@echo "📖 Panchangam Development Commands:"
	@echo ""
	@echo "🏗️  Build Commands:"
	@echo "  make build          - Build the main server"
	@echo "  make build-demo     - Build the demo client"
	@echo "  make build-cli      - Build the CLI client"
	@echo "  make build-simple   - Build the simple sunrise client"
	@echo ""
	@echo "🚀 Run Commands:"
	@echo "  make run            - Start the panchangam server"
	@echo "  make demo           - Run demo client (New York)"
	@echo "  make cli            - Run CLI client validation"
	@echo "  make simple         - Run simple sunrise client (London)"
	@echo "  make demo-london    - Run demo client (London)"
	@echo "  make demo-tokyo     - Run demo client (Tokyo)"
	@echo "  make demo-interactive - Run interactive demo examples"
	@echo ""
	@echo "🧪 Test Commands:"
	@echo "  make test           - Run all tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make test-astronomy - Run astronomy package tests"
	@echo "  make test-validation - Run historical validation tests"
	@echo ""
	@echo "🔧 Development Commands:"
	@echo "  make proto          - Generate protobuf files"
	@echo "  make deps           - Install dependencies"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter"
	@echo "  make check          - Run all checks (fmt, lint, test)"
	@echo "  make dev            - Start development server with auto-reload"
	@echo ""
	@echo "🐳 Docker Commands:"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run Docker container"
	@echo ""
	@echo "🧹 Utility Commands:"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make help           - Show this help"
	@echo ""
	@echo "📚 Quick Start:"
	@echo "  1. make run          # Start server"
	@echo "  2. make cli          # Test CLI client"
	@echo "  3. make simple       # Test simple client"
	@echo "  4. make test-validation # Validate accuracy"