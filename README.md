# ms-gofiber

Clean Go Fiber service using the preferred app composition pattern:

* `cmd/app`: process entrypoint and app-level models.
* `pkg/server`: dependency wiring and Fiber middleware order.
* `router`: route group registration.
* `handler`: Fiber handlers and middleware.
* `internal/domain/todo`: entity, repository port, repository implementation, and service/usecase.
* `pkg/response`: generic response envelope without business-specific codes.

## Run

```bash
make run
```

Default port comes from `APP_PORT`. If it is empty, Fiber uses the address passed by the caller.

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
    "Title": "required"
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
