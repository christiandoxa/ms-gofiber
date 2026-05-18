# API Layer

This directory owns Fiber-facing code and API-facing documentation.

Package layout:

* `dto`: request payload types.
* `handler`: Fiber handlers/controllers.
* `middleware`: Fiber middleware and API error handling.
* `presenter`: response data mapping.
* `respond`: response envelope and HTTP status mapping.
* `router`: route registration.
* `validation`: request struct-level validation rules.

Current service behavior is summarized in the API section of the root `README.md`.

When adding API documentation:

* Keep examples aligned with the response envelope from `api/respond`.
* Document required headers for protected endpoints.
* Keep internal-only endpoints clearly labeled.
* Prefer examples that can be executed with `curl`.
