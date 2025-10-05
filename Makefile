.PHONY: build install test testacc testrace coverage clean fmt fmtcheck vet lint deps docs install-tools help check

default: help

# Build the provider binary
build:
	go build -o terraform-provider-popsink

# Build for all platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/terraform-provider-popsink_linux_amd64
	GOOS=linux GOARCH=arm64 go build -o bin/terraform-provider-popsink_linux_arm64
	GOOS=darwin GOARCH=amd64 go build -o bin/terraform-provider-popsink_darwin_amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/terraform-provider-popsink_darwin_arm64
	GOOS=windows GOARCH=amd64 go build -o bin/terraform-provider-popsink_windows_amd64.exe

# Install the provider locally for testing
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/popsink/popsink/1.0.0/$$(go env GOOS)_$$(go env GOARCH)
	cp terraform-provider-popsink ~/.terraform.d/plugins/registry.terraform.io/popsink/popsink/1.0.0/$$(go env GOOS)_$$(go env GOARCH)/

# Run unit tests
test:
	go test -v ./...

# Run acceptance tests
testacc:
	TF_ACC=1 go test -v ./... -timeout 120m

# Run tests with race detector
testrace:
	go test -race -v ./...

# Run tests with coverage
coverage:
	go test -v -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts and test cache
clean:
	rm -f terraform-provider-popsink
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -testcache

# Format code
fmt:
	gofmt -s -w .
	goimports -w .

# Check if code is formatted
fmtcheck:
	@echo "==> Checking code formatting..."
	@sh -c "'$(CURDIR)/scripts/fmtcheck.sh'"

# Run go vet
vet:
	go vet ./...

# Run linter
lint:
	golangci-lint run --timeout=5m

# Download and tidy dependencies
deps:
	go mod download
	go mod tidy

# Generate documentation
docs:
	go generate ./...

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

# Run all checks (format, vet, lint, test)
check: fmtcheck vet lint test
	@echo "==> All checks passed!"

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the provider binary"
	@echo "  build-all     - Build for all platforms"
	@echo "  install       - Install the provider locally for testing"
	@echo "  test          - Run unit tests"
	@echo "  testacc       - Run acceptance tests"
	@echo "  testrace      - Run tests with race detector"
	@echo "  coverage      - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts and test cache"
	@echo "  fmt           - Format code"
	@echo "  fmtcheck      - Check if code is formatted"
	@echo "  vet           - Run go vet"
	@echo "  lint          - Run golangci-lint"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  docs          - Generate documentation"
	@echo "  install-tools - Install development tools"
	@echo "  check         - Run all checks (format, vet, lint, test)"
	@echo "  help          - Show this help message"
