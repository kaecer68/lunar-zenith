#!/usr/bin/env bash
#
# sync-contracts.sh - Thin wrapper for central workspace sync
# 
# This script delegates to the central workspace-sync engine in destiny-cloud.
#

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
WORKSPACE_ROOT="$(cd "$REPO_ROOT/.." && pwd)"
DESTINY_CLOUD="$WORKSPACE_ROOT/destiny-cloud"

# Build central sync tool if needed
if [[ ! -f "$DESTINY_CLOUD/bin/workspace-sync" ]]; then
  echo "[info] Building central workspace-sync..." >&2
  (cd "$DESTINY_CLOUD" && go build -o bin/workspace-sync ./cmd/workspace-sync)
fi

# Parse mode
MODE=""
if [[ "${1:-}" == "--check" ]]; then
  MODE="-check"
fi

# Get service name from directory
SERVICE_NAME="$(basename "$REPO_ROOT")"

# Forward to central sync engine
exec "$DESTINY_CLOUD/bin/workspace-sync" \
  -registry "$DESTINY_CLOUD/configs/workspace/services.yaml" \
  -contracts "$WORKSPACE_ROOT/destiny-contracts" \
  -workspace "$WORKSPACE_ROOT" \
  -service "$SERVICE_NAME" \
  $MODE
