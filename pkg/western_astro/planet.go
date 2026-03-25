package western_astro

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kaecer68/lunar-zenith/pkg/celestial"
	swisseph "github.com/tejzpr/go-swisseph"
)

type Planet int

const (
	Sun     Planet = 0
	Moon    Planet = 1
	Mercury Planet = 2
	Venus   Planet = 3
	Mars    Planet = 4
	Jupiter Planet = 5
	Saturn  Planet = 6
	Uranus  Planet = 7
	Neptune Planet = 8
	Pluto   Planet = 9
)

func (p Planet) String() string {
	switch p {
	case Sun:
		return "Sun"
	case Moon:
		return "Moon"
	case Mercury:
		return "Mercury"
	case Venus:
		return "Venus"
	case Mars:
		return "Mars"
	case Jupiter:
		return "Jupiter"
	case Saturn:
		return "Saturn"
	case Uranus:
		return "Uranus"
	case Neptune:
		return "Neptune"
	case Pluto:
		return "Pluto"
	default:
		return "Unknown"
	}
}

func (p Planet) NameZh() string {
	switch p {
	case Sun:
		return "太陽"
	case Moon:
		return "月亮"
	case Mercury:
		return "水星"
	case Venus:
		return "金星"
	case Mars:
		return "火星"
	case Jupiter:
		return "木星"
	case Saturn:
		return "土星"
	case Uranus:
		return "天王星"
	case Neptune:
		return "海王星"
	case Pluto:
		return "冥王星"
	default:
		return "未知"
	}
}

func (p Planet) Symbol() string {
	switch p {
	case Sun:
		return "☉"
	case Moon:
		return "☽"
	case Mercury:
		return "☿"
	case Venus:
		return "♀"
	case Mars:
		return "♂"
	case Jupiter:
		return "♃"
	case Saturn:
		return "♄"
	case Uranus:
		return "♅"
	case Neptune:
		return "♆"
	case Pluto:
		return "♇"
	default:
		return "?"
	}
}

func (p Planet) toSwissPlanet() int32 {
	switch p {
	case Sun:
		return swisseph.Sun
	case Moon:
		return swisseph.Moon
	case Mercury:
		return swisseph.Mercury
	case Venus:
		return swisseph.Venus
	case Mars:
		return swisseph.Mars
	case Jupiter:
		return swisseph.Jupiter
	case Saturn:
		return swisseph.Saturn
	case Uranus:
		return swisseph.Uranus
	case Neptune:
		return swisseph.Neptune
	case Pluto:
		return swisseph.Pluto
	default:
		return -1
	}
}

type PlanetaryPosition struct {
	Planet       Planet
	Longitude    float64
	Latitude     float64
	Speed        float64
	IsRetrograde bool
	Distance     float64
	JD           float64
}

func CalculatePosition(planet Planet, t time.Time) (*PlanetaryPosition, error) {
	if err := ensureEphemerisPath(); err != nil {
		return nil, fmt.Errorf("ephemeris initialization failed: %w", err)
	}

	pt := celestial.NewPrecisionTime(t)
	jd := pt.JDE

	result := swisseph.CalcUT(jd, planet.toSwissPlanet(),
		swisseph.FlagSwieph|swisseph.FlagSpeed)

	if result.Flag < 0 {
		return nil, fmt.Errorf("swisseph calculation failed with flag: %d", result.Flag)
	}

	return &PlanetaryPosition{
		Planet:       planet,
		Longitude:    normalizeLongitude(result.Data[0]),
		Latitude:     result.Data[1],
		Speed:        result.Data[3],
		IsRetrograde: result.Data[3] < 0,
		Distance:     result.Data[2],
		JD:           jd,
	}, nil
}

func normalizeLongitude(lon float64) float64 {
	lon = fmod(lon, 360.0)
	if lon < 0 {
		lon += 360.0
	}
	return lon
}

func fmod(x, y float64) float64 {
	return x - float64(int64(x/y))*y
}

var ephemerisPathSet = false

func ensureEphemerisPath() error {
	if ephemerisPathSet {
		return nil
	}

	paths := []string{
		"./data/ephe",
		"../data/ephe",
		"../../data/ephe",
	}

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		paths = append(paths,
			filepath.Join(exeDir, "data/ephe"),
			filepath.Join(exeDir, "../data/ephe"),
		)
	}

	for _, path := range paths {
		if isValidEphemerisPath(path) {
			swisseph.SetEphePath(path)
			ephemerisPathSet = true
			return nil
		}
	}

	swisseph.SetEphePath("")
	ephemerisPathSet = true
	return nil
}

func isValidEphemerisPath(path string) bool {
	seplPath := filepath.Join(path, "sepl_18.se1")
	if _, err := os.Stat(seplPath); err == nil {
		return true
	}
	seplPath = filepath.Join(path, "sepl_12.se1")
	if _, err := os.Stat(seplPath); err == nil {
		return true
	}
	return false
}

func IsAvailable() bool {
	if err := ensureEphemerisPath(); err != nil {
		return false
	}

	_, err := CalculatePosition(Sun, time.Now())
	return err == nil
}

func GetAllPlanets() []Planet {
	return []Planet{Mercury, Venus, Mars, Jupiter, Saturn}
}
