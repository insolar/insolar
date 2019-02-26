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

COPY [ "/docker/genconfig.go",          \
       "/docker/healthcheck.go",        \
       "/docker/insolard.yaml.default", \
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
        go build -o /opt/bin/healthcheck /opt/tmp/healthcheck.go && \
        go build -o /opt/bin/genconfig   /opt/tmp/genconfig.go

#                             RESULT CONTAINER
FROM alpine:3.9 AS insolard
MAINTAINER eugene.blikh@insolar.io

RUN set -x && \
        mkdir -p /var/log /opt/bin /etc/insolar /var/lib/insolar && \
        apk add --no-cache --virtual /opt/.runtime-deps \
            bash    \
            curl

COPY --from=insolard-builder        \
        [ "/opt/bin/insolar",       \
          "/opt/bin/insolard",      \
          "/opt/bin/benchmark",     \
          "/opt/bin/apirequester",  \
          "/opt/bin/healthcheck",   \
          "/opt/bin/genconfig",     \
          "/opt/bin/"]
COPY /docker/docker-entrypoint.sh /opt/bin/docker-entrypoint.sh

EXPOSE 8001 13831 19101 33301
HEALTHCHECK --start-period=3m --interval=1m --timeout=30s CMD [ "/opt/bin/healthcheck" ]
VOLUME [ "/etc/insolar", "/var/log", "/var/lib/insolar" ]

ENTRYPOINT [ "/opt/bin/docker-entrypoint.sh"  ]
CMD []
