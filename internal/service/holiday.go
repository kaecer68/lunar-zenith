package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kaecer68/lunar-zenith/pkg/celestial"
	"github.com/kaecer68/lunar-zenith/pkg/zodiac"
)

// HolidayType 定義假期的類型
type HolidayType string

const (
	TypePublicHoliday HolidayType = "holiday" // 國定假日
	TypeWorkday       HolidayType = "workday" // 補班日
	TypeFestival      HolidayType = "festival"
	TypeCommemoration HolidayType = "commemoration"
	TypeWeekend       HolidayType = "weekend" // 週末
)

const (
	regionTaiwan = "tw"
	regionChina  = "cn"
)

var taiwanFixedObservances = []TaiwanHoliday{
	{Date: "0101", Name: "中華民國開國紀念日", IsHoliday: true, Category: TypePublicHoliday},
	{Date: "0228", Name: "和平紀念日", IsHoliday: true, Category: TypePublicHoliday},
	{Date: "0312", Name: "國父逝世紀念日", IsHoliday: false, Category: TypeCommemoration},
	{Date: "0314", Name: "反侵略日", IsHoliday: false, Category: TypeCommemoration},
	{Date: "0321", Name: "民族平等紀念日", IsHoliday: false, Category: TypeCommemoration},
	{Date: "0329", Name: "革命先烈紀念日", IsHoliday: false, Category: TypeCommemoration},
	{Date: "0407", Name: "言論自由日", IsHoliday: false, Category: TypeCommemoration},
	{Date: "0626", Name: "原住民族抵抗日", IsHoliday: false, Category: TypeCommemoration},
	{Date: "0715", Name: "解嚴紀念日", IsHoliday: false, Category: TypeCommemoration},
	{Date: "0801", Name: "原住民族日", IsHoliday: false, Category: TypeCommemoration},
	{Date: "0815", Name: "終戰紀念日", IsHoliday: false, Category: TypeCommemoration},
	{Date: "0823", Name: "八二三紀念日", IsHoliday: false, Category: TypeCommemoration},
	{Date: "0921", Name: "國家防災日", IsHoliday: false, Category: TypeCommemoration},
	{Date: "0928", Name: "孔子誕辰紀念日", IsHoliday: true, Category: TypePublicHoliday},
	{Date: "1010", Name: "國慶日", IsHoliday: true, Category: TypePublicHoliday},
	{Date: "1024", Name: "臺灣聯合國日", IsHoliday: false, Category: TypeCommemoration},
	{Date: "1025", Name: "臺灣光復暨金門古寧頭大捷紀念日", IsHoliday: true, Category: TypePublicHoliday},
	{Date: "1112", Name: "國父誕辰紀念日", IsHoliday: false, Category: TypeCommemoration},
	{Date: "1225", Name: "行憲紀念日", IsHoliday: true, Category: TypePublicHoliday},
	{Date: "1228", Name: "全國客家日", IsHoliday: false, Category: TypeCommemoration},

	{Date: "0119", Name: "消防節", IsHoliday: false, Category: TypeFestival},
	{Date: "0308", Name: "婦女節", IsHoliday: false, Category: TypeFestival},
	{Date: "0312", Name: "植樹節", IsHoliday: false, Category: TypeFestival},
	{Date: "0329", Name: "青年節", IsHoliday: false, Category: TypeFestival},
	{Date: "0330", Name: "國際醫師節", IsHoliday: false, Category: TypeFestival},
	{Date: "0404", Name: "兒童節", IsHoliday: true, Category: TypePublicHoliday},
	{Date: "0501", Name: "勞動節", IsHoliday: true, Category: TypePublicHoliday},
	{Date: "0512", Name: "國際護理師節", IsHoliday: false, Category: TypeFestival},
	{Date: "0605", Name: "環境日", IsHoliday: false, Category: TypeFestival},
	{Date: "0608", Name: "國家海洋日", IsHoliday: false, Category: TypeFestival},
	{Date: "0615", Name: "警察節", IsHoliday: false, Category: TypeFestival},
	{Date: "0701", Name: "漁民節", IsHoliday: false, Category: TypeFestival},
	{Date: "0808", Name: "父親節", IsHoliday: false, Category: TypeFestival},
	{Date: "0903", Name: "軍人節", IsHoliday: false, Category: TypeFestival},
	{Date: "0909", Name: "國民體育日", IsHoliday: false, Category: TypeFestival},
	{Date: "0928", Name: "教師節", IsHoliday: true, Category: TypePublicHoliday},
	{Date: "1108", Name: "海巡節", IsHoliday: false, Category: TypeFestival},
	{Date: "1112", Name: "中華文化復興節", IsHoliday: false, Category: TypeFestival},
	{Date: "1203", Name: "身心障礙者日", IsHoliday: false, Category: TypeFestival},
	{Date: "1210", Name: "人權日", IsHoliday: false, Category: TypeFestival},
	{Date: "1218", Name: "移民日", IsHoliday: false, Category: TypeFestival},
}

