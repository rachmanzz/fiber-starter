.PHONY: help build run dev test test-v clean spark-init spark-migrate docker-up docker-down

APP_NAME = server
SPARK_BIN = spark

help: ## Display available commands
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the main app and spark CLI binaries
	@echo "Building binaries..."
	go build -ldflags="-w -s" -o $(APP_NAME) ./cmd/server/main.go
	go build -ldflags="-w -s" -o $(SPARK_BIN) ./spark-cli/main.go

run: build ## Build and run the app server
	./$(APP_NAME)

dev: ## Run development server with live reload (using Spark CLI)
	go run ./spark-cli dev

test: ## Run unit tests
	go test ./...

test-v: ## Run unit tests with verbose output
	go test -v ./...

clean: ## Clean built binaries and log files
	rm -f $(APP_NAME) $(SPARK_BIN)
	rm -rf logs/*.log

spark-init: build ## Initialize project module via Spark CLI
	./$(SPARK_BIN) init

spark-migrate: build ## Run database migrations via Spark CLI
	./$(SPARK_BIN) migrate

docker-up: ## Start application and PostgreSQL containers via docker-compose
	docker compose up -d --build

docker-down: ## Stop docker containers
	docker compose down
