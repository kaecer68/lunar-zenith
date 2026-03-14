package celestial

import (
	"math"
)

const (
	Deg2Rad = math.Pi / 180.0
	Rad2Deg = 180.0 / math.Pi
)

// SolarLongitude 計算給定 JDE 的太陽平黃經（精簡版 VSOP87）
// 精度足以對齊節氣計算需求
func SolarLongitude(jde float64) float64 {
	t := (jde - 2451545.0) / 36525.0 // 儒略世紀數

	// 太陽幾何平黃經 (L0)
	l0 := 280.46646 + 36000.76983*t + 0.0003032*t*t
	// 太陽平近點角 (M)
	m := 357.52911 + 35999.05029*t - 0.0001537*t*t - 0.000000058*t*t*t

	// 修正到 0-360
	l0 = math.Mod(l0, 360.0)
	if l0 < 0 {
		l0 += 360.0
	}
	m_rad := math.Mod(m, 360.0) * Deg2Rad

	// 太陽中心差 (C)
	c := (1.914602-0.004817*t-0.000014*t*t)*math.Sin(m_rad) +
		(0.019993-0.000101*t)*math.Sin(2*m_rad) +
		0.000289*math.Sin(3*m_rad)

	// 太陽真黃經 (Theta)
	theta := l0 + c

	// 考慮章動與光行差 (Apparent Longitude)
	omega := 125.04 - 1934.136*t
	lambda := theta - 0.00569 - 0.00478*math.Sin(omega*Deg2Rad)

	return math.Mod(lambda+360.0, 360.0)
}

// GetSolarTermName 獲取節氣名稱
func GetSolarTermName(index int) string {
	if index < 0 || index >= len(SolarTerms) {
		return "未知"
	}
	return SolarTerms[index]
}

// SolarTermInfo 封裝節氣的完整資訊
type SolarTermInfo struct {
	Index     int     // 節氣索引 (0-23)
	Name      string  // 節氣名稱 (如：春分)
	Longitude float64 // 當前太陽黃經
}

// GetSolarTerm 根據 jde 獲取當前節氣資訊
func GetSolarTerm(jde float64) SolarTermInfo {
	lon := SolarLongitude(jde)
	normLon := math.Mod(lon, 360.0)
	if normLon < 0 {
		normLon += 360.0
	}
	index := int(math.Floor(normLon / 15.0))
	
	return SolarTermInfo{
		Index:     index,
		Name:      GetSolarTermName(index),
		Longitude: lon,
	}
}

// EstimateTermTime 搜尋太陽到達指定黃經 (targetLon) 的精確 JDE
// 使用二分搜尋法在指定範圍 (startJDE, endJDE) 內搜尋，精確度約為 1 秒 (0.00001 days)
func EstimateTermTime(targetLon float64, startJDE, endJDE float64) float64 {
	low := startJDE
	high := endJDE
	const precision = 0.00001 // 約 0.86 秒

	for i := 0; i < 50; i++ { // 限制迭代次數
		mid := (low + high) / 2
		lon := SolarLongitude(mid)
		
		// 處理 360/0 度跨越問題
		diff := lon - targetLon
		if diff > 180 {
			diff -= 360
		} else if diff < -180 {
			diff += 360
		}

		if math.Abs(high-low) < precision {
			return mid
		}

		if diff < 0 {
			low = mid
		} else {
			high = mid
		}
	}
	return (low + high) / 2
}

// 節氣名稱列表 (繁體中文)
var SolarTerms = []string{
	"春分", "清明", "穀雨", "立夏", "小滿", "芒種",
	"夏至", "小暑", "大暑", "立秋", "處暑", "白露",
	"秋分", "寒露", "霜降", "立冬", "小雪", "大雪",
	"冬至", "小寒", "大寒", "立春", "雨水", "驚蟄",
}

// FindPreviousWinterSolstice 尋找上一個冬至 (270度) 的 JDE
func FindPreviousWinterSolstice(jde float64) float64 {
	// 粗略定位 270 度的範圍
	for d := jde; d > jde-400; d -= 5 {
		lon := SolarLongitude(d)
		if lon > 265 && lon < 275 {
			return EstimateTermTime(270.0, d-10, d+10)
		}
	}
	return EstimateTermTime(270.0, jde-370.0, jde)
}
