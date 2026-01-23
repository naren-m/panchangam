#!/bin/sh

# Panchangam PostgreSQL Backup Script
# This script creates compressed backups and optionally uploads to S3

set -e

# Configuration
BACKUP_DIR="${BACKUP_DIR:-/backups}"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="${BACKUP_DIR}/panchangam_backup_${TIMESTAMP}.sql.gz"
RETENTION_DAYS="${BACKUP_RETENTION_DAYS:-30}"

# Database connection
PGHOST="${PGHOST:-postgres}"
PGPORT="${PGPORT:-5432}"
PGDATABASE="${PGDATABASE:-panchangam}"
PGUSER="${PGUSER:-panchangam}"

# S3 Configuration (optional)
S3_BUCKET="${S3_BUCKET:-}"
S3_PREFIX="${S3_PREFIX:-backups/postgres}"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Logging function
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1"
}

log "${GREEN}Starting backup process...${NC}"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Check database connection
if ! pg_isready -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" > /dev/null 2>&1; then
    log "${RED}Error: Cannot connect to database${NC}"
    exit 1
fi

log "Database connection successful"

# Create backup
log "Creating backup: $BACKUP_FILE"
if pg_dump -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" \
    --verbose \
    --format=plain \
    --no-owner \
    --no-acl \
    | gzip > "$BACKUP_FILE"; then

    BACKUP_SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    log "${GREEN}Backup created successfully: $BACKUP_FILE ($BACKUP_SIZE)${NC}"
else
    log "${RED}Error: Backup failed${NC}"
    exit 1
fi

# Upload to S3 if configured
if [ -n "$S3_BUCKET" ]; then
    log "Uploading backup to S3: s3://${S3_BUCKET}/${S3_PREFIX}/"

    if command -v aws > /dev/null 2>&1; then
        if aws s3 cp "$BACKUP_FILE" "s3://${S3_BUCKET}/${S3_PREFIX}/$(basename $BACKUP_FILE)"; then
            log "${GREEN}Backup uploaded to S3 successfully${NC}"
        else
            log "${YELLOW}Warning: Failed to upload backup to S3${NC}"
        fi
    else
        log "${YELLOW}Warning: AWS CLI not found, skipping S3 upload${NC}"
    fi
fi

# Clean up old local backups
log "Cleaning up backups older than ${RETENTION_DAYS} days"
find "$BACKUP_DIR" -name "panchangam_backup_*.sql.gz" -type f -mtime +${RETENTION_DAYS} -delete
REMAINING_BACKUPS=$(find "$BACKUP_DIR" -name "panchangam_backup_*.sql.gz" -type f | wc -l)
log "${GREEN}Cleanup complete. Remaining backups: ${REMAINING_BACKUPS}${NC}"

# Clean up old S3 backups if configured
if [ -n "$S3_BUCKET" ] && command -v aws > /dev/null 2>&1; then
    log "Cleaning up old S3 backups"
    CUTOFF_DATE=$(date -d "${RETENTION_DAYS} days ago" +%Y-%m-%d 2>/dev/null || date -v-${RETENTION_DAYS}d +%Y-%m-%d)

    aws s3 ls "s3://${S3_BUCKET}/${S3_PREFIX}/" | while read -r line; do
        BACKUP_DATE=$(echo "$line" | awk '{print $1}')
        BACKUP_NAME=$(echo "$line" | awk '{print $4}')

        if [ "$BACKUP_DATE" \< "$CUTOFF_DATE" ]; then
            log "Deleting old S3 backup: $BACKUP_NAME"
            aws s3 rm "s3://${S3_BUCKET}/${S3_PREFIX}/${BACKUP_NAME}"
        fi
    done
fi

# Create a latest symlink
ln -sf "$BACKUP_FILE" "${BACKUP_DIR}/latest.sql.gz"

log "${GREEN}Backup process completed successfully${NC}"
log "Backup location: $BACKUP_FILE"
log "Latest backup: ${BACKUP_DIR}/latest.sql.gz"
