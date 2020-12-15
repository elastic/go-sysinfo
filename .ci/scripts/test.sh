#!/usr/bin/env bash
set -euxo pipefail

GO111MODULE=off go get -u github.com/elastic/go-licenser

go mod verify
go-licenser -d
go run .ci/scripts/check_format.go
go run .ci/scripts/check_lint.go
go test -v ./...