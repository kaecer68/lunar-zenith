package zodiac

import (
	"github.com/kaecer68/lunar-zenith/pkg/celestial"
	"testing"
	"time"
)

func TestGetAstrologicalPillar(t *testing.T) {
	// 2024 年立春約在 2/4 08:27 UTC
	// 1. 立春前：2024-02-04 04:00 UTC
	tBefore := time.Date(2024, 2, 4, 4, 0, 0, 0, time.UTC)
	ptBefore := celestial.NewPrecisionTime(tBefore)
	pillarBefore := GetAstrologicalPillar(ptBefore)

	// 立春前應為 癸卯年 (2023 屬兔)
	if pillarBefore.Year.String() != "癸卯" {
		t.Errorf("立春前應為 癸卯, got %s", pillarBefore.Year.String())
	}

	// 2. 立春後：2024-02-04 12:00 UTC
	tAfter := time.Date(2024, 2, 4, 12, 0, 0, 0, time.UTC)
	ptAfter := celestial.NewPrecisionTime(tAfter)
	pillarAfter := GetAstrologicalPillar(ptAfter)

	// 立春後應為 甲辰年 (2024 屬龍)
	if pillarAfter.Year.String() != "甲辰" {
		t.Errorf("立春後應為 甲辰, got %s", pillarAfter.Year.String())
	}

	// 3. 驗證月份切換
	// 2024-03-05 驚蟄前 (寅月)
	tYin := time.Date(2024, 3, 5, 2, 0, 0, 0, time.UTC)
	pillarYin := GetAstrologicalPillar(celestial.NewPrecisionTime(tYin))
	if pillarYin.Month.BranchIndex != 2 { // 寅是 2
		t.Errorf("驚蟄前應為 寅月, got %s", pillarYin.Month.String())
	}

	// 2024-03-05 驚蟄後 (卯月)
	tMao := time.Date(2024, 3, 5, 12, 0, 0, 0, time.UTC)
	pillarMao := GetAstrologicalPillar(celestial.NewPrecisionTime(tMao))
	if pillarMao.Month.BranchIndex != 3 { // 卯是 3
		t.Errorf("驚蟄後應為 卯月, got %s", pillarMao.Month.String())
	}
}
