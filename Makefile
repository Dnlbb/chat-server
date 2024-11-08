include chat.env
export

LOCAL_BIN:=$(CURDIR)/bin
GOOSE_CMD=${LOCAL_BIN}/goose


install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	$(LOCAL_BIN)/golangci-lint cache clean
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml



install-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
	go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	@[ -f ${LOCAL_BIN}/goose ] || { \
      		echo "Installing goose..."; \
      		GOBIN=${LOCAL_BIN} go install github.com/pressly/goose/v3/cmd/goose@v3.14.0; \
      }


get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make generate-api


generate-api:
	mkdir -p pkg/chat_v1
	protoc --proto_path=./api/proto/chat_v1 \
	--go_out=pkg/chat_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/chat_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	./api/proto/chat_v1/chat.proto


build:
	GOOS=linux GOARCH=amd64 go build -o auth_linux cmd/main.go


local-migration-status:
	${GOOSE_CMD} -dir ${MIGRATION_DIR} postgres ${PG_DSN} status -v

local-migration-up:
	${GOOSE_CMD} -dir ${MIGRATION_DIR} postgres ${PG_DSN} up -v

local-migration-down:
	${GOOSE_CMD} -dir ${MIGRATION_DIR} postgres ${PG_DSN} down -v

local-migration-create:
	@if [ -z "$(name)" ]; then \
		echo "Please provide a migration name, usage: make local-migration-create name=add_table"; \
		exit 1; \
	fi
	${GOOSE_CMD} -dir ${MIGRATION_DIR} create $(name) sql


docker-up:
	docker compose -f ./deploy/docker-compose.yaml up -d

docker-down:
	docker compose -f ./deploy/docker-compose.yaml down