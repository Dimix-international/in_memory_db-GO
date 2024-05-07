MAIN_FILE = cmd/main.go
OUTPUT_DIR = bin

# Default target
.PHONY: build run

build:
	@go build -o $(OUTPUT_DIR)/in_memory_db-GO $(MAIN_FILE)

run: build
	@./$(OUTPUT_DIR)/in_memory_db-GO

test:
	go test ./...

cover:
	go test -coverprofile=coverage.out ./...

lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

lint:
	golangci-lint run ./... --fix --config=.golangci.yaml

