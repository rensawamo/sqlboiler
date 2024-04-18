.DEFAULT_GOAL := help

.PHONY: help test

SQL_MIGRATE_CLI		= sql-migrate

export UID_GID=$(shell id -u):$(shell id -g)

help: ## show this message
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-8s\033[0m %s\n", $$1, $$2}'

# mysql targets
up: 		## Up
	docker compose -f database/docker-compose.yml up --force-recreate -d

down: 	## Down
	docker compose -f database/docker-compose.yml down

logs:	## Show logs
	docker compose -f database/docker-compose.yml logs -f

login:	## login db
	docker compose -f database/docker-compose.yml exec db /bin/bash

clean: 		## clean
clean: down
	rm -rf database/mysql/data/*

mysql-client:	## connet mysql from mysql cli
	docker compose -f docker-compose.yml exec db /bin/bash -c "mysql -u root -p -D db"

# internal targets
