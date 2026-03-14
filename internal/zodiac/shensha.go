package zodiac

// TwelveOfficers 標注建除十二神
var TwelveOfficers = []string{
	"建", "除", "滿", "平", "定", "執",
	"破", "危", "成", "收", "開", "閉",
}

// GetTwelveOfficer 計算給定月份地支與日期地支的「建除十二神」
// monthBranch: 月份地支索引 (0:子, 1:丑, 2:寅...)
// dayBranch: 日期地支索引
func GetTwelveOfficer(monthBranch int, dayBranch int) string {
	// 算法：以月建為「建」。例如寅月(2)碰到寅日(2)，offset 為 0，即「建」
	offset := (dayBranch - monthBranch) % 12
	if offset < 0 {
		offset += 12
	}
	return TwelveOfficers[offset]
}

// CommonShenSha 常用神煞結構
type CommonShenSha struct {
	Name        string
	Description string
}

// GetYearShenSha 獲取基於年支的常用神煞 (如：生肖、驛馬、桃花)
func GetYearShenSha(yearBranch int) []CommonShenSha {
	res := []CommonShenSha{}
	
	// 範例 1: 驛馬 (Yi Ma)
	// 申子辰馬在寅，寅午戌馬在申，巳酉丑馬在亥，亥卯未馬在巳
	yiMaMap := map[int]string{
		0: "寅", 4: "寅", 8: "寅",  // 申(8), 子(0), 辰(4) -> 寅
		2: "申", 6: "申", 10: "申", // 寅(2), 午(6), 戌(10) -> 申
		5: "亥", 9: "亥", 1: "亥",  // 巳(5), 酉(9), 丑(1) -> 亥
		11: "巳", 3: "巳", 7: "巳", // 亥(11), 卯(3), 未(7) -> 巳
	}
	if b, ok := yiMaMap[yearBranch]; ok {
		res = append(res, CommonShenSha{Name: "年驛馬", Description: "驛馬位在 " + b})
	}
	
	// 範例 2: 桃花 (Peach Blossom)
	// 亥卯未見子，申子辰見酉，寅午戌見卯，巳酉丑見午
	taoHuaMap := map[int]string{
		11: "子", 3: "子", 7: "子",
		8: "酉", 0: "酉", 4: "酉",
		2: "卯", 6: "卯", 10: "卯",
		5: "午", 9: "午", 1: "午",
	}
	if b, ok := taoHuaMap[yearBranch]; ok {
		res = append(res, CommonShenSha{Name: "年桃花", Description: "桃花位在 " + b})
	}

	return res
}
