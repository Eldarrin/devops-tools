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
# Development									#
#################################################

##@ development

development: generate fmt vet golangci

generate:
		go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen \
			-config pkg/api/config.yml pkg/api/migrator.yml

fmt:
		go fmt ./...

vet:
		go vet ./...

golangci:
		golangci-lint run ./...

#################################################
# Build											#
#################################################

##@ build

build: fmt vet golangci backend client docker-build

backend:
		go build -o backend cmd/backend/main.go

client:
		go build -o client cmd/client/main.go

docker-build:
		buildah build -t tools-backend .

#################################################
# Test											#
#################################################

##@ test

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


