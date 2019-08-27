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
	"bytes"
	"math/rand"
	"testing"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func testRPCPacket() *Packet {
	sender, _ := host.NewHostN("127.0.0.1:31337", gen.Reference())
	receiver, _ := host.NewHostN("127.0.0.2:31338", gen.Reference())

	result := NewPacket(sender, receiver, types.RPC, 123)
	result.TraceID = "d6b44f62-7b5e-4249-90c7-ccae194a5baa"
	return result
}

func TestSerializePacket(t *testing.T) {
	msg := testRPCPacket()
	msg.SetRequest(&RPCRequest{Method: "test", Data: []byte{0, 1, 2, 3}})

	_, err := SerializePacket(msg)

	require.NoError(t, err)
}

func TestDeserializePacket(t *testing.T) {
	msg := testRPCPacket()
	msg.SetRequest(&RPCRequest{Method: "test", Data: []byte{0, 1, 2, 3}})

	serialized, _ := SerializePacket(msg)

	var buffer bytes.Buffer

	buffer.Write(serialized)

	deserialized, err := DeserializePacket(log.GlobalLogger, &buffer)

	require.NoError(t, err)
	require.Equal(t, deserialized.Packet, msg)
	require.Equal(t, deserialized.Bytes(), serialized)
}

func TestDeserializeBigPacket(t *testing.T) {
	data := make([]byte, 1024*1024*10)
	rand.Read(data)

	msg := testRPCPacket()
	msg.SetRequest(&RPCRequest{Method: "test", Data: data})

	serialized, err := SerializePacket(msg)
	require.NoError(t, err)

	var buffer bytes.Buffer
	buffer.Write(serialized)

	deserializedMsg, err := DeserializePacket(log.GlobalLogger, &buffer)
	require.NoError(t, err)

	deserializedData := deserializedMsg.GetRequest().GetRPC().Data
	require.EqualValues(t, data, deserializedData)
}

type PacketSuite struct {
	suite.Suite
	sender *host.Host
	packet *Packet
}

func (s *PacketSuite) TestGetType() {
	s.Equal(s.packet.GetType(), types.RPC)
}

func (s *PacketSuite) TestGetData() {
	s.EqualValues(s.packet.GetRequest().GetRPC().Data, []byte{0, 1, 2, 3})
}

func (s *PacketSuite) TestGetRequestID() {
	s.EqualValues(s.packet.GetRequestID(), 123)
}

func TestPacketMethods(t *testing.T) {
	p := testRPCPacket()
	p.SetRequest(&RPCRequest{Method: "test", Data: []byte{0, 1, 2, 3}})

	suite.Run(t, &PacketSuite{
		sender: p.Sender,
		packet: p,
	})
}

func marshalUnmarshal(t *testing.T, p1, p2 *Packet) {
	data, err := p1.Marshal()
	require.NoError(t, err)
	err = p2.Unmarshal(data)
	require.NoError(t, err)
}

func marshalUnmarshalPacketRequest(t *testing.T, request interface{}) (p1, p2 *Packet) {
	p1, p2 = &Packet{}, &Packet{}
	p1.SetRequest(request)
	marshalUnmarshal(t, p1, p2)
	require.NotNil(t, p2.GetRequest())
	return p1, p2
}

func marshalUnmarshalPacketResponse(t *testing.T, response interface{}) (p1, p2 *Packet) {
	p1, p2 = &Packet{}, &Packet{}
	p1.SetResponse(response)
	marshalUnmarshal(t, p1, p2)
	require.NotNil(t, p2.GetResponse())
	return p1, p2
}

func TestPacket_SetRequest(t *testing.T) {
	type SomeData struct {
		someField int
	}
	p := Packet{}
	f := func() {
		p.SetRequest(&SomeData{})
	}
	assert.Panics(t, f)
}

func TestPacket_SetResponse(t *testing.T) {
	type SomeData struct {
		someField int
	}
	p := Packet{}
	f := func() {
		p.SetResponse(&SomeData{})
	}
	assert.Panics(t, f)
}

func TestPacket_GetRequest_GetRPC(t *testing.T) {
	rpc := RPCRequest{Method: "meth", Data: []byte("123")}
	p1, p2 := marshalUnmarshalPacketRequest(t, &rpc)
	require.NotNil(t, p2.GetRequest().GetRPC())
	assert.Equal(t, p1.GetRequest().GetRPC().Method, p2.GetRequest().GetRPC().Method)
	assert.Equal(t, p1.GetRequest().GetRPC().Data, p2.GetRequest().GetRPC().Data)
}

func TestPacket_GetRequest_GetAuthorize(t *testing.T) {
	ss := []byte("onetwothree")
	sign := []byte("abcdefg")
	auth := AuthorizeRequest{AuthorizeData: &AuthorizeData{Certificate: ss, Version: "ver1"}, Signature: sign}
	_, p2 := marshalUnmarshalPacketRequest(t, &auth)
	require.NotNil(t, p2.GetRequest().GetAuthorize())
	require.NotNil(t, p2.GetRequest().GetAuthorize().AuthorizeData)

	assert.Equal(t, ss, p2.GetRequest().GetAuthorize().AuthorizeData.Certificate)
	assert.Equal(t, sign, p2.GetRequest().GetAuthorize().Signature)
	assert.Equal(t, "ver1", p2.GetRequest().GetAuthorize().AuthorizeData.Version)
}

func TestPacket_GetResponse(t *testing.T) {
	response := BasicResponse{}
	_, p2 := marshalUnmarshalPacketResponse(t, &response)
	assert.NotNil(t, p2.GetResponse().GetBasic())
}

func TestPacket_Marshal_0x80(t *testing.T) {
	p := testRPCPacket()
	data, err := p.Marshal()
	require.NoError(t, err)
	assert.EqualValues(t, 0x80, data[0])
}
