package western_astro

import (
	"fmt"
	"math"
	"time"
)

const (
	AspectConjunction = 0.0
	AspectOpposition  = 180.0
	AspectTrine       = 120.0
	AspectSquare      = 90.0
	AspectSextile     = 60.0
	AspectQuincunx    = 150.0

	DefaultOrb = 8.0
)

type AspectType string

const (
	Conjunction AspectType = "合相"
	Opposition  AspectType = "對沖"
	Trine       AspectType = "三合"
	Square      AspectType = "刑克"
	Sextile     AspectType = "六合"
	Quincunx    AspectType = "梅花"
)

type PlanetaryAspect struct {
	Planet1       Planet     `json:"planet1"`
	Planet1Name   string     `json:"planet1_name"`
	Planet1Symbol string     `json:"planet1_symbol"`
	Planet2       Planet     `json:"planet2"`
	Planet2Name   string     `json:"planet2_name"`
	Planet2Symbol string     `json:"planet2_symbol"`
	Aspect        AspectType `json:"aspect"`
	Angle         float64    `json:"angle"`
	Orb           float64    `json:"orb"`
	ExactDate     *time.Time `json:"exact_date,omitempty"`
}

func CalculateAspects(t time.Time) ([]PlanetaryAspect, error) {
	planets := []Planet{Mercury, Venus, Mars, Jupiter, Saturn, Uranus, Neptune}
	positions := make(map[Planet]*PlanetaryPosition)

	for _, p := range planets {
		pos, err := CalculatePosition(p, t)
		if err != nil {
			continue
		}
		positions[p] = pos
	}

	var aspects []PlanetaryAspect

	for i := 0; i < len(planets); i++ {
		for j := i + 1; j < len(planets); j++ {
			p1, p2 := planets[i], planets[j]
			pos1, ok1 := positions[p1]
			pos2, ok2 := positions[p2]
			if !ok1 || !ok2 {
				continue
			}

			aspect := calculateAspectBetween(pos1, pos2)
			if aspect != nil {
				aspects = append(aspects, *aspect)
			}
		}
	}

	return aspects, nil
}

func calculateAspectBetween(pos1, pos2 *PlanetaryPosition) *PlanetaryAspect {
	diff := math.Abs(pos2.Longitude - pos1.Longitude)
	if diff > 180 {
		diff = 360 - diff
	}

	aspects := []struct {
		angle float64
		name  AspectType
	}{
		{AspectConjunction, Conjunction},
		{AspectOpposition, Opposition},
		{AspectTrine, Trine},
		{AspectSquare, Square},
		{AspectSextile, Sextile},
		{AspectQuincunx, Quincunx},
	}

	for _, a := range aspects {
		orb := math.Abs(diff - a.angle)
		if orb <= DefaultOrb {
			return &PlanetaryAspect{
				Planet1:       pos1.Planet,
				Planet1Name:   pos1.Planet.NameZh(),
				Planet1Symbol: pos1.Planet.Symbol(),
				Planet2:       pos2.Planet,
				Planet2Name:   pos2.Planet.NameZh(),
				Planet2Symbol: pos2.Planet.Symbol(),
				Aspect:        a.name,
				Angle:         diff,
				Orb:           orb,
			}
		}
	}

	return nil
}

func GetMajorConjunctions(t time.Time) ([]PlanetaryAspect, error) {
	allAspects, err := CalculateAspects(t)
	if err != nil {
		return nil, err
	}

	var conjunctions []PlanetaryAspect
	for _, a := range allAspects {
		if a.Aspect == Conjunction && a.Orb <= 5.0 {
			conjunctions = append(conjunctions, a)
		}
	}

	return conjunctions, nil
}

func FormatAspectDescription(aspect PlanetaryAspect) string {
	return fmt.Sprintf("%s %s %s (誤差 %.1f°)",
		aspect.Planet1Name,
		aspect.Aspect,
		aspect.Planet2Name,
		aspect.Orb)
}
