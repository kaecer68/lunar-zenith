package zodiac

import (
	"math"
	"testing"
	"time"

	"github.com/kaecer68/lunar-zenith/v4/pkg/celestial"
)

// Test2033DetailedAnalysis 詳細分析 2033 年的閏月計算
func Test2033DetailedAnalysis(t *testing.T) {
	// 估算 2032 年冬至
	ws2032 := findWinterSolsticeForTest2(2032)
	ws2033 := findWinterSolsticeForTest2(2033)

	// 找到含冬至的朔日
	nmWS := findWinterSolsticeMonthStartForTest2(ws2032)

	// 收集所有朔日
	months := collectMonthsForTest2(nmWS, ws2033)

	t.Logf("=== 2033 年朔望月詳細分析 ===")
	t.Logf("冬至年: %.6f 到 %.6f", ws2032, ws2033)
	t.Logf("朔日數量: %d (應為 13 或 14)", len(months))

	monthCount := len(months) - 1
	hasZQ := make([]bool, monthCount)
	onEndDay := make([]bool, monthCount)

	for i := 0; i < monthCount && i < 15; i++ {
		nmStart := months[i]
		nmEnd := months[i+1]

		sLonStart := celestial.SolarLongitude(nmStart)
		sLonEnd := celestial.SolarLongitude(nmEnd)

		has, endDay := monthHasZhongqiWithBoundary2(nmStart, nmEnd)
		hasZQ[i] = has
		onEndDay[i] = endDay

		// 計算中氣日期（如果有）
		zqDate := "N/A"
		if has {
			lonStart := sLonStart
			lonEnd := sLonEnd
			endUnwrapped := lonEnd
			if endUnwrapped < lonStart {
				endUnwrapped += 360.0
			}
			boundary := math.Ceil((lonStart+1e-12)/30.0) * 30.0
			if boundary <= endUnwrapped+1e-12 {
				target := math.Mod(boundary, 360.0)
				termJDE := celestial.EstimateTermTime(target, nmStart, nmEnd)
				zqDate = jdToDate2(termJDE)[:10]
			}
		}

		zqStatus := "無"
		if has {
			zqStatus = "有"
		}
		endDayStatus := ""
		if endDay {
			endDayStatus = "是"
		}

		// 計算農曆月份
		lunarMonth := i
		if i == 0 {
			lunarMonth = 11
		} else if i == 1 {
			lunarMonth = 12
		} else {
			lunarMonth = i - 1
			if lunarMonth > 12 {
				lunarMonth -= 12
			}
		}

		t.Logf("[%2d] %.6f %s %.1f° %.1f° %s %s %s (月序=%d)",
			i, nmStart, jdToDate2(nmStart)[:10], sLonStart, sLonEnd,
			zqStatus, endDayStatus, zqDate, lunarMonth)
	}

	// 統計無中氣月
	noZqCount := 0
	maxRun := 0
	currentRun := 0
	for i := 1; i < monthCount; i++ {
		if !hasZQ[i] {
			noZqCount++
			currentRun++
			if currentRun > maxRun {
				maxRun = currentRun
			}
			status := ""
			if onEndDay[i] {
				status = " (onEndDay)"
			}
			t.Logf("  索引 %d: 無中氣%s", i, status)
		} else {
			currentRun = 0
		}
	}
	t.Logf("總共 %d 個無中氣月，最長連續 %d 個", noZqCount, maxRun)

	// 模擬 buildLeapIndex 邏輯
	leapPos := -1
	for i := 1; i < monthCount; i++ {
		if !hasZQ[i] {
			if maxRun >= 3 && onEndDay[i] {
				t.Logf("  索引 %d: 無中氣但被長鏈規則跳過 (onEndDay=%v)", i, onEndDay[i])
				continue
			}
			leapPos = i
			t.Logf("  索引 %d: 被選為閏月", i)
			break
		}
	}
	if leapPos == -1 {
		t.Log("  未找到閏月 (返回 -1)")
	}
}

func findWinterSolsticeForTest2(year int) float64 {
	t := time.Date(year, 12, 20, 0, 0, 0, 0, time.UTC)
	startJDE := float64(t.Unix())/86400.0 + 2440587.5
	for d := startJDE; d < startJDE+10; d += 0.1 {
		lon := celestial.SolarLongitude(d)
		if lon >= 269.5 && lon <= 270.5 {
			return celestial.EstimateTermTime(270.0, d-1, d+1)
		}
	}
	return startJDE + 1
}

func findWinterSolsticeMonthStartForTest2(ws float64) float64 {
	nmPrev := celestial.PreviousNewMoon(ws)
	nmNext := celestial.FindNewMoon(ws, 1)
	if civilDayUTC82(nmNext) == civilDayUTC82(ws) {
		return nmNext
	}
	return nmPrev
}

func collectMonthsForTest2(nmWS, nextWS float64) []float64 {
	months := make([]float64, 0, 15)
	curr := nmWS
	for curr < nextWS+1 {
		months = append(months, curr)
		curr = celestial.FindNewMoon(curr+15, 1)
	}
	return months
}

func civilDayUTC82(jd float64) int {
	return int(math.Floor(jd + 0.5 + 8.0/24.0))
}

func jdToDate2(jd float64) string {
	unixTime := int64((jd - 2440587.5) * 86400)
	t := time.Unix(unixTime, 0).UTC()
	return t.Format("2006-01-02 15:04")
}

func monthHasZhongqiWithBoundary2(nmStart, nmEnd float64) (has bool, onEndDay bool) {
	lonStart := celestial.SolarLongitude(nmStart)
	lonEnd := celestial.SolarLongitude(nmEnd)
	endUnwrapped := lonEnd
	if endUnwrapped < lonStart {
		endUnwrapped += 360.0
	}

	boundary := math.Ceil((lonStart+1e-12)/30.0) * 30.0
	if boundary > endUnwrapped+1e-12 {
		return false, false
	}

	target := math.Mod(boundary, 360.0)
	termJDE := celestial.EstimateTermTime(target, nmStart, nmEnd)

	startDay := civilDayUTC82(nmStart)
	endDay := civilDayUTC82(nmEnd)
	termDay := civilDayUTC82(termJDE)
	if termDay == endDay {
		return false, true
	}

	return termDay >= startDay && termDay < endDay, false
}
