# 1) build step (approx local build time ~4m w/o cache)
ARG GOLANG_VERSION=1.12.9
FROM golang:${GOLANG_VERSION} AS build

ADD . /go/src/github.com/insolar/insolar
WORKDIR /go/src/github.com/insolar/insolar

# pass build variables as arguments to avoid adding .git directory
ARG BUILD_NUMBER
ARG BUILD_DATE
ARG BUILD_TIME
ARG BUILD_HASH
ARG BUILD_VERSION
# build step
RUN BUILD_NUMBER=${BUILD_NUMBER} \
    BUILD_DATE=${BUILD_DATE} \
    BUILD_TIME=${BUILD_TIME} \
    BUILD_HASH=${BUILD_HASH} \
    BUILD_VERSION=${BUILD_VERSION} \
    make build

# 2) Base image for running tests and binaries distribution images
FROM golang:${GOLANG_VERSION}
WORKDIR /go/src/github.com/insolar/insolar

# tools for functest (launchnet)
RUN RUN apt-get -y update && apt-get -y install jq lsof psmisc && apt-get clean

COPY --from=build /go/src/github.com/insolar/insolar /go/src/github.com/insolar/insolar
# It's ok to have binaries and sources because we need all this dependencies and environment for tests insgorund and running functest
