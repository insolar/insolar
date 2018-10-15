package pulsartestutils

import (
	"net"

	"github.com/stretchr/testify/mock"
)

// MockListener mocks net.Listener interface
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
