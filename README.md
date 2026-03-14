# 🌙 Lunar-Zenith (算曆之巔)

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
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
  - **宗教支持**: 自動換算佛曆 (Buddhist) 與道曆 (Taoist) 年份。
- **⚡ 高性能架構**: 全無狀態 (Stateless) 設計，支持 gRPC 與 REST 雙棧通訊，具備 Zero-Panic 的健壯性。

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

### 2. 啟動 REST API
```bash
./bin/server
```
預設服務將開啟於 `http://localhost:8080`。

### 3. 調用示例
獲取指定日期的完整曆法大禮包：
```bash
curl "http://localhost:8080/v1/calendar?date=2024-02-10"
```

---

## 📝 授權協議 (License)

本專案基於 **[MIT License](LICENSE)** 進行開源。您可以自由地使用、修改及分發，但也請保留原作者信息。

---

## 👨‍💻 作者 (Author)

**Kaecer** 
- GitHub: [@kaecer68](https://github.com/kaecer68)
- 德凱 GoLuck 實用易理作品，旨在將傳統曆法計算現代化、精密化。

> *「算曆之巔，意在精確；天人之際，存乎一心。」*
