
.PHONY: install-linter
install-linter:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: lint
lint:
	@golangci-lint run --config=.golangci.yml

.PHONY: test
test:
	@go test -count=1 -v -race ./...

.PHONY: proto
proto:
	@rm -rf source/shared/proto
	@mkdir source/shared/proto
	@docker run -v $(PWD):/defs namely/protoc-all:1.51_2 -i proto -d proto -o go -l go && \
    mv go/github.com/pipelane/pipelaner/source/shared/proto source/shared  && \
    rm -rf go

.PHONY: install-pkl-go
install-pkl-go:
	go install github.com/apple/pkl-go/cmd/pkl-gen-go@v0.8.1

.PHONY: pkl-generate-go
pkl-generate-go:
	pkl-gen-go pkl/Pipelaner.pkl

.PHONY: pkl-build
pkl-project:
	pkl project package pkl