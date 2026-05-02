package service

import (
	"encoding/json"
	"testing"

	"github.com/kaecer68/lunar-zenith/v4/pkg/celestial"
	"github.com/kaecer68/lunar-zenith/v4/pkg/zodiac"
	"github.com/stretchr/testify/assert"
)

// TestRESTGRPCResponseParity 驗證 REST 與 gRPC 回應的欄位對齊
func TestRESTGRPCResponseParity(t *testing.T) {
	fixture := CalendarResponse{
		GregorianDate: "2026-03-18",
		JulianDay:     2461117.5,
		DeltaT:        69.184,
		Lunar: zodiac.LunarDate{
			Year:        2026,
			Month:       2,
			Day:         1,
			IsLeap:      false,
			YearPillar:  zodiac.Sexagenary{StemIndex: 2, BranchIndex: 6},
			MonthPillar: zodiac.Sexagenary{StemIndex: 7, BranchIndex: 3},
			DayPillar:   zodiac.Sexagenary{StemIndex: 0, BranchIndex: 0},
		},
		Buddhist: "佛曆 2570 年",
		Taoist:   "道曆 4723 年",
		Pillars: zodiac.AstrologicalPillar{
			Year:  zodiac.Sexagenary{StemIndex: 2, BranchIndex: 6},
			Month: zodiac.Sexagenary{StemIndex: 7, BranchIndex: 3},
			Day:   zodiac.Sexagenary{StemIndex: 0, BranchIndex: 0},
			Hour:  zodiac.Sexagenary{StemIndex: 2, BranchIndex: 2},
		},
		SolarTerm: celestial.SolarTermInfo{
			Index:           5,
			Name:            "春分",
			Longitude:       0,
			StartTime:       "2026-03-20T17:45:00+08:00",
			NextTermName:    "清明",
			NextTermTime:    "2026-04-04T22:17:00+08:00",
			IsTransitionDay: false,
		},
		TwelveOfficer: "建",
		ShenSha: []zodiac.CommonShenSha{
			{Name: "年桃花", Description: "桃花位在 卯"},
		},
		Suitable:  []string{"祭祀", "祈福"},
		Avoidable: []string{"動土"},
		Directions: Directions{
			Wealth:  "東南",
			Fortune: "正北",
			Study:   "西南",
			Love:    "正東",
		},
		HolidayInfo: struct {
			IsHoliday bool   `json:"is_holiday"`
			Name      string `json:"name"`
		}{IsHoliday: true, Name: "春節"},
		ChinaHolidayInfo: struct {
			IsHoliday bool   `json:"is_holiday"`
			Name      string `json:"name"`
		}{IsHoliday: false, Name: ""},
		MoonLongitude:  177.2451,
		MoonElongation: 12.4819,
		Mansion: zodiac.MansionInfo{
			Name: "角", Animal: "蛟", FullName: "角木蛟",
			Palace: "東方青龍", Element: "木", Index: 0,
		},
		DailyDeity: zodiac.DailyDeityInfo{
			Name: "青龍", Type: "吉", Desc: "天乙星，天貴星",
		},
		FetalGod: zodiac.FetalGodInfo{
			Position: "門外東南", Description: "甲日胎神在門外東南", Taboo: "忌修門",
		},
		ClashSha: zodiac.ClashShaInfo{
			ClashZodiac: "沖猴", ClashBranch: "申", ShaDirection: "煞北", ShaDesc: "沖猴煞北",
		},
		LunarFestivals: []FestivalInfo{
			{Name: "天公生", Type: "道教", Description: "玉皇大帝聖誕", Priority: 100},
		},
	}

	restResp := toCalendarRESTResponse(fixture)
	grpcResp := toCalendarGRPCResponse(fixture)

	t.Run("basic_fields", func(t *testing.T) {
		assert.Equal(t, restResp.GregorianDate, grpcResp.GregorianDate)
		assert.Equal(t, restResp.JulianDay, grpcResp.JulianDay)
		assert.Equal(t, restResp.DeltaT, grpcResp.DeltaT)
		assert.Equal(t, restResp.LunarDate, grpcResp.LunarDate)
		assert.Equal(t, restResp.Buddhist, grpcResp.Buddhist)
		assert.Equal(t, restResp.Taoist, grpcResp.Taoist)
		assert.Equal(t, restResp.TwelveOfficer, grpcResp.TwelveOfficer)
		assert.Equal(t, restResp.Suitable, grpcResp.Suitable)
		assert.Equal(t, restResp.Avoidable, grpcResp.Avoidable)
		assert.Equal(t, restResp.MoonLongitude, grpcResp.MoonLongitude)
		assert.Equal(t, restResp.MoonElongation, grpcResp.MoonElongation)
	})

	t.Run("lunar_info", func(t *testing.T) {
		assert.Equal(t, restResp.Lunar.Year, int(grpcResp.Lunar.Year))
		assert.Equal(t, restResp.Lunar.Month, int(grpcResp.Lunar.Month))
		assert.Equal(t, restResp.Lunar.Day, int(grpcResp.Lunar.Day))
		assert.Equal(t, restResp.Lunar.IsLeap, grpcResp.Lunar.IsLeap)
		assert.Equal(t, restResp.Lunar.StringValue, grpcResp.Lunar.StringValue)
	})

	t.Run("solar_term", func(t *testing.T) {
		assert.Equal(t, restResp.SolarTerm.Index, int(grpcResp.SolarTerm.Index))
		assert.Equal(t, restResp.SolarTerm.Name, grpcResp.SolarTerm.Name)
		assert.InDelta(t, restResp.SolarTerm.Longitude, grpcResp.SolarTerm.Longitude, 0.0001)
		assert.Equal(t, restResp.SolarTerm.StartTime, grpcResp.SolarTerm.StartTime)
		assert.Equal(t, restResp.SolarTerm.NextTermName, grpcResp.SolarTerm.NextTermName)
		assert.Equal(t, restResp.SolarTerm.NextTermTime, grpcResp.SolarTerm.NextTermTime)
		assert.Equal(t, restResp.SolarTerm.IsTransitionDay, grpcResp.SolarTerm.IsTransitionDay)
	})

	t.Run("lunar_festivals", func(t *testing.T) {
		assert.Equal(t, len(restResp.LunarFestivals), len(grpcResp.LunarFestivals))
		if len(restResp.LunarFestivals) > 0 {
			assert.Equal(t, restResp.LunarFestivals[0].Name, grpcResp.LunarFestivals[0].Name)
			assert.Equal(t, restResp.LunarFestivals[0].Type, grpcResp.LunarFestivals[0].Type)
			assert.Equal(t, restResp.LunarFestivals[0].Description, grpcResp.LunarFestivals[0].Description)
			assert.Equal(t, restResp.LunarFestivals[0].Priority, int(grpcResp.LunarFestivals[0].Priority))
		}
	})

	t.Run("holiday_info", func(t *testing.T) {
		assert.Equal(t, restResp.HolidayInfo.IsHoliday, grpcResp.HolidayInfo.IsHoliday)
		assert.Equal(t, restResp.HolidayInfo.Name, grpcResp.HolidayInfo.Name)
		assert.Equal(t, restResp.ChinaHolidayInfo.IsHoliday, grpcResp.ChinaHolidayInfo.IsHoliday)
		assert.Equal(t, restResp.ChinaHolidayInfo.Name, grpcResp.ChinaHolidayInfo.Name)
	})
}

