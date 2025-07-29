
.PHONY: install-linter
install-linter:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: lint
lint:
	@golangci-lint run --config=.golangci.yml

.PHONY: test
test:
	@go test -count=1 -v ./... -race
.PHONY: proto
proto:
	@rm -rf sources/shared/proto
	@mkdir sources/shared/proto
	@docker run -v $(PWD):/defs namely/protoc-all:1.51_2 -i proto -d proto -o go -l go && \
    mv go/github.com/pipelane/pipelaner/sources/shared/proto sources/shared  && \
    rm -rf go

.PHONY: install-pkl-go
install-pkl-go:
	go install github.com/apple/pkl-go/cmd/pkl-gen-go@v0.11.0

.PHONY: pkl-generate-go
pkl-generate-go:
	rm -rf ./gen
	pkl-gen-go pkl/Pipelaner.pkl

.PHONY: pkl-build
pkl-project:
	pkl project package pkl