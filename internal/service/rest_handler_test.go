package service

import (
	"testing"
	"time"

	"github.com/kaecer68/lunar-zenith/pkg/zodiac"
)

func TestToCalendarRESTResponseAlignsContractFields(t *testing.T) {
	agg := NewAggregator(nil, nil)
	res := agg.GetCalendar(time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC))

	got := toCalendarRESTResponse(res)

	if got.GregorianDate != "2026-03-18" {
		t.Errorf("GregorianDate = %q; want 2026-03-18", got.GregorianDate)
	}
	if got.LunarDate == "" {
		t.Errorf("LunarDate should not be empty")
	}
	if got.Lunar.StringValue != got.LunarDate {
		t.Errorf("Lunar.StringValue = %q; want %q", got.Lunar.StringValue, got.LunarDate)
	}
	if got.SolarTerm == "" {
		t.Errorf("SolarTerm should not be empty")
	}
	if got.SolarTerm != got.SolarTermDetail.Name {
		t.Errorf("SolarTerm = %q; want %q", got.SolarTerm, got.SolarTermDetail.Name)
	}
	if got.Pillars.Year == "" || got.Pillars.Month == "" || got.Pillars.Day == "" || got.Pillars.Hour == "" {
		t.Errorf("pillars should be rendered as strings: %+v", got.Pillars)
	}
	if got.Pillars.Year != zodiac.GetStemBranchName(res.Pillars.Year.StemIndex, res.Pillars.Year.BranchIndex) {
		t.Errorf("Pillars.Year = %q; want %q", got.Pillars.Year, res.Pillars.Year.String())
	}
}
