PACKAGES=$(shell go list ./... | grep -v 'tests')

### Tools needed for development
devtools:
	@echo "Installing devtools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1
	go install mvdan.cc/gofumpt@latest

### Testing
unit_test:
	go test $(PACKAGES)

test:
	go test ./... -covermode=atomic

### Formatting, linting, and vetting
fmt:
	gofumpt -l -w .
	go mod tidy

check:
	golangci-lint run --build-tags "${BUILD_TAG}" --timeout=20m0s