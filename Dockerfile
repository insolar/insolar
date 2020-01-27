# 1) build step (approx local build time ~4m w/o cache)
ARG GOLANG_VERSION=1.12.16
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

FROM debian:buster-slim
WORKDIR /go/src/github.com/insolar/insolar
RUN  set -eux; \
     groupadd -r insolar --gid=999; \
     useradd -r -g insolar --uid=999 --shell=/bin/bash insolar
COPY $PWD/bin/insolar $PWD/bin/insolard $PWD/bin/keeperd $PWD/bin/pulsard /usr/local/bin/
