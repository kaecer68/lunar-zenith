# 🌙 Lunar-Zenith (算曆之巔)

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Version](https://img.shields.io/badge/Version-v4.0.0-blue)](https://github.com/kaecer68/lunar-zenith/releases/tag/v4.0.0)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Precision](https://img.shields.io/badge/Precision-Astronomical-blueviolet)](#)

**Lunar-Zenith** 是一款基於 Go 1.25+ 開發的全新一代**高精度曆法算曆引擎**。它完美融合了現代天體物理演算與傳統東方曆法智慧，專為追求「極致精度」與「台灣曆法標準」的開發者設計。

---

## 🚀 核心亮點 (Key Features)

- **🌌 天文級精度**: 基於 Jean Meeus 的 *Astronomical Algorithms* 與精簡版 VSOP87/ELP2000 理論，提供秒級精確的節氣與朔望月判定。
- **🇹🇼 台灣標準對齊**: 曆法規則、節氣計算與公眾假期完全遵循 **中華民國 (台灣)** 政府公告之標準 (對齊 DATA.GOV.TW)。
- **⛩️ 文化與宗教模型**: 
  - **核心干支**: 完整支持年月日時四柱、五虎遁、五鼠遁及立春精確換年。
  - **神煞系統**: 內建建除十二神、年驛馬、年桃花等常用神煞。
  - **擴充神煞**: 二十八星宿、值神、胎神、沖煞。
  - **宗教支持**: 自動換算佛曆 (Buddhist) 與道曆 (Taoist) 年份，內建台灣重要農曆宗教節日 (玉皇大帝、觀世音、媽祖等)。
- **🎉 節日/紀念日雙區域模型**:
  - 台灣：節日與紀念日可顯示，不以是否放假作為顯示條件。
  - 大陸：以台灣為基礎做規則覆寫（如 6/1 兒童節、8/1 建軍節、6月第3個週日父親節、9/10 教師節）。
  - `is_holiday` 僅表示休假狀態；工作日與週末預設不顯示節日名稱。
- **🪐 西洋占星輸出**:
  - 行星順行/逆行資訊：`western_astro`
  - 行星相位/交匯資訊：`aspects`
- **⚡ 高性能架構**: 全無狀態 (Stateless) 設計，支持 gRPC 與 REST 雙棧通訊，具備 Zero-Panic 的健壯性。
- **🌐 網頁查詢介面**: 內建現代化 Web UI，無需客戶端即可通過瀏覽器查詢完整曆法資訊。

---

## 🏛️ 技術架構 (Architecture)

本專案採用三層核心架構：
1. **Celestial (天體部)**: 處理儒略日 (JD)、Delta-T 修正、太陽/月球物理位置計算。
2. **Zodiac (文化部)**: 將天文數據轉化為干支、農曆日期、神煞、以及宗教曆法序列。
3. **Service (數據部)**: 載入政府 API/JSON 假期數據，並透過聚合器 (Aggregator) 提供統一的服務接口。

---

## 🛠️ 快速上手 (Quick Start)

### 1. 安裝與建置
```bash
# 克隆專案
git clone https://github.com/kaecer68/lunar-zenith.git
cd lunar-zenith

# 下載依賴
go mod tidy

# 編譯服務
go build -o bin/server ./cmd/server/main.go
```

### 2. 啟動服務
```bash
make dev
bash -c 'set -a; . ./.env.ports; set +a; echo "REST=http://localhost:${LUNAR_REST_PORT} gRPC=:${LUNAR_GRPC_PORT}"'
```
服務會使用 `destiny-contracts/runtime/ports.env` 同步出的契約 port 啟動。

> `make dev` 會先同步 `contracts/runtime/ports.env` 到本地 `.env.ports`，
> 再以契約 port 啟動服務（REST/gRPC）。

### 2.1 Port 契約同步工具
```bash
# 同步契約 port 到 .env.ports
make sync-contracts

# 驗證 .env.ports 是否與契約一致（適合 CI）
make verify-contracts

# 清理契約定義的本服務埠號佔用
make dev-clean
```

> PR/CI pipeline 會固定執行 `make verify-contracts` 作為 gate，任何未同步 `.env.ports` 的變更都會被拒絕；此檔案僅能由 `scripts/sync-contracts.sh` 生成，請勿手動修改。

### 3. 使用網頁介面
先執行 `make sync-contracts`，再依 `.env.ports` 內的 `LUNAR_REST_PORT` 訪問 `http://localhost:$LUNAR_REST_PORT`，即可使用圖形化查詢介面：
- 選擇日期查看完整曆法資訊
- 查看農曆、干支、節氣、宜忌、吉神方位
- 使用方向鍵快速切換日期
- 按 T 鍵快速回到今天

### 4. 調用 API 示例
獲取指定日期的完整曆法大禮包：
```bash
bash -c 'set -a; . ./.env.ports; set +a; curl "http://localhost:${LUNAR_REST_PORT}/v1/calendar?date=2024-02-10"'
```

### 4.1 主要 API 欄位

- 行政節日：`holiday_info`（台灣）、`china_holiday_info`（大陸）
- 西洋占星：`western_astro`（順逆行）、`aspects`（相位/交匯）

### 4.2 契約語義與相容性

- `solar_term`：
  目前 REST 與 gRPC 均僅保留單一 `SolarTerm` 物件，包含 `index`、`name`、`longitude`。
- `date`：
  REST query 與 gRPC request 都允許省略；省略時預設為 `Asia/Taipei` 今日，取中午 `12:00` 作為穩定採樣時刻。
- invalid date：
  REST 回 `400 Bad Request`；gRPC 回 `InvalidArgument`。
- 巢狀命名：
  目前 `shen_sha` / `mansion` / `daily_deity` / `fetal_god` / `clash_sha` 的 REST 子物件全面使用 snake_case，與 OpenAPI / proto 語義保持一致。

### 4.3 v4 遷移提示

- 若較舊版 REST client 依賴 `solar_term_detail`，請改讀 `solar_term`。
- 若較舊版 REST client 以字串讀取 `solar_term`，請改為讀取 `solar_term.name`。
- 若舊版 REST client 仍以 `mansion.Name`、`daily_deity.Type`、`clash_sha.ClashZodiac` 等 PascalCase 子鍵取值，請改為 `mansion.name`、`daily_deity.type`、`clash_sha.clash_zodiac`。
- gRPC client 若曾讀取 `solar_term_detail`，請改為直接使用 `solar_term`；其餘巢狀 message 命名本輪無需遷移。

---

## 🔁 契約同步與維護

- OpenAPI 契約（workspace 單一來源）：`../destiny-contracts/openapi/lunar-zenith.yaml`
- gRPC proto 來源（服務內）：`proto/lunar.proto`
- gRPC proto 發佈（供 workspace 治理/同步）：`contracts/proto/lunar.proto`
- 契約治理入口：`make check-docs-consistency`
- 每次新增欄位後請同步：
  1. 更新契約（OpenAPI + `proto/lunar.proto`）
  2. 將 `proto/lunar.proto` 同步到 `contracts/proto/lunar.proto`
  3. 更新 `internal/service/calendar_response.go`（REST/gRPC 映射）
  4. 重新生成 `gen/*.pb.go`
  5. 先執行 `openapi-generator validate -i contracts/openapi/lunar-zenith.yaml`
  6. 最後執行 `make verify-all`

---

## 📝 授權協議 (License)

本專案基於 **[MIT License](LICENSE)** 進行開源。您可以自由地使用、修改及分發，但也請保留原作者信息。

---

## 👨‍💻 作者 (Author)

**Kaecer** 
- GitHub: [@kaecer68](https://github.com/kaecer68)
- Twitter: [@kaecer68](https://twitter.com/kaecer)
- 德凱 GoLuck 實用易理作品，旨在將傳統曆法計算現代化、精密化。
- https://goluck.im/

> *「算曆之巔，意在精確；天人之際，存乎一心。」*
