ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

lint:
	which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.0)
	golangci-lint run --config=$(ROOT)/.golangci.yml $(ROOT)/...

format:
	@which gofumpt || (go install mvdan.cc/gofumpt@latest)
	@gofumpt -l -w $(ROOT)
	@which gci || (go install github.com/daixiang0/gci@latest)
	@gci write $(ROOT)
	@which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.0)
	@golangci-lint run --fix

test:
	go test ./...