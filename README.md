# ms-gofiber — Go Fiber Boilerplate

A production-ready **Go Fiber** boilerplate focused on **clean architecture + DDD**, **first-class observability** (
Elastic **APM** end-to-end), structured **logging** with **welog**, **PostgreSQL (pgx/v5)**, **Redis (go-redis/v9)**,
and a **validator** layer with both field-level and struct-level custom rules.
It also includes a **mandatory self-hit** outbound client example that logs with welog and traces with APM.

---

## Key Features

* **Fiber v2** web server with sane defaults.
* **Observability by default**

    * **Elastic APM**: inbound (Fiber), outbound HTTP, Postgres (pgx/v5), and Redis (via custom hook) — everything
      carrying a `context.Context` is traced.
    * **welog**: request/response logging middleware + per-request logger in `c.Locals("logger")`, and **client logs**
      via `welog.LogFiberClient(...)`.
* **Postgres (pgx/v5)** pool with APM instrumentation.
* **Redis (go-redis/v9)** client with APM hook (`ProcessHook`, `ProcessPipelineHook`, `DialHook`).
* **Validation** with go-playground/validator:

    * **Plain-base** rules (field validations),
    * **Struct-base** rules (cross-field/semantic),
    * **Rule registration model** mirrors your provided pattern.
* **Response remapping via map** (no `switch`) for consistent API envelopes.
* **Self-hit endpoint** demonstrating outbound HTTP with APM + welog client logging.
* **Clean architecture** layering inspired by Fiber recipe: `domain` (entities + ports) → `usecase` (application
  service) → `transport/http` (handler/routes/presenter), with infra in `adapter`.

---

## Tech Stack

* Go `1.22+`
* Fiber `v2`
* Elastic APM modules (`apmfiber`, `apmhttp`, `apmpgxv5`)
* pgx/v5 (with pgxpool)
* go-redis/v9
* welog
* validator/v10
* logrus
* `github.com/joho/godotenv` for `.env` loading

---

## Folder Structure

```
ms-gofiber/
├─ go.mod
├─ go.sum
├─ .env.example
├─ sql/
│  └─ schema.sql
├─ cmd/
│  └─ ms-gofiber/
│     └─ main.go
├─ internal/
│  ├─ app/                 # build Fiber app (middlewares, DI, router)
│  ├─ config/              # config loader (.env via godotenv)
│  ├─ dto/                 # request DTOs validated with validator
│  ├─ domain/
│  │  └─ todo/             # entity + repository port + domain errors
│  ├─ usecase/
│  │  └─ todo/             # application service/use case
│  ├─ adapter/
│  │  ├─ cache/
│  │  │  └─ redis/         # redis cache adapter
│  │  └─ repository/
│  │     └─ postgres/      # pgx/v5 repository implementation
│  ├─ middleware/          # global middlewares (error, headers, request id)
│  ├─ transport/
│  │  └─ http/             # router + handlers + presenter + routes
│  └─ validator/           # custom rules: plainbase + structbase + register
├─ pkg/
│  ├─ apmredis/            # Redis APM hook (trace all commands)
│  ├─ apperror/            # typed app errors
│  ├─ cache/               # Redis client construction
│  ├─ db/                  # Postgres pool construction (APM'd)
│  ├─ httpx/               # APM-wrapped HTTP client + welog client logs
│  └─ respond/             # response envelope + code→HTTP map
```

---

## Getting Started

### Prerequisites

* Go `1.22+`
* PostgreSQL & Redis
* (Optional) Elastic APM Server if you want tracing visualized

### 1) Configure environment

Copy the sample and edit to your needs:

```bash
cp .env.example .env
```

Important keys (see `.env.example`):

* **Server**: `APP_HOST`, `APP_PORT`, `APP_READ_TIMEOUT_SEC`, `APP_WRITE_TIMEOUT_SEC`
* **Postgres**: `PG_URL`, `PG_MAX_CONNS`, `PG_MIN_CONNS`
* **Redis**: `REDIS_ADDR`, `REDIS_DB`, `REDIS_PASSWORD`, `REDIS_DEFAULT_TTL_SEC`
* **Elastic APM**: `ELASTIC_APM_SERVER_URL`, `ELASTIC_APM_SERVICE_NAME`, `ELASTIC_APM_ENVIRONMENT`, etc.

> Note: `godotenv` loads `.env` automatically if present.

### 2) Install deps

