# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@go build -o main.exe cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker-run:
	@docker compose up --build

# Shutdown DB container
docker-down:
	@docker compose down

# Database migrations
migrate-up:
	@echo "Running migrations..."
	@go run cmd/migrate/migrate.go up

migrate-down:
	@echo "Rolling back migrations..."
	@go run cmd/migrate/migrate.go down

# Database seeding
seed:
	@echo "Seeding database..."
	@go run cmd/seed/seed.go

# Setup database (migrate + seed)
db-setup: migrate-up seed

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main.exe main

# Live Reload
watch:
	@D:\Project\GoWorkspace\bin\air.exe

dev: watch

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Help
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make dev            - Run with live reload"
	@echo "  make migrate-up     - Run database migrations"
	@echo "  make migrate-down   - Rollback database migrations"
	@echo "  make seed           - Seed database with initial data"
	@echo "  make db-setup       - Setup database (migrate + seed)"
	@echo "  make docker-run     - Start Docker containers"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make deps           - Install dependencies"
	@echo "  make fmt            - Format code"

.PHONY: all build run clean watch dev docker-run docker-down migrate-up migrate-down seed db-setup deps fmt help
