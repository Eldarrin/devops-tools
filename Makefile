SHELL				= /bin/bash
ARCH				= x86_64
CGO					?=0
CGO_ENABLED			?=0
TARGET_OS			?= linux

GIT_VERSION			= $(shell git describe --tags)
GIT_COMMIT			= $(shell git rev-parse --short HEAD)
DATE				= $(shell date +%Y%m%d)

TEST_CLUSTER_NAME	= test-cluster

ifeq (,$(shell go env GOBIN))
GOBIN				= $(shell go env GOPATH)/bin
else
GOBIN				= $(shell go env GOBIN)
endif

#################################################
# Configure the development environment			#
#################################################

##@ configure

configure: configure-dev

configure-dev: ## Configure the development environment
		go mod download
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.46.2

configure-ci: ## Configure the CI environment
		go mod download
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b . v1.46.2

#################################################
# Development									#
#################################################

##@ development

development: generate fmt vet golangci

check: fmt vet golangci

generate: ## Generate APIs and echo server from OpenAPI spec
		go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen \
			-config pkg/api/config.yml pkg/api/migrator.yml

fmt: ## fmt the code
		go fmt ./...

vet: ## vet the code
		go vet ./...

golangci: ## lint the code
		golangci-lint run ./...

#################################################
# Build											#
#################################################

##@ build

build: check backend client docker-build

backend: ## compile the backend
		go build -o backend cmd/backend/main.go

client: ## compile the client
		go build -o client cmd/client/main.go

docker-build: ## build the docker image
		buildah build -t tools-backend .

manifests: ## generate the kubernetes manifests



#################################################
# Test											#
#################################################

##@ test

test: ## run the tests
		go test -v ./...

##@ test-cluster

#################################################
# Deploy										#
#################################################

##@ deploy


#################################################
# Clean											#
#################################################

##@ clean

#################################################
# Publish										#
#################################################

##@ publish

#################################################
# Help											#
#################################################

##@ help

help: ## Display this help
		@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


