package zodiac

// LunarFestival 農曆節日定義
type LunarFestival struct {
	Name        string // 節日名稱
	Month       int    // 農曆月份
	Day         int    // 農曆日期
	Type        string // 類型：道教/佛教/民間/民俗
	Description string // 說明
	Priority    int    // 優先級（數字越大越重要，用於當日多節日時排序）
}

// LunarFestivals 台灣重要農曆宗教節日清單
// 按農曆月份排序，方便查詢
var LunarFestivals = []LunarFestival{
	// 正月
	{"天公生（玉皇大帝聖誕）", 1, 9, "道教", "玉皇大帝萬壽，天界最高神祇聖誕，子時拜天公", 100},
	{"元宵節", 1, 15, "民俗", "上元節，天官賜福，賞花燈、吃湯圓", 90},
	
	// 二月
	{"頭牙（土地公聖誕）", 2, 2, "道教", "福德正神聖誕，土地公生，商家祭拜祈求財運", 80},
	{"媽祖生（媽祖聖誕）", 3, 23, "道教", "天上聖母媽祖聖誕，台灣最重要的民間信仰節日之一", 100},
	
	// 四月
	{"保生大帝聖誕", 4, 15, "道教", "大道公吳夲聖誕，醫神保生大帝萬夽", 85},
	
	// 五月
	{"佛誕節（浴佛節）", 4, 8, "佛教", "釋迦牟尼佛聖誕，浴佛儀式，華人地區多慶祝農曆四月初八", 95},
	
	// 六月
	{"觀世音菩薩成道日", 6, 19, "佛教", "觀世音菩薩成道紀念日", 90},
	
	// 七月
	{"義民節", 7, 20, "道教", "義民爺祭，新竹褒忠亭義民廟重要祭典", 85},
	{"中元祭", 7, 15, "道教", "中元節，地官赦罪，普渡孤魂，鬼月最重要節日", 95},
	
	// 八月
		// 九月
	{"觀世音菩薩聖誕", 9, 19, "佛教", "觀世音菩薩聖誕，信眾朝拜祈求平安", 90},
	
	// 十月
	{"艋舺青山王祭", 10, 22, "道教", "青山靈安尊王聖誕，艋舺青山宮重大祭典", 85},
	
	// 十二月
	{"臘八節（釋迦牟尼佛成道日）", 12, 8, "佛教", "釋迦牟尼佛成道日，喝臘八粥", 85},
	{"尾牙（土地公聖誕）", 12, 16, "道教", "一年最後一次牙祭，土地公聖誕，商家宴請員工", 85},
}

// GetLunarFestival 根據農曆月日取得節日資訊
// 返回匹配的節日列表（可能有多個）
func GetLunarFestival(lunarMonth, lunarDay int) []LunarFestival {
	var result []LunarFestival
	for _, f := range LunarFestivals {
		if f.Month == lunarMonth && f.Day == lunarDay {
			result = append(result, f)
		}
	}
	return result
}

// GetFestivalTypeColor 返回節日類型的顏色標識（供前端使用）
func GetFestivalTypeColor(festivalType string) string {
	switch festivalType {
	case "道教":
		return "#c9a84c" // 金色
	case "佛教":
		return "#9966cc" // 紫色
	case "民間":
		return "#44b0aa" // 青色
	case "民俗":
		return "#4caa7a" // 綠色
	default:
		return "#9090a8" // 灰色
	}
}