```bash
go mod tidy
```

### 3) Prepare database

```bash
psql "$PG_URL" -f sql/schema.sql
```

### 4) Run the service

```bash
go run ./cmd/ms-gofiber
```

Server listens on `APP_HOST:APP_PORT` (default `0.0.0.0:8080`).

---

## API Reference (Quick)

### Required Headers (for protected endpoints)

Middleware `HeaderGuard` + `ExternalIDGuard` mewajibkan header ini pada endpoint non-skip:

```
X-PARTNER-ID: <alphanumeric, max 36>
CHANNEL-ID: <alphanumeric, max 5>
X-EXTERNAL-ID: <numeric, max 36, unique within TTL>
```

Skip path default:

```
/v1/health
/v1/internal/echo
/v1/client/self-call
```

### Health

```
GET /v1/health
```

→ `{"status":"ok"}`

### Todos

```
POST   /v1/todos           # create
GET    /v1/todos           # list (limit, offset)
GET    /v1/todos/:id       # get by id
PUT    /v1/todos/:id       # update
DELETE /v1/todos/:id       # delete
```

**Sample: Create**

```bash
curl -X POST http://localhost:8080/v1/todos \
  -H "Content-Type: application/json" \
  -H "X-PARTNER-ID: PARTNER123" \
  -H "CHANNEL-ID: CHN01" \
  -H "X-EXTERNAL-ID: 10000000000001" \
  -d '{"title":"My Task","completed":false}'
```

**Response (success envelope)**

```json
{
  "code": "OK",
  "message": "success",
  "data": {
    "id": "d2c8...uuid...",
    "title": "My Task",
    "completed": false,
    "created_at": "2025-11-03T02:15:00Z",
    "updated_at": "2025-11-03T02:15:00Z"
  }
}
```

### Internal echo (service local probe)

```
GET /v1/internal/echo?msg=hello
```

### Self-Hit (mandatory outbound example)

```
GET /v1/client/self-call
```

* Calls `GET /v1/internal/echo` via the APM-wrapped HTTP client.
* Logs client request/response using `welog.LogFiberClient(...)`.
* Returns upstream status and body in the standard success envelope.

### Validator Demo (plain + struct base)

```
POST /v1/internal/validator/prepare-example
```

Body example:

```json
{
  "terminalType": "APP",
  "osType": "ANDROID",
  "osVersion": "14",
  "grantType": "AUTHORIZATION_CODE",
  "paymentMethodType": "DANA",
  "scope": ["SEND_OTP"],
  "transactionTime": "2026-02-23T10:30:00Z",
  "merchantName": "Demo Merchant"
}
```

---

## Response & Error Model

All responses use a consistent envelope (see `pkg/respond/respond.go`):

```json
{
  "code": "OK | BAD_REQUEST | VALIDATION | NOT_FOUND | DB_ERROR | INTERNAL | ...",
  "message": "human friendly message",
  "data": {
    ...
  },
  // optional
  "meta": {
    ...
  },
  // optional
  "fields": {
    "Field": "Tag"
  }
  // on validation errors
}
```

**HTTP status mapping** is driven by a **map** (no switch) to keep concerns separated and predictable.

---

## Validation

* **Field-level** rules in `internal/validator/rule/plainbase` (e.g. `alphanum_with_space`, `authorization_scope`, etc).
* **Struct-level** rules in `internal/validator/rule/structbase`.

Example DTO (validated in handlers):

```go
type TodoUpsertRequest struct {
Title     string `json:"title" validate:"required,min=3,max=100,alphanum_with_space"`
Completed bool   `json:"completed"`
}
```

Struct-level rule examples:

* `TodoUpsertStructRule`: enforce trim + non-blank title.
* `PrepareExampleStructRule`: validasi kombinasi `terminalType` vs `osType`/`osVersion` (adaptasi dari pola
  `PrepareStructRule` di project referensi).

**Rule registration** follows a map-based model:

```go
var customRules = map[string]validator.Func{
// field rules...
}
var customStructRules = []fiber.Map{
{ "func": structbase.TodoUpsertStructRule, "type": dto.TodoUpsertRequest{} },
{ "func": structbase.PrepareExampleStructRule, "type": dto.PrepareExampleRequest{} },
}
```

---

## Logging

* **welog** middleware is mounted globally:

  ```go
  app.Use(welog.NewFiber(fiber.Config{}))
  ```
