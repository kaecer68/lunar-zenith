# Lunar-Zenith Skills Map

**Purpose**: 本文件提供 AI 助手快速了解黃曆計算系統與本專案實現的知識地圖。
**Version**: 1.0.0  
**Updated**: 2026-03-25  
**For**: AI Assistants working on lunar-zenith project

---

## 1. 曆法系統基礎 (Calendar System Fundamentals)

### 1.1 農曆本質

農曆是**陰陽合曆** (Lunisolar Calendar)，非純陰曆：
- **月相週期 (朔望月)**: 29.5306 天 - 決定月份長度 (29或30天)
- **回歸年 (太陽年)**: 365.2422 天 - 決定閏月插入
- **朔 (New Moon)**: 日月黃經差為 0°，農曆月初一
- **望 (Full Moon)**: 日月黃經差為 180°，農曆月十五

### 1.2 閏月規則 (Leap Month Rules)

**核心規則**: 無中氣之月為閏月

```
1. 冬至必定在農曆 11 月
2. 包含 13 個月的歲 (sui，冬至到冬至) 必有閏月
3. 閏月是第一個沒有中氣的月份
4. 中氣 (Major Solar Term): 雨水、春分、穀雨、小滿、夏至、大暑、處暑、秋分、霜降、小雪、冬至、大寒
```

**範例**: 若某月份只有「清明」(節氣) 沒有「穀雨」(中氣)，則此月為閏月。

### 1.3 農曆新年規則

- 傳統: 第二個冬至後的新月
- 例外: 若閏月出現在11或12月後，可能會有變化 (如 2033 年問題)

---

## 2. 天文演算 (Astronomical Calculations)

### 2.1 儒略日 (Julian Day)

**定義**: 從公元前 4713 年 1 月 1 日正午 (UT) 開始計算的連續天數

**計算公式** (Jean Meeus, Astronomical Algorithms Ch. 2):
```
JD = 365.25*(y+4716) + 30.6001*(m+1) + d + b - 1524.5
其中:
  - 若月 <= 2，年-1，月+12
  - a = floor(y/100)
  - b = 2 - a + floor(a/4) (格里曆修正)
```

**代碼位置**: `pkg/celestial/base.go:TimeToJD()`

**重要基準**:
- JD 2451545.0 = 2000-01-01 12:00 TT
- JD 0 = 公元前 4713 年 1 月 1 日正午

### 2.2 Delta-T (TT - UT)

**定義**: 地球時 (Terrestrial Time) 與 世界時 (Universal Time) 的差異

**為何重要**: 
- 天文計算使用 TT (均勻時間)
- 民用時間使用 UT (受地球自轉減速影響)
- 忽略 Delta-T 會導致歷史日期誤差數小時

**估計值**:
- 2000 年: ~64 秒
- 2024 年: ~69 秒
- 線性增長: ~0.5-1 秒/年

**代碼位置**: `pkg/celestial/base.go:EstimateDeltaT()`

**公式**: ΔT = 62.92 + 0.32217*(y-2000) + 0.005589*(y-2000)² (NASA/Espenak, 2005-2050)

### 2.3 太陽位置計算 (VSOP87)

**方法**: VSOP87 (Variations Séculaires des Orbites Planétaires) 精簡版

**計算步驟**:
1. 計算儒略世紀數: T = (JDE - 2451545.0) / 36525
2. 太陽平黃經 L0 = 280.46646 + 36000.76983*T + 0.0003032*T²
3. 太陽平近點角 M = 357.52911 + 35999.05029*T
4. 中心差 C = (1.9146-0.0048*T)*sin(M) + 0.019993*sin(2M) + 0.000289*sin(3M)
5. 太陽真黃經 = L0 + C
6. 章動與光行差修正

**代碼位置**: `pkg/celestial/solar.go:SolarLongitude()`

### 2.4 月球位置計算 (ELP-2000)

**方法**: ELP-2000/82 (Éphéméride Lunaire Parisienne) 精簡版

**主要週期項**:
- 月球平黃經
- 月球平近點角
- 太陽平近點角
- 月球平緯角
- 日月黃經差

**代碼位置**: `pkg/celestial/moon.go:MoonLongitude()`

