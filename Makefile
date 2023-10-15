.PHONY: install-linter
install-linter:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: lint
lint:
	@golangci-lint run --enable=gocritic,gocyclo,gofmt,gosec,misspell,unparam,asciicheck --timeout=30m

.PHONY: test
test:
	@go test -v ./