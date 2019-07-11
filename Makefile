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

ALL_PACKAGES = ./...
MOCKS_PACKAGE = github.com/insolar/insolar/testutils
GOBUILD ?= go build
FUNCTEST_COUNT ?= 1
TESTED_PACKAGES ?= $(shell go list ${ALL_PACKAGES} | grep -v "${MOCKS_PACKAGE}")
COVERPROFILE ?= coverage.txt
TEST_ARGS ?= -timeout 1200s
BUILD_TAGS ?=

CI_GOMAXPROCS ?= 8
CI_TEST_ARGS ?= -p 4

BUILD_NUMBER := $(TRAVIS_BUILD_NUMBER)
BUILD_DATE = $(shell date "+%Y-%m-%d")
BUILD_TIME = $(shell date "+%H:%M:%S")
BUILD_HASH = $(shell git rev-parse --short HEAD)
BUILD_VERSION ?= $(shell git describe --abbrev=0 --tags)

GOPATH ?= `go env GOPATH`
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
	golangci-lint run

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
	./scripts/build/fetchdeps github.com/golang/dep/cmd/dep v0.5.3

.PHONY: install-build-tools
install-build-tools:
	go clean -modcache
	./scripts/build/fetchdeps golang.org/x/tools/cmd/stringer 63e6ed9258fa6cbc90aab9b1eef3e0866e89b874
	./scripts/build/fetchdeps github.com/gojuno/minimock/cmd/minimock 890c67cef23dd06d694294d4f7b1026ed7bac8e6
	./scripts/build/fetchdeps github.com/gogo/protobuf/protoc-gen-gogoslick v1.2.1

.PHONY: install-deps
install-deps: install-godep install-build-tools

