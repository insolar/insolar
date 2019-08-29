//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

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
