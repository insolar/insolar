FROM golang:1.11 as builder

RUN mkdir -p /go/src/github.com/insolar/insolar
ADD . /go/src/github.com/insolar/insolar
WORKDIR /go/src/github.com/insolar/insolar
COPY Gopkg.toml Gopkg.lock ./
ENV BIN_DIR="/go/bin/"

RUN make install-deps pre-build build


#FROM scratch
#COPY --from=builder /go/src/github.com/insolar/insolar /go/src/github.com/insolar/insolar
#COPY --from=builder /go/bin/insolard /insolard
#COPY --from=builder /go/bin/insgorund /insgorund
ENTRYPOINT ["/go/bin/insolard"]