![actions status](https://github.com/demig00d/auth-service/actions/workflows/ci.yml/badge.svg)

## Описание

Проект является реализацией [этого](https://github.com/demig00d/auth-service/blob/master/TASK.md) тествого задания.

## Запуск

```bash
mv .env.template .env
```

```bash
make compose-up
```

## Тесты

```bash
make test
```

## Конфигурация

В задании нет упоминаний настроек для приложения, но для удобства они есть,
поэтому в проекте присутствует сторонняя библиотека - <https://github.com/ilyakaznacheev/cleanenv>

## REST маршруты

* `POST /authorize`
  * принимает GUID
  * возвращает пару Access, Refresh токенов
* `POST /refresh`
  * принимает Refresh токен
  * возвращает пару Access, Refresh токенов
