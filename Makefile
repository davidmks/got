.PHONY: build test lint check clean

build:
	go build -o got .

test:
	go test ./...

lint:
	golangci-lint run

check: lint test

clean:
	rm -f got