### 2.5 二十四節氣 (24 Solar Terms)

**定義**: 太陽黃經每 15° 為一個節氣

| 索引 | 名稱 | 黃經 | 類型 |
|------|------|------|------|
| 0 | 春分 | 0° | 中氣 |
| 1 | 清明 | 15° | 節氣 |
| 2 | 穀雨 | 30° | 中氣 |
| ... | ... | ... | ... |
| 18 | 冬至 | 270° | 中氣 |
| 19 | 小寒 | 285° | 節氣 |
| 20 | 大寒 | 300° | 中氣 |

**區分**:
- **節氣** (Jiéqì): 奇數索引 (立春、驚蟄、清明...)
- **中氣** (Zhōngqì): 偶數索引 (雨水、春分、穀雨...)

**代碼位置**: `pkg/celestial/solar.go`

---

## 3. 干支系統 (Sexagenary Cycle / Ganzhi)

### 3.1 十天干 (Heavenly Stems)

| 索引 | 天干 | 拼音 | 陰陽 | 五行 |
|------|------|------|------|------|
| 0 | 甲 | Jiǎ | 陽 | 木 |
| 1 | 乙 | Yǐ | 陰 | 木 |
| 2 | 丙 | Bǐng | 陽 | 火 |
| 3 | 丁 | Dīng | 陰 | 火 |
| 4 | 戊 | Wù | 陽 | 土 |
| 5 | 己 | Jǐ | 陰 | 土 |
| 6 | 庚 | Gēng | 陽 | 金 |
| 7 | 辛 | Xīn | 陰 | 金 |
| 8 | 壬 | Rén | 陽 | 水 |
| 9 | 癸 | Guǐ | 陰 | 水 |

**代碼**: `pkg/zodiac/sexagenary.go:HeavenlyStems`

### 3.2 十二地支 (Earthly Branches)

| 索引 | 地支 | 拼音 | 生肖 | 時辰 | 陰陽 | 五行 |
|------|------|------|------|------|------|------|
| 0 | 子 | Zǐ | 鼠 | 23:00-01:00 | 陽 | 水 |
| 1 | 丑 | Chǒu | 牛 | 01:00-03:00 | 陰 | 土 |
| 2 | 寅 | Yín | 虎 | 03:00-05:00 | 陽 | 木 |
| 3 | 卯 | Mǎo | 兔 | 05:00-07:00 | 陰 | 木 |
| 4 | 辰 | Chén | 龍 | 07:00-09:00 | 陽 | 土 |
| 5 | 巳 | Sì | 蛇 | 09:00-11:00 | 陰 | 火 |
| 6 | 午 | Wǔ | 馬 | 11:00-13:00 | 陽 | 火 |
| 7 | 未 | Wèi | 羊 | 13:00-15:00 | 陰 | 土 |
| 8 | 申 | Shēn | 猴 | 15:00-17:00 | 陽 | 金 |
| 9 | 酉 | Yǒu | 雞 | 17:00-19:00 | 陰 | 金 |
| 10 | 戌 | Xū | 狗 | 19:00-21:00 | 陽 | 土 |
| 11 | 亥 | Hài | 豬 | 21:00-23:00 | 陰 | 水 |

**代碼**: `pkg/zodiac/sexagenary.go:EarthlyBranches`

### 3.3 年柱計算

**公式**: `(year - 4) mod 60`

**基準**: 西元 4 年為甲子年 (Stem=0, Branch=0)

**代碼**: `pkg/zodiac/sexagenary.go:NewYearSexagenary()`

```go
// 西元 2024 年
offset = (2024 - 4) % 60 = 2020 % 60 = 40
StemIndex = 40 % 10 = 0 (甲)
BranchIndex = 40 % 12 = 4 (辰)
=> 甲辰年
```

### 3.4 月柱計算 (五虎遁)

**規則**: 
- 農曆正月起於 **寅月** (非公曆 1 月)
- 月干由年干決定 (五虎遁)

**五虎遁口訣**:
```
甲己之年丙作首 (年干 甲/己 → 正月 丙寅)
乙庚之歲戊為頭 (年干 乙/庚 → 正月 戊寅)
丙辛之歲尋庚起 (年干 丙/辛 → 正月 庚寅)
丁壬壬位順行流 (年干 丁/壬 → 正月 壬寅)
戊癸之年何方發 甲寅之上好追求 (年干 戊/癸 → 正月 甲寅)
```

