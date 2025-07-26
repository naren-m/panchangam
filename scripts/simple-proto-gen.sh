#!/bin/bash

set -e

echo "Generating Protocol Buffer files..."

# Generate basic protobuf files first
protoc \
    --proto_path=proto \
    --proto_path=/usr/local/include \
    --go_out=proto \
    --go_opt=paths=source_relative \
    --go-grpc_out=proto \
    --go-grpc_opt=paths=source_relative \
    proto/panchangam.proto

echo "Basic protobuf files generated successfully!"
echo "Note: gRPC-gateway generation skipped due to googleapis dependency issues"
echo "The HTTP gateway will be implemented manually using the existing gRPC client"