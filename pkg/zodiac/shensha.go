package zodiac

import "fmt"

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
		0: "寅", 4: "寅", 8: "寅", // 申(8), 子(0), 辰(4) -> 寅
		2: "申", 6: "申", 10: "申", // 寅(2), 午(6), 戌(10) -> 申
		5: "亥", 9: "亥", 1: "亥", // 巳(5), 酉(9), 丑(1) -> 亥
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

// ═══════════════════════════════════════════════════════════
// 二十八星宿 (28 Lunar Mansions)
// ═══════════════════════════════════════════════════════════

// TwentyEightMansions 二十八星宿名稱 (東方青龍、南方朱雀、西方白虎、北方玄武)
var TwentyEightMansions = []struct {
	Name    string // 星宿名
	Animal  string // 對應動物 (角木蛟、亢金龍等)
	Palace  string // 宮位 (東方/南方/西方/北方)
	Element string // 五行屬性
}{
	{"角", "蛟", "東方青龍", "木"}, // 0
	{"亢", "龍", "東方青龍", "金"}, // 1
	{"氐", "貉", "東方青龍", "土"}, // 2
	{"房", "兔", "東方青龍", "日"}, // 3
	{"心", "狐", "東方青龍", "月"}, // 4
	{"尾", "虎", "東方青龍", "火"}, // 5
	{"箕", "豹", "東方青龍", "水"}, // 6
	{"斗", "獬", "北方玄武", "木"}, // 7
	{"牛", "牛", "北方玄武", "金"}, // 8
	{"女", "蝠", "北方玄武", "土"}, // 9
	{"虛", "鼠", "北方玄武", "日"}, // 10
	{"危", "燕", "北方玄武", "月"}, // 11
	{"室", "豬", "北方玄武", "火"}, // 12
	{"壁", "獝", "北方玄武", "水"}, // 13
	{"奎", "狼", "西方白虎", "木"}, // 14
	{"婁", "狗", "西方白虎", "金"}, // 15
	{"胃", "雉", "西方白虎", "土"}, // 16
	{"昴", "雞", "西方白虎", "日"}, // 17
	{"畢", "鳥", "西方白虎", "火"}, // 18
	{"觜", "猴", "西方白虎", "火"}, // 19
	{"參", "猿", "西方白虎", "水"}, // 20
	{"井", "犴", "南方朱雀", "木"}, // 21
	{"鬼", "羊", "南方朱雀", "金"}, // 22
	{"柳", "獐", "南方朱雀", "土"}, // 23
	{"星", "馬", "南方朱雀", "日"}, // 24
	{"張", "鹿", "南方朱雀", "月"}, // 25
	{"翼", "蛇", "南方朱雀", "火"}, // 26
	{"軫", "蚓", "南方朱雀", "水"}, // 27
}

// MansionInfo 星宿資訊
type MansionInfo struct {
	Name     string // 星宿名 (角、亢、氐...)
	Animal   string // 對應動物 (蛟、龍、貉...)
	FullName string // 全名 (角木蛟、亢金龍...)
	Palace   string // 宮位 (東方青龍等)
	Element  string // 五行
	Index    int    // 索引 0-27
}

// GetTwentyEightMansion 根據公曆日期計算當日值日星宿。
//
// 採用「公曆實際天數 A」口徑，對應零基索引公式：
// mansionIndex = (A + 24) % 28
// 其中 A 為公曆絕對天數（1-01-01 為 1）。
func GetTwentyEightMansion(year int, month int, day int) MansionInfo {
	a := gregorianAbsoluteDays(year, month, day)
	mansionIdx := (a + 24) % 28

	m := TwentyEightMansions[mansionIdx]
	return MansionInfo{
		Name:     m.Name,
		Animal:   m.Animal,
		FullName: m.Name + m.Element + m.Animal,
		Palace:   m.Palace,
		Element:  m.Element,
		Index:    mansionIdx,
	}
}

func gregorianAbsoluteDays(year int, month int, day int) int {
	monthDays := [...]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	y := year - 1
	days := y*365 + y/4 - y/100 + y/400

	for m := 1; m < month; m++ {
		days += monthDays[m]
	}
	if month > 2 && isGregorianLeapYear(year) {
		days++
	}

	return days + day
}

