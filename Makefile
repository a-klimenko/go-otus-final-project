BIN := "./bin/rotator"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate:
	go generate ./...

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/rotator

run:
	./deployments/create_env.sh && \
	docker-compose --env-file deployments/.env -f deployments/docker-compose.yaml up -d

down:
	docker-compose -f deployments/docker-compose.yaml down

integration-tests: run
	docker-compose -f deployments/docker-compose.tests.yaml up --build && \
	docker-compose -f deployments/docker-compose.tests.yaml down && \
	docker-compose -f deployments/docker-compose.yaml down

test:
	go test -race -count 100 ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run lint test
