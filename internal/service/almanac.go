package service

import "github.com/kaecer68/lunar-zenith/pkg/zodiac"

// Almanac method note:
// Current implementation uses a deterministic baseline model:
// 1) derive Twelve Officer (建除十二神) from month/day branch,
// 2) map officer -> suitable/avoidable via static table,
// 3) map day stem -> deity directions via static table.
//
// This is intentionally a stable and explainable subset. It is NOT a full
// multi-factor traditional tongshu engine (which may combine additional
// shensha layers, conflict-resolution priorities, and school-specific rules).
// If expanding rules, keep provenance and precedence explicit in SKILLS.md.

// Directions 吉神方位
type Directions struct {
	Wealth  string `json:"wealth"`  // 財神方位
	Fortune string `json:"fortune"` // 福神方位
	Study   string `json:"study"`   // 文曲方位
	Love    string `json:"love"`    // 喜神方位
}

// AlmanacEntry 黃曆宜忌條目
type AlmanacEntry struct {
	Suitable  []string `json:"suitable"`  // 宜
	Avoidable []string `json:"avoidable"` // 忌
}

// AlmanacRuleSource 記錄規則來源，供後續維運與對拍追蹤。
type AlmanacRuleSource struct {
	Layer      string `json:"layer"`      // 規則層（base/shensha/conflict）
	RuleID     string `json:"rule_id"`    // 規則識別
	Provenance string `json:"provenance"` // 來源註記（書目/版本）
	Note       string `json:"note"`       // 補充說明
}

// AlmanacRuleHit 是單一事項的判定記錄。
type AlmanacRuleHit struct {
	Action  string            `json:"action"`  // 事項（例如：嫁娶）
	Verdict string            `json:"verdict"` // suitable/avoidable
	Source  AlmanacRuleSource `json:"source"`
}

// AlmanacDecision 是可追溯的宜忌決策結果。
type AlmanacDecision struct {
	Suitable   []string         `json:"suitable"`
	Avoidable  []string         `json:"avoidable"`
	Directions Directions       `json:"directions"`
	Hits       []AlmanacRuleHit `json:"hits"`
}

// AlmanacContext 控制宜忌決策層的啟用狀態與口徑。
type AlmanacContext struct {
	MonthBranch       int    // 月支索引（0-11）
	DayBranch         int    // 日支索引（0-11）
	EnableShenshaV1   bool   // 是否啟用神煞層 v1
	EnableShenshaV2   bool   // 是否啟用神煞層 v2（主流對齊增強）
	ReferenceProfile  string // 來源口徑標記
	ConflictPolicyTag string // 衝突仲裁標記
}

// BaselineAlmanacContext 回傳與既有版本相容的基線口徑（不啟用神煞覆寫）。
func BaselineAlmanacContext() AlmanacContext {
	return AlmanacContext{
		MonthBranch:       -1,
		DayBranch:         -1,
		EnableShenshaV1:   false,
		EnableShenshaV2:   false,
		ReferenceProfile:  "Baseline-v1",
		ConflictPolicyTag: "avoid-wins",
	}
}

// MainstreamAlmanacContextV2 回傳主流對齊第二版口徑（啟用神煞層 v1+v2）。
func MainstreamAlmanacContextV2(monthBranch int, dayBranch int) AlmanacContext {
	return AlmanacContext{
		MonthBranch:       monthBranch,
		DayBranch:         dayBranch,
		EnableShenshaV1:   true,
		EnableShenshaV2:   true,
		ReferenceProfile:  "Tongshu-Mainstream-v2",
		ConflictPolicyTag: "layered-mainstream-v2",
	}
}

// MainstreamAlmanacContextV1 保留相容；內部映射到 v2 口徑。
func MainstreamAlmanacContextV1(monthBranch int, dayBranch int) AlmanacContext {
	return MainstreamAlmanacContextV2(monthBranch, dayBranch)
}

