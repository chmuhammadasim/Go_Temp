# Makefile for Go Backend

# Variables
BINARY_NAME=server
CMD_DIR=cmd/server
BUILD_DIR=bin

# Build the application
build:
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go

# Run the application
run:
	@echo "Running application..."
	go run $(CMD_DIR)/main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf *.db

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Security check
security:
	@echo "Running security checks..."
	gosec ./...

# Generate API documentation
docs:
	@echo "Generating API documentation..."
	swag init -g cmd/server/main.go

# Development setup
dev-setup: deps
	@echo "Setting up development environment..."
	@cp .env.example .env
	@echo "Please edit .env file with your configuration"

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t go-backend .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env go-backend

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  security      - Run security checks"
	@echo "  docs          - Generate API documentation"
	@echo "  dev-setup     - Setup development environment"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  help          - Show this help message"

.PHONY: build run clean deps test test-coverage fmt lint security docs dev-setup docker-build docker-run help