// TestRESTResponseSchemaAlignment 驗證 REST response 包含所有契約定義的欄位
func TestRESTResponseSchemaAlignment(t *testing.T) {
	fixture := CalendarResponse{
		GregorianDate: "2026-03-18",
		JulianDay:     2461117.5,
		DeltaT:        69.184,
		Lunar: zodiac.LunarDate{
			Year:   2026,
			Month:  2,
			Day:    1,
			IsLeap: false,
		},
		Buddhist: "佛曆 2570 年",
		Taoist:   "道曆 4723 年",
		Pillars:  zodiac.AstrologicalPillar{},
		SolarTerm: celestial.SolarTermInfo{
			Index:     5,
			Name:      "春分",
			Longitude: 0,
		},
		TwelveOfficer: "建",
		ShenSha:       []zodiac.CommonShenSha{},
		Suitable:      []string{},
		Avoidable:     []string{},
		Directions:    Directions{},
		HolidayInfo: struct {
			IsHoliday bool   `json:"is_holiday"`
			Name      string `json:"name"`
		}{IsHoliday: false, Name: ""},
		ChinaHolidayInfo: struct {
			IsHoliday bool   `json:"is_holiday"`
			Name      string `json:"name"`
		}{IsHoliday: false, Name: ""},
		MoonLongitude:  0,
		MoonElongation: 0,
		Mansion:        zodiac.MansionInfo{},
		DailyDeity:     zodiac.DailyDeityInfo{},
		FetalGod:       zodiac.FetalGodInfo{},
		ClashSha:       zodiac.ClashShaInfo{},
		LunarFestivals: nil,
	}

	restResp := toCalendarRESTResponse(fixture)

	jsonBytes, err := json.Marshal(restResp)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	assert.NoError(t, err)

	requiredFields := []string{
		"gregorian_date", "julian_day", "delta_t", "lunar_date", "lunar",
		"buddhist", "taoist", "pillars", "solar_term", "twelve_officer",
		"shen_sha", "suitable", "avoidable", "directions",
		"holiday_info", "china_holiday_info",
		"moon_longitude", "moon_elongation",
		"mansion", "daily_deity", "fetal_god", "clash_sha",
		"lunar_festivals",
	}

	for _, field := range requiredFields {
		assert.Contains(t, result, field, "REST response 缺少契約必填欄位: %s", field)
	}
}

