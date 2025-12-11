# Makefile for CRM Stomatology

.PHONY: start stop test test-integration test-coverage clean lint fmt mock-gen help

# Start project using docker-compose
start:
	@echo "Starting CRM Stomatology..."
	@docker compose up -d --build
	@echo ""
	@echo "Application started!"
	@echo "UI: http://localhost:8080"
	@echo ""
	@echo "Use 'make stop' to stop the application"

# Stop docker-compose
stop:
	@echo "Stopping CRM Stomatology..."
	@docker compose down
	@echo "Application stopped"

# View logs
logs:
	@docker compose logs -f app

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@go vet ./...
	@echo "Linting complete"

# Generate mocks
mock-gen:
	@echo "Generating mocks..."
	@go generate ./gen/...
	@echo "Mocks generated"

# Run unit tests
test:
	@echo "Running unit tests..."
	@go test -v ./...
	@echo "Running integration tests..."
	@go test -v -tags=integration ./...

# Run integration tests only
test-integration:
	@echo "Running integration tests..."
	@go test -v -tags=integration ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -cover ./...

# Clean build artifacts and docker volumes
clean:
	@echo "Cleaning..."
	@docker compose down -v 2>/dev/null || true
	@go clean
	@echo "Clean complete"

# Help
help:
	@echo "Available commands:"
	@echo "  start           - Start project using docker-compose (UI: http://localhost:8080)"
	@echo "  stop            - Stop docker-compose"
	@echo "  logs            - View application logs"
	@echo "  fmt             - Format code"
	@echo "  lint            - Lint code"
	@echo "  mock-gen        - Generate mocks"
	@echo "  test            - Run all tests (unit + integration)"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-coverage   - Run tests with coverage"
	@echo "  clean           - Clean build artifacts and docker volumes"
	@echo "  help            - Show this help"