**算法**: `startStem = (yearStem % 5 * 2 + 2) % 10`

**代碼**: `pkg/zodiac/sexagenary.go:GetMonthSexagenary()`

### 3.5 日柱計算

**基準**: JD 2451545.0 (2000-01-01) 為戊午日

**公式**: `offset = (JD + 0.5 + 49) % 60`

**代碼**: `pkg/zodiac/sexagenary.go:GetDaySexagenary()`

### 3.6 時柱計算 (五鼠遁)

**時辰對應**: 子時 = 23:00-01:00 (跨日!)

**計算時支**: `hourBranch = ((hour + 1) / 2) % 12`

**五鼠遁口訣**:
```
甲己還加甲 (日干 甲/己 → 子時 甲子)
乙庚丙作初 (日干 乙/庚 → 子時 丙子)
丙辛從戊起 (日干 丙/辛 → 子時 戊子)
丁壬庚子居 (日干 丁/壬 → 子時 庚子)
戊癸何方發 壬子是真途 (日干 戊/癸 → 子時 壬子)
```

**算法**: `startStem = (dayStem % 5 * 2) % 10`

**代碼**: `pkg/zodiac/sexagenary.go:GetHourSexagenary()`

---

## 4. 神煞系統 (Shensha System)

### 4.1 建除十二神 (Twelve Officers)

| 索引 | 名稱 | 吉凶 | 主要含義 |
|------|------|------|----------|
| 0 | 建 | - | 建立、開始 |
| 1 | 除 | - | 清除、沐浴 |
| 2 | 滿 | - | 圓滿、豐收 |
| 3 | 平 | - | 平常、平安 |
| 4 | 定 | - | 安定、決定 |
| 5 | 執 | - | 執行、捕捉 |
| 6 | 破 | 凶 | 破壞、損耗 |
| 7 | 危 | - | 危險、不安 |
| 8 | 成 | 吉 | 成就、完成 |
| 9 | 收 | - | 收穫、收藏 |
| 10 | 開 | 吉 | 開啟、開始 |
| 11 | 閉 | - | 閉塞、結束 |

**計算法則**: 
- 以月支為「建」
- `offset = (dayBranch - monthBranch) mod 12`

**代碼**: `pkg/zodiac/shensha.go:GetTwelveOfficer()`

### 4.2 年支神煞

#### 驛馬 (Yi Ma / Travel Star)

**口訣**: 
```
申子辰馬在寅
寅午戌馬在申
巳酉丑馬在亥
亥卯未馬在巳
```

**代碼**: `pkg/zodiac/shensha.go:GetYearShenSha()`

#### 桃花 (Peach Blossom / Romance)

**口訣**:
```
亥卯未見子
申子辰見酉
寅午戌見卯
巳酉丑見午
```

### 4.3 二十八星宿 (28 Lunar Mansions)

**四象分布**:
- **東方青龍**: 角、亢、氐、房、心、尾、箕
- **北方玄武**: 斗、牛、女、虛、危、室、壁
- **西方白虎**: 奎、婁、胃、昴、畢、觜、參
- **南方朱雀**: 井、鬼、柳、星、張、翼、軫

**各月起始**:
```
正月起角、二月起奎、三月起胃、四月起畢
五月起參、六月起鬼、七月起張、八月起角
九月起奎、十月起胃、十一月起畢、十二月起參
```

**計算**: `mansionIdx = (monthStart + dayBranch) % 28`

**代碼**: `pkg/zodiac/shensha.go:GetTwentyEightMansion()`

### 4.4 十二值神 (Daily Deities)

| 日支 | 值神 | 吉凶 |
|------|------|------|
| 子 | 青龍 | 吉 |
| 丑 | 明堂 | 吉 |
| 寅 | 天刑 | 凶 |
| 卯 | 朱雀 | 凶 |
| 辰 | 金匱 | 吉 |
| 巳 | 天德 | 吉 |
| 午 | 白虎 | 凶 |
| 未 | 玉堂 | 吉 |
| 申 | 天牢 | 凶 |
| 酉 | 玄武 | 凶 |
| 戌 | 司命 | 吉 |
| 亥 | 勾陳 | 凶 |

