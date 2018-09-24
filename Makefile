INSOLAR =insolar
INSOLARD=insolard
INSGOCC =insgocc
PULSARD =pulsard
GINSIDER=ginsider-cli

ALL_PACKAGES=./...
COVERPROFILE=coverage.txt

.PHONY: all lint ci-lint metalint clean install-deps install build test test_with_coverage

all: clean install-deps install build test

lint: ci-lint

ci-lint:
	golangci-lint run

metalint:
	gometalinter --vendor $(ALL_PACKAGES)

clean:
	go clean $(ALL_PACKAGES)
	rm -f $(INSOLARD)
	rm -f $(INSOLAR)
	rm -f $(INSGOCC)
	rm -f $(PULSARD)
	rm -f $(GINSIDER)
	rm -f $(COVERPROFILE)

install-deps:
	go get -u github.com/golang/dep/cmd/dep
	go get -u golang.org/x/tools/cmd/stringer

install:
	dep ensure
	go generate -x $(ALL_PACKAGES)

build:
	go build -o $(INSOLARD) cmd/insolard/*
	go build -o $(INSOLAR) cmd/insolar/*
	go build -o $(INSGOCC) cmd/insgocc/*
	go build -o $(PULSARD) cmd/pulsard/*
	go build -o $(GINSIDER) logicrunner/goplugin/ginsider-cli/main.go

test:
	go test -v $(ALL_PACKAGES)

test_with_coverage:
	CGO_ENABLED=1 go test --race --coverprofile=$(COVERPROFILE) --covermode=atomic $(ALL_PACKAGES)