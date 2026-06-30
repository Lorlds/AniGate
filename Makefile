PREFIX ?= $(HOME)/.local
BINDIR ?= $(PREFIX)/bin
MINI_CONFIG ?= configs/anigate.mini.example.json
MAX_CONFIG ?= configs/anigate.max.example.json

.PHONY: build install test vet race verify tools run-http-mini run-http-max run-stdio-mini run-stdio-max clean

build:
	mkdir -p bin
	go build -trimpath -o bin/anigate-mini ./cmd/anigate-mini
	go build -trimpath -o bin/anigate-max ./cmd/anigate-max
	go build -trimpath -o bin/anigate ./cmd/anigate

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
	./bin/anigate-mini tools --config $(MINI_CONFIG)
	./bin/anigate-max tools --config $(MAX_CONFIG)

run-http-mini: build
	./bin/anigate-mini http --addr 127.0.0.1:8787 --config $(MINI_CONFIG)

run-http-max: build
	./bin/anigate-max http --addr 127.0.0.1:8788 --config $(MAX_CONFIG)

run-stdio-mini: build
	./bin/anigate-mini stdio --config $(MINI_CONFIG)

run-stdio-max: build
	./bin/anigate-max stdio --config $(MAX_CONFIG)

clean:
	rm -rf bin dist coverage.out
