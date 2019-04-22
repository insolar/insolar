package transport

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/configuration"
)

func TestNewDatagramTransport(t *testing.T) {
	table := []struct {
		name    string
		cfg     configuration.Transport
		success bool
	}{
		{
			name:    "default config",
			cfg:     configuration.NewHostNetwork().Transport,
			success: true,
		},
		{
			name:    "localhost",
			cfg:     configuration.Transport{Address: "localhost:0"},
			success: true,
		},
		{
			name:    "invalid address",
			cfg:     configuration.Transport{Address: "invalid"},
			success: false,
		},
		{
			name:    "FixedPublicAddress",
			cfg:     configuration.Transport{Address: "localhost:0", FixedPublicAddress: "192.168.1.1"},
			success: true,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			udp, address, err := NewDatagramTransport(test.cfg)
			assert.Equal(t, test.success, err == nil)
			if test.success {
				assert.NoError(t, err)
				assert.NotNil(t, udp)

				_, err = net.ResolveUDPAddr("udp", address)
				assert.NoError(t, err)
			}
		})
	}
}

type testNode struct {
	udp     DatagramTransport
	address string
}

func (testNode) HandleDatagram(address string, buf []byte) {
	log.Println("Datagram from ", address, " data: ", buf)
}

func newTestNode(port int) (*testNode, error) {
	cfg := configuration.NewHostNetwork().Transport
	cfg.Address = fmt.Sprintf("127.0.0.1:%d", port)

	udp, address, err := NewDatagramTransport(cfg)
	if err != nil {
		return nil, err
	}
	result := &testNode{udp: udp, address: address}
	udp.SetDatagramHandler(result)
	return result, nil
}

func TestUdpTransport_SendDatagram(t *testing.T) {
	ctx := context.Background()

	node1, err := newTestNode(0)
	assert.NoError(t, err)
	node2, err := newTestNode(0)
	assert.NoError(t, err)

	err = node1.udp.Start(ctx)
	assert.NoError(t, err)

	err = node2.udp.Start(ctx)
	assert.NoError(t, err)

	err = node1.udp.SendDatagram(ctx, node2.address, []byte{1, 2, 3})
	assert.NoError(t, err)

	err = node2.udp.SendDatagram(ctx, node1.address, []byte{5, 4, 3})
	assert.NoError(t, err)

	err = node1.udp.Stop(ctx)
	assert.NoError(t, err)

	<-time.After(time.Second)
	err = node1.udp.Start(ctx)
	assert.NoError(t, err)

	err = node1.udp.SendDatagram(ctx, node2.address, []byte{1, 2, 3})
	assert.NoError(t, err)

	err = node2.udp.SendDatagram(ctx, node1.address, []byte{5, 4, 3})
	assert.NoError(t, err)

	///
	err = node1.udp.Stop(ctx)
	assert.NoError(t, err)
	err = node2.udp.Stop(ctx)
	assert.NoError(t, err)

}
