# Install development tools (golangci-lint, gofumpt)
devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.1
	go install mvdan.cc/gofumpt@latest

# Tests
integration_test:
	go test ./tests/integration/*

# Format
fmt:
	gofumpt -l -w .
	go mod tidy

# Linter
check:
	golangci-lint run --build-tags "${BUILD_TAG}" --timeout=20m0s

# Run on development
dev:
	# enable gRPC log tools
	export GRPC_GO_LOG_VERBOSITY_LEVEL=99
	export GRPC_GO_LOG_SEVERITY_LEVEL=info

	# run server
	go run cmd/server/server.go

# Build for production
build:
	export GIN_MODE=release
	export ENV=production
	go mod tidy
	go clean -cache
	go build -o ./build/server cmd/server/server.go

# Pre Push Git Hook
pre-push:
	make fmt
	make check
	make test
	make build
	