include .env


LOCAL_BIN := $(CURDIR)/bin
GOBIN := $(LOCAL_BIN)
LOCAL_MIGRATION_DIR := $(MIGRATION_DIR)
LOCAL_MIGRATION_DSN := "host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

GOOSE_CMD := ${LOCAL_BIN}/goose


install-deps:
	@[ -f ${LOCAL_BIN}/goose ] || { \
  		echo "Installing goose..."; \
  		GOBIN=${GOBIN} go install github.com/pressly/goose/v3/cmd/goose@v3.14.0; \
  }

local-migration-status: install-deps
	${GOOSE_CMD} -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up: install-deps
	${GOOSE_CMD} -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down: install-deps
	${GOOSE_CMD} -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

local-migration-create: install-deps
	@if [ -z "$(name)" ]; then \
		echo "Please provide a migration name, usage: make local-migration-create name=add_table"; \
		exit 1; \
	fi
	${GOOSE_CMD} -dir ${LOCAL_MIGRATION_DIR} create $(name) sql
