.PHONY: build test clean lint fmt check mod-tidy all deps

# Default target
all: clean deps lint test build

# Download dependencies
deps:
	go mod tidy
	go mod download

# Build the project
build:
	go build -o bin/ars-client ./...

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-cov:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run linter
lint:
	go vet ./...
	if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping additional linting"; \
	fi

# Format code
fmt:
	go fmt ./...

# Check if code is formatted
check:
	! go fmt ./... | grep -v "^$$"

# Ensure go.mod is tidy
mod-tidy:
	go mod tidy
	git diff --exit-code go.mod go.sum

# Install dev tools
dev-tools:
	go install golang.org/x/lint/golint@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
