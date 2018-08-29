FROM golang:alpine

RUN apk update; apk add git
RUN go get -u github.com/golang/dep/cmd/dep
RUN go install github.com/golang/dep/cmd/dep

RUN mkdir -p /go/src/github.com/insolar/insolar
ADD . /go/src/github.com/insolar/insolar
WORKDIR /go/src/github.com/insolar/insolar
RUN dep ensure
RUN go install github.com/insolar/insolar/cmd/insolar

CMD ["/go/bin/insolar"]
