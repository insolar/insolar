// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package future

import (
	"context"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/stretchr/testify/require"
)

func newPacket() *packet.Packet {
	sender, _ := host.NewHostN("127.0.0.1:31337", gen.Reference())
	receiver, _ := host.NewHostN("127.0.0.2:31338", gen.Reference())
	return packet.NewPacket(sender, receiver, types.Pulse, 123)
}

func TestNewPacketHandler(t *testing.T) {
	ph := NewPacketHandler(NewManager())

	require.IsType(t, ph, &packetHandler{})
}

func TestPacketHandler_Handle_Response(t *testing.T) {
	m := NewManager()
	ph := NewPacketHandler(m)

	req := newPacket()
	req.SetRequest(&packet.PulseRequest{})

	future := m.Create(req)
	resp := newPacket()
	resp.Receiver = req.Sender
	resp.Sender = req.Receiver
	resp.SetResponse(&packet.BasicResponse{})

	receivedPacket := packet.NewReceivedPacket(resp, nil)
	ph.Handle(context.Background(), receivedPacket)

	res, err := future.WaitResponse(time.Minute)

	require.NoError(t, err)
	require.Equal(t, receivedPacket, res)
}

func TestPacketHandler_Handle_NotResponse(t *testing.T) {
	m := NewManager()
	ph := NewPacketHandler(m)

	req := newPacket()
	future := m.Create(req)

	resp := newPacket()
	resp.Receiver = req.Sender
	resp.Sender = req.Receiver

	ph.Handle(context.Background(), packet.NewReceivedPacket(resp, nil))

	_, err := future.WaitResponse(time.Millisecond)

	require.Error(t, err)
	require.Equal(t, err, ErrTimeout)
}

func TestPacketHandler_Handle_NotProcessable(t *testing.T) {
	m := NewManager()
	ph := NewPacketHandler(m)

	req := newPacket()
	req.SetRequest(&packet.PulseRequest{})
	future := m.Create(req)

	resp := newPacket()
	resp.SetResponse(&packet.BasicResponse{})

	ph.Handle(context.Background(), packet.NewReceivedPacket(resp, nil))

	_, err := future.WaitResponse(time.Minute)

	require.Error(t, err)
	require.Equal(t, err, ErrChannelClosed)
}

func TestShouldProcessPacket(t *testing.T) {
	m := NewManager()

	req := newPacket()
	future := m.Create(req)

	resp := newPacket()
	resp.Receiver = req.Sender
	resp.Sender = req.Receiver

	require.True(t, shouldProcessPacket(future, packet.NewReceivedPacket(resp, nil)))
}

func TestShouldProcessPacket_WrongType(t *testing.T) {
	m := NewManager()

	req := newPacket()
	future := m.Create(req)

	resp := newPacket()
	resp.Receiver = req.Sender
	resp.Sender = req.Receiver
	resp.Type = uint32(types.RPC)

	require.False(t, shouldProcessPacket(future, packet.NewReceivedPacket(resp, nil)))
}

func TestShouldProcessPacket_WrongSender(t *testing.T) {
	m := NewManager()

	req := newPacket()
	future := m.Create(req)

	require.False(t, shouldProcessPacket(future, packet.NewReceivedPacket(req, nil)))
}
