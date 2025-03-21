SHELL := /bin/bash

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

generate: fmt lint staticcheck md-lint

fmt: gci
	go mod tidy
	go fmt ./... && go vet ./...
	find . -type f -name '*.go' -a ! -name '*zz_generated*' -exec $(GCI) write -s standard -s default -s "prefix(github.com/fra98/pokedex)" {} \;

lint: golangci-lint
	$(GOLANGCILINT) run --new=false

# Run static check anaylisis tools.
# - nilaway: static analysis tool to detect potential Nil panics in Go code
staticcheck: nilaway
	$(NILAWAY) -include-pkgs github.com/fra98/pokedex ./...

md-lint: markdownlint
	@find . -type f -name '*.md' -a -not -path "./.github/*" \
		-exec $(MARKDOWNLINT) {} +

# Installers:

# Install gci if not available
gci:
ifeq (, $(shell which gci))
	@go install github.com/daixiang0/gci@v0.13.6
GCI=$(GOBIN)/gci
else
GCI=$(shell which gci)
endif

# Install golangci-lint if not available
golangci-lint:
ifeq (, $(shell which golangci-lint))
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8
GOLANGCILINT=$(GOBIN)/golangci-lint
else
GOLANGCILINT=$(shell which golangci-lint)
endif

# Install markdownlint if not available
markdownlint:
ifeq (, $(shell which markdownlint))
	@echo "markdownlint is not installed. Please install it: https://github.com/igorshubovych/markdownlint-cli#installation"
	@exit 1
else
MARKDOWNLINT=$(shell which markdownlint)
endif

# Install nilaway if not available
nilaway:
ifeq (, $(shell which nilaway))
	@go install go.uber.org/nilaway/cmd/nilaway@latest
NILAWAY=$(GOBIN)/nilaway
else
NILAWAY=$(shell which nilaway)
endif


