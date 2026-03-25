
.PHONY: dev run sync-contracts verify-contracts dev-clean test vet build verify-all

GO_CHECK_ENV = CGO_LDFLAGS='-Wl,-w'

dev:
	@chmod +x scripts/sync-contracts.sh
	bash scripts/sync-contracts.sh
	bash -c 'set -a; . ./.env.ports; set +a; go run ./cmd/server/main.go'

run: dev

sync-contracts:
	@chmod +x scripts/sync-contracts.sh
	bash scripts/sync-contracts.sh

verify-contracts:
	@chmod +x scripts/sync-contracts.sh
	bash scripts/sync-contracts.sh --check

test:
	bash -c "$(GO_CHECK_ENV) go test ./..."

vet:
	bash -c "$(GO_CHECK_ENV) go vet ./..."

build:
	bash -c "$(GO_CHECK_ENV) go build ./..."

verify-all: verify-contracts test vet build

dev-clean:
	@chmod +x scripts/dev-clean.sh
	bash scripts/dev-clean.sh