.PHONY: pre-build
pre-build: ensure generate regen-builtin

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
	$(GOBUILD) -o $(BIN_DIR)/$(INSOLARD) ${BUILD_TAGS} -ldflags "${LDFLAGS}" cmd/insolard/*.go

.PHONY: $(INSOLAR)
$(INSOLAR):
	$(GOBUILD) -o $(BIN_DIR)/$(INSOLAR) ${BUILD_TAGS} -ldflags "${LDFLAGS}" cmd/insolar/*.go

.PHONY: $(INSGOCC)
$(INSGOCC): cmd/insgocc/insgocc.go logicrunner/preprocessor
	$(GOBUILD) -o $(BININSGOCC) -ldflags "${LDFLAGS}" cmd/insgocc/*.go

$(BININSGOCC): $(INSGOCC)

.PHONY: $(PULSARD)
$(PULSARD):
	$(GOBUILD) -o $(BIN_DIR)/$(PULSARD) -ldflags "${LDFLAGS}" cmd/pulsard/*.go

.PHONY: $(INSGORUND)
$(INSGORUND):
	CGO_ENABLED=1 $(GOBUILD) -o $(BIN_DIR)/$(INSGORUND) -ldflags "${LDFLAGS}" cmd/insgorund/*.go

.PHONY: $(BENCHMARK)
$(BENCHMARK):
	$(GOBUILD) -o $(BIN_DIR)/$(BENCHMARK) -ldflags "${LDFLAGS}" cmd/benchmark/*.go

.PHONY: $(PULSEWATCHER)
$(PULSEWATCHER):
	$(GOBUILD) -o $(BIN_DIR)/$(PULSEWATCHER) -ldflags "${LDFLAGS}" cmd/pulsewatcher/*.go

.PHONY: $(APIREQUESTER)
$(APIREQUESTER):
	$(GOBUILD) -o $(BIN_DIR)/$(APIREQUESTER) -ldflags "${LDFLAGS}" cmd/apirequester/*.go

.PHONY: $(HEALTHCHECK)
$(HEALTHCHECK):
	$(GOBUILD) -o $(BIN_DIR)/$(HEALTHCHECK) -ldflags "${LDFLAGS}" cmd/healthcheck/*.go

.PHONY: test_unit
test_unit:
	CGO_ENABLED=1 go test $(TEST_ARGS) $(ALL_PACKAGES)

.PHONY: functest
functest:
	CGO_ENABLED=1 go test -test.v $(TEST_ARGS) -tags functest ./functest -count=$(FUNCTEST_COUNT)

.PNONY: functest_race
functest_race:
	make clean
	GOBUILD='go build -race' make build
	FUNCTEST_COUNT=10 make functest

.PHONY: test_func
test_func: functest

.PHONY: test_slow
test_slow:
	CGO_ENABLED=1 go test $(TEST_ARGS) -tags slowtest ./logicrunner/... ./server/internal/...

.PHONY: test
test: test_unit

.PHONY: test_all
test_all: test_unit test_func test_slow

.PHONY: test_with_coverage
test_with_coverage: $(ARTIFACTS_DIR)
	CGO_ENABLED=1 go test $(TEST_ARGS) --coverprofile=$(ARTIFACTS_DIR)/cover.all --covermode=atomic $(TESTED_PACKAGES)
	@cat $(ARTIFACTS_DIR)/cover.all | ./scripts/dev/cover-filter.sh > $(COVERPROFILE)

.PHONY: test_with_coverage_fast
test_with_coverage_fast:
	CGO_ENABLED=1 go test $(TEST_ARGS) -count 1 --coverprofile=$(COVERPROFILE) --covermode=atomic $(ALL_PACKAGES)

$(ARTIFACTS_DIR):
	mkdir -p $(ARTIFACTS_DIR)

.PHONY: ci_test_with_coverage
ci_test_with_coverage:
	GOMAXPROCS=$(CI_GOMAXPROCS) CGO_ENABLED=1 \
		go test $(CI_TEST_ARGS) $(TEST_ARGS) -json -v -count 1 --coverprofile=$(COVERPROFILE) --covermode=atomic -tags slowtest $(ALL_PACKAGES)

.PHONY: ci_test_unit
ci_test_unit:
	GOMAXPROCS=$(CI_GOMAXPROCS) CGO_ENABLED=1 \
		go test $(CI_TEST_ARGS) $(TEST_ARGS) -json -v $(ALL_PACKAGES) -race -count 10 | tee ci_test_unit.json

.PHONY: ci_test_slow
ci_test_slow:
	GOMAXPROCS=$(CI_GOMAXPROCS) CGO_ENABLED=1 \
		go test $(CI_TEST_ARGS) $(TEST_ARGS) -json -v -tags slowtest ./logicrunner/... ./server/internal/... -count 1 | tee -a ci_test_unit.json

.PHONY: ci_test_func
ci_test_func:
	# GOMAXPROCS=2, because we launch at least 5 insolard nodes in functest + 1 pulsar,
	# so try to be more honest with processors allocation.
	GOMAXPROCS=2 CGO_ENABLED=1 INSOLAR_LOG_LEVEL=error \
		go test $(CI_TEST_ARGS) $(TEST_ARGS) -json -tags functest -v ./functest -count 3 -failfast | tee ci_test_func.json

.PHONY: ci_test_integrtest
ci_test_integrtest:
	GOMAXPROCS=$(CI_GOMAXPROCS) CGO_ENABLED=1 \
		go test $(CI_TEST_ARGS) $(TEST_ARGS) -json -tags networktest -v ./network/tests -count=1 | tee ci_test_integrtest.json

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

generate-protobuf:
	protoc -I./vendor -I./ --gogoslick_out=./ network/node/internal/node/node.proto
	protoc -I./vendor -I./ --gogoslick_out=./ insolar/record/record.proto
	protoc -I./vendor -I./ --gogoslick_out=./ --proto_path=${GOPATH}/src insolar/payload/payload.proto
	protoc -I./vendor -I./ --gogoslick_out=./ ledger/object/lifeline.proto
	protoc -I./vendor -I./ --gogoslick_out=./ ledger/object/filamentindex.proto
	protoc -I./vendor -I./ --gogoslick_out=./ insolar/pulse/pulse.proto
	protoc -I./vendor -I./ --gogoslick_out=./ --proto_path=${GOPATH}/src network/hostnetwork/packet/packet.proto

regen-builtin: $(BININSGOCC)
	$(BININSGOCC) regen-builtin

build-track:
	$(GOBUILD) -o $(BIN_DIR)/track ./scripts/cmd/track/track.go
