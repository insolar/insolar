/*
 *    Copyright 2018 Insolar
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

package pulsartestutil

import (
	"net"
	"net/rpc"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/mock"
)

// Mock of listener for pulsar's tests
type MockListener struct {
	mock.Mock
}

func (mock *MockListener) Accept() (net.Conn, error) {
	panic("implement me")
}

func (mock *MockListener) Close() error {
	panic("implement me")
}

func (mock *MockListener) Addr() net.Addr {
	panic("implement me")
}

// Mock of storage for pulsar's tests
type MockStorage struct {
	mock.Mock
}

func (mock *MockStorage) GetLastPulse() (*core.Pulse, error) {
	args := mock.Called()
	return args.Get(0).(*core.Pulse), args.Error(1)
}

func (*MockStorage) SetLastPulse(pulse *core.Pulse) error {
	panic("implement me")
}

func (*MockStorage) SavePulse(pulse *core.Pulse) error {
	return nil
}

func (*MockStorage) Close() error {
	panic("implement me")
}

// MockEntropy generator for pulsar's tests
var MockEntropy = [64]byte{1, 2, 3, 4, 5, 6, 7, 8}

type MockEntropyGenerator struct {
}

func (MockEntropyGenerator) GenerateEntropy() core.Entropy {
	return MockEntropy
}

// Mock of RpcClientWrapper
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
	} else {
		return impl.IsInitFunc()
	}
}

func (*CustomRPCWrapperMock) SetRPCClient(client *rpc.Client) {
}

func (impl *CustomRPCWrapperMock) CreateConnection(connectionType configuration.ConnectionType, connectionAddress string) error {
	if impl.CreateConnectionFunc == nil {
		return nil
	} else {
		return impl.CreateConnectionFunc()
	}
}

func (*CustomRPCWrapperMock) Close() error {
	return nil
}

func (impl *CustomRPCWrapperMock) Go(serviceMethod string, args interface{}, reply interface{}, done chan *rpc.Call) *rpc.Call {
	return impl.Done
}

func (impl *CustomRPCWrapperMock) ResetClient() {
}
