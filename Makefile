.PHONY: help run build clean deps db-create db-drop db-reset db-test up down dev dev-local

# Variables
APP_NAME=zpmeow
DB_NAME=zpmeow
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5432

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: ## Run the application
	go run ./cmd/server/main.go

build: ## Build the application
	go build -o bin/$(APP_NAME) ./cmd/server/main.go

clean: ## Clean build artifacts
	rm -rf bin/

deps: ## Download dependencies
	go mod tidy

swagger: ## Generate swagger documentation
	swag init -g cmd/server/main.go -o docs/

# Docker commands (optional)
up: ## Start services with docker compose
	@if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then \
		docker compose up -d; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose up -d; \
	else \
		echo "❌ docker compose not found. Install Docker or use 'make dev-local'"; \
		exit 1; \
	fi

down: ## Stop services with docker compose
	@if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then \
		docker compose down; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose down; \
	else \
		echo "❌ docker compose not found"; \
		exit 1; \
	fi

dev: up ## Start development environment (docker + app)
	@echo "Waiting for database to be ready..."
	@sleep 5
	@make db-create || echo "Database might already exist"
	@echo "Starting application..."
	@make run

dev-local: ## Start development with local PostgreSQL
	@echo "🚀 Starting development with local PostgreSQL..."
	@make db-create || echo "Database might already exist"
	@echo "Starting application..."
	@make run

# Database commands
db-create: ## Create database
	@echo "Creating database $(DB_NAME)..."
	@if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then \
		docker compose exec -T postgres psql -U $(DB_USER) -d postgres -c "CREATE DATABASE $(DB_NAME);"; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose exec -T postgres psql -U $(DB_USER) -d postgres -c "CREATE DATABASE $(DB_NAME);"; \
	else \
		echo "❌ docker compose not found. Install Docker first."; \
		exit 1; \
	fi

db-drop: ## Drop database
	@echo "Dropping database $(DB_NAME)..."
	@if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then \
		docker compose exec -T postgres psql -U $(DB_USER) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose exec -T postgres psql -U $(DB_USER) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"; \
	else \
		echo "❌ docker compose not found. Install Docker first."; \
		exit 1; \
	fi

db-reset: db-drop db-create ## Reset database (drop and create)

db-test: ## Test database connection
	@echo "Testing database connection..."
	@if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then \
		docker compose exec -T postgres psql -U $(DB_USER) -d $(DB_NAME) -c "SELECT version();"; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose exec -T postgres psql -U $(DB_USER) -d $(DB_NAME) -c "SELECT version();"; \
	else \
		echo "❌ docker compose not found. Install Docker first."; \
		exit 1; \
	fi

# DBGate commands
dbgate-up: ## Start DBGate database management tool
	@echo "🚀 Starting DBGate database management tool..."
	@if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then \
		docker compose up -d dbgate; \
		echo "✅ DBGate started successfully!"; \
		echo "🌐 Access DBGate at: http://localhost:3000"; \
		echo "📊 Database: zpmeow"; \
		echo "👤 User: postgres"; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose up -d dbgate; \
		echo "✅ DBGate started successfully!"; \
		echo "🌐 Access DBGate at: http://localhost:3000"; \
		echo "📊 Database: zpmeow"; \
		echo "👤 User: postgres"; \
	else \
		echo "❌ docker compose not found. Install Docker first."; \
		exit 1; \
	fi

dbgate-down: ## Stop DBGate
	@echo "🛑 Stopping DBGate..."
	@if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then \
		docker compose stop dbgate; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose stop dbgate; \
	else \
		echo "❌ docker compose not found"; \
		exit 1; \
	fi

dbgate-logs: ## Show DBGate logs
	@if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then \
		docker compose logs -f dbgate; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose logs -f dbgate; \
	else \
		echo "❌ docker compose not found"; \
		exit 1; \
	fi
