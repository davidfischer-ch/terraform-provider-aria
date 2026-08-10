#!/usr/bin/env bash
#
# Local pre-push check. Mirrors the GitHub Actions "Tests" workflow (.github/workflows/test.yml):
# build, lint, generate, unit tests.
#
# Usage:
#   ./scripts/check.sh              # run everything
#   ./scripts/check.sh build lint   # run only the named steps
#   SKIP_GENERATE=1 ./scripts/check.sh
#
# Steps: build vet lint generate test
set -euo pipefail

cd "$(dirname "$0")/.."

# golangci-lint release that CI installs (via `version: latest`). Pinned here for reproducibility.
# The build's Go version must be at least the `go` directive in go.mod.
GOLANGCI_LINT_VERSION=${GOLANGCI_LINT_VERSION:-v2.12.2}

# --- Pretty Logging -------------------------------------------------------------------------------

if [ -t 1 ]; then
  bold=$'\e[1m'
  green=$'\e[32m'
  red=$'\e[31m'
  dim=$'\e[2m'
  reset=$'\e[0m'
else
  bold=
  green=
  red=
  dim=
  reset=
fi

step() { printf '\n%s==> %s%s\n' "$bold" "$1" "$reset"; }
ok()   { printf '%s✓ %s%s\n' "$green" "$1" "$reset"; }
die()  { printf '%s✗ %s%s\n' "$red" "$1" "$reset" >&2; exit 1; }

# --- Resolve golangci-lint ------------------------------------------------------------------------

# Prefer the prebuilt release (it matches CI exactly): its Go build version caps the supported `go`
# directive. Install into ./bin when absent or the version differs.
find_golangci() {
  local want=${GOLANGCI_LINT_VERSION#v}
  for bin in "$PWD/bin/golangci-lint" "$(command -v golangci-lint 2>/dev/null || true)"; do
    [ -x "$bin" ] || continue
    if "$bin" version 2>/dev/null | grep -q "$want"; then
      echo "$bin"
      return 0
    fi
  done
  return 1
}

ensure_golangci() {
  if GCL=$(find_golangci); then
    return 0
  fi

  step "Installing golangci-lint $GOLANGCI_LINT_VERSION into ./bin"
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh \
    | sh -s -- -b "$PWD/bin" "$GOLANGCI_LINT_VERSION" \
    || die 'failed to install golangci-lint (network?). Install it manually and re-run.'
  GCL=$PWD/bin/golangci-lint
}

# --- Steps ----------------------------------------------------------------------------------------

run_build() {
  step 'go build ./...'
  go build ./... && ok 'build'
}

run_vet() {
  step 'go vet ./...'
  go vet ./... && ok 'vet'
}

run_lint() {
  ensure_golangci
  step 'golangci-lint config verify'
  "$GCL" config verify && ok 'config valid'
  step 'golangci-lint run ./...'
  "$GCL" run ./... && ok 'lint'
}

run_generate() {

  if [ "${SKIP_GENERATE:-}" = '1' ]; then
    printf '%sskipping generate%s\n' "$dim" "$reset"
    return 0
  fi

  step 'go generate ./... (docs/examples must be up to date)'
  go generate ./...

  # Scope to generate's own outputs to keep unrelated pending changes from tripping this. main.go
  # runs tfplugindocs (writes docs/) and `terraform fmt` (writes examples/). --porcelain also
  # surfaces new untracked files, which `git diff` would miss.
  local dirty
  dirty=$(git status --porcelain -- docs/ examples/)

  if [ -n "$dirty" ]; then
    printf '%s\n' "$dirty"
    die "generate changed docs/ or examples/; run 'go generate ./...' and commit the result"
  fi

  ok 'generate (docs/examples up to date)'
}

run_test() {
  step 'go test ./internal/provider/'
  TF_ACC='' ARIA_HOST='my-aria-instance.net' ARIA_REFRESH_TOKEN='faketokenhere' \
    go test -cover ./internal/provider/ && ok 'tests'
}

# --- Dispatch -------------------------------------------------------------------------------------

steps=("$@"); [ ${#steps[@]} -eq 0 ] && steps=(build vet lint generate test)
for s in "${steps[@]}"; do
  case "$s" in
    build)    run_build ;;
    vet)      run_vet ;;
    lint)     run_lint ;;
    generate) run_generate ;;
    test)     run_test ;;
    *) die "unknown step: $s (valid: build vet lint generate test)" ;;
  esac
done

printf '\n%s%sAll checks passed.%s\n' "$bold" "$green" "$reset"
