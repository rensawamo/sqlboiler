.DEFAULT_GOAL := help

.PHONY: help test

SQL_MIGRATE_CLI		= sql-migrate

export UID_GID=$(shell id -u):$(shell id -g)

help:	## show this message
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-8s\033[0m %s\n", $$1, $$2}'

lint:	## staticcheck
	staticcheck ./...

test:	## go test
	go test ./...

fmt:	## go fmt
	go fmt ./...

generate: ## go genrate
	go generate

boiler-gen: ## sqlboiler generate
	sqlboiler mysql
	
# sql-migrate targets
migrate-up: ## migrate up
	$(SQL_MIGRATE_CLI) up --config ../dbconfig.yml

# sql-migrate targets
migrate-down: ## migrate down
	$(SQL_MIGRATE_CLI) down --config ../dbconfig.yml


migrate-status: ## migrate status
	$(SQL_MIGRATE_CLI)  status --config ../dbconfig.yml

# internal targets