func isGregorianLeapYear(year int) bool {
	if year%400 == 0 {
		return true
	}
	if year%100 == 0 {
		return false
	}
	return year%4 == 0
}

// sexagenaryDayIndex 將日干支索引（天干 0-9, 地支 0-11）還原為六十甲子序號（0-59）。
func sexagenaryDayIndex(dayStem int, dayBranch int) int {
	stem := ((dayStem % 10) + 10) % 10
	branch := ((dayBranch % 12) + 12) % 12

	for i := 0; i < 60; i++ {
		if i%10 == stem && i%12 == branch {
			return i
		}
	}

	// 理論上干支配對必有解；保底返回地支值以避免 panic。
	return branch
}

// ═══════════════════════════════════════════════════════════
// 值神 (Daily Deity - 十二值神輪值)
// ═══════════════════════════════════════════════════════════

// DailyDeities 十二值神
var DailyDeities = []struct {
	Name string
	Type string // 吉/凶/中
	Desc string
}{
	{"青龍", "吉", "天乙星，天貴星，利有攸往"},
	{"明堂", "吉", "貴人星，明輔星，利見大人"},
	{"天刑", "凶", "黑道，天刑星，利用刑獄"},
	{"朱雀", "凶", "黑道，天訴星，利用公事"},
	{"金匱", "吉", "福德星，月仙星，利釋道用事"},
	{"天德", "吉", "寶光星，天德星，百事吉"},
	{"白虎", "凶", "黑道，天殺星，宜出師遠行"},
	{"玉堂", "吉", "少微星，天開星，百事吉"},
	{"天牢", "凶", "黑道，鎮神星，陰人用事吉"},
	{"玄武", "凶", "黑道，獄星，君子用之吉"},
	{"司命", "吉", "鳳輦星，月仙星，從寅至申時用"},
	{"勾陳", "凶", "黑道，地獄星，起造安葬不利"},
}

// DailyDeityInfo 值神資訊
type DailyDeityInfo struct {
	Name string
	Type string
	Desc string
}

// GetDailyDeity 根據日支計算當日值神
// 算法：日支對應值神
// 子日青龍、丑日明堂、寅日天刑、卯日朱雀、辰日金匱、巳日天德
// 午日白虎、未日玉堂、申日天牢、酉日玄武、戌日司命、亥日勾陳
func GetDailyDeity(dayBranch int) DailyDeityInfo {
	deityOrder := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11} // 子開始
	dIdx := deityOrder[dayBranch%12]
	d := DailyDeities[dIdx]
	return DailyDeityInfo{
		Name: d.Name,
		Type: d.Type,
		Desc: d.Desc,
	}
}

// GetDailyDeityByMonthBranch 根據月支+日支計算當日值神（主流黃道黑道口徑）。
// 口訣：子午月青龍起申，丑未月青龍起戌，寅申月青龍起子，
// 卯酉月青龍起寅，辰戌月青龍起辰，巳亥月青龍起午。
func GetDailyDeityByMonthBranch(monthBranch int, dayBranch int) DailyDeityInfo {
	// 每個月支對應「青龍所在日支」起點。
	qStartByMonthBranch := []int{8, 10, 0, 2, 4, 6, 8, 10, 0, 2, 4, 6}
	start := qStartByMonthBranch[((monthBranch%12)+12)%12]
	idx := (dayBranch - start + 12) % 12
	d := DailyDeities[idx]
	return DailyDeityInfo{
		Name: d.Name,
		Type: d.Type,
		Desc: d.Desc,
	}
}

// ═══════════════════════════════════════════════════════════
// 胎神 (Fetal God)
// ═══════════════════════════════════════════════════════════

// FetalGodInfo 胎神資訊
type FetalGodInfo struct {
	Position    string // 胎神位置
	Description string // 詳細說明
	Taboo       string // 禁忌事項
}

// GetFetalGod 根據日干計算當日胎神位置（相容舊介面，僅保留日干層級）。
// Deprecated: 優先使用 GetFetalGodByDayPillar(dayStem, dayBranch)。
func GetFetalGod(dayStem int) FetalGodInfo {
	return GetFetalGodByDayPillar(dayStem, 0)
}

