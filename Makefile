ifeq ($(CI), true)
	GO = go
else
	GO = go1.21.3
endif


all: build

build:
	$(GO) build -ldflags="-s -w" ./cmd/rudolphe


style:
	$(GO) fmt ./...
	golangci-lint run

test:
	$(GO) test ./...

