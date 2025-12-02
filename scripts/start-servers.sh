#!/bin/bash
# Panchangam Server Startup Script
# Starts both the gRPC service and HTTP gateway

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
GRPC_PORT=${GRPC_PORT:-50051}
HTTP_PORT=${HTTP_PORT:-8080}
LOG_LEVEL=${LOG_LEVEL:-info}

echo -e "${GREEN}🚀 Starting Panchangam Services${NC}"
echo -e "${GREEN}================================${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

# Build the binaries
echo -e "${YELLOW}📦 Building binaries...${NC}"
go build -o ./bin/panchangam-server ./cmd/server/main.go
go build -o ./bin/panchangam-gateway ./cmd/gateway/main.go

# Check if builds were successful
if [ ! -f ./bin/panchangam-server ] || [ ! -f ./bin/panchangam-gateway ]; then
    echo -e "${RED}Error: Build failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Build successful${NC}"

# Create PID directory if it doesn't exist
mkdir -p ./tmp/pids

# Start gRPC server
echo -e "${YELLOW}🔧 Starting gRPC server on port ${GRPC_PORT}...${NC}"
./bin/panchangam-server --grpc-port=${GRPC_PORT} --log-level=${LOG_LEVEL} > ./tmp/grpc-server.log 2>&1 &
GRPC_PID=$!
echo $GRPC_PID > ./tmp/pids/grpc-server.pid

# Wait for gRPC server to be ready
echo -e "${YELLOW}⏳ Waiting for gRPC server to start...${NC}"
sleep 2

# Check if gRPC server is running
if ! ps -p $GRPC_PID > /dev/null; then
    echo -e "${RED}Error: gRPC server failed to start${NC}"
    cat ./tmp/grpc-server.log
    exit 1
fi

echo -e "${GREEN}✅ gRPC server started (PID: ${GRPC_PID})${NC}"

# Start HTTP gateway
echo -e "${YELLOW}🌐 Starting HTTP gateway on port ${HTTP_PORT}...${NC}"
./bin/panchangam-gateway --grpc-endpoint=localhost:${GRPC_PORT} --http-port=${HTTP_PORT} --log-level=${LOG_LEVEL} > ./tmp/http-gateway.log 2>&1 &
GATEWAY_PID=$!
echo $GATEWAY_PID > ./tmp/pids/http-gateway.pid

# Wait for HTTP gateway to be ready
echo -e "${YELLOW}⏳ Waiting for HTTP gateway to start...${NC}"
sleep 2

# Check if HTTP gateway is running
if ! ps -p $GATEWAY_PID > /dev/null; then
    echo -e "${RED}Error: HTTP gateway failed to start${NC}"
    cat ./tmp/http-gateway.log
    kill $GRPC_PID 2>/dev/null || true
    exit 1
fi

echo -e "${GREEN}✅ HTTP gateway started (PID: ${GATEWAY_PID})${NC}"

echo ""
echo -e "${GREEN}✨ All services started successfully!${NC}"
echo -e "${GREEN}====================================${NC}"
echo ""
echo -e "📡 ${YELLOW}gRPC Server:${NC} localhost:${GRPC_PORT} (PID: ${GRPC_PID})"
echo -e "🌐 ${YELLOW}HTTP Gateway:${NC} http://localhost:${HTTP_PORT} (PID: ${GATEWAY_PID})"
echo ""
echo -e "${YELLOW}Quick Test:${NC} curl http://localhost:${HTTP_PORT}/api/v1/health"
echo -e "${YELLOW}Stop Services:${NC} ./scripts/stop-servers.sh"
