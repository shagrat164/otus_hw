BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.5

lint: install-lint-deps
	golangci-lint run ./...

migrate-up:
	goose -dir migrations postgres "postgres://calendar:calendar@localhost:5432/calendar?sslmode=disable" up

migrate-down:
	goose -dir migrations postgres "postgres://calendar:calendar@localhost:5432/calendar?sslmode=disable" down

compose-up:
	docker-compose -f deployments/docker-compose.yaml up -d

compose-down:
	docker-compose -f deployments/docker-compose.yaml down

clean:
	rm -rf ./bin

.PHONY: build run build-img run-img version test lint migrate-up migrate-down compose-up compose-down clean
