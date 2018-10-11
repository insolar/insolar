BIN_DIR = bin
INSOLAR = insolar
INSOLARD = insolard
INSGOCC = $(BIN_DIR)/insgocc
PULSARD = pulsard
INSGORUND = insgorund
INSUPDATER = insupdater
INSUPDATESERV = updateserv


ALL_PACKAGES = ./...
COVERPROFILE = coverage.txt

BUILD_NUMBER := $(TRAVIS_BUILD_NUMBER)
BUILD_DATE = $(shell date "+%Y-%m-%d")
BUILD_TIME = $(shell date "+%H:%M:%S")
BUILD_HASH = $(shell git rev-parse --short HEAD)
BUILD_VERSION ?= $(shell git describe --abbrev=0 --tags)

LDFLAGS += -X github.com/insolar/insolar/version.Version=${BUILD_VERSION}
LDFLAGS += -X github.com/insolar/insolar/version.BuildNumber=${BUILD_NUMBER}
LDFLAGS += -X github.com/insolar/insolar/version.BuildDate=${BUILD_DATE}
LDFLAGS += -X github.com/insolar/insolar/version.BuildTime=${BUILD_TIME}
LDFLAGS += -X github.com/insolar/insolar/version.GitHash=${BUILD_HASH}

.PHONY: all lint ci-lint metalint clean install-deps install build test test_with_coverage regen_proxyes

all: clean install-deps install build test

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
	make $(INSOLARD) $(INSOLAR) $(INSGOCC) $(PULSARD) $(INSGORUND) $(INSUPDATER) $(INSUPDATESERV)

$(INSOLARD):
	go build -o $(BIN_DIR)/$(INSOLARD) -ldflags "${LDFLAGS}" cmd/insolard/*.go

$(INSOLAR):
	go build -o $(BIN_DIR)/$(INSOLAR) -ldflags "${LDFLAGS}" cmd/insolar/*.go

$(INSGOCC):
	go build -o $(INSGOCC) -ldflags "${LDFLAGS}" cmd/insgocc/*.go

$(PULSARD):
	go build -o $(BIN_DIR)/$(PULSARD) -ldflags "${LDFLAGS}" cmd/pulsard/*.go

$(INSGORUND):
	go build -o $(BIN_DIR)/$(INSGORUND) -ldflags "${LDFLAGS}" cmd/insgorund/*.go

$(INSUPDATER):
	go build -o $(BIN_DIR)/$(INSUPDATER) -ldflags "${LDFLAGS}" cmd/insupdater/*.go

$(INSUPDATESERV):
	go build -o $(BIN_DIR)/$(INSUPDATESERV) -ldflags "${LDFLAGS}" cmd/updateserv/*.go

test:
	go test -v $(ALL_PACKAGES)

test_with_coverage:
	CGO_ENABLED=1 go test --coverprofile=$(COVERPROFILE) --covermode=atomic $(ALL_PACKAGES)


CONTRACTS = $(wildcard genesis/experiment/*)
regen_proxyes: $(INSGOCC)
	$(foreach c,$(CONTRACTS), $(INSGOCC) proxy genesis/experiment/$(notdir $(c))/$(notdir $(c)).go; )