**代碼**: `pkg/zodiac/shensha.go:GetDailyDeity()`

### 4.5 胎神 (Fetal God)

**計算依據**: 日干

| 日干 | 胎神位置 | 禁忌 |
|------|----------|------|
| 甲 | 門外東南 | 忌修門、動土 |
| 乙 | 碓磨廁外東南 | 忌移動碓磨 |
| 丙 | 廚灶爐外正南 | 忌修廚灶 |
| 丁 | 倉庫廁外東南 | 忌修倉庫 |
| 戊 | 房床廁外正南 | 忌移動房床 |
| 己 | 占門床外正南 | 忌修門、移動床 |
| 庚 | 碓磨廁外正南 | 忌移動碓磨 |
| 辛 | 廚灶廁外西南 | 忌修廚灶 |
| 壬 | 倉庫雞棲外東南 | 忌修倉庫 |
| 癸 | 占門床外東南 | 忌修門、移動床 |

**代碼**: `pkg/zodiac/shensha.go:GetFetalGod()`

### 4.6 沖煞 (Clash & Sha)

**地支相沖**:
```
子午相沖、丑未相沖、寅申相沖
卯酉相沖、辰戌相沖、巳亥相沖
```

**煞方**: 被沖地支的方位

**代碼**: `pkg/zodiac/shensha.go:GetClashSha()`

---

## 5. 黃曆宜忌 (Almanac Auspicious/Inauspicious)

### 5.1 宜忌計算邏輯

**主要依據**:
1. **建除十二神** (最重要)
2. **日干** (輔助判斷)
3. **值神** (吉凶參考)

### 5.2 建除十二神對應宜忌

| 建除 | 宜 | 忌 |
|------|----|----|
| 建 | 出行、上任、會友、上樑 | 嫁娶、開倉、安葬 |
| 除 | 沐浴、清潔、療病、出行 | 嫁娶、安床、開市 |
| 滿 | 嫁娶、祈福、移徙、入宅、開市 | 安葬、出行、求醫 |
| 平 | 修造、動土、安床、修飾垣墻 | 嫁娶、開市、出行 |
| 定 | 嫁娶、祭祀、祈福、求嗣 | 出行、上任、交易 |
| 執 | 祭祀、祈福、捕捉 | 嫁娶、安床、入宅 |
| 破 | 破屋、壞垣、祛病、解除 | 嫁娶、開市、立約、入宅 |
| 危 | 安床、修飾垣墻、拆卸 | 嫁娶、出行、開市、交易 |
| 成 | 嫁娶、祭祀、祈福、開市、交易、立約 | 安葬、動土 |
| 收 | 祭祀、祈福、嫁娶、修造、捕捉 | 開倉、出財、安葬、出行 |
| 開 | 開市、交易、立約、安機械、出行 | 安葬、動土、嫁娶 |
| 閉 | 祭祀、祈福、安葬、修墳 | 開市、出行、嫁娶、入宅 |

**代碼**: `internal/service/almanac.go:AlmanacActivities`

### 5.3 吉神方位

**計算依據**: 日干

| 日干 | 財神 | 喜神 | 福神 | 文曲 |
|------|------|------|------|------|
| 甲 | 東 | 東北 | 東南 | 北 |
| 乙 | 東南 | 西北 | 東 | 東南 |
| 丙 | 南 | 南 | 西 | 東 |
| 丁 | 東南 | 東南 | 西南 | 西南 |
| 戊 | 東 | 東北 | 東南 | 北 |
| 己 | 東北 | 西北 | 東 | 東南 |
| 庚 | 西 | 南 | 西 | 東 |
| 辛 | 西南 | 東南 | 西南 | 西南 |
| 壬 | 北 | 東北 | 東南 | 北 |
| 癸 | 西北 | 西北 | 東 | 東南 |

**代碼**: `internal/service/almanac.go:GetDeityDirections()`

---

## 6. 宗教曆法 (Religious Calendars)

### 6.1 佛曆 (Buddhist Calendar)

**計算**: 西元年 + 544

