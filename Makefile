.PHONY: build test clean install run

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install

# Binary name
BINARY_NAME=accio
BINARY_PATH=./cmd/accio

# Build flags
LDFLAGS=-ldflags "-s -w"

all: test build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) $(BINARY_PATH)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GOBUILD) -o $(BINARY_NAME) $(BINARY_PATH)
	./$(BINARY_NAME)

install:
	$(GOINSTALL) $(LDFLAGS) $(BINARY_PATH)

# Cross-compilation targets
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 $(BINARY_PATH)

build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe $(BINARY_PATH)

build-mac:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 $(BINARY_PATH)

build-all: build-linux build-windows build-mac

# Help target
help:
	@echo "Available targets:"
	@echo "  build       - Build the binary"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  run         - Build and run the binary"
	@echo "  install     - Install the binary"
	@echo "  build-linux - Build for Linux"
	@echo "  build-windows - Build for Windows"
	@echo "  build-mac   - Build for macOS"
	@echo "  build-all   - Build for all platforms"