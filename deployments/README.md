# Panchangam Deployment Guide

This directory contains the files used to run Panchangam outside a local
developer shell. Keep this guide current with the files in this directory.

## What Runs

- Frontend: React app served by Nginx
- API Gateway: HTTP API on port 8080
- gRPC Server: Panchangam calculation service on port 50052
- Redis: cache for repeated calculations
- Observability: Prometheus, Grafana, Jaeger, Loki, and Alertmanager

## Docker Compose

Use Docker Compose for local production-like testing.

```bash
docker-compose -f docker-compose.prod.yml up -d
docker-compose -f docker-compose.prod.yml ps
curl http://localhost:8080/api/v1/health
curl http://localhost:80
```

Stop the stack with:

```bash
docker-compose -f docker-compose.prod.yml down
```

## Kubernetes

Use Kustomize for cluster deployments.

```bash
kustomize build deployments/k8s/overlays/production | kubectl apply -f -
kubectl rollout status deployment/panchangam-grpc -n panchangam
kubectl rollout status deployment/panchangam-gateway -n panchangam
kubectl rollout status deployment/panchangam-frontend -n panchangam
kubectl get pods -n panchangam
```

## Configuration

Set only the values the runtime actually reads.

```bash
REDIS_PASSWORD=<strong-password>
LOG_LEVEL=info
ENVIRONMENT=production
GRAFANA_PASSWORD=<strong-password>
SLACK_WEBHOOK_URL=<optional-webhook-url>
PAGERDUTY_SERVICE_KEY=<optional-service-key>
ALERT_EMAIL=<optional-alert-email>
```

## Monitoring

Key endpoints:

- API health: `http://localhost:8080/api/v1/health`
- Frontend: `http://localhost:80`
- Prometheus: `http://localhost:9090/-/healthy`
- Grafana: `http://localhost:3000/api/health`

Watch these signals during deploys:

- service availability
- request error rate
- p95 and p99 latency
- CPU and memory usage
- Redis availability and cache hit rate

## Troubleshooting

Check service logs:

```bash
kubectl logs -f deployment/panchangam-gateway -n panchangam
kubectl logs -f deployment/panchangam-grpc -n panchangam
docker-compose -f docker-compose.prod.yml logs -f gateway
```

Check Kubernetes events:

```bash
kubectl get events -n panchangam --sort-by='.lastTimestamp'
```

Check Redis from the Redis pod:

```bash
kubectl exec -it deployment/redis -n panchangam -- \
    sh -c 'REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli INFO stats'
```

## Maintenance

- Daily: check alerts and health endpoints
- Weekly: review logs and high-latency requests
- Monthly: review resource limits and dependency updates
- Quarterly: run a disaster recovery drill for the deployment platform
