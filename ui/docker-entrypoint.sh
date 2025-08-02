#!/bin/sh

# Environment configuration for containerized React app
# This script allows runtime environment variable injection

set -e

# Default API endpoint if not provided
API_ENDPOINT=${API_ENDPOINT:-"http://localhost:8080"}

# For remote deployment, detect if we're in a container network
# and use the appropriate endpoint
if [ "$API_ENDPOINT" = "http://gateway:8080" ]; then
    # We're using internal Docker networking, but need external access
    # Try to determine the host IP or use environment variable
    if [ -n "$HOST_IP" ]; then
        API_ENDPOINT="http://${HOST_IP}:8085"
    elif [ -n "$PUBLIC_API_ENDPOINT" ]; then
        API_ENDPOINT="$PUBLIC_API_ENDPOINT"
    else
        # Keep the internal endpoint for now, will be overridden by external config
        API_ENDPOINT="http://gateway:8080"
    fi
fi

# Create runtime config file
cat > /usr/share/nginx/html/runtime-config.js << EOF
window.__RUNTIME_CONFIG__ = {
  API_ENDPOINT: "${API_ENDPOINT}",
  BUILD_TIME: "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  VERSION: "${VERSION:-unknown}"
};
EOF

echo "Runtime configuration created:"
cat /usr/share/nginx/html/runtime-config.js

# Execute the original command
exec "$@"