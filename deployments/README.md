# Panchangam Production Deployment Guide

This directory contains all the infrastructure and deployment configurations for the Panchangam application.

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Prerequisites](#prerequisites)
4. [Quick Start](#quick-start)
5. [Deployment Methods](#deployment-methods)
6. [Configuration](#configuration)
7. [Monitoring](#monitoring)
8. [Backup & Recovery](#backup--recovery)
9. [Troubleshooting](#troubleshooting)

## Overview

Panchangam is deployed as a microservices architecture with the following components:

- **Frontend**: React SPA served by Nginx
- **Backend Gateway**: HTTP/REST API gateway
- **Backend gRPC**: Core gRPC service
- **Database**: PostgreSQL 16 with read replicas
- **Cache**: Redis 7
- **Monitoring**: Prometheus, Grafana, Jaeger, Loki, Alertmanager

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Load Balancer (Nginx)                   │
│                         SSL/TLS Termination                  │
└────────────────┬─────────────────────┬──────────────────────┘
                 │                     │
        ┌────────▼────────┐   ┌───────▼────────┐
        │    Frontend     │   │   API Gateway   │
        │  (React/Nginx)  │   │   (Go HTTP)     │
        └─────────────────┘   └────────┬────────┘
                                       │
                              ┌────────▼────────┐
                              │   gRPC Server   │
                              │   (Go gRPC)     │
                              └────────┬────────┘
                                       │
                    ┌──────────────────┼──────────────────┐
                    │                  │                  │
          ┌─────────▼────────┐  ┌─────▼──────┐  ┌────────▼────────┐
          │   PostgreSQL     │  │   Redis    │  │  OpenTelemetry  │
          │   Primary+Replica│  │   Cache    │  │     (Jaeger)    │
          └──────────────────┘  └────────────┘  └─────────────────┘
                    │
          ┌─────────▼────────┐
          │  Backup Service  │
          │  (Daily + S3)    │
          └──────────────────┘

          ┌──────────────────────────────────────┐
          │       Monitoring Stack               │
          │  Prometheus → Grafana → Alertmanager │
          │  Loki (Logs) → Grafana               │
          └──────────────────────────────────────┘
```

## Prerequisites

### For Kubernetes Deployment

- Kubernetes cluster (1.25+)
- kubectl configured
- kustomize (v4.5+)
- Cert-manager for SSL certificates
- Nginx Ingress Controller

### For Docker Compose Deployment

- Docker (20.10+)
- Docker Compose (v2.0+)
- At least 8GB RAM, 4 CPU cores
- 50GB disk space

### General Requirements

- Domain name with DNS configured
- SSL certificates (or Let's Encrypt)
- SMTP server for alerts (optional)
- S3 bucket for backups (optional)

## Quick Start

### Option 1: Docker Compose (Recommended for Development/Staging)

```bash
# 1. Clone the repository
git clone https://github.com/naren-m/panchangam.git
cd panchangam

# 2. Copy and configure environment variables
cp .env.production .env
# Edit .env and update passwords and secrets

# 3. Deploy the stack
docker-compose -f docker-compose.prod.yml up -d

# 4. Run database migrations
./deployments/migrations/migrate.sh up

# 5. Verify deployment
docker-compose -f docker-compose.prod.yml ps
curl http://localhost:8080/health
curl http://localhost:80
```

### Option 2: Kubernetes (Recommended for Production)

```bash
# 1. Clone the repository
git clone https://github.com/naren-m/panchangam.git
cd panchangam

# 2. Configure secrets
cd deployments/k8s/overlays/production
cp secrets.env.example secrets.env
# Edit secrets.env with production credentials

# 3. Deploy using kustomize
kustomize build . | kubectl apply -f -

# 4. Monitor rollout
kubectl rollout status deployment/panchangam-grpc -n panchangam
kubectl rollout status deployment/panchangam-gateway -n panchangam
kubectl rollout status deployment/panchangam-frontend -n panchangam

# 5. Verify deployment
kubectl get pods -n panchangam
kubectl get svc -n panchangam
kubectl get ingress -n panchangam
```

## Deployment Methods

### Automated Deployment with Script

```bash
# Deploy to staging with Docker Compose
./deployments/scripts/deploy.sh staging docker-compose

# Deploy to production with Kubernetes
./deployments/scripts/deploy.sh production kubernetes
```

### CI/CD Deployment

The application uses GitHub Actions for CI/CD:

- **Push to `develop` branch**: Automatically deploys to staging
- **Push to `main` branch**: Automatically deploys to production (with approval)
- **Manual trigger**: Use workflow_dispatch for ad-hoc deployments

## Configuration

### Environment Variables

Each environment has its own `.env` file:

- `.env.development` - Local development
- `.env.staging` - Staging environment
- `.env.production` - Production environment

Key configuration variables:

```bash
# Database
POSTGRES_DB=panchangam
POSTGRES_USER=panchangam
POSTGRES_PASSWORD=<strong-password>

# Redis
REDIS_PASSWORD=<strong-password>

# Application
LOG_LEVEL=info
ENVIRONMENT=production

# Monitoring
GRAFANA_USER=admin
GRAFANA_PASSWORD=<strong-password>

# Alerts
SLACK_WEBHOOK_URL=<webhook-url>
PAGERDUTY_SERVICE_KEY=<service-key>
```

### Multi-Environment Setup

The project supports three environments using Kustomize overlays:

1. **Development** (`deployments/k8s/overlays/dev`)
   - 1 replica per service
   - Debug logging enabled
   - Relaxed resource limits

2. **Staging** (`deployments/k8s/overlays/staging`)
   - 2 replicas per service
   - Production-like configuration
   - Moderate resource limits

3. **Production** (`deployments/k8s/overlays/production`)
   - 3+ replicas per service
   - Strict resource limits
   - Auto-scaling enabled
   - Production monitoring

## Monitoring

### Access Monitoring Dashboards

- **Grafana**: `https://monitoring.panchangam.app`
  - Username: admin
  - Password: (set in secrets)

- **Prometheus**: `https://monitoring.panchangam.app/prometheus`
- **Jaeger**: `https://monitoring.panchangam.app/jaeger`
- **Alertmanager**: `https://monitoring.panchangam.app/alertmanager`

### Key Metrics

1. **Service Health**
   - Uptime percentage
   - Request rate (RPS)
   - Error rate (%)
   - P95/P99 latency

2. **Resource Utilization**
   - CPU usage (%)
   - Memory usage (%)
   - Disk I/O
   - Network throughput

3. **Business Metrics**
   - API calls per minute
   - Cache hit rate
   - Database query performance
   - Active user sessions

### Alerts

Configured alerts:

- Service down (Critical)
- High error rate > 5% (Warning)
- High response time > 1s (Warning)
- High memory usage > 90% (Warning)
- High CPU usage > 90% (Warning)
- Database connection issues (Critical)
- Redis down (Critical)

## Backup & Recovery

### Automated Backups

Backups run daily at 2 AM UTC:

- Local retention: 30 days
- S3 retention: 90 days
- Includes database dumps and configuration

### Manual Backup

```bash
# Backup database
docker-compose -f docker-compose.prod.yml exec postgres \
    pg_dump -U panchangam panchangam | gzip > backup_$(date +%Y%m%d).sql.gz

# Or use the backup script
./deployments/backup/backup.sh
```

### Restore from Backup

```bash
# Restore latest backup
./deployments/backup/restore.sh latest

# Restore specific backup
./deployments/backup/restore.sh /backups/panchangam_backup_20250118_020000.sql.gz
```

### Disaster Recovery

Recovery Time Objective (RTO): 15 minutes
Recovery Point Objective (RPO): 1 hour

1. **Database Failure**
   ```bash
   # Switch to replica
   kubectl scale deployment postgres-replica --replicas=1
   kubectl patch service postgres-primary -p '{"spec":{"selector":{"role":"replica"}}}'
   ```

2. **Complete Infrastructure Failure**
   ```bash
   # Restore from S3 backup
   aws s3 cp s3://panchangam-backups/latest.sql.gz ./
   ./deployments/backup/restore.sh ./latest.sql.gz

   # Redeploy infrastructure
   ./deployments/scripts/deploy.sh production kubernetes
   ```

## Troubleshooting

### Common Issues

#### 1. Service Won't Start

```bash
# Check logs
kubectl logs -f deployment/panchangam-gateway -n panchangam
# Or for Docker Compose
docker-compose -f docker-compose.prod.yml logs -f gateway

# Check events
kubectl get events -n panchangam --sort-by='.lastTimestamp'
```

#### 2. Database Connection Issues

```bash
# Test database connection
kubectl exec -it deployment/panchangam-gateway -n panchangam -- \
    psql -h postgres-primary -U panchangam -d panchangam

# Check database logs
kubectl logs -f deployment/postgres-primary -n panchangam
```

#### 3. High Memory Usage

```bash
# Check resource usage
kubectl top pods -n panchangam

# Increase limits in deployment manifests
# Edit deployments/k8s/base/backend-deployment.yaml
```

#### 4. SSL Certificate Issues

```bash
# Check cert-manager logs
kubectl logs -f deployment/cert-manager -n cert-manager

# Check certificate status
kubectl describe certificate panchangam-tls -n panchangam

# Force renewal
kubectl delete certificate panchangam-tls -n panchangam
```

### Health Check Endpoints

- Frontend: `http://localhost:80/health`
- API Gateway: `http://localhost:8080/health`
- gRPC Server: Use grpcurl or health check command
- Prometheus: `http://localhost:9090/-/healthy`
- Grafana: `http://localhost:3000/api/health`

### Performance Tuning

1. **Database Optimization**
   - Adjust `shared_buffers` and `effective_cache_size`
   - Enable query logging for slow queries
   - Regular VACUUM and ANALYZE

2. **Cache Optimization**
   - Monitor cache hit rate
   - Adjust Redis `maxmemory` and eviction policy
   - Implement cache warming for frequently accessed data

3. **Application Tuning**
   - Adjust HPA thresholds based on traffic patterns
   - Optimize database queries
   - Enable connection pooling

### Rollback Procedure

```bash
# Kubernetes rollback
kubectl rollout undo deployment/panchangam-gateway -n panchangam
kubectl rollout undo deployment/panchangam-grpc -n panchangam
kubectl rollout undo deployment/panchangam-frontend -n panchangam

# Docker Compose rollback
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d --force-recreate
```

## Security Best Practices

1. **Secrets Management**
   - Use Kubernetes Secrets or external secret managers (Vault, AWS Secrets Manager)
   - Never commit secrets to version control
   - Rotate credentials regularly

2. **Network Security**
   - Enable HTTPS/TLS everywhere
   - Use network policies in Kubernetes
   - Implement rate limiting

3. **Access Control**
   - Follow principle of least privilege
   - Use RBAC in Kubernetes
   - Enable audit logging

4. **Regular Updates**
   - Keep base images updated
   - Apply security patches promptly
   - Run regular vulnerability scans

## Maintenance

### Regular Tasks

- **Daily**: Monitor alerts and metrics
- **Weekly**: Review logs for errors and anomalies
- **Monthly**: Review and optimize resource allocation
- **Quarterly**: Security audit and dependency updates

### Scaling

**Horizontal Scaling** (Automatic with HPA):
- Backend services: 3-10 replicas
- Frontend: 2-5 replicas

**Vertical Scaling** (Manual):
- Update resource limits in deployment manifests
- Apply changes: `kubectl apply -f deployments/k8s/base/`

## Support

For issues and questions:

- GitHub Issues: https://github.com/naren-m/panchangam/issues
- Documentation: https://github.com/naren-m/panchangam/docs
- Email: support@panchangam.app

## License

Copyright © 2025 Panchangam Project
