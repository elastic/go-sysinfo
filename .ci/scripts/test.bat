go install github.com/elastic/go-licenser@latest
go install gotest.tools/gotestsum@latest

go mod verify
go-licenser -d
go run .ci/scripts/check_format.go
go run .ci/scripts/check_lint.go

mkdir -p build
SET OUT_FILE=build\output-report.out
gotestsum --format testname --junitfile build\junit-%GO_VERSION%.xml -- ./...