func (s *HolidayService) getChinaObservanceOverrides(t time.Time) ([]TaiwanHoliday, bool) {
	md := t.Format("0102")
	observances := make([]TaiwanHoliday, 0, 2)

	if md == "0601" {
		observances = append(observances, TaiwanHoliday{Name: "兒童節", IsHoliday: false, Category: TypeFestival})
		return observances, true
	}

	if md == "0801" {
		observances = append(observances, TaiwanHoliday{Name: "建軍節", IsHoliday: false, Category: TypeFestival})
		return observances, true
	}

	if md == thirdSundayOfJune(t.Year()).Format("0102") {
		observances = append(observances, TaiwanHoliday{Name: "父親節", IsHoliday: false, Category: TypeFestival})
		return observances, true
	}

	if md == "0910" {
		observances = append(observances, TaiwanHoliday{Name: "教師節", IsHoliday: false, Category: TypeFestival})
		return observances, true
	}

	if md == "0404" {
		if qingmingDay(t.Year()) == 4 {
			observances = append(observances, TaiwanHoliday{Name: "清明節", IsHoliday: true, Category: TypePublicHoliday})
			return observances, true
		}
		return nil, true
	}

	if md == "0808" || md == "0928" {
		return nil, true
	}

	return nil, false
}

// TaiwanHoliday 封裝台灣政府公告的單條假期資訊
type TaiwanHoliday struct {
	Date        string      `json:"date"`        // 格式: YYYYMMDD
	Name        string      `json:"name"`        // 假期名稱 (如: 春節)
	IsHoliday   bool        `json:"is_holiday"`  // 是否放假
	Description string      `json:"description"` // 備註
	Category    HolidayType `json:"category"`    // 類別
}

// HolidayService 處理假期數據的檢索服務
type HolidayService struct {
	holidays map[string]TaiwanHoliday
	fallback *HolidayService
	region   string
}

// NewHolidayService 創建一個新的假期服務實例
func NewHolidayService() *HolidayService {
	return &HolidayService{
		holidays: make(map[string]TaiwanHoliday),
		region:   regionTaiwan,
	}
}

// NewHolidayServiceWithFallback 創建帶 fallback 的假期服務
func NewHolidayServiceWithFallback(fallback *HolidayService) *HolidayService {
	return &HolidayService{
		holidays: make(map[string]TaiwanHoliday),
		fallback: fallback,
		region:   regionChina,
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

func (s *HolidayService) IsHoliday(dateStr string) (bool, string) {
	if h, ok := s.holidays[dateStr]; ok {
		if h.Category == TypeWorkday {
			return false, ""
		}
		if !h.IsHoliday {
			return false, h.Name
		}
		return true, h.Name
	}

	t, err := time.Parse("20060102", dateStr)
	if err != nil {
		return false, ""
	}

	if s.region == regionChina {
		observances, blockFallback := s.getChinaObservanceOverrides(t)
		if len(observances) > 0 {
			return mergeObservances(observances)
		}
		if blockFallback {
			if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
				return true, ""
			}
			return false, ""
		}
	}

	if s.fallback != nil {
		return s.fallback.IsHoliday(dateStr)
	}

	observances := s.getTaiwanObservances(t)
	if len(observances) > 0 {
		return mergeObservances(observances)
	}

	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		return true, ""
	}

	return false, ""
}

func mergeObservances(observances []TaiwanHoliday) (bool, string) {
	nameSet := make(map[string]struct{})
	names := make([]string, 0, len(observances))
	isHoliday := false

	for _, obs := range observances {
		if obs.Name != "" {
			if _, ok := nameSet[obs.Name]; !ok {
				nameSet[obs.Name] = struct{}{}
				names = append(names, obs.Name)
			}
		}
		if obs.Category != TypeWorkday && obs.IsHoliday {
			isHoliday = true
		}
	}

	return isHoliday, strings.Join(names, "、")
}

