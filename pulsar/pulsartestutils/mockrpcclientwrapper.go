/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

// Package pulsartestutil - test utils for pulsar package
package pulsartestutils

import (
	"net/rpc"

	"github.com/insolar/insolar/configuration"
	"github.com/stretchr/testify/mock"
)

// MockRPCClientWrapper is a mock of RPCClientWrapper interface
// It uses testify.mock
type MockRPCClientWrapper struct {
	mock.Mock
}

func (mock *MockRPCClientWrapper) Lock() {
	mock.Mock.Called()
}

func (mock *MockRPCClientWrapper) Unlock() {
	mock.Mock.Called()
}

func (mock *MockRPCClientWrapper) IsInitialised() bool {
	args := mock.Mock.Called()
	return args.Bool(0)
}

func (mock *MockRPCClientWrapper) SetRPCClient(client *rpc.Client) {
	mock.Mock.Called(client)
}

func (mock *MockRPCClientWrapper) CreateConnection(connectionType configuration.ConnectionType, connectionAddress string) error {
	args := mock.Mock.Called(connectionType, connectionAddress)
	return args.Error(0)
}

func (mock *MockRPCClientWrapper) Close() error {
	args := mock.Mock.Called()
	return args.Error(0)
}

func (mock *MockRPCClientWrapper) ResetClient() {
	mock.Called()
}

func (mock *MockRPCClientWrapper) Go(serviceMethod string, args interface{}, reply interface{}, done chan *rpc.Call) *rpc.Call {
	mockArgs := mock.Mock.Called(serviceMethod, args, reply, done)
	return mockArgs.Get(0).(*rpc.Call)
}
