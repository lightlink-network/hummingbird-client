SHELL=/bin/bash
EXECUTABLE=hb
BUILD_DIR=./build
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64
DARWIN=$(EXECUTABLE)_darwin_amd64
VERSION=$(shell git describe --tags --always)
GO_PATH=$(shell go env GOPATH)
SYSTEM_PATH=$(shell echo $$PATH)
.PHONY: all clean

all: build ## Build hummingbird


build: ## Build binaries
	@echo "--> Starting build for all OS into $(BUILD_DIR) directory"
	@$(MAKE) darwin 
	@$(MAKE) linux 
	@$(MAKE) windows
	@echo "--> Building complete: Hummingbird Version: $(VERSION)"

windows: $(WINDOWS) ## Build for Windows

linux: $(LINUX) ## Build for Linux

darwin: $(DARWIN) ## Build for Darwin (macOS)

$(WINDOWS):
	@mkdir -p $(BUILD_DIR)
	@env GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(WINDOWS) -ldflags="-X main.Version=$(VERSION)" ./cli/main.go

$(LINUX):
	@mkdir -p $(BUILD_DIR)
	@env GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(LINUX) -ldflags="-X main.Version=$(VERSION)" ./cli/main.go

$(DARWIN):
	@mkdir -p $(BUILD_DIR)
	@env GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(DARWIN) -ldflags="-X main.Version=$(VERSION)" ./cli/main.go

install: ## Install binary (mac or linux)
	@echo "--> Installing Hummingbird on your system"
ifeq ($(shell uname -s),Darwin)
	@cp $(BUILD_DIR)/$(DARWIN) $(shell go env GOPATH)/bin/$(EXECUTABLE)
	@echo "--> Installation complete. Run '$(EXECUTABLE) --help' to get started."
	@echo "--> NOTE: Your go path must be in your system path to run the binary via '$(EXECUTABLE)' command."
	@echo "--> Otherwise, you can run the binary via '$(BUILD_DIR)/$(DARWIN) --help' command."
else ifeq ($(shell uname -s),Linux)
	@cp $(BUILD_DIR)$(LINUX) /usr/local/bin/$(EXECUTABLE)
	@echo "--> Installation complete. Run '$(EXECUTABLE) --help' to get started."
else
	@echo "Unsupported operating system. Please install manually."
endif

clean: ## Remove previous build
	@rm -d -r -f $(BUILD_DIR)

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'