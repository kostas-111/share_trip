APP_NAME        := share_trip
BIN_DIR         := $(CURDIR)/bin
MAIN_PKG        := ./cmd/sharetrip

DB_DSN          ?= postgres://postgres:password@localhost:6544/share_trip?sslmode=disable

GOOSE           ?= goose
MIGRATIONS_DIR  ?= ./migrations

GOLANGCI_LINT_VERSION := v2.11.3
MIGRATE_VERSION       := v4.18.1

MIGRATE       := $(BIN_DIR)/migrate
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint

.DEFAULT_GOAL := help
.PHONY: help deps fmt lint test build run up down migrate-up migrate-down migrate-status check

help: ## Показать список команд
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

# ==== Окружение ====
deps: $(GOLANGCI_LINT) $(MIGRATE) ## Подготовить окружение (модули + инструменты)
	go mod download

$(GOLANGCI_LINT):
	@mkdir -p $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

$(MIGRATE):
	@mkdir -p $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)

# ==== Качество кода ====
fmt: ## Форматировать код
	go fmt ./...

lint: $(GOLANGCI_LINT) ## Запустить golangci-lint
	$(GOLANGCI_LINT) run --timeout=5m ./...

test: ## Запустить тесты
	go test -race -cover ./...

# ==== Сборка и запуск ====
build: ## Собрать бинарник
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN_PKG)

run: ## Запустить приложение
	go run $(MAIN_PKG)

# ==== Инфраструктура ====
up: ## Поднять PostgreSQL в Docker
	docker compose up -d

down: ## Остановить контейнеры
	docker compose down

migrate-up: ## Применить все миграции
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

migrate-down: ## Откатить последнюю миграцию
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

migrate-status: ## Показать статус всех миграций
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" status

migrate-version: ## Показать текущую версию схемы
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" version

migrate-reset: ## Откатить ВСЕ миграции (down-to 0)
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down-to 0

migrate-create: ## Создать миграцию: make migrate-create name=create_trips
	$(GOOSE) -dir $(MIGRATIONS_DIR) create $(name) sql

# ==== Полный прогон (как в CI) ====
check: fmt lint test ## fmt + lint + test