// AlmanacActivities 建除十二神對應的宜忌活動（目前為主決策表）
var AlmanacActivities = map[string]AlmanacEntry{
	"建": {
		Suitable:  []string{"出行", "上任", "會友", "上樑"},
		Avoidable: []string{"嫁娶", "開倉", "安葬"},
	},
	"除": {
		Suitable:  []string{"沐浴", "清潔", "療病", "出行"},
		Avoidable: []string{"嫁娶", "安床", "開市"},
	},
	"滿": {
		Suitable:  []string{"嫁娶", "祈福", "移徙", "入宅", "開市"},
		Avoidable: []string{"安葬", "出行", "求醫"},
	},
	"平": {
		Suitable:  []string{"修造", "動土", "安床", "修飾垣墻"},
		Avoidable: []string{"嫁娶", "開市", "出行"},
	},
	"定": {
		Suitable:  []string{"嫁娶", "祭祀", "祈福", "求嗣"},
		Avoidable: []string{"出行", "上任", "交易"},
	},
	"執": {
		Suitable:  []string{"祭祀", "祈福", "捕捉"},
		Avoidable: []string{"嫁娶", "安床", "入宅"},
	},
	"破": {
		Suitable:  []string{"破屋", "壞垣", "祛病", "解除"},
		Avoidable: []string{"嫁娶", "開市", "立約", "入宅"},
	},
	"危": {
		Suitable:  []string{"安床", "修飾垣墻", "拆卸"},
		Avoidable: []string{"嫁娶", "出行", "開市", "交易"},
	},
	"成": {
		Suitable:  []string{"嫁娶", "祭祀", "祈福", "開市", "交易", "立約"},
		Avoidable: []string{"安葬", "動土"},
	},
	"收": {
		Suitable:  []string{"祭祀", "祈福", "嫁娶", "修造", "捕捉"},
		Avoidable: []string{"開倉", "出財", "安葬", "出行"},
	},
	"開": {
		Suitable:  []string{"開市", "交易", "立約", "安機械", "出行"},
		Avoidable: []string{"安葬", "動土", "嫁娶"},
	},
	"閉": {
		Suitable:  []string{"祭祀", "祈福", "安葬", "修墳"},
		Avoidable: []string{"開市", "出行", "嫁娶", "入宅"},
	},
}

// DayOfficers 十二值日星（建除十二神）
var DayOfficers = []string{"建", "除", "滿", "平", "定", "執", "破", "危", "成", "收", "開", "閉"}

// GetAlmanacByOfficer 根據建除十二神獲取宜忌
func GetAlmanacByOfficer(officer string) AlmanacEntry {
	if entry, ok := AlmanacActivities[officer]; ok {
		return entry
	}
	return AlmanacEntry{
		Suitable:  []string{},
		Avoidable: []string{},
	}
}

// GetDeityDirections 根據日干獲取吉神方位
func GetDeityDirections(dayStemIndex int) Directions {
	// 財神方位：根據日干
	wealthDirections := []string{"東", "東南", "南", "東南", "東", "東北", "西", "西南", "北", "西北"}

	// 喜神方位：根據日干
	loveDirections := []string{"東北", "西北", "南", "東南", "東北", "西北", "南", "東南", "東北", "西北"}

	// 福神方位：相對固定，根據日干微調
	fortuneDirections := []string{"東南", "東", "西", "西南", "東南", "東", "西", "西南", "東南", "東"}

	// 文曲方位：根據日干
	studyDirections := []string{"北", "東南", "東", "西南", "北", "東南", "東", "西南", "北", "東南"}

	return Directions{
		Wealth:  wealthDirections[dayStemIndex%10],
		Love:    loveDirections[dayStemIndex%10],
		Fortune: fortuneDirections[dayStemIndex%10],
		Study:   studyDirections[dayStemIndex%10],
	}
}

// CalculateAlmanac 計算黃曆宜忌和吉神方位
// 注意：目前僅以建除主表 + 日干方位推導；其他神煞資訊不直接改寫宜忌列表。
func CalculateAlmanac(officer string, dayStemIndex int) (suitable, avoidable []string, directions Directions) {
	decision := CalculateAlmanacDetailed(officer, dayStemIndex)
	return decision.Suitable, decision.Avoidable, decision.Directions
}

// CalculateAlmanacDetailed 計算可追溯的宜忌決策。
// 目前採 Baseline v1：
// - base layer: 建除主表
// - shensha layer: 保留擴充入口（預設不覆寫）
// - conflict layer: 同事項衝突時，忌優先於宜
func CalculateAlmanacDetailed(officer string, dayStemIndex int) AlmanacDecision {
	return CalculateAlmanacDetailedWithContext(officer, dayStemIndex, BaselineAlmanacContext())
}

