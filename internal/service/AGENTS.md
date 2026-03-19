# internal/service - Service Layer

**Domain**: API aggregation, gRPC/REST handlers, business logic orchestration

## OVERVIEW

Service layer that aggregates celestial/zodiac calculations with holiday data. Provides unified API via gRPC and REST endpoints.

## STRUCTURE

```
internal/service/
├── aggregator.go      # 服務聚合器：組合天文，曆法，假期數據
├── grpc_server.go     # gRPC 服務器實現
├── rest_handler.go    # REST handler (Gin)
├── holiday.go         # 假期服務：JSON 載入，查詢
├── almanac.go         # 曆書計算：完整日曆大禮包
└── *_test.go          # 單元測試
```

## WHERE TO LOOK

| Task | File | Notes |
|------|------|-------|
| API 聚合 | `aggregator.go` | 組合 HolidayService + 曆法計算 |
| gRPC 實現 | `grpc_server.go` | `LunarServiceServer` |
| REST 路由 | `rest_handler.go` | Gin handler 註冊 |
| 假期載入 | `holiday.go:LoadFromJSON()` | JSON 文件載入 |
| 曆書計算 | `almanac.go` | 完整日曆響應組裝 |

## CODE MAP

| Symbol | Type | Role |
|--------|------|------|
| `HolidayService` | struct | 假期數據載入/查詢 |
| `Aggregator` | struct | 服務聚合器 (依賴 HolidayService) |
| `GrpcServer` | struct | gRPC 服務器實現 |
| `RestHandler` | struct | REST handler (Gin) |
| `NewAggregator(holidaySvc)` | constructor | 創建聚合器 |
| `RegisterRoutes(r *gin.Engine)` | method | REST 路由註冊 |

## CONVENTIONS

- **Dependency Injection**: `Aggregator` 依賴 `HolidayService` (構造函數注入)
- **Error Handling**: 返回 `error`, 不 panic (與 `cmd/` 區分)
- **Stateless**: 服務無狀態，支持並發請求
- **JSON Charset**: REST 響應標頭 `Content-Type: application/json; charset=utf-8`

## ANTI-PATTERNS

- ❌ 直接實例化 `Aggregator` (必須用 `NewAggregator()`)
- ❌ 在 service 層調用 `log.Fatal()` (返回 error)
- ❌ 繞過 `HolidayService` 直接讀取 JSON
- ❌ 修改 protobuf 生成的 interface (實現 `LunarServiceServer` 即可)

## API ENDPOINTS

- **REST**: `GET /v1/calendar?date=YYYY-MM-DD` → 完整曆法數據
- **gRPC**: `LunarService.GetCalendar()`
- **Health**: `GET /` → project, version, status

## DEPENDENCIES

- **Internal**: `pkg/celestial/`, `pkg/zodiac/`
- **External**: Gin, gRPC, `configs/holidays_2024_sample.json`

## TESTING

```bash
go test ./internal/service/... -v
```

## CONFIGURATION

| Env Var | Default | Description |
|---------|---------|-------------|
| `GRPC_PORT` | 50051 | gRPC 端口 |
| HTTP | 8080 | REST 端口 (hardcoded) |

## NOTES

- Holiday data: 2024 sample only (expand needed)
- Contract sync: API must match `contracts/openapi/lunar-zenith.yaml`
- Stateless, concurrent-safe
