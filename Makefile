.DEFAULT_GOAL := run 
.PHONY: fmt vet build test

fmt:
	@go fmt ./...
vet: fmt
	@go vet ./...
build: vet
	@go build -o bin/bcpg
test:
	@go test -v ./...
run:	build
	@./bin/bcpg
