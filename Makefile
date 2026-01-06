.PHONY: build install clean test fmt vet lint help

# Build variables
BINARY_NAME=getho
MAIN_PATH=./cmd
BUILD_DIR=./bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-s -w"
BUILD_FLAGS=-trimpath

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

install: ## Install the binary to $GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $$(go env GOPATH)/bin/$(BINARY_NAME) $(MAIN_PATH)

clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) -s -w .

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

lint: fmt vet ## Run linter (fmt + vet)

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

mod-verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	$(GOMOD) verify

mod-vendor: ## Vendor dependencies
	@echo "Vendoring dependencies..."
	$(GOMOD) vendor

