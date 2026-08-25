.PHONY: adapter-test fmt fmt-check test test-race vet verify verify-all

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

fmt-check:
	test -z "$$(gofmt -l $$(find . -name '*.go' -not -path './vendor/*'))"

test:
	go test ./...

test-race:
	go test -race ./...

vet:
	go vet ./...

verify: fmt-check vet test test-race

adapter-test:
	cd adapters/goconfig && go test -race ./...

verify-all: verify adapter-test
