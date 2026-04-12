package service

import (
	"time"

	lunarv1 "github.com/kaecer68/lunar-zenith/v4/gen"
	"github.com/kaecer68/lunar-zenith/v4/pkg/western_astro"
)

const invalidCalendarDateMessage = "Invalid date format, use YYYY-MM-DD"

type fourPillarsResponse struct {
	Year  string `json:"year"`
	Month string `json:"month"`
	Day   string `json:"day"`
	Hour  string `json:"hour"`
}

type calendarRESTResponse struct {
	GregorianDate    string                          `json:"gregorian_date"`
	JulianDay        float64                         `json:"julian_day"`
	DeltaT           float64                         `json:"delta_t"`
	LunarDate        string                          `json:"lunar_date"`
	Lunar            zodiacLunarResponse             `json:"lunar"`
	Buddhist         string                          `json:"buddhist"`
	Taoist           string                          `json:"taoist"`
	Pillars          fourPillarsResponse             `json:"pillars"`
	SolarTerm        solarTermResponse               `json:"solar_term"`
	TwelveOfficer    string                          `json:"twelve_officer"`
	ShenSha          []shenShaResponse               `json:"shen_sha"`
	Suitable         []string                        `json:"suitable"`
	Avoidable        []string                        `json:"avoidable"`
	Directions       Directions                      `json:"directions"`
	HolidayInfo      holidayInfoResponse             `json:"holiday_info"`
	ChinaHolidayInfo holidayInfoResponse             `json:"china_holiday_info"`
	MoonLongitude    float64                         `json:"moon_longitude"`
	MoonElongation   float64                         `json:"moon_elongation"`
	Mansion          mansionResponse                 `json:"mansion"`
	DailyDeity       dailyDeityResponse              `json:"daily_deity"`
	FetalGod         fetalGodResponse                `json:"fetal_god"`
	ClashSha         clashShaResponse                `json:"clash_sha"`
	LunarFestivals   []FestivalInfo                  `json:"lunar_festivals"`
	WesternAstro     []western_astro.RetrogradeInfo  `json:"western_astro"`
	Aspects          []western_astro.PlanetaryAspect `json:"aspects"`
}

type zodiacLunarResponse struct {
	Year        int    `json:"year"`
	Month       int    `json:"month"`
	Day         int    `json:"day"`
	IsLeap      bool   `json:"is_leap"`
	StringValue string `json:"string_value"`
}

type solarTermResponse struct {
	Index     int     `json:"index"`
	Name      string  `json:"name"`
	Longitude float64 `json:"longitude"`
}

type holidayInfoResponse struct {
	IsHoliday bool   `json:"is_holiday"`
	Name      string `json:"name"`
}

type shenShaResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type mansionResponse struct {
	Name     string `json:"name"`
	Animal   string `json:"animal"`
	FullName string `json:"full_name"`
	Palace   string `json:"palace"`
	Element  string `json:"element"`
	Index    int    `json:"index"`
}

type dailyDeityResponse struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Desc string `json:"desc"`
}

type fetalGodResponse struct {
	Position    string `json:"position"`
	Description string `json:"description"`
	Taboo       string `json:"taboo"`
}

type clashShaResponse struct {
	ClashZodiac  string `json:"clash_zodiac"`
	ClashBranch  string `json:"clash_branch"`
	ShaDirection string `json:"sha_direction"`
	ShaDesc      string `json:"sha_desc"`
}

func toCalendarRESTResponse(res CalendarResponse) calendarRESTResponse {
	solarTerm := solarTermResponse{
		Index:     res.SolarTerm.Index,
		Name:      res.SolarTerm.Name,
		Longitude: res.SolarTerm.Longitude,
	}
	shenSha := make([]shenShaResponse, 0, len(res.ShenSha))
	for _, item := range res.ShenSha {
		shenSha = append(shenSha, shenShaResponse{
			Name:        item.Name,
			Description: item.Description,
		})
	}

	return calendarRESTResponse{
		GregorianDate: res.GregorianDate,
		JulianDay:     res.JulianDay,
		DeltaT:        res.DeltaT,
		LunarDate:     res.Lunar.String(),
		Lunar: zodiacLunarResponse{
			Year:        res.Lunar.Year,
			Month:       res.Lunar.Month,
			Day:         res.Lunar.Day,
			IsLeap:      res.Lunar.IsLeap,
			StringValue: res.Lunar.String(),
		},
		Buddhist: res.Buddhist,
		Taoist:   res.Taoist,
		Pillars: fourPillarsResponse{
			Year:  res.Pillars.Year.String(),
			Month: res.Pillars.Month.String(),
			Day:   res.Pillars.Day.String(),
			Hour:  res.Pillars.Hour.String(),
		},
		SolarTerm:     solarTerm,
		TwelveOfficer: res.TwelveOfficer,
		ShenSha:       shenSha,
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
		Mansion: mansionResponse{
			Name:     res.Mansion.Name,
			Animal:   res.Mansion.Animal,
			FullName: res.Mansion.FullName,
			Palace:   res.Mansion.Palace,
			Element:  res.Mansion.Element,
			Index:    res.Mansion.Index,
		},
		DailyDeity: dailyDeityResponse{
			Name: res.DailyDeity.Name,
			Type: res.DailyDeity.Type,
			Desc: res.DailyDeity.Desc,
		},
		FetalGod: fetalGodResponse{
			Position:    res.FetalGod.Position,
			Description: res.FetalGod.Description,
			Taboo:       res.FetalGod.Taboo,
		},
		ClashSha: clashShaResponse{
			ClashZodiac:  res.ClashSha.ClashZodiac,
			ClashBranch:  res.ClashSha.ClashBranch,
			ShaDirection: res.ClashSha.ShaDirection,
			ShaDesc:      res.ClashSha.ShaDesc,
		},
		LunarFestivals: res.LunarFestivals,
		WesternAstro:   res.WesternAstro,
		Aspects:        res.Aspects,
	}
}

