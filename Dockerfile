FROM registry.insolar.io/builder as insolard
ADD . /go/src/github.com/insolar/insolar
ENV BIN_DIR=/go/bin
RUN rm -Rf /go/src/github.com/insolar/insolar/vendor/* && mv /go/vendor /go/src/github.com/insolar/insolar/
RUN make build

