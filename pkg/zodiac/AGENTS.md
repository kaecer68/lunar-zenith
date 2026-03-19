# pkg/zodiac - Cultural Calendar Module

**Domain**: Traditional Chinese calendar calculations, sexagenary cycle, shensha (神煞)

## OVERVIEW

Converts astronomical data (JD, solar terms) into cultural calendar representations: Ganzhi, lunar dates, zodiac animals, religious calendars.

## STRUCTURE

```
pkg/zodiac/
├── sexagenary.go      # 干支計算：年/月/日/時四柱，五虎遁，五鼠遁
├── lunar_engine.go    # 農曆核心：朔望月，閏月判定 (TODO)
├── lunar_date.go      # 農曆日期轉換
├── alignment.go       # 曆法對齊：節氣，星座
├── shensha.go         # 神煞系統：建除十二神，年驛馬，年桃花
├── religious.go       # 宗教曆：佛曆，道曆
└── calendar/          # (reserved for calendar-specific logic)
```

## WHERE TO LOOK

| Task | File | Notes |
|------|------|-------|
| 年干支計算 | `sexagenary.go:NewYearSexagenary()` | 基準：西元 4 年甲子 |
| 日干支計算 | `sexagenary.go:GetDaySexagenary()` | JD 基準，2000-01-01 戊午日 |
| 月干支 (五虎遁) | `sexagenary.go:GetMonthSexagenary()` | 寅月起算 |
| 時干支 (五鼠遁) | `sexagenary.go:GetHourSexagenary()` | 子時 23:00-01:00 |
| 農曆轉換 | `lunar_engine.go` | 閏月判定待完善 |
| 神煞計算 | `shensha.go` | 建除十二神，年支神煞 |
| 宗教年份 | `religious.go` | 佛曆 (Buddhist), 道曆 (Taoist) |

## CONVENTIONS

- **Stem/Branch Index**: 天干 0-9 (甲 - 癸), 地支 0-11 (子 - 亥)
- **Month Base**: 農曆月以寅月 (正月) 為 1，非公曆 1 月
- **Hour Branch**: 子時跨日 (23:00-01:00), `GetHourBranch()` 處理
- **Animal Mapping**: `ZodiacAnimals[]` 直接對應地支索引

## ANTI-PATTERNS

- ❌ 使用公曆月份直接計算月干支 (必須轉為農曆寅月起算)
- ❌ 忽略閏月判定 (目前 `lunar_engine.go:59` TODO，2024 無閏月可用)
- ❌ 修改基準常量 (西元 4 年甲子，2000-01-01 戊午日)

## KEY FUNCTIONS

```go
// 年干支
NewYearSexagenary(year int) Sexagenary

// 日干支 (JD 輸入)
GetDaySexagenary(jd float64) Sexagenary

// 月干支 (五虎遁：年干 + 月份)
GetMonthSexagenary(yearStem int, month int) Sexagenary

// 時干支 (五鼠遁：日干 + 時支)
GetHourSexagenary(dayStem int, hourBranch int) Sexagenary

// 時支 (24 小時制轉地支)
GetHourBranch(hour int) int
```

## KNOWN LIMITATIONS

1. **Leap Month**: `lunar_engine.go:59` - 閏月判定待完善 (無中氣月邏輯)
   - 目前 2024 無閏月，MVP 可用
   - 2023, 2025 等有閏月年份會出錯

2. **Test Coverage**: 邊界條件測試不足 (閏年，歷史日期)

## DEPENDENCIES

- `pkg/celestial/`: JD 計算，Delta-T 修正
- `internal/service/`: 假期數據聚合

## TESTING

```bash
go test ./pkg/zodiac/... -v
```

Test files: `*_test.go` for each module (table-driven tests)
