#!/usr/bin/env bash
#
# dial.sh - one stop dev script for the Dial project
#
# Usage:
#  ./scripts/dial.sh <command>
#
#  Commands:
#
#  build - Build the CLI/GUI binary
#  run - Build and run the binary (passes through extra args)
#  test - Run go test with race detector + coverage
#  lint - Run go vet + golangcli-lint (if installed) + gofmt check
#  fmt - Run gofmt -w on the whole repo
#  tidy - Run go mod tidy
#  sum - Run go mod verify + print go.sum status
#  check - Run fmt-check + vet + lint + test (what CI should run)
#  clean - Remove build artifacts
#  all - tidy + fmt + lint + test + build
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

BIN_DIR="$ROOT_DIR/bin"
BIN_NAME="dial"
BIN_PATH="$BIN_DIR/$BIN_NAME"

# -- colors --
if [[ -t 1 ]]; then
    BOLD="$(tput bold)"; GREEN="$(tput setaf 2)"; RED="$(tput setaf 1)"
    YELLOW="$(tput setaf 3)"; RESET="$(tput sgr0)"
else
    BOLD=""; GREEN=""; RED=""; YELLOW=""; RESET=""
fi

info() { echo "${BOLD}${YELLOW}==>${RESET} $*"; }
ok() { echo "${BOLD}${GREEN}==>${RESET} $*"; }
fail() { echo "${BOLD}${RED}==>${RESET} $*" >&2; }

cmd_build() {
    info "Building $BIN_NAME..."
    mkdir -p "$BIN_DIR"
    go build -o "$BIN_PATH"
    ok "Built $BIN_PATH"
}

cmd_run() {
    cmd_build
    info "Running $BIN_NAME $*"
    "$BIN_PATH" "$@"
}

cmd_test() {
    info "Running tests..."
    go test ./... -race -covermode=atomic -coverprofile=coverage.out
    ok "Tests passed"
    if command -v go >/dev/null 2>&1; then
        echo
        go tool cover -func=coverage.out | tail -n 1
    fi
}

cmd_fmt() {
    info "Formatting code..."
    gofmt -l -w .
    ok "Formatting complete"
}

cmd_fmt_check() {
    info "Checking formatting..."
    UNFORMATTED="$(gofmt -l .)"
    if [[ -n "$UNFORMATTED" ]]; then
        fail "The following files are not gofmt'd:"
        echo "$UNFORMATTED"
        return 1
    fi
    ok "All files formatted correctly"
}

cmd_vet() {
    info "Running go vet..."
    go vet ./...
    ok "go vet passed"
}

cmd_lint() {
    cmd_fmt_check
    cmd_vet
    if command -v golangci-lint >/dev/null 2>&1; then
        info "Running golangci-lint"
        golangci-lint run ./...
        ok "golangci-lint passed"
    else
        info "golangci-lint not installed, skipping (install: https://golangci-lint.run/usage/install)"
    fi
}

cmd_tidy() {
    info "Running go mod tidy..."
    go mod tidy
    ok "go.mod / go.sum tidied"
}

cmd_sum() {
    info "Verifying go.sum..."
    go mod verify
    ok "go.sum verified"
}

cmd_check() {
    cmd_lint
    cmd_test
}

cmd_clean() {
    info "Cleaning build artifacts..."
    rm -rf "$BIN_DIR" coverage.out
    ok "Clean"
}

cmd_all() {
    cmd_tidy
    cmd_fmt
    cmd_lint
    cmd_test
    cmd_build
    ok "All steps complete"
}

usage() {
    sed -n '2,20p' "$0" | sed 's/^# \{0,1\}//'
}

if [[ $# -lt 1 ]]; then
    usage
    exit 1
fi

COMMAND="$1"
shift || true

case "$COMMAND" in
    build) cmd_build "$@" ;;
    run) cmd_run "$@" ;;
    test) cmd_test "$@" ;;
    fmt) cmd_fmt "$@" ;;
    fmt-check) cmd_fmt_check "$@" ;;
    vet) cmd_vet "$@" ;;
    lint) cmd_lint "$@" ;;
    tidy) cmd_tidy "$@" ;;
    sum) cmd_sum "$@" ;;
    check) cmd_check "$@" ;;
    clean) cmd_clean "$@" ;;
    all) cmd_all "$@" ;;
    *)
        fail "Unknown command: $COMMAND"
        usage
        exit 1
        ;;
esac
