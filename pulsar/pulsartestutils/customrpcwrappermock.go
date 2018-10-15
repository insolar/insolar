package pulsartestutil

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
