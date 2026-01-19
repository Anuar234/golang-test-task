# Go Numbers Service

Мини‑микросервис на Go: принимает число по REST, пишет его в Postgres и возвращает отсортированный список всех чисел.

## Что внутри

- Go + net/http без лишней магии.
- Postgres для хранения чисел.
- Docker + docker compose, запуск в одну кнопку.
- Тесты на хендлер.

## Быстрый старт

```bash
docker compose up --build
```

Сервис слушает `http://localhost:8080`.

## API

### POST /numbers

Request:

```json
{"number": 3}
```

Response:

```json
{"numbers": [1,2,3]}
```

### GET /healthz

Возвращает `200` когда сервис жив.

## Примеры вызова

### PowerShell

```powershell
Invoke-WebRequest -Method Post -Uri http://localhost:8080/numbers `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{"number":3}'
```

### curl (Windows)

```bash
curl.exe -X POST http://localhost:8080/numbers -H "Content-Type: application/json" -d "{\"number\":3}"
```

## Конфиг (env)

- `HTTP_ADDR` — адрес/порт сервера (по умолчанию `:8080`).
- `DB_HOST` — хост Postgres (по умолчанию `db`).
- `DB_PORT` — порт Postgres (по умолчанию `5432`).
- `DB_USER` — пользователь (по умолчанию `app`).
- `DB_PASSWORD` — пароль (по умолчанию `app`).
- `DB_NAME` — база (по умолчанию `numbers`).
- `DB_SSLMODE` — sslmode (по умолчанию `disable`).

## Локальные тесты

```bash
go test ./...
```

## Структура проекта

- `cmd/server/main.go` — запуск сервиса.
- `internal/config` — конфиг из env.
- `internal/httpapi` — http хендлеры.
- `internal/storage` — работа с Postgres.
- `docker-compose.yml` — база + апп.

## Примечания

- Сообщение `locale: not found` в логах Postgres на Alpine — это ок, на работу не влияет.
