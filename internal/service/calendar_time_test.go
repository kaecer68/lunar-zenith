package service

import "testing"

func TestResolveCalendarQueryTimeEmptyDateDefaultsToTaipeiNoon(t *testing.T) {
	got, err := resolveCalendarQueryTime("")
	if err != nil {
		t.Fatalf("resolveCalendarQueryTime(\"\") error = %v", err)
	}
	if got.Location() != taipeiLoc {
		t.Errorf("location = %v; want %v", got.Location(), taipeiLoc)
	}
	if got.Hour() != 12 || got.Minute() != 0 || got.Second() != 0 {
		t.Errorf("time = %02d:%02d:%02d; want 12:00:00", got.Hour(), got.Minute(), got.Second())
	}
}
