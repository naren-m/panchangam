# Panchangam Deployment Guide

Complete guide for deploying the Panchangam application with gRPC service, HTTP gateway, and frontend UI.

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Quick Start](#quick-start)
3. [Architecture](#architecture)
4. [Development Setup](#development-setup)
5. [Production Deployment](#production-deployment)
6. [Testing](#testing)
7. [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Software
- **Go 1.21+**: For backend services
- **Node.js 18+** and **npm/yarn**: For frontend UI
- **Protocol Buffers** compiler (optional): For proto regeneration

### Optional Tools
- **grpcurl**: For testing gRPC endpoints
- **Docker**: For containerized deployment
- **kubectl**: For Kubernetes deployment

## Quick Start

### Start All Services (Development)

```bash
# Clone and navigate to project
cd panchangam

# Start gRPC server and HTTP gateway
./scripts/start-servers.sh

# In a new terminal, start the frontend
cd ui
npm install
npm run dev
```

Access the application:
- Frontend UI: http://localhost:5173
- HTTP API: http://localhost:8080/api/v1/panchangam
- gRPC Service: localhost:50051

### Stop Services

```bash
./scripts/stop-servers.sh
```

## Architecture

```
┌─────────────┐         ┌─────────────┐         ┌──────────────┐
│   Browser   │────────▶│ HTTP Gateway│────────▶│ gRPC Service │
│   (React)   │         │  (Port 8080)│         │  (Port 50051)│
└─────────────┘         └─────────────┘         └──────────────┘
                              │                          │
                              │                          ▼
                              │                   ┌──────────────┐
                              │                   │  Astronomy   │
                              │                   │ Calculations │
                              ▼                   └──────────────┘
                        ┌─────────────┐
                        │    CORS &   │
                        │   Logging   │
                        └─────────────┘
```

### Component Roles

1. **gRPC Service** (`cmd/server/main.go`):
   - Core astronomical calculations
   - All 5 Panchangam elements (Tithi, Nakshatra, Yoga, Karana, Vara)
   - Swiss Ephemeris integration
   - OpenTelemetry observability

2. **HTTP Gateway** (`cmd/gateway/main.go`):
   - REST API wrapper around gRPC
   - CORS handling for web clients
   - Request/response transformation
   - Health check endpoint

3. **Frontend UI** (`ui/`):
   - React application with TypeScript
   - Calendar visualization
   - Real-time Panchangam data display
   - Location-based calculations

## Development Setup

### 1. Backend Services

#### Build Services
```bash
# Build gRPC server
go build -o bin/panchangam-server cmd/server/main.go

# Build HTTP gateway
go build -o bin/panchangam-gateway cmd/gateway/main.go
```

#### Run Services Manually

**Terminal 1 - gRPC Server:**
```bash
./bin/panchangam-server \
  --grpc-port=50051 \
  --log-level=info
```

**Terminal 2 - HTTP Gateway:**
```bash
./bin/panchangam-gateway \
  --grpc-endpoint=localhost:50051 \
  --http-port=8080 \
  --log-level=info
```

### 2. Frontend Development

```bash
cd ui

# Install dependencies
npm install

# Run development server with hot reload
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### 3. Environment Configuration

#### Backend (.env)
```bash
# gRPC Server
GRPC_PORT=50051
LOG_LEVEL=info

# HTTP Gateway
HTTP_PORT=8080
GRPC_ENDPOINT=localhost:50051
```

#### Frontend (ui/.env)
```bash
VITE_API_BASE_URL=http://localhost:8080
```

## Production Deployment

### Option 1: Systemd Services

#### 1. Create Service Files

**/etc/systemd/system/panchangam-grpc.service:**
```ini
[Unit]
Description=Panchangam gRPC Service
After=network.target

[Service]
Type=simple
User=panchangam
WorkingDirectory=/opt/panchangam
ExecStart=/opt/panchangam/bin/panchangam-server --grpc-port=50051
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

**/etc/systemd/system/panchangam-gateway.service:**
```ini
[Unit]
Description=Panchangam HTTP Gateway
After=panchangam-grpc.service
Requires=panchangam-grpc.service

[Service]
Type=simple
User=panchangam
WorkingDirectory=/opt/panchangam
ExecStart=/opt/panchangam/bin/panchangam-gateway \
  --grpc-endpoint=localhost:50051 \
  --http-port=8080
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

#### 2. Enable and Start Services
```bash
sudo systemctl enable panchangam-grpc
sudo systemctl enable panchangam-gateway

sudo systemctl start panchangam-grpc
sudo systemctl start panchangam-gateway

# Check status
sudo systemctl status panchangam-grpc
sudo systemctl status panchangam-gateway
```

### Option 2: Docker Deployment

#### 1. Build Docker Images

**Dockerfile (gRPC Server):**
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o panchangam-server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/panchangam-server /usr/local/bin/
EXPOSE 50051
CMD ["panchangam-server", "--grpc-port=50051"]
```

**docker-compose.yml:**
```yaml
version: '3.8'

services:
  grpc-server:
    build:
      context: .
      dockerfile: Dockerfile.grpc
    ports:
      - "50051:50051"
    environment:
      - LOG_LEVEL=info
    restart: unless-stopped

  http-gateway:
    build:
      context: .
      dockerfile: Dockerfile.gateway
    ports:
      - "8080:8080"
    environment:
      - GRPC_ENDPOINT=grpc-server:50051
      - LOG_LEVEL=info
    depends_on:
      - grpc-server
    restart: unless-stopped

  frontend:
    build:
      context: ./ui
      dockerfile: Dockerfile
    ports:
      - "80:80"
    environment:
      - VITE_API_BASE_URL=http://localhost:8080
    depends_on:
      - http-gateway
    restart: unless-stopped
```

#### 2. Deploy with Docker Compose
```bash
docker-compose up -d
```

### Option 3: Kubernetes Deployment

See `k8s/` directory for complete Kubernetes manifests.

**Quick deploy:**
```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/grpc-service.yaml
kubectl apply -f k8s/http-gateway.yaml
kubectl apply -f k8s/frontend.yaml
```

## Testing

### 1. Health Checks

**HTTP Gateway:**
```bash
curl http://localhost:8080/api/v1/health
```

Expected response:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "service": "panchangam-gateway",
  "version": "1.0.0"
}
```

**gRPC Service (requires grpcurl):**
```bash
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
```

### 2. API Testing

**Get Panchangam Data:**
```bash
curl "http://localhost:8080/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946&tz=Asia/Kolkata"
```

Expected response structure:
```json
{
  "date": "2024-01-15",
  "tithi": "Chaturthi (4)",
  "nakshatra": "Uttara Bhadrapada (26)",
  "yoga": "Siddha (21)",
  "karana": "Gara (6)",
  "sunrise_time": "01:15:32",
  "sunset_time": "12:41:47",
  "events": [...]
}
```

### 3. Performance Testing

```bash
# Load test with 100 concurrent requests
ab -n 1000 -c 100 "http://localhost:8080/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946&tz=Asia/Kolkata"

# Expected: >100 requests/second
# Expected: <200ms average response time
```

### 4. Integration Testing

```bash
# Run Go integration tests
go test ./... -v -run Integration

# Run frontend tests
cd ui
npm test
```

## Troubleshooting

### Common Issues

#### 1. gRPC Server Not Starting
```bash
# Check if port is in use
lsof -i :50051

# Check server logs
tail -f ./tmp/grpc-server.log

# Verify Go modules
go mod tidy
go mod verify
```

#### 2. HTTP Gateway Connection Failed
```bash
# Verify gRPC server is running
grpcurl -plaintext localhost:50051 list

# Check gateway logs
tail -f ./tmp/http-gateway.log

# Test direct gRPC connection
grpcurl -plaintext localhost:50051 panchangam.Panchangam/Get
```

#### 3. CORS Errors in Frontend
The gateway is configured with CORS for:
- http://localhost:5173 (Vite dev server)
- http://localhost:3000 (React dev server)

To add more origins, edit `gateway/server.go`:
```go
AllowedOrigins: []string{
    "http://localhost:5173",
    "http://localhost:3000",
    "https://your-domain.com", // Add your domain
},
```

#### 4. Performance Issues

**Check service performance:**
```bash
# Test individual components
go test ./services/panchangam/... -bench=. -benchmem

# Monitor resource usage
htop
# or
docker stats  # if using Docker
```

**Common fixes:**
- Increase ephemeris cache size
- Enable HTTP caching headers
- Use CDN for frontend assets
- Scale horizontally with load balancer

### Log Locations

- **Development**: `./tmp/*.log`
- **Systemd**: `journalctl -u panchangam-grpc -f`
- **Docker**: `docker logs <container_id>`
- **Kubernetes**: `kubectl logs <pod-name>`

### Debug Mode

Enable debug logging:
```bash
# Backend
./bin/panchangam-server --log-level=debug
./bin/panchangam-gateway --log-level=debug

# Frontend
VITE_LOG_LEVEL=debug npm run dev
```

## Monitoring & Observability

The application includes OpenTelemetry integration for:
- **Distributed Tracing**: Track requests across services
- **Metrics**: Performance, error rates, latency
- **Logs**: Structured logging with context

### Exporting Telemetry

Configure OTEL exporters in environment:
```bash
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
export OTEL_SERVICE_NAME=panchangam
export OTEL_TRACES_EXPORTER=otlp
export OTEL_METRICS_EXPORTER=otlp
```

### Recommended Tools
- **Jaeger**: Distributed tracing
- **Prometheus**: Metrics collection
- **Grafana**: Visualization dashboards

## Production Checklist

Before deploying to production:

- [ ] All tests passing (`go test ./...`)
- [ ] Performance benchmarks meet targets
- [ ] SSL/TLS certificates configured
- [ ] Firewall rules configured
- [ ] Monitoring and alerting set up
- [ ] Backup strategy in place
- [ ] Log rotation configured
- [ ] Rate limiting configured
- [ ] Load testing completed
- [ ] Security audit performed
- [ ] Documentation updated
- [ ] Rollback procedure tested

## Support

For issues and questions:
- GitHub Issues: https://github.com/naren-m/panchangam/issues
- Documentation: [README.md](README.md)
- Feature Coverage: [FEATURES.md](FEATURES.md)

## License

See [LICENSE](LICENSE) file for details.
