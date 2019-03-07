IMAGE = insolar/genesis
VERSION ?= v0.8.4

.PHONY: build
build:
	docker build --build-arg VERSION=$(VERSION) -t $(IMAGE):$(VERSION) .

.PHONY: push
push:
	docker push $(IMAGE):$(VERSION)

.PHONY: run
run:
	docker run -e GENESIS_OUTPUT_DIR=/data --rm -it $(IMAGE):$(VERSION) sh
