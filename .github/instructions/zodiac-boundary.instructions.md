---
description: "Use when editing lunar calendar conversion, leap month logic, Chinese calendar boundary handling, sexagenary mapping, or pkg/zodiac date outputs. Enforces boundary-year verification and JDE-based astronomical correctness."
name: "Zodiac Boundary Validation"
applyTo: "pkg/zodiac/**/*.go"
---
# Zodiac Boundary Validation

Scope: calendar-boundary and leap-month correctness for `pkg/zodiac` changes.

- Treat leap-month and year-boundary logic as high-risk changes.
- Do not claim all-year correctness without explicit boundary-year validation.
- If outputs depend on astronomical positions, ensure upstream calculations use JDE (not UT directly): `JDE = JD + DeltaT/86400`.
- For changes touching `LunarDate.Month`, `LunarDate.Day`, `IsLeap`, or lunar new-year boundaries, run both:
  - `go test ./pkg/zodiac -run "TestLunarEngineLeap|Test.*Leap" -v`
  - `go test ./pkg/zodiac -run "TestLunarEngineEdge|Test.*Edge" -v`
- Keep table-driven tests and use `t.Errorf()` for non-fatal case failures.
- When adding or updating edge cases, prefer adding explicit year/date fixtures instead of broad assumptions.

References:
- Root defaults: `AGENTS.md`
- Area-specific details: `pkg/zodiac/AGENTS.md`
- Boundary tests: `pkg/zodiac/lunar_engine_leap_test.go`, `pkg/zodiac/lunar_engine_edge_test.go`
