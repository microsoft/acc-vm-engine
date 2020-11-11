
build: generate
	go build ./cmd/acc-vm-engine/

generate: bootstrap
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor
	go generate $(GOFLAGS) -v ./pkg/engine/fileloader.go

.PHONY: bootstrap
bootstrap:
	which go-bindata || go get github.com/go-bindata/go-bindata/...

.PHONY: test
test:
	./acc-vm-engine generate -c ./test/tvm-ub1804.json -o ./_output/tvm
	./acc-vm-engine generate -c ./test/cvm-ub1804.json -o ./_output/cvm

.PHONY: clean
clean:
	rm -f ./acc-vm-engine
