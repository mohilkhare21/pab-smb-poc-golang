# Makefile for Multi-Tenant Admin Portal Backend

# Variables
BINARY_NAME=admin-portal-server
BUILD_DIR=build
MAIN_FILE=cmd/server/main.go

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(shell git describe --tags --always --dirty)"

.PHONY: all build clean test deps run dev docker-build docker-run help

# Default target
all: clean build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run in development mode with hot reload (requires air)
dev:
	@echo "Running in development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/go-delve/delve/cmd/dlv@latest

# Lint the code
lint:
	@echo "Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	$(GOCMD) vet ./...

# Generate API documentation (if using swag)
docs:
	@echo "Generating API documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g $(MAIN_FILE); \
	else \
		echo "swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(BINARY_NAME)

# Create .env file from example
setup-env:
	@echo "Setting up environment file..."
	@if [ ! -f .env ]; then \
		cp configs/env.example .env; \
		echo "Created .env file from example. Please update with your configuration."; \
	else \
		echo ".env file already exists."; \
	fi

# Database migrations (placeholder for future use)
migrate:
	@echo "Running database migrations..."
	@echo "Migrations not implemented yet."

# Seed database with test data
seed:
	@echo "Seeding database with test data..."
	@echo "Seed data not implemented yet."

# Show help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  deps          - Install dependencies"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Run in development mode with hot reload"
	@echo "  install-tools - Install development tools"
	@echo "  lint          - Lint the code"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  docs          - Generate API documentation"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  setup-env     - Create .env file from example"
	@echo "  migrate       - Run database migrations"
	@echo "  seed          - Seed database with test data"
	@echo "  help          - Show this help message"