**說明**: 佛教紀年從釋迦牟尼佛涅槃年開始

**代碼**: `pkg/zodiac/religious.go`

### 6.2 道曆 (Taoist Calendar)

**計算**: 西元年 + 2697

**說明**: 道教紀年從黃帝即位年開始

**代碼**: `pkg/zodiac/religious.go`

---

## 6.5 行政節日與區域差異 (TW/CN Observance Rules)

### 6.5.1 核心原則

1. **節日/紀念日顯示優先**：只要是節日或紀念日，`holiday_info.name` 應顯示名稱，不以是否放假作為顯示條件。  
2. **放假狀態獨立**：`is_holiday` 僅表示是否休假。  
3. **工作日與週末預設不顯示名稱**：補班、一般工作日、週末預設名稱為空字串。  
4. **大陸規則採 fallback 覆寫**：先走大陸差異規則，再 fallback 至台灣共通規則。

### 6.5.2 台灣規則

- 固定日期節日/紀念日：如 `4/4 兒童節`、`8/1 原住民族日`、`8/8 父親節`、`9/28 教師節`。
- 動態規則節日：如母親節（五月第二個週日）、祖父母節（八月第四個週日）、清明節（4/4 或 4/5）、農曆節日（春節、端午、中秋、除夕等）。
- 同日多節日：以 `、` 合併顯示，例如 `國父逝世紀念日、植樹節`。

### 6.5.3 大陸覆寫規則（相對台灣）

- `6/1` 為兒童節（覆寫台灣 `4/4` 兒童節顯示）。
- `8/1` 為建軍節（覆寫台灣 `8/1` 原住民族日）。
- 父親節為**六月第三個週日**（覆寫台灣 `8/8`）。
- 教師節為 `9/10`（覆寫台灣 `9/28`）。
- 不應顯示「一般工作日」等工作日名稱。

### 6.5.4 代碼落點

- `internal/service/holiday.go`
  - `IsHoliday()`：入口，處理 JSON 覆寫、區域規則與 fallback。
  - `getTaiwanObservances()`：台灣節日/紀念日規則。
  - `getChinaObservanceOverrides()`：大陸差異規則。
- `internal/service/holiday_test.go`
  - 應覆蓋：非放假節日顯示、多節日同日、TW/CN 覆寫差異、工作日與週末名稱清理。

---

## 6.6 西洋占星：順逆行與相位/交匯

### 6.6.1 順行/逆行判定

- 依行星黃經瞬時速度 `speed` 判斷：`speed < 0` 視為逆行。
- API 輸出核心欄位：
  - `planet`、`name_zh`、`symbol`
  - `is_retrograde`
  - `longitude`、`speed`
  - 可選：`next_station_date`、`station_type`

**代碼**:
- `pkg/western_astro/retrograde.go:GetRetrogradeInfo()`
- `pkg/western_astro/retrograde.go:GetAllRetrogradeInfo()`

### 6.6.2 相位/交匯計算

- 支援主要相位角：`0°` 合相、`60°` 六合、`90°` 刑克、`120°` 三合、`150°` 梅花、`180°` 對沖。
- 預設容許誤差（orb）`8°`；重大交匯可使用更嚴格 orb。
- 相位輸出欄位：`planet1/planet2`、中文名、符號、`aspect`、`angle`、`orb`。

**代碼**:
- `pkg/western_astro/aspects.go:CalculateAspects()`
- `pkg/western_astro/aspects.go:GetMajorConjunctions()`

---

## 6.7 契約同步技能（REST + gRPC）

### 6.7.1 Contract-First 準則

1. 先更新契約：`contracts/openapi/lunar-zenith.yaml`、`api/v1/lunar.proto`。  
2. 再同步服務層：`internal/service/rest_handler.go`、`internal/service/grpc_server.go`。  
3. 最後驗證：`go test ./...`、`make verify-contracts`。

### 6.7.2 gRPC 同步實務

- 修改 `api/v1/lunar.proto` 後，重新生成：

```bash
PATH="$PATH:$(go env GOPATH)/bin" protoc \
  --proto_path=api/v1 \
  --go_out=api/v1 --go_opt=paths=source_relative \
  --go-grpc_out=api/v1 --go-grpc_opt=paths=source_relative \
  api/v1/lunar.proto
```

