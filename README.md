# ms-gofiber

Clean Go Fiber service using the preferred app composition pattern:

* `cmd/app`: process entrypoint and app-level models.
* `pkg/server`: dependency wiring and Fiber middleware order.
* `router`: route group registration.
* `handler`: Fiber handlers and middleware.
* `external/domain/echo`: generic outbound client example with repository/service split.
* `internal/domain/todo`: entity, repository port, repository implementation, and service/usecase.
* `pkg/infrastructure`: generic SQLite database and in-memory cache helpers.
* `pkg/rule`: generic plain-base and struct-base validator rules.
* `pkg/response`: generic response envelope without business-specific codes.

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

Todo endpoints require a generic client header:

```bash
X-CLIENT-ID: demo-client
```

Create todo:

```bash
curl -X POST http://localhost:8080/v1/todos \
  -H 'Content-Type: application/json' \
  -H 'X-CLIENT-ID: demo-client' \
  -d '{"title":"write tests","completed":false}'
```

List todos:

```bash
curl http://localhost:8080/v1/todos \
  -H 'X-CLIENT-ID: demo-client'
```

External echo example:

```bash
curl 'http://localhost:8080/v1/external/echo?target=https://example.com' \
  -H 'X-CLIENT-ID: demo-client'
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
