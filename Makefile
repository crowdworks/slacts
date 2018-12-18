.DEFAULT_GOAL := help

# https://gist.github.com/tadashi-aikawa/da73d277a3c1ec6767ed48d1335900f3
.PHONY: $(shell grep -E '^[a-zA-Z_-]+:' $(MAKEFILE_LIST) | sed 's/://')

install: ## install slacts command
	go install github.com/crowdworks/slacts/cmd/slacts

test: ## run tests for all package
	go test -v ./...

coverage: ## measure tests coverage and generate coverage profile html
	go test -coverprofile coverprofile -v ./...
	go tool cover -html coverprofile -o coverprofile.html

govet: ## exec go vet checks all for package
	go vet ./...

golint: ## exec golint checks all for package
	golint ./...

goimports: ## exec goimports checks all for package
	goimports -l .

goimports-fix: ## fix format by goimports
	goimports -w .

lint: lint-docker lint-golang ## lint all

lint-docker: ## lint Dockerfile
	hadolint Dockerfile

lint-golang: ## lint golang code
	golangci-lint run ./...

# https://postd.cc/auto-documented-makefile/
help: ## show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
