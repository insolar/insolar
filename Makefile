.PHONY: lint test ci-lint

ci-lint:
	golangci-lint --tests=0 run

metalint:
	gometalinter --vendor ./...

test:
	go test -v ./...
