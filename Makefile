PACKAGES=$(shell go list ./ ...)
VERSION=$(shell git describe --always --tags --dirty)

define echotask
	@tput setaf 6
	@echo -n "  $1"
	@tput setaf 8
	@echo -n " - "
	@tput sgr0
	@echo $2
endef

define echoversion
	@echo -n "  $1 "
	@tput setaf 5
	@echo $2
	@tput sgr0
endef

help:
	@echo
	$(call echoversion,"Project crytocurrency",$(VERSION))
	@echo
	$(call echotask,"deps","install all dependencies")
	$(call echotask,"format","formats code with gofumpt")
	$(call echotask,"formatcheck","checks if code is formatted with gofumpt")
	$(call echotask,"lint","run all linters")
	$(call echotask,"test","run all tests")
	$(call echotask,"test_ci","run tests for the ci job")
	@echo

test: ## Run tests using gotestsum.
	@gotestsum \
	    --format=dots-v2 -- \
	    -timeout=30000ms \
	    -covermode=set \
	    -coverprofile=.coverage.out ${PACKAGES}

test_ci: ## Run tests using normal test runner for ci output
	@go test -coverpkg ./... \
	    -coverprofile .coverage.out ${PACKAGES} && go tool cover -func=.coverage.out

lint: ## Check the code using various linters and static checkers.
	golangci-lint run --timeout=5m ./...

deps: ## Install all dependencies for this project
	go mod download
	@make install_lint
	@make install_gofumpt
	@make install_gotestsum

install_lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b $(shell go env GOPATH)/bin v1.49.0

install_gotestsum:
	go install gotest.tools/gotestsum@v1.8.2

install_gofumpt:
	go install mvdan.cc/gofumpt@v0.3.1

format:
	gofumpt -l -w .

formatcheck:
	test `gofumpt -l . | wc -l` -eq 0


.PHONY: help test lint install build

.DEFAULT_GOAL := help
