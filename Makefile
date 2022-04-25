
PHONY: build
build: build-catalog build-updater

PHONY: build-catalog
build-catalog:
	go build -ldflags "-s -w" -o build/catalog ./cmd/catalog/

PHONY: build-updater
build-updater:
	go build -ldflags "-s -w" -o build/updater ./cmd/updater/

PHONY: test
test:
	go test ./...

