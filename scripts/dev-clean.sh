#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
PORTS_FILE="$REPO_ROOT/.env.ports"
CONTRACT_PORTS_FILE="$REPO_ROOT/contracts/runtime/ports.env"

LUNAR_GRPC_PORT="${LUNAR_GRPC_PORT:-50051}"
LUNAR_REST_PORT="${LUNAR_REST_PORT:-8080}"

if [[ -f "$PORTS_FILE" ]]; then
  # shellcheck disable=SC1090
  source "$PORTS_FILE"
elif [[ -f "$CONTRACT_PORTS_FILE" ]]; then
  # shellcheck disable=SC1090
  source "$CONTRACT_PORTS_FILE"
fi

ports=("$LUNAR_GRPC_PORT" "$LUNAR_REST_PORT")
for port in "${ports[@]}"; do
  pids="$(lsof -tiTCP:"$port" -sTCP:LISTEN || true)"
  if [[ -n "$pids" ]]; then
    echo "[dev-clean] 清理 port $port: $pids"
    kill $pids 2>/dev/null || true
  fi
done
