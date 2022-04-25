
PHONY: build
build: build-catalog build-updater

PHONY: build-catalog
build-catalog:
	go build -ldflags "-s -w" -o bin/catalog ./cmd/catalog/

PHONY: build-updater
build-updater:
	go build -ldflags "-s -w" -o bin/updater ./cmd/updater/

PHONY: test
test:
	go test ./...

