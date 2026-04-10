package zodiac

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/kaecer68/lunar-zenith/pkg/celestial"
)

// TestLunarEngine_HistoricalEdgeFixtures 追蹤目前已知的模型邊界日期。
//
// 這些案例先以 skip 保留，避免被遺忘；每次測試會輸出關鍵診斷資料，
// 供後續把 skip 轉為 hard assert 時直接對照。
func TestLunarEngine_HistoricalEdgeFixtures(t *testing.T) {
	engine := &LunarEngine{}

	cases := []struct {
		name      string
		date      time.Time
		note      string
		wantMonth int
		wantDay   int
		wantLeap  bool
		strict    bool
	}{
		{
			name:      "2001-cny-boundary",
			date:      time.Date(2001, 1, 24, 12, 0, 0, 0, time.UTC),
			note:      "已修正為民用日判月起點一致案例，保留 diagnostics 供邊界回歸追蹤",
			wantMonth: 1,
			wantDay:   1,
			wantLeap:  false,
			strict:    true,
		},
		{
			name:      "2020-leap-4th-day1-boundary",
			date:      time.Date(2020, 5, 23, 12, 0, 0, 0, time.UTC),
			note:      "已修正為官方一致案例，保留 diagnostics 供邊界回歸追蹤",
			wantMonth: 4,
			wantDay:   1,
			wantLeap:  true,
			strict:    true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			pt := celestial.NewPrecisionTime(tt.date)
			got := engine.GetLunarDate(pt.JD)

			t.Logf("boundary date=%s => month=%d day=%d leap=%v", tt.date.Format("2006-01-02"), got.Month, got.Day, got.IsLeap)
			t.Logf("diagnostics:\n%s", debugLunarInternals(pt.JD))

			if got.Month != tt.wantMonth || got.Day != tt.wantDay || got.IsLeap != tt.wantLeap {
				if tt.strict {
					t.Errorf("date %s: got month=%d day=%d leap=%v; want month=%d day=%d leap=%v",
						tt.date.Format("2006-01-02"), got.Month, got.Day, got.IsLeap,
						tt.wantMonth, tt.wantDay, tt.wantLeap)
					return
				}
				t.Skipf("TODO: %s；%s", tt.date.Format("2006-01-02"), tt.note)
				return
			}

			if !tt.strict {
				t.Skipf("TODO: %s；%s", tt.date.Format("2006-01-02"), tt.note)
			}
		})
	}
}

func debugLunarInternals(jd float64) string {
	year, month, _ := celestial.JDToDate(jd)
	deltaTSeconds := celestial.EstimateDeltaT(time.Date(year, time.Month(month), 15, 0, 0, 0, 0, time.UTC))
	jde := jd + deltaTSeconds/86400.0

	nm0 := monthStartNewMoonForCivilDay(jd, jde)
	ws := celestial.FindPreviousWinterSolstice(nm0 - 1e-6)
	nmWS := winterSolsticeMonthStart(ws)

	leapPos, allMonths := buildLeapIndex(nmWS, ws)
	pos := 0
	for i, nm := range allMonths {
		if math.Abs(nm-nm0) < 0.5 {
			pos = i
			break
		}
	}

	start := pos - 2
	if start < 0 {
		start = 0
	}
	end := pos + 2
	if end >= len(allMonths) {
		end = len(allMonths) - 1
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("JD=%.6f JDE=%.6f deltaT=%.2fs\n", jd, jde, deltaTSeconds))
	b.WriteString(fmt.Sprintf("nm0=%.6f ws=%.6f nmWS=%.6f\n", nm0, ws, nmWS))
	b.WriteString(fmt.Sprintf("leapPos=%d pos=%d months=%d\n", leapPos, pos, len(allMonths)))
	for i := start; i <= end; i++ {
		nm := allMonths[i]
		lon := celestial.SolarLongitude(nm)
		mark := ""
		if i == pos {
			mark = " <- current"
		}
		if i == leapPos {
			mark += " <- leap"
		}
		b.WriteString(fmt.Sprintf("  [%02d] nm=%.6f lon=%.2f%s\n", i, nm, lon, mark))
	}

	return strings.TrimRight(b.String(), "\n")
}
