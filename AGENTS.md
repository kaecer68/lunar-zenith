# Lunar-Zenith Project Knowledge Base

**Go Version:** 1.25.6 | **Module:** github.com/kaecer68/lunar-zenith

## OVERVIEW

高精度曆法算曆引擎。天文級精度節氣/朔望月計算，台灣曆法標準對齊。REST + gRPC 雙棧服務。

## STRUCTURE

```
lunar-zenith/
├── cmd/server/         # Server entry point
├── cmd/test_grpc/      # gRPC test client
├── api/v1/             # Protobuf generated (DO NOT EDIT)
├── pkg/celestial/      # Astronomical: JD, Delta-T, solar/lunar positions
├── pkg/zodiac/         # Cultural calendar: sexagenary, lunar, shensha, festivals
├── pkg/western_astro/  # Western astrology: retrograde, aspects, planets
├── internal/service/   # Service layer: aggregation, handlers
├── internal/webui/     # Web UI static files
└── configs/            # Holiday JSON data
```

## BUILD & TEST COMMANDS

```bash
# Build
make build                    # Build all packages
make build BIN=bin/server     # Build specific binary
go build -o bin/server ./cmd/server/main.go

# Test
make test                     # Run all tests
go test ./...                 # All packages
go test -v ./pkg/celestial/   # Specific package with verbose

# Single test pattern
go test -v -run TestTimeToJD ./pkg/celestial/
go test -v -run TestNewPrecisionTime ./pkg/celestial/
go test -v -run TestNewYearSexagenary ./pkg/zodiac/
go test -v -run TestGetDaySexagenary ./pkg/zodiac/
go test -v -run TestSolarLongitude ./pkg/celestial/
go test -v -run TestLunarEngine ./pkg/zodiac/

# Lint
make vet                      # Run go vet

# Development
make dev                      # Sync contracts and run server
make sync-contracts           # Sync ports.env from contracts/
make verify-contracts         # CI check - verify ports synced
make dev-clean                # Clean up port listeners

# Verify all CI checks
make verify-all               # Runs verify-contracts + test + vet + build
```

## CODE STYLE

### Formatting
- Standard `gofmt` formatting
- Use tabs for indentation
- No line length limit (keep readable)

### Naming
- **Types:** PascalCase (`PrecisionTime`, `HolidayService`)
- **Functions:** PascalCase exported, camelCase unexported
- **Constructors:** `New{Type}()` pattern
- **Getters:** `Get{Property}()` or just `{Property}()`
- **Variables:** Short when obvious (`t`, `jd`, `err`)
- **Packages:** Single word, match directory (`celestial`, `zodiac`)

### Imports
Order: stdlib → external → internal (project)
```go
import (
    "math"
    "time"
    
    "github.com/gin-gonic/gin"
    "google.golang.org/grpc"
    
    "github.com/kaecer68/lunar-zenith/pkg/celestial"
)
```

### Error Handling
- **Zero-Panic Rule:** Only `cmd/` uses `log.Fatal()` or `panic()`
- Always wrap errors with context: `fmt.Errorf("...: %w", err)`
- Never naked returns

### Types
```go
// Document purpose with comment
type PrecisionTime struct {
    UT     time.Time // 民用協調世界時
    JD     float64   // 儒略日
    DeltaT float64   // TT - UT 差值（秒）
}

// Enum pattern
type HolidayType string
const (
    TypePublicHoliday HolidayType = "holiday"
    TypeWorkday       HolidayType = "workday"
)
```

### Testing
```go
func TestFunctionName(t *testing.T) {
    got := FunctionName(input)
    want := expected
    if got != want {
        t.Errorf("FunctionName(...) = %v; want %v", got, want)
    }
}
```
- Use table-driven tests for multiple cases
- Use `t.Errorf()`, not `t.Fatalf()` for non-fatal failures
- Test error cases explicitly

## CONVENTIONS

- **Contract-first:** API changes → update OpenAPI → implement
- **Stateless:** Services stateless, support horizontal scaling
- **繁體中文:** Calendar names use Traditional Chinese
- **Precision:** `float64` for all astronomical calculations
- **CGO:** Some packages need `CGO_LDFLAGS='-Wl,-w'` (see Makefile)

## ANTI-PATTERNS

- ❌ Directly edit `api/v1/*.go` (protoc generated)
- ❌ Use `panic()` or `log.Fatal()` outside `cmd/`
- ❌ Use UT directly for solar/lunar (must use JDE)
- ❌ Ignore Delta-T correction
- ❌ Manual edit of `.env.ports` (use `make sync-contracts`)

## QUICK REFERENCE

| Task | Location |
|------|----------|
| Server entry | `cmd/server/main.go` |
| REST handler | `internal/service/rest_handler.go` |
| gRPC server | `internal/service/grpc_server.go` |
| Astronomical | `pkg/celestial/*.go` |
| Calendar logic | `pkg/zodiac/*.go` |
| Western astro | `pkg/western_astro/*.go` |
| Holiday data | `configs/holidays_*.json` |
| API contract | `contracts/openapi/lunar-zenith.yaml` |

## RUNTIME PORT CONTRACT

- All REST/gRPC ports defined in `contracts/runtime/ports.env`
- Run `make sync-contracts` to generate `.env.ports`
- Run `make verify-contracts` in CI to check sync
- Never manually edit `.env.ports`
