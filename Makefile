.PHONY: test build docker

all: test build

test:
	go test -coverprofile=coverage.out ./...

build:
	go build

docker:
	docker build -t tauffredou/nextver .

.DEFAULT_GOAL: all
