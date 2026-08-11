.PHONY: build test lint clean help changelog tools generate render-check

APP         = markdown
VERSION     = $(shell git describe --tags --abbrev=0)
GIT_REVISION := $(shell git rev-parse HEAD)
GO          = go
GO_BUILD    = $(GO) build
GO_TEST     = $(GO) test -v
GO_LIST     = $(GO) list
GO_TOOL     = $(GO) tool
GO_INSTALL  = $(GO) install
GOOS        = ""
GOARCH      = ""
GO_PKGROOT  = ./...
GO_PACKAGES = $(shell $(GO_LIST) $(GO_PKGROOT))
GO_LDFLAGS  = 

clean: ## Clean project
	-rm -rf $(APP) coverage.out coverage.html

test: ## Start unit test for server
	env GOOS=$(GOOS) $(GO_TEST) -cover -coverpkg=$(shell $(GO_LIST) ./... | grep -v '/doc/' | paste -sd,) -coverprofile=coverage.out $(shell $(GO_LIST) ./... | grep -v '/doc/')
	$(GO_TOOL) cover -html=coverage.out -o coverage.html

lint: ## Run linter
	golangci-lint run

generate: ## Regenerate the sample documents under doc/ (CheckAutoGenerateFiles verifies these)
	$(GO) generate ./...

render-check: ## Render every mermaid diagram committed in this repository (requires node)
	cd scripts/mermaid-check && npm ci
	node scripts/mermaid-check/selftest.mjs
	git ls-files -z '*.md' | node scripts/mermaid-check/check.mjs --stdin0

.DEFAULT_GOAL := help
help: ## Show help message  
	@grep -E '^[0-9a-zA-Z_-]+[[:blank:]]*:.*?## .*$$' $(MAKEFILE_LIST) | sort \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1;32m%-15s\033[0m %s\n", $$1, $$2}'

changelog: ## Generate changelog
	ghch --format markdown > CHANGELOG.md

tools: ## Install dependency tools 
	$(GO_INSTALL) github.com/Songmu/ghch/cmd/ghch@latest