- 避免手改 `api/v1/*.pb.go`，以生成結果為準。

### 6.7.3 REST/gRPC 一致性檢查

- 檢查 OpenAPI path 與實際路由一致（本專案為 `/v1/calendar`）。
- 檢查 REST JSON 與 gRPC message 欄位集一致（如 `holiday_info`、`china_holiday_info`、`western_astro`、`aspects`）。
- 新增欄位後，需同步更新 README / SKILLS 等維運文檔。

---

## 7. 代碼架構總覽 (Code Architecture)

### 7.1 模組劃分

```
┌─────────────────────────────────────────────────────────┐
│  API Layer (gRPC + REST)                                 │
│  internal/service/grpc_server.go                         │
│  internal/service/rest_handler.go                        │
├─────────────────────────────────────────────────────────┤
│  Service Layer (Aggregator)                              │
│  internal/service/aggregator.go - 組合所有數據            │
│  internal/service/almanac.go - 黃曆宜忌計算               │
│  internal/service/holiday.go - 假期數據                   │
├─────────────────────────────────────────────────────────┤
│  Zodiac Module (Cultural Calendar)                       │
│  pkg/zodiac/sexagenary.go - 干支計算                     │
│  pkg/zodiac/shensha.go - 神煞系統                        │
│  pkg/zodiac/lunar_engine.go - 農曆日期                   │
│  pkg/zodiac/religious.go - 宗教曆法                      │
│  pkg/zodiac/festival.go - 農曆節日                       │
├─────────────────────────────────────────────────────────┤
│  Celestial Module (Astronomical)                         │
│  pkg/celestial/base.go - JD, Delta-T                     │
│  pkg/celestial/solar.go - 太陽位置、節氣                  │
│  pkg/celestial/moon.go - 月球位置、朔望                   │
└─────────────────────────────────────────────────────────┘
```

### 7.2 關鍵類型與函數

**celestial 模組**:
```go
// 高精度時間封裝
PrecisionTime {
    UT: time.Time    // 世界時
    JD: float64      // 儒略日
    JDE: float64     // 曆書儒略日
    DeltaT: float64  // TT-UT 差值
}

// 核心函數
NewPrecisionTime(t) - 創建高精度時間
TimeToJD(t) - 時間轉 JD
SolarLongitude(jde) - 太陽黃經
MoonLongitude(jde) - 月球黃經
MoonPhase(jde) - 日月黃經差
FindNewMoon(jd, dir) - 尋找朔日
GetSolarTerm(jde) - 獲取節氣
```

**zodiac 模組**:
```go
// 干支結構
Sexagenary {
    StemIndex: int   // 0-9
    BranchIndex: int // 0-11
}

// 核心函數
NewYearSexagenary(year) - 年柱
GetMonthSexagenary(yearStem, month) - 月柱
GetDaySexagenary(jd) - 日柱
GetHourSexagenary(dayStem, hourBranch) - 時柱
GetHourBranch(hour) - 時支計算
GetTwelveOfficer(monthBranch, dayBranch) - 建除十二神
GetYearShenSha(yearBranch) - 年支神煞
GetTwentyEightMansion(month, dayStem, dayBranch) - 二十八星宿
GetDailyDeity(dayBranch) - 值神
GetFetalGod(dayStem) - 胎神
GetClashSha(dayBranch) - 沖煞
```

**service 模組**:
```go
// 聚合器
Aggregator {
    HolidaySvc: *HolidayService
    LunarEng: *LunarEngine
}

// 核心函數
NewAggregator(holidaySvc) - 創建聚合器
GetCalendar(t) - 獲取完整曆法數據
CalculateAlmanac(officer, dayStem) - 計算宜忌
GetDeityDirections(dayStem) - 吉神方位
```

---

## 8. 已知限制與注意事項 (Known Limitations)

### 8.1 閏月計算

**問題**: `lunar_engine.go:59` 閏月判定尚未完全實現

**影響**:
- 2024 無閏月，計算正確
- 2023 (閏二月)、2025 (閏六月) 等有閏月年份可能出錯

