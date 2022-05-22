#!make
include build.env
export $(shell sed 's/=.*//' build.env)
GIT_COMMIT := $(shell git describe --always --long --dirty)
PROJECT_NAME := $(shell basename "$$PWD")

.DEFAULT_GOAL := default

.PHONY: default
default: fmt build run #tests

#
# *** fmt ***
#
.PHONY: fmt
fmt:
	@find . -maxdepth 1 -iname "*go" -exec go fmt {} +
	@for file in `find internal/ -iname "*go"`; do echo $$file; go fmt $$file; done;

#
# *** build steps ***
#
.PHONY: build-executable
build: go-mod build-executable

go-mod:
	@go mod vendor
	@go mod verify

build-executable:
	@rm -f ${EXECUTABLE}
	@go build -o ${EXECUTABLE} -ldflags "-X main.AppVersion=${GIT_COMMIT}" .
	@echo "wrote binary to ${EXECUTABLE}"

#
# *** example runs ****
#
.PHONY: run
run: run-help run-func

run-help:
	@./${EXECUTABLE} --help
	@echo "\n____________________________ \n"

run-func:
	@./${EXECUTABLE} --file challenge/logs.csv --limit 1 --field message scan
	@echo "\n____________________________"
	@./${EXECUTABLE} --file challenge/logs.csv --limit 1 --field message stats --urlpath login
	@echo "\n____________________________"

#
# *** tests ****
#
.PHONY: tests
tests: unit-tests app-tests

.PHONY: unit-tests
unit-tests:
	@go test -cover -failfast -short "./.../types"
	@echo "\n____________________________"
	@go test -cover -failfast -short "./.../app"
	@echo "\n____________________________"

.PHONY: app-tests
app-tests:
	@go test -cover -parallel 1 -failfast -short "."
	@echo "\n____________________________"
