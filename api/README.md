# API Notes

This directory is reserved for API-facing documentation such as OpenAPI files, request examples, and endpoint notes.

Current service behavior is summarized in the API section of the root `README.md`.

When adding API documentation:

* Keep examples aligned with the response envelope from `pkg/respond`.
* Document required headers for protected endpoints.
* Keep internal-only endpoints clearly labeled.
* Prefer examples that can be executed with `curl`.