var fetalGodBySexagenaryDay = []string{
	"占門碓外東南", "碓磨廁外東南", "廚灶爐外正南", "倉庫門外正南", "房床棲外正南", "占門床外正南", "占碓磨外正南", "廚灶廁外西南", "倉庫爐外西南", "房床門外西南",
	"門雞棲外西南", "碓磨床外西南", "廚灶碓外西南", "倉庫廁外西南", "房床爐外正南", "占大門外正南", "碓磨棲外正西", "廚灶床外正西", "倉庫碓外西北", "房床廁外西北",
	"占門爐外西北", "門碓磨外西北", "廚灶棲外西北", "倉庫床外西北", "房床碓外正北", "占門廁外正北", "碓磨爐外正北", "廚灶門外正北", "倉庫棲外正北", "占房床房內北",
	"占門碓房內北", "碓磨廁房內北", "廚灶爐房內北", "門倉庫房內北", "房床棲房內南", "占門床房內南", "占碓磨房內南", "廚灶廁房內南", "倉庫爐房內南", "房床門房內南",
	"門雞棲房內北", "碓磨床房內北", "廚灶碓房內北", "倉庫廁房內北", "房床爐房內北", "占大門外東北", "碓磨棲外東北", "廚灶床外東北", "倉庫碓外東北", "房床廁外東北",
	"占門爐外東北", "門碓磨外正東", "廚灶棲外正東", "倉庫床外正東", "房床碓外正東", "占門廁外正東", "碓磨爐外東南", "廚灶門外東南", "倉庫棲外東南", "占房床外東南",
}

// GetFetalGodByDayPillar 根據日干+日支計算胎神方位（60 甲子定表口徑）。
func GetFetalGodByDayPillar(dayStem int, dayBranch int) FetalGodInfo {
	idx := sexagenaryDayIndex(dayStem, dayBranch)
	pos := fetalGodBySexagenaryDay[idx]
	desc := fmt.Sprintf("胎神占方：%s", pos)
	tabu := fmt.Sprintf("忌在%s方位動土、修造、敲打、搬移重物", pos)

	return FetalGodInfo{
		Position:    pos,
		Description: desc,
		Taboo:       tabu,
	}
}

// ═══════════════════════════════════════════════════════════
// 沖煞 (Clash & Sha)
// ═══════════════════════════════════════════════════════════

// ClashShaInfo 沖煞資訊
type ClashShaInfo struct {
	ClashZodiac  string // 沖生肖 (如：沖猴)
	ClashBranch  string // 沖地支 (如：申)
	ShaDirection string // 煞方向 (如：煞北)
	ShaDesc      string // 煞說明
}

// GetClashSha 根據日支計算當日沖煞
// 算法：
//  1. 子午相沖、丑未相沖、寅申相沖、卯酉相沖、辰戌相沖、巳亥相沖
//  2. 煞方採主流通用四正位映射（依日支）：
//     子辰申 -> 煞南；丑巳酉 -> 煞東；寅午戌 -> 煞北；卯未亥 -> 煞西
func GetClashSha(dayBranch int) ClashShaInfo {
	// 相沖關係：子(0)<->午(6), 丑(1)<->未(7), 寅(2)<->申(8), 卯(3)<->酉(9), 辰(4)<->戌(10), 巳(5)<->亥(11)
	clashMap := map[int]int{
		0: 6, 6: 0, // 子午沖
		1: 7, 7: 1, // 丑未沖
		2: 8, 8: 2, // 寅申沖
		3: 9, 9: 3, // 卯酉沖
		4: 10, 10: 4, // 辰戌沖
		5: 11, 11: 5, // 巳亥沖
	}

	// 主流日支煞方（四正位）
	shaByDayBranch := map[int]string{
		0:  "南", // 子
		1:  "東", // 丑
		2:  "北", // 寅
		3:  "西", // 卯
		4:  "南", // 辰
		5:  "東", // 巳
		6:  "北", // 午
		7:  "西", // 未
		8:  "南", // 申
		9:  "東", // 酉
		10: "北", // 戌
		11: "西", // 亥
	}

	clashIdx := clashMap[dayBranch]
	clashBranch := EarthlyBranches[clashIdx]
	clashAnimal := ZodiacAnimals[clashIdx]
	shaDir := shaByDayBranch[dayBranch]

	return ClashShaInfo{
		ClashZodiac:  "沖" + clashAnimal,
		ClashBranch:  clashBranch,
		ShaDirection: "煞" + shaDir,
		ShaDesc:      shaDir + "方諸事不宜",
	}
}
