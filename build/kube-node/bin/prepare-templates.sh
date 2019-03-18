#!/bin/bash
set -e
set -u
EXTERNAL_IP=$EXTERNAL_IP
INSOLAR_VERSION=${INSOLAR_VERSION:-"v0.8.4"}

mkdir -p .out
envsubst < templates/service.tmpl.yaml > .out/service.yaml
envsubst < templates/ve-sts.tmpl.yaml > .out/ve-sts.yaml
