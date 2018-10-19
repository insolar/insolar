FROM alpine:3.4

MAINTAINER Alexey Smirnov "alexey.smirnov@insolar.io"

RUN apk add --no-cache \
        python3 \
    && pip3 install \
        jinja2-cli[yaml] \
        PyYAML 

VOLUME ["/templates", "/out", "/data"]

CMD ["jinja2", "--help"]
