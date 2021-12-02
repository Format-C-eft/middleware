GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.17","$(shell printf "$(GO_VERSION_SHORT)\n1.17" | sort -V | head -1)")
$(error NEED GO VERSION >= 1.17. Found: $(GO_VERSION_SHORT))
endif

export GO111MODULE = on

SERVICE_NAME = middleware
SERVICE_PATH = Format-C-eft/middleware
# SERVICE_VERSION = 1.0.0.1

COMMIT_HASH = $(shell git rev-parse HEAD)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
TIME_BUILD = $(shell date +"%FT%H:%M%z")

GOARCH = amd64
BUILD_DIR = ./bin

LDFLAGS = -ldflags "-X 'github.com/$(SERVICE_PATH)/internal/config.branch=$(BRANCH)' \
					-X 'github.com/$(SERVICE_PATH)/internal/config.commitHash=$(COMMIT_HASH)'\
					-X 'github.com/$(SERVICE_PATH)/internal/config.timeBuild=$(TIME_BUILD)'"

.PHONY: run
run:
	go run cmd/middleware/main.go

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: deps
deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0

.PHONY: build
build: build-linux build-darwin

.PHONY: build-linux
build-linux:
	go mod download && CGO_ENABLED=0 \
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILD_DIR}/middleware-linux-${GOARCH} ./cmd/middleware/main.go;

.PHONY: build-darwin
build-darwin:
	go mod download && CGO_ENABLED=0 \
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILD_DIR}/middleware-darwin-${GOARCH} ./cmd/middleware/main.go;

PHONY: build-docker
build-docker:
	go mod download && CGO_ENABLED=0 \
	go build ${LDFLAGS} -o ${BUILD_DIR}/middleware-docker ./cmd/middleware/main.go;

.PHONY: test
test:
	go test -v -race -timeout 30s -coverprofile cover.out ./...
	go tool cover -func cover.out | grep total | awk '{print $$3}'