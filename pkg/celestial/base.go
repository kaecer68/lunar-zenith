package celestial

import (
	"math"
	"time"
)

var taipeiLocation = func() *time.Location {
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return time.FixedZone("CST", 8*3600)
	}
	return loc
}()

// PrecisionTime 封裝了天文計算所需的高精度時間結構
type PrecisionTime struct {
	UT     time.Time // 民用協調世界時 (Universal Time)
	JD     float64   // 儒略日 (Julian Day in UT)
	JDE    float64   // 儒略曆元 (Julian Ephemeris Day in TT)
	DeltaT float64   // TT - UT 的差值（秒）
}

// NewPrecisionTime 根據給定的 Go Time 創建高精度時間對象
func NewPrecisionTime(t time.Time) *PrecisionTime {
	pt := &PrecisionTime{
		UT: t.UTC(),
	}
	pt.JD = TimeToJD(pt.UT)
	pt.DeltaT = EstimateDeltaT(pt.UT)
	pt.JDE = pt.JD + (pt.DeltaT / 86400.0)
	return pt
}

// TimeToJD 將 Go 的時間對象轉換為儒略日 (Julian Day)
// 算法參考 Jean Meeus "Astronomical Algorithms" 第二章
func TimeToJD(t time.Time) float64 {
	y := float64(t.Year())
	m := float64(t.Month())
	d := float64(t.Day()) + float64(t.Hour())/24.0 + float64(t.Minute())/1440.0 + float64(t.Second())/86400.0

	if m <= 2 {
		y--
		m += 12
	}

	a := math.Floor(y / 100)
	b := 2 - a + math.Floor(a/4)

	jd := math.Floor(365.25*(y+4716)) + math.Floor(30.6001*(m+1)) + d + b - 1524.5
	return jd
}

// JDToDate 將儒略日 (Julian Day) 轉換為公曆日期
// 算法參考 Jean Meeus "Astronomical Algorithms" 第七章
func JDToDate(jd float64) (year, month, day int) {
	jd += 0.5
	z := math.Floor(jd)
	f := jd - z
	var a float64
	if z < 2299161 {
		a = z
	} else {
		alpha := math.Floor((z - 1867216.25) / 36524.25)
		a = z + 1 + alpha - math.Floor(alpha/4)
	}
	b := a + 1524
	c := math.Floor((b - 122.1) / 365.25)
	d := math.Floor(365.25 * c)
	e := math.Floor((b - d) / 30.6001)

	dayFrac := b - d - math.Floor(30.6001*e) + f
	day = int(dayFrac)

	if e < 14 {
		month = int(e - 1)
	} else {
		month = int(e - 13)
	}

	if month > 2 {
		year = int(c - 4716)
	} else {
		year = int(c - 4715)
	}
	return
}

// JDToTime 將儒略日轉換為指定時區的 time.Time
// jd: 儒略日 (UT)
// loc: 目標時區（如 Asia/Taipei）
func JDToTime(jd float64, loc *time.Location) time.Time {
	y, m, d := JDToDate(jd)
dayFrac := jd + 0.5 - math.Floor(jd+0.5)
	hour := int(dayFrac * 24)
	min := int((dayFrac*24 - float64(hour)) * 60)
	sec := int(((dayFrac*24-float64(hour))*60 - float64(min)) * 60)

	// dayFrac ∈ [0,1) ⇒ hour ∈ [0,23]，無需跨日調整
	return time.Date(y, time.Month(m), d, hour, min, sec, 0, loc)
}

// EstimateDeltaT 估算 TT 與 UT 之間的差值（秒）
// 採用 NASA/Espenak & Meeus 多段多項式擬合，覆蓋 1800–2150 年
// 參考: https://eclipsewise.com/help/deltatpoly2014.html
func EstimateDeltaT(t time.Time) float64 {
	y := float64(t.Year()) + (float64(t.Month())-0.5)/12.0

	switch {
	case y < 1800:
		// 延伸公式：Stephenson & Houlden (1986) 適用 948–1600，此處做保守外推
		u := (y - 1820) / 100
		return -20 + 32*u*u

	case y < 1860:
		u := y - 1820
		return 124.07 - 0.9246*u + 0.003141*u*u

	case y < 1900:
		u := y - 1860
		return 7.62 + 0.5737*u - 0.251754*u*u +
			0.01680668*u*u*u -
			0.0004473624*u*u*u*u +
			u*u*u*u*u/233174

	case y < 1920:
		u := y - 1900
		return -2.79 + 1.494119*u - 0.0598939*u*u +
			0.0061966*u*u*u - 0.000197*u*u*u*u

	case y < 1941:
		u := y - 1920
		return 21.20 + 0.84493*u - 0.076100*u*u + 0.0020936*u*u*u

	case y < 1961:
		u := y - 1950
		return 29.07 + 0.407*u - u*u/233 + u*u*u/2547

	case y < 1986:
		u := y - 1975
		return 45.45 + 1.067*u - u*u/260 - u*u*u/718

	case y < 2005:
		u := y - 2000
		return 63.86 + 0.3345*u - 0.060374*u*u +
			0.0017275*u*u*u + 0.000651814*u*u*u*u +
			0.00002373599*u*u*u*u*u

	case y < 2050:
		u := y - 2000
		return 62.92 + 0.32217*u + 0.005589*u*u

	case y < 2150:
		return -20 + 32*((y-1820)/100)*((y-1820)/100) - 0.5628*(2150-y)

	default:
		u := (y - 1820) / 100
		return -20 + 32*u*u
	}
}
