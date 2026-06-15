# Deployment Infrastructure Summary

This file summarizes the current deployment shape. It is intentionally short so
it stays useful when the deployment files change.

## Current Components

- `docker-compose.prod.yml`: production-like local stack
- `deployments/k8s/base`: shared Kubernetes manifests
- `deployments/k8s/overlays`: environment-specific Kustomize overlays
- `deployments/scripts/deploy.sh`: deploy helper for Compose and Kubernetes
- `deployments/prometheus`: metrics scrape and alert rules
- `deployments/grafana`: dashboards and data source provisioning
- `deployments/alertmanager`: alert routing
- `deployments/loki`: log storage config
- `deployments/nginx`: frontend and load-balancer config

## Runtime Services

- frontend
- API gateway
- gRPC server
- Redis
- Prometheus
- Grafana
- Jaeger
- Loki
- Alertmanager

## Simplification Notes

- The application runtime currently uses Redis for caching.
- The deployment stack does not start unused SQL storage services.
- Deploy commands only build, roll out, check health, and run smoke tests.
- Operational docs point to the gateway health endpoint:
  `http://localhost:8080/api/v1/health`.

## Verification

Use these local checks before publishing deployment changes:

```bash
python3 test/test_deployment_health_config.py
docker compose -f docker-compose.prod.yml config
bash -n deployments/scripts/deploy.sh
```

If `kustomize` is installed, also run:

```bash
kustomize build deployments/k8s/base
```
