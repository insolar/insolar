export GO111MODULE ?= on
export GOSUMDB ?= sum.golang.org
export GOFLAGS ?= -mod=vendor

BIN_DIR ?= bin
ARTIFACTS_DIR ?= .artifacts
INSOLAR = insolar
INSOLARD = insolard
INSGOCC = insgocc
PULSARD = pulsard
TESTPULSARD = testpulsard
PULSEWATCHER = pulsewatcher
BACKUPMANAGER = backupmanager
HEALTHCHECK = healthcheck
KEEPERD = keeperd
BADGER = badger
HEAVY_BADGER_TOOL= heavy-badger

ALL_PACKAGES = ./...
MOCKS_PACKAGE = github.com/insolar/insolar/testutils
GOBUILD ?= go build
GOTEST ?= go test

TEST_COUNT ?= 1
TEST_PARALLEL ?= 1
TEST_ARGS ?= -timeout 1200s
GOMAXPROCS ?= 0
TESTED_PACKAGES ?= $(shell go list ${ALL_PACKAGES} | grep -v "${MOCKS_PACKAGE}")
COVERPROFILE ?= coverage.txt
BUILD_TAGS ?=

BUILD_NUMBER := $(TRAVIS_BUILD_NUMBER)
# skip git parsing commands if no git
ifneq ("$(wildcard ./.git)", "")
	BUILD_DATE ?= $(shell ./scripts/dev/git-date-time.sh -d)
	BUILD_TIME ?= $(shell ./scripts/dev/git-date-time.sh -t)
	BUILD_HASH ?= $(shell git rev-parse --short HEAD)
	BUILD_VERSION ?= $(shell git describe --tags)
endif
DOCKER_BASE_IMAGE_TAG ?= $(BUILD_VERSION)
DOCKER_IMAGE_TAG ?= $(DOCKER_BASE_IMAGE_TAG)

GOPATH ?= `go env GOPATH`
LDFLAGS += -X github.com/insolar/insolar/version.Version=${BUILD_VERSION}
LDFLAGS += -X github.com/insolar/insolar/version.BuildNumber=${BUILD_NUMBER}
LDFLAGS += -X github.com/insolar/insolar/version.BuildDate=${BUILD_DATE}
LDFLAGS += -X github.com/insolar/insolar/version.BuildTime=${BUILD_TIME}
LDFLAGS += -X github.com/insolar/insolar/version.GitHash=${BUILD_HASH}

BININSGOCC=$(BIN_DIR)/$(INSGOCC)

SLOW_PKGS = ./logicrunner/... ./server/internal/... ./cmd/backupmanager/... ./ledger/light/integration/... ./ledger/heavy/executor/integration/...  ./ledger/heavy/integration/... ./virtual/integration ./api

.PHONY: all
all: vendor submodule clean pre-build build ## cleanup, install deps, (re)generate all code and build all binaries

.PHONY: submodule
submodule: ## init git submodule
	git submodule init
	git submodule update

.PHONY: lint
lint: ## lint
	golangci-lint run

.PHONY: metalint
metalint: ## run gometalinter
	gometalinter --vendor $(ALL_PACKAGES)

.PHONY: clean
clean: ## run all cleanup tasks
	go clean $(ALL_PACKAGES)
	rm -f $(COVERPROFILE)
	rm -rf $(BIN_DIR)
	./insolar-scripts/insolard/launchnet.sh -l

.PHONY: install-build-tools
install-build-tools: ## install tools for codegen
	./scripts/build/ls-tools.go | xargs -tI % go install -v %

.PHONY: install-deps
install-deps: install-build-tools ## install dep and codegen tools

.PHONY: pre-build
pre-build: install-deps generate regen-builtin ## install dependencies, (re)generates all code

.PHONY: generate
generate: ## run go generate
	go generate -x $(ALL_PACKAGES)

.PHONY: test_git_no_changes
test_git_no_changes: ## checks if no git changes in project dir (for CI Codegen task)
	ci/scripts/git_diff_without_comments.sh

.PHONY: ensure
ensure: ## does nothing (keep it until all direct calls of `make ensure` have will be removed)
	echo 'All dependencies are already in ./vendor! Run `make vendor` manually if needed'

.PHONY: vendor
vendor: ## update vendor dependencies
	rm -rf vendor
	go mod vendor

.PHONY: build
build: $(BIN_DIR) $(INSOLARD) $(INSOLAR) $(INSGOCC) $(PULSARD) $(TESTPULSARD) $(HEALTHCHECK) ## build all binaries
build: $(PULSEWATCHER) $(BACKUPMANAGER) $(KEEPERD) $(HEAVY_BADGER_TOOL)
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

