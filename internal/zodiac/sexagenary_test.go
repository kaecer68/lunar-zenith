package zodiac

import (
	"testing"
)

func TestNewYearSexagenary(t *testing.T) {
	// 2024 應為 甲辰
	got := NewYearSexagenary(2024).String()
	want := "甲辰"
	if got != want {
		t.Errorf("NewYearSexagenary(2024) = %s; want %s", got, want)
	}

	// 1984 應為 甲子 (基準)
	got2 := NewYearSexagenary(1984).String()
	want2 := "甲子"
	if got2 != want2 {
		t.Errorf("NewYearSexagenary(1984) = %s; want %s", got2, want2)
	}
}

func TestGetDaySexagenary(t *testing.T) {
	// 2000-01-01 (JD 2451545.0) 應為 戊午
	got := GetDaySexagenary(2451545.0).String()
	want := "戊午"
	if got != want {
		t.Errorf("GetDaySexagenary(2451545.0) = %s; want %s", got, want)
	}
}

func TestGetMonthSexagenary(t *testing.T) {
	// 甲年 (Stem 0) 的 1 月 (寅月) 應為 丙寅 (五虎遁)
	got := GetMonthSexagenary(0, 1).String()
	want := "丙寅"
	if got != want {
		t.Errorf("GetMonthSexagenary(甲, 寅月) = %s; want %s", got, want)
	}
}

func TestGetHourSexagenary(t *testing.T) {
	// 甲日 (Stem 0) 的 子時 (Branch 0) 應為 甲子 (五鼠遁)
	got := GetHourSexagenary(0, 0).String()
	want := "甲子"
	if got != want {
		t.Errorf("GetHourSexagenary(甲日, 子時) = %s; want %s", got, want)
	}

	// 戊日 (Stem 4) 的 午時 (Branch 6) 應為 戊午
	got2 := GetHourSexagenary(4, 6).String()
	want2 := "戊午"
	if got2 != want2 {
		t.Errorf("GetHourSexagenary(戊日, 午時) = %s; want %s", got2, want2)
	}
}
