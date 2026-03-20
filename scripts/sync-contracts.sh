#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

DEFAULT_CONTRACT_PORTS_FILE="$REPO_ROOT/contracts/runtime/ports.env"
if [[ ! -f "$DEFAULT_CONTRACT_PORTS_FILE" ]]; then
  DEFAULT_CONTRACT_PORTS_FILE="$REPO_ROOT/../destiny-contracts/runtime/ports.env"
fi

CONTRACT_PORTS_FILE="${CONTRACT_PORTS_FILE:-$DEFAULT_CONTRACT_PORTS_FILE}"
TARGET_PORTS_FILE="${TARGET_PORTS_FILE:-$REPO_ROOT/.env.ports}"

MODE="sync"
if [[ "${1:-}" == "--check" ]]; then
  MODE="check"
fi

required_keys=("LUNAR_GRPC_PORT" "LUNAR_REST_PORT")

get_value() {
  local key="$1"
  local value
  value="$(grep -E "^${key}=" "$CONTRACT_PORTS_FILE" | tail -n1 | cut -d= -f2-)"
  if [[ -z "$value" || ! "$value" =~ ^[0-9]+$ ]]; then
    echo "[contracts] 欄位錯誤: $key=$value" >&2
    exit 1
  fi
  echo "$value"
}

validate_contract_file() {
  if [[ ! -f "$CONTRACT_PORTS_FILE" ]]; then
    echo "[contracts] 找不到契約檔: $CONTRACT_PORTS_FILE" >&2
    exit 1
  fi
  for key in "${required_keys[@]}"; do
    if ! grep -qE "^${key}=" "$CONTRACT_PORTS_FILE"; then
      echo "[contracts] 缺少必要欄位: $key" >&2
      exit 1
    fi
  done
}

render_target_file() {
  local grpc_port rest_port
  grpc_port="$(get_value "LUNAR_GRPC_PORT")"
  rest_port="$(get_value "LUNAR_REST_PORT")"

  cat <<EOF
# 此檔案由 scripts/sync-contracts.sh 生成，請勿手動修改
# 契約來源：$CONTRACT_PORTS_FILE

LUNAR_GRPC_PORT=$grpc_port
LUNAR_REST_PORT=$rest_port

GRPC_PORT=$grpc_port
REST_PORT=$rest_port
EOF
}

validate_contract_file

tmp="$(mktemp)"
trap 'rm -f "$tmp"' EXIT
render_target_file > "$tmp"

if [[ "$MODE" == "check" ]]; then
  if [[ ! -f "$TARGET_PORTS_FILE" ]]; then
    echo "[contracts] 缺少本地檔案: $TARGET_PORTS_FILE" >&2
    exit 1
  fi
  if ! cmp -s "$tmp" "$TARGET_PORTS_FILE"; then
    echo "[contracts] 失敗：.env.ports 與契約不同步" >&2
    exit 1
  fi
  echo "[contracts] 驗證成功：.env.ports 與契約一致"
  exit 0
fi

cp "$tmp" "$TARGET_PORTS_FILE"
echo "[contracts] 已同步到: $TARGET_PORTS_FILE"
