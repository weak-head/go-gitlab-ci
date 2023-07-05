PROJECT_NAME := "gogin"
BUILD_FOLDER := "build"
PKG := "git.lothric.net/examples/go/gogin"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: dep swagger lint test race msan coverage coverhtml build clean help

dep: ## Install dependencies
	@go get -v -d ./...
	@go get -u golang.org/x/lint/golint
	@go install github.com/swaggo/swag/cmd/swag@latest

swagger: dep ## Generate OpenAPI documentation
	@swag fmt
	@swag init -q \
		-g ./cmd/gogin/main.go \
		-o ./api

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}

race: dep ## Run data race detector
	@go test -race -short ${PKG_LIST}

msan: dep ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

coverage: ## Generate global code coverage report
	@./scripts/coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	@./scripts/coverage.sh html;

build: dep ## Build the binary file
	@go build -o ./"${BUILD_FOLDER}"/"${PROJECT_NAME}" -v ./cmd/"${PROJECT_NAME}"

clean: ## Remove previous build and coverage reports
	@rm -rf "${BUILD_FOLDER}"
	@rm -rf "coverage"

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
