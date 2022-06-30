go install github.com/redanthrax/go-licenser@latest

go mod verify
go-licenser -d
go run .ci/scripts/check_format.go
go run .ci/scripts/check_lint.go

mkdir -p build
SET OUT_FILE=build\output-report.out
go test "./..." -v > %OUT_FILE% | type %OUT_FILE%

go install github.com/jstemmer/go-junit-report@latest
go-junit-report > build\junit-%GO_VERSION%.xml < %OUT_FILE%
