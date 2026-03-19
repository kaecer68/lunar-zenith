# pkg/celestial - Astronomical Calculation Module

**Domain**: Celestial mechanics, Julian Day calculations, Delta-T estimation, solar/lunar position

## OVERVIEW

High-precision astronomical calculations based on Jean Meeus's *Astronomical Algorithms*. Provides JD/JDE conversion, Delta-T estimation, and solar/lunar position computations.

## STRUCTURE

```
pkg/celestial/
├── base.go       # PrecisionTime, JD↔Date conversion, Delta-T estimation
├── solar.go      # Solar position, solar terms (節氣)
├── moon.go       # Lunar position, phases, new moon detection
└── *_test.go     # Unit tests for each module
```

## WHERE TO LOOK

| Task | File | Notes |
|------|------|-------|
| JD 轉換 | `base.go:TimeToJD()` | Meeus Algorithm, Ch. 2 |
| JD→日期 | `base.go:JDToDate()` | Meeus Algorithm, Ch. 7 |
| Delta-T 估算 | `base.go:EstimateDeltaT()` | NASA/Espenak 2005-2050 |
| 太陽位置 | `solar.go` | VSOP87 精簡版 |
| 節氣計算 | `solar.go` | 太陽黃經 15°倍數 |
| 月球位置 | `moon.go` | ELP2000 精簡版 |
| 朔望月 | `moon.go` | 新月 (初一) 判定 |

## CODE MAP

| Symbol | Type | Role |
|--------|------|------|
| `PrecisionTime` | struct | 封裝 UT/JD/JDE/DeltaT |
| `NewPrecisionTime(t time.Time)` | constructor | 從 Go time 創建高精度時間 |
| `TimeToJD(t time.Time) float64` | func | time → 儒略日 |
| `JDToDate(jd float64) (y,m,d)` | func | 儒略日 → 日期 |
| `EstimateDeltaT(t time.Time) float64` | func | TT - UT (秒) |

## CONVENTIONS

- **UT vs TT**: 
  - UT (Universal Time): 民用時，`time.Time` 輸入
  - TT (Terrestrial Time): 曆書時，JDE = JD + DeltaT/86400
  - DeltaT: TT - UT (秒), 2024 年約 69s

- **JD Convention**: JD 0.5 起算 (正午), `JDToDate()` 自動 +0.5 校正

- **Precision**: `float64` for all JD calculations (15-16 significant digits)

## ANTI-PATTERNS

- ❌ 使用 UT 直接計算太陽/月球位置 (必須用 JDE)
- ❌ 忽略 Delta-T 修正 (歷史日期誤差可達數小時)
- ❌ 修改 JD 基準常量 (JD 2451545.0 = 2000-01-01 12:00 TT)

## DEPENDENCIES

- **External**: None (pure Go)
- **Internal**: Used by `pkg/zodiac/` for lunar conversion

## TESTING

```bash
go test ./pkg/celestial/... -v
```

- Stateless, thread-safe, no benchmarks yet
