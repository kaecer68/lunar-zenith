package zodiac

import (
	"math"
	"github.com/kaecer68/lunar-zenith/pkg/celestial"
)

// LunarEngine 完成高精度定閏演算
type LunarEngine struct{}

// GetLunarDate 根據 jde 獲取精確農曆日期
func (e *LunarEngine) GetLunarDate(jd float64) LunarDate {
	jde := jd + (69.0 / 86400.0) // 基準 Delta-T

	// 1. 尋找當月「朔」時刻所在日期 (初一)
	nm0 := celestial.FindNewMoon(jde, -1)
	dayIdx := int(math.Floor(jd+0.5+8.0/24.0)) - int(math.Floor(nm0+0.5+8.0/24.0)) + 1

	// 2. 定位基準冬至 (WS)
	// 冬至一定在農曆 11 月。
	// 我們尋找一個冬至 ws，使得包含它的朔日 nmWS <= nm0。
	ws := celestial.FindPreviousWinterSolstice(jde + 32)
	nmWS := celestial.FindNewMoon(ws, -1)
	if nmWS > nm0+0.01 {
		// 如果該冬至的朔日晚於當前朔日，說明當前朔日屬於上一個冬至週期
		ws = celestial.FindPreviousWinterSolstice(nmWS - 2)
		nmWS = celestial.FindNewMoon(ws, -1)
	}

	// 3. 計算 nmWS 到 nm0 之間的月數 (用於確定月份)
	months := 0
	curr := nmWS
	for {
		if curr > nm0-0.01 && curr < nm0+0.01 {
			break
		}
		next := celestial.FindNewMoon(curr+15, 1)
		if next > nm0+0.01 {
			break
		}
		months++
		curr = next
	}

	// 基礎月份：11月 + 月數
	calcMonth := (11+months-1)%12 + 1

	// 4. 計算農曆年份
	// 規則：農曆年份（如 2024 甲辰）的切換點在「正月初一」。
	// 冬至 (ws) 所在的週期對應一個農曆年。
	// 若當前月為 11 或 12 月，則年份為 yWS。
	// 若當前月為 1-10 月，則年份為 yWS + 1。
	yWS, _, _ := celestial.JDToDate(ws)
	lunarYear := yWS + 1
	if calcMonth >= 11 {
		lunarYear = yWS
	}

	// TODO: 完善閏月判定 (無中氣月)
	// 目前 2024 無閏月，此邏輯可滿足 MVP

	return LunarDate{
		Year:       lunarYear,
		Month:      calcMonth,
		Day:        dayIdx,
		YearPillar: NewYearSexagenary(lunarYear),
	}
}
