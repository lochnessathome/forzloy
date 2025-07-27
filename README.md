### Запуск

## Сборка веб-приложения

docker build .

## Запуск всех контейнеров

docker-compose up

Чтобы посмотреть адрес для подключения с локальной машины, нужно выполнить: docker inspect forzloy_default И в секции Containers найти billing.


### Примеры использования

## Регистрация

curl -d '{"login":"lolk", "password":"12345678"}' -H "Content-Type: application/json" -X POST http://localhost:8080/api/auth/register

## Логин

curl -d '{"login":"lolk", "password":"12345678"}' -H "Content-Type: application/json" -X POST http://localhost:8080/api/auth/login

## Покупка отчёта

curl -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMyIsImV4cCI6MTc1ODczMzAxNCwibmJmIjoxNzUzNTQ5MDE0LCJpYXQiOjE3NTM1NDkwMTR9.LOHHbwQ6jq_NvCV18x5uV2FVY1uET8aOEsj4fI2KHI8" -X POST http://localhost:8080/api/reports/66a4b08e-1e2a-41df-b957-cfbb87b0cde8/purchase

### Для разработки

Опционально можно установить https://github.com/golang-migrate/migrate/tree/master/cmd/migrate , это консольная утилита, помогающая создавать файлы миграций.
