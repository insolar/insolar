#!/usr/bin/env bash

export GOPATH=`go env GOPATH`
go get -u github.com/gojuno/minimock/cmd/minimock
cd $GOPATH/src/github.com/gojuno/minimock/cmd/minimock
git checkout 890c67cef23dd06d694294d4f7b1026ed7bac8e6
go install