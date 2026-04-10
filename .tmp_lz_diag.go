package main

import (
	"fmt"
	"math"
	"time"

	"github.com/kaecer68/lunar-zenith/pkg/celestial"
)

func monthHasZhongqi(nmStart, nmEnd float64) bool {
	lonStart := celestial.SolarLongitude(nmStart)
	lonEnd := celestial.SolarLongitude(nmEnd)
	gridStart := math.Floor(lonStart/30.0) * 30.0
	gridEnd := math.Floor(lonEnd/30.0) * 30.0
	if lonEnd < lonStart {
		return true
	}
	return gridEnd > gridStart
}

func buildLeapIndex(nmWS, ws float64) (int, []float64) {
	nextWS := celestial.FindNextWinterSolstice(ws + 20)
	months := make([]float64, 0, 15)
	curr := nmWS
	for curr < nextWS+1 {
		months = append(months, curr)
		curr = celestial.FindNewMoon(curr+15, 1)
	}
	if len(months) <= 12 {
		return -1, months
	}
	for i := 1; i < len(months)-1; i++ {
		if !monthHasZhongqi(months[i], months[i+1]) {
			return i, months
		}
	}
	return -1, months
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func diagnose(date string) {
	t, _ := time.Parse("2006-01-02", date)
	loc, _ := time.LoadLocation("Asia/Taipei")
	tt := time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, loc)
	jd := celestial.TimeToJD(tt.UTC())

	year, month, _ := celestial.JDToDate(jd)
	deltaT := celestial.EstimateDeltaT(time.Date(year, time.Month(month), 15, 0, 0, 0, 0, time.UTC))
	jde := jd + deltaT/86400.0

	nm0 := celestial.PreviousNewMoon(jde)
	ws := celestial.FindPreviousWinterSolstice(jde + 32)
	nmWS := celestial.PreviousNewMoon(ws)
	if nmWS > nm0+0.01 {
		ws = celestial.FindPreviousWinterSolstice(nmWS - 2)
		nmWS = celestial.FindNewMoon(ws, -1)
	}

	leapPos, allMonths := buildLeapIndex(nmWS, ws)
	pos := 0
	for i, nm := range allMonths {
		if math.Abs(nm-nm0) < 0.5 {
			pos = i
			break
		}
	}

	fmt.Printf("%s jd=%.6f jde=%.6f nm0=%.6f pos=%d leapPos=%d months=%d\n", date, jd, jde, nm0, pos, leapPos, len(allMonths))
	for i := max(0, pos-2); i <= min(len(allMonths)-1, pos+3); i++ {
		mark := ""
		if i == pos {
			mark += " <-current"
		}
		if i == leapPos {
			mark += " <-leap"
		}
		fmt.Printf("  [%02d] %.6f lon=%.2f%s\n", i, allMonths[i], celestial.SolarLongitude(allMonths[i]), mark)
	}
}

func main() {
	for _, d := range []string{"2025-07-25", "2020-05-23", "2001-05-23", "2014-10-24", "1984-11-23", "2033-12-22"} {
		diagnose(d)
	}
}
