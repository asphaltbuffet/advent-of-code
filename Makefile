SHELL := /bin/bash

.DEFAULT_GOAL := dev
.PHONY: all
all: ## build pipeline
all: mod inst gen build spell lint test

.PHONY: dev
dev: ## dev build - no tools adjustments
dev: gen build spell lint test

.PHONY: ci
ci: ## CI build pipeline
ci: all diff

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## remove files created during build pipeline
	$(call print-target)
	rm -rf dist || true
	rm -rf bin || true
	rm -f '"$(shell go env GOCACHE)/../golangci-lint"'
	go clean -i -cache -testcache -modcache -fuzzcache -x

.PHONY: mod
mod: ## go mod tidy
	$(call print-target)
	go mod tidy
	cd tools && go mod tidy

.PHONY: inst
inst: ## go install tools
	$(call print-target)
	cd tools && go install $(shell cd tools && go list -f '{{ join .Imports " " }}' -tags=tools)

.PHONY: gen
gen: ## go generate
	$(call print-target)
	go generate ./...

.PHONY: build
build: ## goreleaser build
build:
	$(call print-target)
	goreleaser build --rm-dist --single-target --snapshot

.PHONY: spell
spell: ## misspell
	$(call print-target)
	misspell -error -locale=US -w **.md

.PHONY: lint
lint: ## golangci-lint
	$(call print-target)
	golangci-lint run --fix

.PHONY: test
test: ## go test
	$(call print-target)
	mkdir -p bin || true
	go test -race -covermode=atomic -coverprofile=bin/coverage.out -coverpkg=./... ./...
	go tool cover -html=bin/coverage.out -o bin/coverage.html

.PHONY: diff
diff: ## git diff
	$(call print-target)
	git diff --exit-code
	RES=$$(git status --porcelain) ; if [ -n "$$RES" ]; then echo $$RES && exit 1 ; fi

define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef