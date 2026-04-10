package service

import (
	"context"
	"testing"
	"time"

	lunarv1 "github.com/kaecer68/lunar-zenith/gen"
	"github.com/kaecer68/lunar-zenith/pkg/western_astro"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestToCalendarGRPCResponseAlignsSharedFields(t *testing.T) {
	agg := NewAggregator(nil, nil)
	res := agg.GetCalendar(time.Date(2025, 10, 6, 12, 0, 0, 0, time.UTC))

	rest := toCalendarRESTResponse(res)
	grpcRes := toCalendarGRPCResponse(res)

	if grpcRes.LunarDate != rest.LunarDate {
		t.Errorf("LunarDate = %q; want %q", grpcRes.LunarDate, rest.LunarDate)
	}
	if grpcRes.SolarTerm == nil {
		t.Fatalf("SolarTerm should not be nil")
	}
	if grpcRes.SolarTerm.Name != rest.SolarTerm.Name {
		t.Errorf("SolarTerm.Name = %q; want %q", grpcRes.SolarTerm.Name, rest.SolarTerm.Name)
	}
	if grpcRes.MoonLongitude != rest.MoonLongitude {
		t.Errorf("MoonLongitude = %v; want %v", grpcRes.MoonLongitude, rest.MoonLongitude)
	}
	if grpcRes.MoonElongation != rest.MoonElongation {
		t.Errorf("MoonElongation = %v; want %v", grpcRes.MoonElongation, rest.MoonElongation)
	}
	if len(rest.LunarFestivals) == 0 {
		t.Fatalf("expected lunar festivals for parity fixture")
	}
	if len(grpcRes.LunarFestivals) != len(rest.LunarFestivals) {
		t.Fatalf("len(LunarFestivals) = %d; want %d", len(grpcRes.LunarFestivals), len(rest.LunarFestivals))
	}
	for i, want := range rest.LunarFestivals {
		got := grpcRes.LunarFestivals[i]
		if got.Priority != int32(want.Priority) {
			t.Errorf("LunarFestivals[%d].Priority = %d; want %d", i, got.Priority, want.Priority)
		}
	}
}

func TestGrpcServerGetCalendarInvalidDateUsesInvalidArgument(t *testing.T) {
	server := NewGrpcServer(NewAggregator(nil, nil))

	_, err := server.GetCalendar(context.Background(), &lunarv1.GetCalendarRequest{Date: "2026/04/11"})
	if err == nil {
		t.Fatalf("expected error for invalid date")
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("status.Code(err) = %v; want %v", status.Code(err), codes.InvalidArgument)
	}
	if status.Convert(err).Message() != invalidCalendarDateMessage {
		t.Errorf("status message = %q; want %q", status.Convert(err).Message(), invalidCalendarDateMessage)
	}
}

func TestGrpcServerGetCalendarEmptyDateUsesTaipeiToday(t *testing.T) {
	server := NewGrpcServer(NewAggregator(nil, nil))
	want, err := resolveCalendarQueryTime("")
	if err != nil {
		t.Fatalf("resolveCalendarQueryTime(\"\") error = %v", err)
	}

	resp, err := server.GetCalendar(context.Background(), &lunarv1.GetCalendarRequest{})
	if err != nil {
		t.Fatalf("GetCalendar() error = %v", err)
	}
	if resp.GregorianDate != want.Format("2006-01-02") {
		t.Errorf("GregorianDate = %q; want %q", resp.GregorianDate, want.Format("2006-01-02"))
	}
}

func TestToCalendarGRPCResponseFormatsNextStationDateAsRFC3339(t *testing.T) {
	agg := NewAggregator(nil, nil)
	res := agg.GetCalendar(time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC))
	next := time.Date(2026, 4, 15, 8, 0, 0, 0, time.FixedZone("CST", 8*3600))
	res.WesternAstro = []western_astro.RetrogradeInfo{
		{
			Planet:          western_astro.Mercury,
			NameZh:          "水星",
			Symbol:          "☿",
			IsRetrograde:    true,
			Longitude:       262.2119,
			Speed:           -0.0987,
			NextStationDate: &next,
			StationType:     "station_direct",
		},
	}

	grpcRes := toCalendarGRPCResponse(res)
	if len(grpcRes.WesternAstro) != 1 {
		t.Fatalf("len(WesternAstro) = %d; want 1", len(grpcRes.WesternAstro))
	}

	got := grpcRes.WesternAstro[0].GetNextStationDate()
	want := next.Format(time.RFC3339)
	if got != want {
		t.Errorf("NextStationDate = %q; want %q", got, want)
	}
	if _, err := time.Parse(time.RFC3339, got); err != nil {
		t.Errorf("NextStationDate should parse as RFC3339, got %q: %v", got, err)
	}
}
