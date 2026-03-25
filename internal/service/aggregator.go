package service

import (
	"time"

	"github.com/kaecer68/lunar-zenith/pkg/celestial"
	"github.com/kaecer68/lunar-zenith/pkg/western_astro"
	"github.com/kaecer68/lunar-zenith/pkg/zodiac"
)

// FestivalInfo 節日資訊
type FestivalInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// CalendarResponse 曆法全家桶：聚合所有維度的數據
type CalendarResponse struct {
	// 1. 基礎數據
	GregorianDate string  `json:"gregorian_date"` // YYYY-MM-DD
	JulianDay     float64 `json:"julian_day"`
	DeltaT        float64 `json:"delta_t"`

	// 2. 農曆與宗教
	Lunar    zodiac.LunarDate `json:"lunar"`
	Buddhist string           `json:"buddhist"`
	Taoist   string           `json:"taoist"`

	// 3. 四柱八字
	Pillars zodiac.AstrologicalPillar `json:"pillars"`

	// 4. 天文節氣
	SolarTerm celestial.SolarTermInfo `json:"solar_term"`

	// 5. 神煞與建除
	TwelveOfficer string                 `json:"twelve_officer"`
	ShenSha       []zodiac.CommonShenSha `json:"shen_sha"`

	// 6. 黃曆宜忌 (v1.3.0 新增)
	Suitable   []string   `json:"suitable"`   // 宜
	Avoidable  []string   `json:"avoidable"`  // 忌
	Directions Directions `json:"directions"` // 吉神方位

	// 7. 行政假期
	HolidayInfo struct {
		IsHoliday bool   `json:"is_holiday"`
		Name      string `json:"name"`
	} `json:"holiday_info"`

	// 7.5 大陸行政假期
	ChinaHolidayInfo struct {
		IsHoliday bool   `json:"is_holiday"`
		Name      string `json:"name"`
	} `json:"china_holiday_info"`

	// 8. 精密天文數值 (UI 擴充)
	MoonLongitude  float64 `json:"moon_longitude"`  // 月球黃經 (度)
	MoonElongation float64 `json:"moon_elongation"` // 日月黃經差 (度, 0=朔 180=望)

	// 9. 擴充神煞 (二十八星宿、值神、胎神、沖煞)
	Mansion    zodiac.MansionInfo    `json:"mansion"`     // 二十八星宿
	DailyDeity zodiac.DailyDeityInfo `json:"daily_deity"` // 值神
	FetalGod   zodiac.FetalGodInfo   `json:"fetal_god"`   // 胎神
	ClashSha   zodiac.ClashShaInfo   `json:"clash_sha"`   // 沖煞

	// 10. 農曆宗教節日
	LunarFestivals []FestivalInfo `json:"lunar_festivals"` // 當日農曆節日列表

	WesternAstro []western_astro.RetrogradeInfo  `json:"western_astro"`
	Aspects      []western_astro.PlanetaryAspect `json:"aspects"`
}

// Aggregator 聚合服務
type Aggregator struct {
	HolidaySvc      *HolidayService
	ChinaHolidaySvc *HolidayService // 大陸假期服務
	LunarEng        *zodiac.LunarEngine
}

// NewAggregator 創建聚合器
func NewAggregator(h *HolidayService, chinaH *HolidayService) *Aggregator {
	return &Aggregator{
		HolidaySvc:      h,
		ChinaHolidaySvc: chinaH,
		LunarEng:        &zodiac.LunarEngine{},
	}
}

// GetCalendar 獲取指定時間的完整曆法數據包
func (a *Aggregator) GetCalendar(t time.Time) CalendarResponse {
	pt := celestial.NewPrecisionTime(t)

	// 1. 獲取四柱
	pillars := zodiac.GetAstrologicalPillar(pt)

	// 2. 獲取農曆
	lunar := a.LunarEng.GetLunarDate(pt.JD)
	lunar.YearPillar = pillars.Year
	lunar.MonthPillar = pillars.Month
	lunar.DayPillar = pillars.Day

	// 3. 獲取宗教
	rel := zodiac.GetReligiousCalendar(t.Year())

	// 4. 獲取節氣
	st := celestial.GetSolarTerm(pt.JDE)

	// 5. 神煞
	officer := zodiac.GetTwelveOfficer(pillars.Month.BranchIndex, pillars.Day.BranchIndex)
	ss := zodiac.GetYearShenSha(pillars.Year.BranchIndex)

	// 6. 黃曆宜忌 (v1.3.0 新增) - 基於建除十二神和日干
	suitable, avoidable, directions := CalculateAlmanac(officer, pillars.Day.StemIndex)

	// 7. 精密月球位置
	moonLon := celestial.MoonLongitude(pt.JDE)
	moonElong := celestial.MoonPhase(pt.JDE)

	// 8. 擴充神煞 (二十八星宿、值神、胎神、沖煞)
	// 計算星宿需要節氣月份 (寅月起算 1-12)
	solarMonth := zodiac.GetSolarMonth(st.Longitude)
	mansion := zodiac.GetTwentyEightMansion(solarMonth, pillars.Day.StemIndex, pillars.Day.BranchIndex)
	dailyDeity := zodiac.GetDailyDeity(pillars.Day.BranchIndex)
	fetalGod := zodiac.GetFetalGod(pillars.Day.StemIndex)
	clashSha := zodiac.GetClashSha(pillars.Day.BranchIndex)

	// 9. 檢查農曆節日
	var festivals []FestivalInfo
	lunarFests := zodiac.GetLunarFestival(lunar.Month, lunar.Day)
	for _, f := range lunarFests {
		festivals = append(festivals, FestivalInfo{
			Name:        f.Name,
			Type:        f.Type,
			Description: f.Description,
		})
	}

	westernAstroInfos, _ := western_astro.GetAllRetrogradeInfo(t)
	aspects, _ := western_astro.CalculateAspects(t)

	res := CalendarResponse{
		GregorianDate:  t.Format("2006-01-02"),
		JulianDay:      pt.JD,
		DeltaT:         pt.DeltaT,
		Lunar:          lunar,
		Buddhist:       rel.FormatBuddhist(),
		Taoist:         rel.FormatTaoist(),
		Pillars:        pillars,
		SolarTerm:      st,
		TwelveOfficer:  officer,
		ShenSha:        ss,
		Suitable:       suitable,
		Avoidable:      avoidable,
		Directions:     directions,
		MoonLongitude:  moonLon,
		MoonElongation: moonElong,
		Mansion:        mansion,
		DailyDeity:     dailyDeity,
		FetalGod:       fetalGod,
		ClashSha:       clashSha,
		LunarFestivals: festivals,
		WesternAstro:   westernAstroInfos,
		Aspects:        aspects,
	}

	// 6. 假期 (若已加載)
	if a.HolidaySvc != nil {
		isHol, name := a.HolidaySvc.IsHoliday(t.Format("20060102"))
		res.HolidayInfo.IsHoliday = isHol
		res.HolidayInfo.Name = name
	}

	// 6.5 大陸假期 (若已加載)
	if a.ChinaHolidaySvc != nil {
		isHol, name := a.ChinaHolidaySvc.IsHoliday(t.Format("20060102"))
		res.ChinaHolidayInfo.IsHoliday = isHol
		res.ChinaHolidayInfo.Name = name
	}

	return res
}
