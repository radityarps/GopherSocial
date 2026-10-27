# Automatically import and export .env if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: all build run test clean

# Default target
all: build

# Run the API server
run:
	@go run ./cmd/api

# Build the API binary
build:
	@go build -o bin/api ./cmd/api

# Run tests
test:
	@go test -v ./...

# Clean build artifacts
clean:
	@rm -rf bin/
