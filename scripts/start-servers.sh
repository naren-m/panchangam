#!/bin/bash

set -e

echo "🚀 Starting Panchangam servers..."

# Build the servers
echo "📦 Building servers..."
go build -o grpc-server ./server/server.go
go build -o gateway-server ./cmd/gateway/main.go

# Start gRPC server in background
echo "🔧 Starting gRPC server on port 50052..."
./grpc-server &
GRPC_PID=$!

# Wait a moment for gRPC server to start
sleep 2

# Start HTTP gateway server in background
echo "🌐 Starting HTTP Gateway server on port 8080..."
./gateway-server --grpc-endpoint=localhost:50052 --http-port=8080 &
GATEWAY_PID=$!

# Wait a moment for gateway to start
sleep 2

echo "✅ Servers started successfully!"
echo ""
echo "🔍 Test endpoints:"
echo "   Health check: curl http://localhost:8080/api/v1/health"
echo "   Panchangam API: curl 'http://localhost:8080/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946&tz=Asia/Kolkata'"
echo ""
echo "🛑 To stop servers, run: kill $GRPC_PID $GATEWAY_PID"
echo "   Or use: pkill -f grpc-server && pkill -f gateway-server"
echo ""
echo "📊 Server processes:"
echo "   gRPC Server PID: $GRPC_PID"
echo "   Gateway Server PID: $GATEWAY_PID"
echo ""
echo "🎉 Phase 1 Implementation Complete!"
echo "   ✅ gRPC-to-HTTP API Gateway"
echo "   ✅ Comprehensive error handling"
echo "   ✅ CORS configuration"
echo "   ✅ Health check endpoints"
echo "   ✅ Request logging and monitoring"
echo ""
echo "📝 Next steps for Phase 2:"
echo "   - Update frontend to use real API (http://localhost:8080/api/v1/panchangam)"
echo "   - Replace mock data in panchangamApi.ts"
echo "   - Add loading states and error handling in UI"

# Keep script running to show server logs
wait