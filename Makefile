BIN					:= deployer
OUTPUT_DIR     		:= build
TMP_DIR             := .tmp
RELEASE_VER   		:= $(shell git rev-parse --short HEAD)
RELEASE_BRANCH      := $(shell git rev-parse --abbrev-ref HEAD | sed 's/\//\-/g')
DOCKER_TAG          := $(RELEASE_BRANCH)-$(RELEASE_VER)
BUILD_TIME 			:= $(shell date +%Y-%m-%dT%T%z)
TEST_PACKAGES  		:= $(shell go list ./...)
INTERNAL_PKG_DIR	:= github.com/vmogilev/$(BIN)/internal

.PHONY: help
.DEFAULT_GOAL := help

gomod/init: ## init go mod
	go mod init github.com/vmogilev/$(BIN)

gomod/vendor: ## vendor Go dependencies
	## see https://github.com/golang/go/issues/35164#issuecomment-546503518
	## without GOSUMDB=off we get:
	##   verifying github.com/xxxx/xxxxxx@v1.0.1: github.com/xxxx/xxxxxx@v1.0.1: reading https://sum.golang.org/lookup/github.com/xxxx/xxxxxx@v1.0.1: 410 Gone
	GOSUMDB=off go mod tidy
	GOSUMDB=off go mod vendor -v

test/all: test/fmt test/codecov test/lint ## Does all of non e2e testing

test: ## Perform unit tests
	@echo "+ $@"
	GOFLAGS=-mod=vendor go test -v -count 1 -cover $(TEST_PACKAGES)

test/fmt: ## Check if all files (excluding vendor) conform to fmt
	@echo "+ $@"
	test -z $(shell echo $(shell go fmt $(shell go list ./...)) | tr -d "[:space:]")

test/codecov:
	@echo "+ $@"
	GOFLAGS=-mod=vendor go test -v -count 1 -coverprofile=coverage.txt -covermode=atomic $(TEST_PACKAGES)

test/lint: ## Verifies `golint` passes
	@echo "+ $@"
	@golangci-lint run

build: clean ## Build binary and output to /build
	GOOS=darwin CGO_ENABLED=0 GOFLAGS=-mod=vendor go build -v -installsuffix cgo -ldflags "-X $(INTERNAL_PKG_DIR)/version.Number=$(RELEASE_VER) -X $(INTERNAL_PKG_DIR)/version.BuildTime=$(BUILD_TIME)" -o $(OUTPUT_DIR)/$(BIN)-darwin .
	@echo "created: $(OUTPUT_DIR)/$(BIN)-darwin"
	GOOS=linux CGO_ENABLED=0 GOFLAGS=-mod=vendor go build -a -installsuffix cgo -ldflags "-X $(INTERNAL_PKG_DIR)/version.Number=$(RELEASE_VER) -X $(INTERNAL_PKG_DIR)/version.BuildTime=$(BUILD_TIME)" -o $(OUTPUT_DIR)/$(BIN)-linux .
	@echo "created: $(OUTPUT_DIR)/$(BIN)-linux"

clean: ## Removing binary in output dir and stop and remove the containers. removed unused vendor pkgs
	$(RM) $(OUTPUT_DIR)/$(BIN)
	$(RM) $(OUTPUT_DIR)/$(BIN)-darwin
	$(RM) $(OUTPUT_DIR)/$(BIN)-linux

docker/build: ## builds the latest docker image
	docker build -t $(BIN):latest .

docker/run: ## runs the docker image
	docker run --rm -i $(BIN):latest execute

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
