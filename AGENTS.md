# Lunar-Zenith Project Guidelines

高精度曆法算曆引擎。根目錄 `AGENTS.md` 提供全域預設；若工作位於較深目錄，優先遵循該目錄下更接近的 `AGENTS.md`。

## Scope

- Root defaults apply across the repo.
- Area-specific overrides live in `pkg/celestial/AGENTS.md`, `pkg/zodiac/AGENTS.md`, `internal/service/AGENTS.md`, and `contracts/AGENTS.md`.
- Link to detailed docs instead of copying them: see `README.md` for product and quick start, `PRD.md` for scope, `SKILLS.md` for calendar/astronomy domain knowledge, and `contracts/README.md` for contract workflow.

## Build And Test

- Use `make sync-contracts` before `make dev` whenever port contracts may have changed.
- Main validation commands: `make test`, `make vet`, `make build`, `make verify-all`.
- `make dev` syncs contracts, loads `.env.ports`, and runs `cmd/server/main.go`.
- `make dev-clean` is recovery-only for stale listeners.
- The Makefile runs Go commands with `CGO_LDFLAGS='-Wl,-w'`; keep parity with that environment when reproducing build or test behavior manually.

## Architecture

- `cmd/server/` is the bootstrap layer. Only `cmd/` may use `log.Fatal()` or `panic()`.
- `internal/service/` aggregates domain logic and exposes REST + gRPC handlers.
- `pkg/celestial/` owns Julian Day, Delta-T, solar terms, and lunar phase calculations.
- `pkg/zodiac/` converts astronomical values into lunar calendar, sexagenary, shensha, and religious calendar outputs.
- `pkg/western_astro/` contains western astrology calculations such as retrograde and aspects.
- `contracts/` is the source of truth for runtime ports and API contracts.
- `gen/` and generated protobuf files are generated artifacts; do not hand-edit them.

## Conventions

- Contract-first: update OpenAPI and protobuf contracts before changing handlers or generated APIs.
- Do not manually edit `.env.ports`; use `make sync-contracts` and `make verify-contracts`.
- Keep services stateless and prefer constructor injection patterns already used in `internal/service/aggregator.go`.
- Use Traditional Chinese for calendar-domain names and user-facing cultural labels.
- Use `float64` for astronomical calculations.
- Wrap errors with context and return errors instead of panicking outside `cmd/`.
- For solar and lunar calculations, use JDE/Delta-T corrected time where required; do not use UT directly for astronomical position calculations.
- Treat leap-month handling as a known risk area in lunar-calendar work; validate leap-month years explicitly and consult `pkg/zodiac/AGENTS.md` before making correctness claims for those cases.

## Testing

- Prefer focused package tests while iterating, then finish with the narrowest relevant `make` or `go test` validation.
- Follow existing table-driven test style and use `t.Errorf()` for non-fatal case reporting.
- When changing contracts or transport layers, verify both REST/gRPC code paths and contract sync behavior.

## Where To Look

- `cmd/server/main.go` for startup, dependency wiring, and port loading.
- `internal/service/aggregator.go` for service composition.
- `internal/service/rest_handler.go` and `internal/service/grpc_server.go` for transport behavior.
- `pkg/celestial/base.go`, `pkg/celestial/solar.go`, and `pkg/celestial/moon.go` for astronomy logic.
- `pkg/zodiac/sexagenary.go` and `pkg/zodiac/lunar_engine.go` for calendar conversions.
- `contracts/runtime/ports.env`, `contracts/openapi/lunar-zenith.yaml`, and `proto/lunar.proto` for contract-driven changes.

## Anti-Patterns

- Do not edit generated files in `gen/` or generated protobuf outputs directly.
- Do not bypass contract updates when adding or changing API fields.
- Do not use UT directly for solar or lunar position calculations.
- Do not call `panic()` or `log.Fatal()` outside `cmd/`.
- Do not treat `.env.ports` as a source file.
