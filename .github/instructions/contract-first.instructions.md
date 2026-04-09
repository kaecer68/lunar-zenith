---
description: "Use when changing API contracts, transport handlers, protobuf/OpenAPI fields, or service responses in internal/service and contracts-related files. Enforces contract-first order and sync verification."
name: "Contract-First Workflow"
applyTo:
   - "internal/service/**/*.go"
   - "contracts/**/*"
   - "proto/**/*.proto"
---
# Contract-First Workflow

Scope: API schema, generated models, and transport responses.

- Always update contracts before implementation changes.
- Source of truth for API shape is OpenAPI and protobuf contracts, not handler structs.
- Do not hand-edit generated outputs in `gen/`.

Required order:
1. Update contracts (`contracts/openapi/lunar-zenith.yaml` and/or `proto/lunar.proto`).
2. Sync contract outputs and ports (`make sync-contracts`).
3. Update transport mappings (`internal/service/rest_handler.go`, `internal/service/grpc_server.go`) and service composition if needed.
4. Run validation:
   - `make verify-contracts`
   - `make test`

Additional checks:
- If request/response fields changed, verify both REST and gRPC code paths.
- Keep field semantics consistent with existing naming and domain language (Traditional Chinese labels for calendar-facing content).
- Wrap returned errors with context; do not use `panic()` or `log.Fatal()` outside `cmd/`.

References:
- Root policy: `AGENTS.md`
- Contracts workflow: `contracts/README.md`
- Service layer defaults: `internal/service/AGENTS.md`
