# Automatically import and export .env if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: all build run test clean docker-up docker-down docker-logs docker-db

# Default target
all: build

# Run the API server directly on host
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
	@rm -rf bin/ tmp/

# Docker Compose targets
docker-up:
	@docker compose up -d --build

docker-down:
	@docker compose down

docker-logs:
	@docker compose logs -f

docker-db:
	@docker compose up -d db
