#!/usr/bin/env bash
#
# install tools for codegen
#

go clean -modcache

./scripts/build/fetchdeps github.com/golang/dep/cmd/dep v0.5.3
./scripts/build/fetchdeps golang.org/x/tools/cmd/stringer 63e6ed9258fa6cbc90aab9b1eef3e0866e89b874
./scripts/build/fetchdeps github.com/gojuno/minimock/cmd/minimock v2.1.8
./scripts/build/fetchdeps github.com/gogo/protobuf/protoc-gen-gogoslick v1.2.1
