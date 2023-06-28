PROJECT_NAME := fastagi_server
APP_NAME ?= fast-agi

.PHONY: all lint fmt test cover yaegi_test dep
all: fmt lint test cover

lint:
	golangci-lint run --fix
fmt:
	gofmt -s -w ./
test: ## Run data race detector
	go test -v -timeout 600s ./...
cover: ## Run unittests
	go test -cover -short -coverprofile=cover.out ./...
	go tool cover -html=cover.out
yaegi_test:
	yaegi test -v .
dep: ## Get the dependencies
	go mod download

.DEFAULT_GOAL := all
