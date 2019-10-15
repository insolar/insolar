#!/usr/bin/env bash
#
# install tools for codegen
#

go clean -modcache
export GO111MODULE=on


./scripts/build/fetchdeps github.com/gogo/protobuf/protoc-gen-gogoslick v1.2.1
./scripts/build/fetchdeps github.com/golang/protobuf/tree/master/protoc-gen-go v1.3.2

./scripts/build/fetchdeps golang.org/x/tools/cmd/stringer 63e6ed9258fa6cbc90aab9b1eef3e0866e89b874
./scripts/build/fetchdeps github.com/gojuno/minimock/cmd/minimock v3.0.5
