.PHONY: fmt fmt-check test test-race vet verify

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
