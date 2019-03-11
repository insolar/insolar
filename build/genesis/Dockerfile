FROM alpine:3.9 AS fetcher
ARG VERSION=v0.8.4
WORKDIR /fetch
RUN wget https://github.com/insolar/insolar/releases/download/$VERSION/insolar-node.tar.gz
RUN tar xzf insolar-node.tar.gz
RUN ls -la /fetch

FROM alpine:3.9

COPY --from=fetcher /fetch/genesis /genesis
ADD init.sh .
