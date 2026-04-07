#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
PORTS_FILE="$REPO_ROOT/.env.ports"
SYNC_SCRIPT="$REPO_ROOT/scripts/sync-contracts.sh"

if [[ ! -x "$SYNC_SCRIPT" ]]; then
  chmod +x "$SYNC_SCRIPT"
fi

bash "$SYNC_SCRIPT" >/dev/null

if [[ ! -f "$PORTS_FILE" ]]; then
  echo "[dev-clean] 缺少 $PORTS_FILE，請先執行 scripts/sync-contracts.sh" >&2
  exit 1
fi

# shellcheck disable=SC1090
source "$PORTS_FILE"

if [[ -z "${LUNAR_GRPC_PORT:-}" || -z "${LUNAR_REST_PORT:-}" ]]; then
  echo "[dev-clean] .env.ports 缺少必要欄位 LUNAR_GRPC_PORT 或 LUNAR_REST_PORT" >&2
  exit 1
fi

ports=("$LUNAR_GRPC_PORT" "$LUNAR_REST_PORT")
for port in "${ports[@]}"; do
  pids="$(lsof -tiTCP:"$port" -sTCP:LISTEN || true)"
  if [[ -n "$pids" ]]; then
    echo "[dev-clean] 清理 port $port: $pids"
    kill $pids 2>/dev/null || true
  fi
done
