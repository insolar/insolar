FROM golang:1.11 as builder

RUN go get -u github.com/golang/dep/cmd/dep && go install github.com/golang/dep/cmd/dep

#RUN mkdir -p /go/src/github.com/insolar/insolar
#ADD . /go/src/github.com/insolar/insolar
WORKDIR /go/src/github.com/insolar/insolar
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix nocgo github.com/insolar/insolar/cmd/insolar

FROM scratch
COPY --from=builder /go/bin/insolar /insolar
ENTRYPOINT ["/insolar"]
