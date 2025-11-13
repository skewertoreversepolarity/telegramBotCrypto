# Makefile для Crypto Bot

.PHONY: build run test clean deps

# Сборка проекта
build:
	go build -o crypto-bot main.go

# Запуск проекта
run:
	go run main.go

# Установка зависимостей
deps:
	go mod tidy
	go mod download

# Тестирование (если будут тесты)
test:
	go test ./...

# Тестирование Telegram бота
test-bot:
	go run cmd/test/test_telegram.go

# Проверка конфигурации
check-config:
	go run cmd/test/check_config.go

# Очистка скомпилированных файлов
clean:
	rm -f crypto-bot

# Создание .env файла из примера
setup:
	cp env.example .env
	@echo "Не забудьте отредактировать .env файл с вашими настройками!"

# Проверка кода
lint:
	golangci-lint run

# Форматирование кода
fmt:
	go fmt ./...

# Показать помощь
help:
	@echo "Доступные команды:"
	@echo "  build   - Собрать проект"
	@echo "  run     - Запустить проект"
	@echo "  deps    - Установить зависимости"
	@echo "  test    - Запустить тесты"
	@echo "  clean   - Очистить скомпилированные файлы"
	@echo "  setup   - Создать .env файл из примера"
	@echo "  lint    - Проверить код"
	@echo "  fmt     - Отформатировать код"
	@echo "  help    - Показать эту справку"
