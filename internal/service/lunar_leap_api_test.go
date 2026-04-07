package service

import (
	"context"
	"testing"
	"time"

	lunarv1 "github.com/kaecer68/lunar-zenith/gen"
)

func TestCalendarAPI_LeapMonthFieldsStayAligned(t *testing.T) {
	agg := NewAggregator(nil, nil)
	grpcServer := NewGrpcServer(agg)

	cases := []struct {
		name string
		date time.Time
	}{
		{
			name: "verified-no-leap-baseline-2024-02-10",
			date: time.Date(2024, 2, 10, 12, 0, 0, 0, time.UTC),
		},
		{
			name: "risk-year-2023-fixture-placeholder",
			date: time.Date(2023, 8, 16, 12, 0, 0, 0, time.UTC),
		},
		{
			name: "risk-year-2025-fixture-placeholder",
			date: time.Date(2025, 8, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name: "complex-year-2033-fixture-placeholder",
			date: time.Date(2033, 12, 22, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			queryDate := time.Date(tt.date.Year(), tt.date.Month(), tt.date.Day(), 0, 0, 0, 0, time.UTC)
			aggRes := agg.GetCalendar(queryDate)
			restRes := toCalendarRESTResponse(aggRes)

			grpcRes, err := grpcServer.GetCalendar(context.Background(), &lunarv1.GetCalendarRequest{Date: queryDate.Format("2006-01-02")})
			if err != nil {
				t.Fatalf("GetCalendar() error = %v", err)
			}

			if restRes.Lunar.Month != aggRes.Lunar.Month || restRes.Lunar.Day != aggRes.Lunar.Day || restRes.Lunar.IsLeap != aggRes.Lunar.IsLeap {
				t.Errorf("REST lunar fields drifted from aggregator for %s: got month=%d day=%d leap=%v; want month=%d day=%d leap=%v",
					tt.date.Format("2006-01-02"),
					restRes.Lunar.Month, restRes.Lunar.Day, restRes.Lunar.IsLeap,
					aggRes.Lunar.Month, aggRes.Lunar.Day, aggRes.Lunar.IsLeap)
			}

			if restRes.Lunar.StringValue != aggRes.Lunar.String() || restRes.LunarDate != aggRes.Lunar.String() {
				t.Errorf("REST lunar string fields drifted from aggregator for %s", tt.date.Format("2006-01-02"))
			}

			if grpcRes.Lunar == nil {
				t.Fatalf("gRPC lunar payload should not be nil for %s", tt.date.Format("2006-01-02"))
			}

			if int(grpcRes.Lunar.Month) != aggRes.Lunar.Month || int(grpcRes.Lunar.Day) != aggRes.Lunar.Day || grpcRes.Lunar.IsLeap != aggRes.Lunar.IsLeap {
				t.Errorf("gRPC lunar fields drifted from aggregator for %s: got month=%d day=%d leap=%v; want month=%d day=%d leap=%v",
					tt.date.Format("2006-01-02"),
					grpcRes.Lunar.Month, grpcRes.Lunar.Day, grpcRes.Lunar.IsLeap,
					aggRes.Lunar.Month, aggRes.Lunar.Day, aggRes.Lunar.IsLeap)
			}

			if grpcRes.Lunar.StringValue != aggRes.Lunar.String() {
				t.Errorf("gRPC lunar string field drifted from aggregator for %s", tt.date.Format("2006-01-02"))
			}
		})
	}
}

func TestCalendarAPI_LeapMonthRiskYearFixtures(t *testing.T) {
	cases := []struct {
		name string
		date time.Time
		note string
	}{
		{
			name: "2023-risk-year-api-fixture",
			date: time.Date(2023, 8, 16, 12, 0, 0, 0, time.UTC),
			note: "assert authoritative month/day/isLeap across aggregator, REST, gRPC, and any affected festival output",
		},
		{
			name: "2025-risk-year-api-fixture",
			date: time.Date(2025, 8, 1, 12, 0, 0, 0, time.UTC),
			note: "fill after confirming a trusted official or astronomical source for leap-month output",
		},
		{
			name: "2033-complex-year-api-fixture",
			date: time.Date(2033, 12, 22, 12, 0, 0, 0, time.UTC),
			note: "this year is known to be tricky; do not replace this skip with guessed values",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Skipf("TODO: add verified service/API expectations for %s; %s", tt.date.Format("2006-01-02"), tt.note)
		})
	}
}