* Per-request logger is available on the context:

  ```go
  c.Locals("logger").(*logrus.Entry).Error("Something went wrong")
  ```
* **Client logs** for outbound calls use:

  ```go
  welog.LogFiberClient(c, reqModel, resModel)
  ```

  (used inside `pkg/httpx/client.go` and self-hit handler)

---

## Observability & APM

* **Inbound**: `apmfiber.Middleware()` creates a transaction for each HTTP request (and recovers panics).
* **Outbound HTTP**: `apmhttp.WrapClient(...)` automatically creates spans for downstream calls (context propagated from
  `c.UserContext()`).
* **Postgres**: `apmpgxv5.Instrument(...)` instruments pgx/v5; queries executed with the request context are traced.
* **Redis**: Custom **hook** creates spans for `ProcessHook`, `ProcessPipelineHook`, and `DialHook`.

> Ensure your APM environment variables are set (see `.env.example`).
> By design, **every operation that carries `context.Context`** in handlers/services/repos/cache/client is either traced
> automatically or wrapped in a custom APM span.

---

## Clean Architecture / DDD

* **Domain** (`internal/domain`): entities + repository ports + domain-level errors, tanpa ketergantungan framework.
* **Usecase** (`internal/usecase`): business flow aplikasi, orkestrasi repo + cache melalui interface.
* **Adapters** (`internal/adapter/...`): implementasi konkret DB/cache untuk memenuhi port/usecase.
* **Transport** (`internal/transport/http`): handlers + routes + presenter, untuk parsing request dan formatting
  response.
* **Cross-cutting** (`pkg/...`): DB/Redis constructors, APM hooks, HTTP client, response mapping, app errors.

---

## Configuration

Loaded via **godotenv** (falls back to OS env if `.env` is missing).
Key values (see `.env.example`):

```env
APP_NAME=ms-gofiber
APP_HOST=0.0.0.0
APP_PORT=8080
APP_READ_TIMEOUT_SEC=10
APP_WRITE_TIMEOUT_SEC=10

PG_URL=postgres://postgres:postgres@localhost:5432/ms_gofiber?sslmode=disable
PG_MAX_CONNS=10
PG_MIN_CONNS=2

REDIS_ADDR=localhost:6379
REDIS_DB=0
REDIS_PASSWORD=
REDIS_DEFAULT_TTL_SEC=60

# Elastic APM (optional but recommended)
ELASTIC_APM_SERVER_URL=http://localhost:8200
ELASTIC_APM_SERVICE_NAME=ms-gofiber
ELASTIC_APM_ENVIRONMENT=local
ELASTIC_APM_RECORDING=true
```

---

## Extending

* **New domain**: create `internal/domain/<name>` (entity + repository port), lanjutkan `internal/usecase/<name>`, lalu
  implement adapter(s) di `internal/adapter` dan expose lewat handlers + routes.
* **New validators**: implement in `internal/validator/rule/plainbase` or `structbase`, register in
  `internal/validator/rule/register.go`.
* **New outbound client**: add helper in `pkg/httpx` or a domain-specific gateway; always pass `c.UserContext()` and log
  with `welog.LogFiberClient(...)`.

---

## Troubleshooting

* **APM shows no traces**: verify APM env vars, server URL, and that requests hit the service. Ensure outbound calls use
  `c.UserContext()`.
* **No welog client logs**: confirm `welog.NewFiber(...)` is mounted and your outbound path calls
  `welog.LogFiberClient(...)` (already wired in `pkg/httpx` and self-hit).
* **DB errors**: ensure the schema is loaded (`sql/schema.sql`) and `PG_URL` is reachable.
* **Redis errors**: verify `REDIS_ADDR` and permissions.
* **Validation errors**: see `fields` object in the error envelope for tag names.

---

## License

MIT License — see [LICENSE](LICENSE) for details.

---

## Quick Commands

```bash
# install deps
go mod tidy

# run locally
go run ./cmd/ms-gofiber

# hit endpoints
curl http://localhost:8080/v1/health
curl -X POST http://localhost:8080/v1/todos \
  -H 'Content-Type: application/json' \
  -H 'X-PARTNER-ID: PARTNER123' \
  -H 'CHANNEL-ID: CHN01' \
  -H 'X-EXTERNAL-ID: 10000000000001' \
  -d '{"title":"Demo Todo","completed":false}'
curl http://localhost:8080/v1/client/self-call
```

Happy building!
