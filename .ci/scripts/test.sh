#!/usr/bin/env bash
set -euxo pipefail

go install gotest.tools/gotestsum@latest

# Run the tests
export OUT_FILE="build/test-report.out"
mkdir -p build
gotestsum --format testname --junitfile "build/junit-${GO_VERSION}.xml" -- -tags integration ./...
