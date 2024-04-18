.DEFAULT_GOAL := help

.PHONY: help test

SQL_MIGRATE_CLI		= sql-migrate

export UID_GID=$(shell id -u):$(shell id -g)

help: ## show this message
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-8s\033[0m %s\n", $$1, $$2}'

buf-gen: ## buf generate code
	buf generate proto && go mod tidy

go-gen: ## go genrate
	go generate

# sql-migrate targets
migrate-up:	## migrate up
	$(SQL_MIGRATE_CLI) up --config ../dbconfig.yml

# sql-migrate targets
migrate-down:	## migrate down
	$(SQL_MIGRATE_CLI) down --config ../dbconfig.yml

migrate-status:	## migrate status
	$(SQL_MIGRATE_CLI)  status --config ../dbconfig.yml


deps: ## install dependencies
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	go install github.com/envoyproxy/protoc-gen-validate@v0.9.0
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@v1.11.1
	go install github.com/bufbuild/buf/cmd/buf@v1.15.0