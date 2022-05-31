#!/usr/bin/env bash
set -euxo pipefail

go install github.com/elastic/go-licenser@latest

go mod verify
go-licenser -d
go run .ci/scripts/check_format.go
go run .ci/scripts/check_lint.go

# Run the tests
set +e
export OUT_FILE="build/test-report.out"
mkdir -p build
go test "./..." -v 2>&1 | tee ${OUT_FILE}
status=$?
go install github.com/jstemmer/go-junit-report@latest
go-junit-report > "build/junit-${GO_VERSION}.xml" < ${OUT_FILE}

exit ${status}
