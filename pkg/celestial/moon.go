package celestial

import (
	"math"
	"os"
	"path/filepath"
	"sync"

	swisseph "github.com/tejzpr/go-swisseph"
)

const synodicMonth = 29.530588853

var (
	moonEpheInitOnce sync.Once
	moonEpheReady    bool
)

// MoonLongitude 計算給定 JDE 的月球黃經 (精簡版 ELP2000-82)
func MoonLongitude(jde float64) float64 {
	if lon, ok := moonLongitudeSwiss(jde); ok {
		return lon
	}
	return moonLongitudeApprox(jde)
}

func moonLongitudeApprox(jde float64) float64 {
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

func moonLongitudeSwiss(jde float64) (float64, bool) {
	moonEpheInitOnce.Do(func() {
		moonEpheReady = initMoonEphemerisPath()
	})
	if !moonEpheReady {
		return 0, false
	}

	res := swisseph.CalcUT(jde, swisseph.Moon, swisseph.FlagSwieph)
	if res.Flag < 0 {
		return 0, false
	}

	lon := math.Mod(res.Data[0], 360.0)
	if lon < 0 {
		lon += 360.0
	}
	return lon, true
}

func initMoonEphemerisPath() bool {
	paths := []string{
		"./data/ephe",
		"../data/ephe",
		"../../data/ephe",
	}

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		paths = append(paths,
			filepath.Join(exeDir, "data/ephe"),
			filepath.Join(exeDir, "../data/ephe"),
		)
	}

	for _, path := range paths {
		if isMoonEphemerisPath(path) {
			swisseph.SetEphePath(path)
			return true
		}
	}

	swisseph.SetEphePath("")
	return true
}

func isMoonEphemerisPath(path string) bool {
	if _, err := os.Stat(filepath.Join(path, "semo_18.se1")); err == nil {
		return true
	}
	if _, err := os.Stat(filepath.Join(path, "semo_12.se1")); err == nil {
		return true
	}
	return false
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
	if direction < 0 {
		return PreviousNewMoon(jde)
	}

	const epsilon = 1e-9
	k := math.Floor((jde - 2451550.09766) / synodicMonth)
	candidate := newMoonByK(k)
	if candidate <= jde+epsilon {
		k++
	}
	return newMoonByK(k)
}

// meanNewMoonJDE 根據 Meeus 第 49 章公式計算第 k 個朔日的 JDE（k=0 對應 2000-01-06.0 UT）。
func meanNewMoonJDE(k float64) float64 {
	t := k / 1236.85
	return 2451550.09766 +
		29.530588853*k +
		0.00015437*t*t -
		0.000000150*t*t*t +
		0.00000000073*t*t*t*t
}

