make lint — Запускает golangci-lint run. Если ошибок нет, выводит сообщение.
make test — Запускает тесты go test ./....
make build — Собирает Docker-образы с docker-compose build.
make up — Запускает сервисы в фоне через docker-compose up -d.
make run — Выполняет все шаги (lint, test, build, up) и пишет, что сервис запущен.

Если golangci-lint или docker-compose не установлены, их нужно поставить.

http://localhost:8080/swagger/