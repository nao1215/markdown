.PHONY: build test clean

APP         = markdown
VERSION     = $(shell git describe --tags --abbrev=0)
GIT_REVISION := $(shell git rev-parse HEAD)
GO          = go
GO_BUILD    = $(GO) build
GO_TEST     = $(GO) test -v
GO_TOOL     = $(GO) tool
GOOS        = ""
GOARCH      = ""
GO_PKGROOT  = ./...
GO_PACKAGES = $(shell $(GO_LIST) $(GO_PKGROOT))
GO_LDFLAGS  = 

clean: ## Clean project
	-rm -rf $(APP) coverage.out coverage.html

test: ## Start unit test for server
	env GOOS=$(GOOS) $(GO_TEST) -cover -coverpkg=$(GO_PKGROOT) -coverprofile=coverage.out $(GO_PKGROOT)
	$(GO_TOOL) cover -html=coverage.out -o coverage.html
	
.DEFAULT_GOAL := help
help: ## Show help message  
	@grep -E '^[0-9a-zA-Z_-]+[[:blank:]]*:.*?## .*$$' $(MAKEFILE_LIST) | sort \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1;32m%-15s\033[0m %s\n", $$1, $$2}'
