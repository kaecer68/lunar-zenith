# Lunar-Zenith Project Guidelines

高精度曆法算曆引擎。根目錄 `AGENTS.md` 提供全域預設；若工作位於較深目錄，優先遵循該目錄下更接近的 `AGENTS.md`。

## Scope

- Root defaults apply across the repo.
- Area-specific overrides live in `pkg/celestial/AGENTS.md`, `pkg/zodiac/AGENTS.md`, `internal/service/AGENTS.md`, `contracts/AGENTS.md`, and `contracts/openapi/AGENTS.md`.
- `contracts/.github/copilot-instructions.md` defines additional contract-workflow defaults inside the `contracts/` subtree.

## Documentation Map

- Link to detailed docs instead of copying them.
- `README.md`: product overview, quick start, API usage examples.
- `PRD.md`: product scope and release evolution.
- `SKILLS.md`: domain rules and astronomy/calendar background.
- `contracts/README.md`: contract-first workflow and contract sync responsibilities.

## Build And Test

- Preferred local sequence:
	1. `make sync-contracts`
	2. `make dev`
- Validation commands: `make check-docs-consistency`, `make test`, `make vet`, `make build`, `make verify-all`.
- `make verify-all` is the main integration gate (`verify-contracts` + `check-docs-consistency` + `test` + `vet` + `build`).
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
- For gRPC contract governance, keep `proto/lunar.proto` and `contracts/proto/lunar.proto` byte-for-byte synchronized.
- Keep `cmd/server/main.go`, `contracts/openapi/lunar-zenith.yaml`, and the README version badge aligned for every release.
- Keep REST nested JSON fields in snake_case; do not reintroduce PascalCase child keys under `shen_sha`, `mansion`, `daily_deity`, `fetal_god`, or `clash_sha`.
- Do not manually edit `.env.ports`; use `make sync-contracts` and `make verify-contracts`.
- Keep services stateless and prefer constructor injection patterns already used in `internal/service/aggregator.go`.
- Use Traditional Chinese for calendar-domain names and user-facing cultural labels.
- Use `float64` for astronomical calculations.
- Wrap errors with context and return errors instead of panicking outside `cmd/`.
- For solar and lunar calculations, use JDE/Delta-T corrected time where required (`JDE = JD + DeltaT/86400`); do not use UT directly for astronomical position calculations.
- Treat leap-month handling as a known risk area in lunar-calendar work; validate leap-month years explicitly and consult `pkg/zodiac/AGENTS.md` before making correctness claims for those cases.
- Do not hard-code year-specific lunar/leap fixes in production logic; year tables are for verification fixtures only.

## Testing

- Prefer focused package tests while iterating, then finish with the narrowest relevant `make` or `go test` validation.
- Follow existing table-driven test style and use `t.Errorf()` for non-fatal case reporting.
- When changing contracts or transport layers, verify both REST/gRPC code paths and contract sync behavior.
- For leap-month or calendar-boundary changes, explicitly run `pkg/zodiac/lunar_engine_leap_test.go` and `pkg/zodiac/lunar_engine_edge_test.go`.
- For leap-month or solar-term boundary changes that affect API output, also run `internal/service/lunar_leap_api_test.go` and finish with `make verify-all`.

## Common Pitfalls

- Do not run `make dev` before `make sync-contracts` when contract ports may be stale.
- Do not update `proto/lunar.proto` without syncing `contracts/proto/lunar.proto` and rerunning `make check-docs-consistency`.
- Do not document or ship PascalCase nested REST examples; treat them as stale contract drift.
- Do not hand-edit `.env.ports`; regenerate via `make sync-contracts`.
- Do not claim leap-month correctness without explicit boundary-year validation.
- Do not claim leap-month/solar-term correctness unless boundary fixtures are hard assertions (not skip/todo) and cross-layer tests pass.
- If ports are already occupied, use `make dev-clean` before restarting local services.

## Where To Look

- `cmd/server/main.go` for startup, dependency wiring, and port loading.
- `internal/service/aggregator.go` for service composition.
- `internal/service/rest_handler.go` and `internal/service/grpc_server.go` for transport behavior.
- `pkg/celestial/base.go`, `pkg/celestial/solar.go`, and `pkg/celestial/moon.go` for astronomy logic.
- `pkg/zodiac/sexagenary.go` and `pkg/zodiac/lunar_engine.go` for calendar conversions.
- `contracts/runtime/ports.env`, `contracts/openapi/lunar-zenith.yaml`, and `proto/lunar.proto` for contract-driven changes.
- `scripts/sync-contracts.sh` for `.env.ports` generation and contract-sync behavior.

## Anti-Patterns

- Do not edit generated files in `gen/` or generated protobuf outputs directly.
- Do not bypass contract updates when adding or changing API fields.
- Do not use UT directly for solar or lunar position calculations.
- Do not call `panic()` or `log.Fatal()` outside `cmd/`.
- Do not treat `.env.ports` as a source file.
