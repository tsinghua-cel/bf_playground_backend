.PHONY: default bfbackend docker run

GOBIN = $(shell pwd)/build/bin
GO ?= latest
GOFILES_NOVENDOR := $(shell go list -f "{{.Dir}}" ./...)
TAG=latest

default: bfbackend

all: bfbackend

run:
	@docker compose up -d

down:
	@docker compose down


bfbackend:
	go build -o=${GOBIN}/$@ -gcflags "all=-N -l" .
	@echo "Done building."

docker:
	docker build -t bfbackend:${TAG} .

clean:
	rm -fr build/*
