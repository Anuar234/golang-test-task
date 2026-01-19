# Go Numbers Service

Simple REST microservice that accepts a number, stores it in PostgreSQL, and returns the sorted list of all numbers.

## Run

```bash
docker compose up --build
```

Service listens on `http://localhost:8080`.

## API

`POST /numbers`

Request:

```json
{"number": 3}
```

Response:

```json
{"numbers": [1,2,3]}
```

`GET /healthz` returns 200 when the service is up.

## Local tests

```bash
go test ./...
```
