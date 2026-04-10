package service

import (
	"reflect"
	"testing"
)

func TestCalculateAlmanacDetailed_BaselineModel(t *testing.T) {
	decision := CalculateAlmanacDetailed("建", 0)

	wantSuitable := []string{"出行", "上任", "會友", "上樑"}
	wantAvoidable := []string{"嫁娶", "開倉", "安葬"}

	if !reflect.DeepEqual(decision.Suitable, wantSuitable) {
		t.Errorf("Suitable mismatch: got %v want %v", decision.Suitable, wantSuitable)
	}
	if !reflect.DeepEqual(decision.Avoidable, wantAvoidable) {
		t.Errorf("Avoidable mismatch: got %v want %v", decision.Avoidable, wantAvoidable)
	}
	if decision.Directions.Wealth != "東" {
		t.Errorf("Wealth direction mismatch: got %q want %q", decision.Directions.Wealth, "東")
	}
	if len(decision.Hits) != len(wantSuitable)+len(wantAvoidable) {
		t.Errorf("Hits count mismatch: got %d", len(decision.Hits))
	}
}

func TestResolveAlmanacConflicts_AvoidableWins(t *testing.T) {
	hits := []AlmanacRuleHit{
		{
			Action:  "祭祀",
			Verdict: "suitable",
			Source:  AlmanacRuleSource{Layer: "base", RuleID: "base-1"},
		},
		{
			Action:  "祭祀",
			Verdict: "avoidable",
			Source:  AlmanacRuleSource{Layer: "shensha", RuleID: "sha-1"},
		},
		{
			Action:  "出行",
			Verdict: "suitable",
			Source:  AlmanacRuleSource{Layer: "base", RuleID: "base-2"},
		},
	}

	selected, suitable, avoidable := resolveAlmanacConflicts(hits)

	if !reflect.DeepEqual(suitable, []string{"出行"}) {
		t.Errorf("Suitable mismatch: got %v", suitable)
	}
	if !reflect.DeepEqual(avoidable, []string{"祭祀"}) {
		t.Errorf("Avoidable mismatch: got %v", avoidable)
	}

	if len(selected) != 2 {
		t.Fatalf("selected count mismatch: got %d want 2", len(selected))
	}
	for _, hit := range selected {
		if hit.Action == "祭祀" {
			if hit.Verdict != "avoidable" {
				t.Errorf("Conflict verdict mismatch for 祭祀: got %s", hit.Verdict)
			}
			if hit.Source.Layer != "shensha" {
				t.Errorf("Conflict source mismatch for 祭祀: got %s", hit.Source.Layer)
			}
		}
	}
}

func TestCalculateAlmanacDetailedWithContext_ShenshaDisabledKeepsBaseline(t *testing.T) {
	decision := CalculateAlmanacDetailedWithContext("建", 0, BaselineAlmanacContext())

	wantSuitable := []string{"出行", "上任", "會友", "上樑"}
	wantAvoidable := []string{"嫁娶", "開倉", "安葬"}

	if !reflect.DeepEqual(decision.Suitable, wantSuitable) {
		t.Errorf("Suitable mismatch: got %v want %v", decision.Suitable, wantSuitable)
	}
	if !reflect.DeepEqual(decision.Avoidable, wantAvoidable) {
		t.Errorf("Avoidable mismatch: got %v want %v", decision.Avoidable, wantAvoidable)
	}
}

func TestCalculateAlmanacDetailedWithContext_ShenshaEnabledAddsAvoidable(t *testing.T) {
	// 寅日對應天刑（凶），預期神煞層新增保守避事項。
	decision := CalculateAlmanacDetailedWithContext("建", 0, MainstreamAlmanacContextV1(2, 2))

	if !contains(decision.Avoidable, "交易") {
		t.Errorf("Expected 交易 in avoidable, got %v", decision.Avoidable)
	}
	if !contains(decision.Avoidable, "入宅") {
		t.Errorf("Expected 入宅 in avoidable, got %v", decision.Avoidable)
	}

	foundShensha := false
	for _, hit := range decision.Hits {
		if hit.Action == "交易" && hit.Source.Layer == "shensha" {
			foundShensha = true
			break
		}
	}
	if !foundShensha {
		t.Errorf("Expected shensha hit for 交易, got hits=%v", decision.Hits)
	}
}

