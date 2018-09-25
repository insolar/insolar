INSOLAR =insolar
INSOLARD=insolard
INSGOCC =insgocc
PULSARD =pulsard

ALL_PACKAGES=./...
COVERPROFILE=coverage.txt

BUILD_NUMBER:=$(TRAVIS_BUILD_NUMBER)
BUILD_DATE = $(shell date "+%Y-%m-%d")
BUILD_TIME = $(shell date "+%H:%M:%S")
BUILD_HASH = $(shell git rev-parse --short HEAD)

LDFLAGS += -X github.com/insolar/insolar/version.Version=0.3.0
LDFLAGS += -X github.com/insolar/insolar/version.BuildNumber=${BUILD_NUMBER}
LDFLAGS += -X github.com/insolar/insolar/version.BuildDate=${BUILD_DATE}
LDFLAGS += -X github.com/insolar/insolar/version.BuildTime=${BUILD_TIME}
LDFLAGS += -X github.com/insolar/insolar/version.GitHash=${BUILD_HASH}

.PHONY: all lint ci-lint metalint clean install-deps install build test test_with_coverage

all: clean install-deps install build test

lint: ci-lint

ci-lint:
	golangci-lint run $(ALL_PACKAGES)

metalint:
	gometalinter --vendor $(ALL_PACKAGES)

clean:
	go clean $(ALL_PACKAGES)
	rm -f $(INSOLARD)
	rm -f $(INSOLAR)
	rm -f $(INSGOCC)
	rm -f $(PULSARD)
	rm -f $(COVERPROFILE)

install-deps:
	go get -u github.com/golang/dep/cmd/dep
	go get -u golang.org/x/tools/cmd/stringer

install:
	dep ensure
	go generate -x $(ALL_PACKAGES)

build: $(INSOLARD) $(INSOLAR) $(INSGOCC) $(PULSARD)

$(INSOLARD):
	go build -o $(INSOLARD) -ldflags "${LDFLAGS}" cmd/insolard/*

$(INSOLAR):
	go build -o $(INSOLAR) -ldflags "${LDFLAGS}" cmd/insolar/*

$(INSGOCC):
	go build -o $(INSGOCC) -ldflags "${LDFLAGS}" cmd/insgocc/*

$(PULSARD):
	go build -o $(PULSARD) -ldflags "${LDFLAGS}" cmd/pulsard/*

test:
	go test -v $(ALL_PACKAGES)

test_with_coverage:
	CGO_ENABLED=1 go test --race --coverprofile=$(COVERPROFILE) --covermode=atomic $(ALL_PACKAGES)
