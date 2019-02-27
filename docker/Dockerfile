#                             BUILDER CONTAINER
FROM golang:1.11.5-alpine3.9 AS insolard-builder
MAINTAINER eugene.blikh@insolar.io

ENV FOREGO_VERSION=4bb213294ff63ed4fc2f095d57e425d9ca18a42b \
    FOREGO_DOWNLOAD_URL=github.com/ddollar/forego

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
       "/opt/tmp/" ]

COPY . /go/src/github.com/insolar/insolar
WORKDIR /go/src/github.com/insolar/insolar

RUN set -x && \
        : "--------- insolard bulding --------------" && \
        : 'cd "/go/src/github.com/insolar/insolar"' && \
        make install-deps ensure && \
        make insolar insolard benchmark apirequester insgorund && \
        cp bin/insolar bin/insolard bin/benchmark bin/apirequester bin/insgorund /opt/bin && \
        : "--------- forego bulding --------------" && \
	    go get ${FOREGO_DOWNLOAD_URL} && \
	    cd "${GOPATH}/src/${FOREGO_DOWNLOAD_URL}" && \
	    git checkout ${FOREGO_VERSION} && \
	    make build && \
	    cp forego /opt/bin && \
        : "--------- helper utilities --------------" && \
        go get gopkg.in/yaml.v2 && \
        go build -o /opt/bin/healthcheck /opt/tmp/healthcheck.go && \
        go build -o /opt/bin/genconfig   /opt/tmp/genconfig.go

#                             INSOLARD CONTAINER
FROM alpine:3.9 AS insolard
MAINTAINER eugene.blikh@insolar.io

RUN set -x && \
        mkdir -p /var/log /opt/bin /etc/insolar /var/lib/insolar && \
        apk add --no-cache --virtual /opt/.runtime-deps \
            bash    \
            curl

COPY [ "/docker/insolard.yaml.default", \
       "/opt/tmp/" ]

COPY --from=insolard-builder        \
        [ "/opt/bin/insolar",       \
          "/opt/bin/insolard",      \
          "/opt/bin/benchmark",     \
          "/opt/bin/apirequester",  \
          "/opt/bin/healthcheck",   \
          "/opt/bin/genconfig",     \
          "/opt/bin/forego",        \
          "/opt/bin/insgorund",     \
          "/opt/bin/"]
COPY /docker/insolard-entrypoint.sh /opt/bin/docker-entrypoint.sh

EXPOSE 8001 13831 13832 19191 7777
HEALTHCHECK --start-period=3m --interval=1m --timeout=30s CMD [ "/opt/bin/healthcheck" ]
VOLUME [ "/etc/insolar", "/var/log", "/var/lib/insolar" ]

ENTRYPOINT [ "/opt/bin/insolard-entrypoint.sh"  ]
CMD []

#                               GENESIS BUILDER
FROM insolard-builder AS insolard-genesis-builder
MAINTAINER ivan.serezhkin@insolar.io

WORKDIR /opt
RUN set -x && \
        cp -rf /go/src/github.com/insolar/insolar/keys /opt || : && \
        /go/src/github.com/insolar/insolar/scripts/build/genesis/genesis.sh

#                              GENESIS CONTAINER
FROM alpine:3.9 AS genesis
MAINTAINER ivan.serezhkin@insolar.io

RUN mkdir -p /opt/genesis
COPY --from=insolard-genesis-builder    \
        [ "/opt/data.tgz",              \
          "/opt/keys.tgz",              \
          "/opt/genesis/" ]
VOLUME /opt/output
CMD ['cp', '/opt/genesis', '/opt/output/genesis']

#                             INSGORUND CONTAINER
FROM alpine:3.9 AS insgorund
MAINTAINER ivan.serezhkin@insolar.io

RUN mkdir -p /tmp/insgorund /opt/bin/
COPY --from=insolard-genesis-builder /opt/bin/insgorund /insgorund
COPY /docker/insgorund-entrypoint.sh /opt/bin/docker-entrypoint.sh

EXPOSE 7777 8002
VOLUME [ "/var/log" ]

ENTRYPOINT [ "/opt/bin/docker-entrypoint.sh" ]
CMD []
