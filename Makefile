# ZPMeow WhatsApp API - Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOMOD=$(GOCMD) mod
BINARY_NAME=zpmeow
MAIN_PATH=./cmd/server
BIN_DIR=bin

.PHONY: build run tidy up down docs clean

# Build the application
build:
	mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	$(GOCMD) run $(MAIN_PATH)

# Download and tidy dependencies
tidy:
	$(GOMOD) tidy

# Start services with Docker Compose
up:
	docker compose up -d

# Stop Docker Compose services with volumes
down:
	docker compose down -v

# Generate Swagger documentation
docs:
	@if command -v swag > /dev/null; then \
		swag init -g $(MAIN_PATH)/main.go -o ./docs; \
	else \
		echo "Swagger not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# Clean build files
clean:
	rm -rf $(BIN_DIR)
