include .env
MIGRATION_PATH = ./cmd/migrate/migrations

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATION_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATION_PATH) -database=$(DSN) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATION_PATH) -database=$(DSN) down $(filter-out $@,$(MAKECMDGOALS))
	
.PHONY: migration-version
migration-version:
	@migrate -path=$(MIGRATION_PATH) -database=$(DSN) version

.PHONY: migration-rollback
migration-rollback:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is not set. Use 'make migration-rollback VERSION=<version>'"; \
		exit 1; \
	fi
	@migrate -path=$(MIGRATION_PATH) -database=$(DSN) force $(VERSION)

.PHONY: seed
seed:
	@go run cmd/migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt

.PHONY: test
test:
	@go test -v ./...