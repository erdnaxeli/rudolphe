all:
	shards build

doc:
	crystal doc

init-dev:
	shards install

lint:
	crystal tool format
	./bin/ameba src spec

static:
	docker run --rm -it -v ${PWD}:/workspace -w /workspace crystallang/crystal:1.2.2-alpine ./build_static.sh

test:
	crystal spec  --error-trace

.PHONY: test
