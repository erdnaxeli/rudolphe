all: compat

build:
	go build -ldflags="-s -w" ./cmd/rudolphe

compat:
	docker run --rm -it \
		-v ${PWD}:/src \
		-v gopkgmod:/go/pkg/mod \
		-v gobuildcache:/root/.cache/go-build \
		golang:1.22-bullseye \
		bash -c "\
			cd /src && \
			git config --global --add safe.directory /src && \
			go build -ldflags='-s -w' ./cmd/rudolphe \
		"

style:
	go fmt ./...
	golangci-lint run

test:
	$(GO) test ./...

