# pkg/western_astro — Western Astrology Module

**Domain**: Western astrology calculations: planetary positions, retrograde motion, aspects/conjunctions

## OVERVIEW

Computes western astrological data including planetary positions, retrograde periods, and major aspects. Uses simplified ephemeris calculations.

## STRUCTURE

```
pkg/western_astro/
├── planet.go      # Planet enum, position calculation, planetary metadata
├── retrograde.go  # Retrograde detection, station dates
├── aspects.go     # Aspect calculation, conjunctions, orb detection
└── planet_test.go # Unit tests
```

## WHERE TO LOOK

| Task | File | Notes |
|------|------|-------|
| Planet position | `planet.go:CalculatePosition()` | Ecliptic longitude, latitude, speed |
| Retrograde info | `retrograde.go:GetRetrogradeInfo()` | Station dates, retrograde span |
| All retrogrades | `retrograde.go:GetAllRetrogradeInfo()` | Batch for all planets |
| Aspects | `aspects.go:CalculateAspects()` | Major aspects with orbs |
| Conjunctions | `aspects.go:GetMajorConjunctions()` | Planet-to-planet conjunctions |

## CODE MAP

| Symbol | Type | Role |
|--------|------|------|
| `Planet` | enum | Planet identifiers (Sun..Pluto) |
| `PlanetaryPosition` | struct | Longitude, latitude, speed, sign |
| `CalculatePosition()` | func | Compute position for given planet/time |
| `RetrogradeInfo` | struct | Station direct/retro dates, status |
| `GetRetrogradeInfo()` | func | Retrograde status for single planet |
| `PlanetaryAspect` | struct | Aspect type, orb, applying/separating |
| `CalculateAspects()` | func | All major aspects at given time |

## CONVENTIONS

- **Planet enum**: `Planet` uses `iota`, 0=Sun through 9=Pluto
- **Longitude**: 0-360° ecliptic longitude, 0° = Aries
- **Sign names**: English in code (`Aries`, `Taurus`...), Chinese via `PlanetNameZh()`
- **Orb tolerance**: Default ±6° for major aspects, ±8° for conjunctions
- **Time input**: Standard `time.Time`, converted to JD internally

## ANTI-PATTERNS

- ❌ Hard-code planet indices (use `Planet` enum)
- ❌ Ignore speed sign for retrograde detection (negative = retrograde)
- ❌ Use UT directly without JD conversion for position calc

## TESTING

```bash
go test ./pkg/western_astro/... -v
```

## DEPENDENCIES

- **Internal**: `pkg/celestial/` (JD conversion, solar longitude)
- **External**: None (pure Go)
