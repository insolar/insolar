FROM centos:7.6.1810 AS insolard

RUN mkdir -p /var/log /opt/bin /tmp/insgorund /opt/insgorund && \
        yum update --assumeyes --skip-broken && \
        yum install bash curl && \
        yum clean all

COPY [ "/docker/testnet-node/insolard.yaml.default", "/opt/tmp/" ]
COPY [ "/docker/testnet-node/healthcheck.sh", "/opt/bin/" ]

COPY --from=insolar-base          \
       [ "/go/src/github.com/insolar/insolar/bin/insolar",      \
         "/go/src/github.com/insolar/insolar/bin/insolard",     \
         "/go/src/github.com/insolar/insolar/bin/benchmark",    \
         "/go/src/github.com/insolar/insolar/bin/apirequester", \
                                  \
         "/opt/bin/"]
COPY /docker/testnet-node/insolard-entrypoint.sh /opt/bin/docker-entrypoint.sh

EXPOSE 7900 7901 8001 18182 19191
HEALTHCHECK --start-period=5m --interval=1m --timeout=30s --retries=5 \
   CMD [ "/opt/bin/healthcheck.sh" ]
VOLUME [ "/etc/insolar", "/var/lib/insolar" ]

ENTRYPOINT [ "/opt/bin/docker-entrypoint.sh"  ]
CMD []

