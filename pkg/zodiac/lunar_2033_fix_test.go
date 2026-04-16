package zodiac

import (
	"math"
	"testing"
	"time"

	"github.com/kaecer68/lunar-zenith/v4/pkg/celestial"
)

// Test2033LeapMonthVerification 驗證 2033 年「無閏七月」
// 這是經典的「2033年問題」：2032-2033 冬至年只有 12 個朔望月，無閏月
// 2033-2034 冬至年有閏十一月（非閏七月）
func Test2033LeapMonthVerification(t *testing.T) {
	engine := &LunarEngine{}

	// 2033 年無閏七月（2032-2033 冬至年只有 12 個朔望月）
	testCases := []struct {
		name      string
		date      time.Time
		wantMonth int
		wantDay   int
		wantLeap  bool
	}{
		{
			name:      "2033-07-26-normal-7th-day1",
			date:      time.Date(2033, 7, 26, 12, 0, 0, 0, time.UTC),
			wantMonth: 7,
			wantDay:   1,
			wantLeap:  false, // 正常七月初一
		},
		{
			name:      "2033-08-24-normal-7th-day30",
			date:      time.Date(2033, 8, 24, 12, 0, 0, 0, time.UTC),
			wantMonth: 7,
			wantDay:   30,
			wantLeap:  false, // 七月三十
		},
		{
			name:      "2033-08-25-normal-8th-day1",
			date:      time.Date(2033, 8, 25, 12, 0, 0, 0, time.UTC),
			wantMonth: 8,
			wantDay:   1,
			wantLeap:  false, // 八月初一（無閏七月！）
		},
		{
			name:      "2033-08-26-normal-8th-day2",
			date:      time.Date(2033, 8, 26, 12, 0, 0, 0, time.UTC),
			wantMonth: 8,
			wantDay:   2,
			wantLeap:  false, // 八月初二
		},
		{
			name:      "2033-09-23-normal-9th-day1",
			date:      time.Date(2033, 9, 23, 12, 0, 0, 0, time.UTC),
			wantMonth: 9,
			wantDay:   1,
			wantLeap:  false, // 九月初一（無閏八月！）
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pt := celestial.NewPrecisionTime(tc.date)
			got := engine.GetLunarDate(pt.JD)

			if got.Month != tc.wantMonth || got.Day != tc.wantDay || got.IsLeap != tc.wantLeap {
				t.Errorf("date %s: got month=%d day=%d leap=%v; want month=%d day=%d leap=%v",
					tc.date.Format("2006-01-02"), got.Month, got.Day, got.IsLeap,
					tc.wantMonth, tc.wantDay, tc.wantLeap)
			}
		})
	}
}

// Test2033LeapMonthDiagnosis 診斷 2033 年的閏月計算
func Test2033LeapMonthDiagnosis(t *testing.T) {
	// 估算 2032 年冬至
	ws2032 := findWinterSolsticeForTest(2032)
	t.Logf("2032 年冬至: JD %.6f", ws2032)

	// 找到含冬至的朔日
	nmWS := findWinterSolsticeMonthStartForTest(ws2032)
	t.Logf("含冬至的朔日: JD %.6f", nmWS)

	// 找到下一個冬至
	ws2033 := findWinterSolsticeForTest(2033)
	t.Logf("2033 年冬至: JD %.6f", ws2033)

	// 收集所有朔日
	months := collectMonthsForTest(nmWS, ws2033)
	t.Logf("朔日數量: %d", len(months))

	for i, m := range months {
		t.Logf("  [%2d] JD %.6f", i, m)
	}

	// 檢查哪個月是閏月
	leapIdx, _ := buildLeapIndex(nmWS, ws2032)
	t.Logf("buildLeapIndex 返回的閏月索引: %d", leapIdx)

	if leapIdx >= 0 && leapIdx < len(months) {
		t.Logf("閏月朔日: JD %.6f", months[leapIdx])
	}
}

// 輔助函數
func findWinterSolsticeForTest(year int) float64 {
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

func findWinterSolsticeMonthStartForTest(ws float64) float64 {
	nmPrev := celestial.PreviousNewMoon(ws)
	nmNext := celestial.FindNewMoon(ws, 1)

	if civilDayUTC8ForTest(nmNext) == civilDayUTC8ForTest(ws) {
		return nmNext
	}
	return nmPrev
}

func civilDayUTC8ForTest(jd float64) int {
	return int(math.Floor(jd + 0.5 + 8.0/24.0))
}

func collectMonthsForTest(nmWS, nextWS float64) []float64 {
	months := make([]float64, 0, 15)
	curr := nmWS
	for curr < nextWS+1 {
		months = append(months, curr)
		curr = celestial.FindNewMoon(curr+15, 1)
	}
	return months
}
