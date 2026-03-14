package zodiac

import (
	"fmt"
)

// ReligiousCalendar 封裝宗教曆法資訊
type ReligiousCalendar struct {
	BuddhistYear int // 佛曆年份 (Gregorian + 544)
	TaoistYear   int // 道曆年份 (Gregorian + 2697)
}

// GetReligiousCalendar 獲取指定西元年份的宗教曆法年份
func GetReligiousCalendar(gregorianYear int) ReligiousCalendar {
	return ReligiousCalendar{
		BuddhistYear: gregorianYear + 544,
		TaoistYear:   gregorianYear + 2697,
	}
}

// FormatBuddhist 返回佛曆格式化字符串 (如: 佛曆 2567 年)
func (r ReligiousCalendar) FormatBuddhist() string {
	return fmt.Sprintf("佛曆 %d 年", r.BuddhistYear)
}

// FormatTaoist 返回道曆格式化字符串 (如: 道曆 4721 年)
func (r ReligiousCalendar) FormatTaoist() string {
	return fmt.Sprintf("道曆 %d 年", r.TaoistYear)
}
