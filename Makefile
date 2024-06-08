MAIN_FILE_VENOM = cmd/venom/main.go
MAIN_FILE_CLI= cmd/cli/main.go
OUTPUT_DIR_VENOM = bin/venom
OUTPUT_DIR_CLI = bin/cli

# Default target
.PHONY: build run

build-venom:
	@go build -o $(OUTPUT_DIR_VENOM)/in_memory_db-GO $(MAIN_FILE_VENOM)

run-venom: build-venom
	@./$(OUTPUT_DIR_VENOM)/in_memory_db-GO --config="config/config.yaml"

build-cli:
	@go build -o $(OUTPUT_DIR_CLI)/in_memory_db-GO $(MAIN_FILE_CLI)

run-cli: build-cli
	@./$(OUTPUT_DIR_CLI)/in_memory_db-GO --address="localhost:3223"

test:
	go test ./...

cover:
	go test -coverprofile=coverage.out ./...

lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

lint:
	golangci-lint run ./... --fix --config=.golangci.yaml

