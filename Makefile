.PHONY: all build test filedialogdemo

build:
	go build -v dinky.go

test:
	go test -v ./...

filedialogdemo:
	go build -v -o filedialogdemo ./cmd/filedialogdemo
