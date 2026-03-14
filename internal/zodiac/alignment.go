package zodiac

import (
	"math"
	"github.com/kaecer68/lunar-zenith/internal/celestial"
)

// GetSolarMonth 根據太陽黃經確定節氣月 (1:寅, 2:卯 ... 11:子, 12:丑)
// 在命理學中，月份是以「節」開始計算的
func GetSolarMonth(lon float64) int {
	// 將黃經偏移，使得立春 (315度) 對應到度數 0
	shifted := math.Mod(lon+45.0, 360.0)
	if shifted < 0 {
		shifted += 360.0
	}
	// 每 30 度為一個節氣月
	monthIndex := int(math.Floor(shifted/30.0)) + 1
	return monthIndex
}

// AstrologicalPillar 封裝了特定時刻的年月日時四柱
type AstrologicalPillar struct {
	Year  Sexagenary
	Month Sexagenary
	Day   Sexagenary
	Hour  Sexagenary
}

// GetAstrologicalPillar 獲取指定 JDE 時刻的命理四柱 (基於節氣月)
func GetAstrologicalPillar(pt *celestial.PrecisionTime) AstrologicalPillar {
	lon := celestial.SolarLongitude(pt.JDE)
	
	// 1. 年柱 (注意：立春前算上一年)
	year := pt.UT.Year()
	// 如果在立春前，年柱算上一年
	// 立春點為 315 度
	if lon < 315.0 && lon > 45.0 {
		// 這裡有個邏輯跨越：1月1日到立春點之間，lon 是在 [285, 315) 左右。
		// 所以如果 lon > 280 (冬至後) 且 < 315，表示還沒到立春。
	}
	
	// 更精準判斷：
	// 若 lon < 315 && 月份是 1月或 2月，則視為上一年
	calcYear := year
	if (pt.UT.Month() == 1 || pt.UT.Month() == 2) && lon < 315.0 {
		calcYear--
	}
	yearPillar := NewYearSexagenary(calcYear)
	
	// 2. 月柱 (基于節氣月)
	monthIdx := GetSolarMonth(lon)
	monthPillar := GetMonthSexagenary(yearPillar.StemIndex, monthIdx)
	
	// 3. 日柱
	dayPillar := GetDaySexagenary(pt.JD)
	
	// 4. 時柱 (考慮早子、晚子)
	hour := pt.UT.Hour()
	// 如果是 23:00 後，算在隔天的早子時 (命理習慣)
	// 但 GetDaySexagenary 已經用了 JD，JD 的 23:00 (+0.5) 會自動推移
	hourBranch := GetHourBranch(hour)
	hourPillar := GetHourSexagenary(dayPillar.StemIndex, hourBranch)
	
	return AstrologicalPillar{
		Year:  yearPillar,
		Month: monthPillar,
		Day:   dayPillar,
		Hour:  hourPillar,
	}
}
