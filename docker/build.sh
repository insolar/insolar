#!/usr/bin/env bash
# Script Name: build.sh
# Description: the script hepls to build insolar docker image,
#              just run this script from anywhere

set -v
set -x

SCRIPT_WORKDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
BUILD_NUMBER=1
BUILD_DATE=$(date +"%D")
BUILD_TIME=$(date +"%H-%M-%S")
BUILD_HASH=$(git rev-parse HEAD)
BUILD_VERSION=$(git describe --abbrev=0 --tags | awk -Fv '{ print $2 }')

docker build \
  -t insolar \
  -f "${SCRIPT_WORKDIR}"/Dockerfile \
  --build-arg BUILD_NUMBER="${BUILD_NUMBER}" \
  --build-arg BUILD_DATE="${BUILD_DATE}" \
  --build-arg BUILD_TIME="${BUILD_TIME}" \
  --build-arg BUILD_HASH="${BUILD_HASH}" \
  --build-arg BUILD_VERSION="${BUILD_VERSION}" \
  --pull \
  "${SCRIPT_WORKDIR}"/..
