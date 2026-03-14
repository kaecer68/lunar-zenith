package service

import (
	"testing"
)

func TestHolidayService(t *testing.T) {
	s := NewHolidayService()
	// 注意：路徑是相對於 go test 執行位置或使用絕對路徑
	// 為了測試穩定性，我們手動注入一條數據
	s.holidays["20240101"] = TaiwanHoliday{
		Date: "20240101",
		Name: "元旦",
		IsHoliday: true,
	}
	s.holidays["20240217"] = TaiwanHoliday{
		Date: "20240217",
		Name: "補班",
		IsHoliday: false,
	}

	// 測試放假
	isHol, name := s.IsHoliday("20240101")
	if !isHol || name != "元旦" {
		t.Errorf("20240101 應該是放假的元旦, got %v, %s", isHol, name)
	}

	// 測試補班
	isHol2, name2 := s.IsHoliday("20240217")
	if isHol2 || name2 != "補班" {
		t.Errorf("20240217 應該是補班不放假, got %v, %s", isHol2, name2)
	}

	// 測試非記錄日期
	isHol3, _ := s.IsHoliday("20240520")
	if isHol3 {
		t.Errorf("未記錄日期不應預設為假期")
	}
}
