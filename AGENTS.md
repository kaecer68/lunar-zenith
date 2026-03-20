# Lunar-Zenith Project Knowledge Base

**Generated:** 2026-03-20  
**Go Version:** 1.25+  
**Contract Version:** lunar-zenith v1.4.0

## OVERVIEW

高精度曆法算曆引擎。天文級精度節氣/朔望月計算，台灣曆法標準對齊。REST + gRPC 雙棧服務。

## STRUCTURE

```
lunar-zenith/
├── cmd/              # Entry points (server, test_grpc)
├── api/v1/           # Protobuf generated code (DO NOT EDIT)
├── pkg/celestial/    # Astronomical: JD, Delta-T, solar/lunar positions
├── pkg/zodiac/       # Cultural calendar: sexagenary, lunar, shensha
├── internal/service/ # Service layer: aggregation, gRPC/REST handlers
├── internal/webui/   # Web UI static files
└── configs/          # Holiday JSON data
```

## BUILD COMMANDS

```bash
# Build server binary
go build -o bin/server ./cmd/server/main.go

# Run server
./bin/server        # REST:8080, gRPC:50051 (env: GRPC_PORT)

# Run tests
go test ./...                    # All tests
go test ./pkg/celestial/...      # Specific package
go test -v ./pkg/zodiac/...      # Verbose output

# Run single test
go test -v -run TestTimeToJD ./pkg/celestial/
go test -v -run TestNewYearSexagenary ./pkg/zodiac/

# Build with optimizations
go build -ldflags="-s -w" -o bin/server ./cmd/server/main.go
```

## CODE STYLE

### Formatting
- Standard Go formatting (`gofmt`)
- No line length limit, but keep readable
- Use tabs for indentation

### Naming Conventions

**Types:** PascalCase, descriptive
```go
type PrecisionTime struct { }
type HolidayService struct { }
type Sexagenary struct { }
```

**Functions:**
- Exported: PascalCase (`NewPrecisionTime`, `GetDaySexagenary`)
- Unexported: camelCase (none found in this codebase)
- Constructor pattern: `New{type}()`
- Getters: `Get{Property}()` or just `{Property}()`

**Variables:**
- Local: short when obvious (`t`, `jd`, `err`)
- Constants: CamelCase or PascalCase for exported
```go
var HeavenlyStems = []string{"甲", "乙", ...}
const TypePublicHoliday HolidayType = "holiday"
```

**Packages:**
- Match directory name exactly
- Single word when possible: `celestial`, `zodiac`, `service`

### Imports
Standard Go import grouping (no blank lines needed in simple cases):
```go
import (
    "math"
    "time"
)
```

External imports:
```go
import (
    "encoding/json"
    "fmt"
    "os"
    
    "github.com/gin-gonic/gin"
    lunarv1 "github.com/kaecer68/lunar-zenith/api/v1"
    "google.golang.org/grpc"
)
```

Order: stdlib → external → internal (project imports)

### Error Handling

**Zero-Panic Rule:** Only `cmd/` may use `log.Fatal()`. All other packages return `error`.

**Error wrapping with context:**
```go
if err != nil {
    return fmt.Errorf("failed to read holiday file: %w", err)
}
```

**No naked returns:** Always return explicit values.

### Types

**Structs:** Document purpose with comment
```go
// PrecisionTime 封裝了天文計算所需的高精度時間結構
type PrecisionTime struct {
    UT     time.Time // 民用協調世界時
    JD     float64   // 儒略日
    JDE    float64   // 儒略曆元
    DeltaT float64   // TT - UT 差值（秒）
}
```

**Constants for enums:**
```go
type HolidayType string
const (
    TypePublicHoliday HolidayType = "holiday"
    TypeWorkday       HolidayType = "workday"
)
```

### Testing

**Test file naming:** `*_test.go` in same package

**Test function naming:** `Test{FunctionName}`
```go
func TestTimeToJD(t *testing.T) {
    t1 := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
    got := TimeToJD(t1)
    want := 2451545.0
    if got != want {
        t.Errorf("TimeToJD(...) = %f; want %f", got, want)
    }
}
```

**Test patterns:**
- Table-driven tests for multiple cases
- Use `t.Errorf()` not `t.Fatalf()` for non-fatal failures
- Test error cases explicitly

## CONVENTIONS

- **Contract-first:** API changes → update OpenAPI → implement
- **Stateless:** All services stateless, support horizontal scaling
- **Zero-Panic:** Only `cmd/` uses `panic`/`log.Fatal`. Return errors.
- **繁體中文:** All calendar names use Traditional Chinese
- **Precision:** `float64` for all astronomical calculations (15-16 digits)

## ANTI-PATTERNS

- ❌ Directly edit `api/v1/*.go` (protoc generated)
- ❌ Add API fields not defined in contract
- ❌ Access `HolidayService` outside `internal/service`
- ❌ Use UT directly for solar/lunar calculations (must use JDE)
- ❌ Ignore Delta-T correction (hours of error on historical dates)
- ❌ Use `panic()` outside `cmd/`

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| API contract | `contracts/openapi/lunar-zenith.yaml` | Update first, then generate |
| Server entry | `cmd/server/main.go` | Gin + gRPC + Web UI |
| Astronomical | `pkg/celestial/*.go` | JD, Delta-T, positions |
| Calendar logic | `pkg/zodiac/*.go` | Sexagenary, lunar, shensha |
| Business logic | `internal/service/*.go` | Aggregator, HolidayService |
| Web UI | `internal/webui/static/` | Built-in web interface |
| Holiday data | `configs/holidays_*.json` | Taiwan/China holidays |

## CODE MAP

| Symbol | File | Role |
|--------|------|------|
| `PrecisionTime` | `pkg/celestial/base.go` | High-precision time wrapper |
| `NewYearSexagenary()` | `pkg/zodiac/sexagenary.go` | Year sexagenary (base year 4 AD) |
| `GetDaySexagenary()` | `pkg/zodiac/sexagenary.go` | Day sexagenary (JD based) |
| `Aggregator` | `internal/service/aggregator.go` | Service aggregator |
| `HolidayService` | `internal/service/holiday.go` | Holiday JSON loader/query |

## RUNNING SINGLE TESTS

```bash
# Pattern: go test -v -run <TestName> <package>
go test -v -run TestTimeToJD ./pkg/celestial/
go test -v -run TestNewPrecisionTime ./pkg/celestial/
go test -v -run TestNewYearSexagenary ./pkg/zodiac/
go test -v -run TestGetDaySexagenary ./pkg/zodiac/
go test -v -run TestSolarLongitude ./pkg/celestial/
```

## RELATED DOCS

- `contracts/README.md` - Contract layer documentation
- `contracts/TASK-BOARD.md` - Cross-service task board
- `PRD.md` - Product requirements
- `pkg/celestial/AGENTS.md` - Celestial module guide
- `pkg/zodiac/AGENTS.md` - Zodiac module guide
- `internal/service/AGENTS.md` - Service layer guide
