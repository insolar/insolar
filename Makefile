BIN_DIR ?= bin
ARTIFACTS_DIR ?= .artifacts
INSOLAR = insolar
INSOLARD = insolard
INSGOCC = insgocc
PULSARD = pulsard
INSGORUND = insgorund
BENCHMARK = benchmark
PULSEWATCHER = pulsewatcher
APIREQUESTER = apirequester
HEALTHCHECK = healthcheck
RECORDBUILDER = protoc-gen-gorecord

ALL_PACKAGES = ./...
MOCKS_PACKAGE = github.com/insolar/insolar/testutils
TESTED_PACKAGES ?= $(shell go list ${ALL_PACKAGES} | grep -v "${MOCKS_PACKAGE}")
COVERPROFILE ?= coverage.txt
TEST_ARGS ?= -timeout 1200s
BUILD_TAGS ?=

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

BININSGOCC=$(BIN_DIR)/$(INSGOCC)


.PHONY: all
all: clean install-deps pre-build build

.PHONY: lint
lint: ci-lint

.PHONY: ci-lint
ci-lint:
	golangci-lint run --new-from-rev=c8f94b7f41b9ae0d2b7ed618d37358b78f479bee

.PHONY: metalint
metalint:
	gometalinter --vendor $(ALL_PACKAGES)

.PHONY: clean
clean:
	go clean $(ALL_PACKAGES)
	rm -f $(COVERPROFILE)
	rm -rf $(BIN_DIR)
	./scripts/insolard/launchnet.sh -l


.PHONY: install-godep
install-godep:
	./scripts/build/fetchdeps github.com/golang/dep/cmd/dep 22125cfaa6ddc71e145b1535d4b7ee9744fefff2

.PHONY: install-build-tools
install-build-tools:
	go clean -modcache
	./scripts/build/fetchdeps golang.org/x/tools/cmd/stringer 63e6ed9258fa6cbc90aab9b1eef3e0866e89b874
	./scripts/build/fetchdeps github.com/gojuno/minimock/cmd/minimock 890c67cef23dd06d694294d4f7b1026ed7bac8e6
	./scripts/build/fetchdeps github.com/gogo/protobuf/protoc-gen-gogoslick v1.2.1

.PHONY: install-deps
install-deps: install-godep install-build-tools

.PHONY: pre-build
pre-build: ensure generate

.PHONY: generate
generate:
	GOPATH=`go env GOPATH` go generate -x $(ALL_PACKAGES)

.PHONY: test_git_no_changes
test_git_no_changes:
	ci/scripts/git_diff_without_comments.sh

.PHONY: ensure
ensure:
	dep ensure

.PHONY: build
build: $(BIN_DIR) $(INSOLARD) $(INSOLAR) $(INSGOCC) $(PULSARD) $(INSGORUND) $(HEALTHCHECK) $(BENCHMARK) $(APIREQUESTER) $(PULSEWATCHER)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