// TestCalendarAPI_LeapMonthFestivalConsistency 驗證節日輸出與農曆月日保持一致。
//
// GetLunarFestival 與 HolidayService 均純用 lunar.Month / lunar.Day 索引，
// 不額外處理 IsLeap。因此若閏月年導致月份計算偏移，節日可能消失或落在
// 錯誤日期，但這個錯誤不會被一般的欄位一致性測試捕捉到。
// 這組測試的目的是：固定「festival 輸出對應的農曆月日」，確保任何
// lunar_engine 的改動都不會讓已知高優先度節日在非閏月年默默消失。
func TestCalendarAPI_LeapMonthFestivalConsistency(t *testing.T) {
	agg := NewAggregator(nil, nil)

	t.Run("verified-festivals-no-leap-2024", func(t *testing.T) {
		cases := []struct {
			name              string
			date              time.Time
			wantFestivalNames []string // 至少有哪些節日必須出現
		}{
			{
				name:              "mid-autumn-2024-09-17",
				date:              time.Date(2024, 9, 17, 12, 0, 0, 0, time.UTC), // 農曆八月十五
				wantFestivalNames: []string{"中秋節"},
			},
			{
				name:              "dragon-boat-2024-06-10",
				date:              time.Date(2024, 6, 10, 12, 0, 0, 0, time.UTC), // 農曆五月初五
				wantFestivalNames: []string{"端午節"},
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				res := agg.GetCalendar(tt.date)
				festNames := make(map[string]bool, len(res.LunarFestivals))
				for _, f := range res.LunarFestivals {
					festNames[f.Name] = true
				}
				for _, want := range tt.wantFestivalNames {
					if !festNames[want] {
						t.Errorf("date %s (lunar %d/%d): expected festival %q in lunar_festivals, got %v",
							tt.date.Format("2006-01-02"), res.Lunar.Month, res.Lunar.Day, want, res.LunarFestivals)
					}
				}
			})
		}
	})

	// festival-leap-month-risk-year-fixtures
	// 使用公認節日日期做雙重驗證：農曆月日正確 + 節日出現在 lunar_festivals。
	// GetLunarFestival 純用 lunar.Month/Day 索引，不處理 IsLeap。
	// 若閏月導致月份偏移，節日會從 lunar_festivals 消失，此處斷言可捕捉該錯誤。
	t.Run("festival-2023-dragon-boat-2023-06-22", func(t *testing.T) {
		date := time.Date(2023, 6, 22, 12, 0, 0, 0, time.UTC)
		res := agg.GetCalendar(date)

		// 農曆月日必須正確（2023 有閏二月，月份若偏移會是 6/5 而非 5/5）
		if res.Lunar.Month != 5 || res.Lunar.Day != 5 {
			t.Errorf("2023-06-22 (端午節): lunar got %d/%d; want 5/5", res.Lunar.Month, res.Lunar.Day)
		}

		// 節日必須出現
		festNames := make(map[string]bool, len(res.LunarFestivals))
		for _, f := range res.LunarFestivals {
			festNames[f.Name] = true
		}
		if !festNames["端午節"] {
			t.Errorf("2023-06-22: expected 端午節 in lunar_festivals (lunar %d/%d), got %v",
				res.Lunar.Month, res.Lunar.Day, res.LunarFestivals)
		}
	})

	t.Run("festival-2025-mid-autumn-2025-10-06", func(t *testing.T) {
		date := time.Date(2025, 10, 6, 12, 0, 0, 0, time.UTC)
		res := agg.GetCalendar(date)

		// 農曆月日必須正確（2025 有閏六月，月份若偏移會是 9/15 而非 8/15）
		if res.Lunar.Month != 8 || res.Lunar.Day != 15 {
			t.Errorf("2025-10-06 (中秋節): lunar got %d/%d; want 8/15", res.Lunar.Month, res.Lunar.Day)
		}

		// 節日必須出現
		festNames := make(map[string]bool, len(res.LunarFestivals))
		for _, f := range res.LunarFestivals {
			festNames[f.Name] = true
		}
		if !festNames["中秋節"] {
			t.Errorf("2025-10-06: expected 中秋節 in lunar_festivals (lunar %d/%d), got %v",
				res.Lunar.Month, res.Lunar.Day, res.LunarFestivals)
		}
	})
}
