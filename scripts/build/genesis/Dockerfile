#                             BUILDER CONTAINER
FROM golang:1.11.5-alpine3.9 AS genesis-builder
MAINTAINER eugene.blikh@insolar.io
MAINTAINER ivan.serezhkin@insolar.io

RUN set -x && \
        mkdir -p /opt/bin /opt/config /opt/tmp /opt/log && \
        apk update     && \
        apk add --no-cache --virtual /opt/.build-deps \
            git           \
            make          \
            alpine-sdk    \
            bash

COPY . /go/src/github.com/insolar/insolar
WORKDIR /go/src/github.com/insolar/insolar

RUN make install-deps ensure && go get gopkg.in/yaml.v2

RUN mkdir /work
WORKDIR /work
RUN /go/src/github.com/insolar/insolar/scripts/build/genesis/genesis.sh

#                             RESULT CONTAINER
FROM alpine:3.9 AS insgorund
MAINTAINER ivan.serezhkin@insolar.io

RUN mkdir /cache
COPY --from=genesis-builder /work/bin/insgorund /insgorund

EXPOSE 7777 8002

ENTRYPOINT /insgorund               \
            -d /cache               \
            --listen=localhost:7777 \
            --rpc=insolard:7778     \
            --metrics=localhost:8002

CMD []
