BUILD_VERSION ?= $(shell git describe --abbrev=0 --tags)
BIN_DIR ?= bin/${BUILD_VERSION}
INSOLAR = insolar
INSOLARD = insolard
INSGOCC = $(BIN_DIR)/insgocc
PULSARD = pulsard
INSGORUND = insgorund
UPDATESERV = updateserv
LOADANALYZER = loadanalyzer

ALL_PACKAGES = ./...
COVERPROFILE = coverage.txt

BUILD_NUMBER := $(TRAVIS_BUILD_NUMBER)
BUILD_DATE = $(shell date "+%Y-%m-%d")
BUILD_TIME = $(shell date "+%H:%M:%S")
BUILD_HASH = $(shell git rev-parse --short HEAD)

LDFLAGS += -X github.com/insolar/insolar/version.Version=${BUILD_VERSION}
LDFLAGS += -X github.com/insolar/insolar/version.BuildNumber=${BUILD_NUMBER}
LDFLAGS += -X github.com/insolar/insolar/version.BuildDate=${BUILD_DATE}
LDFLAGS += -X github.com/insolar/insolar/version.BuildTime=${BUILD_TIME}
LDFLAGS += -X github.com/insolar/insolar/version.GitHash=${BUILD_HASH}

.PHONY: all lint ci-lint metalint clean install-deps pre-build build test test_with_coverage regen-proxies

all: clean install-deps pre-build build test

lint: ci-lint

ci-lint:
	golangci-lint run $(ALL_PACKAGES)

metalint:
	gometalinter --vendor $(ALL_PACKAGES)

clean:
	go clean $(ALL_PACKAGES)
	rm -f $(COVERPROFILE)
	rm -rf $(BIN_DIR) 

install-deps:
	go get -u github.com/golang/dep/cmd/dep
	go get -u golang.org/x/tools/cmd/stringer

pre-build:
	dep ensure
	go generate -x $(ALL_PACKAGES)

build: 
	mkdir -p $(BIN_DIR)
	make $(INSOLARD) $(INSOLAR) $(INSGOCC) $(PULSARD) $(INSGORUND) $(UPDATESERV)

$(INSOLARD):
	go build -o $(BIN_DIR)/$(INSOLARD) -ldflags "${LDFLAGS}" cmd/insolard/*.go

$(INSOLAR):
	go build -o $(BIN_DIR)/$(INSOLAR) -ldflags "${LDFLAGS}" cmd/insolar/*.go

$(INSGOCC): cmd/insgocc/insgocc.go logicrunner/goplugin/preprocessor
	go build -o $(INSGOCC) -ldflags "${LDFLAGS}" cmd/insgocc/*.go

$(PULSARD):
	go build -o $(BIN_DIR)/$(PULSARD) -ldflags "${LDFLAGS}" cmd/pulsard/*.go

$(UPDATESERV):
	go build -o $(BIN_DIR)/$(UPDATESERV) -ldflags "${LDFLAGS}" cmd/updateserv/*.go

$(INSGORUND):
	go build -o $(BIN_DIR)/$(INSGORUND) -ldflags "${LDFLAGS}" cmd/insgorund/*.go

$(LOADANALYZER):
	go build -o $(BIN_DIR)/$(LOADANALYZER) -ldflags "${LDFLAGS}" cmd/loadanalyzer/*.go

test:
	go test -v $(ALL_PACKAGES)

test_with_coverage:
	CGO_ENABLED=1 go test --coverprofile=$(COVERPROFILE) --covermode=atomic $(ALL_PACKAGES)


CONTRACTS = $(wildcard application/contract/*)
regen-proxies: $(INSGOCC)
	$(foreach c,$(CONTRACTS), $(INSGOCC) proxy application/contract/$(notdir $(c))/$(notdir $(c)).go; )

docker-insolard:
	docker build --tag insolar/insolard -f ./docker/Dockerfile.insolard .

docker-pulsar:
	docker build --tag insolar/pulsar -f ./docker/Dockerfile.pulsar .

docker-insgorund:
	docker build --tag insolar/insgorund -f ./docker/Dockerfile.insgorund .


docker: docker-insolard docker-pulsar docker-insgorund

