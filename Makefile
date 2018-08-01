.PHONY: lint test ci-lint

ci-lint:
	golangci-lint run

metalint:
	gometalinter --vendor ./...

test:
	go test -v ./...
