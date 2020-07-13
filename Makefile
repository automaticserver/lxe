# Shared Makefile for go projects

# This Make should be usable with all go code related projects and is setup as generic as possible. There are some
# hints about this file:
# - Lowercase variables, or variables with lowercase suffixes, are not intended to be used directly. They're only used
#   to build the resulting variable and can change any time!
# - Meta targets must always be at the top of each section to rapidly see where dependencies are added

# CONFIG
# Some configuration variables used in various places

VERSION_REGEX = ^v[0-9]+\.[0-9]+\.[0-9]+


# COMMON
# Global defined shorthand variables for various common fields

DATE_ISO_8601 = $(shell date -Is)
DATE_RFC_2822 = $(shell date -R)
BINARIES =


# GIT STATE
# Some generic git state info all targets can use

GIT_HASH_SHORT = $(or $(CI_COMMIT_SHORT_SHA), $(shell git describe --always))
GIT_HASH_LONG = $(or $(CI_COMMIT_SHA), $(shell git rev-parse $(GIT_HASH_SHORT)))
GIT_TREESTATE = clean
GIT_UNSTAGED = $(shell git diff --quiet >/dev/null 2>&1; [ $$? -eq 1 ] && echo "1")
GIT_STAGED= $(shell git diff --cached --quiet >/dev/null 2>&1; [ $$? -eq 1 ] && echo "1")
GIT_UNTRACKED = $(shell git status --porcelain 2>/dev/null | grep "^??" >/dev/null; [ $$? -eq 0 ] && echo "1")
# I don't get how I can make an IF with a logical OR
ifeq ($(GIT_UNSTAGED), 1)
  GIT_TREESTATE = dirty
endif
ifeq ($(GIT_STAGED), 1)
  GIT_TREESTATE = dirty
endif
ifeq ($(GIT_UNTRACKED), 1)
  GIT_TREESTATE = dirty
endif


# VERSION INFO
# Generic version data based on the pushed tag and git state

# A software version is always only the version number, which matches the tag it's build from, otherwise it's always
# 0.0.0. Suffixes like "-rc1" must be allowed for version tags. Set CI_COMMIT_TAG to provide a tagged version.
TAG_VERSION_valid := $(shell bash -c '[[ "$(CI_COMMIT_TAG)" =~ $(VERSION_REGEX) ]] && echo "$(CI_COMMIT_TAG)" | cut -c 2-')
VERSION := $(or $(TAG_VERSION_valid), 0.0.0)
BUILD_NUMBER := $(or $(CI_PIPELINE_IID), undef)
BUILD_DATE := $(DATE_RFC_2822)


##@ General
# show help info about this Makefile

# The help will print out all targets with their descriptions organized bellow their categories. The categories are
# represented by `##@` and the target descriptions by `##`. The awk commands is responsable to read the entire set of
# makefiles included in this invocation, looking for lines of the file as xyz: ## something, and then pretty-format the
# target and help. Then, if there's a line with ##@ something, that gets pretty-printed as a category.
.PHONY: help
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# don't add a body in these, sort alphabetically

.PHONY: info
info: ## Show informations related to the project. Version variants, git state, code/build/package info, ...
.PHONY: build
build: ## Build all code artifacts
.PHONY: install
install: ## Install all build artifacts
.PHONY: generate
generate: ## Update/Generate resources over the whole project
.PHONY: test
test: ## Local tests over the whole project
.PHONY: lint
lint: ## Check/Lint over the whole project
.PHONY: clean
clean: ## Clean all artifacts

# already link to version common target
info: info-common
info-common:
	@echo PROJECT: "$(PROJECT)"
	@echo VERSION: "$(VERSION)"
	@echo GIT_HASH_SHORT: "$(GIT_HASH_SHORT)"
	@echo GIT_HASH_LONG: "$(GIT_HASH_LONG)"
	@echo GIT_TREESTATE: "$(GIT_TREESTATE)"
	@echo BUILD_NUMBER: "$(BUILD_NUMBER)"
	@echo BUILD_DATE: "$(BUILD_DATE)"
	@echo BINARIES: "$(BINARIES)"

