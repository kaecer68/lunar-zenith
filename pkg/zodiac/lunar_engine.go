package zodiac

import (
	"math"
	"time"

	"github.com/kaecer68/lunar-zenith/pkg/celestial"
)

// LunarEngine 完成高精度定閏演算
type LunarEngine struct{}

// monthHasZhongqi 判斷從 nmStart 到 nmEnd（下一朔）之間是否包含至少一個中氣。
// 中氣是太陽黃經為 30 的倍數的節氣（雨水=330°、春分=0°、穀雨=30°…冬至=270°、大寒=300°）。
// 即黃經 0, 30, 60, 90 … 330，共 12 個。
func monthHasZhongqi(nmStart, nmEnd float64) bool {
	lonStart := celestial.SolarLongitude(nmStart)
	lonEnd := celestial.SolarLongitude(nmEnd)

	// 將黃經量化到中氣格（每 30 度一格）
	gridStart := math.Floor(lonStart/30.0) * 30.0
	gridEnd := math.Floor(lonEnd/30.0) * 30.0

	// 跨越 360/0 邊界：lonEnd < lonStart
	if lonEnd < lonStart {
		// 至少跨越了一個中氣格
		return true
	}
	// 同一格或進入了下一格，代表中間有中氣
	return gridEnd > gridStart
}

// buildLeapIndex 根據冬至朔日 nmWS 及當年冬至 ws，計算緊接的一個「冬至年」（13 個月）中
// 哪個月索引（0-based，0=11月）是閏月；若該年無閏月則回傳 -1。
//
// 規則（Jean Meeus / 香港天文台標準）：
//  1. 冬至到次冬至若包含 13 個朔望月，則其中第一個「無中氣」的月為閏月。
//  2. 閏月放在前一個正常月的後面，共用前一個月的月序，以 IsLeap=true 標記。
func buildLeapIndex(nmWS, ws float64) (leapPos int, months []float64) {
	// 從 ws+20 開始尋找「下一個」冬至（此時太陽已過 284°，安全遠離 265–275° 帶）
	nextWS := celestial.FindNextWinterSolstice(ws + 20)

	// 收集此冬至週期所有朔日
	months = make([]float64, 0, 15)
	curr := nmWS
	for curr < nextWS+1 {
		months = append(months, curr)
		curr = celestial.FindNewMoon(curr+15, 1)
	}

	// 若朔日數量 <= 12，本年無閏月
	if len(months) <= 12 {
		return -1, months
	}

	// 找第一個「無中氣」的月（index 1 之後，跳過 11 月本身）
	for i := 1; i < len(months)-1; i++ {
		if !monthHasZhongqi(months[i], months[i+1]) {
			return i, months
		}
	}
	return -1, months
}

// GetLunarDate 根據 JD (UT) 獲取精確農曆日期
func (e *LunarEngine) GetLunarDate(jd float64) LunarDate {
	// 根據 JD 推算目標年份，取對應的 Delta-T（秒 → 日）做 JDE 修正
	year, month, _ := celestial.JDToDate(jd)
	deltaTSeconds := celestial.EstimateDeltaT(time.Date(year, time.Month(month), 15, 0, 0, 0, 0, time.UTC))
	jde := jd + deltaTSeconds/86400.0

	// 1. 朔日（當月初一）與日序
	// 使用 PreviousNewMoon（Meeus 均值 + 窄窗口二分）取代 FindNewMoon(jde,-1)，
	// 避免在「jde 剛好在朔日後幾小時」時二分法收斂到前一個朔日的臨界問題。
	nm0 := celestial.PreviousNewMoon(jde)
	dayIdx := int(math.Floor(jd+0.5+8.0/24.0)) - int(math.Floor(nm0+0.5+8.0/24.0)) + 1

	// 2. 定位「含冬至」那個朔日 nmWS
	ws := celestial.FindPreviousWinterSolstice(jde + 32)
	nmWS := celestial.PreviousNewMoon(ws)
	if nmWS > nm0+0.01 {
		ws = celestial.FindPreviousWinterSolstice(nmWS - 2)
		nmWS = celestial.FindNewMoon(ws, -1)
	}

	// 3. 建立本冬至年的閏月資訊，以及所有朔日序列
	leapPos, allMonths := buildLeapIndex(nmWS, ws)

	// 4. 找 nm0 在 allMonths 中的位置
	pos := 0
	for i, nm := range allMonths {
		if math.Abs(nm-nm0) < 0.5 {
			pos = i
			break
		}
	}

	// 5. 計算農曆月份與是否閏月
	//    allMonths[0] = 農曆 11 月朔日
	//    leapPos > 0 代表 allMonths[leapPos] 是閏月
	isLeap := false
	calcMonth := 0

	if leapPos > 0 && pos == leapPos {
		// 當前月就是閏月：月序與前一個月相同
		isLeap = true
		// 前一個月的「正常月號」= (11 + leapPos - 1 - 1) % 12 + 1
		// 更安全：用 pos-1 的正常月號
		normalPos := pos - 1
		// normalPos 之前有幾個非閏月
		leapsBefore := 0
		if leapPos > 0 && normalPos >= leapPos {
			leapsBefore = 1
		}
		calcMonth = (11+normalPos-leapsBefore-1)%12 + 1
	} else {
		// 計算 pos 之前有幾個閏月（最多 1 個）
		leapsBefore := 0
		if leapPos > 0 && pos > leapPos {
			leapsBefore = 1
		}
		calcMonth = (11+pos-leapsBefore-1)%12 + 1
	}

	// 6. 農曆年份：以正月初一為切換點
	yWS, _, _ := celestial.JDToDate(ws)
	lunarYear := yWS + 1
	if calcMonth >= 11 {
		lunarYear = yWS
	}

	return LunarDate{
		Year:       lunarYear,
		Month:      calcMonth,
		Day:        dayIdx,
		IsLeap:     isLeap,
		YearPillar: NewYearSexagenary(lunarYear),
	}
}
