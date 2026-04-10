package zodiac

import (
	"testing"
)

func TestGetTwentyEightMansion_ByGregorianDate(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		month    int
		day      int
		wantFull string
	}{
		{name: "2024-04-04 對應角木蛟", year: 2024, month: 4, day: 4, wantFull: "角木蛟"},
		{name: "2025-05-05 對應心月狐", year: 2025, month: 5, day: 5, wantFull: "心月狐"},
		{name: "2026-06-06 對應女土蝠", year: 2026, month: 6, day: 6, wantFull: "女土蝠"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTwentyEightMansion(tt.year, tt.month, tt.day)
			if got.FullName != tt.wantFull {
				t.Errorf("GetTwentyEightMansion(%d,%d,%d)=%s; want %s", tt.year, tt.month, tt.day, got.FullName, tt.wantFull)
			}
		})
	}
}

func TestGetTwelveOfficer(t *testing.T) {
	// 寅月 (2) 碰到 寅日 (2) 應該是 「建」
	got := GetTwelveOfficer(2, 2)
	want := "建"
	if got != want {
		t.Errorf("GetTwelveOfficer(寅月, 寅日) = %s; want %s", got, want)
	}

	// 寅月 (2) 碰到 卯日 (3) 應該是 「除」
	got2 := GetTwelveOfficer(2, 3)
	want2 := "除"
	if got2 != want2 {
		t.Errorf("GetTwelveOfficer(寅月, 卯日) = %s; want %s", got2, want2)
	}

	// 寅月 (2) 碰到 申日 (8) 應該是 「破」 (沖月建)
	// 建(0) 除(1) 滿(2) 平(3) 定(4) 執(5) 破(6)
	// (8-2) = 6
	got3 := GetTwelveOfficer(2, 8)
	want3 := "破"
	if got3 != want3 {
		t.Errorf("GetTwelveOfficer(寅月, 申日) = %s; want %s", got3, want3)
	}
}

func TestGetYearShenSha(t *testing.T) {
	// 辰年 (4, 龍) 的驛馬應在 寅
	ss := GetYearShenSha(4)
	foundYiMa := false
	for _, s := range ss {
		if s.Name == "年驛馬" && s.Description == "驛馬位在 寅" {
			foundYiMa = true
		}
	}
	if !foundYiMa {
		t.Error("辰年驛馬判定錯誤")
	}

	// 午年 (6, 馬) 的桃花應在 卯
	ss2 := GetYearShenSha(6)
	foundTaoHua := false
	for _, s := range ss2 {
		if s.Name == "年桃花" && s.Description == "桃花位在 卯" {
			foundTaoHua = true
		}
	}
	if !foundTaoHua {
		t.Error("午年桃花判定錯誤")
	}
}

func TestGetClashSha_MainstreamDirection(t *testing.T) {
	tests := []struct {
		name      string
		dayBranch int
		wantClash string
		wantSha   string
	}{
		{name: "寅日沖猴煞北", dayBranch: 2, wantClash: "沖猴", wantSha: "煞北"},
		{name: "子日沖馬煞南", dayBranch: 0, wantClash: "沖馬", wantSha: "煞南"},
		{name: "卯日沖雞煞西", dayBranch: 3, wantClash: "沖雞", wantSha: "煞西"},
		{name: "酉日沖兔煞東", dayBranch: 9, wantClash: "沖兔", wantSha: "煞東"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetClashSha(tt.dayBranch)
			if got.ClashZodiac != tt.wantClash {
				t.Errorf("GetClashSha(%d) clash = %s; want %s", tt.dayBranch, got.ClashZodiac, tt.wantClash)
			}
			if got.ShaDirection != tt.wantSha {
				t.Errorf("GetClashSha(%d) sha = %s; want %s", tt.dayBranch, got.ShaDirection, tt.wantSha)
			}
		})
	}
}

func TestGetDailyDeityByMonthBranch(t *testing.T) {
	tests := []struct {
		name        string
		monthBranch int
		dayBranch   int
		wantName    string
	}{
		{name: "丑月戌日為青龍", monthBranch: 1, dayBranch: 10, wantName: "青龍"},
		{name: "寅月申日為天牢", monthBranch: 2, dayBranch: 8, wantName: "天牢"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDailyDeityByMonthBranch(tt.monthBranch, tt.dayBranch)
			if got.Name != tt.wantName {
				t.Errorf("GetDailyDeityByMonthBranch(%d, %d) = %s; want %s", tt.monthBranch, tt.dayBranch, got.Name, tt.wantName)
			}
		})
	}
}

func TestGetFetalGodByDayPillar(t *testing.T) {
	tests := []struct {
		name      string
		dayStem   int
		dayBranch int
		wantPos   string
	}{
		{name: "戊戌日房床棲房內南", dayStem: 4, dayBranch: 10, wantPos: "房床棲房內南"},
		{name: "己酉日占大門外東北", dayStem: 5, dayBranch: 9, wantPos: "占大門外東北"},
		{name: "辛亥日廚灶床外東北", dayStem: 7, dayBranch: 11, wantPos: "廚灶床外東北"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFetalGodByDayPillar(tt.dayStem, tt.dayBranch)
			if got.Position != tt.wantPos {
				t.Errorf("GetFetalGodByDayPillar(%d,%d)=%s; want %s", tt.dayStem, tt.dayBranch, got.Position, tt.wantPos)
			}
		})
	}
}
