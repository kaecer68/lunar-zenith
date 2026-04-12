//go:build ignore
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kaecer68/lunar-zenith/v4/pkg/celestial"
	"github.com/kaecer68/lunar-zenith/v4/pkg/zodiac"
)

// 1900-2100 年閏月資料（來源：香港天文台、紫金山天文台）
var verifiedLeapMonths = []struct {
	year      int
	leapMonth int // 閏月（1-12）
}{
	{1900, 8}, // 閏八月
	{1903, 5}, // 閏五月
	{1906, 4}, // 閏四月
	{1909, 2}, // 閏二月
	{1911, 6}, // 閏六月
	{1914, 5}, // 閏五月
	{1917, 2}, // 閏二月
	{1919, 7}, // 閏七月
	{1922, 5}, // 閏五月
	{1925, 4}, // 閏四月
	{1928, 2}, // 閏二月
	{1930, 6}, // 閏六月
	{1933, 5}, // 閏五月
	{1936, 3}, // 閏三月
	{1938, 7}, // 閏七月
	{1941, 6}, // 閏六月
	{1944, 4}, // 閏四月
	{1947, 2}, // 閏二月
	{1949, 7}, // 閏七月
	{1952, 5}, // 閏五月
	{1955, 3}, // 閏三月
	{1957, 8}, // 閏八月
	{1960, 6}, // 閏六月
	{1963, 4}, // 閏四月
	{1966, 3}, // 閏三月
	{1968, 7}, // 閏七月
	{1971, 5}, // 閏五月
	{1974, 4}, // 閏四月
	{1976, 8}, // 閏八月
	{1979, 6}, // 閏六月
	{1982, 4}, // 閏四月
	{1984, 10}, // 閏十月
	{1987, 6}, // 閏六月
	{1990, 5}, // 閏五月
	{1993, 3}, // 閏三月
	{1995, 8}, // 閏八月
	{1998, 5}, // 閏五月
	{2001, 4}, // 閏四月
	{2004, 2}, // 閏二月
	{2006, 7}, // 閏七月
	{2009, 5}, // 閏五月
	{2012, 4}, // 閏四月
	{2014, 9}, // 閏九月
	{2017, 6}, // 閏六月
	{2020, 4}, // 閏四月
	{2023, 2}, // 閏二月
	{2025, 6}, // 閏六月
	{2028, 5}, // 閏五月
	{2031, 3}, // 閏三月
	{2033, 7}, // 閏七月（經典邊界案例）
	{2036, 6}, // 閏六月
	{2039, 5}, // 閏五月
	{2042, 2}, // 閏二月
	{2044, 7}, // 閏七月
	{2047, 5}, // 閏五月
	{2050, 3}, // 閏三月
	{2052, 8}, // 閏八月
	{2055, 6}, // 閏六月
	{2058, 5}, // 閏五月
	{2061, 3}, // 閏三月
	{2063, 7}, // 閏七月
	{2066, 5}, // 閏五月
	{2069, 4}, // 閏四月
	{2071, 8}, // 閏八月
	{2074, 6}, // 閏六月
	{2077, 4}, // 閏四月
	{2080, 3}, // 閏三月
	{2082, 7}, // 閏七月
	{2085, 5}, // 閏五月
	{2088, 4}, // 閏四月
	{2090, 8}, // 閏八月
	{2093, 6}, // 閏六月
	{2096, 4}, // 閏四月
	{2099, 2}, // 閏二月
}

func main() {
	engine := &zodiac.LunarEngine{}

	fmt.Println("=== 1900-2100 年閏月驗證 ===")
	fmt.Println()

	passed := 0
	failed := 0
	failedYears := []int{}

	for _, vm := range verifiedLeapMonths {
		// 找到該年閏月的初一
		leapMonthFound := verifyLeapMonth(engine, vm.year, vm.leapMonth)

		if leapMonthFound {
			fmt.Printf("✓ %d年 閏%d月: 正確\n", vm.year, vm.leapMonth)
			passed++
		} else {
			fmt.Printf("✗ %d年 閏%d月: 錯誤\n", vm.year, vm.leapMonth)
			failed++
			failedYears = append(failedYears, vm.year)
		}
	}

	fmt.Println()
	fmt.Println("=== 結果統計 ===")
	fmt.Printf("總共: %d 個閏月年份\n", len(verifiedLeapMonths))
	fmt.Printf("通過: %d\n", passed)
	fmt.Printf("失敗: %d\n", failed)

	if len(failedYears) > 0 {
		fmt.Printf("\n失敗年份: %v\n", failedYears)
		os.Exit(1)
	} else {
		fmt.Println("\n✓ 所有閏月驗證通過！")
	}
}

func verifyLeapMonth(engine *zodiac.LunarEngine, year, leapMonth int) bool {
	// 估算閏月時間（農曆月份約對應西曆月份+0-2個月）
	// 農曆月份通常比西曆晚 0-2 個月，閏月也是如此
	// 例如：閏二月通常在 3-4 月
	estimatedStartMonth := leapMonth
	if estimatedStartMonth < 3 {
		estimatedStartMonth = 3 // 至少從3月開始
	}

	startDate := time.Date(year, time.Month(estimatedStartMonth-2), 1, 12, 0, 0, 0, time.UTC)
	endDate := time.Date(year, time.Month(estimatedStartMonth+3), 1, 12, 0, 0, 0, time.UTC)

	foundAny := false
	foundDay1 := false
	maxDay := 0

	for d := startDate; d.Before(endDate); d = d.Add(24 * time.Hour) {
		pt := celestial.NewPrecisionTime(d)
		lunar := engine.GetLunarDate(pt.JD)

		if lunar.Month == leapMonth {
			foundAny = true
			if lunar.Day == 1 {
				foundDay1 = true
				if lunar.IsLeap {
					return true // 找到閏月初一
				}
			}
			if lunar.Day > maxDay {
				maxDay = lunar.Day
			}
		}
	}

	// 調試輸出
	if !foundAny {
		fmt.Printf("  [除錯] %d年: 未找到 %d月\n", year, leapMonth)
	} else if !foundDay1 {
		fmt.Printf("  [除錯] %d年: 找到 %d月但無初一，最大日期=%d\n", year, leapMonth, maxDay)
	} else {
		fmt.Printf("  [除錯] %d年: 找到 %d月初一但非閏月\n", year, leapMonth)
	}

	return false
}
