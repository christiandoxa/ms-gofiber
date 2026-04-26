# Developer Tooling

## Baseline Commands

Use these commands from the repository root:

```bash
make fmt
make test
make coverage
make lint
make sonar
make run
```

Equivalent raw commands:

```bash
go fmt ./...
go test ./...
go test -coverprofile=coverage.out ./...
golangci-lint run ./...
sonar-scanner
go run ./cmd/ms-gofiber
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

Install `golangci-lint` before running `make lint`.

## SonarQube

`sonar-project.properties` provides a headless scanner baseline. Run `make sonar` only after `sonar-scanner`,
`SONAR_HOST_URL`, and `SONAR_TOKEN` are configured for your SonarQube project.

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
