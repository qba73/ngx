.PHONY: dox test vet check cover tidy

help: ## Show help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\033[36m\033[0m\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-24s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

dox: ## Run tests with gotestdox
	@gotestdox  >/dev/null 2>&1 || go install github.com/bitfield/gotestdox/cmd/gotestdox@latest ;
	gotestdox

check: ## Run staticcheck analyzer
	@staticcheck -version >/dev/null 2>&1 || go install honnef.co/go/tools/cmd/staticcheck@2022.1;
	staticcheck ./...

test: ## Run tests
	go test -race -shuffle=on ./...

vet: ## Run go vet
	go vet ./...

cover: ## Run unit tests and generate test coverage report
	go test -race -v ./... -count=1 -cover -covermode=atomic -coverprofile=coverage.out
	go tool cover -html coverage.out

tidy: ## Run go mod tidy
	go mod tidy

