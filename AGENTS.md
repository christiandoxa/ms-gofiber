# Repository Instructions

## Ownership

This repository can have multiple concurrent editors. Do not revert or overwrite changes that you did not make.

Codex-owned paths for documentation and developer tooling:

* `README.md`
* `AGENTS.md`
* `Makefile`
* `.golangci.yml`
* `.vscode/extensions.json`
* `.vscode/settings.json`
* `docs/**`
* `api/**`

Do not edit Go source files, `go.mod`, or `go.sum` unless the user explicitly expands the scope.

## Architecture Rules

Follow clean architecture boundaries:

* Domain code owns entities and repository ports. It must not import Fiber, SQL, Redis, APM, logging, configuration, or transport packages.
* Application usecases orchestrate domain behavior through interfaces. They must not depend on framework controllers or concrete storage clients.
* Adapters translate external inputs and outputs. Controllers, presenters, SQL repositories, Redis repositories, and HTTP clients belong outside the domain.
* Cross-cutting packages under `pkg/` must stay generic and reusable. Do not hide business rules there.

## Design Pattern Policy

Use patterns only when they simplify a real dependency or variation point:

* Repository pattern is required at domain/application boundaries.
* Presenter/mapper functions should keep response shaping outside usecases.
* Middleware is the correct place for request-wide concerns such as headers, request ID, security headers, error handling, logging, and tracing.
* Strategy-style maps are preferred for stable code-to-handler mappings, such as response status mapping and validator registration.
* Avoid adding factories, builders, observers, or generic abstractions until there are at least two concrete implementations or a clear test seam.

## Baseline Checks

Use the documented tooling before handoff when relevant:

```bash
make fmt
make test
make race
make coverage
make lint
```

If a tool is unavailable locally, report the missing tool and the command that could not run.

Coverage is a hard baseline gate. Keep total statement coverage at `100.0%`; use focused tests and gomonkey-based seams
for external or hard-to-trigger branches instead of weakening production code.
