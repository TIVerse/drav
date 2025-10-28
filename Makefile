.PHONY: all fmt lint test bench cover build examples docs clean help

# Variables
BINARY_NAME=drav
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags "-s -w"
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

all: fmt lint test build ## Run all checks and build

fmt: ## Format code
	@echo "==> Formatting code..."
	@gofmt -s -w .
	@$(GO) mod tidy

lint: ## Run linters
	@echo "==> Running linters..."
	@golangci-lint run ./...

test: ## Run tests
	@echo "==> Running tests..."
	@$(GO) test $(GOFLAGS) -race -timeout=5m ./...

bench: ## Run benchmarks
	@echo "==> Running benchmarks..."
	@$(GO) test -bench=. -benchmem -run=^$$ ./...

cover: ## Generate coverage report
	@echo "==> Generating coverage report..."
	@$(GO) test -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report: $(COVERAGE_HTML)"

build: ## Build CLI binary
	@echo "==> Building $(BINARY_NAME)..."
	@$(GO) build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/drav

examples: ## Build all examples
	@echo "==> Building examples..."
	@for dir in examples/*/; do \
		echo "Building $$dir..."; \
		(cd "$$dir" && $(GO) build -o ../../bin/$$(basename $$dir) .); \
	done

docs: ## Build documentation
	@echo "==> Building docs..."
	@cd docs && mkdocs build

clean: ## Clean build artifacts
	@echo "==> Cleaning..."
	@rm -rf bin/ dist/ build/
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@rm -f *.prof *.pprof
	@find . -name "*.test" -delete

install: ## Install CLI binary
	@echo "==> Installing $(BINARY_NAME)..."
	@$(GO) install $(LDFLAGS) ./cmd/drav

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
