.PHONY: install-linter
install-linter:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: lint
lint:
	@golangci-lint run --config=.golangci.yml

.PHONY: test
test:
	@go test -count=1 -v ./

.PHONY: proto
proto:
	@rm -rf source/shared/proto
	@mkdir source/shared/proto
	@docker run -v $(PWD):/defs namely/protoc-all:1.51_2 -i proto -d proto -o go -l go && \
    mv go/github.com/pipelane/pipelaner/source/shared/proto source/shared  && \
    rm -rf go