// TestSolarTermObjectNotString 驗證 solar_term 始終為物件
func TestSolarTermObjectNotString(t *testing.T) {
	fixture := CalendarResponse{
		GregorianDate: "2026-03-18",
		SolarTerm: celestial.SolarTermInfo{
			Index:     5,
			Name:      "春分",
			Longitude: 0,
		},
	}

	restResp := toCalendarRESTResponse(fixture)

	jsonBytes, _ := json.Marshal(restResp)
	var result map[string]interface{}
	json.Unmarshal(jsonBytes, &result)

	solarTerm, ok := result["solar_term"].(map[string]interface{})
	assert.True(t, ok, "solar_term 應該是 object，不是 string")
	assert.Equal(t, "春分", solarTerm["name"])
	assert.Equal(t, float64(5), solarTerm["index"])
}

// TestLunarInfoContainsIsLeap 驗證 lunar 物件包含 is_leap
func TestLunarInfoContainsIsLeap(t *testing.T) {
	fixture := CalendarResponse{
		Lunar: zodiac.LunarDate{
			Year:   2025,
			Month:  6,
			Day:    15,
			IsLeap: true,
		},
	}

	restResp := toCalendarRESTResponse(fixture)
	assert.True(t, restResp.Lunar.IsLeap, "lunar.is_leap 應該正確反映閏月狀態")

	grpcResp := toCalendarGRPCResponse(fixture)
	assert.True(t, grpcResp.Lunar.IsLeap, "gRPC lunar.is_leap 應該正確反映閏月狀態")
}

// TestLunarFestivalPriorityExists 驗證 lunar_festivals 包含 priority
func TestLunarFestivalPriorityExists(t *testing.T) {
	fixture := CalendarResponse{
		LunarFestivals: []FestivalInfo{
			{Name: "天公生", Type: "道教", Description: "玉皇大帝聖誕", Priority: 100},
			{Name: "普通節", Type: "民俗", Description: "次要", Priority: 50},
		},
	}

	restResp := toCalendarRESTResponse(fixture)
	assert.Equal(t, 2, len(restResp.LunarFestivals))
	assert.Equal(t, 100, restResp.LunarFestivals[0].Priority)

	grpcResp := toCalendarGRPCResponse(fixture)
	assert.Equal(t, 2, len(grpcResp.LunarFestivals))
	assert.Equal(t, int32(100), grpcResp.LunarFestivals[0].Priority)
}
