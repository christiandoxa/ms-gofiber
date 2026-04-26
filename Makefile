GO ?= go
GOLANGCI_LINT ?= $(shell command -v golangci-lint 2>/dev/null || printf '%s/bin/golangci-lint' "$$($(GO) env GOPATH)")
APP_PKG ?= ./cmd/ms-gofiber
SONAR_SCANNER ?= sonar-scanner

.PHONY: help tidy fmt test coverage lint sonar run verify

help:
	@printf '%s\n' 'Targets:'
	@printf '%s\n' '  tidy    - sync Go module files'
	@printf '%s\n' '  fmt     - format Go packages'
	@printf '%s\n' '  test    - run all Go tests'
	@printf '%s\n' '  coverage - run tests with coverage profile'
	@printf '%s\n' '  lint    - run golangci-lint'
	@printf '%s\n' '  sonar   - run Sonar scanner with sonar-project.properties'
	@printf '%s\n' '  run     - run the service locally'
	@printf '%s\n' '  verify  - run fmt, test, and lint'

tidy:
	$(GO) mod tidy

fmt:
	$(GO) fmt ./...

test:
	$(GO) test ./...

coverage:
	$(GO) test -coverprofile=coverage.out ./...

lint:
	$(GOLANGCI_LINT) run ./...

sonar: coverage
	$(SONAR_SCANNER)

run:
	$(GO) run $(APP_PKG)

verify: fmt test lint
