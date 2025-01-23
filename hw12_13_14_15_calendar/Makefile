MIGRATIONS = "./migrations"
BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

LOCAL_BIN:=$(CURDIR)/bin

.PHONY: install-deps
install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.2
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.24.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.24.0

.PHONY: install-mockery
install-mockery:
	go install github.com/vektra/mockery/v2@v2.50.0

.PHONY: generate
generate:
	rm -rf internal/pb
	mkdir -p internal/pb

	protoc \
		--proto_path=api/ \
		--go_out=internal/pb \
		--go_opt=paths=source_relative \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		--go-grpc_out=internal/pb \
		--go-grpc_opt=paths=source_relative \
        --plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
        --grpc-gateway_out=internal/pb --grpc-gateway_opt=paths=source_relative \
        --plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
        api/*.proto


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

