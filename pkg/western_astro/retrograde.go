package western_astro

import (
	"fmt"
	"time"

	"github.com/kaecer68/lunar-zenith/pkg/celestial"
)

type RetrogradeInfo struct {
	Planet          Planet     `json:"planet"`
	NameZh          string     `json:"name_zh"`
	Symbol          string     `json:"symbol"`
	IsRetrograde    bool       `json:"is_retrograde"`
	Longitude       float64    `json:"longitude"`
	Speed           float64    `json:"speed"`
	NextStationDate *time.Time `json:"next_station_date,omitempty"`
	StationType     string     `json:"station_type,omitempty"`
}

func GetRetrogradeInfo(planet Planet, t time.Time) (*RetrogradeInfo, error) {
	pos, err := CalculatePosition(planet, t)
	if err != nil {
		return nil, fmt.Errorf("calculate position: %w", err)
	}

	info := &RetrogradeInfo{
		Planet:       planet,
		NameZh:       planet.NameZh(),
		Symbol:       planet.Symbol(),
		IsRetrograde: pos.IsRetrograde,
		Longitude:    pos.Longitude,
		Speed:        pos.Speed,
	}

	stationRetro, stationDirect, err := FindStationDates(planet, t, 365)
	if err != nil {
		return nil, fmt.Errorf("find station dates: %w", err)
	}

	if info.IsRetrograde {
		if stationDirect != nil {
			info.NextStationDate = stationDirect
			info.StationType = "station_direct"
		}
	} else {
		if stationRetro != nil {
			info.NextStationDate = stationRetro
			info.StationType = "station_retrograde"
		}
	}

	return info, nil
}

func FindStationDates(planet Planet, startTime time.Time, searchDays int) (*time.Time, *time.Time, error) {
	pos, err := CalculatePosition(planet, startTime)
	if err != nil {
		return nil, nil, err
	}

	var stationRetro, stationDirect *time.Time

	if pos.IsRetrograde {
		stationDirect = findNextStationDirect(planet, startTime, searchDays)
		stationRetro = findPreviousStationRetrograde(planet, startTime, searchDays)
	} else {
		stationRetro = findNextStationRetrograde(planet, startTime, searchDays)
		stationDirect = findPreviousStationDirect(planet, startTime, searchDays)
	}

	return stationRetro, stationDirect, nil
}

func findNextStationRetrograde(planet Planet, startTime time.Time, maxDays int) *time.Time {
	dayStep := 5
	currentTime := startTime
	endTime := startTime.AddDate(0, 0, maxDays)

	prevPos, _ := CalculatePosition(planet, currentTime)
	if prevPos == nil {
		return nil
	}

	for currentTime.Before(endTime) {
		currentTime = currentTime.AddDate(0, 0, dayStep)
		pos, err := CalculatePosition(planet, currentTime)
		if err != nil {
			continue
		}

		if !prevPos.IsRetrograde && pos.IsRetrograde {
			stationTime := bisectStationDate(planet, currentTime.AddDate(0, 0, -dayStep), currentTime, true)
			return &stationTime
		}
		prevPos = pos
	}

	return nil
}

func findNextStationDirect(planet Planet, startTime time.Time, maxDays int) *time.Time {
	dayStep := 5
	currentTime := startTime
	endTime := startTime.AddDate(0, 0, maxDays)

	prevPos, _ := CalculatePosition(planet, currentTime)
	if prevPos == nil {
		return nil
	}

	for currentTime.Before(endTime) {
		currentTime = currentTime.AddDate(0, 0, dayStep)
		pos, err := CalculatePosition(planet, currentTime)
		if err != nil {
			continue
		}

		if prevPos.IsRetrograde && !pos.IsRetrograde {
			stationTime := bisectStationDate(planet, currentTime.AddDate(0, 0, -dayStep), currentTime, false)
			return &stationTime
		}
		prevPos = pos
	}

	return nil
}

func findPreviousStationRetrograde(planet Planet, startTime time.Time, maxDays int) *time.Time {
	dayStep := 5
	currentTime := startTime
	endTime := startTime.AddDate(0, 0, -maxDays)

	prevPos, _ := CalculatePosition(planet, currentTime)
	if prevPos == nil {
		return nil
	}

	for currentTime.After(endTime) {
		currentTime = currentTime.AddDate(0, 0, -dayStep)
		pos, err := CalculatePosition(planet, currentTime)
		if err != nil {
			continue
		}

		if !pos.IsRetrograde && prevPos.IsRetrograde {
			stationTime := bisectStationDate(planet, currentTime, currentTime.AddDate(0, 0, dayStep), true)
			return &stationTime
		}
		prevPos = pos
	}

	return nil
}

func findPreviousStationDirect(planet Planet, startTime time.Time, maxDays int) *time.Time {
	dayStep := 5
	currentTime := startTime
	endTime := startTime.AddDate(0, 0, -maxDays)

	prevPos, _ := CalculatePosition(planet, currentTime)
	if prevPos == nil {
		return nil
	}

	for currentTime.After(endTime) {
		currentTime = currentTime.AddDate(0, 0, -dayStep)
		pos, err := CalculatePosition(planet, currentTime)
		if err != nil {
			continue
		}

		if pos.IsRetrograde && !prevPos.IsRetrograde {
			stationTime := bisectStationDate(planet, currentTime, currentTime.AddDate(0, 0, dayStep), false)
			return &stationTime
		}
		prevPos = pos
	}

	return nil
}

func bisectStationDate(planet Planet, startTime, endTime time.Time, lookingForRetrograde bool) time.Time {
	low := startTime
	high := endTime
	maxIterations := 20

	for i := 0; i < maxIterations; i++ {
		mid := low.Add(high.Sub(low) / 2)

		if mid.Equal(low) || mid.Equal(high) {
			return mid
		}

		pos, err := CalculatePosition(planet, mid)
		if err != nil {
			return mid
		}

		if lookingForRetrograde {
			if pos.IsRetrograde {
				high = mid
			} else {
				low = mid
			}
		} else {
			if pos.IsRetrograde {
				low = mid
			} else {
				high = mid
			}
		}
	}

	return low.Add(high.Sub(low) / 2)
}

func GetAllRetrogradeInfo(t time.Time) ([]RetrogradeInfo, error) {
	planets := GetAllPlanets()
	results := make([]RetrogradeInfo, 0, len(planets))

	for _, planet := range planets {
		info, err := GetRetrogradeInfo(planet, t)
		if err != nil {
			continue
		}
		results = append(results, *info)
	}

	return results, nil
}

func TimeToJD(t time.Time) float64 {
	pt := celestial.NewPrecisionTime(t)
	return pt.JDE
}
