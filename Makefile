.PHONY: help build-terminal build-base golangci-lint init.golangci-lint test goimports

.DEFAULT_GOAL := help

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test: ## Run Go code tests
	go test -v -shuffle=on -covermode=atomic

goimports: ## Run goimports
	goimports -w .
