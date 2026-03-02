# =============================================================================
# Переменные
# =============================================================================
OUTPUT := ./bin/app
GO_LINT_VERSION := 2.7.2
GO_FILE := ./main.go
<<<<<<< HEAD
GO_EXE := ./app
=======
>>>>>>> 4dd41921cf8b5341806fb2b0f67c5ff30d264719

# =============================================================================
# Справка
# =============================================================================
.PHONY: help
help: ## Показать справку
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# =============================================================================
# Разработка
# =============================================================================
.PHONY: run
<<<<<<< HEAD
run: ## Сборка и запуск приложения
	go build -o ${OUTPUT} ${GO_FILE} && go run ${GO_FILE}
=======
run: ## Запустить приложение
	go run ${GO_FILE}

.PHONY: build
build: ## Сборка приложения
	go build -o ${OUTPUT} ${GO_FILE}
>>>>>>> 4dd41921cf8b5341806fb2b0f67c5ff30d264719

.PHONY: test
test: ## Запуск тестов
	go test -count=1 -v ./...

<<<<<<< HEAD
.PHONY: clear
clear: ## Удаление бинарного файла
	rm -f ${OUTPUT} ${GO_EXE}

=======
>>>>>>> 4dd41921cf8b5341806fb2b0f67c5ff30d264719
# =============================================================================
# Качество кода
# =============================================================================
.PHONY: lint
lint: ## Запуск линтера
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v${GO_LINT_VERSION} run

.PHONY: lint-fix
lint-fix: ## Запуск линтера с автофиксом
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v${GO_LINT_VERSION} run --fix

# =============================================================================
# Окружение (Docker)
# =============================================================================
.PHONY: up
up: ## Поднять docker окружение (PostgreSQL)
	docker compose up -d

.PHONY: down
down: ## Остановить docker окружение
	docker compose down --remove-orphans

.PHONY: logs
logs: ## Показать логи docker контейнеров
	docker compose logs -f

# =============================================================================
# Зависимости
# =============================================================================
.PHONY: deps
deps: ## Загрузить зависимости
	go mod tidy
	go mod download

.PHONY: mod-check
mod-check: ## Проверка актуальности go.mod/go.sum
	go mod tidy
	@FILES="go.mod"; [ -f go.sum ] && FILES="$$FILES go.sum"; git diff --exit-code -- $$FILES || (echo "go.mod/go.sum не синхронизированы. Запустите 'go mod tidy'" && exit 1)

# =============================================================================
# CI
# =============================================================================
.PHONY: ci
ci: ## Запустить все CI проверки
	@echo "=== Mod Check ==="
	go mod tidy
	@FILES="go.mod"; [ -f go.sum ] && FILES="$$FILES go.sum"; git diff --exit-code -- $$FILES || (echo "go.mod/go.sum не синхронизированы" && exit 1)
	@echo ""
	@echo "=== Build ==="
	@mkdir -p ./bin
<<<<<<< HEAD
	go build -o ./bin/app && go run ./main.go
=======
	go build -o ./bin/ -v ./...
>>>>>>> 4dd41921cf8b5341806fb2b0f67c5ff30d264719
	@echo ""
	@echo "=== Test ==="
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@echo ""
	@echo "=== Lint ==="
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v${GO_LINT_VERSION} run --timeout=10m
	@echo ""
	@echo "CI passed!"
