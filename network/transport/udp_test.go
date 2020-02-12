// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package transport

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type testNode struct {
	udp     DatagramTransport
	address string
}

func (t *testNode) HandleDatagram(ctx context.Context, address string, buf []byte) {
	inslogger.FromContext(ctx).Info("Handle Datagram ", buf)
}

func newTestNode(port int) (*testNode, error) {
	cfg := configuration.NewHostNetwork().Transport
	cfg.Address = fmt.Sprintf("127.0.0.1:%d", port)

	node := &testNode{}
	udp, err := NewFactory(cfg).CreateDatagramTransport(node)
	if err != nil {
		return nil, err
	}
	node.udp = udp

	err = node.udp.Start(context.Background())
	if err != nil {
		return nil, err
	}

	node.address = udp.Address()
	return node, nil
}

func TestUdpTransport_SendDatagram(t *testing.T) {
	ctx := context.Background()

	node1, err := newTestNode(0)
	assert.NoError(t, err)
	node2, err := newTestNode(0)
	assert.NoError(t, err)

	err = node1.udp.SendDatagram(ctx, node2.address, []byte{1, 2, 3})
	assert.NoError(t, err)

	err = node2.udp.SendDatagram(ctx, node1.address, []byte{5, 4, 3})
	assert.NoError(t, err)

	err = node1.udp.Stop(ctx)
	assert.NoError(t, err)

	err = node1.udp.Start(ctx)
	assert.NoError(t, err)

	err = node1.udp.SendDatagram(ctx, node2.address, []byte{1, 2, 3})
	assert.NoError(t, err)

	err = node2.udp.SendDatagram(ctx, node1.address, []byte{5, 4, 3})
	assert.NoError(t, err)

	err = node1.udp.Stop(ctx)
	assert.NoError(t, err)
	err = node2.udp.Stop(ctx)
	assert.NoError(t, err)
}

func TestUdpTransport_SendDatagram_Error(t *testing.T) {
	cfg := configuration.NewHostNetwork().Transport
	cfg.Address = fmt.Sprintf("127.0.0.1:%d", 0)

	node := &testNode{}
	udp, err := NewFactory(cfg).CreateDatagramTransport(node)
	require.NoError(t, err)

	ctx := context.Background()
	err = udp.SendDatagram(ctx, udp.Address(), []byte{1, 2, 3})
	require.EqualError(t, err, "failed to send datagram: transport is not started")

	err = udp.Start(ctx)
	require.NoError(t, err)

	err = udp.SendDatagram(ctx, udp.Address(), []byte{1, 2, 3})
	require.NoError(t, err)

	err = udp.Stop(ctx)
	require.NoError(t, err)
}
