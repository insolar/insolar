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

	deserialized, _, err := DeserializePacket(log.GlobalLogger(), &buffer)

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

	deserializedMsg, _, err := DeserializePacket(log.GlobalLogger(), &buffer)
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
