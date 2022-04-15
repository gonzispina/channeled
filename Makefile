# cross parameters
SHELL:=/bin/bash -O extglob
BINARY=sigi-api
VERSION=0.1.0

LDFLAGS=-ldflags "-X main.Version=${VERSION}"

# Build step, generates the binary.
build: clean generate
	@go build ${LDFLAGS} -o ${BINARY} cmd/web/main/*.go
