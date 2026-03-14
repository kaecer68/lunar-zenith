package service

import (
	"encoding/json"
	"fmt"
	"os"
)

// HolidayType 定義假期的類型
type HolidayType string

const (
	TypePublicHoliday HolidayType = "holiday"     // 國定假日
	TypeWorkday       HolidayType = "workday"     // 補班日
	TypeWeekend       HolidayType = "weekend"     // 週末
)

// TaiwanHoliday 封裝台灣政府公告的單條假期資訊
type TaiwanHoliday struct {
	Date        string      `json:"date"`         // 格式: YYYYMMDD
	Name        string      `json:"name"`         // 假期名稱 (如: 春節)
	IsHoliday   bool        `json:"is_holiday"`   // 是否放假
	Description string      `json:"description"`  // 備註
	Category    HolidayType `json:"category"`     // 類別
}

// HolidayService 處理假期數據的檢索服務
type HolidayService struct {
	holidays map[string]TaiwanHoliday
}

// NewHolidayService 創建一個新的假期服務實例
func NewHolidayService() *HolidayService {
	return &HolidayService{
		holidays: make(map[string]TaiwanHoliday),
	}
}

// LoadFromJSON 從 JSON 文件載入假期數據
func (s *HolidayService) LoadFromJSON(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read holiday file: %w", err)
	}

	var list []TaiwanHoliday
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("failed to unmarshal holiday data: %w", err)
	}

	for _, h := range list {
		s.holidays[h.Date] = h
	}
	return nil
}

// IsHoliday 檢查指定日期 (YYYYMMDD) 是否為假期
func (s *HolidayService) IsHoliday(dateStr string) (bool, string) {
	if h, ok := s.holidays[dateStr]; ok {
		return h.IsHoliday, h.Name
	}
	// 這裡未來可以增加自動判定週末的邏輯
	return false, ""
}
