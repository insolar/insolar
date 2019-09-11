# 1) build step (approx local build time ~4m w/o cache)
ARG GOLANG_VERSION=1.12.9
FROM golang:${GOLANG_VERSION} AS build

RUN mkdir -p /go/src/github.com/insolar/insolar
WORKDIR /go/src/github.com/insolar/insolar

# Install required tools step (adds ~400mb)
RUN echo -n "<DISK USAGE> before install_build_tools: " && du -sh / 2>/dev/null || true
ADD ./scripts/build ./scripts/build
RUN ./scripts/build/install_build_tools.sh

# Install go dependencies step (adds ~800mb)
RUN echo -n "<DISK USAGE> before vendor: " && du -sh / 2>/dev/null || true
ADD Gopkg.* ./
RUN dep ensure --vendor-only

# adds sources (~15mb)
RUN echo -n "<DISK USAGE> before insolar sources copy: " && du -sh / 2>/dev/null || true
ADD . /go/src/github.com/insolar/insolar/

RUN echo -n "<DISK USAGE> before build: " && du -sh / 2>/dev/null || true
# pass build variables as arguments to avoid adding .git directory
ARG BUILD_NUMBER
ARG BUILD_DATE
ARG BUILD_TIME
ARG BUILD_HASH
ARG BUILD_VERSION
# build step: adds 400mb
# not really cachebale because it uses BUILD_TIME variable
# (we can mock BUILD_TIME, but it's not worth it, besacuse build time is around 1 min)
RUN BUILD_NUMBER=${BUILD_NUMBER} \
    BUILD_DATE=${BUILD_DATE} \
    BUILD_TIME=${BUILD_TIME} \
    BUILD_HASH=${BUILD_HASH} \
    BUILD_VERSION=${BUILD_VERSION} \
    make build
RUN echo -n "<DISK USAGE> after build: " && du -sh / 2>/dev/null || true

# 2) Base image for running test tests and binaries distribution iamges
#    (approx local build time ~30s)
FROM golang:${GOLANG_VERSION}
WORKDIR /go/src/github.com/insolar/insolar

# adds 20mb
RUN echo -n "<DISK USAGE> before apt-get: " && du -sh / 2>/dev/null || true
# tools for functest (launchnet)
RUN apt-get -y update && apt-get -y install jq lsof psmisc && apt-get clean

# adds 350mb (>250mb – binaries)
RUN echo -n "<DISK USAGE> before insolar files copy: " && du -sh / 2>/dev/null || true
COPY --from=build /go/src/github.com/insolar/insolar /go/src/github.com/insolar/insolar
#RUN ls -al ./bin && ls -al ./
# now here we have 1.2GB (it's ok, because we need all dependencies and environment for tests insgorund)
#RUN echo -n "<DISK USAGE> result:" && du -sh / 2>/dev/null || true
