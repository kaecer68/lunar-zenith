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

// GetTwentyEightMansion 根據日柱干支計算當日值日星宿
// 算法：以日支為基準，配合節氣月份推算
// 正月起角，二月起奎，三月起胃，四月起畢，五月起參，六月起鬼
// 七月起張，八月起角，九月起奎，十月起胃，十一月起畢，十二月起參
// 每月30天，按日支順推
func GetTwentyEightMansion(month int, dayStem int, dayBranch int) MansionInfo {
	// 各月起始星宿索引
	monthStartMansion := []int{
		0,  // 正月 - 角
		14, // 二月 - 奎
		16, // 三月 - 胃
		18, // 四月 - 畢
		20, // 五月 - 參
		22, // 六月 - 鬼
		25, // 七月 - 張
		0,  // 八月 - 角
		14, // 九月 - 奎
		16, // 十月 - 胃
		18, // 十一月 - 畢
		20, // 十二月 - 參
	}

	// 計算偏移：使用地支索引
	startIdx := monthStartMansion[(month-1)%12]
	offset := dayBranch
	mansionIdx := (startIdx + offset) % 28

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

// ═══════════════════════════════════════════════════════════
// 胎神 (Fetal God)
// ═══════════════════════════════════════════════════════════

// FetalGodInfo 胎神資訊
type FetalGodInfo struct {
	Position    string // 胎神位置
	Description string // 詳細說明
	Taboo       string // 禁忌事項
}

// GetFetalGod 根據日干計算當日胎神位置
// 算法：甲己之日占門房，乙庚碓磨莫移方
//
//	丙辛廚灶莫相干，丁壬倉庫忌修弄
//	戊癸房床若移動，不招瘟疫也瘟癀
func GetFetalGod(dayStem int) FetalGodInfo {
	positions := []struct {
		Pos  string
		Desc string
		Tabu string
	}{
		{"門外東南", "甲日胎神在門外東南", "忌修門、動土"},         // 甲
		{"碓磨廁外東南", "乙日胎神在碓磨廁外東南", "忌移動碓磨、修廁"},   // 乙
		{"廚灶爐外正南", "丙日胎神在廚灶爐外正南", "忌修廚灶、動爐"},    // 丙
		{"倉庫廁外東南", "丁日胎神在倉庫廁外東南", "忌修倉庫、動廁"},    // 丁
		{"房床廁外正南", "戊日胎神在房床廁外正南", "忌移動房床、修廁"},   // 戊
		{"占門床外正南", "己日胎神在占門床外正南", "忌修門、移動床"},    // 己
		{"碓磨廁外正南", "庚日胎神在碓磨廁外正南", "忌移動碓磨、修廁"},   // 庚
		{"廚灶廁外西南", "辛日胎神在廚灶廁外西南", "忌修廚灶、動廁"},    // 辛
		{"倉庫雞棲外東南", "壬日胎神在倉庫雞棲外東南", "忌修倉庫、動雞棲"}, // 壬
		{"占門床外東南", "癸日胎神在占門床外東南", "忌修門、移動床"},    // 癸
	}

	p := positions[dayStem%10]
	return FetalGodInfo{
		Position:    p.Pos,
		Description: p.Desc,
		Taboo:       p.Tabu,
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
// 算法：子午相沖、丑未相沖、寅申相沖、卯酉相沖、辰戌相沖、巳亥相沖
//
//	沖即為煞位，如子午相沖，午日沖子(鼠)，煞在北方
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

	// 地支對應方向
	branchDirection := map[int]string{
		0:  "北",  // 子
		1:  "東北", // 丑
		2:  "東北", // 寅
		3:  "東",  // 卯
		4:  "東南", // 辰
		5:  "東南", // 巳
		6:  "南",  // 午
		7:  "西南", // 未
		8:  "西南", // 申
		9:  "西",  // 酉
		10: "西北", // 戌
		11: "西北", // 亥
	}

	clashIdx := clashMap[dayBranch]
	clashBranch := EarthlyBranches[clashIdx]
	clashAnimal := ZodiacAnimals[clashIdx]
	shaDir := branchDirection[clashIdx]

	return ClashShaInfo{
		ClashZodiac:  "沖" + clashAnimal,
		ClashBranch:  clashBranch,
		ShaDirection: "煞" + shaDir,
		ShaDesc:      clashBranch + "方(" + shaDir + ")諸事不宜",
	}
}
