//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// +build !functest

package api

import (
	"errors"
	"net/http"
)

// ContractService is a service that provides ability to add custom contracts
type ContractService struct {
	runner *Runner
}

// NewContractService is dummy for NewContractService in contract.go that hidden under build tag
func NewContractService(runner *Runner) *ContractService {
	return &ContractService{runner: runner}
}

type DummyArgs struct{}
type DummyReply struct{}

type UploadReply struct {
	PrototypeRef string `json:"PrototypeRef"`
}

func (s *ContractService) Upload(r *http.Request, args *DummyArgs, reply *DummyReply) error {
	return errors.New("method allowed only in build with functest tag")
}

func (s *ContractService) CallConstructor(r *http.Request, args *DummyArgs, reply *DummyReply) error {
	return errors.New("method allowed only in build with functest tag")
}

func (s *ContractService) CallMethod(r *http.Request, args *DummyArgs, reply *DummyReply) error {
	return errors.New("method allowed only in build with functest tag")
}
