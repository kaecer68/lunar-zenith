package zodiac

import (
	"testing"
)

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
