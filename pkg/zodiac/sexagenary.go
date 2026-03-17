package zodiac

// 天干 (Heavenly Stems)
var HeavenlyStems = []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}

// 地支 (Earthly Branches)
var EarthlyBranches = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

// 生肖 (Zodiac Animals)
var ZodiacAnimals = []string{"鼠", "牛", "虎", "兔", "龍", "蛇", "馬", "羊", "猴", "雞", "狗", "豬"}

// Sexagenary 封裝一個干支組合
type Sexagenary struct {
	StemIndex   int // 天干索引 (0-9)
	BranchIndex int // 地支索引 (0-11)
}

// String 返回干支的繁體中文名稱 (如: 甲子)
func (s Sexagenary) String() string {
	return HeavenlyStems[s.StemIndex] + EarthlyBranches[s.BranchIndex]
}

// Animal 返回地支對應的生肖
func (s Sexagenary) Animal() string {
	return ZodiacAnimals[s.BranchIndex]
}

// NewYearSexagenary 根據陽曆年份計算年干支
// 基準：西元 4 年是「甲子」年
func NewYearSexagenary(year int) Sexagenary {
	// 偏移量計算，處理西元前年份
	var offset int
	if year >= 0 {
		offset = (year - 4) % 60
	} else {
		offset = (year - 3) % 60
	}
	if offset < 0 {
		offset += 60
	}

	return Sexagenary{
		StemIndex:   offset % 10,
		BranchIndex: offset % 12,
	}
}

// GetHourBranch 根據 24 小時制獲取地支索引
// 23:00 - 01:00 是 子時
func GetHourBranch(hour int) int {
	return ((hour + 1) / 2) % 12
}

// GetDaySexagenary 根據儒略日計算日干支
// 基準：JD 2451545.0 (2000-01-01) 是「甲子」日 (實際上 2000/1/1 是戊午日，這裡需要精確校對基準)
// 修正基準：2000/01/01 是 戊午日 (Stem:4, Branch:6)
// JD 2451545.0 偏移量計算
func GetDaySexagenary(jd float64) Sexagenary {
	// 修正偏移：JD 0 是 癸丑 (實際上干支計日是連續的，不受格里曆變革影響)
	// 已知 2000年1月1日 (JD 2451545) 是 戊午日 (4, 6)
	// 算法：(JD + 0.5 + 基準偏移) % 60
	offset := int(jd+0.5+49) % 60
	if offset < 0 {
		offset += 60
	}
	return Sexagenary{
		StemIndex:   offset % 10,
		BranchIndex: offset % 12,
	}
}

// GetMonthSexagenary 根據年干與月份獲取月干支 (根據五虎遁)
// month: 1 (寅月) to 12 (丑月)
// yearStem: 年干索引 (0:甲, 1:乙...)
func GetMonthSexagenary(yearStem int, month int) Sexagenary {
	// 五虎遁：甲己之年丙作首...
	// 正月(寅月)的天干：(yearStem % 5 * 2 + 2) % 10
	startStem := (yearStem%5*2 + 2) % 10
	monthStem := (startStem + (month - 1)) % 10
	monthBranch := (month + 1) % 12 // 寅月是 2

	return Sexagenary{
		StemIndex:   monthStem,
		BranchIndex: monthBranch,
	}
}

// GetHourSexagenary 根據日干與小時獲取時干支 (根據五鼠遁)
// hourBranch: 地支索引 (0:子, 1:丑...)
// dayStem: 日干索引
func GetHourSexagenary(dayStem int, hourBranch int) Sexagenary {
	// 五鼠遁：甲己還加甲...
	startStem := (dayStem % 5 * 2) % 10
	hourStem := (startStem + hourBranch) % 10

	return Sexagenary{
		StemIndex:   hourStem,
		BranchIndex: hourBranch,
	}
}

// GetStemBranchName 根據天干和地支索引獲取干支名稱
func GetStemBranchName(stemIndex, branchIndex int) string {
	if stemIndex < 0 || stemIndex >= len(HeavenlyStems) {
		return "未知"
	}
	if branchIndex < 0 || branchIndex >= len(EarthlyBranches) {
		return "未知"
	}
	return HeavenlyStems[stemIndex] + EarthlyBranches[branchIndex]
}
