FROM golang:1.11-alpine3.9
MAINTAINER eugene.blikh@insolar.io

ENV INSOLAR_VERSION=v0.8.1-0-gca1b45d52 \
    INSOLAR_DOWNLOAD_URL=github.com/insolar/insolar \
    FOREGO_VERSION=4bb213294ff63ed4fc2f095d57e425d9ca18a42b \
    FOREGO_DOWNLOAD_URL=github.com/ddollar/forego

RUN set -x && \
        mkdir -p /opt/bin /opt/config /opt/tmp /opt/log && \
        apk update     && \
        apk add --no-cache --virtual /opt/.build-deps \
            git           \
            make          \
            alpine-sdk && \
        apk add --no-cache --virtual /opt/.runtime-deps \
            curl \
            bash

COPY genconfig.go healthcheck.go insolard.yaml.default /opt/tmp/

RUN set -x && \
        : "--------- insolard bulding --------------" && \
        go get ${INSOLAR_DOWNLOAD_URL} && \
        cd "${GOPATH}/src/${INSOLAR_DOWNLOAD_URL}" && \
        git checkout ${INSOLAR_VERSION} && \
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
        go build -o /opt/bin/genconfig   /opt/tmp/genconfig.go   && \
        : "--------- remove obsolete files --------------" && \
        rm -rf "${GOPATH}/src/${INSOLAR_DOWNLOAD_URL}" && \
        rm -rf "${GOPATH}/src/${FOREGO_DOWNLOAD_URL}" && \
        rm -rf "${GOPATH}/pkg/" && \
        apk del /opt/.build-deps

COPY docker-entrypoint.sh /opt/bin/docker-entrypoint.sh

HEALTHCHECK --start-period=3m --interval=1m --timeout=30s CMD [ "/opt/bin/healthcheck" ]

EXPOSE 8001 13831 19101 33301
VOLUME /opt/config
WORKDIR "/opt"
CMD [ "/opt/bin/docker-entrypoint.sh"  ]
