.PHONY: help run migrate seed test fmt build docker-up docker-down

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

run: ## Run the API server
	go run ./cmd/api

migrate: ## Run database migrations
	go run ./cmd/migrate

seed: ## Seed initial database data
	go run ./cmd/seed

test: ## Run unit tests
	go test -v ./...

fmt: ## Format Go source files
	gofmt -w .

build: ## Build all packages
	go build ./...

docker-up: ## Build and start containers
	docker compose up -d --build

docker-down: ## Stop and remove containers
	docker compose down
