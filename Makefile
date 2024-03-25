PACKAGES=$(shell go list ./... | grep -v 'tests')

# Install tools are needed for development (golangci-lint, gofumpt)
devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1
	go install mvdan.cc/gofumpt@latest

# GoLang Unit-Test
test:
	go test ./... -covermode=atomic

# Format
fmt:
	gofumpt -l -w .
	go mod tidy

# Linter
check:
	golangci-lint run --build-tags "${BUILD_TAG}" --timeout=20m0s
