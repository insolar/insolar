#!/usr/bin/env bash
#
# install tools for codegen
#

go clean -modcache

./scripts/build/fetchdeps github.com/golang/dep/cmd/dep v0.5.3
./scripts/build/fetchdeps golang.org/x/tools/cmd/stringer gopls/v0.1.7
./scripts/build/fetchdeps github.com/gojuno/minimock/cmd/minimock v2.1.8
./scripts/build/fetchdeps github.com/gogo/protobuf/protoc-gen-gogoslick v1.2.1
