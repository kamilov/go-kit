.PHONY: default
default: help

.PHONY: help
help: ## help information
	@grep -E '^[a-zA-Z_\/-]+:.*?## .*$$' $(FILE) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: test
test: ## run tests
	@go test -cover -covermode=count -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

.PHONY: lint
lint: ## run linter
	@golangci-lint run ./... -c .golangci.yaml -v