.PHONY: $(INSOLARD)
$(INSOLARD):
	$(GOBUILD) -o $(BIN_DIR)/$(INSOLARD) ${BUILD_TAGS} -ldflags "${LDFLAGS}" cmd/insolard/*.go

.PHONY: $(INSOLAR)
$(INSOLAR):
	$(GOBUILD) -o $(BIN_DIR)/$(INSOLAR) ${BUILD_TAGS} -ldflags "${LDFLAGS}" application/cmd/insolar/*.go

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

.PHONY: $(PULSEWATCHER)
$(PULSEWATCHER):
	$(GOBUILD) -o $(BIN_DIR)/$(PULSEWATCHER) -ldflags "${LDFLAGS}" cmd/pulsewatcher/*.go

.PHONY: $(BACKUPMANAGER)
$(BACKUPMANAGER):
	$(GOBUILD) -o $(BIN_DIR)/$(BACKUPMANAGER) -ldflags "${LDFLAGS}" cmd/backupmanager/*.go

.PHONY: $(HEALTHCHECK)
$(HEALTHCHECK):
	$(GOBUILD) -o $(BIN_DIR)/$(HEALTHCHECK) -ldflags "${LDFLAGS}" cmd/healthcheck/*.go

.PHONY: $(KEEPERD)
$(KEEPERD):
	$(GOBUILD) -o $(BIN_DIR)/$(KEEPERD) -ldflags "${LDFLAGS}" cmd/keeperd/*.go

.PHONY: $(HEAVY_BADGER_TOOL)
$(HEAVY_BADGER_TOOL):
	$(GOBUILD) -o $(BIN_DIR)/$(HEAVY_BADGER_TOOL) ./cmd/heavy-badger/

.PHONY: test_unit
test_unit: ## run all unit tests
	GOMAXPROCS=$(GOMAXPROCS) CGO_ENABLED=1 \
		$(GOTEST) -count=$(TEST_COUNT) -p=$(TEST_PARALLEL) $(TEST_ARGS) $(ALL_PACKAGES)

.PHONY: test
test: test_unit ## alias for test_unit

.PHONY: functest
functest: ## run functests
	GOMAXPROCS=$(GOMAXPROCS) CGO_ENABLED=1 \
		$(GOTEST) -test.v -p=$(TEST_PARALLEL) $(TEST_ARGS) -tags "functest" ./application/functest -count=$(TEST_COUNT) -failfast
	GOMAXPROCS=$(GOMAXPROCS) CGO_ENABLED=1 \
    		$(GOTEST) -test.v -p=$(TEST_PARALLEL) $(TEST_ARGS) -tags "bloattest" ./application/functest -count=1 -failfast
	GOMAXPROCS=$(GOMAXPROCS) CGO_ENABLED=1 \
		$(GOTEST) -test.v -p=$(TEST_PARALLEL) $(TEST_ARGS) -tags "functest" ./applicationbase/functest -count=$(TEST_COUNT) -failfast
	GOMAXPROCS=$(GOMAXPROCS) CGO_ENABLED=1 \
        		$(GOTEST) -test.v -p=$(TEST_PARALLEL) $(TEST_ARGS) -tags "functest_endless_abandon" ./application/functest -count=1 -failfast

.PHONY: test_func
test_func: functest ## alias for functest

.PNONY: functest_race
functest_race: ## run functest 10 times with -race flag
	make clean
	GOBUILD='go build -race' make build
	TEST_COUNT=10 make functest

.PHONY: test_slow
test_slow: ## run tests with slowtest tag
	CGO_ENABLED=1 GOMAXPROCS=$(GOMAXPROCS) \
		$(GOTEST) $(TEST_ARGS) -p=$(TEST_PARALLEL) -tags slowtest $(ALL_PACKAGES) -count=$(TEST_COUNT)

.PHONY: test_all
test_all: test_unit test_func test_slow ## run all tests (unit, func, slow)

.PHONY: test_with_coverage
test_with_coverage: $(ARTIFACTS_DIR) ## run unit tests with generation of coverage file
	CGO_ENABLED=1 $(GOTEST) $(TEST_ARGS) -tags coverage --coverprofile=$(ARTIFACTS_DIR)/cover.all --covermode=count $(TESTED_PACKAGES)
	@cat $(ARTIFACTS_DIR)/cover.all | ./scripts/dev/cover-filter.sh > $(COVERPROFILE)

.PHONY: test_with_coverage_fast
test_with_coverage_fast: ## ???
	CGO_ENABLED=1 $(GOTEST) $(TEST_ARGS) -tags coverage -count=$(TEST_COUNT) --coverprofile=$(COVERPROFILE) --covermode=count $(ALL_PACKAGES)

$(ARTIFACTS_DIR):
	mkdir -p $(ARTIFACTS_DIR)

.PHONY: test-with-coverage
test-with-coverage: ## run unit tests with coverage, outputs json to stdout (CI)
	GOMAXPROCS=$(GOMAXPROCS) CGO_ENABLED=1 \
		$(GOTEST) $(TEST_ARGS) -json -v -count=$(TEST_COUNT) -p=$(TEST_PARALLEL) --coverprofile=$(COVERPROFILE) --covermode=count -tags 'coverage' $(ALL_PACKAGES)

.PHONY: test-integrtest
test-integrtest: ## run networktest 1 time
	GOMAXPROCS=$(GOMAXPROCS) CGO_ENABLED=1 \
		$(GOTEST) $(TEST_ARGS) -tags networktest -v ./network/tests -count=$(TEST_COUNT)

.PHONY: regen-proxies
CONTRACTS = $(wildcard applicationbase/contract/*)
regen-proxies: $(BININSGOCC) ## regen contracts proxies
	$(foreach c, $(CONTRACTS), $(BININSGOCC) proxy applicationbase/contract/$(notdir $(c))/$(notdir $(c)).go; )

.PHONY: generate-protobuf
generate-protobuf: ## generate protobuf structs
	protoc -I./vendor -I./ --gogoslick_out=./ network/node/internal/node/node.proto
	protoc -I./vendor -I./ --gogoslick_out=./ insolar/record/record.proto
	protoc -I./vendor -I./ --gogoslick_out=./ insolar/jet/jet.proto
	protoc -I./vendor -I./ --gogoslick_out=./ insolar/node.proto
	protoc -I./vendor -I./ --gogoslick_out=./ ledger/drop/drop.proto
	protoc -I./vendor -I./ --gogoslick_out=./ insolar/record/record.proto
	protoc -I./vendor -I./ --gogoslick_out=./ --proto_path=${GOPATH}/src insolar/payload/payload.proto
	protoc -I./vendor -I./ --gogoslick_out=./ insolar/pulse/pulse.proto
	protoc -I./vendor -I./ --gogoslick_out=./ --proto_path=${GOPATH}/src network/hostnetwork/packet/packet.proto
	protoc -I./vendor -I./ --gogoslick_out=./ --proto_path=${GOPATH}/src network/consensus/adapters/candidate/profile.proto
	protoc -I./vendor -I./ --gogoslick_out=./ network/servicenetwork/watermillmessage.proto
	protoc -I./vendor -I./ --gogoslick_out=./ ledger/heavy/executor/jetinfo.proto
	protoc -I./vendor -I./ --gogoslick_out=./ instrumentation/instracer/span_data.proto
	protoc -I./vendor -I/usr/local/include -I./ \
    		-I$(GOPATH)/src \
    		--gogoslick_out=plugins=grpc:./  \
    		ledger/heavy/exporter/record_exporter.proto
	protoc -I./vendor -I/usr/local/include -I./ \
    		-I$(GOPATH)/src \
    		--gogoslick_out=plugins=grpc:./  \
    		ledger/heavy/exporter/pulse_exporter.proto


.PHONY: regen-builtin
regen-builtin: $(BININSGOCC) ## regenerate builtin contracts code
	$(BININSGOCC) regen-builtin -c application/builtin/contract -i github.com/insolar/insolar/application/builtin/contract
# 	rebuild wrapper for panicAsLogicalError to make panic logical error
	$(BININSGOCC) wrapper -p -o application/builtin/contract/panicAsLogicalError/panicAsLogicalError.wrapper.go application/builtin/contract/panicAsLogicalError/panicAsLogicalError.go -m builtin
	$(BININSGOCC) regen-builtin -c applicationbase/builtin/contract -i github.com/insolar/insolar/applicationbase/builtin/contract

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

.PHONY: docker_base_build
docker_base_build: ## build base image with source dependencies and compiled binaries
	docker build -t insolar/insolar-base:$(DOCKER_BASE_IMAGE_TAG) \
		--build-arg BUILD_DATE="$(BUILD_DATE)" \
		--build-arg BUILD_TIME="$(BUILD_TIME)" \
		--build-arg BUILD_NUMBER="$(BUILD_NUMBER)" \
		--build-arg BUILD_HASH="$(BUILD_HASH)" \
		--build-arg BUILD_VERSION="$(BUILD_VERSION)" \
		-f ./scripts/kube/bootstrap/Dockerfile .
	docker tag insolar/insolar-base:$(DOCKER_BASE_IMAGE_TAG) insolar/insolar-base:latest
	docker images "insolar/insolar-base"

.PHONY: docker_build
docker_build: ## build image with binaries and files required for kubernetes deployment.
	docker build -t insolar/insolar:$(DOCKER_IMAGE_TAG) \
		--build-arg BUILD_DATE="$(BUILD_DATE)" \
		--build-arg BUILD_TIME="$(BUILD_TIME)" \
		--build-arg BUILD_NUMBER="$(BUILD_NUMBER)" \
		--build-arg BUILD_HASH="$(BUILD_HASH)" \
		--build-arg BUILD_VERSION="$(BUILD_VERSION)" \
		-f Dockerfile .
	docker images "insolar/insolar"

.PHONY: docker_clean
docker_clean: ## removes intermediate docker image layers w/o tags (beware: it clean up space, but resets caches)
	docker image prune -f


.PHONY: application-git-config
application-git-config: ## change remote repo url from HTTPS to SSH in 'application' submodule
	perl -i.bak -pe 's!url\s*=\s*https?://github.com/!url = git\@github.com:!' .git/modules/application/config

.PHONY: help
help: ## display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
