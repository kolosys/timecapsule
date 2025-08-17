.PHONY: test test-coverage build demo clean lint format

# Default target
all: test build

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover ./...

# Run tests with coverage report
test-coverage-html:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build the demo
build:
	go build -o bin/demo cmd/demo/main.go

# Run the demo
demo:
	go run cmd/demo/main.go

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Format code
format:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod tidy
	go mod download

# Generate documentation
docs:
	godoc -http=:6060

# Run benchmarks
bench:
	go test -bench=. -benchmem ./...

# Check for security vulnerabilities
security:
	gosec ./...

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Help
help:
	@echo "Available targets:"
	@echo "  test              - Run tests"
	@echo "  test-coverage     - Run tests with coverage"
	@echo "  test-coverage-html- Run tests with HTML coverage report"
	@echo "  build             - Build the demo"
	@echo "  demo              - Run the demo"
	@echo "  clean             - Clean build artifacts"
	@echo "  format            - Format code"
	@echo "  lint              - Run linter"
	@echo "  deps              - Install dependencies"
	@echo "  docs              - Start documentation server"
	@echo "  bench             - Run benchmarks"
	@echo "  security          - Check for security vulnerabilities"
	@echo "  install-tools     - Install development tools"
	@echo "  help              - Show this help"
