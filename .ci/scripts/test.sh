#!/usr/bin/env bash
set -euxo pipefail

GO111MODULE=off go get -u github.com/elastic/go-licenser

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
go get -v -u github.com/jstemmer/go-junit-report
go-junit-report > "build/junit-${GO_VERSION}.xml" < ${OUT_FILE}

exit ${status}
