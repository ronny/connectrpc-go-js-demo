all: checks test binaries

GO_BUILD_FLAGS ?=
BIN_DIR ?= ./bin/

binaries:
	mkdir -p ${BIN_DIR}
	CGO_ENABLED=0 go build ${GO_BUILD_FLAGS} -o ${BIN_DIR} ./cmd/...

test:
	env TZ=UTC go test -shuffle=on -race -covermode=atomic -coverprofile=coverage.txt -count=1 -timeout=5m  ./...

lint:
	golangci-lint run --fix
	buf lint

format:
	golangci-lint fmt
	buf format -w

tidy:
	go mod tidy

generate:
	go generate ./...
	buf build
	buf generate

checks: tidy format generate lint

ci-checks: checks
	git diff --exit-code

.PHONY: all binaries test lint format tidy generate checks ci-checks
