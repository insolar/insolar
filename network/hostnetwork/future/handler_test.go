/*
 *    Copyright 2019 Insolar Technologies
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

package future

import (
	"context"
	"testing"
	"time"

	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func newPacket() *packet.Packet {
	sender, _ := host.NewHostN("127.0.0.1:31337", testutils.RandomRef())
	receiver, _ := host.NewHostN("127.0.0.2:31338", testutils.RandomRef())
	builder := packet.NewBuilder(sender)
	p := builder.
		Receiver(receiver).
		Type(packet.TestPacket).
		Request(&packet.RequestTest{[]byte{0, 1, 2, 3}}).
		RequestID(network.RequestID(123)).
		Build()
	return p
}

func TestNewPacketHandler(t *testing.T) {
	ph := NewPacketHandler(NewManager())

	require.IsType(t, ph, &packetHandler{})
}

func TestPacketHandler_Handle_Response(t *testing.T) {
	m := NewManager()
	ph := NewPacketHandler(m)

	req := newPacket()
	future := m.Create(req)

	resp := newPacket()
	resp.Receiver = req.Sender
	resp.Sender = req.Receiver
	resp.IsResponse = true

	ph.Handle(context.Background(), resp)

	res, err := future.WaitResponse(time.Millisecond)

	require.NoError(t, err)
	require.Equal(t, resp, res)
}

func TestPacketHandler_Handle_NotResponse(t *testing.T) {
	m := NewManager()
	ph := NewPacketHandler(m)

	req := newPacket()
	future := m.Create(req)

	resp := newPacket()
	resp.Receiver = req.Sender
	resp.Sender = req.Receiver

	ph.Handle(context.Background(), resp)

	_, err := future.WaitResponse(time.Millisecond)

	require.Error(t, err)
	require.Equal(t, err, ErrTimeout)
}

func TestPacketHandler_Handle_NotProcessable(t *testing.T) {
	m := NewManager()
	ph := NewPacketHandler(m)

	req := newPacket()
	future := m.Create(req)

	resp := newPacket()
	resp.IsResponse = true

	ph.Handle(context.Background(), resp)

	_, err := future.WaitResponse(time.Millisecond)

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

	require.True(t, shouldProcessPacket(future, resp))
}

func TestShouldProcessPacket_WrongType(t *testing.T) {
	m := NewManager()

	req := newPacket()
	future := m.Create(req)

	resp := newPacket()
	resp.Receiver = req.Sender
	resp.Sender = req.Receiver
	resp.Type = types.RPC

	require.False(t, shouldProcessPacket(future, resp))
}

func TestShouldProcessPacket_WrongSender(t *testing.T) {
	m := NewManager()

	req := newPacket()
	future := m.Create(req)

	require.False(t, shouldProcessPacket(future, req))
}

func TestShouldProcessPacket_WrongSenderPing(t *testing.T) {
	m := NewManager()

	req := newPacket()
	future := m.Create(req)

	resp := newPacket()
	resp.Type = types.Ping

	require.False(t, shouldProcessPacket(future, resp))
}
