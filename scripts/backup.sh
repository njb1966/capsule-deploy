#!/usr/bin/env bash
# GemCities — Restic backup to Backblaze B2
#
# Backs up:
#   /srv/capsules                    — all user capsule files
#   /var/lib/capsule-service/        — SQLite database
#
# Schedule: run daily at 3:00 AM via cron.
# Install: cp backup.sh /usr/local/bin/gemcities-backup
#          chmod 700 /usr/local/bin/gemcities-backup
#          chown root:root /usr/local/bin/gemcities-backup
#
# Add to root crontab (crontab -e):
#   0 3 * * * /usr/local/bin/gemcities-backup >> /var/log/gemcities-backup.log 2>&1
#
# First-time setup:
#   1. Set B2_ACCOUNT_ID, B2_ACCOUNT_KEY, and RESTIC_PASSWORD in /etc/capsule-service/env
#   2. Run: restic -r b2:gemcities-backups:/data init
#   3. Run this script manually to verify the first backup succeeds.

set -euo pipefail

# ── Configuration ────────────────────────────────────────────────────────────

RESTIC_REPO="b2:gemcities-backups:/data"
LOG_FILE="/var/log/gemcities-backup.log"
ALERT_EMAIL="abuse@gemcities.com"   # receives failure alerts

# Backup sources
BACKUP_PATHS=(
    /srv/capsules
    /var/lib/capsule-service
)

# Retention policy
KEEP_DAILY=7
KEEP_WEEKLY=4
KEEP_MONTHLY=6

# ── Load secrets ─────────────────────────────────────────────────────────────
# Expects B2_ACCOUNT_ID, B2_ACCOUNT_KEY, RESTIC_PASSWORD in the environment.
# If running from cron, source the env file directly:

ENV_FILE="/etc/capsule-service/env"
if [[ -f "$ENV_FILE" ]]; then
    # shellcheck disable=SC1090
    source "$ENV_FILE"
fi

: "${B2_ACCOUNT_ID:?B2_ACCOUNT_ID not set}"
: "${B2_ACCOUNT_KEY:?B2_ACCOUNT_KEY not set}"
: "${RESTIC_PASSWORD:?RESTIC_PASSWORD not set}"

export B2_ACCOUNT_ID B2_ACCOUNT_KEY RESTIC_PASSWORD

# ── Backup ───────────────────────────────────────────────────────────────────

TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
echo "[$TIMESTAMP] Starting backup..."

if restic -r "$RESTIC_REPO" backup "${BACKUP_PATHS[@]}" --tag gemcities; then
    echo "[$TIMESTAMP] Backup completed successfully."
else
    echo "[$TIMESTAMP] ERROR: Backup failed." >&2
    echo "GemCities backup FAILED at $TIMESTAMP. Check $LOG_FILE on the server." \
        | mail -s "GemCities backup failure" "$ALERT_EMAIL" 2>/dev/null || true
    exit 1
fi

# ── Prune old snapshots ───────────────────────────────────────────────────────

echo "[$TIMESTAMP] Pruning old snapshots..."
restic -r "$RESTIC_REPO" forget \
    --keep-daily  "$KEEP_DAILY"   \
    --keep-weekly "$KEEP_WEEKLY"  \
    --keep-monthly "$KEEP_MONTHLY" \
    --prune \
    --tag gemcities

echo "[$TIMESTAMP] Done."
