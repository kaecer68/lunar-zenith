package service

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
	"time"
)

type almanacParityFixture struct {
	Name              string
	Date              string
	ExpectedOfficer   string
	ExpectedSuitable  []string
	ExpectedAvoidable []string
	ExpectedClash     string
	ExpectedSha       string
}

// 主流口徑對拍樣本（可逐步擴充到 30 天以上）。
// 預設僅輸出差異；若設置 STRICT_TONGSHU_PARITY=1 則差異會使測試失敗。
var mainstreamParityFixtures = []almanacParityFixture{
	{
		Name:              "gemini-sample-2026-04-10",
		Date:              "2026-04-10",
		ExpectedSuitable:  []string{"嫁娶", "祈福", "求嗣", "開光", "出行", "開市", "立券交易", "入宅", "移徙", "安床", "動土", "破土", "拆卸"},
		ExpectedAvoidable: []string{"祭祀", "探病", "入殮", "安葬"},
		ExpectedClash:     "沖猴",
		ExpectedSha:       "煞北",
	},
}

func TestAlmanacMainstreamParityFixtures(t *testing.T) {
	strict := os.Getenv("STRICT_TONGSHU_PARITY") == "1"
	extraPath := os.Getenv("ALMANAC_PARITY_FIXTURES_JSON")
	fixtures, err := loadAllParityFixtures(extraPath)
	if err != nil {
		t.Fatalf("load parity fixtures failed: %v", err)
	}

	agg := NewAggregator(nil, nil)
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		t.Fatalf("load timezone failed: %v", err)
	}

	for _, fx := range fixtures {
		fx := fx
		t.Run(fx.Name, func(t *testing.T) {
			d, err := time.ParseInLocation("2006-01-02", fx.Date, loc)
			if err != nil {
				t.Fatalf("parse date failed: %v", err)
			}
			r := agg.GetCalendar(time.Date(d.Year(), d.Month(), d.Day(), 12, 0, 0, 0, loc))

			missingSuitable := diffContains(fx.ExpectedSuitable, r.Suitable)
			missingAvoidable := diffContains(fx.ExpectedAvoidable, r.Avoidable)
			extraSuitable := diffContains(r.Suitable, fx.ExpectedSuitable)
			extraAvoidable := diffContains(r.Avoidable, fx.ExpectedAvoidable)

			var problems []string
			if len(missingSuitable) > 0 {
				problems = append(problems, "missing suitable: "+strings.Join(missingSuitable, ", "))
			}
			if len(missingAvoidable) > 0 {
				problems = append(problems, "missing avoidable: "+strings.Join(missingAvoidable, ", "))
			}
			if len(extraSuitable) > 0 {
				problems = append(problems, "unexpected suitable: "+strings.Join(extraSuitable, ", "))
			}
			if len(extraAvoidable) > 0 {
				problems = append(problems, "unexpected avoidable: "+strings.Join(extraAvoidable, ", "))
			}
			if fx.ExpectedOfficer != "" && r.TwelveOfficer != fx.ExpectedOfficer {
				problems = append(problems, "officer mismatch: got "+r.TwelveOfficer+" want "+fx.ExpectedOfficer)
			}
			if fx.ExpectedClash != "" && r.ClashSha.ClashZodiac != fx.ExpectedClash {
				problems = append(problems, "clash mismatch: got "+r.ClashSha.ClashZodiac+" want "+fx.ExpectedClash)
			}
			if fx.ExpectedSha != "" && r.ClashSha.ShaDirection != fx.ExpectedSha {
				problems = append(problems, "sha mismatch: got "+r.ClashSha.ShaDirection+" want "+fx.ExpectedSha)
			}

			if len(problems) == 0 {
				return
			}

			t.Logf("parity diff (%s): %s", fx.Date, strings.Join(problems, " | "))
			t.Logf("actual suitable=%v", r.Suitable)
			t.Logf("actual avoidable=%v", r.Avoidable)
			t.Logf("actual clash=%s sha=%s", r.ClashSha.ClashZodiac, r.ClashSha.ShaDirection)
			if strict {
				t.Fatalf("strict parity enabled and fixture mismatch found")
			}
		})
	}
}

func diffContains(expect []string, actual []string) []string {
	actualSet := make(map[string]struct{}, len(actual))
	for _, a := range actual {
		actualSet[a] = struct{}{}
	}
	missing := make([]string, 0)
	for _, e := range expect {
		if _, ok := actualSet[e]; !ok {
			missing = append(missing, e)
		}
	}
	sort.Strings(missing)
	return missing
}

func loadAllParityFixtures(extraPath string) ([]almanacParityFixture, error) {
	fixtures := make([]almanacParityFixture, 0, len(mainstreamParityFixtures))
	fixtures = append(fixtures, mainstreamParityFixtures...)

	if extraPath == "" {
		return fixtures, nil
	}

	extra, err := loadParityFixturesFromJSON(extraPath)
	if err != nil {
		return nil, err
	}
	fixtures = append(fixtures, extra...)
	return fixtures, nil
}

func loadParityFixturesFromJSON(path string) ([]almanacParityFixture, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read fixtures file %s: %w", path, err)
	}

	var fixtures []almanacParityFixture
	if err := json.Unmarshal(raw, &fixtures); err != nil {
		return nil, fmt.Errorf("parse fixtures file %s: %w", path, err)
	}
	return fixtures, nil
}
