.PHONY: metalint ci-lint test

ci-lint:
	golangci-lint --tests=0 run

metalint:
	gometalinter --vendor ./...

test:
	go test -v ./...
