package zodiac

import (
	"math"
	"time"

	"github.com/kaecer68/lunar-zenith/v4/pkg/celestial"
)

// LunarEngine 完成高精度定閏演算
type LunarEngine struct{}

// monthHasZhongqi 判斷從 nmStart 到 nmEnd（下一朔）之間是否包含至少一個中氣。
// 中氣是太陽黃經為 30 的倍數的節氣（雨水=330°、春分=0°、穀雨=30°…冬至=270°、大寒=300°）。
// 即黃經 0, 30, 60, 90 … 330，共 12 個。
func monthHasZhongqi(nmStart, nmEnd float64) bool {
	has, _ := monthHasZhongqiWithBoundary(nmStart, nmEnd)
	return has
}

func monthHasZhongqiWithBoundary(nmStart, nmEnd float64) (has bool, onEndDay bool) {
	lonStart := celestial.SolarLongitude(nmStart)
	lonEnd := celestial.SolarLongitude(nmEnd)
	endUnwrapped := lonEnd
	if endUnwrapped < lonStart {
		endUnwrapped += 360.0
	}

	// 在一個朔望月內太陽黃經增量約 29 度，至多跨越一個中氣界線。
	boundary := math.Ceil((lonStart+1e-12)/30.0) * 30.0
	if boundary > endUnwrapped+1e-12 {
		return false, false
	}

	target := math.Mod(boundary, 360.0)
	termJDE := celestial.EstimateTermTime(target, nmStart, nmEnd)

	// 以 UTC+8 民用日為歸屬基準：[朔日, 次朔日) 視為本月。
	// 若中氣與次朔落在同一民用日，則歸到下月，不算本月中氣。
	startDay := civilDayUTC8(nmStart)
	endDay := civilDayUTC8(nmEnd)
	termDay := civilDayUTC8(termJDE)
	if termDay == endDay {
		return false, true
	}

	return termDay >= startDay && termDay < endDay, false
}

func civilDayUTC8(jd float64) int {
	return int(math.Floor(jd + 0.5 + 8.0/24.0))
}

func monthStartNewMoonForCivilDay(jd, jde float64) float64 {
	nmPrev := celestial.PreviousNewMoon(jde)
	nmNext := celestial.FindNewMoon(jde, 1)

	queryDay := civilDayUTC8(jd)
	if civilDayUTC8(nmNext) == queryDay {
		return nmNext
	}
	return nmPrev
}

func winterSolsticeMonthStart(ws float64) float64 {
	nmPrev := celestial.PreviousNewMoon(ws)
	nmNext := celestial.FindNewMoon(ws, 1)

	// 農曆月以 UTC+8 民用日切換：若冬至與下一朔落在同一天，
	// 則該日已屬於下一個月，11 月錨點應取下一朔。
	if civilDayUTC8(nmNext) == civilDayUTC8(ws) {
		return nmNext
	}
	return nmPrev
}

// buildLeapIndex 根據冬至朔日 nmWS 及當年冬至 ws，計算緊接的一個「冬至年」（13 個月）中
// 哪個月索引（0-based，0=11月）是閏月；若該年無閏月則回傳 -1。
//
// 規則（Jean Meeus / 香港天文台標準）：
//  1. 冬至到次冬至若包含 13 個朔望月，則其中第一個「無中氣」的月為閏月。
//  2. 閏月放在前一個正常月的後面，共用前一個月的月序，以 IsLeap=true 標記。
//
// 注意：2033年問題
//  2032-2033 冬至年只有 12 個朔望月，即使出現無中氣月也不置閏。
//  2033-2034 冬至年有 13 個朔望月，閏月為閏十一月（首個無中氣月）。
//  關鍵原則：冬至必須落在十一月。
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

	// 若朔日數量 <= 12，本年無閏月（2032-2033 冬至年即此案例）
	if len(months) <= 12 {
		return -1, months
	}

	monthCount := len(months) - 1
	hasZQ := make([]bool, monthCount)
	onEndDay := make([]bool, monthCount)
	maxEndDayRun := 0
	currentRun := 0
	for i := 1; i < monthCount; i++ {
		has, endDayBoundary := monthHasZhongqiWithBoundary(months[i], months[i+1])
		hasZQ[i] = has
		onEndDay[i] = endDayBoundary
		if !has && endDayBoundary {
			currentRun++
			if currentRun > maxEndDayRun {
				maxEndDayRun = currentRun
			}
		} else {
			currentRun = 0
		}
	}

	// 若同一冬至年出現長鏈（>=3）「中氣落在次朔同一民用日」案例，
	// 視為邊界歸屬連鎖歧義：避免將這些月份當作首個無中氣月。
	// 這可消除 2033-2034 冬至年（閏十一月）的提前假閏判定。
	for i := 1; i < monthCount; i++ {
		if !hasZQ[i] {
			if maxEndDayRun >= 3 && onEndDay[i] {
				continue
			}
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
	nm0 := monthStartNewMoonForCivilDay(jd, jde)
	dayIdx := int(math.Floor(jd+0.5+8.0/24.0)) - int(math.Floor(nm0+0.5+8.0/24.0)) + 1

	// 2. 定位「含冬至」那個朔日 nmWS
	// 鎖定「當前朔月起點之前」最近的冬至，避免同一朔月內因跨越冬至
	// 在不同日期落到不同週期，造成同一年月序跳變（1984/2033 類型）。
	ws := celestial.FindPreviousWinterSolstice(nm0 - 1e-6)
	nmWS := winterSolsticeMonthStart(ws)

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
	if pos >= len(allMonths)-1 {
		pos = len(allMonths) - 2
	}

	// 5. 生成本冬至年的月序（含閏月）與農曆年映射。
	// allMonths[0] 固定為含冬至月（農曆 11 月），其農曆年為 yWS。
	yWS, _, _ := celestial.JDToDate(ws)
	monthCount := len(allMonths) - 1
	monthNums := make([]int, monthCount)
	monthLeaps := make([]bool, monthCount)
	monthYears := make([]int, monthCount)

	monthNums[0] = 11
	monthLeaps[0] = false
	monthYears[0] = yWS

	for i := 1; i < monthCount; i++ {
		if leapPos > 0 && i == leapPos {
			// 閏月：沿用上一月月號，年份不變。
			monthNums[i] = monthNums[i-1]
			monthLeaps[i] = true
			monthYears[i] = monthYears[i-1]
			continue
		}

		nextMonth := monthNums[i-1]%12 + 1
		monthNums[i] = nextMonth
		monthLeaps[i] = false
		monthYears[i] = monthYears[i-1]
		if nextMonth == 1 {
			// 只有進入「非閏正月」才切換農曆年。
			monthYears[i]++
		}
	}

	calcMonth := monthNums[pos]
	isLeap := monthLeaps[pos]
	lunarYear := monthYears[pos]

	return LunarDate{
		Year:       lunarYear,
		Month:      calcMonth,
		Day:        dayIdx,
		IsLeap:     isLeap,
		YearPillar: NewYearSexagenary(lunarYear),
	}
}
