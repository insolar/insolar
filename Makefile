INSOLAR = insolar
INSOLARD = insolard
INSGOCC = insgocc
PULSARD = pulsard
INSGORUND =insgorund
BIN_DIR =bin

ALL_PACKAGES=./...
COVERPROFILE=coverage.txt

.PHONY: all lint ci-lint metalint clean install-deps install build test test_with_coverage

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
	make $(INSOLARD) $(INSOLAR) $(INSGOCC) $(PULSARD) $(INSGORUND)

$(INSOLARD):
	go build -o $(BIN_DIR)/$(INSOLARD) cmd/insolard/*.go

$(INSOLAR):
	go build -o $(BIN_DIR)/$(INSOLAR) cmd/insolar/*.go

$(INSGOCC):
	go build -o $(BIN_DIR)/$(INSGOCC) cmd/insgocc/*.go

$(PULSARD):
	go build -o $(BIN_DIR)/$(PULSARD) cmd/pulsard/*.go

$(INSGORUND):
	go build -o $(BIN_DIR)/$(INSGORUND) cmd/insgorund/*.go

test:
	go test -v $(ALL_PACKAGES)

test_with_coverage:
	CGO_ENABLED=1 go test --race --coverprofile=$(COVERPROFILE) --covermode=atomic $(ALL_PACKAGES)
