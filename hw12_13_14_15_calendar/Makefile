MIGRATIONS = "./migrations"
BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

.PHONY: build 
build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

.PHONY: run
run: build
	$(BIN) --config ./configs/config.toml

.PHONY: build-img
build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

.PHONY: 
run-img: build-img
	docker run $(DOCKER_IMG)

.PHONY: version 
version: build
	$(BIN) version

.PHONY: test
test:
	go test -race ./internal/...

.PHONY: install-lint-deps
install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.60.2

.PHONY: lint
lint: install-lint-deps
	golangci-lint run  --config=.golangci.yml ./...

.PHONY: local-migration-up
local-migration-up:
	migrate -database $(CALENDAR_DSN) -path $(MIGRATIONS) up

.PHONY: local-migration-down
local-migration-down:
	migrate -database $(CALENDAR_DSN) -path $(MIGRATIONS) down

