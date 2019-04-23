#                             BUILDER CONTAINER
FROM golang:1.11.5-alpine3.9 AS insolard-builder
MAINTAINER eugene.blikh@insolar.io

RUN set -x && \
        mkdir -p /opt/bin /opt/config /opt/tmp /opt/log && \
        apk update     && \
        apk add --no-cache --virtual /opt/.build-deps \
            git           \
            make          \
            alpine-sdk    \
            bash

COPY [ "/docker/genconfig/genconfig.go", \
       "/opt/tmp/" ]

COPY . /go/src/github.com/insolar/insolar
WORKDIR /go/src/github.com/insolar/insolar

RUN set -x && \
        : "--------- insolard bulding --------------" && \
        : 'cd "/go/src/github.com/insolar/insolar"' && \
        make install-deps ensure && \
        make insolar insolard benchmark apirequester && \
        cp bin/insolar bin/insolard bin/benchmark bin/apirequester /opt/bin && \
        : "--------- helper utilities --------------" && \
        go get gopkg.in/yaml.v2 && \
        go build -o /opt/bin/genconfig   /opt/tmp/genconfig.go

#                             INSOLARD CONTAINER
FROM alpine:3.9 AS insolard
MAINTAINER eugene.blikh@insolar.io

RUN set -x && \
        mkdir -p /var/log /opt/bin /etc/insolar /var/lib/insolar && \
        apk add --no-cache --virtual /opt/.runtime-deps \
            bash    \
            curl

COPY [ "/docker/insolard.yaml.default", "/opt/tmp/" ]
COPY [ "/docker/healthcheck.sh", "/opt/bin/" ]

COPY --from=insolard-builder        \
        [ "/opt/bin/insolar",       \
          "/opt/bin/insolard",      \
          "/opt/bin/benchmark",     \
          "/opt/bin/apirequester",  \
          "/opt/bin/genconfig",     \
          "/opt/bin/"]
COPY docker/insolard-entrypoint.sh /opt/bin/docker-entrypoint.sh

EXPOSE 7900 7901 8001 18182 19191
HEALTHCHECK --start-period=5m --interval=1m --timeout=30s --retries=5 \
    CMD [ "/opt/bin/healthcheck.sh" ]
VOLUME [ "/etc/insolar", "/var/lib/insolar" ]

ENTRYPOINT [ "/opt/bin/docker-entrypoint.sh"  ]
CMD []

#                               GENESIS BUILDER
FROM insolard-builder AS insolard-genesis-builder
MAINTAINER ivan.serezhkin@insolar.io

WORKDIR /opt
RUN set -x && \
        cp -rf /go/src/github.com/insolar/insolar/keys /opt || true && \
        /go/src/github.com/insolar/insolar/scripts/build/genesis/genesis.sh

#                              GENESIS CONTAINER
FROM alpine:3.9 AS genesis
MAINTAINER ivan.serezhkin@insolar.io

RUN mkdir -p /opt/genesis
COPY --from=insolard-genesis-builder    \
        [ "/opt/data.tgz",              \
          "/opt/keys.tgz",              \
          "/opt/genesis/" ]
VOLUME [ "/opt/output" ]
CMD ['cp', '/opt/genesis', '/opt/output/genesis']

#                             INSGORUND CONTAINER
FROM alpine:3.9 AS insgorund
MAINTAINER ivan.serezhkin@insolar.io

RUN set -x && \
        mkdir -p /var/log /opt/bin /tmp/insgorund /opt/insgorund && \
        apk add --no-cache --virtual /opt/.runtime-deps \
            bash    \
            curl

COPY --from=insolard-genesis-builder /opt/bin/insgorund /opt/bin/insgorund
COPY docker/insgorund-entrypoint.sh /opt/bin/docker-entrypoint.sh

EXPOSE 7777 8002

ENTRYPOINT [ "/opt/bin/docker-entrypoint.sh" ]
CMD []
