.PHONY: help build run test test-unit test-integration clean swagger install

help:
	@echo "Available commands:"
	@echo "  make install           - Install dependencies"
	@echo "  make build            - Build the application"
	@echo "  make run              - Run the application"
	@echo "  make test             - Run all tests"
	@echo "  make test-unit        - Run unit tests"
	@echo "  make test-integration - Run integration tests"
	@echo "  make swagger          - Generate Swagger documentation"
	@echo "  make clean            - Clean build artifacts"

install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	go install github.com/swaggo/swag/cmd/swag@latest

build:
	@echo "Building application..."
	CGO_ENABLED=1 go build -o bin/api cmd/api/main.go

run:
	@echo "Running application..."
	go run cmd/api/main.go

test:
	@echo "Running all tests..."
	go test -v ./...

test-unit:
	@echo "Running unit tests..."
	go test -v ./tests/unit/...

test-integration:
	@echo "Running integration tests..."
	go test -v ./tests/integration/...

test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/api/main.go -o docs

clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf docs/
	rm -f coverage.out coverage.html
	rm -f database/*.db

docker-build:
	@echo "Building Docker image..."
	docker build -t cliente-api:latest .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 cliente-api:latest
