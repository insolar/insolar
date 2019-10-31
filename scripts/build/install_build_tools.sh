#!/usr/bin/env bash
#
# install tools for codegen
#

export GO111MODULE=on

go install -mod=vendor golang.org/x/tools/cmd/stringer
go install -mod=vendor github.com/gogo/protobuf/protoc-gen-gogoslick
go install -mod=vendor github.com/gojuno/minimock/v3/cmd/minimock
go install -mod=vendor github.com/golang/protobuf/protoc-gen-go
go install -mod=vendor github.com/dgraph-io/badger/badger
