.PHONY: build
build:
	go build ./cmd/acc-vm-engine/

.PHONY: test
test:
	./acc-vm-engine generate -c ./test/tvm-ub1804.json

.PHONY: generate
generate:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor
	go generate $(GOFLAGS) -v ./pkg/engine/fileloader.go

.PHONY: bootstrap
bootstrap:
ifndef HAS_GOBINDATA
	go get github.com/go-bindata/go-bindata/...
endif
