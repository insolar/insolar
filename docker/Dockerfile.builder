FROM golang:1.11

RUN mkdir -p /go/src/github.com/insolar/insolar
ADD . /go/src/github.com/insolar/insolar
WORKDIR /go/src/github.com/insolar/insolar

ENV BIN_DIR="/go/bin"
ENV CGO_ENABLED=1
ENV GOOS=linux

RUN apt-get update && apt-get -y install jq lsof vim && apt-get clean && rm -Rf /go/src/github.com/insolar/insolar/vendor/* && make install-deps && make pre-build && make test && mv /go/src/github.com/insolar/insolar/vendor /go/
