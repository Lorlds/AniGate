PREFIX ?= $(HOME)/.local
BINDIR ?= $(PREFIX)/bin
BINARY ?= anigate
CONFIG ?= configs/anigate.example.json

.PHONY: build install test vet race verify tools run-http run-stdio clean

build:
	mkdir -p bin
	go build -trimpath -o bin/$(BINARY) ./cmd/anigate

install:
	./scripts/install.sh

test:
	go test ./...

vet:
	go vet ./...

race:
	go test -race ./...

verify:
	./scripts/verify.sh

tools: build
	./bin/$(BINARY) tools --config $(CONFIG)

run-http: build
	./bin/$(BINARY) http --addr 127.0.0.1:8787 --config $(CONFIG)

run-stdio: build
	./bin/$(BINARY) stdio --config $(CONFIG)

clean:
	rm -rf bin dist coverage.out
