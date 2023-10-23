ifeq ($(CI), true)
	GO = go
else
	GO = go1.21.3
endif


all: compat

build:
	$(GO) build -ldflags="-s -w" ./cmd/rudolphe

compat:
	docker run --rm -it \
		-v ${PWD}:/src \
		-v gopkgmod:/go/pkg/mod \
		-v gobuildcache:/root/.cache/go-build \
		golang:1.21-bullseye \
		bash -c "\
			cd /src && \
			git config --global --add safe.directory /src && \
			go build -ldflags='-s -w' ./cmd/rudolphe \
		"

style:
	$(GO) fmt ./...
	golangci-lint run

test:
	$(GO) test ./...

