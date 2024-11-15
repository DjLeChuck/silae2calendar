.DEFAULT_GOAL := help
.PHONY:help fmt vet tidy build

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## fmt: format code
fmt:
	go fmt ./...

## vet: format code and execute got vet
vet: fmt
	go vet ./...

## tidy: format code and tidy modfile
tidy:
	go fmt ./...
	go mod tidy -v

## build: build the application
build:
	GOARCH=amd64 GOOS=darwin go build -o dist/silae2calendar-darwin
	GOARCH=amd64 GOOS=linux go build -o dist/silae2calendar-linux
