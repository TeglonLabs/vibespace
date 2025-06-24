# vibespace MCP Experience Makefile
# Production-grade Go project automation

# Project metadata
PROJECT_NAME := vibespace-mcp
PACKAGE := github.com/bmorphism/vibespace-mcp-go
BINARY_NAME := vibespace-mcp
MAIN_PACKAGE := ./cmd/server

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS := -s -w \
	-X main.version=$(VERSION) \
	-X main.commit=$(COMMIT) \
	-X main.buildTime=$(BUILD_TIME)

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOVET := $(GOCMD) vet
GOFMT := gofmt

# Directories
BIN_DIR := bin
DIST_DIR := dist
COVERAGE_DIR := coverage

# Platform detection
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

ifeq ($(UNAME_S),Darwin)
	GOOS := darwin
else ifeq ($(UNAME_S),Linux)
	GOOS := linux
else ifeq ($(OS),Windows_NT)
	GOOS := windows
else
	GOOS := linux
endif

ifeq ($(UNAME_M),x86_64)
	GOARCH := amd64
else ifeq ($(UNAME_M),arm64)
	GOARCH := arm64
else ifeq ($(UNAME_M),aarch64)
	GOARCH := arm64
else
	GOARCH := amd64
endif

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

# Default target
.DEFAULT_GOAL := help

# Phony targets
.PHONY: help build clean test coverage lint fmt vet deps check install \
        release docker docker-build docker-push cross-compile \
        security bench profile tools update-deps

## help: Show this help message
help:
	@echo "$(BLUE)$(PROJECT_NAME) Build System$(NC)"
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2 } \
		/^##@/ { printf "\n$(BLUE)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

## build: Build the binary for current platform
build: deps
	@echo "$(BLUE)Building $(PROJECT_NAME)...$(NC)"
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "$(GREEN)Build complete: $(BIN_DIR)/$(BINARY_NAME)$(NC)"

## clean: Remove build artifacts
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	$(GOCLEAN)
	rm -rf $(BIN_DIR) $(DIST_DIR) $(COVERAGE_DIR)
	rm -f coverage.out coverage.html *.prof *.bench
	@echo "$(GREEN)Clean complete$(NC)"

## install: Install the binary to GOPATH/bin
install: build
	@echo "$(BLUE)Installing $(PROJECT_NAME)...$(NC)"
	$(GOCMD) install -ldflags="$(LDFLAGS)" $(MAIN_PACKAGE)
	@echo "$(GREEN)Install complete$(NC)"

##@ Development

## deps: Download and verify dependencies
deps:
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) verify

## tidy: Clean up dependencies
tidy:
	@echo "$(BLUE)Tidying dependencies...$(NC)"
	$(GOMOD) tidy

## fmt: Format Go code
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GOFMT) -s -w .
	@if command -v goimports >/dev/null 2>&1; then \
		echo "$(BLUE)Running goimports...$(NC)"; \
		goimports -local $(PACKAGE) -w .; \
	else \
		echo "$(YELLOW)goimports not found, skipping...$(NC)"; \
	fi

## vet: Run go vet
vet:
	@echo "$(BLUE)Running go vet...$(NC)"
	$(GOVET) ./...

## lint: Run golangci-lint
lint:
	@echo "$(BLUE)Running golangci-lint...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=10m; \
	else \
		echo "$(RED)golangci-lint not found. Installing...$(NC)"; \
		$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run --timeout=10m; \
	fi

## check: Run all checks (fmt, vet, lint)
check: fmt vet lint
	@echo "$(GREEN)All checks passed$(NC)"

##@ Testing

## test: Run tests
test:
	@echo "$(BLUE)Running tests...$(NC)"
	$(GOTEST) -race -v ./...

## test-short: Run short tests
test-short:
	@echo "$(BLUE)Running short tests...$(NC)"
	$(GOTEST) -short -race ./...

## coverage: Generate test coverage report
coverage:
	@echo "$(BLUE)Generating coverage report...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	$(GOCMD) tool cover -func=coverage.out
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

## bench: Run benchmarks
bench:
	@echo "$(BLUE)Running benchmarks...$(NC)"
	$(GOTEST) -bench=. -benchmem -count=3 ./... | tee benchmark.txt

##@ Security

## security: Run security checks
security:
	@echo "$(BLUE)Running security checks...$(NC)"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "$(YELLOW)gosec not found, installing...$(NC)"; \
		$(GOCMD) install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "$(YELLOW)govulncheck not found, installing...$(NC)"; \
		$(GOCMD) install golang.org/x/vuln/cmd/govulncheck@latest; \
		govulncheck ./...; \
	fi

##@ Release

## cross-compile: Build for multiple platforms
cross-compile: deps
	@echo "$(BLUE)Cross-compiling for multiple platforms...$(NC)"
	@mkdir -p $(DIST_DIR)
	@for os in linux windows darwin; do \
		for arch in amd64 arm64; do \
			if [ "$$os" = "windows" ] && [ "$$arch" = "arm64" ]; then continue; fi; \
			ext=""; \
			if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
			echo "$(YELLOW)Building for $$os/$$arch...$(NC)"; \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 \
				$(GOBUILD) -ldflags="$(LDFLAGS)" \
				-o $(DIST_DIR)/$(BINARY_NAME)-$$os-$$arch$$ext $(MAIN_PACKAGE); \
		done; \
	done
	@echo "$(GREEN)Cross-compilation complete$(NC)"

## release: Create release archives
release: cross-compile
	@echo "$(BLUE)Creating release archives...$(NC)"
	@cd $(DIST_DIR) && \
	for file in $(BINARY_NAME)-*; do \
		if [[ "$$file" == *windows* ]]; then \
			zip "$${file%.exe}.zip" "$$file"; \
		else \
			tar -czf "$$file.tar.gz" "$$file"; \
		fi; \
	done
	@cd $(DIST_DIR) && sha256sum * > checksums.txt
	@echo "$(GREEN)Release archives created in $(DIST_DIR)$(NC)"

##@ Docker

## docker-build: Build Docker image
docker-build:
	@echo "$(BLUE)Building Docker image...$(NC)"
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg COMMIT=$(COMMIT) \
		-t $(PROJECT_NAME):$(VERSION) \
		-t $(PROJECT_NAME):latest \
		.
	@echo "$(GREEN)Docker image built: $(PROJECT_NAME):$(VERSION)$(NC)"

## docker-run: Run Docker container
docker-run: docker-build
	@echo "$(BLUE)Running Docker container...$(NC)"
	docker run --rm -it $(PROJECT_NAME):$(VERSION)

##@ Tools

## tools: Install development tools
tools:
	@echo "$(BLUE)Installing development tools...$(NC)"
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOCMD) install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	$(GOCMD) install golang.org/x/vuln/cmd/govulncheck@latest
	$(GOCMD) install golang.org/x/tools/cmd/goimports@latest
	$(GOCMD) install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "$(GREEN)Development tools installed$(NC)"

## update-deps: Update dependencies to latest versions
update-deps:
	@echo "$(BLUE)Updating dependencies...$(NC)"
	$(GOCMD) get -u ./...
	$(GOMOD) tidy
	@echo "$(GREEN)Dependencies updated$(NC)"

##@ Information

## version: Show version information
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Platform: $(GOOS)/$(GOARCH)"

## info: Show project information
info:
	@echo "$(BLUE)Project Information$(NC)"
	@echo "Name: $(PROJECT_NAME)"
	@echo "Package: $(PACKAGE)"
	@echo "Binary: $(BINARY_NAME)"
	@echo "Main Package: $(MAIN_PACKAGE)"
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Target Platform: $(GOOS)/$(GOARCH)"
	@echo "Go Version: $(shell $(GOCMD) version)"

# Include custom targets if they exist
-include Makefile.local
