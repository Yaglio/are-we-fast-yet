#!/bin/bash
set -e # make script fail on first error
SCRIPT_PATH="$(dirname "$0")"
source "$SCRIPT_PATH/../script.inc"
pushd "$SCRIPT_PATH" > /dev/null

if [[ "$1" = "style" ]]
then
  INFO Check Format
  UNFORMATTED="$(gofmt -l .)"
  if [[ -n "$UNFORMATTED" ]]; then
    ERR "gofmt found unformatted files:"
    echo "$UNFORMATTED"
    exit 1
  fi

  INFO Run go vet
  go vet ./...

  INFO Run Unit Tests
  go test ./...
else
  INFO Build Go Benchmarks
  # ADR-011: report on amd64 with the default GOAMD64=v1 (no FMA contraction)
  GOAMD64=v1 go build -o harness .
fi
