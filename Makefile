# Makefile for auto-git

BINARY_NAME=auto-git
VERSION?=0.1.0
BUILD_DIR=bin
DIST_DIR=dist
TEMP_DIR=$(DIST_DIR)/auto-git-$(VERSION)
BUILD_TIME=$(shell date +%Y-%m-%dT%H:%M:%S)
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

.PHONY: build clean package

# Build the application
build:
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	go clean
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@rm -f $(BINARY_NAME)
	@echo "Clean complete"

# Package source code as tar.gz
package:
	@echo "Packaging source code..."
	@rm -rf $(TEMP_DIR)
	@mkdir -p $(TEMP_DIR)
	@cp main.go go.mod go.sum Makefile LICENSE README.md .gitignore $(TEMP_DIR)/
	@cp -r internal $(TEMP_DIR)/
	@cp -r Formula $(TEMP_DIR)/
	@cd $(DIST_DIR) && tar -czf auto-git-$(VERSION).tar.gz auto-git-$(VERSION)
	@rm -rf $(TEMP_DIR)
	@echo "Package created: $(DIST_DIR)/auto-git-$(VERSION).tar.gz"
