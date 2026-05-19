# Architecture Rules

## Layer Responsibilities

* `cmd/app`: process bootstrap and app-level models.
* `pkg/server`: dependency wiring, Fiber config, middleware order, and route registration.
* `router`: route group construction.
* `handler`: HTTP input/output translation.
* `external/domain/<feature>`: outbound adapter repositories and services.
* `internal/domain/<feature>/model`: feature entities and DTOs.
* `internal/domain/<feature>/repository`: repository ports and persistence implementation.
* `internal/domain/<feature>/service`: usecase orchestration.
* `pkg/*`: reusable cross-cutting helpers only.

## Dependency Direction

* `cmd/app` calls `pkg/server`.
* `pkg/server` wires repositories, services, middleware, and routers.
* `router` calls handler factory functions with `*cmd/app/model.Service`.
* `handler` parses Fiber input, validates DTOs, calls services, and shapes responses.
* `service` coordinates repositories and feature behavior.
* `model` stays free of Fiber, persistence, logging, APM, and configuration.

## Current Composition

The current service uses:

* SQLite-backed todo persistence under `pkg/infrastructure/database`.
* Generic in-memory TTL cache under `pkg/infrastructure/cache`.
* Generic outbound HTTP client under `pkg/client`.
* External echo adapter under `external/domain/echo`.
* Cache flush repository/service under `internal/domain/cache`.
* External ID duplicate guard under `internal/domain/externalid`.
* Neutral response-code remapping under `internal/domain/remapping`.
* Todo repository and service/usecase under `internal/domain/todo`.
* Generic request validation through `internal/domain/reqvalidator/service`.
* Plain-base and struct-base validator rules under `pkg/rule`.
* Generic response envelope and neutral response-code model through `pkg/response` and `pkg/responsecode`.

## Pattern Policy

Use interfaces at real boundaries only:

* Repository interfaces belong at service/repository boundaries.
* Handler factory functions receive `*cmd/app/model.Service`.
* Middleware owns request-wide concerns such as header validation.
* DTO mapping stays outside usecases.

Avoid extra factories, builders, service locators, or internal-code mappings unless there is a concrete variation point.
