// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build !functest

package api

import (
	"errors"
	"net/http"

	"github.com/insolar/rpc/v2"
)

// FuncTestContractService is a service that provides ability to add custom contracts
type FuncTestContractService struct {
	runner *Runner
}

// NewFuncTestContractService is dummy for NewFuncTestContractService in functest_contract.go that hidden under build tag
func NewFuncTestContractService(runner *Runner) *FuncTestContractService {
	return &FuncTestContractService{runner: runner}
}

type DummyArgs struct{}
type DummyReply struct{}

type UploadReply struct {
	PrototypeRef string `json:"PrototypeRef"`
}

func (s *FuncTestContractService) Upload(r *http.Request, args *DummyArgs, requestBody *rpc.RequestBody, reply *DummyReply) error {
	return errors.New("method allowed only in build with functest tag")
}

func (s *FuncTestContractService) CallConstructor(r *http.Request, args *DummyArgs, requestBody *rpc.RequestBody, reply *DummyReply) error {
	return errors.New("method allowed only in build with functest tag")
}

func (s *FuncTestContractService) CallMethod(r *http.Request, args *DummyArgs, requestBody *rpc.RequestBody, reply *DummyReply) error {
	return errors.New("method allowed only in build with functest tag")
}
