package celestial

import (
	"math"
)

// MoonLongitude 計算給定 JDE 的月球黃經 (精簡版 ELP2000-82)
func MoonLongitude(jde float64) float64 {
	t := (jde - 2451545.0) / 36525.0 // 儒略世紀數

	// 月球平黃經 (Mean Longitude)
	lPrime := 218.3164477 + 481267.88123421*t
	// 月球平近點角 (Mean Anomaly)
	mPrime := 134.9633964 + 477198.8675055*t
	// 太陽平近點角 (Mean Anomaly)
	m := 357.5291092 + 35999.0502909*t
	// 月球平緯角 (Mean Argument of Latitude)
	f := 93.2720950 + 483202.0175233*t
	// 日月黃經差 (Mean Elongation)
	d := 297.8501921 + 445267.1114034*t

	// 修正到 0-360
	lPrime = math.Mod(lPrime, 360.0)
	if lPrime < 0 {
		lPrime += 360.0
	}

	// ELP2000 週期項 (簡化版本，精度足以判定初一)
	lambda := lPrime +
		6.288774*math.Sin(mPrime*Deg2Rad) +
		1.274027*math.Sin((2*d-mPrime)*Deg2Rad) +
		0.658314*math.Sin(2*d*Deg2Rad) +
		0.213118*math.Sin(2*mPrime*Deg2Rad) -
		0.185116*math.Sin(m*Deg2Rad) -
		0.114332*math.Sin(2*f*Deg2Rad)

	return math.Mod(lambda+360.0, 360.0)
}

// MoonPhase 計算日月黃經差 (Elongation)
// 0 度表示「朔」(New Moon)
func MoonPhase(jde float64) float64 {
	sLon := SolarLongitude(jde)
	mLon := MoonLongitude(jde)
	diff := mLon - sLon
	return math.Mod(diff+360.0, 360.0)
}

// FindNewMoon 搜尋距離指定 jde 最近的前一個或後一個「朔」(New Moon)
// direction: -1 (搜尋前一個), 1 (搜尋下一個)
func FindNewMoon(jde float64, direction float64) float64 {
	low := jde
	high := jde + direction*30.0 // 朔望月約 29.53 天

	if direction < 0 {
		low, high = high, low
	}

	const precision = 0.00001 // 約 1 秒

	for i := 0; i < 40; i++ {
		mid := (low + high) / 2
		phase := MoonPhase(mid)

		// 我們尋找相位相接近 0 的點
		// 注意：跨越 360/0 的問題
		if phase > 180 {
			phase -= 360
		}

		if math.Abs(high-low) < precision {
			return mid
		}

		// 根據相位方向調整邊界 (這裡需要考慮 direction)
		if phase < 0 {
			low = mid
		} else {
			high = mid
		}
	}
	return (low + high) / 2
}
