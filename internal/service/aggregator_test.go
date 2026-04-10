package service

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestAggregator_GetCalendar(t *testing.T) {
	hSvc := NewHolidayService()
	// 模擬元旦假期數據
	hSvc.holidays["20240101"] = TaiwanHoliday{
		Date:      "20240101",
		Name:      "開國紀念日",
		IsHoliday: true,
	}

	chinaSvc := NewHolidayService()

	agg := NewAggregator(hSvc, chinaSvc)

	// 2024-01-01 12:00
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	res := agg.GetCalendar(testTime)

	// 驗證聚合結果
	if res.GregorianDate != "2024-01-01" {
		t.Errorf("GregorianDate error: got %s", res.GregorianDate)
	}

	if res.HolidayInfo.Name != "開國紀念日" {
		t.Errorf("Holiday aggregation failed: got %s", res.HolidayInfo.Name)
	}

	if res.Buddhist == "" {
		t.Error("Buddhist calendar missing")
	}

	// 輸出 JSON 看看外觀 (輔助偵錯)
	data, _ := json.MarshalIndent(res, "", "  ")
	fmt.Printf("AGGREGATED PROOF:\n%s\n", string(data))
}

func TestAggregator_GetCalendar_OfficerAndDailyDeityParitySamples(t *testing.T) {
	agg := NewAggregator(nil, nil)

	tests := []struct {
		date        string
		wantOfficer string
		wantDeity   string
	}{
		{date: "2024-04-04", wantOfficer: "破", wantDeity: "白虎"},
		{date: "2025-05-05", wantOfficer: "執", wantDeity: "金匱"},
		{date: "2026-06-06", wantOfficer: "執", wantDeity: "朱雀"},
	}

	for _, tt := range tests {
		t.Run(tt.date, func(t *testing.T) {
			queryTime, err := resolveCalendarQueryTime(tt.date)
			if err != nil {
				t.Fatalf("resolveCalendarQueryTime(%s) failed: %v", tt.date, err)
			}
			res := agg.GetCalendar(queryTime)

			if res.TwelveOfficer != tt.wantOfficer {
				t.Errorf("%s officer = %s; want %s", tt.date, res.TwelveOfficer, tt.wantOfficer)
			}
			if res.DailyDeity.Name != tt.wantDeity {
				t.Errorf("%s daily deity = %s; want %s", tt.date, res.DailyDeity.Name, tt.wantDeity)
			}
		})
	}
}
