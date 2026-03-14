# Lunar-Zenith (算曆之巔): 全新一代高精度曆法算曆引擎

## 1. 核心定位 (Core Positioning)
- **基礎環境**: 基於 Go 1.25.6，最低相容版本定為 Go 1.25.0。
- **標準語系**: 全案採用「繁體中文」(Traditional Chinese) 作為原生語言。
- **文化基準**: 曆法規則、節氣計算與公眾假期完全對準 **台灣 (Taiwan/ROC)** 政府公布之標準。
- **精度標準**: 天文級運算，對齊紫金山天文台/瑞士編算表 (Swiss Ephemeris) 精度。

## 2. 功能範疇 (Functional Scope)
- **時間跨度**: 西元 1900 年 01 月 01 日 至 2150 年 12 月 31 日。
- **核心能力**:
  - 精確天文算曆 (JD, Delta-T, 朔望月, 定氣)。
  - 台灣法定節假日與補班課逻辑（對準行政院人事行政總處 DATA.GOV.TW 數據）。
  - 干支紀年/月/日/時、建除十二神、神煞計算模型。
  - **宗教曆法**: 支持佛曆 (Buddhist Calendar) 與道曆 (Taoist Calendar) 的精確演算。
- **服務接口**:
  - **gRPC**: 提供高效能、強型別的服務間通訊。
  - **REST (OpenAPI 3.0)**: 提供標準 Web API，並自動生成 Swagger 定制化文檔。

## 3. 核心架構 (Core Architecture)
- **Celestial (天體部 - 公開)**: 位於 `pkg/celestial`，處理天文計算法、天體物理位置與朔望月。
- **Zodiac (文化部 - 公開)**: 位於 `pkg/zodiac`，處理台灣民俗曆法、干支轉化、神煞模型及宗教曆法。
- **Service (服務部 - 內部)**: 位於 `internal/service`，聚合曆法運算並載入行政假期數據。

## 4. 性能與安全指標 (Security & Performance)
- **並發安全**: 全無狀態 (Stateless) 設計，內部使用讀寫鎖保護計算緩存。
- **Zero-Panic**: 嚴格錯誤返回路徑，排除所有隱式陣列越界風險。
- **高能效率**: 減少內存分配 (Allocations)，確保在高併發環境下的穩定反應。

## 5. 數據來源與維護 (Data & Maintenance)
- **假期來源**: 行政院人事行政總處公開資料 (DATA.GOV.TW)。
- **更新機制**: 透過外部外掛式 JSON 定義，支持 Zero-Downtime 年度更新。
