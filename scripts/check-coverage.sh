#!/usr/bin/env bash
set -euo pipefail

coverage_file="${1:-coverage.out}"
required="${2:-100.0}"

if [[ ! -f "$coverage_file" ]]; then
  printf 'coverage file not found: %s\n' "$coverage_file" >&2
  exit 1
fi

actual="$(go tool cover -func="$coverage_file" | awk '/^total:/ {gsub("%", "", $3); print $3}')"
if [[ -z "$actual" ]]; then
  printf 'failed to read total coverage from %s\n' "$coverage_file" >&2
  exit 1
fi

if [[ "$actual" != "$required" ]]; then
  printf 'coverage %s%% does not meet required %s%%\n' "$actual" "$required" >&2
  exit 1
fi

printf 'coverage %s%%\n' "$actual"
