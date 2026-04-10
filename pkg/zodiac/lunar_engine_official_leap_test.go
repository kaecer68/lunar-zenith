package zodiac

import (
	"testing"
	"time"

	"github.com/kaecer68/lunar-zenith/pkg/celestial"
)

// TestLunarEngine_OfficialLeapMonthSpotChecks
// 使用官方閏月表做抽查驗證：
// - 對每個指定農曆年，掃描該公曆年每日（台北中午）
// - 取得第一個 IsLeap=true 的農曆月份，並比對官方閏月
//
// 這個測試僅作為演算法驗證，不參與月序計算邏輯。
func TestLunarEngine_OfficialLeapMonthSpotChecks(t *testing.T) {
	engine := &LunarEngine{}
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		t.Fatalf("load timezone failed: %v", err)
	}

	cases := []struct {
		year          int
		wantLeapMonth int
	}{
		{year: 1984, wantLeapMonth: 10},
		{year: 2001, wantLeapMonth: 4},
		{year: 2014, wantLeapMonth: 9},
		{year: 2020, wantLeapMonth: 4},
		{year: 2025, wantLeapMonth: 6},
		{year: 2033, wantLeapMonth: 11},
	}

	for _, tt := range cases {
		t.Run(time.Date(tt.year, 1, 1, 0, 0, 0, 0, time.UTC).Format("2006"), func(t *testing.T) {
			got := detectLeapMonthByLunarYear(engine, tt.year, loc)
			if got != tt.wantLeapMonth {
				t.Errorf("year %d: got leap month=%d, want=%d", tt.year, got, tt.wantLeapMonth)
			}
		})
	}
}

func detectLeapMonthByLunarYear(engine *LunarEngine, lunarYear int, loc *time.Location) int {
	start := time.Date(lunarYear, 1, 1, 12, 0, 0, 0, loc)
	end := time.Date(lunarYear+1, 3, 1, 12, 0, 0, 0, loc)
	for tm := start; tm.Before(end); tm = tm.AddDate(0, 0, 1) {
		jd := celestial.TimeToJD(tm.UTC())
		lunar := engine.GetLunarDate(jd)
		if lunar.Year == lunarYear && lunar.IsLeap {
			return lunar.Month
		}
	}
	return 0
}
