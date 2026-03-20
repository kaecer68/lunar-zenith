
.PHONY: dev run sync-contracts verify-contracts dev-clean

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

dev-clean:
	@chmod +x scripts/dev-clean.sh
	bash scripts/dev-clean.sh
