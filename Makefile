# Install development tools (golangci-lint, gofumpt)
devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.1
	go install mvdan.cc/gofumpt@latest

# Tests
test:
	go test ./... 

# Format
fmt:
	gofumpt -l -w .
	go mod tidy

# Linter
check:
	golangci-lint run --build-tags "${BUILD_TAG}" --timeout=20m0s

# Run on development
dev:
	go run cmd/server/server.go

# Build for production
build:
	export GIN_MODE=release
	export ENV=production
	go mod tidy
	go clean -cache
	go build -o ./build/server cmd/server/server.go
