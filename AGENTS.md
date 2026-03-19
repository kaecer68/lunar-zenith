# Lunar-Zenith Project Knowledge Base

**Generated:** 2026-03-18  
**Commit:** see `git log -1 --oneline`  
**Branch:** see `git branch --show-current`  
**Contract Version:** lunar-zenith v1.4.0

## OVERVIEW

高精度曆法算曆引擎 (Go 1.25+)。天文級精度節氣/朔望月計算，台灣曆法標準對齊。REST + gRPC 雙棧服務。

## STRUCTURE

```
lunar-zenith/
├── cmd/              # 入口：server (REST+gRPC), test_grpc
├── api/v1/           # Protobuf 生成代碼 (DO NOT EDIT)
├── pkg/celestial/    # 天體計算：儒略日，Delta-T, 太陽/月球位置
├── pkg/zodiac/       # 文化曆法：干支，農曆，神煞，宗教曆
├── internal/service/ # 服務層：API 聚合，gRPC/REST handler
├── configs/          # 假期數據 JSON
└── contracts/        # symlink → destiny-contracts (OpenAPI 契約)
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| API 契約 | `contracts/openapi/lunar-zenith.yaml` | 先更新契約，再 `make generate` |
| 啟動入口 | `cmd/server/main.go` | Gin + gRPC + Web UI 並行 |
| 天文計算 | `pkg/celestial/*.go` | JD 轉換，Delta-T 估算 |
| 干支/農曆 | `pkg/zodiac/*.go` | 五虎遁，五鼠遁，建除十二神 |
| 業務邏輯 | `internal/service/*.go` | Aggregator, HolidayService |
| 網頁介面 | `internal/webui/static/index.html` | 內建 Web UI |
| 假期數據 | `configs/holidays_2024_sample.json` | 台灣公眾假期 |
| 神煞系統 | `pkg/zodiac/shensha.go` | 二十八星宿、值神、胎神、沖煞 |
| 宗教節日 | `pkg/zodiac/festival.go` | 農曆節日清單 |

## CODE MAP

| Symbol | Type | Location | Role |
|--------|------|----------|------|
| `PrecisionTime` | struct | `pkg/celestial/base.go` | 高精度時間封裝 (JD/JDE/DeltaT) |
| `NewYearSexagenary` | func | `pkg/zodiac/sexagenary.go` | 年干支計算 (基準：西元 4 年甲子) |
| `GetDaySexagenary` | func | `pkg/zodiac/sexagenary.go` | 日干支計算 (JD 基準) |
| `Aggregator` | struct | `internal/service/aggregator.go` | 服務聚合器 |
| `HolidayService` | struct | `internal/service/holiday.go` | 假期載入/查詢 |
| `WebUI` | package | `internal/webui/*.go` | 網頁查詢介面 |

## CONVENTIONS

- **Contract-first**: API 變更 → 更新 OpenAPI → `make generate` → 實現
- **Stateless**: 所有服務無狀態，支持水平擴展
- **Zero-Panic**: `cmd/` 外禁止 `panic()`，錯誤返回 `error`
- **繁體中文**: 所有曆法名稱使用繁體中文

## ANTI-PATTERNS (THIS PROJECT)

- ❌ 直接修改 `api/v1/*.go` (protoc 生成)
- ❌ 添加契約未定義的 API 欄位
- ❌ `internal/service` 外直接訪問 `HolidayService`
- ❌ 忽略 Delta-T 修正 (天文計算必須使用 JDE)

## UNIQUE STYLES

- **雙時制**: UT (民用時) / TT (曆書時) 並行，Delta-T 橋接
- **基準校正**: 干支計算使用歷史基準 (西元 4 年甲子，2000-01-01 戊午日)
- **台灣優先**: 假期/曆法規則遵循中華民國政府公告標準，內建重要宗教節日
- **網頁優先**: 內建 Web UI，根路徑 `/` 提供圖形化查詢介面
- **完整神煞**: 二十八星宿、值神、胎神、沖煞、農曆節日 (v1.4.0)

## COMMANDS

```bash
# 編譯
go build -o bin/server ./cmd/server/main.go

# 測試
go test ./...

# 契約生成 (需 Makefile，暫無)
# openapi-generator generate -i contracts/openapi/lunar-zenith.yaml -g go-server -o api/v1/

# 啟動服務
./bin/server  # REST:8080, gRPC:50051
```

## NOTES

- **Leap Month**: 閏月判定待完善 (`pkg/zodiac/lunar_engine.go:59` TODO)
- **Generated Code**: `api/v1/` 為 protoc 生成，禁止手動修改
- **Contract Sync**: `contracts/` 為 symlink，需確保 `destiny-contracts` 存在

## RELATED DOCS

- `contracts/README.md` - 契約層完整文檔
- `contracts/TASK-BOARD.md` - 跨服務任務看板
- `contracts/HANDOFF.md` - AI 交接報告模板
- `PRD.md` - 產品需求文檔
