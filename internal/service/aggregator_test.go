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

	agg := NewAggregator(hSvc)

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
