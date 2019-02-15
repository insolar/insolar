FROM pre:latest
ADD . /go/src/github.com/insolar/insolar

ENV BIN_DIR=/go/bin
ENV CGO_ENABLED=1
ENV GOOS=linux

RUN rm -Rf /go/src/github.com/insolar/insolar/vendor/* && mv /go/vendor /go/src/github.com/insolar/insolar/
RUN make build

EXPOSE 8080
EXPOSE 19191
ENTRYPOINT ["/go/bin/insolard"]