// CalculateAlmanacWithContext 計算含規則上下文的宜忌結果（相容舊回傳格式）。
func CalculateAlmanacWithContext(officer string, dayStemIndex int, ctx AlmanacContext) (suitable, avoidable []string, directions Directions) {
	decision := CalculateAlmanacDetailedWithContext(officer, dayStemIndex, ctx)
	return decision.Suitable, decision.Avoidable, decision.Directions
}

// CalculateAlmanacDetailedWithContext 計算可追溯決策（含規則層上下文）。
func CalculateAlmanacDetailedWithContext(officer string, dayStemIndex int, ctx AlmanacContext) AlmanacDecision {
	directions := GetDeityDirections(dayStemIndex)
	hits := buildBaseOfficerHits(officer)
	hits = applyShenshaOverrides(hits, ctx, officer)
	selectedHits, suitable, avoidable := resolveAlmanacConflicts(hits)

	return AlmanacDecision{
		Suitable:   suitable,
		Avoidable:  avoidable,
		Directions: directions,
		Hits:       selectedHits,
	}
}

func buildBaseOfficerHits(officer string) []AlmanacRuleHit {
	entry := GetAlmanacByOfficer(officer)
	hits := make([]AlmanacRuleHit, 0, len(entry.Suitable)+len(entry.Avoidable))

	for _, action := range entry.Suitable {
		hits = append(hits, AlmanacRuleHit{
			Action:  action,
			Verdict: "suitable",
			Source: AlmanacRuleSource{
				Layer:      "base",
				RuleID:     "officer:" + officer,
				Provenance: "Baseline-v1/TwelveOfficerTable",
				Note:       "建除主表直接映射",
			},
		})
	}

	for _, action := range entry.Avoidable {
		hits = append(hits, AlmanacRuleHit{
			Action:  action,
			Verdict: "avoidable",
			Source: AlmanacRuleSource{
				Layer:      "base",
				RuleID:     "officer:" + officer,
				Provenance: "Baseline-v1/TwelveOfficerTable",
				Note:       "建除主表直接映射",
			},
		})
	}

	return hits
}

func applyShenshaOverrides(hits []AlmanacRuleHit, ctx AlmanacContext, officer string) []AlmanacRuleHit {
	if !ctx.EnableShenshaV1 || ctx.DayBranch < 0 || ctx.DayBranch > 11 {
		return hits
	}

	deity := zodiac.GetDailyDeity(ctx.DayBranch)
	if ctx.MonthBranch >= 0 && ctx.MonthBranch <= 11 {
		deity = zodiac.GetDailyDeityByMonthBranch(ctx.MonthBranch, ctx.DayBranch)
	}
	if deity.Type == "凶" {
		for _, action := range []string{"嫁娶", "開市", "交易", "入宅"} {
			hits = append(hits, AlmanacRuleHit{
				Action:  action,
				Verdict: "avoidable",
				Source: AlmanacRuleSource{
					Layer:      "shensha",
					RuleID:     "daily-deity:type:xiong",
					Provenance: ctx.ReferenceProfile + "/DailyDeityType",
					Note:       "值神為凶，收斂為保守避事",
				},
			})
		}
	}

	if deity.Type == "吉" {
		for _, action := range []string{"祈福", "求嗣"} {
			hits = append(hits, AlmanacRuleHit{
				Action:  action,
				Verdict: "suitable",
				Source: AlmanacRuleSource{
					Layer:      "shensha",
					RuleID:     "daily-deity:type:ji",
					Provenance: ctx.ReferenceProfile + "/DailyDeityType",
					Note:       "值神為吉，增加常見吉事（不直接推祭祀）",
				},
			})
		}
	}

	if ctx.EnableShenshaV2 {
		hits = applyMainstreamV2Overrides(hits, ctx, officer)
	}

	return hits
}

