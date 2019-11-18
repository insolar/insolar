#!/usr/bin/env bash
set -e
set -x

cd build/_go_tools
go run ./ls-imports/main.go -u -f tools.go | xargs -tI % go install -v %
