VERSION=$(shell git describe --always --tags --dirty)
LDFLAGS=-ldflags "-s -w -X main.version=${VERSION}"
BUILD_PARAMS=CGO_ENABLED=0
TEST=$(shell go list ./... | grep -v /test/)
ENTRYPOINT=cmd/*

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
	$(call echotask,"run","run the project")
	$(call echotask,"deps","install all dependencies")
	$(call echotask,"format","formats code with gofumpt")
	$(call echotask,"formatcheck","checks if code is formatted with gofumpt")
	$(call echotask,"lint","run all linters")
	$(call echotask,"test","run all tests")
	$(call echotask,"test_html","run tests showing coverage in the browser")
	$(call echotask,"test_ci","run tests using normal test runner for ci output")
	$(call echotask,"build","compile project for the current platform")
	$(call echotask,"build_all","compile project for all supported platforms")
	$(call echotask,"build_amd64","compile project for amd64")
	$(call echotask,"build_arm64v8","compile project for arm64v8")
	$(call echotask,"build_windows","compile project for windows")
	$(call echotask,"docker_build","build all docker containers")
	$(call echotask,"docker_amd64","create amd64 container")
	$(call echotask,"docker_arm64v8","create arm64v8 container")
	$(call echotask,"clean","clean the build directory")
	@echo

run:
	cd cmd && go run .

test:
	@gotestsum \
	    --format=dots-v2 -- \
	    -timeout=30000ms \
	    -covermode=set \
	    -coverprofile=.coverage.out ${TEST}

test_html:
	@$(MAKE) test
	@go tool cover -html=.coverage.out

test_ci:
	@go test -coverpkg ./... \
	    -coverprofile .coverage.out ${TEST} && go tool cover -func=.coverage.out

build_all: build_amd64 build_arm64v8 build_windows

build: # Create a production binary for current platform.
	@${BUILD_PARAMS} go build ${LDFLAGS} -o \
	    build/crypto-${VERSION}-$(shell go env GOHOSTOS)-$(shell go env GOHOSTARCH) ${ENTRYPOINT}

build_amd64: # Create Linux AMD64 binary.
	@GOARCH=amd64 GOOS=linux ${BUILD_PARAMS} go build ${LDFLAGS} \
		-o build/crypto-${VERSION}-linux-amd64 ${ENTRYPOINT}
	@touch build/crypto-linux-amd64
	@ln -sf crypto-${VERSION}-linux-amd64 build/crypto-linux-amd64

build_arm64v8: # Create default linux arm (armv8) binary.
	@GOARCH=arm GOOS=linux ${BUILD_PARAMS} go build ${LDFLAGS} \
		-o build/crypto-${VERSION}-linux-arm64v8 ${ENTRYPOINT}
	@touch build/crypto-linux-arm64v8
	@ln -sf crypto-${VERSION}-linux-arm64v8 build/crypto-linux-arm64v8

build_windows: # Create Windows binaries.
	@GOARCH=386 GOOS=windows ${BUILD_PARAMS} go build ${LDFLAGS} -o \
	    build/crypto-${VERSION}-windows-386.exe ${ENTRYPOINT}
	@GOARCH=amd64 GOOS=windows ${BUILD_PARAMS} go build ${LDFLAGS} -o \
	    build/crypto-${VERSION}-windows-amd64.exe ${ENTRYPOINT}

docker_build: docker_amd64 docker_arm64v8

docker_amd64: # Build amd64 docker container
	docker build . -f ./docker/amd64/Dockerfile -t blockchain-node-amd64:${VERSION}

docker_arm64v8: # Build arm64v8 docker container
	docker build . -f ./docker/arm64v8/Dockerfile -t blockchain-node-arm64v8:${VERSION}

lint: # Check the code using various linters and static checkers.
	golangci-lint run --timeout=5m ./...

deps: # Install all dependencies for this project
	go mod download && go mod verify
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

clean:
	rm -rf build

.PHONY: build help test lint deps format clean run

.DEFAULT_GOAL := help
