//go:build ignore
package main

import (
	"fmt"
	"math"
	"time"

	"github.com/kaecer68/lunar-zenith/v4/pkg/celestial"
	"github.com/kaecer68/lunar-zenith/v4/pkg/zodiac"
)

func main() {
	fmt.Println("=== 2033 年閏月診斷 ===")
	fmt.Println()

	engine := &zodiac.LunarEngine{}

	// 測試 2033 年的關鍵日期
	testDates := []string{
		"2033-08-20", // 七月廿六
		"2033-08-21", //
		"2033-08-22", //
		"2033-08-23", //
		"2033-08-24", //
		"2033-08-25", // 應為閏七月初一
		"2033-08-26", //
		"2033-08-27", //
		"2033-08-28", //
		"2033-08-29", //
		"2033-08-30", // 八月初一（實際返回）
	}

	fmt.Println("日期		農曆結果")
	fmt.Println("-------------------------------")
	for _, dateStr := range testDates {
		t, _ := time.Parse("2006-01-02", dateStr)
		pt := celestial.NewPrecisionTime(t)
		lunar := engine.GetLunarDate(pt.JD)

		leapStr := ""
		if lunar.IsLeap {
			leapStr = "閏"
		}
		fmt.Printf("%s\t%d年%s%d月%d日\n",
			dateStr, lunar.Year, leapStr, lunar.Month, lunar.Day)
	}

	// 計算冬至和朔日
	fmt.Println()
	fmt.Println("=== 2032-2033 冬至年分析 ===")

	// 2032 年冬至
	ws2032 := findWinterSolstice(2032)
	fmt.Printf("2032 年冬至: JD %.6f (%s)\n", ws2032, jdToDate(ws2032))

	// 2033 年冬至
	ws2033 := findWinterSolstice(2033)
	fmt.Printf("2033 年冬至: JD %.6f (%s)\n", ws2033, jdToDate(ws2033))

	// 找到含冬至的朔日
	nmWS := findWinterSolsticeMonthStart(ws2032)
	fmt.Printf("含 2032 冬至的朔日: JD %.6f (%s)\n", nmWS, jdToDate(nmWS))

	// 收集所有朔日
	fmt.Println()
	fmt.Println("朔日列表 (從 2032 冬至開始):")
	months := collectMonths(nmWS, ws2033)
	for i, m := range months {
		fmt.Printf("  [%2d] JD %.6f (%s)\n", i, m, jdToDate(m))
	}

	fmt.Printf("\n總共 %d 個朔日\n", len(months))

	// 分析每個月的朔望月和相關節氣
	if len(months) > 1 {
		fmt.Println()
		fmt.Println("=== 每月朔望月分析 ===")
		fmt.Println("月索引\t朔日\t\t朔日黃經\t次朔黃經\t中氣?\t中氣日期")
		fmt.Println("----------------------------------------------------------------")

		for i := 0; i < len(months)-1 && i < 15; i++ {
			nmStart := months[i]
			nmEnd := months[i+1]

			sLonStart := celestial.SolarLongitude(nmStart)
			sLonEnd := celestial.SolarLongitude(nmEnd)

			hasZQ, _ := monthHasZhongqi(nmStart, nmEnd)
			zqStatus := "無"
			if hasZQ {
				zqStatus = "有"
			}

			// 計算這個月對應的農曆月份（從11月開始）
			lunarMonth := (i) % 12
			if lunarMonth == 0 {
				lunarMonth = 12
			}
			if i == 0 {
				lunarMonth = 11 // 第一個月是11月
			}

			fmt.Printf("[%2d]\t%s\t%.2f°\t%.2f°\t%s\t月序=%d\n",
				i, jdToDate(nmStart)[:10], sLonStart, sLonEnd, zqStatus, lunarMonth)
		}
	}
}

func findWinterSolstice(year int) float64 {
	// 估算冬至時間（約 12 月 21-22 日）
	t := time.Date(year, 12, 20, 0, 0, 0, 0, time.UTC)
	startJDE := float64(t.Unix())/86400.0 + 2440587.5

	// 尋找黃經 270°
	for d := startJDE; d < startJDE+10; d += 0.1 {
		lon := celestial.SolarLongitude(d)
		if lon >= 269.5 && lon <= 270.5 {
			return celestial.EstimateTermTime(270.0, d-1, d+1)
		}
	}
	return startJDE + 1
}

func findWinterSolsticeMonthStart(ws float64) float64 {
	nmPrev := celestial.PreviousNewMoon(ws)
	nmNext := celestial.FindNewMoon(ws, 1)

	// 若冬至與下一朔落在同一天，取下一朔
	if civilDayUTC8(nmNext) == civilDayUTC8(ws) {
		return nmNext
	}
	return nmPrev
}

func civilDayUTC8(jd float64) int {
	return int(math.Floor(jd + 0.5 + 8.0/24.0))
}

func collectMonths(nmWS, nextWS float64) []float64 {
	months := make([]float64, 0, 15)
	curr := nmWS
	for curr < nextWS+1 {
		months = append(months, curr)
		curr = celestial.FindNewMoon(curr+15, 1)
	}
	return months
}

func jdToDate(jd float64) string {
	unixTime := int64((jd - 2440587.5) * 86400)
	t := time.Unix(unixTime, 0).UTC()
	return t.Format("2006-01-02 15:04")
}

func monthHasZhongqi(nmStart, nmEnd float64) (has bool, onEndDay bool) {
	sLonStart := celestial.SolarLongitude(nmStart)
	sLonEnd := celestial.SolarLongitude(nmEnd)
	endUnwrapped := sLonEnd
	if endUnwrapped < sLonStart {
		endUnwrapped += 360.0
	}

	boundary := math.Ceil((sLonStart+1e-12)/30.0) * 30.0
	if boundary > endUnwrapped+1e-12 {
		return false, false
	}

	target := math.Mod(boundary, 360.0)
	termJDE := celestial.EstimateTermTime(target, nmStart, nmEnd)

	startDay := civilDayUTC8(nmStart)
	endDay := civilDayUTC8(nmEnd)
	termDay := civilDayUTC8(termJDE)

	if termDay == endDay {
		return false, true
	}

	return termDay >= startDay && termDay < endDay, false
}
