# Makefile for kratos app

GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)

BUILD_TIME:=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
API_VERSION := v1
ROOT:=$(shell go list -m -f '{{.Dir}}')
GO_MODULE_NAME:=$(shell go list -m)
GO_PKG_NAME:=$(shell go list)
SERVICE_NAME:=$(subst $(GO_MODULE_NAME)/app/,,$(GO_PKG_NAME))
API_PATH=$(ROOT)/api/$(SERVICE_NAME)/$(API_VERSION)
API_PROTO_FILES=$(shell find $(API_PATH) -name *.proto)

VERSION_PKG_NAME=$(GO_PKG_NAME)/internal/version

# migration
MIGRATION_PATH = ./migrations
MIGRATE_CMD = migrate

# service version
ifeq ($(SERVICE_NAME),)
	VERSION=$(shell git describe --tags --always --dirty)
else
	MODULE_NAME=$(subst /,-,$(SERVICE_NAME))
	VERSION=$(shell git describe --tags --always --dirty --match "app/$(MODULE_NAME)/*")
endif

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
endif

.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal \
	       --proto_path=$(ROOT)/third_party \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: build
# build
build:
	mkdir -p bin/ && \
		go build \
		-ldflags \
		"-X $(VERSION_PKG_NAME).Version=$(VERSION) -X $(VERSION_PKG_NAME).BuildTime=$(BUILD_TIME)" \
		-o ./bin/ ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: generate
# generate
generate:
	go generate ./...
	go mod tidy

.PHONY: all
# generate all
all:
	make config;
	make generate;

.PHONY: migration
migration:
	@if [ -z "$(name)" ]; then \
		echo "❌ 请指定 name 参数，如: make migrate-create name=create_user"; \
		exit 1; \
	fi
	@migrate create -ext sql -dir $(MIGRATION_PATH) -seq "$(name)"
	@echo "✅ Migration created: $(name)"

.PHONY: api
# generate api proto
api:
	protoc --proto_path=$(ROOT)/api \
	       --proto_path=$(ROOT)/third_party \
 	       --go_out=paths=source_relative:$(ROOT)/api \
 	       --go-http_out=paths=source_relative:$(ROOT)/api \
 	       --go-grpc_out=paths=source_relative:$(ROOT)/api \
		   --go-errors_out=paths=source_relative:$(ROOT)/api \
		   --validate_out=paths=source_relative,lang=go:$(ROOT)/api \
	       --openapi_out=fq_schema_naming=true,default_response=false,version=$(VERSION):$(API_PATH) \
	       $(API_PROTO_FILES)

.PHONY: debug-vars
debug-vars:
	@echo "--- Makefile Variables Debug ---"
	@echo "ROOT:            $(ROOT)"
	@echo "GO_MODULE_NAME:  $(GO_MODULE_NAME)"
	@echo "GO_PKG_NAME:     $(GO_PKG_NAME)"
	@echo "SERVICE_NAME:    $(SERVICE_NAME)"
	@echo "API_PATH:        $(API_PATH)"
	@echo "API_PROTO_FILES: $(API_PROTO_FILES)"
	@echo "--------------------------------"

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
