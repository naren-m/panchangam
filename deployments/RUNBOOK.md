# Panchangam Operations Runbook

This runbook provides step-by-step procedures for common operational tasks and incident response.

## Table of Contents

1. [Emergency Contacts](#emergency-contacts)
2. [Incident Response](#incident-response)
3. [Common Operations](#common-operations)
4. [Troubleshooting Playbooks](#troubleshooting-playbooks)
5. [Maintenance Procedures](#maintenance-procedures)

## Emergency Contacts

### On-Call Rotation
- Primary: [Contact Info]
- Secondary: [Contact Info]
- Manager: [Contact Info]

### External Services
- Cloud Provider Support: [Contact Info]
- Database Support: [Contact Info]
- Security Team: [Contact Info]

## Incident Response

### Severity Levels

**P0 - Critical (Response Time: < 15 minutes)**
- Complete service outage
- Data loss or corruption
- Security breach

**P1 - High (Response Time: < 1 hour)**
- Partial service degradation affecting > 50% of users
- Performance issues causing timeouts
- Monitoring system failure

**P2 - Medium (Response Time: < 4 hours)**
- Minor service degradation
- Non-critical feature failure
- Elevated error rates

**P3 - Low (Response Time: < 24 hours)**
- Cosmetic issues
- Minor bugs
- Documentation errors

### Incident Response Process

1. **Acknowledge**: Respond to alert within SLA
2. **Assess**: Determine severity and impact
3. **Communicate**: Update status page and stakeholders
4. **Mitigate**: Implement immediate fixes or workarounds
5. **Resolve**: Apply permanent fix
6. **Document**: Write post-mortem

## Common Operations

### 1. Deploy New Version

**Prerequisites:**
- Code reviewed and approved
- Tests passing in CI/CD
- Staging deployment successful

**Procedure:**
```bash
# 1. Verify current state
kubectl get pods -n panchangam

# 2. Deploy new version
./deployments/scripts/deploy.sh production kubernetes

# 3. Monitor rollout
kubectl rollout status deployment/panchangam-gateway -n panchangam
kubectl rollout status deployment/panchangam-grpc -n panchangam

# 4. Run smoke tests
curl https://api.panchangam.app/health
curl https://panchangam.app

# 5. Monitor metrics for 30 minutes
# Check Grafana dashboard for errors, latency, resource usage
```

**Rollback Procedure:**
```bash
kubectl rollout undo deployment/panchangam-gateway -n panchangam
kubectl rollout undo deployment/panchangam-grpc -n panchangam
kubectl rollout undo deployment/panchangam-frontend -n panchangam
```

### 2. Scale Services

**Horizontal Scaling (Production):**
```bash
# Scale backend services
kubectl scale deployment/panchangam-grpc --replicas=5 -n panchangam
kubectl scale deployment/panchangam-gateway --replicas=5 -n panchangam

# Scale frontend
kubectl scale deployment/panchangam-frontend --replicas=3 -n panchangam

# Verify scaling
kubectl get pods -n panchangam -w
```

**Note:** HPA will automatically scale within configured limits. Manual scaling overrides HPA temporarily.

### 3. Database Operations

**Run Migrations:**
```bash
# Connect to migration pod
kubectl run migration \
    --namespace=panchangam \
    --image=migrate/migrate:latest \
    --rm -i --restart=Never \
    --command -- migrate \
    -path /migrations \
    -database "postgres://user:pass@postgres-primary:5432/panchangam?sslmode=disable" \
    up

# Or use the script
./deployments/migrations/migrate.sh up
```

**Create Database Backup:**
```bash
# Manual backup
kubectl exec -it deployment/postgres-primary -n panchangam -- \
    pg_dump -U panchangam panchangam | gzip > backup_$(date +%Y%m%d_%H%M%S).sql.gz

# Or use automated backup script
kubectl exec -it deployment/backup -n panchangam -- /backup.sh
```

**Restore Database:**
```bash
# Copy backup to pod
kubectl cp backup.sql.gz panchangam/postgres-primary-xxx:/tmp/

# Restore
kubectl exec -it deployment/postgres-primary -n panchangam -- \
    gunzip -c /tmp/backup.sql.gz | psql -U panchangam panchangam
```

### 4. Certificate Renewal

**Check Certificate Expiry:**
```bash
kubectl describe certificate panchangam-tls -n panchangam
```

**Force Renewal:**
```bash
# Delete existing certificate
kubectl delete certificate panchangam-tls -n panchangam

# Cert-manager will automatically create new one
# Wait 2-3 minutes and verify
kubectl get certificate -n panchangam
```

### 5. Clear Cache

**Redis Cache:**
```bash
# Connect to Redis
kubectl exec -it deployment/redis -n panchangam -- redis-cli

# Clear all keys
FLUSHALL

# Clear specific pattern
KEYS panchangam:cache:*
DEL panchangam:cache:*

# Exit
EXIT
```

### 6. View Logs

**Application Logs:**
```bash
# Gateway logs
kubectl logs -f deployment/panchangam-gateway -n panchangam

# gRPC server logs
kubectl logs -f deployment/panchangam-grpc -n panchangam

# All containers in a pod
kubectl logs -f pod/panchangam-gateway-xxx -n panchangam --all-containers

# Previous logs (after crash)
kubectl logs deployment/panchangam-gateway -n panchangam --previous
```

**Aggregate Logs (Loki):**
```bash
# Access Grafana Explore
# URL: https://monitoring.panchangam.app/explore
# Query: {namespace="panchangam", app="panchangam"}
```

## Troubleshooting Playbooks

### Playbook 1: Service is Down

**Symptoms:**
- Health check failures
- 502/503 errors
- Alert: "ServiceDown"

**Investigation:**
```bash
# 1. Check pod status
kubectl get pods -n panchangam

# 2. Describe failing pod
kubectl describe pod <pod-name> -n panchangam

# 3. Check recent events
kubectl get events -n panchangam --sort-by='.lastTimestamp' | tail -20

# 4. Check logs
kubectl logs <pod-name> -n panchangam --tail=100
```

**Common Causes & Solutions:**

**Cause: Out of Memory (OOMKilled)**
```bash
# Increase memory limits
# Edit deployment and update resources.limits.memory
kubectl edit deployment/panchangam-gateway -n panchangam

# Or apply updated manifest
kubectl apply -f deployments/k8s/base/backend-deployment.yaml
```

**Cause: Image Pull Failure**
```bash
# Check image exists
docker pull ghcr.io/naren-m/panchangam-backend:latest

# Verify registry credentials
kubectl get secret -n panchangam | grep docker

# Update image pull secret if needed
kubectl create secret docker-registry ghcr-secret \
    --docker-server=ghcr.io \
    --docker-username=<username> \
    --docker-password=<token> \
    -n panchangam
```

**Cause: Failed Health Checks**
```bash
# Test health check manually
kubectl exec -it deployment/panchangam-gateway -n panchangam -- \
    curl http://localhost:8080/health

# Increase probe timeouts if needed
# Edit deployment and update livenessProbe/readinessProbe
```

### Playbook 2: High Response Time

**Symptoms:**
- P95 latency > 1s
- Alert: "HighResponseTime"
- User complaints of slow performance

**Investigation:**
```bash
# 1. Check current metrics
# Grafana: Response Time dashboard

# 2. Check resource utilization
kubectl top pods -n panchangam

# 3. Check database performance
kubectl exec -it deployment/postgres-primary -n panchangam -- \
    psql -U panchangam -c "SELECT * FROM pg_stat_activity WHERE state = 'active';"

# 4. Check cache hit rate
kubectl exec -it deployment/redis -n panchangam -- \
    redis-cli INFO stats | grep hit_rate
```

**Common Causes & Solutions:**

**Cause: Low Cache Hit Rate**
```bash
# Warm up cache for popular queries
# Implement cache pre-loading in application

# Increase cache size
# Edit Redis deployment and increase maxmemory
kubectl edit deployment/redis -n panchangam
```

**Cause: Slow Database Queries**
```bash
# Enable slow query logging
kubectl exec -it deployment/postgres-primary -n panchangam -- \
    psql -U panchangam -c "ALTER DATABASE panchangam SET log_min_duration_statement = 1000;"

# Check for missing indexes
kubectl exec -it deployment/postgres-primary -n panchangam -- \
    psql -U panchangam -c "SELECT schemaname, tablename, indexname FROM pg_indexes WHERE schemaname = 'public';"

# Run ANALYZE
kubectl exec -it deployment/postgres-primary -n panchangam -- \
    psql -U panchangam -c "ANALYZE;"
```

**Cause: Insufficient Resources**
```bash
# Scale horizontally
kubectl scale deployment/panchangam-gateway --replicas=5 -n panchangam

# Or increase resource limits
kubectl edit deployment/panchangam-gateway -n panchangam
```

### Playbook 3: Database Connection Issues

**Symptoms:**
- "Too many connections" errors
- Database timeout errors
- Alert: "HighDatabaseConnections"

**Investigation:**
```bash
# 1. Check active connections
kubectl exec -it deployment/postgres-primary -n panchangam -- \
    psql -U panchangam -c "SELECT count(*) FROM pg_stat_activity;"

# 2. Check connection by application
kubectl exec -it deployment/postgres-primary -n panchangam -- \
    psql -U panchangam -c "SELECT application_name, count(*) FROM pg_stat_activity GROUP BY application_name;"

# 3. Check max connections
kubectl exec -it deployment/postgres-primary -n panchangam -- \
    psql -U panchangam -c "SHOW max_connections;"
```

**Solutions:**
```bash
# Increase max_connections
kubectl exec -it deployment/postgres-primary -n panchangam -- \
    psql -U panchangam -c "ALTER SYSTEM SET max_connections = 200;"

# Restart PostgreSQL to apply
kubectl rollout restart deployment/postgres-primary -n panchangam

# Implement connection pooling in application
# Use PgBouncer or similar connection pooler
```

### Playbook 4: Disk Space Full

**Symptoms:**
- "No space left on device" errors
- Database write failures
- Log rotation issues

**Investigation:**
```bash
# Check PVC usage
kubectl get pvc -n panchangam

# Check disk usage in pod
kubectl exec -it deployment/postgres-primary -n panchangam -- df -h

# Check largest files/directories
kubectl exec -it deployment/postgres-primary -n panchangam -- \
    du -sh /var/lib/postgresql/data/* | sort -h
```

**Solutions:**
```bash
# Clean up old backups
kubectl exec -it deployment/backup -n panchangam -- \
    find /backups -type f -mtime +30 -delete

# Increase PVC size
kubectl edit pvc postgres-pvc -n panchangam
# Update storage size and apply

# Clean up old logs (if applicable)
kubectl exec -it deployment/postgres-primary -n panchangam -- \
    find /var/log -name "*.log" -mtime +7 -delete
```

## Maintenance Procedures

### Planned Downtime Procedure

**Before Maintenance:**
```bash
# 1. Schedule maintenance window
# 2. Notify stakeholders via status page
# 3. Create backup
./deployments/backup/backup.sh

# 4. Scale down to maintenance mode (optional)
kubectl scale deployment/panchangam-gateway --replicas=0 -n panchangam
```

**During Maintenance:**
```bash
# Perform maintenance tasks
# e.g., database upgrades, schema changes, etc.
```

**After Maintenance:**
```bash
# 1. Restore services
kubectl scale deployment/panchangam-gateway --replicas=3 -n panchangam

# 2. Verify health
kubectl get pods -n panchangam
curl https://api.panchangam.app/health

# 3. Monitor metrics for 1 hour
# 4. Update status page
```

### Database Upgrade Procedure

```bash
# 1. Create backup
./deployments/backup/backup.sh

# 2. Test upgrade on staging first
# Deploy to staging and verify

# 3. Schedule production maintenance window

# 4. Update PostgreSQL version in deployment
# Edit deployments/k8s/base/postgres-deployment.yaml
# Change image: postgres:16-alpine to postgres:17-alpine

# 5. Apply changes
kubectl apply -f deployments/k8s/base/postgres-deployment.yaml

# 6. Monitor rollout
kubectl rollout status deployment/postgres-primary -n panchangam

# 7. Verify database version
kubectl exec -it deployment/postgres-primary -n panchangam -- \
    psql -U panchangam -c "SELECT version();"

# 8. Run smoke tests
```

### Certificate Rotation

```bash
# 1. Backup current certificates
kubectl get secret panchangam-tls -n panchangam -o yaml > tls-backup.yaml

# 2. Update certificate
# If using cert-manager, delete and it will auto-renew
kubectl delete certificate panchangam-tls -n panchangam

# 3. Wait for new certificate
kubectl get certificate -n panchangam -w

# 4. Verify new certificate
kubectl describe certificate panchangam-tls -n panchangam
```

### Scaling for Traffic Spikes

**Proactive Scaling (before expected spike):**
```bash
# 1. Increase replicas
kubectl scale deployment/panchangam-grpc --replicas=10 -n panchangam
kubectl scale deployment/panchangam-gateway --replicas=10 -n panchangam

# 2. Update HPA max replicas temporarily
kubectl edit hpa panchangam-grpc-hpa -n panchangam
# Change maxReplicas to higher value

# 3. Monitor performance
# Check Grafana dashboards

# 4. After spike, scale back down
kubectl scale deployment/panchangam-grpc --replicas=3 -n panchangam
kubectl scale deployment/panchangam-gateway --replicas=3 -n panchangam
```

## Post-Incident Review Template

After resolving an incident, complete this template:

```markdown
# Post-Incident Review: [Incident Title]

**Date:** YYYY-MM-DD
**Severity:** P0/P1/P2/P3
**Duration:** X hours Y minutes
**Impact:** Number of affected users, revenue impact, etc.

## Timeline
- HH:MM - Incident detected
- HH:MM - Team notified
- HH:MM - Root cause identified
- HH:MM - Mitigation applied
- HH:MM - Service restored
- HH:MM - Incident closed

## Root Cause
[Detailed explanation of what caused the incident]

## Impact
[Detailed impact on users, business, etc.]

## Resolution
[How the incident was resolved]

## Action Items
1. [ ] Action item 1 (Owner: Name, Due: Date)
2. [ ] Action item 2 (Owner: Name, Due: Date)

## Lessons Learned
- What went well
- What could be improved
- Preventive measures
```

## Monitoring Checklist

**Daily:**
- [ ] Check alert notifications
- [ ] Review Grafana dashboards
- [ ] Check error rates
- [ ] Verify backup completion

**Weekly:**
- [ ] Review slow query logs
- [ ] Check disk usage trends
- [ ] Review security alerts
- [ ] Update runbook if needed

**Monthly:**
- [ ] Review incident reports
- [ ] Optimize database indexes
- [ ] Review and update alerts
- [ ] Capacity planning review

**Quarterly:**
- [ ] Disaster recovery drill
- [ ] Security audit
- [ ] Dependency updates
- [ ] Performance optimization review
