package pulsartestutil

import (
	"net"

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

func (MockStorage) SetLastPulse(pulse *core.Pulse) error {
	panic("implement me")
}

func (MockStorage) SavePulse(pulse *core.Pulse) error {
	panic("implement me")
}

// Mock of entropy generator for pulsar's tests
var MockEntropy = [64]byte{1, 2, 3, 4, 5, 6, 7, 8}

type MockEntropyGenerator struct {
}

func (MockEntropyGenerator) GenerateEntropy() core.Entropy {
	return MockEntropy
}
