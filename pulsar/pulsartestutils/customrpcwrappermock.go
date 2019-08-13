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

// Package pulsartestutil - test utils for pulsar package
package pulsartestutils

import (
	"net/rpc"

	"github.com/insolar/insolar/configuration"
)

// CustomRPCWrapperMock is a mock of RPCClientWrapper interface
// It uses hand-created mocks
type CustomRPCWrapperMock struct {
	Done                 *rpc.Call
	IsInitFunc           func() bool
	CreateConnectionFunc func() error
}

func (*CustomRPCWrapperMock) Lock() {
}

func (*CustomRPCWrapperMock) Unlock() {
}

func (impl *CustomRPCWrapperMock) IsInitialised() bool {
	if impl.IsInitFunc == nil {
		return false
	}
	return impl.IsInitFunc()
}

func (*CustomRPCWrapperMock) SetRPCClient(client *rpc.Client) {
}

func (impl *CustomRPCWrapperMock) CreateConnection(connectionType configuration.ConnectionType, connectionAddress string) error {
	if impl.CreateConnectionFunc == nil {
		return nil
	}
	return impl.CreateConnectionFunc()
}

func (*CustomRPCWrapperMock) Close() error {
	return nil
}

func (impl *CustomRPCWrapperMock) Go(serviceMethod string, args interface{}, reply interface{}, done chan *rpc.Call) *rpc.Call {
	return impl.Done
}

func (impl *CustomRPCWrapperMock) ResetClient() {
}