func applyMainstreamV2Overrides(hits []AlmanacRuleHit, ctx AlmanacContext, officer string) []AlmanacRuleHit {
	type officerPack struct {
		positive []string
		negative []string
		neutral  []string
		note     string
	}

	packs := map[string]officerPack{
		"開": {
			positive: []string{"嫁娶", "開市", "交易", "立券交易", "入宅", "移徙", "安床", "動土", "破土", "拆卸", "開光", "祈福", "求嗣"},
			negative: []string{"祭祀", "探病", "入殮"},
			neutral:  []string{"交易", "立約", "安機械"},
			note:     "開日主流對齊增強",
		},
		"成": {
			positive: []string{"嫁娶", "祭祀", "祈福", "開市", "交易", "立約", "入宅", "安床", "出行"},
			negative: []string{"安葬", "探病"},
			note:     "成日主流對齊增強",
		},
		"滿": {
			positive: []string{"嫁娶", "祈福", "移徙", "入宅", "開市", "安床", "求嗣"},
			negative: []string{"安葬", "探病", "入殮"},
			note:     "滿日主流對齊增強",
		},
	}

	pack, ok := packs[officer]
	if !ok {
		return hits
	}

	for _, action := range pack.positive {
		hits = append(hits, AlmanacRuleHit{
			Action:  action,
			Verdict: "suitable_strong",
			Source: AlmanacRuleSource{
				Layer:      "shensha",
				RuleID:     "mainstream-v2:officer:" + officer + ":positive-pack",
				Provenance: ctx.ReferenceProfile + "/" + officer + "DayPack",
				Note:       pack.note + "：提升常見民用吉事",
			},
		})
	}

	for _, action := range pack.negative {
		hits = append(hits, AlmanacRuleHit{
			Action:  action,
			Verdict: "avoidable_strong",
			Source: AlmanacRuleSource{
				Layer:      "shensha",
				RuleID:     "mainstream-v2:officer:" + officer + ":negative-pack",
				Provenance: ctx.ReferenceProfile + "/" + officer + "DayPack",
				Note:       pack.note + "：收斂常見忌事",
			},
		})
	}

	for _, action := range pack.neutral {
		hits = append(hits, AlmanacRuleHit{
			Action:  action,
			Verdict: "neutral_strong",
			Source: AlmanacRuleSource{
				Layer:      "shensha",
				RuleID:     "mainstream-v2:officer:" + officer + ":neutral-pack",
				Provenance: ctx.ReferenceProfile + "/" + officer + "DayPack",
				Note:       pack.note + "：移除非目標清單事項",
			},
		})
	}

	return hits
}

func resolveAlmanacConflicts(hits []AlmanacRuleHit) ([]AlmanacRuleHit, []string, []string) {
	type selected struct {
		hit   AlmanacRuleHit
		score int
	}

	selectedByAction := make(map[string]selected, len(hits))
	actionOrder := make([]string, 0, len(hits))

	for _, hit := range hits {
		score := verdictPriority(hit.Verdict)
		if cur, ok := selectedByAction[hit.Action]; ok {
			if score > cur.score {
				selectedByAction[hit.Action] = selected{hit: hit, score: score}
			}
			continue
		}

		actionOrder = append(actionOrder, hit.Action)
		selectedByAction[hit.Action] = selected{hit: hit, score: score}
	}

	selectedHits := make([]AlmanacRuleHit, 0, len(actionOrder))
	suitable := make([]string, 0, len(actionOrder))
	avoidable := make([]string, 0, len(actionOrder))

	for _, action := range actionOrder {
		sel := selectedByAction[action]
		selectedHits = append(selectedHits, sel.hit)
		if sel.hit.Verdict == "neutral_strong" {
			continue
		}
		if sel.hit.Verdict == "avoidable" || sel.hit.Verdict == "avoidable_strong" {
			avoidable = append(avoidable, action)
			continue
		}
		suitable = append(suitable, action)
	}

	return selectedHits, suitable, avoidable
}

func verdictPriority(verdict string) int {
	switch verdict {
	case "avoidable_strong":
		return 4
	case "neutral_strong":
		return 3
	case "suitable_strong":
		return 2
	case "avoidable":
		return 1
	case "suitable":
		return 0
	default:
		return -1
	}
}

// GetDayOfficerName 獲取十二值日星名稱
func GetDayOfficerName(index int) string {
	if index >= 0 && index < len(DayOfficers) {
		return DayOfficers[index]
	}
	return "未知"
}

// GetTwelveOfficerEnhanced 獲取十二值日星（增強版，返回索引）
func GetTwelveOfficerEnhanced(monthBranch, dayBranch int) (string, int) {
	index := ((dayBranch - monthBranch) + 12) % 12
	return DayOfficers[index], index
}

// GetOfficerFromName 根據名稱獲取建除十二神索弓|
func GetOfficerFromName(name string) int {
	for i, officer := range DayOfficers {
		if officer == name {
			return i
		}
	}
	return -1
}