.PHONY: get
get: nothing ## Print provided make variables space-delimited. Care, variables may output nothing or multiple entries.
	@echo $(foreach arg, $(ARGS), $($(arg)))

.PHONY: do
do: ## Execute command as it were part of Makefile (for debugging)
	@>&2 echo "$(ARGS)"
	@$(ARGS)

# Nothing is faking a target which always needs to be done so we don't have the "nothing to be done" output
.PHONY: nothing
nothing:
	@:

# For these targets remember all further arguments and override the bodies of thee matching arguments with nothing
ARGK = get run do
ifneq ($(filter $(firstword $(MAKECMDGOALS)),$(ARGK)),)
  ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(ARGS): nothing;@:)
endif


##@ Go
# all targets related to go code and it's environment

info: info-go
build: build-go
clean: clean-go
install: install-go
generate: generate-go
test: test-go-unit
lint: lint-go

export GO111MODULE = on

# by not using `go list` we can save time not waiting for all modules to download and allow other make targets to use the project name without having go
GO_MODULENAME = $(shell grep '^module' go.mod | cut -d' ' -f2)
PROJECT = $(notdir $(GO_MODULENAME))
PACKAGE_NAME ?= $(PROJECT)

LDFLAG_OPTIONS = -X "$(GO_MODULENAME)/cli/version.Version=$(VERSION)" \
                 -X "$(GO_MODULENAME)/cli/version.GitCommit=$(GIT_HASH_LONG)" \
                 -X "$(GO_MODULENAME)/cli/version.GitTreeState=$(GIT_TREESTATE)" \
                 -X "$(GO_MODULENAME)/cli/version.BuildNumber=$(BUILD_NUMBER)" \
                 -X "$(GO_MODULENAME)/cli/version.BuildDate=$(BUILD_DATE)" \
                 -X "$(GO_MODULENAME)/cli/version.PackageName=$(PACKAGE_NAME)"
ldflags = all='$(LDFLAG_OPTIONS) -s -w'

GO_BINARIES := $(patsubst cmd/%,bin/%, $(wildcard cmd/*))
BINARIES := $(BINARIES) $(GO_BINARIES)
# $(wildcard) is not good enough for recursive file lookup
GO_SOURCES := go.mod $(shell find * -type f -name '*.go' ! -name '*_test.go' ! -name '.*' ! -wholename 'vendor/*')

info-go:
	@echo GO_MODULENAME: "$(GO_MODULENAME)"
#	@echo GO_SOURCES: "$(GO_SOURCES)"

# all cmd/* packages can be build to bin/*
bin/%: cmd/% $(GO_SOURCES)
	go build -o $@ -ldflags=$(ldflags) $(GO_MODULENAME)/$<

install-go: build-go
	mkdir -p $(DESTDIR)/usr/bin
	cp bin/* $(DESTDIR)/usr/bin/

build-go: $(GO_BINARIES) ## Build the binaries

generate-go: dep-go ## Update/Generate code over the whole code source
	go generate ./...
	@echo "$(GO_SOURCES)" | tr ' ' '\n' | fixtures/go/packageimportcomments.sh

lint-go: ## Lint the source code
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run

test-go-unit: ## Unit-Test the code
	go test ./... -coverprofile go.coverprofile
	go tool cover -func go.coverprofile

test-go-race: ## Race-Test the code
	go test -race ./...

test-go-msan: ## Memory-Test the code
	go test -msan ./...

dep-go: ## Download code dependencies
	go mod download

vendor-go: dep-go go.mod go.sum ## Export go modules to vendor
	go mod vendor

clean-go: clean-vendor-go
	@rm -r bin || true

clean-vendor-go:
	@rm -r vendor || true

# Run executes the go program directly. First argument is which cmd, add further arguments after a '--'
run: ## Execute program directly. `make run <cmd> [-- <arguments>]`
# The first word is the cmd go package, then all assignment-args, then the remaining words
	$(eval cmd := $(word 1, $(ARGS)) $(MAKEOVERRIDES) $(wordlist 2,$(words $(ARGS)),$(ARGS)))
# we only want output to stdout from the program
	$(eval line := go run -ldflags=$(ldflags) $(GO_MODULENAME)/cmd/$(cmd))
	@>&2 echo "$(line)"
	@$(line)
