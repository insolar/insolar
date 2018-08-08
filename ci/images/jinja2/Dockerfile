FROM alpine:3.4

MAINTAINER Alexey Smirnov "alexey.smirnov@insolar.io"

RUN apk add --no-cache \
        python \
        py-pip \
    && apk add --no-cache --virtual .build-deps \
        python-dev \
    && pip install --upgrade pip \
    && pip install \
        jinja2-cli[yaml] \
        PyYAML \
    && apk del .build-deps

VOLUME ["/templates", "/out", "/data"]

CMD ["jinja2", "--help"]