**解決方向**:
實現「無中氣月」判定邏輯：
1. 計算每個朔望月的兩個節氣
2. 若月份只含節氣不含中氣，則為閏月
3. 閏月使用前一月名稱加「閏」字

### 8.2 Delta-T 精度

**問題**: 當前僅支援 2000-2100 年

**影響**: 歷史日期或遠期日期誤差較大

**解決方向**: 
- 使用表驅動方式載入 NASA 歷史數據
- 或實現更精確的多項式擬合

### 8.3 時區處理

**問題**: 系統固定使用 UTC+8 (東八區)

**影響**: 其他時區用戶可能看到錯誤日期

**解決方向**: 支援時區參數輸入

### 8.4 曆法改革

**注意**: 中國曆法在 1645 年改革 (使用定氣)
- 1645 年前使用平氣 (平均太陽運動)
- 1645 年後使用定氣 (真實太陽位置)

**當前實現**: 僅支援定氣，歷史日期 (1645年前) 可能有誤差

---

## 9. 修改指南 (Modification Guidelines)

### 9.1 修改天文計算

**正確做法**:
1. 確認修改基於 Jean Meeus 算法或 VSOP87/ELP-2000
2. 保持 `float64` 精度
3. 區分 JD 與 JDE (加上 Delta-T)
4. 更新對應測試

**禁止**:
- 使用近似公式替代精確算法
- 忽略 Delta-T 修正
- 修改 JD 基準常量

### 9.2 修改干支計算

**正確做法**:
1. 保持現有基準 (4年甲子、2000-01-01戊午)
2. 確認五虎遁/五鼠遁算法
3. 考慮農曆月份 (寅月起算)

**禁止**:
- 修改基準常量
- 直接使用公曆月份計算月柱
- 忽略子時跨日處理

### 9.3 修改神煞系統

**正確做法**:
1. 參考傳統曆書驗證規則
2. 保持建除十二神計算邏輯
3. 新增神煞需說明計算依據

**注意**: 神煞規則在不同地區/流派可能有差異

### 9.4 修改宜忌規則

**正確做法**:
1. 基於建除十二神為主
2. 參考傳統通書 (如：協紀辨方書)
3. 保持數據驅動 (map 結構)

**注意**: 不同通書宜忌可能有差異，需說明採用標準

---

## 10. 驗證檢查清單 (Validation Checklist)

在進行任何修改前，確認：

- [ ] 理解黃曆基本規則 (陰陽合曆)
- [ ] 確認天文計算使用 JDE (非 JD)
- [ ] 確認干支計算使用正確基準
- [ ] 確認農曆月份起於寅月
- [ ] 確認閏月處理邏輯 (如適用)
- [ ] 驗證節氣計算 (VSOP87)
- [ ] 驗證朔日計算 (ELP-2000)
- [ ] 確認神煞計算規則
- [ ] 更新相關測試用例
- [ ] 檢查與現有 AGENTS.md 的一致性

---

## 11. 參考資源 (References)

### 11.1 主要參考書籍

1. **Jean Meeus - "Astronomical Algorithms" (2nd Edition)**
   - 第 2 章: Julian Day
   - 第 7 章: Julian Day to Calendar Date
   - 第 22 章: Nutation and the Obliquity of the Ecliptic
   - 第 25 章: Solar Coordinates
   - 第 47 章: Position of the Moon

2. **P. Bretagnon & G. Francou - "VSOP87 Planetary Solutions"**
   - 用於太陽黃經計算

3. **M. Chapront-Touzé & J. Chapront - "ELP 2000-82"**
   - 用於月球位置計算

### 11.2 線上資源

- [香港天文台 - 二十四節氣](https://www.hko.gov.hk/tc/gts/time/24solarterms.htm)
- [Chinese Calendar Rules - ytliu0](https://ytliu0.github.io/ChineseCalendar/rules.html)
- [NASA - Delta T](https://eclipse.gsfc.nasa.gov/SEcat5/deltat.html)

### 11.3 傳統曆書

- 協紀辨方書 (清代)
- 星曆考原 (清代)
- 台灣中央氣象署農民曆

---

**End of Skills Map**

*本文件由 AI 助手維護，任何修改應確保知識準確性與代碼一致性。*
