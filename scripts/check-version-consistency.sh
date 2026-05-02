#!/bin/bash
# check-version-consistency.sh
# 檢查 lunar-zenith 版本號在各文件中是否一致

set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION_FILE="$PROJECT_ROOT/VERSION"
README_FILE="$PROJECT_ROOT/README.md"
OPENAPI_FILE="$PROJECT_ROOT/contracts/openapi/lunar-zenith.yaml"

# 1. 讀取 VERSION 文件（唯一來源）
if [[ ! -f "$VERSION_FILE" ]]; then
    echo "❌ VERSION 文件不存在: $VERSION_FILE"
    exit 1
fi

VERSION=$(cat "$VERSION_FILE" | tr -d '[:space:]')
echo "📄 VERSION 文件版本: $VERSION"

# 2. 檢查 README.md badge
README_VERSION=$(sed -n 's/.*Version-\([^-]*\)-.*/\1/p' "$README_FILE" | head -1 || true)
if [[ -z "$README_VERSION" ]]; then
    echo "⚠️  無法從 README.md 解析版本號"
elif [[ "$README_VERSION" != "$VERSION" ]]; then
    echo "❌ README.md 版本不一致: $README_VERSION (期望: $VERSION)"
    exit 1
else
    echo "✅ README.md 版本一致: $README_VERSION"
fi

# 3. 檢查 OpenAPI info.version（注意: OpenAPI 標準無 v 前綴）
OPENAPI_VERSION=$(sed -n 's/^  version: \(.*\)/\1/p' "$OPENAPI_FILE" | head -1 || true)
EXPECTED_OPENAPI="${VERSION#v}" # 移除 v 前綴

if [[ -z "$OPENAPI_VERSION" ]]; then
    echo "⚠️  無法從 OpenAPI 文件解析版本號"
elif [[ "$OPENAPI_VERSION" != "$EXPECTED_OPENAPI" ]]; then
    echo "❌ OpenAPI info.version 不一致: $OPENAPI_VERSION (期望: $EXPECTED_OPENAPI)"
    exit 1
else
    echo "✅ OpenAPI info.version 一致: $OPENAPI_VERSION"
fi

# 4. 檢查 OpenAPI HealthStatus schema example
# 在 HealthStatus schema 區域內查找 version 的 example
EXAMPLE_VERSION=$(awk '/HealthStatus:/{flag=1} flag && /version:/{getline; if(/example:/){split($0, a, /"/); print a[2]; exit}}' "$OPENAPI_FILE" || true)
if [[ -n "$EXAMPLE_VERSION" && "$EXAMPLE_VERSION" != "$VERSION" ]]; then
    echo "❌ OpenAPI HealthStatus example 不一致: $EXAMPLE_VERSION (期望: $VERSION)"
    exit 1
else
    echo "✅ OpenAPI HealthStatus example 一致: ${EXAMPLE_VERSION:-未設置}"
fi

# 5. 檢查 go.mod 模塊路徑
GO_MOD_VERSION=$(sed -n 's/^module .*\/v\([0-9]*\).*/\1/p' "$PROJECT_ROOT/go.mod" || true)
MAJOR_VERSION="${VERSION#v}"     # 移除 v 前綴
MAJOR_VERSION="${MAJOR_VERSION%%.*}" # 取主版本號
if [[ -n "$GO_MOD_VERSION" && "$GO_MOD_VERSION" != "$MAJOR_VERSION" ]]; then
    echo "❌ go.mod 主版本不一致: v$GO_MOD_VERSION (期望: v$MAJOR_VERSION)"
    exit 1
else
    echo "✅ go.mod 主版本一致: v$GO_MOD_VERSION"
fi

echo ""
echo "🎉 所有版本號檢查通過！"
