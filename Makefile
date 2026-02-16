# Makefile for Agent Proxy

BINARY_NAME=agent-proxy
FRONTEND_DIR=web/dashboard
BACKEND_STATIC_DIR=internal/api/dist

.PHONY: all build build-frontend build-backend clean

all: build

build: build-frontend build-backend

build-frontend:
	@echo "Building React frontend..."
	cd $(FRONTEND_DIR) && npm install && npm run build
	@echo "Copying frontend assets to backend..."
	rm -rf $(BACKEND_STATIC_DIR)
	mkdir -p $(BACKEND_STATIC_DIR)
	cp -r $(FRONTEND_DIR)/dist/* $(BACKEND_STATIC_DIR)/

build-backend:
	@echo "Building Go binary..."
	go build -o $(BINARY_NAME) ./cmd/agent-proxy

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -rf $(BACKEND_STATIC_DIR)
	rm -rf $(FRONTEND_DIR)/dist

run: build
	./$(BINARY_NAME)
