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
| 農曆轉換 | `lunar_engine.go` | 已實作無中氣月定閏，先看邊界測試與 `IsLeap` 輸出 |
| 閏月輸出 | `lunar_date.go` | `LunarDate.IsLeap` 可用；邊界年份請搭配測試驗證 |
| 二十八星宿 | `shensha.go:GetTwentyEightMansion()` | 公曆實際天數口徑；公式：mansionIndex = (A + 24) % 28 |
| 胎神 (60甲子) | `shensha.go:GetFetalGodByDayPillar()` | 60甲子日柱定表；精確映射，推薦新介面 |
| 胎神 (日干) | `shensha.go:GetFetalGod()` | 舊介面，已棄用；僅保留相容層 |
| 建除十二神 | `shensha.go:GetTwelveOfficer()` | 月支決策，依(dayBranch - monthBranch) % 12 |
| 沖煞 | `shensha.go:GetClashSha()` | 日支相沖 + 煞方四正位映射 |
| 宗教年份 | `religious.go` | 佛曆 (Buddhist), 道曆 (Taoist) |

## CONVENTIONS

- **Stem/Branch Index**: 天干 0-9 (甲 - 癸), 地支 0-11 (子 - 亥)
- **Month Base**: 農曆月以寅月 (正月) 為 1，非公曆 1 月
- **Hour Branch**: 子時跨日 (23:00-01:00), `GetHourBranch()` 處理
- **Animal Mapping**: `ZodiacAnimals[]` 直接對應地支索引
- **Correctness Claims**: 若功能涉及農曆月序、`IsLeap`、農曆年切換、節日落點或對外 API 欄位，只有在明確驗證過對應閏月年份後，才可宣稱「正確」或「已支援」

## ANTI-PATTERNS

- ❌ 使用公曆月份直接計算月干支 (必須轉為農曆寅月起算)
- ❌ 忽略閏月與邊界年份驗證（尤其 2001、2020、2033 這類高風險日期）
- ❌ 在未驗證閏月年份前，宣稱農曆日期、節日落點或 `IsLeap` 輸出對所有年份正確
- ❌ 修改基準常量 (西元 4 年甲子，2000-01-01 戊午日)
- ❌ 使用舊方式計算二十八星宿：基於每月起始日期 + 日支 (已廢棄；改用公曆實際天數公式)
- ❌ 爲胎神計算使用日干單獨變量 (必須使用 60甲子日柱，即 dayStem + dayBranch 組合)

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

// 二十八星宿 (公曆實際天數口徑)
GetTwentyEightMansion(year, month, day int) MansionInfo

// 胎神 (60甲子日柱定表，推薦使用)
GetFetalGodByDayPillar(dayStem, dayBranch int) FetalGodInfo

// 胎神 (僅日干層級，已棄用)
GetFetalGod(dayStem int) FetalGodInfo  // Deprecated: use GetFetalGodByDayPillar()

// 建除十二神
GetTwelveOfficer(monthBranch, dayBranch int) string

// 沖煞
GetClashSha(dayBranch int) ClashShaInfo
```

## KNOWN LIMITATIONS

1. **Leap Month Edge Cases**: 定閏主邏輯已上線，但邊界年份仍需嚴格驗證
   - 2023、2025 主風險年案例已納入並通過
   - 2001-01-24、2020-05-23、2033 等邊界案例目前仍應視為高風險；不可憑推測宣稱正確
   - 若修改會影響 `LunarDate.Month`、`LunarDate.IsLeap`、農曆新年、節日映射或任何對外 API 欄位，必須先補對應年份測試或明確標註未支援範圍
   - 在未完成上述驗證前，正確說法應是「目前僅對已驗證年份有信心」，而不是「已支援所有年份農曆/閏月」

2. **Test Coverage**: 邊界條件測試不足 (閏年，歷史日期)

## DEPENDENCIES

- `pkg/celestial/`: JD 計算，Delta-T 修正
- `internal/service/`: 假期數據聚合

## TESTING

```bash
go test ./pkg/zodiac/... -v
```

Test files: `*_test.go` for each module (table-driven tests)

- When touching lunar conversion logic, add or run focused cases that cover at least one no-leap year and one leap-month year.
- If the change affects public responses, verify downstream behavior in `internal/service/` instead of checking `pkg/zodiac/` only.
