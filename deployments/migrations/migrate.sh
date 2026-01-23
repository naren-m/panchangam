#!/bin/sh

# Database Migration Script for Panchangam
# Usage: ./migrate.sh [up|down|create|version|force]

set -e

# Configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-panchangam}"
DB_USER="${DB_USER:-panchangam}"
DB_PASSWORD="${DB_PASSWORD:-panchangam123}"
MIGRATIONS_DIR="${MIGRATIONS_DIR:-./deployments/migrations}"

# Construct database URL
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if migrate is installed
if ! command -v migrate &> /dev/null; then
    echo "${RED}Error: golang-migrate is not installed${NC}"
    echo "Install it with:"
    echo "  macOS: brew install golang-migrate"
    echo "  Linux: curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && mv migrate /usr/local/bin/"
    exit 1
fi

# Function to run migrations up
migrate_up() {
    echo "${GREEN}Running migrations up...${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up
    echo "${GREEN}Migrations completed successfully${NC}"
}

# Function to run migrations down
migrate_down() {
    echo "${YELLOW}Rolling back migrations...${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" down
    echo "${GREEN}Rollback completed${NC}"
}

# Function to create new migration
migrate_create() {
    if [ -z "$2" ]; then
        echo "${RED}Error: Please provide a migration name${NC}"
        echo "Usage: ./migrate.sh create migration_name"
        exit 1
    fi
    echo "${GREEN}Creating new migration: $2${NC}"
    migrate create -ext sql -dir "$MIGRATIONS_DIR" -seq "$2"
    echo "${GREEN}Migration files created${NC}"
}

# Function to check migration version
migrate_version() {
    echo "${GREEN}Current migration version:${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" version
}

# Function to force migration version
migrate_force() {
    if [ -z "$2" ]; then
        echo "${RED}Error: Please provide a version number${NC}"
        echo "Usage: ./migrate.sh force <version>"
        exit 1
    fi
    echo "${YELLOW}Forcing migration to version: $2${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" force "$2"
    echo "${GREEN}Version forced${NC}"
}

# Function to show migration status
migrate_status() {
    echo "${GREEN}Migration status:${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" version
}

# Main command handler
case "${1:-up}" in
    up)
        migrate_up
        ;;
    down)
        migrate_down
        ;;
    create)
        migrate_create "$@"
        ;;
    version|status)
        migrate_version
        ;;
    force)
        migrate_force "$@"
        ;;
    *)
        echo "Usage: $0 {up|down|create <name>|version|force <version>}"
        echo ""
        echo "Commands:"
        echo "  up              - Apply all pending migrations"
        echo "  down            - Rollback all migrations"
        echo "  create <name>   - Create a new migration file"
        echo "  version         - Show current migration version"
        echo "  force <version> - Force set migration version (use with caution)"
        exit 1
        ;;
esac
