#!/bin/bash
# Panchangam Server Stop Script

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}🛑 Stopping Panchangam Services${NC}"

# Stop HTTP gateway
if [ -f ./tmp/pids/http-gateway.pid ]; then
    GATEWAY_PID=$(cat ./tmp/pids/http-gateway.pid)
    if ps -p $GATEWAY_PID > /dev/null; then
        echo -e "${YELLOW}Stopping HTTP gateway (PID: ${GATEWAY_PID})...${NC}"
        kill $GATEWAY_PID
        echo -e "${GREEN}✅ HTTP gateway stopped${NC}"
    fi
    rm ./tmp/pids/http-gateway.pid
fi

# Stop gRPC server
if [ -f ./tmp/pids/grpc-server.pid ]; then
    GRPC_PID=$(cat ./tmp/pids/grpc-server.pid)
    if ps -p $GRPC_PID > /dev/null; then
        echo -e "${YELLOW}Stopping gRPC server (PID: ${GRPC_PID})...${NC}"
        kill $GRPC_PID
        echo -e "${GREEN}✅ gRPC server stopped${NC}"
    fi
    rm ./tmp/pids/grpc-server.pid
fi

echo -e "${GREEN}✨ All services stopped${NC}"
