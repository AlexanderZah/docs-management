# Docs Management Service

## Описание

Сервис для хранения и управления документами с поддержкой прав доступа и кешированием в Redis.  
Используется PostgreSQL для хранения данных и Redis для кеша. 

## Как запустить

1. Клонируйте репозиторий:

```bash
git clone https://github.com/AlexanderZah/docs-management.git
cd docs-management
```

2. Запустите сервисы с помощью Docker Compose:
```bash
docker-compose -f deployments/docker-compose.yaml up --build
```

3. Приложение будет доступно по адресу:
```
http://localhost:8082
```
4. Чтобы остановить контейнеры, выполните:
```
docker-compose -f deployments/docker-compose.yaml down
```