func TestCalculateAlmanacDetailedWithContext_MainstreamV2KaiDayOverrides(t *testing.T) {
	decision := CalculateAlmanacDetailedWithContext("開", 0, MainstreamAlmanacContextV1(2, 2))

	// v2 強覆寫應使開日常見民用事項可被提升為宜。
	for _, action := range []string{"嫁娶", "開市", "入宅", "動土", "立券交易", "開光"} {
		if !contains(decision.Suitable, action) {
			t.Errorf("Expected %s in suitable, got %v", action, decision.Suitable)
		}
	}

	for _, action := range []string{"交易", "立約", "安機械"} {
		if contains(decision.Suitable, action) || contains(decision.Avoidable, action) {
			t.Errorf("Expected %s suppressed, got suitable=%v avoidable=%v", action, decision.Suitable, decision.Avoidable)
		}
	}

	for _, action := range []string{"祭祀", "探病", "入殮", "安葬"} {
		if !contains(decision.Avoidable, action) {
			t.Errorf("Expected %s in avoidable, got %v", action, decision.Avoidable)
		}
	}
}

func TestResolveAlmanacConflicts_SuitableStrongBeatsAvoidable(t *testing.T) {
	hits := []AlmanacRuleHit{
		{
			Action:  "動土",
			Verdict: "avoidable",
			Source:  AlmanacRuleSource{Layer: "base", RuleID: "base-avoid"},
		},
		{
			Action:  "動土",
			Verdict: "suitable_strong",
			Source:  AlmanacRuleSource{Layer: "shensha", RuleID: "v2-pack"},
		},
	}

	_, suitable, avoidable := resolveAlmanacConflicts(hits)
	if !contains(suitable, "動土") {
		t.Errorf("Expected 動土 in suitable, got suitable=%v avoidable=%v", suitable, avoidable)
	}
	if contains(avoidable, "動土") {
		t.Errorf("Did not expect 動土 in avoidable, got suitable=%v avoidable=%v", suitable, avoidable)
	}
}

func TestResolveAlmanacConflicts_NeutralStrongSuppressesAction(t *testing.T) {
	hits := []AlmanacRuleHit{
		{
			Action:  "立約",
			Verdict: "suitable",
			Source:  AlmanacRuleSource{Layer: "base", RuleID: "base-suit"},
		},
		{
			Action:  "立約",
			Verdict: "neutral_strong",
			Source:  AlmanacRuleSource{Layer: "shensha", RuleID: "v2-neutral"},
		},
	}

	_, suitable, avoidable := resolveAlmanacConflicts(hits)
	if contains(suitable, "立約") || contains(avoidable, "立約") {
		t.Errorf("Expected 立約 suppressed, got suitable=%v avoidable=%v", suitable, avoidable)
	}
}

func TestCalculateAlmanacDetailedWithContext_MainstreamV2ChengDayPack(t *testing.T) {
	decision := CalculateAlmanacDetailedWithContext("成", 0, MainstreamAlmanacContextV2(2, 2))

	for _, action := range []string{"嫁娶", "開市", "交易", "立約", "祈福"} {
		if !contains(decision.Suitable, action) {
			t.Errorf("Expected %s in suitable for 成日, got %v", action, decision.Suitable)
		}
	}
	for _, action := range []string{"安葬", "探病"} {
		if !contains(decision.Avoidable, action) {
			t.Errorf("Expected %s in avoidable for 成日, got %v", action, decision.Avoidable)
		}
	}
}

func TestCalculateAlmanacDetailedWithContext_MainstreamV2ManDayPack(t *testing.T) {
	decision := CalculateAlmanacDetailedWithContext("滿", 0, MainstreamAlmanacContextV2(2, 2))

	for _, action := range []string{"嫁娶", "祈福", "入宅", "安床"} {
		if !contains(decision.Suitable, action) {
			t.Errorf("Expected %s in suitable for 滿日, got %v", action, decision.Suitable)
		}
	}
	for _, action := range []string{"安葬", "探病", "入殮"} {
		if !contains(decision.Avoidable, action) {
			t.Errorf("Expected %s in avoidable for 滿日, got %v", action, decision.Avoidable)
		}
	}
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
