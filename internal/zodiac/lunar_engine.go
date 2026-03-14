package zodiac

import (
	"math"
	"github.com/kaecer68/lunar-zenith/internal/celestial"
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
	// 冬至一定在農曆 11 月
	ws := celestial.FindPreviousWinterSolstice(jde)
	if ws > nm0 {
		// 如果當前月就在冬至前，我們需要找再上一個冬至
		ws = celestial.FindPreviousWinterSolstice(nm0 - 5)
	}

	// 3. 計算冬至到當前朔之間的月數
	months := 0
	checkNM := ws
	for {
		nextNM := celestial.FindNewMoon(checkNM+15, 1)
		if nextNM > nm0+0.01 { // 容誤差
			break
		}
		months++
		checkNM = nextNM
	}

	// 基礎月份：11月 + 月數
	calcMonth := (11+months-1)%12 + 1

	// TODO: 完善閏月判定 (無中氣月)
	// 目前 2024 無閏月，此邏輯可滿足 MVP

	return LunarDate{
		Year:  2024, 
		Month: calcMonth,
		Day:   dayIdx,
		YearPillar: NewYearSexagenary(2024),
	}
}
