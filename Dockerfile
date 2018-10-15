#FROM golang:1.11 as builder
#RUN mkdir -p /go/src/github.com/insolar/insolar
#WORKDIR /go/src/github.com/insolar/insolar
#ENV BIN_DIR="/go/bin"
#ENV CGO_ENABLED=0
#ENV GOOS=linux
#ADD Makefile /go/src/github.com/insolar/insolar/Makefile
#RUN make install-deps



# insolard
FROM registry.ins.world/builder as insolard
ADD . /go/src/github.com/insolar/insolar
RUN make pre-build build

