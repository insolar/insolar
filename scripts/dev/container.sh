#!/usr/bin/env sh

# Running stuff in docker container. Can be used to simulate low resourse machine.
# For example, to limit cpu to 0.5, pass `--cpus=0.5` to the script. To limit memory pass `-m 512mb`.

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
IMAGE_NAME=insolar_dev
CONTAINER_NAME=insolar_dev

docker build -t $IMAGE_NAME -f $DIR/Dockerfile .
docker run -it \
  --name $CONTAINER_NAME \
  --rm \
  -u $UID \
  -v `pwd`:/go/src/github.com/insolar/insolar \
  -v `go env GOCACHE`:/var/gocache \
  -e GOCACHE=/var/gocache \
  $@ $IMAGE_NAME /bin/bash
