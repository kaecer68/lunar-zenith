package western_astro

import (
	"fmt"
	"testing"
	"time"
)

func TestPlanetString(t *testing.T) {
	tests := []struct {
		planet Planet
		want   string
	}{
		{Mercury, "Mercury"},
		{Venus, "Venus"},
		{Mars, "Mars"},
		{Jupiter, "Jupiter"},
		{Saturn, "Saturn"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.planet.String(); got != tt.want {
				t.Errorf("Planet.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanetNameZh(t *testing.T) {
	tests := []struct {
		planet Planet
		want   string
	}{
		{Mercury, "水星"},
		{Venus, "金星"},
		{Mars, "火星"},
		{Jupiter, "木星"},
		{Saturn, "土星"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.planet.NameZh(); got != tt.want {
				t.Errorf("Planet.NameZh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanetSymbol(t *testing.T) {
	symbols := map[Planet]string{
		Mercury: "☿",
		Venus:   "♀",
		Mars:    "♂",
		Jupiter: "♃",
		Saturn:  "♄",
	}

	for planet, want := range symbols {
		t.Run(planet.String(), func(t *testing.T) {
			if got := planet.Symbol(); got != want {
				t.Errorf("Planet.Symbol() = %v, want %v", got, want)
			}
		})
	}
}

func TestCalculatePosition(t *testing.T) {
	if !IsAvailable() {
		t.Skip("Swiss Ephemeris not available - skipping test")
	}

	tests := []struct {
		name   string
		planet Planet
		date   time.Time
	}{
		{
			name:   "Mercury on 2024-01-01",
			planet: Mercury,
			date:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:   "Venus on 2024-01-01",
			planet: Venus,
			date:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:   "Mars on 2024-01-01",
			planet: Mars,
			date:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := CalculatePosition(tt.planet, tt.date)
			if err != nil {
				t.Fatalf("CalculatePosition() error = %v", err)
			}

			if pos.Longitude < 0 || pos.Longitude >= 360 {
				t.Errorf("Longitude out of range [0, 360): %f", pos.Longitude)
			}

			if pos.Distance <= 0 {
				t.Errorf("Distance should be positive: %f", pos.Distance)
			}

			if pos.Planet != tt.planet {
				t.Errorf("Planet mismatch: got %v, want %v", pos.Planet, tt.planet)
			}
		})
	}
}

func TestRetrogradeDetection(t *testing.T) {
	if !IsAvailable() {
		t.Skip("Swiss Ephemeris not available - skipping test")
	}

	tests := []struct {
		name           string
		planet         Planet
		date           time.Time
		wantRetrograde bool
	}{
		{
			name:           "Mercury NOT retrograde on 2024-01-15",
			planet:         Mercury,
			date:           time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
			wantRetrograde: false,
		},
		{
			name:           "Mercury IS retrograde on 2024-04-10",
			planet:         Mercury,
			date:           time.Date(2024, 4, 10, 12, 0, 0, 0, time.UTC),
			wantRetrograde: true,
		},
		{
			name:           "Venus NOT retrograde on 2024-01-15",
			planet:         Venus,
			date:           time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
			wantRetrograde: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := CalculatePosition(tt.planet, tt.date)
			if err != nil {
				t.Fatalf("CalculatePosition() error = %v", err)
			}

			if pos.IsRetrograde != tt.wantRetrograde {
				t.Errorf("IsRetrograde = %v, want %v", pos.IsRetrograde, tt.wantRetrograde)
			}
		})
	}
}

func TestFindStationDates(t *testing.T) {
	if !IsAvailable() {
		t.Skip("Swiss Ephemeris not available - skipping test")
	}

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	stationRetro, stationDirect, err := FindStationDates(Mercury, start, 365)
	if err != nil {
		t.Fatalf("FindStationDates() error = %v", err)
	}

	if stationRetro == nil {
		t.Error("Expected station retrograde date, got nil")
	}

	if stationDirect == nil {
		t.Error("Expected station direct date, got nil")
	}

	if stationDirect == nil {
		t.Error("Expected station direct date, got nil")
	}
}

func TestGetRetrogradeInfo(t *testing.T) {
	if !IsAvailable() {
		t.Skip("Swiss Ephemeris not available - skipping test")
	}

	date := time.Date(2024, 4, 10, 12, 0, 0, 0, time.UTC)
	info, err := GetRetrogradeInfo(Mercury, date)
	if err != nil {
		t.Fatalf("GetRetrogradeInfo() error = %v", err)
	}

	if !info.IsRetrograde {
		t.Error("Expected Mercury to be retrograde on 2024-04-10")
	}

	if info.NameZh != "水星" {
		t.Errorf("NameZh = %v, want 水星", info.NameZh)
	}

	if info.Symbol != "☿" {
		t.Errorf("Symbol = %v, want ☿", info.Symbol)
	}

	if info.NextStationDate == nil {
		t.Error("Expected NextStationDate to be set")
	}
}

func TestGetAllRetrogradeInfo(t *testing.T) {
	if !IsAvailable() {
		t.Skip("Swiss Ephemeris not available - skipping test")
	}

	date := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	infos, err := GetAllRetrogradeInfo(date)
	if err != nil {
		t.Fatalf("GetAllRetrogradeInfo() error = %v", err)
	}

	if len(infos) != 5 {
		t.Errorf("Expected 5 planets, got %d", len(infos))
	}

	for _, info := range infos {
		if info.NameZh == "" {
			t.Error("Empty NameZh in result")
		}
		if info.Symbol == "" {
			t.Error("Empty Symbol in result")
		}
	}
}

func TestNormalizeLongitude(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0},
		{360, 0},
		{720, 0},
		{-90, 270},
		{450, 90},
		{180, 180},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%.1f", tt.input), func(t *testing.T) {
			result := normalizeLongitude(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeLongitude(%.1f) = %.1f, want %.1f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsAvailable(t *testing.T) {
	available := IsAvailable()
	t.Logf("Swiss Ephemeris available: %v", available)
}
