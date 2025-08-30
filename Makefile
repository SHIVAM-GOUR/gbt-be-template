# MAIN COMMANDS
# make dev (Start development with hot reload)
# make dev-local (Start development locally without Docker)
# make test (Run tests)
# make build (Build for production)
# make docker-run (Use Docker - full stack)
# make docker-db (Start only PostgreSQL database)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=gbt-be-template
BINARY_UNIX=$(BINARY_NAME)_unix

# Docker parameters
DOCKER_IMAGE=gbt-be-template
DOCKER_TAG=latest

# Database parameters
DB_URL=postgres://postgres:password@localhost:5432/gbt_template?sslmode=disable

.PHONY: help build clean test deps run dev docker-build docker-run docker-stop migrate-up migrate-down migrate-create air-install

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	$(GOBUILD) -o bin/$(BINARY_NAME) -v ./cmd/app

build-linux: ## Build the application for Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_UNIX) -v ./cmd/app

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)
	rm -f bin/$(BINARY_UNIX)

test: ## Run tests
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

run: build ## Build and run the application
	./bin/$(BINARY_NAME)

dev: ## Run the application with Air (hot reload)
	air

dev-local: ## Run locally without Docker (requires local PostgreSQL)
	@echo "Starting local development..."
	@echo "Make sure PostgreSQL is running and database exists"
	@echo "Run 'make migrate-up' if you haven't run migrations"
	air

docker-build: ## Build Docker image
	sudo docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run application in Docker with docker-compose
	sudo docker-compose up

docker-stop: ## Stop Docker containers
	sudo docker-compose down

docker-logs: ## Show Docker logs
	sudo docker-compose logs -f

docker-db: ## Start only PostgreSQL database
	sudo docker-compose -f docker-compose.db.yml up -d

migrate-up: ## Run database migrations up
	migrate -path migrations -database "$(DB_URL)" up

migrate-down: ## Run database migrations down
	migrate -path migrations -database "$(DB_URL)" down

migrate-create: ## Create a new migration file (usage: make migrate-create name=migration_name)
	migrate create -ext sql -dir migrations -seq $(name)

air-install: ## Install Air for hot reload
	go install github.com/air-verse/air@latest

migrate-install: ## Install golang-migrate
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

tools: air-install migrate-install ## Install development tools

setup: deps tools ## Setup development environment
	@echo "Development environment setup complete!"
	@echo "Run 'make dev' to start the application with hot reload"

# Database operations
db-create: ## Create database
	docker exec -it gbt-postgres createdb -U postgres gbt_template

db-drop: ## Drop database
	docker exec -it gbt-postgres dropdb -U postgres gbt_template

db-reset: db-drop db-create migrate-up ## Reset database

# Linting and formatting
fmt: ## Format Go code
	go fmt ./...

lint: ## Run golangci-lint
	golangci-lint run

# Security
security: ## Run security checks
	gosec ./...
