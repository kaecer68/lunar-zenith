#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SKILLS_FILE="$ROOT_DIR/SKILLS.md"

if [[ ! -f "$SKILLS_FILE" ]]; then
  echo "[docs-check] ERROR: SKILLS.md not found at $SKILLS_FILE" >&2
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

# Required knowledge anchors (must stay in sync with current implementation)
require_in_skills "## 6.8 UI 與資料鍵名一致性（避免 undefined）"
require_in_skills "internal/webui/static/index.html"
require_in_skills 'name`, `type`, `description'
require_in_skills "NewAggregator(holidaySvc, chinaHolidaySvc)"
require_in_skills "Asia/Taipei"

# Stale anchor that should not return in skills map
forbid_in_skills "web/static/index.html"

# Enforce single UI source in repository
if [[ -f "$ROOT_DIR/web/static/index.html" ]]; then
  echo "[docs-check] ERROR: duplicate UI source exists: web/static/index.html" >&2
  exit 1
fi

echo "[docs-check] OK: skills map and UI-source consistency checks passed"