// trueNewMoonJDE 計算第 k 個朔日的真朔 JDE。
// 公式來源: Jean Meeus, Astronomical Algorithms (2nd ed.), Chapter 49.
func trueNewMoonJDE(k float64) float64 {
	t := k / 1236.85
	t2 := t * t
	t3 := t2 * t
	t4 := t3 * t

	e := 1 - 0.002516*t - 0.0000074*t2

	m := normalizeDeg(2.5534 + 29.10535670*k - 0.0000014*t2 - 0.00000011*t3)
	mPrime := normalizeDeg(201.5643 + 385.81693528*k + 0.0107582*t2 + 0.00001238*t3 - 0.000000058*t4)
	f := normalizeDeg(160.7108 + 390.67050284*k - 0.0016118*t2 - 0.00000227*t3 + 0.000000011*t4)
	omega := normalizeDeg(124.7746 - 1.56375588*k + 0.0020672*t2 + 0.00000215*t3)

	jde := meanNewMoonJDE(k)

	// 主週期修正項（新月）
	jde += -0.40720*math.Sin(mPrime*Deg2Rad) +
		0.17241*e*math.Sin(m*Deg2Rad) +
		0.01608*math.Sin(2*mPrime*Deg2Rad) +
		0.01039*math.Sin(2*f*Deg2Rad) +
		0.00739*e*math.Sin((mPrime-m)*Deg2Rad) -
		0.00514*e*math.Sin((mPrime+m)*Deg2Rad) +
		0.00208*e*e*math.Sin(2*m*Deg2Rad) -
		0.00111*math.Sin((mPrime-2*f)*Deg2Rad) -
		0.00057*math.Sin((mPrime+2*f)*Deg2Rad) +
		0.00056*e*math.Sin((2*mPrime+m)*Deg2Rad) -
		0.00042*math.Sin(3*mPrime*Deg2Rad) +
		0.00042*e*math.Sin((m+2*f)*Deg2Rad) +
		0.00038*e*math.Sin((m-2*f)*Deg2Rad) -
		0.00024*e*math.Sin((2*mPrime-m)*Deg2Rad) -
		0.00017*math.Sin(omega*Deg2Rad) -
		0.00007*math.Sin((mPrime+2*m)*Deg2Rad) +
		0.00004*math.Sin((2*mPrime-2*f)*Deg2Rad) +
		0.00004*math.Sin(3*m*Deg2Rad) +
		0.00003*math.Sin((mPrime+m-2*f)*Deg2Rad) +
		0.00003*math.Sin((2*mPrime+2*f)*Deg2Rad) -
		0.00003*math.Sin((mPrime+m+2*f)*Deg2Rad) +
		0.00003*math.Sin((mPrime-m+2*f)*Deg2Rad) -
		0.00002*math.Sin((mPrime-m-2*f)*Deg2Rad) -
		0.00002*math.Sin((3*mPrime+m)*Deg2Rad) +
		0.00002*math.Sin(4*mPrime*Deg2Rad)

	// 行星攝動修正項
	a1 := normalizeDeg(299.77 + 0.107408*k - 0.009173*t2)
	a2 := normalizeDeg(251.88 + 0.016321*k)
	a3 := normalizeDeg(251.83 + 26.651886*k)
	a4 := normalizeDeg(349.42 + 36.412478*k)
	a5 := normalizeDeg(84.66 + 18.206239*k)
	a6 := normalizeDeg(141.74 + 53.303771*k)
	a7 := normalizeDeg(207.14 + 2.453732*k)
	a8 := normalizeDeg(154.84 + 7.306860*k)
	a9 := normalizeDeg(34.52 + 27.261239*k)
	a10 := normalizeDeg(207.19 + 0.121824*k)
	a11 := normalizeDeg(291.34 + 1.844379*k)
	a12 := normalizeDeg(161.72 + 24.198154*k)
	a13 := normalizeDeg(239.56 + 25.513099*k)
	a14 := normalizeDeg(331.55 + 3.592518*k)

	jde += 0.000325*math.Sin(a1*Deg2Rad) +
		0.000165*math.Sin(a2*Deg2Rad) +
		0.000164*math.Sin(a3*Deg2Rad) +
		0.000126*math.Sin(a4*Deg2Rad) +
		0.000110*math.Sin(a5*Deg2Rad) +
		0.000062*math.Sin(a6*Deg2Rad) +
		0.000060*math.Sin(a7*Deg2Rad) +
		0.000056*math.Sin(a8*Deg2Rad) +
		0.000047*math.Sin(a9*Deg2Rad) +
		0.000042*math.Sin(a10*Deg2Rad) +
		0.000040*math.Sin(a11*Deg2Rad) +
		0.000037*math.Sin(a12*Deg2Rad) +
		0.000035*math.Sin(a13*Deg2Rad) +
		0.000023*math.Sin(a14*Deg2Rad)

	return jde
}

func normalizeDeg(v float64) float64 {
	v = math.Mod(v, 360.0)
	if v < 0 {
		v += 360.0
	}
	return v
}

// PreviousNewMoon 以 Meeus 均值公式估算再用 ±2 天窄窗口二分求精，
// 穩定回傳 jde 之前最近的朔日 JDE。
// 解決 FindNewMoon(jde,-1) 在「jde 剛好在朔日後幾小時」時收斂到
// 前一個月朔日的二分法臨界問題。
func PreviousNewMoon(jde float64) float64 {
	const epsilon = 1e-9
	k := math.Floor((jde - 2451550.09766) / synodicMonth)
	candidate := newMoonByK(k)
	if candidate > jde+epsilon {
		k--
	}
	for {
		next := newMoonByK(k + 1)
		if next > jde+epsilon {
			break
		}
		k++
	}
	return newMoonByK(k)
}

func newMoonByK(k float64) float64 {
	guess := trueNewMoonJDE(k)
	return refineNewMoonAround(guess)
}

func refineNewMoonAround(guess float64) float64 {
	phaseSigned := func(t float64) float64 {
		p := MoonPhase(t)
		if p > 180 {
			p -= 360
		}
		return p
	}

	low := guess - 1.5
	high := guess + 1.5
	const scanStep = 0.125

	prevT := low
	prevS := phaseSigned(prevT)
	found := false
	for t := low + scanStep; t <= high+1e-9; t += scanStep {
		s := phaseSigned(t)
		if prevS <= 0 && s >= 0 {
			low = prevT
			high = t
			found = true
			break
		}
		prevT = t
		prevS = s
	}

	if !found {
		return guess
	}

	const precision = 0.00001
	for i := 0; i < 60; i++ {
		mid := (low + high) / 2.0
		s := phaseSigned(mid)
		if math.Abs(high-low) < precision {
			return mid
		}
		if s < 0 {
			low = mid
		} else {
			high = mid
		}
	}

	return (low + high) / 2.0
}
