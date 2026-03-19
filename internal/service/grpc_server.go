package service

import (
	"context"
	"time"

	lunarv1 "github.com/kaecer68/lunar-zenith/api/v1"
	"github.com/kaecer68/lunar-zenith/pkg/celestial"
	"github.com/kaecer68/lunar-zenith/pkg/zodiac"
)

// GrpcServer 實現 gRPC LunarService 服務器
type GrpcServer struct {
	lunarv1.UnimplementedLunarServiceServer
	Aggregator *Aggregator
}

// NewGrpcServer 創建 gRPC 服務器
func NewGrpcServer(agg *Aggregator) *GrpcServer {
	return &GrpcServer{
		Aggregator: agg,
	}
}

// GetCalendar 獲取完整曆法數據（gRPC 實現）
func (s *GrpcServer) GetCalendar(ctx context.Context, req *lunarv1.GetCalendarRequest) (*lunarv1.GetCalendarResponse, error) {
	// 1. 解析日期
	var t time.Time
	var err error

	if req.Date == "" {
		t = time.Now()
	} else {
		t, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			return nil, err
		}
	}

	// 2. 獲取曆法數據
	pt := celestial.NewPrecisionTime(t)

	// 3. 獲取四柱
	pillars := zodiac.GetAstrologicalPillar(pt)

	// 4. 獲取農曆
	lunar := s.Aggregator.LunarEng.GetLunarDate(pt.JD)

	// 5. 獲取宗教曆法
	rel := zodiac.GetReligiousCalendar(t.Year())

	// 6. 獲取節氣（關鍵數據）
	st := celestial.GetSolarTerm(pt.JDE)

	// 7. 獲取神煞
	officer := zodiac.GetTwelveOfficer(pillars.Month.BranchIndex, pillars.Day.BranchIndex)
	ss := zodiac.GetYearShenSha(pillars.Year.BranchIndex)

	// 8. 黃曆宜忌 (v1.3.0 新增)
	suitable, avoidable, directions := CalculateAlmanac(officer, pillars.Day.StemIndex)

	// 9. 擴充神煞 (v1.4.0 新增)
	solarMonth := zodiac.GetSolarMonth(st.Longitude)
	mansion := zodiac.GetTwentyEightMansion(solarMonth, pillars.Day.StemIndex, pillars.Day.BranchIndex)
	dailyDeity := zodiac.GetDailyDeity(pillars.Day.BranchIndex)
	fetalGod := zodiac.GetFetalGod(pillars.Day.StemIndex)
	clashSha := zodiac.GetClashSha(pillars.Day.BranchIndex)

	// 10. 農曆節日 (v1.4.0 新增)
	var lunarFestivals []*lunarv1.LunarFestival
	lunarFests := zodiac.GetLunarFestival(lunar.Month, lunar.Day)
	for _, f := range lunarFests {
		lunarFestivals = append(lunarFestivals, &lunarv1.LunarFestival{
			Name:        f.Name,
			Type:        f.Type,
			Description: f.Description,
			Priority:    int32(f.Priority),
		})
	}

	// 11. 構建響應
	res := &lunarv1.GetCalendarResponse{
		GregorianDate: t.Format("2006-01-02"),
		JulianDay:     pt.JD,
		DeltaT:        pt.DeltaT,
		Lunar: &lunarv1.LunarInfo{
			Year:        int32(lunar.Year),
			Month:       int32(lunar.Month),
			Day:         int32(lunar.Day),
			IsLeap:      lunar.IsLeap,
			StringValue: lunar.String(),
		},
		Buddhist: rel.FormatBuddhist(),
		Taoist:   rel.FormatTaoist(),
		Pillars: &lunarv1.Pillars{
			Year:  zodiac.GetStemBranchName(pillars.Year.StemIndex, pillars.Year.BranchIndex),
			Month: zodiac.GetStemBranchName(pillars.Month.StemIndex, pillars.Month.BranchIndex),
			Day:   zodiac.GetStemBranchName(pillars.Day.StemIndex, pillars.Day.BranchIndex),
			Hour:  zodiac.GetStemBranchName(pillars.Hour.StemIndex, pillars.Hour.BranchIndex),
		},
		// 確保 SolarTerm 正確填充
		SolarTerm: &lunarv1.SolarTerm{
			Index:     int32(st.Index),
			Name:      st.Name,
			Longitude: st.Longitude,
		},
		TwelveOfficer: officer,
		Suitable:      suitable,
		Avoidable:     avoidable,
		Directions: &lunarv1.Directions{
			Wealth:  directions.Wealth,
			Fortune: directions.Fortune,
			Study:   directions.Study,
			Love:    directions.Love,
		},
		// v1.4.0 擴充神煞
		Mansion: &lunarv1.Mansion{
			Name:     mansion.Name,
			Animal:   mansion.Animal,
			FullName: mansion.FullName,
			Palace:   mansion.Palace,
			Element:  mansion.Element,
			Index:    int32(mansion.Index),
		},
		DailyDeity: &lunarv1.DailyDeity{
			Name: dailyDeity.Name,
			Type: dailyDeity.Type,
			Desc: dailyDeity.Desc,
		},
		FetalGod: &lunarv1.FetalGod{
			Position:    fetalGod.Position,
			Description: fetalGod.Description,
			Taboo:       fetalGod.Taboo,
		},
		ClashSha: &lunarv1.ClashSha{
			ClashZodiac:  clashSha.ClashZodiac,
			ClashBranch:  clashSha.ClashBranch,
			ShaDirection: clashSha.ShaDirection,
			ShaDesc:      clashSha.ShaDesc,
		},
		LunarFestivals: lunarFestivals,
	}

	// 12. 添加神煞
	for _, s := range ss {
		res.ShenSha = append(res.ShenSha, &lunarv1.ShenSha{
			Name:        s.Name,
			Description: s.Description,
		})
	}

	// 13. 假期信息
	if s.Aggregator.HolidaySvc != nil {
		isHol, name := s.Aggregator.HolidaySvc.IsHoliday(t.Format("20060102"))
		res.HolidayInfo = &lunarv1.HolidayInfo{
			IsHoliday: isHol,
			Name:      name,
		}
	}

	// 14. 大陸假期信息 (v1.4.0 新增)
	if s.Aggregator.ChinaHolidaySvc != nil {
		isHol, name := s.Aggregator.ChinaHolidaySvc.IsHoliday(t.Format("20060102"))
		res.ChinaHolidayInfo = &lunarv1.HolidayInfo{
			IsHoliday: isHol,
			Name:      name,
		}
	}

	return res, nil
}
