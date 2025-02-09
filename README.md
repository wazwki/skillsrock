## Использование приложения

#### — Запускает линтер  
```sh
make lint
```

#### — Запускает тесты  
```sh
make test
```

#### — Собирает Docker-образы  
```sh
make build
```

#### — Запускает приложение в фоне  
```sh
make up
```

#### — Выполняет все шаги (lint, test, build, up)  
```sh
make run
```

#### — Выполняет все шаги и запускает приложение без генерации и проверки jwt-токенов  
```sh
make run DEBUG=true
```


#### Если golangci-lint или docker-compose не установлены, их нужно поставить.  

#### http://localhost:8080/swagger/ - адрес swagger, для использования запускать приложение в debug режиме  

---
### Регистрация  
```sh
curl -X 'POST' \
  'http://localhost:8080/api/v1/auth/register' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "test",
  "password": "test"
}'
```

### Логин  
```sh
curl -X 'POST' \
  'http://localhost:8080/api/v1/auth/login' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "test",
  "password": "test"
}'
```

### Создание задачи  
```sh
curl -X 'POST' \
  'http://localhost:8080/api/v1/tasks' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "description": "test1",
  "due_date": "2025-12-12 15:04:05",
  "priority": "low",
  "status": "pending",
  "title": "test1"
}'
```

### Получение задач
```sh
curl -X 'GET' \
  'http://localhost:8080/api/v1/tasks' \
  -H 'accept: application/json'
```

#### Поиск по названию  
```sh
curl -X 'GET' \
  'http://localhost:8080/api/v1/tasks?name=test1' \
  -H 'accept: application/json'
```

#### Фильтрация  
```sh
curl -X 'GET' \
  'http://localhost:8080/api/v1/tasks?status=pending&sort_by=low&priority=low' \
  -H 'accept: application/json'
```

### Обновление задачи
```sh
curl -X 'PUT' \
  'http://localhost:8080/api/v1/tasks/1' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "description": "test1 - updated",
  "due_date": "2025-12-12 15:04:05",
  "priority": "low",
  "status": "pending",
  "title": "test1"
}'
```

### Удаление задачи  
```sh
curl -X 'DELETE' \
  'http://localhost:8080/api/v1/tasks/1' \
  -H 'accept: application/json'
```

### Экспорт задач
```sh
curl -X 'GET' \
  'http://localhost:8080/api/v1/tasks/export' \
  -H 'accept: application/json'
```

### Импорт задач
```sh
curl -X 'POST' \
  'http://localhost:8080/api/v1/tasks/import' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '[
{
  "description": "test_imported",
  "due_date": "2025-12-12 15:04:05",
  "priority": "low",
  "status": "pending",
  "title": "test_imported"
}
]'
```

### Получение аналитики
```sh
curl -X 'GET' \
  'http://localhost:8080/api/v1/analytics' \
  -H 'accept: application/json'
```