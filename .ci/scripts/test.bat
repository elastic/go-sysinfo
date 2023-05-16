go install gotest.tools/gotestsum@latest

mkdir -p build
SET OUT_FILE=build\output-report.out
gotestsum --format testname --junitfile build\junit-%GO_VERSION%.xml -- ./...
