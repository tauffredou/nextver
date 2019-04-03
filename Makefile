.PHONY: test build docker

all: test build

test:
	go test ./...

build:
	go build

docker:
	docker build -t tauffredou/nextver .

.DEFAULT_GOAL: all
