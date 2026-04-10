#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SKILLS_FILE="$ROOT_DIR/SKILLS.md"
README_FILE="$ROOT_DIR/README.md"
SERVER_MAIN="$ROOT_DIR/cmd/server/main.go"
OPENAPI_FILE="$ROOT_DIR/contracts/openapi/lunar-zenith.yaml"
LUNAR_PROTO_SRC="$ROOT_DIR/proto/lunar.proto"
LUNAR_PROTO_CONTRACT="$ROOT_DIR/contracts/proto/lunar.proto"
UI_FILE="$ROOT_DIR/internal/webui/static/index.html"

if [[ ! -f "$SKILLS_FILE" ]]; then
  echo "[docs-check] ERROR: SKILLS.md not found at $SKILLS_FILE" >&2
  exit 1
fi

if [[ ! -f "$README_FILE" || ! -f "$SERVER_MAIN" || ! -f "$OPENAPI_FILE" ]]; then
  echo "[docs-check] ERROR: required version-tracking files missing" >&2
  exit 1
fi

require_in_skills() {
  local needle="$1"
  if ! grep -Fq "$needle" "$SKILLS_FILE"; then
    echo "[docs-check] ERROR: missing in SKILLS.md -> $needle" >&2
    exit 1
  fi
}

forbid_in_skills() {
  local needle="$1"
  if grep -Fq "$needle" "$SKILLS_FILE"; then
    echo "[docs-check] ERROR: forbidden stale text in SKILLS.md -> $needle" >&2
    exit 1
  fi
}

forbid_in_ui() {
  local needle="$1"
  if grep -Fq "$needle" "$UI_FILE"; then
    echo "[docs-check] ERROR: forbidden stale UI contract fallback -> $needle" >&2
    exit 1
  fi
}

# Required knowledge anchors (must stay in sync with current implementation)
require_in_skills "## 6.8 UI 與資料鍵名一致性（避免 undefined）"
require_in_skills "internal/webui/static/index.html"
require_in_skills 'name`, `type`, `description'
require_in_skills 'full_name`, `clash_zodiac`, `sha_direction'
require_in_skills "NewAggregator(holidaySvc, chinaHolidaySvc)"
require_in_skills "Asia/Taipei"
require_in_skills "contracts/proto/lunar.proto"
require_in_skills 'solar_term` 在 REST / gRPC 皆為物件'
require_in_skills 'solar_term_detail` 已移除'
require_in_skills 'REST 巢狀欄位統一為 snake_case'

# Stale anchor that should not return in skills map
forbid_in_skills "web/static/index.html"
forbid_in_skills "仍沿用 PascalCase 欄位"

# Enforce single UI source in repository
if [[ -f "$ROOT_DIR/web/static/index.html" ]]; then
  echo "[docs-check] ERROR: duplicate UI source exists: web/static/index.html" >&2
  exit 1
fi

if [[ ! -f "$UI_FILE" ]]; then
  echo "[docs-check] ERROR: UI file missing: $UI_FILE" >&2
  exit 1
fi

forbid_in_ui "f?.Name ??"
forbid_in_ui "f?.Type ??"
forbid_in_ui "f?.Description ??"
forbid_in_ui "lunar.Month || lunar.month"
forbid_in_ui "lunar.Day || lunar.day"
forbid_in_ui "lunar.IsLeap || lunar.is_leap"
forbid_in_ui "lunar.Year || lunar.year"
forbid_in_ui "pillars.Year || pillars.year"
forbid_in_ui "pillars.Month || pillars.month"
forbid_in_ui "pillars.Day || pillars.day"
forbid_in_ui "pillars.Hour || pillars.hour"
forbid_in_ui "StemIndex"
forbid_in_ui "BranchIndex"
forbid_in_ui "ld.Day"
forbid_in_ui "ld.IsLeap"
forbid_in_ui "ld.Month"

# Enforce gRPC contract publication and sync
if [[ ! -f "$LUNAR_PROTO_SRC" ]]; then
  echo "[docs-check] ERROR: source proto missing: $LUNAR_PROTO_SRC" >&2
  exit 1
fi

if [[ ! -f "$LUNAR_PROTO_CONTRACT" ]]; then
  echo "[docs-check] ERROR: contract proto missing: $LUNAR_PROTO_CONTRACT" >&2
  echo "[docs-check] HINT: cp proto/lunar.proto contracts/proto/lunar.proto" >&2
  exit 1
fi

if ! cmp -s "$LUNAR_PROTO_SRC" "$LUNAR_PROTO_CONTRACT"; then
  echo "[docs-check] ERROR: contracts/proto/lunar.proto is out of sync with proto/lunar.proto" >&2
  diff -u "$LUNAR_PROTO_CONTRACT" "$LUNAR_PROTO_SRC" >&2 || true
  echo "[docs-check] HINT: cp proto/lunar.proto contracts/proto/lunar.proto" >&2
  exit 1
fi

server_version="$(rg -o 'const serviceVersion = \"v[0-9.]+\"' "$SERVER_MAIN" | awk -F'\"' 'NR==1{print $2}')"
openapi_version="$(rg -o '^  version: [0-9.]+' "$OPENAPI_FILE" | awk '{print $2}' | head -n1)"
readme_badge_version="$(rg -o 'Version-v[0-9.]+-blue' "$README_FILE" | head -n1 | cut -d- -f2)"

if [[ -z "$server_version" || -z "$openapi_version" || -z "$readme_badge_version" ]]; then
  echo "[docs-check] ERROR: failed to parse one or more version anchors" >&2
  exit 1
fi

if [[ "v$openapi_version" != "$server_version" || "$readme_badge_version" != "$server_version" ]]; then
  echo "[docs-check] ERROR: version drift detected" >&2
  echo "[docs-check] server=$server_version openapi=v$openapi_version readme=$readme_badge_version" >&2
  exit 1
fi

echo "[docs-check] OK: skills map, UI-source, version, and proto contract sync checks passed"