func (s *HolidayService) getTaiwanObservances(t time.Time) []TaiwanHoliday {
	observances := make([]TaiwanHoliday, 0, 8)
	md := t.Format("0102")

	for _, rule := range taiwanFixedObservances {
		if rule.Date == md {
			observances = append(observances, rule)
		}
	}

	if md == secondSundayOfMay(t.Year()).Format("0102") {
		observances = append(observances, TaiwanHoliday{Name: "母親節", IsHoliday: false, Category: TypeFestival})
	}

	if md == fourthSundayOfAugust(t.Year()).Format("0102") {
		observances = append(observances, TaiwanHoliday{Name: "祖父母節", IsHoliday: false, Category: TypeFestival})
	}

	if md == fmt.Sprintf("04%02d", qingmingDay(t.Year())) {
		observances = append(observances, TaiwanHoliday{Name: "清明節", IsHoliday: true, Category: TypePublicHoliday})
	}

	pt := celestial.NewPrecisionTime(t)
	lunar := (&zodiac.LunarEngine{}).GetLunarDate(pt.JD)

	if lunar.Month == 1 && lunar.Day == 1 {
		observances = append(observances,
			TaiwanHoliday{Name: "春節", IsHoliday: true, Category: TypePublicHoliday},
			TaiwanHoliday{Name: "道教節", IsHoliday: false, Category: TypeFestival},
		)
	}
	if lunar.Month == 1 && lunar.Day == 15 {
		observances = append(observances, TaiwanHoliday{Name: "元宵節", IsHoliday: false, Category: TypeFestival})
	}
	if lunar.Month == 4 && lunar.Day == 8 {
		observances = append(observances, TaiwanHoliday{Name: "佛陀誕辰日", IsHoliday: false, Category: TypeFestival})
	}
	if lunar.Month == 5 && lunar.Day == 5 {
		observances = append(observances, TaiwanHoliday{Name: "端午節", IsHoliday: true, Category: TypePublicHoliday})
	}
	if lunar.Month == 7 && lunar.Day == 15 {
		observances = append(observances, TaiwanHoliday{Name: "中元節", IsHoliday: false, Category: TypeFestival})
	}
	if lunar.Month == 8 && lunar.Day == 15 {
		observances = append(observances, TaiwanHoliday{Name: "中秋節", IsHoliday: true, Category: TypePublicHoliday})
	}
	if lunar.Month == 9 && lunar.Day == 9 {
		observances = append(observances, TaiwanHoliday{Name: "重陽節", IsHoliday: false, Category: TypeFestival})
	}
	if lunar.Month == 12 && isLunarMonthEnd(t) {
		observances = append(observances, TaiwanHoliday{Name: "除夕", IsHoliday: true, Category: TypePublicHoliday})
	}

	return observances
}

func secondSundayOfMay(year int) time.Time {
	t := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	for t.Weekday() != time.Sunday {
		t = t.AddDate(0, 0, 1)
	}
	return t.AddDate(0, 0, 7)
}

func fourthSundayOfAugust(year int) time.Time {
	t := time.Date(year, 8, 1, 0, 0, 0, 0, time.UTC)
	for t.Weekday() != time.Sunday {
		t = t.AddDate(0, 0, 1)
	}
	return t.AddDate(0, 0, 21)
}

func thirdSundayOfJune(year int) time.Time {
	t := time.Date(year, 6, 1, 0, 0, 0, 0, time.UTC)
	for t.Weekday() != time.Sunday {
		t = t.AddDate(0, 0, 1)
	}
	return t.AddDate(0, 0, 14)
}

func qingmingDay(year int) int {
	v := float64(year%100)*0.2422 + 4.81
	day := int(v) - (year%100)/4
	if day < 4 {
		return 4
	}
	if day > 5 {
		return 5
	}
	return day
}

func isLunarMonthEnd(t time.Time) bool {
	ptToday := celestial.NewPrecisionTime(t)
	ptTomorrow := celestial.NewPrecisionTime(t.AddDate(0, 0, 1))
	eng := &zodiac.LunarEngine{}
	today := eng.GetLunarDate(ptToday.JD)
	tomorrow := eng.GetLunarDate(ptTomorrow.JD)

	if today.Month == 12 && tomorrow.Month == 1 {
		return true
	}

	return today.Month == 12 && tomorrow.Month == 12 && tomorrow.Day < today.Day
}
