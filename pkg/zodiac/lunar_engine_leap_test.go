package zodiac

import (
	"testing"
	"time"

	"github.com/kaecer68/lunar-zenith/pkg/celestial"
)

func TestLunarEngine_LeapMonthCoverage(t *testing.T) {
	engine := &LunarEngine{}

	t.Run("verified-no-leap-year-2024-baseline", func(t *testing.T) {
		cases := []struct {
			name      string
			date      time.Time
			wantMonth int
			wantDay   int
			wantLeap  bool
		}{
			{
				name:      "lunar-new-year",
				date:      time.Date(2024, 2, 10, 12, 0, 0, 0, time.UTC),
				wantMonth: 1,
				wantDay:   1,
				wantLeap:  false,
			},
			{
				name:      "buddha-birthday",
				date:      time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC),
				wantMonth: 4,
				wantDay:   8,
				wantLeap:  false,
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				pt := celestial.NewPrecisionTime(tt.date)
				got := engine.GetLunarDate(pt.JD)

				if got.Month != tt.wantMonth || got.Day != tt.wantDay || got.IsLeap != tt.wantLeap {
					t.Errorf("date %s: got month=%d day=%d leap=%v; want month=%d day=%d leap=%v",
						tt.date.Format("2006-01-02"), got.Month, got.Day, got.IsLeap,
						tt.wantMonth, tt.wantDay, tt.wantLeap)
				}
			})
		}
	})

	// leap-month-risk-year-fixtures
	// 使用公認節日日期做驗證基準：
	//   2023-06-22 = 端午節 (農曆五月初五)，2023 有閏二月
	//   2025-10-06 = 中秋節 (農曆八月十五)，2025 有閏六月
	// 引擎若未正確計算閏月，後續月份會偏移 1，導致下方斷言失敗。
	// 這些失敗是預期行為，明確記錄 lunar_engine.go:59 TODO 的影響範圍。
	t.Run("2023-dragon-boat-lunar-5-5", func(t *testing.T) {
		pt := celestial.NewPrecisionTime(time.Date(2023, 6, 22, 12, 0, 0, 0, time.UTC))
		got := engine.GetLunarDate(pt.JD)
		if got.Month != 5 || got.Day != 5 || got.IsLeap {
			t.Errorf("2023-06-22 (端午節，2023 有閏二月): got month=%d day=%d leap=%v; want month=5 day=5 leap=false",
				got.Month, got.Day, got.IsLeap)
		}
	})

	t.Run("2025-mid-autumn-lunar-8-15", func(t *testing.T) {
		pt := celestial.NewPrecisionTime(time.Date(2025, 10, 6, 12, 0, 0, 0, time.UTC))
		got := engine.GetLunarDate(pt.JD)
		if got.Month != 8 || got.Day != 15 || got.IsLeap {
			t.Errorf("2025-10-06 (中秋節，2025 有閏六月): got month=%d day=%d leap=%v; want month=8 day=15 leap=false",
				got.Month, got.Day, got.IsLeap)
		}
	})

	t.Run("2033-complex-year-fixture", func(t *testing.T) {
		t.Skipf("TODO: 2033 是已知複雜年，需來自可信天文或官方曆法來源的授權真值才可啟用")
	})
}

// TestLunarEngine_HistoricalDecadeCoverage 驗證演算法在 1900+ 各十年的精確性。
// 測試資料來源：
//   - 春節（正月初一）：維基百科中國農曆年份列表 + 香港天文台曆書
//   - 中秋節（八月十五）：中央氣象局、台灣農民曆、中國科學院紫金山天文台
//   - 閏月 IsLeap：直接對照當年農曆曆書確認
func TestLunarEngine_HistoricalDecadeCoverage(t *testing.T) {
	engine := &LunarEngine{}

	cases := []struct {
		name      string
		date      time.Time
		wantMonth int
		wantDay   int
		wantLeap  bool
	}{
		// ── 春節錨點（正月初一必在所有閏月之前，不受閏月偏移影響）──
		{
			// 己丑年，1949 有閏七月
			name:      "1949-cny",
			date:      time.Date(1949, 1, 29, 12, 0, 0, 0, time.UTC),
			wantMonth: 1, wantDay: 1, wantLeap: false,
		},
		{
			// 甲子年，1984 有閏十月
			name:      "1984-cny",
			date:      time.Date(1984, 2, 2, 12, 0, 0, 0, time.UTC),
			wantMonth: 1, wantDay: 1, wantLeap: false,
		},
		{
			// 庚子年，2020 有閏四月
			name:      "2020-cny",
			date:      time.Date(2020, 1, 25, 12, 0, 0, 0, time.UTC),
			wantMonth: 1, wantDay: 1, wantLeap: false,
		},
		// ── 閏四月年中秋節（閏月之後第四個月，強力驗證後置月份偏移）──
		{
			// 2001 有閏四月；八月初一 = 9/17，加 14 天 = 10/01
			name:      "2001-mid-autumn-after-leap-4th",
			date:      time.Date(2001, 10, 1, 12, 0, 0, 0, time.UTC),
			wantMonth: 8, wantDay: 15, wantLeap: false,
		},
		{
			// 2012 有閏四月；八月初一 = 9/16，加 14 天 = 9/30
			name:      "2012-mid-autumn-after-leap-4th",
			date:      time.Date(2012, 9, 30, 12, 0, 0, 0, time.UTC),
			wantMonth: 8, wantDay: 15, wantLeap: false,
		},
		{
			// 2020 有閏四月；八月初一 = 9/17，加 14 天 = 10/01
			name:      "2020-mid-autumn-after-leap-4th",
			date:      time.Date(2020, 10, 1, 12, 0, 0, 0, time.UTC),
			wantMonth: 8, wantDay: 15, wantLeap: false,
		},
		// ── IsLeap=true 識別（直接驗證閏月月份標記）──
		{
			// 2020 正常四月初一（閏月前，IsLeap=false）
			name:      "2020-regular-4th-day1",
			date:      time.Date(2020, 4, 23, 12, 0, 0, 0, time.UTC),
			wantMonth: 4, wantDay: 1, wantLeap: false,
		},
		{
			// 2023 閏二月初一；新月 3/21 17:23 UTC → CST 3/22 01:23
			name:      "2023-leap-2nd-day1",
			date:      time.Date(2023, 3, 22, 12, 0, 0, 0, time.UTC),
			wantMonth: 2, wantDay: 1, wantLeap: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			pt := celestial.NewPrecisionTime(tt.date)
			got := engine.GetLunarDate(pt.JD)

			if got.Month != tt.wantMonth || got.Day != tt.wantDay || got.IsLeap != tt.wantLeap {
				t.Errorf("date %s: got month=%d day=%d leap=%v; want month=%d day=%d leap=%v",
					tt.date.Format("2006-01-02"), got.Month, got.Day, got.IsLeap,
					tt.wantMonth, tt.wantDay, tt.wantLeap)
			}
		})
	}
}
