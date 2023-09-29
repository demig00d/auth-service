![actions status](https://github.com/demig00d/auth-service/actions/workflows/ci.yml/badge.svg)

# Описание

Проект является реализацией [этого](https://github.com/demig00d/auth-service/blob/master/TASK.md) тестового задания.

# Запуск

## 1. Создать файл с настройками

> [!info]
> В задании нет упоминаний настроек для приложения, но для удобства они есть,
> поэтому в проекте присутствует сторонняя библиотека - <https://github.com/ilyakaznacheev/cleanenv>

```bash
mv env.template .env
```

## 2.1 Запустить сервис

```bash
go run cmd/app/main.go
```

## 2.2. Запустить сервис c mongo через make и docker compose

```bash
make compose-up
```

# Тесты

```bash
go test -v -cover -race ./internal/...
```

или через `make`

```bash
make test
```

# REST маршруты

* `POST /authorize`
  * принимает GUID
  * возвращает пару Access, Refresh токенов
* `POST /refresh`
  * принимает GUID и Refresh токен
  * возвращает пару Access, Refresh токенов
