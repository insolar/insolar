FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o ginsider-cli main.go
CMD ["/app/ginsider-cli"]

