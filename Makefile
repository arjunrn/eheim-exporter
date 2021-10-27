.PHONY: lint build

all: build

lint:
	golangci-lint run -c .golangci.yaml 

build:
	go build -o bin/eheim-exporter main.go

test:
	go test ./...
