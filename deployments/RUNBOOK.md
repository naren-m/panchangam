# Panchangam Operations Runbook

Use this runbook for common production operations. Keep commands simple and
prefer health checks before deeper debugging.

## Incident Response

1. Acknowledge the alert.
2. Check impact with the health endpoint and dashboards.
3. Mitigate the issue or roll back.
4. Record what happened and what changed.

## Deploy A New Version

```bash
kubectl get pods -n panchangam
./deployments/scripts/deploy.sh production kubernetes
kubectl rollout status deployment/panchangam-gateway -n panchangam
kubectl rollout status deployment/panchangam-grpc -n panchangam
kubectl rollout status deployment/panchangam-frontend -n panchangam
curl https://api.panchangam.app/api/v1/health
curl https://panchangam.app
```

Rollback:

```bash
kubectl rollout undo deployment/panchangam-gateway -n panchangam
kubectl rollout undo deployment/panchangam-grpc -n panchangam
kubectl rollout undo deployment/panchangam-frontend -n panchangam
```

## Scale Services

```bash
kubectl scale deployment/panchangam-grpc --replicas=5 -n panchangam
kubectl scale deployment/panchangam-gateway --replicas=5 -n panchangam
kubectl scale deployment/panchangam-frontend --replicas=3 -n panchangam
kubectl get pods -n panchangam -w
```

The HPA may change replica counts after manual scaling.

## Clear Redis Cache

Clear all Redis keys:

```bash
kubectl exec -it deployment/redis -n panchangam -- \
    sh -c 'REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli FLUSHALL'
```

Clear only Panchangam keys:

```bash
kubectl exec -it deployment/redis -n panchangam -- \
    sh -c 'export REDISCLI_AUTH="$REDIS_PASSWORD"; redis-cli --scan --pattern "panchangam:*" | xargs -r redis-cli DEL'
```

Check Redis stats:

```bash
kubectl exec -it deployment/redis -n panchangam -- \
    sh -c 'REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli INFO stats'
```

## View Logs

```bash
kubectl logs -f deployment/panchangam-gateway -n panchangam
kubectl logs -f deployment/panchangam-grpc -n panchangam
kubectl logs -f deployment/panchangam-frontend -n panchangam
```

For Docker Compose:

```bash
docker-compose -f docker-compose.prod.yml logs -f gateway
docker-compose -f docker-compose.prod.yml logs -f grpc-server
docker-compose -f docker-compose.prod.yml logs -f frontend
```

## High Response Time

Check the simple causes first:

```bash
kubectl top pods -n panchangam
kubectl logs -f deployment/panchangam-gateway -n panchangam
kubectl exec -it deployment/redis -n panchangam -- \
    sh -c 'REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli INFO stats'
```

If resources are saturated, scale the busy service:

```bash
kubectl scale deployment/panchangam-gateway --replicas=5 -n panchangam
```

## Service Will Not Start

```bash
kubectl describe pod -l app=panchangam-gateway -n panchangam
kubectl logs -f deployment/panchangam-gateway -n panchangam
kubectl get events -n panchangam --sort-by='.lastTimestamp'
curl http://localhost:8080/api/v1/health
```

Common checks:

- image tag exists
- required secrets exist
- Redis is ready
- CPU and memory limits are not too low
- probes use `/api/v1/health`

## Certificate Rotation

```bash
kubectl describe certificate panchangam-tls -n panchangam
kubectl delete certificate panchangam-tls -n panchangam
kubectl get certificate -n panchangam
```

## Maintenance Checklist

Daily:

- check alerts
- review Grafana dashboards
- verify `https://api.panchangam.app/api/v1/health`

Weekly:

- review application errors
- check resource trends
- update this runbook when commands change

Monthly:

- review alert thresholds
- review dependency updates
- run a small rollback drill in staging