.PHONY: $(INSOLARD)
$(INSOLARD):
	go build -o $(BIN_DIR)/$(INSOLARD) ${BUILD_TAGS} -ldflags "${LDFLAGS}" cmd/insolard/*.go

.PHONY: $(INSOLAR)
$(INSOLAR):
	go build -o $(BIN_DIR)/$(INSOLAR) ${BUILD_TAGS} -ldflags "${LDFLAGS}" cmd/insolar/*.go

.PHONY: $(INSGOCC)
$(INSGOCC): cmd/insgocc/insgocc.go logicrunner/preprocessor
	go build -o $(BININSGOCC) -ldflags "${LDFLAGS}" cmd/insgocc/*.go

$(BININSGOCC): $(INSGOCC)

.PHONY: $(PULSARD)
$(PULSARD):
	go build -o $(BIN_DIR)/$(PULSARD) -ldflags "${LDFLAGS}" cmd/pulsard/*.go

.PHONY: $(INSGORUND)
$(INSGORUND):
	CGO_ENABLED=1 go build -o $(BIN_DIR)/$(INSGORUND) -ldflags "${LDFLAGS}" cmd/insgorund/*.go

.PHONY: $(BENCHMARK)
$(BENCHMARK):
	go build -o $(BIN_DIR)/$(BENCHMARK) -ldflags "${LDFLAGS}" cmd/benchmark/*.go

.PHONY: $(PULSEWATCHER)
$(PULSEWATCHER):
	go build -o $(BIN_DIR)/$(PULSEWATCHER) -ldflags "${LDFLAGS}" cmd/pulsewatcher/*.go

.PHONY: $(APIREQUESTER)
$(APIREQUESTER):
	go build -o $(BIN_DIR)/$(APIREQUESTER) -ldflags "${LDFLAGS}" cmd/apirequester/*.go

.PHONY: $(HEALTHCHECK)
$(HEALTHCHECK):
	go build -o $(BIN_DIR)/$(HEALTHCHECK) -ldflags "${LDFLAGS}" cmd/healthcheck/*.go

.PHONY: functest
functest:
	CGO_ENABLED=1 go test $(TEST_ARGS) -tags functest ./functest -count=1

.PHONY: test
test:
	CGO_ENABLED=1 go test $(TEST_ARGS) $(ALL_PACKAGES)

.PHONY: test_fast
test_fast:
	go test $(TEST_ARGS) -count 1 -v $(ALL_PACKAGES)

$(ARTIFACTS_DIR):
	mkdir -p $(ARTIFACTS_DIR)

.PHONY: test_with_coverage
test_with_coverage: $(ARTIFACTS_DIR)
	CGO_ENABLED=1 go test $(TEST_ARGS) --coverprofile=$(ARTIFACTS_DIR)/cover.all --covermode=atomic $(TESTED_PACKAGES)
	@cat $(ARTIFACTS_DIR)/cover.all | ./scripts/dev/cover-filter.sh > $(COVERPROFILE)

.PHONY: test_with_coverage_fast
test_with_coverage_fast:
	CGO_ENABLED=1 go test $(TEST_ARGS) -count 1 --coverprofile=$(COVERPROFILE) --covermode=atomic $(ALL_PACKAGES)

.PHONY: ci_test_with_coverage
ci_test_with_coverage:
	CGO_ENABLED=1 go test $(TEST_ARGS) -count 1 -parallel 4 --coverprofile=$(COVERPROFILE) --covermode=atomic -v $(ALL_PACKAGES) | tee unit.file

.PHONY: ci_test_func
ci_test_func:
	CGO_ENABLED=1 go test $(TEST_ARGS) -tags functest -v ./functest -count=1 | tee func.file

.PHONY: ci_test_integrtest
ci_test_integrtest:
	CGO_ENABLED=1 go test $(TEST_ARGS) -tags networktest -v ./network/servicenetwork -count=1 | tee integr.file


.PHONY: regen-proxies
CONTRACTS = $(wildcard application/contract/*)
regen-proxies: $(BININSGOCC)
	$(foreach c, $(CONTRACTS), $(BININSGOCC) proxy application/contract/$(notdir $(c))/$(notdir $(c)).go; )

.PHONY: docker-pulsar
docker-pulsar:
	docker build --tag insolar/pulsar -f ./docker/Dockerfile.pulsar .

.PHONY: docker-insolard
docker-insolard:
	docker build --target insolard --tag insolar/insolard -f ./docker/Dockerfile .

.PHONY: docker-genesis
docker-genesis:
	docker build --target genesis --tag insolar/genesis -f ./docker/Dockerfile .

.PHONY: docker-insgorund
docker-insgorund:
	docker build --target insgorund --tag insolar/insgorund -f ./docker/Dockerfile .

.PHONY: docker
docker: docker-insolard docker-genesis docker-insgorund

$(RECORDBUILDER):
	go build -o $(BIN_DIR)/$(RECORDBUILDER) -ldflags "${LDFLAGS}" cmd/protobuf-record-gen/*.go

generate-protobuf:
	protoc -I./vendor -I./ --gogoslick_out=./ network/node/internal/node/node.proto
	PATH="$(BIN_DIR):$(PATH)" protoc -I./vendor -I./ --gorecord_out=./ insolar/record/record.proto

regen-builtin: $(BININSGOCC)
	$(BININSGOCC) regen-builtin
