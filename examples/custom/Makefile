.PHONY: install-pkl-go
install-pkl-go:
	go install github.com/apple/pkl-go/cmd/pkl-gen-go@v0.11.0

.PHONY: example-generate
pkl-generate-go:
	pkl-gen-go pkl/custom.pkl
	rm -rf github.com/