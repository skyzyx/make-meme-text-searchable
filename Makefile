#-------------------------------------------------------------------------------
# Running `make` will show the list of subcommands that will run.

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(dir $(mkfile_path))

#-------------------------------------------------------------------------------
# Global stuff.

GO=$(shell which go)

# Determine which version of `echo` to use. Use version from coreutils if available.
ECHOCHECK := $(shell command -v /usr/local/opt/coreutils/libexec/gnubin/echo 2> /dev/null)
ifdef ECHOCHECK
    ECHO=/usr/local/opt/coreutils/libexec/gnubin/echo
else
    ECHO=echo
endif

#-------------------------------------------------------------------------------
# Running `make` will show the list of subcommands that will run.

all: help

.PHONY: help
## help: [help]* Prints this help message.
help:
	@ $(ECHO) "Usage:"
	@ $(ECHO) ""
	@ sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /' | \
		while IFS= read -r line; do \
			if [[ "$$line" == *"]*"* ]]; then \
				$(ECHO) -e "\033[1;33m$$line\033[0m"; \
			else \
				$(ECHO) "$$line"; \
			fi; \
		done

#-------------------------------------------------------------------------------
# Install

# Private
.PHONY: _install-go-deps
_install-go-deps:
	$(GO) install github.com/quasilyte/go-consistent@latest
	$(GO) install github.com/jgautheron/goconst/cmd/goconst@latest
	$(GO) install mvdan.cc/gofumpt@latest
	$(GO) install github.com/pavius/impi/cmd/impi@latest

.PHONY: install-deps-linux
## install-deps-linux: [deps]* Installs the tools for Linux.
install-deps-linux: _install-go-deps
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin

.PHONY: install-deps-mac
## install-deps-mac: [deps]* Installs the tools for macOS. Requires Homebrew.
install-deps-mac: _install-go-deps
	brew install golangci-lint

#-------------------------------------------------------------------------------
# Build/Run

.PHONY: fetch
## fetch: [build] Fetches remote data that is required for building.
fetch:
	wget -qq https://data.iana.org/TLD/tlds-alpha-by-domain.txt

.PHONY: tidy
## tidy: [build] Updates go.mod and downloads dependencies.
tidy:
	$(GO) mod tidy -go=1.17 -v
	$(GO) mod download -x
	$(GO) get -v ./...

.PHONY: generate
## generate: [build] Reads the raw list and generates a UTF-8 JSON document.
generate:
	go run main.go | tee tlds.json

.PHONY: build
## build: [build]* Compiles the final artifacts.
build: fetch generate clean-files

#-------------------------------------------------------------------------------
# Clean

.PHONY: clean-files
## clean-files: [clean] Clean temporary files.
clean-files:
	rm -f tlds-alpha-by-domain.txt

.PHONY: clean-go
## clean-go: [clean] Clean Go's module cache.
clean-go:
	$(GO) clean -i -r -x -testcache -modcache -cache

.PHONY: clean
## clean: [clean]* Runs ALL non-Golang cleaning tasks.
clean: clean-files

#-------------------------------------------------------------------------------
# Linting

.PHONY: fmt
## fmt: [lint] Runs `gofumpt` against all Golang files.
fmt:
	@ echo " "
	@ echo "=====> Running gofumpt..."
	gofumpt -s -w *.go

.PHONY: golint
## golint: [lint] Runs `golangci-lint` against all Golang files.
golint:
	@ echo " "
	@ echo "=====> Running golangci-lint..."
	golangci-lint run --fix *.go

.PHONY: goupdate
## goupdate: [lint] Runs `go-mod-outdated` to check for out-of-date packages.
goupdate:
	@ echo " "
	@ echo "=====> Running go-mod-outdated..."
	$(GO) list -u -m -json all | go-mod-outdated -update -direct -style markdown

.PHONY: goconsistent
## goconsistent: [lint] Runs `go-consistent` to ensure consistent patterns.
goconsistent:
	@ echo " "
	@ echo "=====> Running go-consistent..."
	- go-consistent -v ./...

.PHONY: goimportorder
## goimportorder: [lint] Runs `impi` to verify that import order is consistent.
goimportorder:
	@ echo " "
	@ echo "=====> Running impi..."
	- impi \
		--local $(shell head -n1 < go.mod | cut -d' ' -f2) \
		--ignore-generated=true \
		--scheme=stdLocalThirdParty \
		./...

.PHONY: goconst
## goconst: [lint] Runs `goconst` to identify opportunities for constants.
goconst:
	@ echo " "
	@ echo "=====> Running goconst..."
	- goconst -match-constant -numbers ./...

.PHONY: markdownlint
## markdownlint: [lint] Runs `markdownlint` against all Markdown documents.
markdownlint:
	@ echo " "
	@ echo "=====> Running Markdownlint..."
	- npx -y markdownlint-cli --fix '*.md' --ignore 'node_modules'

.PHONY: lint
## lint: [lint]* Runs ALL linting/validation tasks.
lint: markdownlint fmt golint goupdate goimportorder goconst
