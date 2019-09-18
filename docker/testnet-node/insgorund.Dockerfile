FROM centos:7.6.1810 AS insgorund

RUN mkdir -p /var/log /opt/bin /tmp/insgorund /opt/insgorund && \
        yum update --assumeyes --skip-broken && \
        yum install bash curl && \
        yum clean all

COPY --from=insolar-base /go/src/github.com/insolar/insolar/bin/insgorund /opt/bin/insgorund
COPY docker/testnet-node/insgorund-entrypoint.sh /opt/bin/docker-entrypoint.sh

EXPOSE 7777 8002

ENTRYPOINT [ "/opt/bin/docker-entrypoint.sh" ]
CMD []
