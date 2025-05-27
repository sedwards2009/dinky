.PHONY: all build test filedialogdemo scrollbardemo

build:
	go build -v dinky.go

test:
	go test -v ./...

filedialogdemo:
	go build -v -o filedialogdemo ./cmd/filedialogdemo

scrollbardemo:
	go build -v -o scrollbardemo ./cmd/scrollbardemo
