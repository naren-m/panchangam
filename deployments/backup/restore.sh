#!/bin/sh

# Panchangam PostgreSQL Restore Script
# This script restores from a backup file or S3

set -e

# Configuration
BACKUP_DIR="${BACKUP_DIR:-/backups}"
BACKUP_FILE="${1:-}"

# Database connection
PGHOST="${PGHOST:-postgres}"
PGPORT="${PGPORT:-5432}"
PGDATABASE="${PGDATABASE:-panchangam}"
PGUSER="${PGUSER:-panchangam}"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Logging function
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# Usage function
usage() {
    echo "Usage: $0 <backup_file>"
    echo ""
    echo "Examples:"
    echo "  $0 /backups/panchangam_backup_20250118_020000.sql.gz"
    echo "  $0 latest  (uses the latest backup)"
    echo ""
    exit 1
}

# Check if backup file is specified
if [ -z "$BACKUP_FILE" ]; then
    log "${RED}Error: No backup file specified${NC}"
    usage
fi

# Handle 'latest' keyword
if [ "$BACKUP_FILE" = "latest" ]; then
    BACKUP_FILE="${BACKUP_DIR}/latest.sql.gz"
fi

# Check if backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    log "${RED}Error: Backup file not found: $BACKUP_FILE${NC}"
    exit 1
fi

log "${YELLOW}WARNING: This will restore the database from backup${NC}"
log "${YELLOW}Current database will be dropped and recreated${NC}"
log "Backup file: $BACKUP_FILE"
log ""
read -p "Are you sure you want to continue? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    log "Restore cancelled"
    exit 0
fi

# Check database connection
if ! pg_isready -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" > /dev/null 2>&1; then
    log "${RED}Error: Cannot connect to database${NC}"
    exit 1
fi

log "Database connection successful"

# Drop existing database
log "${YELLOW}Dropping existing database...${NC}"
psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d postgres -c "DROP DATABASE IF EXISTS $PGDATABASE;"

# Create new database
log "Creating new database..."
psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d postgres -c "CREATE DATABASE $PGDATABASE OWNER $PGUSER;"

# Restore from backup
log "Restoring from backup..."
if gunzip -c "$BACKUP_FILE" | psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" > /dev/null 2>&1; then
    log "${GREEN}Restore completed successfully${NC}"
else
    log "${RED}Error: Restore failed${NC}"
    exit 1
fi

log "${GREEN}Database restored successfully from $BACKUP_FILE${NC}"
