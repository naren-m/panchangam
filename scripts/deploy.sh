#!/bin/bash

# Deployment script for Panchangam application
set -euo pipefail

# Default values
ENVIRONMENT="staging"
VERSION=""
REGISTRY="ghcr.io"
NAMESPACE="panchangam"
FORCE_DEPLOY="false"
DRY_RUN="false"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Usage function
usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Deploy Panchangam application to specified environment.

OPTIONS:
    -e, --environment ENV    Deployment environment (staging/production) [default: staging]
    -v, --version VERSION    Version tag to deploy [required]
    -r, --registry REGISTRY  Container registry [default: ghcr.io]
    -n, --namespace NS       Kubernetes namespace [default: panchangam]
    -f, --force             Force deployment without confirmation
    -d, --dry-run           Show what would be deployed without executing
    -h, --help              Show this help message

EXAMPLES:
    $0 -e staging -v v1.2.3
    $0 -e production -v v1.2.3 --force
    $0 --dry-run -e staging -v latest

EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -e|--environment)
            ENVIRONMENT="$2"
            shift 2
            ;;
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -r|--registry)
            REGISTRY="$2"
            shift 2
            ;;
        -n|--namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -f|--force)
            FORCE_DEPLOY="true"
            shift
            ;;
        -d|--dry-run)
            DRY_RUN="true"
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            usage
            exit 1
            ;;
    esac
done

# Validation
if [ -z "$VERSION" ]; then
    echo -e "${RED}Error: Version is required${NC}"
    usage
    exit 1
fi

if [[ ! "$ENVIRONMENT" =~ ^(staging|production)$ ]]; then
    echo -e "${RED}Error: Environment must be 'staging' or 'production'${NC}"
    exit 1
fi

# Helper functions
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

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    local missing_tools=()
    
    if ! command -v docker &> /dev/null; then
        missing_tools+=("docker")
    fi
    
    if ! command -v kubectl &> /dev/null; then
        missing_tools+=("kubectl")
    fi
    
    if ! command -v helm &> /dev/null; then
        missing_tools+=("helm")
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        log_error "Missing required tools: ${missing_tools[*]}"
        exit 1
    fi
    
    log_success "All prerequisites satisfied"
}

# Validate images exist
validate_images() {
    log_info "Validating container images..."
    
    local backend_image="${REGISTRY}/panchangam-backend:${VERSION}"
    local frontend_image="${REGISTRY}/panchangam-frontend:${VERSION}"
    
    if [ "$DRY_RUN" = "false" ]; then
        if ! docker manifest inspect "$backend_image" &> /dev/null; then
            log_error "Backend image not found: $backend_image"
            exit 1
        fi
        
        if ! docker manifest inspect "$frontend_image" &> /dev/null; then
            log_error "Frontend image not found: $frontend_image"
            exit 1
        fi
    fi
    
    log_success "Container images validated"
}

# Generate Kubernetes manifests
generate_manifests() {
    log_info "Generating Kubernetes manifests..."
    
    local backend_image="${REGISTRY}/panchangam-backend:${VERSION}"
    local frontend_image="${REGISTRY}/panchangam-frontend:${VERSION}"
    
    # Create temporary directory for manifests
    local manifest_dir="/tmp/panchangam-deploy-${VERSION}"
    mkdir -p "$manifest_dir"
    
    # Generate backend deployment
    cat > "$manifest_dir/backend-deployment.yaml" << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: panchangam-backend
  namespace: ${NAMESPACE}
  labels:
    app: panchangam-backend
    version: ${VERSION}
    environment: ${ENVIRONMENT}
spec:
  replicas: $( [ "$ENVIRONMENT" = "production" ] && echo "3" || echo "2" )
  selector:
    matchLabels:
      app: panchangam-backend
  template:
    metadata:
      labels:
        app: panchangam-backend
        version: ${VERSION}
    spec:
      containers:
      - name: gateway
        image: ${backend_image}
        command: ["/panchangam-gateway"]
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: HTTP_PORT
          value: "8080"
        - name: GRPC_ENDPOINT
          value: "localhost:50052"
        - name: LOG_LEVEL
          value: "$( [ "$ENVIRONMENT" = "production" ] && echo "warn" || echo "info" )"
        livenessProbe:
          httpGet:
            path: /api/v1/health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /api/v1/health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "128Mi"
            cpu: "200m"
      - name: grpc-server
        image: ${backend_image}
        command: ["/panchangam-grpc"]
        ports:
        - containerPort: 50052
          name: grpc
        env:
        - name: PORT
          value: "50052"
        - name: LOG_LEVEL
          value: "$( [ "$ENVIRONMENT" = "production" ] && echo "warn" || echo "info" )"
        livenessProbe:
          exec:
            command: ["/panchangam-grpc", "--health-check"]
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          exec:
            command: ["/panchangam-grpc", "--health-check"]
          initialDelaySeconds: 5
          periodSeconds: 10
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "128Mi"
            cpu: "200m"
---
apiVersion: v1
kind: Service
metadata:
  name: panchangam-backend
  namespace: ${NAMESPACE}
  labels:
    app: panchangam-backend
spec:
  selector:
    app: panchangam-backend
  ports:
  - port: 80
    targetPort: 8080
    name: http
  - port: 50052
    targetPort: 50052
    name: grpc
EOF

    # Generate frontend deployment
    cat > "$manifest_dir/frontend-deployment.yaml" << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: panchangam-frontend
  namespace: ${NAMESPACE}
  labels:
    app: panchangam-frontend
    version: ${VERSION}
    environment: ${ENVIRONMENT}
spec:
  replicas: $( [ "$ENVIRONMENT" = "production" ] && echo "3" || echo "2" )
  selector:
    matchLabels:
      app: panchangam-frontend
  template:
    metadata:
      labels:
        app: panchangam-frontend
        version: ${VERSION}
    spec:
      containers:
      - name: frontend
        image: ${frontend_image}
        ports:
        - containerPort: 80
          name: http
        env:
        - name: API_ENDPOINT
          value: "http://panchangam-backend"
        - name: VERSION
          value: "${VERSION}"
        livenessProbe:
          httpGet:
            path: /health
            port: 80
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health
            port: 80
          initialDelaySeconds: 5
          periodSeconds: 10
        resources:
          requests:
            memory: "32Mi"
            cpu: "25m"
          limits:
            memory: "64Mi"
            cpu: "100m"
---
apiVersion: v1
kind: Service
metadata:
  name: panchangam-frontend
  namespace: ${NAMESPACE}
  labels:
    app: panchangam-frontend
spec:
  selector:
    app: panchangam-frontend
  ports:
  - port: 80
    targetPort: 80
    name: http
EOF

    # Generate ingress if production
    if [ "$ENVIRONMENT" = "production" ]; then
        cat > "$manifest_dir/ingress.yaml" << EOF
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: panchangam-ingress
  namespace: ${NAMESPACE}
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/force-ssl-redirect: "true"
spec:
  tls:
  - hosts:
    - panchangam.example.com
    - api.panchangam.example.com
    secretName: panchangam-tls
  rules:
  - host: panchangam.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: panchangam-frontend
            port:
              number: 80
  - host: api.panchangam.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: panchangam-backend
            port:
              number: 80
EOF
    fi
    
    echo "$manifest_dir"
}

