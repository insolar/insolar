.PHONY: all
all: insolard insgorund ## build insolard and insgorund docker images

.PHONY: insolard
insolard: ## build insolard docker image
	docker build --target insolard --tag insolar/insolard -f insolard.Dockerfile ../../.

.PHONY: insgorund
insgorund: ## build insgorund docker image
	docker build --target insgorund --tag insolar/insgorund -f insgorund.Dockerfile ../../.

