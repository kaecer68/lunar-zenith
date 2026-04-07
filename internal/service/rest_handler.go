package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RestHandler 處理 HTTP 請求
type RestHandler struct {
	Aggregator *Aggregator
}

type fourPillarsResponse struct {
	Year  string `json:"year"`
	Month string `json:"month"`
	Day   string `json:"day"`
	Hour  string `json:"hour"`
}

type calendarRESTResponse struct {
	GregorianDate    string              `json:"gregorian_date"`
	JulianDay        float64             `json:"julian_day"`
	DeltaT           float64             `json:"delta_t"`
	LunarDate        string              `json:"lunar_date"`
	Lunar            zodiacLunarResponse `json:"lunar"`
	Buddhist         string              `json:"buddhist"`
	Taoist           string              `json:"taoist"`
	Pillars          fourPillarsResponse `json:"pillars"`
	SolarTerm        string              `json:"solar_term"`
	SolarTermDetail  solarTermDetail     `json:"solar_term_detail"`
	TwelveOfficer    string              `json:"twelve_officer"`
	ShenSha          interface{}         `json:"shen_sha"`
	Suitable         []string            `json:"suitable"`
	Avoidable        []string            `json:"avoidable"`
	Directions       Directions          `json:"directions"`
	HolidayInfo      holidayInfoResponse `json:"holiday_info"`
	ChinaHolidayInfo holidayInfoResponse `json:"china_holiday_info"`
	MoonLongitude    float64             `json:"moon_longitude"`
	MoonElongation   float64             `json:"moon_elongation"`
	Mansion          interface{}         `json:"mansion"`
	DailyDeity       interface{}         `json:"daily_deity"`
	FetalGod         interface{}         `json:"fetal_god"`
	ClashSha         interface{}         `json:"clash_sha"`
	LunarFestivals   []FestivalInfo      `json:"lunar_festivals"`
	WesternAstro     interface{}         `json:"western_astro"`
	Aspects          interface{}         `json:"aspects"`
}

type zodiacLunarResponse struct {
	Year        int    `json:"year"`
	Month       int    `json:"month"`
	Day         int    `json:"day"`
	IsLeap      bool   `json:"is_leap"`
	StringValue string `json:"string_value"`
}

type solarTermDetail struct {
	Index     int     `json:"index"`
	Name      string  `json:"name"`
	Longitude float64 `json:"longitude"`
}

type holidayInfoResponse struct {
	IsHoliday bool   `json:"is_holiday"`
	Name      string `json:"name"`
}

// NewRestHandler 創建 REST 處理器
func NewRestHandler(agg *Aggregator) *RestHandler {
	return &RestHandler{Aggregator: agg}
}

// RegisterRoutes 註冊路由
func (h *RestHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/v1/calendar", h.GetCalendar)
}

// GetCalendar 獲取曆法數據
// Query: ?date=2024-03-14
func (h *RestHandler) GetCalendar(c *gin.Context) {
	dateStr := c.Query("date")
	var t time.Time
	var err error

	if dateStr == "" {
		t = time.Now()
	} else {
		t, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
			return
		}
	}

	res := h.Aggregator.GetCalendar(t)
	c.JSON(http.StatusOK, toCalendarRESTResponse(res))
}

func toCalendarRESTResponse(res CalendarResponse) calendarRESTResponse {
	return calendarRESTResponse{
		GregorianDate: res.GregorianDate,
		JulianDay:     res.JulianDay,
		DeltaT:        res.DeltaT,
		LunarDate:     res.Lunar.String(),
		Pillars: fourPillarsResponse{
			Year:  res.Pillars.Year.String(),
			Month: res.Pillars.Month.String(),
			Day:   res.Pillars.Day.String(),
			Hour:  res.Pillars.Hour.String(),
		},
		SolarTerm: res.SolarTerm.Name,
		Lunar: zodiacLunarResponse{
			Year:        res.Lunar.Year,
			Month:       res.Lunar.Month,
			Day:         res.Lunar.Day,
			IsLeap:      res.Lunar.IsLeap,
			StringValue: res.Lunar.String(),
		},
		Buddhist: res.Buddhist,
		Taoist:   res.Taoist,
		SolarTermDetail: solarTermDetail{
			Index:     res.SolarTerm.Index,
			Name:      res.SolarTerm.Name,
			Longitude: res.SolarTerm.Longitude,
		},
		TwelveOfficer: res.TwelveOfficer,
		ShenSha:       res.ShenSha,
		Suitable:      res.Suitable,
		Avoidable:     res.Avoidable,
		Directions:    res.Directions,
		HolidayInfo: holidayInfoResponse{
			IsHoliday: res.HolidayInfo.IsHoliday,
			Name:      res.HolidayInfo.Name,
		},
		ChinaHolidayInfo: holidayInfoResponse{
			IsHoliday: res.ChinaHolidayInfo.IsHoliday,
			Name:      res.ChinaHolidayInfo.Name,
		},
		MoonLongitude:  res.MoonLongitude,
		MoonElongation: res.MoonElongation,
		Mansion:        res.Mansion,
		DailyDeity:     res.DailyDeity,
		FetalGod:       res.FetalGod,
		ClashSha:       res.ClashSha,
		LunarFestivals: res.LunarFestivals,
		WesternAstro:   res.WesternAstro,
		Aspects:        res.Aspects,
	}
}
