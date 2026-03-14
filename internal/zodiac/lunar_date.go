package zodiac

import (
	"fmt"
)

// LunarDate 封裝繁體中文農曆日期資訊
type LunarDate struct {
	Year        int    // 農曆年份
	Month       int    // 農曆月份 (1-12)
	Day         int    // 農曆日期 (1-30)
	IsLeap      bool   // 是否為閏月
	YearPillar  Sexagenary
	MonthPillar Sexagenary
	DayPillar   Sexagenary
}

// String 返回農曆日期的繁體中文描述
func (l LunarDate) String() string {
	leapStr := ""
	if l.IsLeap {
		leapStr = "閏"
	}
	return fmt.Sprintf("農曆 %s 年 %s%s 月 %s", 
		l.YearPillar.String(), 
		leapStr, monthName(l.Month), dayName(l.Day))
}

func monthName(m int) string {
	names := []string{"", "正", "二", "三", "四", "五", "六", "七", "八", "九", "十", "十一", "臘"}
	if m < 1 || m >= len(names) {
		return "未知"
	}
	return names[m]
}

func dayName(d int) string {
	prefix := []string{"初", "十", "廿", "三"}
	units := []string{"十", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
	
	if d <= 0 || d > 30 {
		return "未知"
	}

	if d == 10 { return "初十" }
	if d == 20 { return "二十" }
	if d == 30 { return "三十" }
	
	p := (d - 1) / 10
	u := d % 10
	return prefix[p] + units[u]
}
