.PHONY: build test race vet run fmt clean

BINARY := bin/controller

build:
	go build -o $(BINARY) ./cmd/controller

test:
	go test ./...

race:
	go test -race ./...

vet:
	go vet ./...

fmt:
	gofmt -w cmd internal

run:
	go run ./cmd/controller

clean:
	rm -rf bin
