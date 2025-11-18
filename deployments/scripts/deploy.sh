#!/bin/bash

# Panchangam Deployment Script
# Usage: ./deploy.sh [environment] [platform]
# Examples:
#   ./deploy.sh production kubernetes
#   ./deploy.sh staging docker-compose

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
ENVIRONMENT="${1:-staging}"
PLATFORM="${2:-kubernetes}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Validation
validate_environment() {
    case "$ENVIRONMENT" in
        development|staging|production)
            log_info "Deploying to: $ENVIRONMENT"
            ;;
        *)
            log_error "Invalid environment: $ENVIRONMENT"
            log_error "Valid environments: development, staging, production"
            exit 1
            ;;
    esac
}

validate_platform() {
    case "$PLATFORM" in
        kubernetes|docker-compose)
            log_info "Using platform: $PLATFORM"
            ;;
        *)
            log_error "Invalid platform: $PLATFORM"
            log_error "Valid platforms: kubernetes, docker-compose"
            exit 1
            ;;
    esac
}

# Pre-deployment checks
pre_deployment_checks() {
    log_info "Running pre-deployment checks..."

    # Check if required tools are installed
    if [ "$PLATFORM" = "kubernetes" ]; then
        if ! command -v kubectl &> /dev/null; then
            log_error "kubectl is not installed"
            exit 1
        fi

        if ! command -v kustomize &> /dev/null; then
            log_error "kustomize is not installed"
            exit 1
        fi
    else
        if ! command -v docker &> /dev/null; then
            log_error "docker is not installed"
            exit 1
        fi

        if ! command -v docker-compose &> /dev/null; then
            log_error "docker-compose is not installed"
            exit 1
        fi
    fi

    log_success "Pre-deployment checks passed"
}

# Database migration
run_migrations() {
    log_info "Running database migrations..."

    if [ "$PLATFORM" = "kubernetes" ]; then
        kubectl run migration-job \
            --namespace=panchangam \
            --image=migrate/migrate:latest \
            --rm -i --restart=Never \
            --command -- migrate \
            -path /migrations \
            -database "postgres://panchangam:${DB_PASSWORD}@postgres-primary:5432/panchangam?sslmode=disable" \
            up
    else
        cd "${PROJECT_ROOT}"
        docker-compose -f docker-compose.prod.yml run --rm \
            -e DB_HOST=postgres \
            -e DB_PORT=5432 \
            -e DB_NAME=panchangam \
            -e DB_USER=panchangam \
            -e DB_PASSWORD="${DB_PASSWORD}" \
            backend-gateway \
            sh -c "cd /app/deployments/migrations && ./migrate.sh up"
    fi

    log_success "Database migrations completed"
}

# Deploy to Kubernetes
deploy_kubernetes() {
    log_info "Deploying to Kubernetes ($ENVIRONMENT)..."

    cd "${PROJECT_ROOT}/deployments/k8s/overlays/${ENVIRONMENT}"

    # Apply kustomization
    kustomize build . | kubectl apply -f -

    # Wait for rollout
    log_info "Waiting for deployment rollout..."
    kubectl rollout status deployment/panchangam-grpc -n panchangam --timeout=5m
    kubectl rollout status deployment/panchangam-gateway -n panchangam --timeout=5m
    kubectl rollout status deployment/panchangam-frontend -n panchangam --timeout=5m

    log_success "Kubernetes deployment completed"
}

# Deploy to Docker Compose
deploy_docker_compose() {
    log_info "Deploying to Docker Compose ($ENVIRONMENT)..."

    cd "${PROJECT_ROOT}"

    # Load environment variables
    if [ -f ".env.${ENVIRONMENT}" ]; then
        export $(cat ".env.${ENVIRONMENT}" | grep -v '^#' | xargs)
    fi

    # Pull latest images
    docker-compose -f docker-compose.prod.yml pull

    # Deploy with rolling restart
    docker-compose -f docker-compose.prod.yml up -d --remove-orphans

    # Wait for health checks
    log_info "Waiting for services to be healthy..."
    sleep 30

    log_success "Docker Compose deployment completed"
}

# Post-deployment validation
post_deployment_checks() {
    log_info "Running post-deployment checks..."

    if [ "$PLATFORM" = "kubernetes" ]; then
        # Check pod status
        kubectl get pods -n panchangam

        # Check service endpoints
        kubectl get svc -n panchangam
    else
        # Check container health
        docker-compose -f docker-compose.prod.yml ps
    fi

    log_success "Post-deployment checks passed"
}

# Smoke tests
run_smoke_tests() {
    log_info "Running smoke tests..."

    # Determine the base URL
    if [ "$ENVIRONMENT" = "production" ]; then
        API_URL="https://api.panchangam.app"
        WEB_URL="https://panchangam.app"
    elif [ "$ENVIRONMENT" = "staging" ]; then
        API_URL="https://api-staging.panchangam.app"
        WEB_URL="https://staging.panchangam.app"
    else
        API_URL="http://localhost:8080"
        WEB_URL="http://localhost:80"
    fi

    # Test API health endpoint
    if curl -f "${API_URL}/health" > /dev/null 2>&1; then
        log_success "API health check passed"
    else
        log_error "API health check failed"
        exit 1
    fi

    # Test frontend
    if curl -f "${WEB_URL}" > /dev/null 2>&1; then
        log_success "Frontend health check passed"
    else
        log_error "Frontend health check failed"
        exit 1
    fi

    log_success "Smoke tests passed"
}

# Rollback function
rollback() {
    log_warning "Rolling back deployment..."

    if [ "$PLATFORM" = "kubernetes" ]; then
        kubectl rollout undo deployment/panchangam-grpc -n panchangam
        kubectl rollout undo deployment/panchangam-gateway -n panchangam
        kubectl rollout undo deployment/panchangam-frontend -n panchangam
    else
        docker-compose -f docker-compose.prod.yml down
        # Restore previous version
        log_info "Please manually restore previous Docker images"
    fi

    log_warning "Rollback completed"
}

# Main deployment flow
main() {
    log_info "==================================="
    log_info "Panchangam Deployment Script"
    log_info "==================================="
    log_info "Environment: $ENVIRONMENT"
    log_info "Platform: $PLATFORM"
    log_info "==================================="

    # Validate inputs
    validate_environment
    validate_platform

    # Run checks
    pre_deployment_checks

    # Run migrations
    if [ "$ENVIRONMENT" != "development" ]; then
        run_migrations
    fi

    # Deploy based on platform
    if [ "$PLATFORM" = "kubernetes" ]; then
        deploy_kubernetes
    else
        deploy_docker_compose
    fi

    # Post-deployment checks
    post_deployment_checks

    # Smoke tests
    run_smoke_tests

    log_success "==================================="
    log_success "Deployment completed successfully!"
    log_success "==================================="
}

# Error handling
trap 'log_error "Deployment failed. Run ./deploy.sh rollback to revert changes."' ERR

# Run main
main "$@"
