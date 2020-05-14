.PHONY: all dep lint test tidy coverage tidy dl

.DEFAULT_GOAL := test

all: dep lint test coverage

dep: gen ## Get the dependencies
	@go get -v -d ./...

lint: ## Lint and security checks
	@golangci-lint run

test: ## Run unittests
	@go test ./... -v -tags=test

race: dep ## Run data race detector
	@go test ./... -race -tags=test

coverage: ## Generate global code coverage report
	@go test ./... -cover -tags=test

tidy: ## Remove previous built binnay and keep latest lean packages
	@go mod tidy

update: ## Update dependencies version
	@go get -u ./...

dl: ## Download golang related command line tools for building pipeline
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOPATH}/bin latest
