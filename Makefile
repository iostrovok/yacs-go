DIRS := ${sort ${dir $(shell find . | grep "/*\\.go")}}
SUBPACKAGES := ${sort ${dir $(shell find ./go | grep "/*\\.go")}}


# Global GOLANG modules for CI/CD
# TODO: use tools for this
BUILD_DEPS := gopkg.in/check.v1 \
	github.com/golang/lint/golint \
	github.com/xeipuuv/gojsonschema


GOBIN := $(CURDIR)/bin/
EXE_PARSER := $(GOBIN)/yacsgo
EXE_SERVER := $(GOBIN)/server
# GOPATH_ENV := $(GOPATH)
GOPATH_ENV := $(GOPATH):$(CURDIR)/go/

ENV := GOPATH=$(GOPATH_ENV) GOBIN=$(GOBIN)

export GOPATH=$(GOPATH_ENV)


##
## List of commands:
##


## default:
all: clean deps fmt lint test build


linux: clean deps fmt lint test build-linux


# Remove build and vendor directories
clean:
	@echo "======================================================================"
	@echo 'MAKE: clean: yacsgo...'
	@rm -rf $(EXE_PARSER) $(EXE_SERVER)
	@rm -rf ./test_out ./test_tmp_checked


# Installing build dependencies. You will need to run this once manually when you clone the repo
deps:
	@echo "======================================================================"
	@echo 'MAKE: install...'
	@$(ENV) go get -v $(BUILD_DEPS)


# Build exe file and suppoting files
build: clean
	@echo "======================================================================"
	@echo 'MAKE: build...'
	mkdir -p $(GOBIN)
	$(ENV) CGO_ENABLED=0 go build -o $(EXE_PARSER) ./yacsgo.go
	$(ENV) CGO_ENABLED=0 go build -o $(EXE_SERVER) ./server.go

# Build exe file and suppoting files
build-linux: clean
	@echo "======================================================================"
	@echo 'MAKE: build...'
	mkdir -p $(GOBIN)
	$(ENV) CGO_ENABLED=0 GOOS=linux go build -o $(EXE_PARSER) ./yacsgo.go
	$(ENV) CGO_ENABLED=0 GOOS=linux go build -o $(EXE_SERVER) ./yacsgo.go


# Full tests
tests: fmt lint test


test:
	@echo "======================================================================"
	@echo "Run race test for " $(SUBPACKAGES)
	@go test -race $(SUBPACKAGES)
	for dir in $(SUBPACKAGES); do \
        echo "go test " $$dir; \
        $(ENV) go test -cover $$dir/; \
	done


lint:
	@echo "======================================================================"
	@echo "Run golint..."
	for dir in $(DIRS); do \
		echo "golint " $$dir; \
        $(ENV) $(GOBIN)/golint $$dir/; \
	done


fmt:
	@echo "======================================================================"
	for dir in $(DIRS); do \
        echo "go fmt " $$dir; \
        $(ENV) go fmt $$dir/*.go; \
	done
