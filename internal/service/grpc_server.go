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

	// 8. 構建響應
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
		// 確保 SolarTerm 正確填充（這是關鍵需求）
		SolarTerm: &lunarv1.SolarTerm{
			Index:     int32(st.Index),
			Name:      st.Name,
			Longitude: st.Longitude,
		},
		TwelveOfficer: officer,
	}

	// 9. 添加神煞
	for _, s := range ss {
		res.ShenSha = append(res.ShenSha, &lunarv1.ShenSha{
			Name:        s.Name,
			Description: s.Description,
		})
	}

	// 10. 假期信息
	if s.Aggregator.HolidaySvc != nil {
		isHol, name := s.Aggregator.HolidaySvc.IsHoliday(t.Format("20060102"))
		res.HolidayInfo = &lunarv1.HolidayInfo{
			IsHoliday: isHol,
			Name:      name,
		}
	}

	return res, nil
}
