/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package packet

import (
	"errors"
	"testing"

	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestBuilder_Build_RequestPacket(t *testing.T) {
	sender, _ := host.NewHostN("127.0.0.1:31337", testutils.RandomRef())
	receiver, _ := host.NewHostN("127.0.0.2:31338", testutils.RandomRef())
	builder := NewBuilder(sender)
	m := builder.Receiver(receiver).Type(TestPacket).Request(&RequestTest{[]byte{0, 1, 2, 3}}).Build()

	expectedPacket := &Packet{
		Sender:        sender,
		RemoteAddress: sender.Address.String(),
		Receiver:      receiver,
		Type:          TestPacket,
		Data:          &RequestTest{[]byte{0, 1, 2, 3}},
		IsResponse:    false,
		Error:         nil,
	}
	require.Equal(t, expectedPacket, m)
}

func TestBuilder_Build_ResponsePacket(t *testing.T) {
	sender, _ := host.NewHostN("127.0.0.1:31337", testutils.RandomRef())
	receiver, _ := host.NewHostN("127.0.0.2:31338", testutils.RandomRef())
	builder := NewBuilder(sender)
	m := builder.Receiver(receiver).Type(TestPacket).Response(&ResponseTest{42}).Build()

	expectedPacket := &Packet{
		Sender:        sender,
		RemoteAddress: sender.Address.String(),
		Receiver:      receiver,
		Type:          TestPacket,
		Data:          &ResponseTest{42},
		IsResponse:    true,
		Error:         nil,
	}
	require.Equal(t, expectedPacket, m)
}

func TestBuilder_Build_ErrorPacket(t *testing.T) {
	sender, _ := host.NewHostN("127.0.0.1:31337", testutils.RandomRef())
	receiver, _ := host.NewHostN("127.0.0.2:31338", testutils.RandomRef())
	builder := NewBuilder(sender)
	m := builder.Receiver(receiver).Type(TestPacket).Response(&ResponseTest{}).Error(errors.New("test error")).Build()

	expectedPacket := &Packet{
		Sender:        sender,
		RemoteAddress: sender.Address.String(),
		Receiver:      receiver,
		Type:          TestPacket,
		Data:          &ResponseTest{},
		IsResponse:    true,
		Error:         errors.New("test error"),
	}
	require.Equal(t, expectedPacket, m)
}
