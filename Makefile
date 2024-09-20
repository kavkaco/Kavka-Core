# Install development tools
devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.1
	go install mvdan.cc/gofumpt@latest
	go install github.com/bufbuild/buf/cmd/buf@v1.33.0

# Tests
unit_test:
	go test -cover $(shell go list ./... | grep -v /tests)

integration_test:
	go test -cover -coverpkg ./internal/service/... ./tests/integration/service

test:
	make unit_test
	make integration_test
	make e2e_test

# Format
fmt:
	gofumpt -l -w .

# Linter
linter:
	golangci-lint run --build-tags "${BUILD_TAG}" --timeout=20m0s

# Run on development
dev:
	export KAVKA_ENV=development
	go run cmd/server/main.go

# Build for production
build:
	export KAVKA_ENV=production
	go build -o ./build/server cmd/server/main.go

gen_protobuf:
	buf generate --path ./protobuf