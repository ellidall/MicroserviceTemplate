all: build test check

modules:
	go mod tidy

build: modules
	go build -o bin/microservicetemplate ./cmd/microservicetemplate

test:
	go test ./...

check:
	golangci-lint run