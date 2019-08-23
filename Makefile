BIN_DIR ?= bin
ARTIFACTS_DIR ?= .artifacts
INSOLAR = insolar
INSOLARD = insolard
INSGOCC = insgocc
PULSARD = pulsard
TESTPULSARD = testpulsard
INSGORUND = insgorund
BENCHMARK = benchmark
PULSEWATCHER = pulsewatcher
BACKUPMERGER = backupmerger
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
all: clean install-deps pre-build build ## cleanup, install deps, (re)generate all code and build all binaries

.PHONY: lint
lint: ci-lint ## alias for ci-lint

.PHONY: ci-lint
ci-lint: ## CI lint
	golangci-lint run

.PHONY: metalint
metalint: ## run gometalinter
	gometalinter --vendor $(ALL_PACKAGES)

.PHONY: clean
clean: ## run all cleanup tasks
	go clean $(ALL_PACKAGES)
	rm -f $(COVERPROFILE)
	rm -rf $(BIN_DIR)
	./scripts/insolard/launchnet.sh -l

.PHONY: install-build-tools
install-build-tools: ## install tools for codegen
	./scripts/build/install_build_tools.sh

.PHONY: install-deps
install-deps: install-build-tools ## install dep and codegen tools

.PHONY: pre-build
pre-build: ensure generate regen-builtin ## install dependencies, (re)generates all code

.PHONY: generate
generate: ## run go generate
	GOPATH=`go env GOPATH` go generate -x $(ALL_PACKAGES)

.PHONY: test_git_no_changes
test_git_no_changes: ## checks if no git changes in project dir (for CI Codegen task)
	ci/scripts/git_diff_without_comments.sh

.PHONY: ensure
ensure: ## install all dependencies
	dep ensure

.PHONY: build
build: $(BIN_DIR) $(INSOLARD) $(INSOLAR) $(INSGOCC) $(PULSARD) $(TESTPULSARD) $(INSGORUND) $(HEALTHCHECK) $(BENCHMARK) $(APIREQUESTER) $(PULSEWATCHER) $(BACKUPMERGER) ## build all binaries

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

