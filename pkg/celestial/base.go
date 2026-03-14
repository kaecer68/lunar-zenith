package celestial

import (
	"math"
	"time"
)

// PrecisionTime 封裝了天文計算所需的高精度時間結構
type PrecisionTime struct {
	UT  time.Time // 民用協調世界時 (Universal Time)
	JD  float64   // 儒略日 (Julian Day in UT)
	JDE float64   // 儒略曆元 (Julian Ephemeris Day in TT)
	DeltaT float64 // TT - UT 的差值（秒）
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

// EstimateDeltaT 估算 TT 與 UT 之間的差值
// 這裡暫時使用一個簡化的二次多項式擬合算法，未來將優化為表驅動
func EstimateDeltaT(t time.Time) float64 {
	y := float64(t.Year()) + (float64(t.Month())-0.5)/12.0
	t_val := (y - 2000) / 100
	// 根據 NASA/Espenak 擬合公式 (2005-2050)
	// ΔT = 62.92 + 0.32217 * (y - 2000) + 0.005589 * (y - 2000)^2
	if y >= 2000 && y <= 2100 {
		dy := y - 2000
		return 62.92 + 0.32217*dy + 0.005589*dy*dy
	}
	// 預設返回當前大致值 (2024年約為 69s)
	_ = t_val
	return 69.0
}
