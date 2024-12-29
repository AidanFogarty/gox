PROJECT = $(shell basename $(PWD))

GOPATH ?= $(shell go env GOPATH)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

.DEFAULT_GOAL := help

## build: Build the binary.
build:
	go build -ldflags="-w -s"

## test: Run all the tests in the project
test:
	go test ./... -v

## fmt: Formats all Go code in the project
fmt:
	go fmt ./...

vet:
	go vet ./...

## ci: Run linter, formatter and tests
ci: build fmt vet test

# Thanks to: https://github.com/azer/go-makefile-example
.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo