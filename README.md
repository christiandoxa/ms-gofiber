# ms-gofiber

Clean Go Fiber service using the preferred app composition pattern:

* `cmd/app`: process entrypoint and app-level models.
* `pkg/server`: dependency wiring and Fiber middleware order.
* `router`: route group registration.
* `handler`: Fiber handlers and middleware.
* `external/domain/echo`: generic outbound client example with repository/service split.
* `internal/domain/cache`: cache repository/service for flush usecase.
* `internal/domain/externalid`: duplicate external ID guard using generic cache.
* `internal/domain/remapping`: generic response-code remapping repository/service.
* `internal/domain/todo`: entity, repository port, repository implementation, and service/usecase.
* `pkg/infrastructure`: generic SQLite database and in-memory cache helpers.
* `pkg/rule`: generic plain-base and struct-base validator rules.
* `pkg/response` and `pkg/responsecode`: generic response envelope and neutral response-code model.

## Run

```bash
make run
```

Default port comes from `APP_PORT`. `DATABASE_PATH` defaults to `:memory:`.

## API

Health check:

```bash
curl http://localhost:8080/v1/health
```

Non-health endpoints require generic client headers:

```bash
X-CLIENT-ID: demo-client
X-EXTERNAL-ID: request-001
```

Create todo:

```bash
curl -X POST http://localhost:8080/v1/todos \
  -H 'Content-Type: application/json' \
  -H 'X-CLIENT-ID: demo-client' \
  -H 'X-EXTERNAL-ID: request-001' \
  -d '{"title":"write tests","completed":false}'
```

List todos:

```bash
curl http://localhost:8080/v1/todos \
  -H 'X-CLIENT-ID: demo-client' \
  -H 'X-EXTERNAL-ID: request-002'
```

External echo example:

```bash
curl 'http://localhost:8080/v1/external/echo?target=https://example.com' \
  -H 'X-CLIENT-ID: demo-client' \
  -H 'X-EXTERNAL-ID: request-003'
```

Flush generic cache:

```bash
curl http://localhost:8080/v1/flush/cache
```

## Response Shape

Success:

```json
{
  "status": "success",
  "data": {}
}
```

Error:

```json
{
  "status": "error",
  "message": "validation failed",
  "fields": {
    "title": "required"
  }
}
```

## Checks

```bash
make fmt
make test
make race
make coverage
make lint
```

Coverage baseline is `100.0%`.
