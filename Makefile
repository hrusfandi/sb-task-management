include .env
export

.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: db-setup
db-setup: ## Create database if not exists
	@./scripts/setup_db.sh

.PHONY: run
run: ## Run the application
	go run main.go

.PHONY: build
build: ## Build the application
	go build -o bin/app main.go

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: migrate-up
migrate-up: ## Run database migrations
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

.PHONY: migrate-down
migrate-down: ## Rollback database migrations
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

.PHONY: migrate-create
migrate-create: ## Create a new migration file (usage: make migrate-create name=migration_name)
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: install-migrate
install-migrate: ## Install golang-migrate CLI tool
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: deps
deps: ## Install dependencies
	go mod download
	go mod tidy

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/