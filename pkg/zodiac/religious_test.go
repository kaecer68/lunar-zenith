package zodiac

import (
	"testing"
)

func TestReligiousCalendar(t *testing.T) {
	// 以 2024 年為例
	// 佛曆: 2024 + 544 = 2568
	// 道曆: 2024 + 2697 = 4721
	rel := GetReligiousCalendar(2024)
	
	if rel.BuddhistYear != 2568 {
		t.Errorf("2024 佛曆計算錯誤: got %d, want 2568", rel.BuddhistYear)
	}
	
	if rel.TaoistYear != 4721 {
		t.Errorf("2024 道曆計算錯誤: got %d, want 4721", rel.TaoistYear)
	}
	
	if rel.FormatBuddhist() != "佛曆 2568 年" {
		t.Errorf("佛曆格式化錯誤: %s", rel.FormatBuddhist())
	}
	
	if rel.FormatTaoist() != "道曆 4721 年" {
		t.Errorf("道曆格式化錯誤: %s", rel.FormatTaoist())
	}
}
