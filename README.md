# Distributed Rate Limiter

A distributed API rate limiter built in Go. It exposes a rate-limit decision
endpoint backed by Redis, manages per-client rate-limit policies in
PostgreSQL, and supports three interchangeable limiting algorithms:
**fixed window**, **sliding window**, and **token bucket**. Prometheus
metrics are collected for every request via a global middleware.

## Features

- **Pluggable rate-limiting algorithms** — fixed window, sliding window, and
  token bucket, each implemented as an atomic Redis Lua script (see
  `internal/limiter/lua/`) to avoid race conditions under concurrent load.
- **Per-client configuration** — every client has its own algorithm, limit,
  and window, stored in PostgreSQL and managed through a CRUD API.
- **API key authentication** — the rate-limit check endpoint is protected by
  an `X-API-Key` header, resolved to a client via middleware.
- **Prometheus instrumentation** — request count and latency are recorded for
  every route through a single global middleware and exposed at `/metrics`.
- **Health checks** — `/health` verifies connectivity to both PostgreSQL and
  Redis.
- **Dockerized local environment** — Postgres, Redis, and the API server run
  together via `docker-compose`, with live reload powered by
  [Air](https://github.com/air-verse/air).

## Tech Stack

| Layer            | Technology                          |
|-------------------|--------------------------------------|
| Language           | Go 1.26                              |
| HTTP framework     | [Gin](https://github.com/gin-gonic/gin) |
| Database           | PostgreSQL (via [GORM](https://gorm.io) / pgx) |
| Cache / limiter store | Redis                              |
| Config             | [Viper](https://github.com/spf13/viper) (`.env`) |
| Metrics            | [Prometheus client](https://github.com/prometheus/client_golang) |
| Live reload (dev)  | [Air](https://github.com/air-verse/air) |

## Architecture

```
cmd/server            entrypoint, graceful shutdown
internal/
  bootstrap            wires config, DB, Redis, routes into an App
  config               env-based configuration (Viper)
  database             GORM/Postgres connection
  redis                Redis connection + Lua script runner
  routes               route registration
  middleware           API key auth, Prometheus instrumentation
  client               client CRUD (policy management)
  check                rate-limit check endpoint
  limiter              rate-limiting algorithms (fixed window, sliding
                       window, token bucket) + embedded Lua scripts
  health               liveness/readiness endpoint
  metrics              Prometheus metric definitions
  httpx                shared JSON response envelope
initdb.d               Postgres schema bootstrap SQL
```

A **client** represents an API consumer with its own rate-limit policy
(algorithm, limit, window, API key). The **check** endpoint looks up the
caller's client via their API key, runs the configured algorithm's Lua script
against Redis, and returns whether the request is allowed.

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/) and Docker Compose
- Go 1.26+ (only needed for running outside Docker, or running tests)

### Run with Docker Compose (recommended)

```bash
cp .env.example .env
# fill in POSTGRES_PASSWORD and REDIS_PASSWORD in .env
docker-compose up --build
```

This starts Postgres, Redis, and the API server (with live reload) on
`http://localhost:8080`.

### Run locally

```bash
cp .env.example .env
# point POSTGRES_HOST / REDIS_ADDR at your local instances
go mod download
go run ./cmd/server
```

### Environment variables

| Variable            | Description                          | Required |
|----------------------|----------------------------------------|----------|
| `POSTGRES_HOST`       | Postgres host                          | yes      |
| `POSTGRES_PORT`       | Postgres port                          | yes      |
| `POSTGRES_USER`       | Postgres user                          | yes      |
| `POSTGRES_PASSWORD`   | Postgres password                      | yes      |
| `POSTGRES_DB`         | Postgres database name                 | yes      |
| `REDIS_ADDR`          | Redis address (default `localhost:6379`) | no     |
| `REDIS_PASSWORD`      | Redis password                         | no       |

## API Reference

The full API is documented as an OpenAPI 3.0 spec at
[`docs/openapi.yaml`](docs/openapi.yaml). Paste it into
[editor.swagger.io](https://editor.swagger.io/) or run
`npx @redocly/cli preview-docs docs/openapi.yaml` for an interactive view.

All responses use a shared JSON envelope:

```json
{ "success": true, "message": "...", "data": { } }
```

### Health

```
GET /health
```
Returns `200` if both Postgres and Redis are reachable, `503` otherwise.

### Client management

```
POST   /api/v1/clients        create a client
GET    /api/v1/clients        list all clients
GET    /api/v1/clients/:id    get a client by ID
PUT    /api/v1/clients/:id    update a client's policy
DELETE /api/v1/clients/:id    delete a client
```

Create/update request body:

```json
{
  "clientId": "payments-service",
  "algorithm": "fixed_window",
  "limit": 100,
  "window_seconds": 60
}
```

`algorithm` is one of `fixed_window`, `sliding_window`, `token_bucket`.
Creating a client returns a generated `api_key`, used to authenticate
`/check` requests.

### Rate-limit check

```
POST /api/v1/check
Header: X-API-Key: <client api key>
```

Request body:

```json
{ "identifier": "user-123" }
```

`identifier` scopes the limit within the client (e.g. per end-user, per IP).
Response:

```json
{
  "success": true,
  "message": "request allowed",
  "data": { "allowed": true, "remaining": 42, "reset_at": 1735689600 }
}
```

Returns `429 Too Many Requests` (with a `Retry-After` header) when the limit
is exceeded.

### Metrics

```
GET /metrics
```

Prometheus scrape endpoint. Exposes `rate_limiter_requests_total` (counter)
and `rate_limiter_request_duration_seconds` (histogram), both labeled by
`method`, `path`, and `status`.

## Testing

```bash
go test ./...
```

Unit tests cover all three rate-limiting algorithms under
`internal/test/limiter`.

## Roadmap

- [ ] Return per-client usage metrics (allowed vs. rejected) alongside
      generic HTTP metrics
- [x] Client-facing API documentation (OpenAPI/Swagger) — see [`docs/openapi.yaml`](docs/openapi.yaml)
- [ ] Distributed tracing
