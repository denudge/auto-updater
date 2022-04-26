include .env

export POSTGRES_USER

PHONY: build
build:
	make build-catalog
	make build-updater

PHONY: build-catalog
build-catalog:
	go build -ldflags "-s -w" -o bin/catalog ./cmd/catalog/

PHONY: build-updater
build-updater:
	go build -ldflags "-s -w" -o bin/updater ./cmd/updater/

.env:
	cp .env.dist .env

PHONY: start
start:
	docker-compose up -d

PHONY: stop
stop:
	docker-compose down --remove-orphans

PHONY: restart
restart:
	make stop
	make start

PHONY: logs
logs:
	docker-compose logs --tail=200 -f

PHONY: psql
psql:
	docker-compose exec postgres psql -U ${POSTGRES_USER}

PHONY: test
test:
	go test ./...
