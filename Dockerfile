# =============================================================================
# Panchangam Multi-Stage Dockerfile
# Builds: Go Backend (gRPC + Gateway) + React Frontend (Nginx)
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: Build Go Backend
# -----------------------------------------------------------------------------
FROM golang:1.23-alpine AS go-builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Declare build arguments for multi-arch support
ARG TARGETARCH

# Build the gRPC server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s" \
    -o /bin/grpc-server \
    ./cmd/grpc-server

# Build the gateway
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s" \
    -o /bin/gateway \
    ./cmd/gateway

# -----------------------------------------------------------------------------
# Stage 2: Build React Frontend
# -----------------------------------------------------------------------------
FROM node:18-alpine AS ui-builder

WORKDIR /app/ui

# Install dependencies first (for better caching)
COPY ui/package*.json ./
RUN npm ci --silent

# Copy source and build
COPY ui/ ./
RUN npm run build

# -----------------------------------------------------------------------------
# Stage 3: Production Image
# -----------------------------------------------------------------------------
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    nginx \
    curl \
    supervisor \
    && mkdir -p /var/log/supervisor /var/run/nginx /var/cache/nginx \
    && mkdir -p /var/lib/nginx/logs /var/lib/nginx/tmp/client_body \
    && mkdir -p /var/lib/nginx/tmp/proxy /var/lib/nginx/tmp/fastcgi \
    && mkdir -p /var/lib/nginx/tmp/uwsgi /var/lib/nginx/tmp/scgi \
    && chmod -R 755 /var/lib/nginx \
    && chown -R root:root /var/lib/nginx /var/cache/nginx /var/log/nginx /var/run/nginx

# Create app user
RUN addgroup -g 1001 -S panchangam && \
    adduser -u 1001 -S panchangam -G panchangam

WORKDIR /app

# Copy Go binaries
COPY --from=go-builder /bin/grpc-server /app/grpc-server
COPY --from=go-builder /bin/gateway /app/gateway

# Copy frontend build
COPY --from=ui-builder /app/ui/dist /usr/share/nginx/html

# Copy nginx configuration
COPY ui/nginx.conf /etc/nginx/nginx.conf

# Copy supervisor configuration
COPY <<EOF /etc/supervisor/conf.d/panchangam.conf
[supervisord]
nodaemon=true
user=root
logfile=/var/log/supervisor/supervisord.log
pidfile=/var/run/supervisord.pid

[program:grpc-server]
command=/app/grpc-server
directory=/app
autostart=true
autorestart=true
startsecs=5
startretries=3
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0
stderr_logfile=/dev/stderr
stderr_logfile_maxbytes=0

[program:gateway]
command=/app/gateway --grpc-endpoint=localhost:50052 --http-port=8080
directory=/app
autostart=true
autorestart=true
startsecs=5
startretries=3
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0
stderr_logfile=/dev/stderr
stderr_logfile_maxbytes=0

[program:nginx]
command=/usr/sbin/nginx -g "daemon off;"
autostart=true
autorestart=true
startsecs=5
startretries=3
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0
stderr_logfile=/dev/stderr
stderr_logfile_maxbytes=0
EOF

# Expose ports
# 80   - Nginx (frontend)
# 8080 - Gateway (REST API)
# 50052 - gRPC server
EXPOSE 80 8080 50052

# Health check (checks both nginx and gateway)
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:80/health 2>/dev/null || \
        curl -f http://localhost:8080/health 2>/dev/null || exit 1

# Start all services via supervisor
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/panchangam.conf"]
