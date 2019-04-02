.PHONY: test build docker

all: test build

test:
	go test ./...

build:
	go build

docker:
	docker built -t tauffredou/nextver .

.DEFAULT_GOAL: all
