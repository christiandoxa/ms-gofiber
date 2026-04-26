# Clean Architecture Rules

## Layer Responsibilities

This service follows a clean architecture shape:

* `internal/app/domain`: domain entities and repository interfaces.
* `internal/app/application/usecase`: application flows that coordinate domain logic through interfaces.
* `internal/app/adapter/controller`: Fiber handlers and route wiring.
* `internal/app/adapter/dto`: request DTOs for transport validation.
* `internal/app/adapter/presenter`: response mapping from domain/application data to API payloads.
* `internal/app/adapter/repository`: concrete persistence and cache implementations.
* `internal/app/adapter/validation`: request-specific struct-level validation rules.
* `pkg/`: reusable infrastructure helpers such as database setup, Redis setup, HTTP client wrapping, APM hooks, app errors, and response envelopes.

## Dependency Direction

Dependencies must point inward:

* Controllers may call usecases and presenters.
* Usecases may call domain interfaces.
* Repository implementations may satisfy domain interfaces.
* Domain code must not import adapters, frameworks, logging, APM, SQL, Redis, or configuration.
* Usecases must not import Fiber, SQL drivers, Redis clients, HTTP clients, or presenter packages.

## Boundary Rules

Keep business decisions close to the domain or usecase that owns them:

* Request parsing and header validation belong in controllers or middleware.
* Domain invariants belong in domain entities or usecases.
* Storage details belong in repository implementations.
* Response status and envelope mapping belong in presenter/respond layers.
* Observability, logging, and transport concerns must not leak into domain entities.
* Best-effort infrastructure failures must be explicit and observable, even when they do not fail the primary usecase flow.

## Design Pattern Usage Policy

Patterns are allowed when they reduce coupling or clarify ownership:

* Use repository interfaces to isolate usecases from SQLite, Redis, or future persistence choices.
* Use presenters/mappers to keep response shape out of usecases.
* Use middleware for request-wide cross-cutting behavior.
* Use map-based registration for stable mappings such as response codes, validators, or handler tables.
* Use small constructors for dependency wiring when they make tests simpler.

Avoid pattern cargo culting:

* Do not add abstract factories, builders, observers, or service locators without a concrete variation point.
* Do not create interfaces for a single implementation unless needed for a boundary, test seam, or external dependency.
* Do not place business rules in generic `pkg/` helpers.
* Do not let DTOs become domain entities.

## Change Checklist

Before adding or moving code, verify:

* The package belongs to the correct layer.
* Imports do not point from inner layers to outer layers.
* Business rules are testable without Fiber, SQL, Redis, or external services.
* New adapters implement existing ports instead of changing usecases to know adapter details.
* New cross-cutting helpers are generic and do not encode feature-specific policy.
