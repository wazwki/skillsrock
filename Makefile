.PHONY: lint test build up

lint:
	golangci-lint run
	@echo "Линтер пройден"

test:
	go test ./...
	@echo "Тесты пройдены"

build:
	docker-compose build

up:
	docker-compose up -d

run: lint test build up
	@echo "Сервис успешно запущен!"
