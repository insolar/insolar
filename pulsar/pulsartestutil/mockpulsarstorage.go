package pulsartestutil

import (
	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/mock"
)

// MockPulsarStorage mocks PulsarStorage interface
type MockPulsarStorage struct {
	mock.Mock
}

func (mock *MockPulsarStorage) GetLastPulse() (*core.Pulse, error) {
	args := mock.Called()
	return args.Get(0).(*core.Pulse), args.Error(1)
}

func (*MockPulsarStorage) SetLastPulse(pulse *core.Pulse) error {
	panic("implement me")
}

func (*MockPulsarStorage) SavePulse(pulse *core.Pulse) error {
	return nil
}

func (*MockPulsarStorage) Close() error {
	panic("implement me")
}
