# Developer Tooling

## Baseline Commands

Use these commands from the repository root:

```bash
make fmt
make test
make race
make coverage
make lint
make sonar
make run
```

Equivalent raw commands:

```bash
go fmt ./...
go test -gcflags=all=-l ./...
go test -race -gcflags=all=-l ./...
go test -gcflags=all=-l -covermode=count -coverprofile=coverage.out ./...
bash scripts/check-coverage.sh coverage.out 100.0
golangci-lint run ./...
sonar-scanner
go run ./cmd/app
```

`make tidy` runs `go mod tidy`. Use it only when dependency changes are expected, because it can modify `go.mod` and `go.sum`.

## Linting

The baseline lint configuration is `.golangci.yml`. It enables a focused starter set:

* `govet`
* `staticcheck`
* `errcheck`
* `gocognit`
* `ineffassign`
* `unused`
* `bodyclose`
* `misspell`

`errcheck` runs with `check-blank: true`, so error returns must be handled explicitly instead of assigned to `_`.

Install `golangci-lint` before running `make lint`.

## Coverage

`make coverage` is a hard gate. It fails unless total statement coverage is exactly `100.0%`.

`GO_TEST_FLAGS` defaults to `-gcflags=all=-l` so gomonkey-based tests are stable and do not require production logic changes.

## SonarQube

`sonar-project.properties` provides a headless scanner baseline. Run `make sonar` only after `sonar-scanner`,
`SONAR_HOST_URL`, and `SONAR_TOKEN` are configured for your SonarQube project.
The scanner is scoped to `cmd`, `external`, `handler`, `internal`, `pkg`, and `router` so generated reports, docs, SQL, and workspace files do not pollute source analysis.

## VSCode

Recommended extensions are listed in `.vscode/extensions.json`:

* Go (`golang.go`)
* SonarQube for IDE / SonarLint (`sonarsource.sonarlint-vscode`)
* YAML support (`redhat.vscode-yaml`)
* EditorConfig (`editorconfig.editorconfig`)

`.vscode/settings.json` includes a placeholder SonarLint connected-mode project:

```json
{
  "sonarlint.connectedMode.project": {
    "connectionId": "your-sonarqube-connection-id",
    "projectKey": "ms-gofiber"
  }
}
```

Replace `connectionId` with the local SonarQube or SonarCloud connection name configured in VSCode.

## Handoff Expectations

For documentation-only changes, review rendered Markdown when practical. For tooling changes, run the relevant command and report any unavailable local tools.
