package celestial

import (
	"math"
	"testing"
	"time"
)

func TestSolarLongitude(t *testing.T) {
	// 2024 年春分時刻約在 3/20 03:06 UTC，太陽黃經應接近 0 度
	t1 := time.Date(2024, 3, 20, 3, 6, 0, 0, time.UTC)
	pt := NewPrecisionTime(t1)
	lon := SolarLongitude(pt.JDE)
	
	// 容許誤差在 0.05 度以內（精簡版 VSOP87）
	if math.Abs(lon-0.0) > 0.05 && math.Abs(lon-360.0) > 0.05 {
		t.Errorf("2024 春分太陽黃經 = %f; want 接近 0.0", lon)
	}
}

func TestGetSolarTerm(t *testing.T) {
	// 測試 2024 立春 (2/4 接近 315 度)
	t1 := time.Date(2024, 2, 4, 8, 27, 0, 0, time.UTC)
	pt := NewPrecisionTime(t1)
	info := GetSolarTerm(pt.JDE)
	
	// 依照我們 SolarTerms 數組定義：春分=0, ... 立春=21, 雨水=22, 驚蟄=23
	// 立春是 315 度，315/15 = 21
	if info.Index != 21 {
		t.Errorf("2024 立春判定錯誤 = %d (%s); want 21 (立春)", info.Index, info.Name)
	}
}