.PHONY: $(TESTPULSARD)
$(TESTPULSARD):
	$(GOBUILD) -o $(BIN_DIR)/$(TESTPULSARD) -ldflags "${LDFLAGS}" cmd/testpulsard/*.go

.PHONY: $(INSGORUND)
$(INSGORUND):
	CGO_ENABLED=1 $(GOBUILD) -o $(BIN_DIR)/$(INSGORUND) -ldflags "${LDFLAGS}" cmd/insgorund/*.go

.PHONY: $(BENCHMARK)
$(BENCHMARK):
	$(GOBUILD) -o $(BIN_DIR)/$(BENCHMARK) -ldflags "${LDFLAGS}" cmd/benchmark/*.go

.PHONY: $(PULSEWATCHER)
$(PULSEWATCHER):
	$(GOBUILD) -o $(BIN_DIR)/$(PULSEWATCHER) -ldflags "${LDFLAGS}" cmd/pulsewatcher/*.go

.PHONY: $(BACKUPMERGER)
$(BACKUPMERGER):
	$(GOBUILD) -o $(BIN_DIR)/$(BACKUPMERGER) -ldflags "${LDFLAGS}" cmd/backupmerger/*.go

.PHONY: $(APIREQUESTER)
$(APIREQUESTER):
	$(GOBUILD) -o $(BIN_DIR)/$(APIREQUESTER) -ldflags "${LDFLAGS}" cmd/apirequester/*.go

.PHONY: $(HEALTHCHECK)
$(HEALTHCHECK):
	$(GOBUILD) -o $(BIN_DIR)/$(HEALTHCHECK) -ldflags "${LDFLAGS}" cmd/healthcheck/*.go

.PHONY: test_unit
test_unit: ## run all unit tests
	CGO_ENABLED=1 go test $(TEST_ARGS) $(ALL_PACKAGES)

.PHONY: functest
functest: ## run functest FUNCTEST_COUNT times
	CGO_ENABLED=1 go test -test.v $(TEST_ARGS) -tags functest ./functest -count=$(FUNCTEST_COUNT)

.PNONY: functest_race
functest_race: ## run functest 10 times with -race flag
	make clean
	GOBUILD='go build -race' make build
	FUNCTEST_COUNT=10 make functest

.PHONY: test_func
test_func: functest ## alias for functest

.PHONY: test_slow
test_slow: ## run tests with slowtest tag
	CGO_ENABLED=1 go test $(TEST_ARGS) -tags slowtest ./logicrunner/... ./server/internal/... ./ledger/light/integration/...

.PHONY: test
test: test_unit ## alias for test_unit

.PHONY: test_all
test_all: test_unit test_func test_slow ## run all tests (unit, func, slow)

.PHONY: test_with_coverage
test_with_coverage: $(ARTIFACTS_DIR) ## run unit tests with generation of coverage file
	CGO_ENABLED=1 go test $(TEST_ARGS) --coverprofile=$(ARTIFACTS_DIR)/cover.all --covermode=atomic $(TESTED_PACKAGES)
	@cat $(ARTIFACTS_DIR)/cover.all | ./scripts/dev/cover-filter.sh > $(COVERPROFILE)

.PHONY: test_with_coverage_fast
test_with_coverage_fast: ## ???
	CGO_ENABLED=1 go test $(TEST_ARGS) -count 1 --coverprofile=$(COVERPROFILE) --covermode=atomic $(ALL_PACKAGES)

$(ARTIFACTS_DIR):
	mkdir -p $(ARTIFACTS_DIR)

.PHONY: ci_test_with_coverage
ci_test_with_coverage: ## run unit tests with coverage, outputs json to stdout (CI)
	GOMAXPROCS=$(CI_GOMAXPROCS) CGO_ENABLED=1 \
		go test $(CI_TEST_ARGS) $(TEST_ARGS) -json -v -count 1 --coverprofile=$(COVERPROFILE) --covermode=atomic -tags slowtest $(ALL_PACKAGES)

.PHONY: ci_test_unit
ci_test_unit: ## run unit tests 10 times and -race flag, redirects json output to file (CI)
	GOMAXPROCS=$(CI_GOMAXPROCS) CGO_ENABLED=1 \
		go test $(CI_TEST_ARGS) $(TEST_ARGS) -json -v $(ALL_PACKAGES) -race -count 10 | tee ci_test_unit.json

.PHONY: ci_test_slow
ci_test_slow: ## run slow tests just once, redirects json output to file (CI)
	GOMAXPROCS=$(CI_GOMAXPROCS) CGO_ENABLED=1 \
		go test $(CI_TEST_ARGS) $(TEST_ARGS) -json -v -tags slowtest ./logicrunner/... ./server/internal/... ./ledger/light/integration/... -count 1 | tee -a ci_test_unit.json

.PHONY: ci_test_func
ci_test_func: ## run functest 3 times, redirects json output to file (CI)
	# GOMAXPROCS=2, because we launch at least 5 insolard nodes in functest + 1 pulsar,
	# so try to be more honest with processors allocation.
	GOMAXPROCS=$(CI_GOMAXPROCS) CGO_ENABLED=1  \
		go test $(CI_TEST_ARGS) $(TEST_ARGS) -json -tags functest -v ./functest -count 3 -failfast | tee ci_test_func.json

.PHONY: ci_test_integrtest
ci_test_integrtest: ## run networktest 1 time, redirects json output to file (CI)
	GOMAXPROCS=$(CI_GOMAXPROCS) CGO_ENABLED=1 \
		go test $(CI_TEST_ARGS) $(TEST_ARGS) -json -tags networktest -v ./network/tests -count=1 | tee ci_test_integrtest.json

.PHONY: regen-proxies
CONTRACTS = $(wildcard application/contract/*)
regen-proxies: $(BININSGOCC) ## regen contracts proxies
	$(foreach c, $(CONTRACTS), $(BININSGOCC) proxy application/contract/$(notdir $(c))/$(notdir $(c)).go; )

.PHONY: docker-insolard
docker-insolard: ## build insolard docker image
	docker build --target insolard --tag insolar/insolard -f ./docker/Dockerfile .

.PHONY: docker-insgorund
docker-insgorund: ## build insgorund docker image
	docker build --target insgorund --tag insolar/insgorund -f ./docker/Dockerfile .

.PHONY: docker
docker: docker-insolard docker-insgorund ## build insolard and insgorund docker images

.PHONY: generate-protobuf
generate-protobuf: ## generate protobuf structs
	protoc -I./vendor -I./ --gogoslick_out=./ network/node/internal/node/node.proto
	protoc -I./vendor -I./ --gogoslick_out=./ insolar/record/record.proto
	protoc -I./vendor -I./ --gogoslick_out=./ insolar/jet/jet.proto
	protoc -I./vendor -I./ --gogoslick_out=./ --proto_path=${GOPATH}/src insolar/payload/payload.proto
	protoc -I./vendor -I./ --gogoslick_out=./ insolar/pulse/pulse.proto
	protoc -I./vendor -I./ --gogoslick_out=./ --proto_path=${GOPATH}/src network/hostnetwork/packet/packet.proto
	protoc -I./vendor -I./ --gogoslick_out=./ --proto_path=${GOPATH}/src network/consensus/adapters/candidate/profile.proto
		protoc -I/usr/local/include -I./ \
    		-I$(GOPATH)/src \
    		--gogoslick_out=plugins=grpc:./  \
    		ledger/heavy/exporter/record_exporter.proto
		protoc -I/usr/local/include -I./ \
    		-I$(GOPATH)/src \
    		--gogoslick_out=plugins=grpc:./  \
    		ledger/heavy/exporter/pulse_exporter.proto


.PHONY: regen-builtin
regen-builtin: $(BININSGOCC) ## regenerate builtin contracts code
	$(BININSGOCC) regen-builtin

.PHONY: build-track
build-track: ## build logs event tracker tool
	$(GOBUILD) -o $(BIN_DIR)/track ./scripts/cmd/track/track.go

.PHONY: generate-introspector-proto
generate-introspector-proto: ## generate grpc api code and mocks for introspector
	protoc -I/usr/local/include -I./ \
		-I$(GOPATH)/src \
		-I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:./  \
		--grpc-gateway_out=logtostderr=true:. \
		--swagger_out=logtostderr=true:. \
		instrumentation/introspector/introproto/*.proto
	GOPATH=`go env GOPATH` go generate -x ./instrumentation/introspector

.PHONY: prepare-inrospector-proto
prepare-inrospector-proto: ## install tools required for grpc development
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get -u github.com/golang/protobuf/protoc-gen-go

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
