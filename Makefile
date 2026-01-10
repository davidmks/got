.PHONY: build test lint check clean setup

build:
	go build -o got .

test:
	go test ./...

lint:
	golangci-lint run

check: lint test

clean:
	rm -f got

setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	pre-commit install
