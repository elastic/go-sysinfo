#!/usr/bin/env bash
set -xeuo pipefail

export CGO_ENABLED=0

# Test that there are no compilation issues when not using CGO. This does
# not imply that all of these targets are supported without CGO. It's only
# a sanity check for build tag issues.
GOOS=aix     GOARCH=ppc64    go build ./...
GOOS=darwin  GOARCH=amd64    go build ./...
GOOS=darwin  GOARCH=arm64    go build ./...
GOOS=linux   GOARCH=386      go build ./...
GOOS=linux   GOARCH=amd64    go build ./...
GOOS=linux   GOARCH=arm      go build ./...
GOOS=linux   GOARCH=arm64    go build ./...
GOOS=linux   GOARCH=mips     go build ./...
GOOS=linux   GOARCH=mips64   go build ./...
GOOS=linux   GOARCH=mips64le go build ./...
GOOS=linux   GOARCH=mipsle   go build ./...
GOOS=linux   GOARCH=ppc64    go build ./...
GOOS=linux   GOARCH=ppc64le  go build ./...
GOOS=linux   GOARCH=riscv64  go build ./...
GOOS=linux   GOARCH=s390x    go build ./...
GOOS=windows GOARCH=amd64    go build ./...
GOOS=windows GOARCH=arm      go build ./...
GOOS=windows GOARCH=arm64    go build ./...
