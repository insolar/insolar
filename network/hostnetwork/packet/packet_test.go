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

package packet

import (
	"encoding/gob"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gob.Register(&RequestTest{})
}

// func TestSerializePacket(t *testing.T) {
// 	sender, _ := host.NewHostN("127.0.0.1:31337", testutils.RandomRef())
// 	receiver, _ := host.NewHostN("127.0.0.2:31338", testutils.RandomRef())
// 	builder := NewBuilder(sender)
// 	msg := builder.Receiver(receiver).Type(TestPacket).Request(&RequestTest{[]byte{0, 1, 2, 3}}).Build()
//
// 	_, err := SerializePacket(msg)
//
// 	require.NoError(t, err)
// }
//
// func TestDeserializePacket(t *testing.T) {
// 	sender, _ := host.NewHostN("127.0.0.1:31337", testutils.RandomRef())
// 	receiver, _ := host.NewHostN("127.0.0.2:31338", testutils.RandomRef())
// 	builder := NewBuilder(sender)
// 	msg := builder.Receiver(receiver).Type(TestPacket).Request(&RequestTest{[]byte{0, 1, 2, 3}}).Build()
//
// 	serialized, _ := SerializePacket(msg)
//
// 	var buffer bytes.Buffer
//
// 	buffer.Write(serialized)
//
// 	deserialized, err := DeserializePacket(&buffer)
//
// 	require.NoError(t, err)
// 	require.Equal(t, deserialized, msg)
// }
//
// func TestDeserializeBigPacket(t *testing.T) {
// 	hostOne, _ := host.NewHost("127.0.0.1:31337")
//
// 	data := make([]byte, 1024*1024*10)
// 	rand.Read(data)
//
// 	builder := NewBuilder(hostOne)
// 	msg := builder.Receiver(hostOne).Type(TestPacket).Request(&RequestTest{data}).Build()
//
// 	serialized, err := SerializePacket(msg)
// 	require.NoError(t, err)
//
// 	var buffer bytes.Buffer
// 	buffer.Write(serialized)
//
// 	deserializedMsg, err := DeserializePacket(&buffer)
// 	require.NoError(t, err)
//
// 	deserializedData := deserializedMsg.Data.(*RequestTest).Data
// 	require.Equal(t, data, deserializedData)
// }
//
// type PacketSuite struct {
// 	suite.Suite
// 	sender *host.Host
// 	packet *Packet
// }
//
// func (s *PacketSuite) TestGetSender() {
// 	s.Equal(s.packet.GetSender(), s.sender.NodeID)
// }
//
// func (s *PacketSuite) TestGetSenderHost() {
// 	s.Equal(s.packet.GetSenderHost(), s.sender)
// }
//
// func (s *PacketSuite) TestGetType() {
// 	s.Equal(s.packet.GetType(), TestPacket)
// }
//
// func (s *PacketSuite) TestGetData() {
// 	s.Equal(s.packet.GetData(), &RequestTest{[]byte{0, 1, 2, 3}})
// }
//
// func (s *PacketSuite) TestGetRequestID() {
// 	s.Equal(s.packet.GetRequestID(), types.RequestID(123))
// }
//
// func TestPacketMethods(t *testing.T) {
// 	sender, _ := host.NewHostN("127.0.0.1:31337", testutils.RandomRef())
// 	receiver, _ := host.NewHostN("127.0.0.2:31338", testutils.RandomRef())
// 	builder := NewBuilder(sender)
// 	p := builder.
// 		Receiver(receiver).
// 		Type(TestPacket).
// 		Request(&RequestTest{[]byte{0, 1, 2, 3}}).
// 		RequestID(types.RequestID(123)).
// 		Build()
//
// 	suite.Run(t, &PacketSuite{
// 		sender: sender,
// 		packet: p,
// 	})
// }

func marshalUnmarshal(t *testing.T, p1, p2 *PacketBackend) {
	data, err := p1.Marshal()
	require.NoError(t, err)
	err = p2.Unmarshal(data)
	require.NoError(t, err)
}

func marshalUnmarshalPacketRequest(t *testing.T, request interface{}) (p1, p2 *PacketBackend) {
	p1, p2 = &PacketBackend{}, &PacketBackend{}
	p1.SetRequest(request)
	marshalUnmarshal(t, p1, p2)
	require.NotNil(t, p2.GetRequest())
	return p1, p2
}

func marshalUnmarshalPacketResponse(t *testing.T, response interface{}) (p1, p2 *PacketBackend) {
	p1, p2 = &PacketBackend{}, &PacketBackend{}
	p1.SetResponse(response)
	marshalUnmarshal(t, p1, p2)
	require.NotNil(t, p2.GetResponse())
	return p1, p2
}

func TestPacketBackend_SetRequest(t *testing.T) {
	type SomeData struct {
		someField int
	}
	p := PacketBackend{}
	f := func() {
		p.SetRequest(&SomeData{})
	}
	assert.Panics(t, f)
}

func TestPacketBackend_SetResponse(t *testing.T) {
	type SomeData struct {
		someField int
	}
	p := PacketBackend{}
	f := func() {
		p.SetResponse(&SomeData{})
	}
	assert.Panics(t, f)
}

func TestPacketBackend_GetRequest_GetPing(t *testing.T) {
	ping := Ping{}
	_, p2 := marshalUnmarshalPacketRequest(t, &ping)
	assert.NotNil(t, p2.GetRequest().GetPing())
}

func TestPacketBackend_GetRequest_GetRPC(t *testing.T) {
	rpc := RPCRequest{Method: "meth", Data: []byte("123")}
	p1, p2 := marshalUnmarshalPacketRequest(t, &rpc)
	require.NotNil(t, p2.GetRequest().GetRPC())
	assert.Equal(t, p1.GetRequest().GetRPC().Method, p2.GetRequest().GetRPC().Method)
	assert.Equal(t, p1.GetRequest().GetRPC().Data, p2.GetRequest().GetRPC().Data)
}

func TestPacketBackend_GetRequest_GetAuthorize(t *testing.T) {
	ss := []byte("onetwothree")
	auth := AuthorizeRequest{Certificate: ss}
	_, p2 := marshalUnmarshalPacketRequest(t, &auth)
	require.NotNil(t, p2.GetRequest().GetAuthorize())
	assert.Equal(t, ss, p2.GetRequest().GetAuthorize().Certificate)
}

func TestPacketBackend_GetResponse(t *testing.T) {
	cascade := BasicResponse{}
	_, p2 := marshalUnmarshalPacketResponse(t, &cascade)
	assert.NotNil(t, p2.GetResponse().GetBasic())
}
