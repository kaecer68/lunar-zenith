package celestial

import (
	"testing"
	"time"
)

func TestTimeToJD(t *testing.T) {
	// 基準測試：2000 年 1 月 1 日 12:00:00 UTC 應該是 JD 2451545.0
	t1 := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	got := TimeToJD(t1)
	want := 2451545.0
	if got != want {
		t.Errorf("TimeToJD(2000-01-01 12:00:00) = %f; want %f", got, want)
	}

	// 測試 1970 紀元
	t2 := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	got2 := TimeToJD(t2)
	want2 := 2440587.5
	if got2 != want2 {
		t.Errorf("TimeToJD(1970-01-01 00:00:00) = %f; want %f", got2, want2)
	}
}

func TestNewPrecisionTime(t *testing.T) {
	now := time.Now()
	pt := NewPrecisionTime(now)
	if pt.JD == 0 {
		t.Error("JD calculation resulted in 0")
	}
	if pt.JDE <= pt.JD {
		t.Errorf("Expected JDE (%f) to be greater than JD (%f) due to DeltaT", pt.JDE, pt.JD)
	}
}
