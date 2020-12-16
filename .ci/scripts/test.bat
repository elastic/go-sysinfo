SET GO111MODULE=off
go get -u github.com/elastic/go-licenser
SET GO111MODULE=on

go mod verify
go-licenser -d
go run .ci/scripts/check_format.go
go run .ci/scripts/check_lint.go

mkdir -p build
SET OUT_FILE=build\output-report.out
go test "./..." -v > %OUT_FILE% | type %OUT_FILE%

go get -v -u github.com/jstemmer/go-junit-report
go-junit-report > build\junit-%GO_VERSION%.xml < %OUT_FILE%
