#!/usr/bin/env bash
# GemCities — deploy script
#
# Pulls latest code from GitHub and deploys both the frontend and backend.
# Run on the VPS as root: /usr/local/bin/gemcities-deploy
#
# Usage:
#   gemcities-deploy          # deploy both frontend and backend
#   gemcities-deploy frontend # deploy frontend only
#   gemcities-deploy backend  # deploy backend only

set -euo pipefail

EDITOR_REPO="/opt/capsule-editor"
SERVICE_REPO="/opt/capsule-service"
FRONTEND_DEST="/srv/capsule-editor/public"
BINARY_DEST="/usr/local/bin/capsule-service"
SERVICE_NAME="capsule-service"

TARGET="${1:-all}"

deploy_frontend() {
    echo "==> Deploying frontend..."
    git -C "$EDITOR_REPO" pull --ff-only
    rsync -a --delete \
        --exclude='.git' \
        "$EDITOR_REPO/" "$FRONTEND_DEST/"
    echo "    Frontend deployed."
}

deploy_backend() {
    echo "==> Building backend..."
    git -C "$SERVICE_REPO" pull --ff-only
    (cd "$SERVICE_REPO" && go build -o /tmp/capsule-service-new .)
    echo "    Stopping service..."
    systemctl stop "$SERVICE_NAME"
    cp /tmp/capsule-service-new "$BINARY_DEST"
    rm /tmp/capsule-service-new
    echo "    Starting service..."
    systemctl start "$SERVICE_NAME"
    sleep 1
    if systemctl is-active --quiet "$SERVICE_NAME"; then
        echo "    Backend deployed and running."
    else
        echo "    ERROR: service failed to start. Check: journalctl -u $SERVICE_NAME -n 50"
        exit 1
    fi
}

case "$TARGET" in
    frontend) deploy_frontend ;;
    backend)  deploy_backend  ;;
    all)      deploy_frontend; deploy_backend ;;
    *)
        echo "Usage: $0 [frontend|backend|all]"
        exit 1
        ;;
esac

echo "==> Done."
