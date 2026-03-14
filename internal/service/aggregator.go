package service

import (
	"time"
	"github.com/kaecer68/lunar-zenith/pkg/celestial"
	"github.com/kaecer68/lunar-zenith/pkg/zodiac"
)

// CalendarResponse 曆法全家桶：聚合所有維度的數據
type CalendarResponse struct {
	// 1. 基礎數據
	GregorianDate string  `json:"gregorian_date"` // YYYY-MM-DD
	JulianDay     float64 `json:"julian_day"`
	DeltaT        float64 `json:"delta_t"`

	// 2. 農曆與宗教
	Lunar    zodiac.LunarDate      `json:"lunar"`
	Buddhist string                `json:"buddhist"`
	Taoist   string                `json:"taoist"`

	// 3. 四柱八字
	Pillars zodiac.AstrologicalPillar `json:"pillars"`

	// 4. 天文節氣
	SolarTerm celestial.SolarTermInfo `json:"solar_term"`

	// 5. 神煞與建除
	TwelveOfficer string                 `json:"twelve_officer"`
	ShenSha       []zodiac.CommonShenSha `json:"shen_sha"`

	// 6. 行政假期
	HolidayInfo struct {
		IsHoliday bool   `json:"is_holiday"`
		Name      string `json:"name"`
	} `json:"holiday_info"`
}

// Aggregator 聚合服務
type Aggregator struct {
	HolidaySvc *HolidayService
	LunarEng   *zodiac.LunarEngine
}

// NewAggregator 創建聚合器
func NewAggregator(h *HolidayService) *Aggregator {
	return &Aggregator{
		HolidaySvc: h,
		LunarEng:   &zodiac.LunarEngine{},
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
	
	res := CalendarResponse{
		GregorianDate: t.Format("2006-01-02"),
		JulianDay:     pt.JD,
		DeltaT:        pt.DeltaT,
		Lunar:         lunar,
		Buddhist:      rel.FormatBuddhist(),
		Taoist:        rel.FormatTaoist(),
		Pillars:       pillars,
		SolarTerm:     st,
		TwelveOfficer: officer,
		ShenSha:       ss,
	}
	
	// 6. 假期 (若已加載)
	if a.HolidaySvc != nil {
		isHol, name := a.HolidaySvc.IsHoliday(t.Format("20060102"))
		res.HolidayInfo.IsHoliday = isHol
		res.HolidayInfo.Name = name
	}
	
	return res
}
