#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Generating Protocol Buffer files...${NC}"

# Create necessary directories
mkdir -p proto/panchangam

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo -e "${RED}Error: protoc is not installed${NC}"
    echo "Please install Protocol Buffers compiler:"
    echo "  macOS: brew install protobuf"
    echo "  Ubuntu: sudo apt install protobuf-compiler"
    exit 1
fi

# Check if required Go plugins are installed
echo -e "${YELLOW}Checking Go protobuf plugins...${NC}"

if ! command -v protoc-gen-go &> /dev/null; then
    echo "Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

if ! command -v protoc-gen-grpc-gateway &> /dev/null; then
    echo "Installing protoc-gen-grpc-gateway..."
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
fi

if ! command -v protoc-gen-openapiv2 &> /dev/null; then
    echo "Installing protoc-gen-openapiv2..."
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
fi

# Find Go installation and module paths
GO_PATH=$(go env GOPATH)
GO_MOD_PATH=$(go env GOMODCACHE)

# Common include paths for googleapis
GOOGLEAPIS_PATH=""
if [ -d "$GO_MOD_PATH/github.com/googleapis/googleapis@"* ]; then
    GOOGLEAPIS_PATH=$(find "$GO_MOD_PATH/github.com/googleapis/googleapis@"* -type d -name "googleapis" | head -1)
fi

if [ -z "$GOOGLEAPIS_PATH" ]; then
    echo -e "${YELLOW}Warning: googleapis not found in module cache. Trying to download...${NC}"
    go mod download github.com/googleapis/googleapis
    GOOGLEAPIS_PATH=$(find "$GO_MOD_PATH/github.com/googleapis/googleapis@"* -type d -name "googleapis" | head -1)
fi

# Alternative: try grpc-gateway path
if [ -z "$GOOGLEAPIS_PATH" ]; then
    GRPC_GATEWAY_PATH=$(find "$GO_MOD_PATH/github.com/grpc-ecosystem/grpc-gateway"* -type d -name "third_party" | head -1)
    if [ -n "$GRPC_GATEWAY_PATH" ]; then
        GOOGLEAPIS_PATH="$GRPC_GATEWAY_PATH/googleapis"
    fi
fi

if [ -z "$GOOGLEAPIS_PATH" ]; then
    echo -e "${RED}Error: Could not find googleapis proto files${NC}"
    echo "Please ensure googleapis are available:"
    echo "  go get google.golang.org/genproto/googleapis/api/annotations"
    exit 1
fi

echo -e "${GREEN}Found googleapis at: $GOOGLEAPIS_PATH${NC}"

# Generate the protocol buffer files
echo -e "${YELLOW}Generating protobuf files...${NC}"

protoc \
    --proto_path=proto \
    --proto_path="$GOOGLEAPIS_PATH" \
    --go_out=proto \
    --go_opt=paths=source_relative \
    --go-grpc_out=proto \
    --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=proto \
    --grpc-gateway_opt=paths=source_relative \
    --openapiv2_out=proto \
    --openapiv2_opt=logtostderr=true \
    proto/panchangam.proto

echo -e "${GREEN}Protocol buffer files generated successfully!${NC}"

# List generated files
echo -e "${YELLOW}Generated files:${NC}"
find proto -name "*.go" -o -name "*.json" -o -name "*.swagger.json" | sort

echo -e "${GREEN}Done!${NC}"