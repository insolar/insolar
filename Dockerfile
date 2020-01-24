FROM debian:buster-slim
RUN  set -eux; \
     groupadd -r insolar --gid=999; \
     useradd -r -g insolar --uid=999 --shell=/bin/bash insolar; \
     apt-get update; \
     apt-get install -y ca-certificates curl dumb-init gnupg openssl; \
     apt-get clean; \
     rm -rf /var/lib/apt/lists/*
COPY $PWD/artifact/insolar $PWD/artifact/insolard $PWD/artifact/keeperd $PWD/artifact/pulsard /usr/local/bin/
