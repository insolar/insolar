# 1) build step (approx local build time ~4m w/o cache)
ARG GOLANG_VERSION=1.12.9
FROM golang:${GOLANG_VERSION} AS build

WORKDIR /go/src/github.com/insolar/insolar

# Install required tools step (adds ~400mb - 11.09.2019)
ADD ./scripts/build ./scripts/build
RUN ./scripts/build/install_build_tools.sh

# Install go dependencies step (adds ~800mb - 11.09.2019)
ADD Gopkg.* ./
RUN dep ensure --vendor-only

# adds sources (adds ~15mb 11.09.2019)
ADD . /go/src/github.com/insolar/insolar/

# pass build variables as arguments to avoid adding .git directory
ARG BUILD_NUMBER
ARG BUILD_DATE
ARG BUILD_TIME
ARG BUILD_HASH
ARG BUILD_VERSION
# build step
# (adds ~400mb - 11.09.2019)
RUN BUILD_NUMBER=${BUILD_NUMBER} \
    BUILD_DATE=${BUILD_DATE} \
    BUILD_TIME=${BUILD_TIME} \
    BUILD_HASH=${BUILD_HASH} \
    BUILD_VERSION=${BUILD_VERSION} \
    make build

# 2) Base image for running test tests and binaries distribution iamges
#    approx local build time ~30s (11.09.2019)
FROM golang:${GOLANG_VERSION}
WORKDIR /go/src/github.com/insolar/insolar

# tools for functest (launchnet)
# adds 20mb (11.09.2019)
RUN apt-get -y update && apt-get -y install jq lsof psmisc && apt-get clean

# adds 350mb, >250mb is binaries (11.09.2019)
COPY --from=build /go/src/github.com/insolar/insolar /go/src/github.com/insolar/insolar
# now here we have ~1.2GB (11.09.2019)
# It's ok, because we need all this dependencies and environment for tests insgorund and running functest
