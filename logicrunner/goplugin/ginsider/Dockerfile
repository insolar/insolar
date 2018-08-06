FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build ginsider.go
CMD ["/app/ginsider"]