# Deploy to Kubernetes
deploy_to_kubernetes() {
    local manifest_dir="$1"
    
    log_info "Deploying to Kubernetes..."
    
    # Create namespace if it doesn't exist
    kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -
    
    if [ "$DRY_RUN" = "true" ]; then
        log_info "DRY RUN: Would apply the following manifests:"
        kubectl apply -f "$manifest_dir" --dry-run=client
    else
        kubectl apply -f "$manifest_dir"
        
        # Wait for deployments to be ready
        log_info "Waiting for deployments to be ready..."
        kubectl wait --for=condition=available --timeout=300s deployment/panchangam-backend -n "$NAMESPACE"
        kubectl wait --for=condition=available --timeout=300s deployment/panchangam-frontend -n "$NAMESPACE"
    fi
}

# Run smoke tests
run_smoke_tests() {
    log_info "Running smoke tests..."
    
    if [ "$DRY_RUN" = "true" ]; then
        log_info "DRY RUN: Would run smoke tests"
        return
    fi
    
    # Get service URL
    local service_url
    if [ "$ENVIRONMENT" = "production" ]; then
        service_url="https://api.panchangam.example.com"
    else
        service_url="http://$(kubectl get svc panchangam-backend -n "$NAMESPACE" -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo 'localhost:8080')"
    fi
    
    # Test health endpoint
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -f -s "$service_url/api/v1/health" > /dev/null; then
            log_success "Health check passed"
            break
        fi
        
        log_info "Health check attempt $attempt/$max_attempts failed, retrying..."
        sleep 10
        ((attempt++))
    done
    
    if [ $attempt -gt $max_attempts ]; then
        log_error "Health check failed after $max_attempts attempts"
        return 1
    fi
    
    # Test API endpoint
    local api_response
    api_response=$(curl -s "$service_url/api/v1/panchangam?date=2024-01-15&lat=12.9716&lng=77.5946" || echo "")
    
    if echo "$api_response" | grep -q "tithi"; then
        log_success "API smoke test passed"
    else
        log_error "API smoke test failed"
        return 1
    fi
}

# Main deployment flow
main() {
    log_info "Starting deployment to $ENVIRONMENT environment"
    log_info "Version: $VERSION"
    log_info "Registry: $REGISTRY"
    log_info "Namespace: $NAMESPACE"
    
    if [ "$DRY_RUN" = "true" ]; then
        log_warning "DRY RUN MODE - No changes will be made"
    fi
    
    # Confirmation for production
    if [ "$ENVIRONMENT" = "production" ] && [ "$FORCE_DEPLOY" = "false" ] && [ "$DRY_RUN" = "false" ]; then
        echo -n "Are you sure you want to deploy to production? (yes/no): "
        read -r confirmation
        if [ "$confirmation" != "yes" ]; then
            log_info "Deployment cancelled"
            exit 0
        fi
    fi
    
    check_prerequisites
    validate_images
    
    local manifest_dir
    manifest_dir=$(generate_manifests)
    
    deploy_to_kubernetes "$manifest_dir"
    
    if [ "$DRY_RUN" = "false" ]; then
        run_smoke_tests
        log_success "Deployment completed successfully!"
        
        # Cleanup
        rm -rf "$manifest_dir"
        
        # Display service information
        log_info "Service information:"
        kubectl get pods -n "$NAMESPACE" -l app=panchangam-backend
        kubectl get pods -n "$NAMESPACE" -l app=panchangam-frontend
        kubectl get svc -n "$NAMESPACE"
    fi
}

# Run main function
main "$@"