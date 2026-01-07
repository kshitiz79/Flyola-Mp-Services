# Flyola Services Backend Makefile

# Variables
BINARY_NAME=flyola-services
MAIN_PATH=./cmd/server
BUILD_DIR=./bin

# Default target
.PHONY: all
all: clean build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	@go run $(MAIN_PATH)/main.go

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
.PHONY: dev
dev:
	@echo "Starting development server with hot reload..."
	@air

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Test database connectivity
.PHONY: test-db
test-db:
	@echo "Testing database connectivity..."
	@go run cmd/test-db/main.go

# Create database
.PHONY: db-create
db-create:
	@echo "Creating database..."
	@mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS flyola_services CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# Drop database (use with caution!)
.PHONY: db-drop
db-drop:
	@echo "⚠️  WARNING: This will delete the database!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		mysql -u root -p -e "DROP DATABASE IF EXISTS flyola_services;"; \
		echo "\n✅ Database dropped"; \
	else \
		echo "\n❌ Cancelled"; \
	fi

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Linting code..."
	@golangci-lint run

# Tidy dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download

# Docker build
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	@docker build -t flyola-services:latest .

# Docker run
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 flyola-services:latest

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  dev           - Run with hot reload (requires air)"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  test-db       - Test database connectivity"
	@echo "  db-create     - Create database"
	@echo "  db-drop       - Drop database (with confirmation)"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code (requires golangci-lint)"
	@echo "  tidy          - Tidy dependencies"
	@echo "  deps          - Install dependencies"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  help          - Show this help message"