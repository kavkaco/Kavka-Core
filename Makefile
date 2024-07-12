# Install development tools
devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.1
	go install mvdan.cc/gofumpt@latest
	go install github.com/bufbuild/buf/cmd/buf@v1.33.0

# Tests
unit_test:
	go test $(shell go list ./... | grep -v /tests)

integration_test:
	go test ./tests/integration/*

test:
	make unit_test
	make integration_test
	make e2e_test

# Format
fmt:
	gofumpt -l -w .
	buf format -w 

# Linter
linter:
	golangci-lint run --build-tags "${BUILD_TAG}" --timeout=20m0s

buf_linter:
	buf lint

# Run on development
dev:
	# set env
	export KAVKA_ENV=development

	# run server
	go run cmd/server/server.go

# Build for production
build:
	# set env
	export KAVKA_ENV=production

	# build
	go build -o ./build/server cmd/server/server.go

gen_protobuf:
	buf generate --path ./protobuf