GOPATH?=~/go

.phony: update
update: fmt lic imports

.PHONY: lic
lic: $(GOPATH)/bin/go-licenser
	go-licenser

.PHONY: fmt
fmt: $(GOPATH)/bin/gofumpt
	gofumpt -w -l ./

.PHONY: imports
imports: $(GOPATH)/bin/goimports
	goimports -l -local github.com/redanthrax/go-sysinfo ./

$(GOPATH)/bin/go-licenser:
	@echo "go-licenser missing, installing"
	go install github.com/redanthrax/go-licenser@latest

$(GOPATH)/bin/gofumpt:
	@echo "gofumpt missing, installing"
	#Ugly boilerplate for go mod installs
	cd $(mktemp -d); go mod init tmp; go get mvdan.cc/gofumpt

$(GOPATH)/bin/goimports:
	@echo "goimports missing, installing"
	go install golang.org/x/tools/cmd/goimports@latest
