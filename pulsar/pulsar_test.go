package pulsar

import (
	"net"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockListener struct {
	mock.Mock
}

func (mock *mockListener) Accept() (net.Conn, error) {
	panic("implement me")
}

func (mock *mockListener) Close() error {
	panic("implement me")
}

func (mock *mockListener) Addr() net.Addr {
	panic("implement me")
}

func TestNewPulsar_WithoutNeighbours(t *testing.T) {
	assert := assert.New(t)
	config := configuration.Pulsar{ConnectionType: "testType", ListenAddress: "listedAddress"}
	actualConnectionType := ""
	actualAddress := ""

	result := NewPulsar(config, func(connectionType string, address string) (net.Listener, error) {
		actualConnectionType = connectionType
		actualAddress = address
		return &mockListener{}, nil
	})

	assert.Equal("testType", actualConnectionType)
	assert.Equal("listedAddress", actualAddress)
	assert.IsType(result.Sock, &mockListener{})
	assert.NotNil(result.PrivateKey)
}

func TestNewPulsar_WithtNeighbours(t *testing.T) {
	assert := assert.New(t)
	config := configuration.Pulsar{
		ConnectionType: "testType",
		ListenAddress:  "listedAddress",
		NodesAddresses: []*configuration.PulsarNodeAddress{
			{ConnectionType: "tcp", Address: "first"},
			{ConnectionType: "pct", Address: "second"},
		},
	}

	result := NewPulsar(config, func(connectionType string, address string) (net.Listener, error) {
		return &mockListener{}, nil
	})

	assert.Equal(2, len(result.Neighbours))

	assert.Equal("tcp", result.Neighbours["first"].ConnectionType.String())
	assert.Equal("pct", result.Neighbours["second"].ConnectionType.String())
}
