package service

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

// AlmanacActivities 建除十二神對應的宜忌活動
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
func CalculateAlmanac(officer string, dayStemIndex int) (suitable, avoidable []string, directions Directions) {
	almanac := GetAlmanacByOfficer(officer)
	directions = GetDeityDirections(dayStemIndex)
	return almanac.Suitable, almanac.Avoidable, directions
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
