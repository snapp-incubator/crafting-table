BASH_PATH:=$(shell which bash)
SHELL=$(BASH_PATH)
ROOT := $(shell realpath $(dir $(lastword $(MAKEFILE_LIST))))
APP := crafting-table

check-goimport:
	which goimports || GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

format:
	- find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R goimports -w R
	- find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R gofmt -s

vendor:
	- go mod tidy
	- go mod vendor

build: format vendor
	- go build

it: build
	- sh integration_test.sh $(ROOT)


querybuilder-example:
	- go run main.go qb example/querybuilder.go