func toCalendarGRPCResponse(res CalendarResponse) *lunarv1.GetCalendarResponse {
	rest := toCalendarRESTResponse(res)
	resp := &lunarv1.GetCalendarResponse{
		GregorianDate:    rest.GregorianDate,
		JulianDay:        rest.JulianDay,
		DeltaT:           rest.DeltaT,
		LunarDate:        rest.LunarDate,
		Lunar:            toProtoLunar(rest.Lunar),
		Buddhist:         rest.Buddhist,
		Taoist:           rest.Taoist,
		Pillars:          toProtoPillars(rest.Pillars),
		SolarTerm:        toProtoSolarTerm(rest.SolarTerm),
		TwelveOfficer:    rest.TwelveOfficer,
		Suitable:         rest.Suitable,
		Avoidable:        rest.Avoidable,
		Directions:       toProtoDirections(rest.Directions),
		HolidayInfo:      toProtoHolidayInfo(rest.HolidayInfo),
		ChinaHolidayInfo: toProtoHolidayInfo(rest.ChinaHolidayInfo),
		MoonLongitude:    rest.MoonLongitude,
		MoonElongation:   rest.MoonElongation,
		Mansion: &lunarv1.Mansion{
			Name:     rest.Mansion.Name,
			Animal:   rest.Mansion.Animal,
			FullName: rest.Mansion.FullName,
			Palace:   rest.Mansion.Palace,
			Element:  rest.Mansion.Element,
			Index:    int32(rest.Mansion.Index),
		},
		DailyDeity: &lunarv1.DailyDeity{
			Name: rest.DailyDeity.Name,
			Type: rest.DailyDeity.Type,
			Desc: rest.DailyDeity.Desc,
		},
		FetalGod: &lunarv1.FetalGod{
			Position:    rest.FetalGod.Position,
			Description: rest.FetalGod.Description,
			Taboo:       rest.FetalGod.Taboo,
		},
		ClashSha: &lunarv1.ClashSha{
			ClashZodiac:  rest.ClashSha.ClashZodiac,
			ClashBranch:  rest.ClashSha.ClashBranch,
			ShaDirection: rest.ClashSha.ShaDirection,
			ShaDesc:      rest.ClashSha.ShaDesc,
		},
	}

	for _, s := range rest.ShenSha {
		resp.ShenSha = append(resp.ShenSha, &lunarv1.ShenSha{
			Name:        s.Name,
			Description: s.Description,
		})
	}

	for _, f := range rest.LunarFestivals {
		resp.LunarFestivals = append(resp.LunarFestivals, &lunarv1.LunarFestival{
			Name:        f.Name,
			Type:        f.Type,
			Description: f.Description,
			Priority:    int32(f.Priority),
		})
	}

	for _, info := range rest.WesternAstro {
		item := &lunarv1.WesternAstroInfo{
			Planet:       int32(info.Planet),
			NameZh:       info.NameZh,
			Symbol:       info.Symbol,
			IsRetrograde: info.IsRetrograde,
			Longitude:    info.Longitude,
			Speed:        info.Speed,
		}
		if info.NextStationDate != nil {
			next := info.NextStationDate.Format(time.RFC3339)
			item.NextStationDate = &next
		}
		if info.StationType != "" {
			typ := info.StationType
			item.StationType = &typ
		}
		resp.WesternAstro = append(resp.WesternAstro, item)
	}

	for _, aspect := range rest.Aspects {
		item := &lunarv1.PlanetaryAspect{
			Planet1:       int32(aspect.Planet1),
			Planet1Name:   aspect.Planet1Name,
			Planet1Symbol: aspect.Planet1Symbol,
			Planet2:       int32(aspect.Planet2),
			Planet2Name:   aspect.Planet2Name,
			Planet2Symbol: aspect.Planet2Symbol,
			Aspect:        string(aspect.Aspect),
			Angle:         aspect.Angle,
			Orb:           aspect.Orb,
		}
		if aspect.ExactDate != nil {
			exact := aspect.ExactDate.Format(time.RFC3339)
			item.ExactDate = &exact
		}
		resp.Aspects = append(resp.Aspects, item)
	}

	return resp
}

func toProtoLunar(l zodiacLunarResponse) *lunarv1.LunarInfo {
	return &lunarv1.LunarInfo{
		Year:        int32(l.Year),
		Month:       int32(l.Month),
		Day:         int32(l.Day),
		IsLeap:      l.IsLeap,
		StringValue: l.StringValue,
	}
}

func toProtoPillars(p fourPillarsResponse) *lunarv1.Pillars {
	return &lunarv1.Pillars{
		Year:  p.Year,
		Month: p.Month,
		Day:   p.Day,
		Hour:  p.Hour,
	}
}

func toProtoSolarTerm(s solarTermResponse) *lunarv1.SolarTerm {
	return &lunarv1.SolarTerm{
		Index:     int32(s.Index),
		Name:      s.Name,
		Longitude: s.Longitude,
	}
}

func toProtoDirections(d Directions) *lunarv1.Directions {
	return &lunarv1.Directions{
		Wealth:  d.Wealth,
		Fortune: d.Fortune,
		Study:   d.Study,
		Love:    d.Love,
	}
}

func toProtoHolidayInfo(info holidayInfoResponse) *lunarv1.HolidayInfo {
	return &lunarv1.HolidayInfo{
		IsHoliday: info.IsHoliday,
		Name:      info.Name,
	}
}
