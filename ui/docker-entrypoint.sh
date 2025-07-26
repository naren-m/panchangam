#!/bin/sh

# Environment configuration for containerized React app
# This script allows runtime environment variable injection

set -e

# Default API endpoint if not provided
API_ENDPOINT=${API_ENDPOINT:-"http://localhost:8080"}

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