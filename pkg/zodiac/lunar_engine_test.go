package zodiac

import (
	"github.com/kaecer68/lunar-zenith/pkg/celestial"
	"testing"
	"time"
)

func TestLunarEngine_GetLunarDate(t *testing.T) {
	engine := &LunarEngine{}

	// 2024-02-10 應該是農曆正月初一
	t1 := time.Date(2024, 2, 10, 12, 0, 0, 0, time.UTC)
	pt1 := celestial.NewPrecisionTime(t1)
	ld1 := engine.GetLunarDate(pt1.JD)

	if ld1.Month != 1 || ld1.Day != 1 {
		t.Errorf("2024-02-10 應為正月初一, got %d月%d日", ld1.Month, ld1.Day)
	}

	// 2024-05-15 (陽曆) 應該是農曆四月初八 (佛誕日)
	t2 := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	pt2 := celestial.NewPrecisionTime(t2)
	ld2 := engine.GetLunarDate(pt2.JD)

	if ld2.Month != 4 || ld2.Day != 8 {
		t.Errorf("2024-05-15 應為四月初八, got %d月%d日", ld2.Month, ld2.Day)
	}

	// 額外驗證：1972, 1990 及冬至
	bugTests := []struct {
		date       time.Time
		wantYear   int
		wantPillar string
	}{
		{time.Date(1972, 6, 8, 12, 0, 0, 0, time.UTC), 1972, "壬子"},
		{time.Date(1990, 6, 15, 12, 0, 0, 0, time.UTC), 1990, "庚午"},
		{time.Date(2023, 12, 22, 12, 0, 0, 0, time.UTC), 2023, "癸卯"},
	}

	for _, tt := range bugTests {
		pt := celestial.NewPrecisionTime(tt.date)
		ld := engine.GetLunarDate(pt.JD)
		if ld.Year != tt.wantYear {
			t.Errorf("Date %s: got year %d, want %d", tt.date.Format("2006-01-02"), ld.Year, tt.wantYear)
		}
		if ld.YearPillar.String() != tt.wantPillar {
			t.Errorf("Date %s: got pillar %s, want %s", tt.date.Format("2006-01-02"), ld.YearPillar.String(), tt.wantPillar)
		}
	}
}
