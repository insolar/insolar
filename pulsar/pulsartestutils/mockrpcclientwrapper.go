package pulsartestutil

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
