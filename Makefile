VERSION := $(shell ./getversion.sh)
LDFLAGS := -X "github.com/baotingfang/go-pivnet-client/cmd.Version=$(VERSION)"

.PHONY: list
list:
	@sh -c "$(MAKE) -p no_targets__ 2>/dev/null | \
	awk -F':' '/^[a-zA-Z0-9][^\$$#\/\\t=]*:([^=]|$$)/ {split(\$$1,A,/ /);for(i in A)print A[i]}' | \
	grep -v Makefile | \
	grep -v '%' | \
	grep -v '__\$$' | \
	sort"

.PHONY: depend
depend:
	go get golang.org/x/tools/cmd/stringer
	go get github.com/onsi/ginkgo/ginkgo
	GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.25.1
	GO111MODULE=off go get -u github.com/maxbrunsfeld/counterfeiter

.PHONY: generate
generate:
	go generate ./...

.PHONY: build
build:
	go build -ldflags '$(LDFLAGS)' -o pivnet-client ./entry/go-pivnet-client.go

.PHONY: unit
unit:
	go test -cover -race ./...

.PHONY: clean
clean:
	rm -f ./pivnet-client

.PHONY: check
check:
	gofmt -w .
	go test -cover -count=1 ./...
	golangci-lint run -v
