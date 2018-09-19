FROM golang:1.11 as builder

RUN go get -u github.com/golang/dep/cmd/dep && go install github.com/golang/dep/cmd/dep

#RUN mkdir -p /go/src/github.com/insolar/insolar
#ADD . /go/src/github.com/insolar/insolar
WORKDIR /go/src/github.com/insolar/insolar
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix nocgo github.com/insolar/insolar/cmd/insolard

#FROM scratch
#COPY --from=builder /go/src/github.com/insolar/insolar /go/src/github.com/insolar/insolar
#COPY --from=builder /go/bin/insolard /insolard
#COPY --from=builder /go/src/github.com/insolar/insolar/logicrunner/goplugin/ginsider-cli/ginsider-cli /ginsider-cli
ENTRYPOINT ["/go/bin/insolard"]