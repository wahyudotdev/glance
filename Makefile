# Makefile for Agent Proxy

BINARY_NAME=glance
FRONTEND_DIR=web/dashboard
BACKEND_STATIC_DIR=internal/apiserver/dist
BUILD_DIR=build
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Build flags
LDFLAGS=-ldflags="-s -w -X glance/internal/config.Version=$(VERSION)"

.PHONY: all build build-frontend build-backend clean build-all build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-windows-amd64 test test-coverage test-coverage-percent load-test lint

all: build

build: build-frontend build-backend

build-frontend:
	@echo "Building React frontend..."
	rm -rf $(FRONTEND_DIR)/dist
	cd $(FRONTEND_DIR) && npm install --no-save && npm run build
	@echo "Copying frontend assets to backend..."
	rm -rf $(BACKEND_STATIC_DIR)
	mkdir -p $(BACKEND_STATIC_DIR)
	cp -r $(FRONTEND_DIR)/dist/* $(BACKEND_STATIC_DIR)/
	touch $(BACKEND_STATIC_DIR)/.gitkeep

build-backend:
	@echo "Building Go binary for current platform..."
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/glance

build-all: build-frontend build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-windows-amd64

build-darwin-amd64:
	@echo "Building for Darwin AMD64..."
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/glance

build-darwin-arm64:
	@echo "Building for Darwin ARM64..."
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/glance

build-linux-amd64:
	@echo "Building for Linux AMD64..."
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/glance

build-windows-amd64:
	@echo "Building for Windows AMD64..."
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/glance

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -rf $(BACKEND_STATIC_DIR)
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

test:
	@echo "Running tests..."
	go test ./internal/...

test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./internal/...
	@go tool cover -func=coverage.out | grep total | awk '{print "Total Coverage: " $$3}'
	go tool cover -html=coverage.out

test-coverage-percent:
	@echo "Calculating test coverage..."
	@go test -coverprofile=coverage.out ./internal/... > /dev/null
	@go tool cover -func=coverage.out | grep total | awk '{print "Total Coverage: " $$3}'

load-test:
	@echo "Running load test (100 reqs, 10 concurrency)..."
	go run scripts/load-test.go -n 100 -c 10

lint: lint-backend lint-frontend

lint-backend:
	@echo "Running backend linter..."
	golangci-lint run

lint-frontend:
	@echo "Running frontend linter..."
	cd $(FRONTEND_DIR) && npm run lint

run: build
	./$(BINARY_NAME)
