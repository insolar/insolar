/*
 *    Copyright 2018 Insolar
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

package packet

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/id"
	"github.com/stretchr/testify/assert"
)

func TestNewPingPacket(t *testing.T) {
	senderAddress, _ := host.NewAddress("127.0.0.1:31337")
	sender := host.NewHost(senderAddress)
	sender.ID, _ = id.NewID()
	receiverAddress, _ := host.NewAddress("127.0.0.2:31338")
	receiver := host.NewHost(receiverAddress)
	receiver.ID, _ = id.NewID()

	m := NewPingPacket(sender, receiver)

	expectedPacket := &Packet{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypePing,
	}
	assert.Equal(t, expectedPacket, m)
}

func TestPacket_IsValid(t *testing.T) {
	builder := NewBuilder()

	correctPacket := builder.Type(TypeRPC).Request(&RequestDataRPC{"test", [][]byte{}}).Build()
	assert.True(t, correctPacket.IsValid())

	badtPacket := builder.Type(TypeStore).Request(&RequestDataRPC{"test", [][]byte{}}).Build()
	assert.False(t, badtPacket.IsValid())
}

func TestNewAuthPacket(t *testing.T) {
	addr1, _ := host.NewAddress("127.0.0.1:55551")
	addr2, _ := host.NewAddress("127.0.0.1:55552")
	sender := host.NewHost(addr1)
	receiver := host.NewHost(addr2)

	msg := NewAuthPacket(BeginAuth, sender, receiver)
	assert.True(t, msg.IsValid())
}

func TestNewCheckNodePrivPacket(t *testing.T) {
	addr1, _ := host.NewAddress("127.0.0.1:55551")
	addr2, _ := host.NewAddress("127.0.0.1:55552")
	sender := host.NewHost(addr1)
	receiver := host.NewHost(addr2)

	msg := NewCheckNodePrivPacket(sender, receiver, "test string")
	assert.True(t, msg.IsValid())
}

func TestNewCheckOriginPacket(t *testing.T) {
	addr1, _ := host.NewAddress("127.0.0.1:55551")
	addr2, _ := host.NewAddress("127.0.0.1:55552")
	sender := host.NewHost(addr1)
	receiver := host.NewHost(addr2)

	msg := NewCheckOriginPacket(sender, receiver)
	assert.True(t, msg.IsValid())
}

func TestNewKnownOuterHostsPacket(t *testing.T) {
	addr1, _ := host.NewAddress("127.0.0.1:55551")
	addr2, _ := host.NewAddress("127.0.0.1:55552")
	sender := host.NewHost(addr1)
	receiver := host.NewHost(addr2)

	msg := NewKnownOuterHostsPacket(sender, receiver, 1)
	assert.True(t, msg.IsValid())
}

func TestNewObtainIPPacket(t *testing.T) {
	addr1, _ := host.NewAddress("127.0.0.1:55551")
	addr2, _ := host.NewAddress("127.0.0.1:55552")
	sender := host.NewHost(addr1)
	receiver := host.NewHost(addr2)

	msg := NewObtainIPPacket(sender, receiver)
	assert.True(t, msg.IsValid())
}

func TestNewRelayPacket(t *testing.T) {
	addr1, _ := host.NewAddress("127.0.0.1:55551")
	addr2, _ := host.NewAddress("127.0.0.1:55552")
	sender := host.NewHost(addr1)
	receiver := host.NewHost(addr2)

	msg := NewRelayPacket(StartRelay, sender, receiver)
	assert.True(t, msg.IsValid())
}

func TestNewRelayOwnershipPacket(t *testing.T) {
	addr1, _ := host.NewAddress("127.0.0.1:55551")
	addr2, _ := host.NewAddress("127.0.0.1:55552")
	sender := host.NewHost(addr1)
	receiver := host.NewHost(addr2)

	msg := NewRelayOwnershipPacket(sender, receiver, true)
	assert.True(t, msg.IsValid())
}

func TestPacket_IsValid_Ok(t *testing.T) {
	tests := []struct {
		name       string
		packetType packetType
		data       interface{}
	}{
		{"TypePing", TypePing, nil},
		{"TypeFindHost", TypeFindHost, &RequestDataFindHost{}},
		{"TypeFindValue", TypeFindValue, &RequestDataFindValue{}},
		{"TypeStore", TypeStore, &RequestDataStore{}},
		{"TypeRPC", TypeRPC, &RequestDataRPC{"test", [][]byte{}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := NewBuilder()
			packet := builder.Type(test.packetType).Request(test.data).Build()
			assert.True(t, packet.IsValid())
		})
	}
}

func TestPacket_IsValid_Fail(t *testing.T) {
	tests := []struct {
		name       string
		packetType packetType
		data       interface{}
	}{
		{"incorrect request", TypeStore, &RequestDataRPC{"test", [][]byte{}}},
		{"incorrect type", packetType(1337), &RequestDataFindHost{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := NewBuilder()
			packet := builder.Type(test.packetType).Request(test.data).Build()
			assert.False(t, packet.IsValid())
		})
	}
}

func TestPacket_IsForMe(t *testing.T) {
	senderAddress, _ := host.NewAddress("127.0.0.1:31337")
	sender := host.NewHost(senderAddress)
	sender.ID, _ = id.NewID()
	receiverAddress, _ := host.NewAddress("127.0.0.2:31338")
	receiver := host.NewHost(receiverAddress)
	receiver.ID, _ = id.NewID()
	builder := NewBuilder()
	origin, _ := host.NewOrigin([]id.ID{receiver.ID}, receiver.Address)

	myPacket := builder.Receiver(receiver).Build()
	notMyPacket := builder.Receiver(sender).Build()

	assert.True(t, myPacket.IsForMe(*origin))
	assert.False(t, notMyPacket.IsForMe(*origin))
}

func TestSerializePacket(t *testing.T) {
	senderAddress, _ := host.NewAddress("127.0.0.1:31337")
	sender := host.NewHost(senderAddress)
	sender.ID, _ = id.NewID()
	receiverAddress, _ := host.NewAddress("127.0.0.2:31338")
	receiver := host.NewHost(receiverAddress)
	receiver.ID, _ = id.NewID()
	builder := NewBuilder()
	msg := builder.Sender(sender).Receiver(receiver).Type(TypeFindHost).Request(&RequestDataFindHost{receiver.ID.Bytes()}).Build()

	_, err := SerializePacket(msg)

	assert.NoError(t, err)
}

func TestDeserializePacket(t *testing.T) {
	senderAddress, _ := host.NewAddress("127.0.0.1:31337")
	sender := host.NewHost(senderAddress)
	sender.ID, _ = id.NewID()
	receiverAddress, _ := host.NewAddress("127.0.0.2:31338")
	receiver := host.NewHost(receiverAddress)
	receiver.ID, _ = id.NewID()
	builder := NewBuilder()
	msg := builder.Sender(sender).Receiver(receiver).Type(TypeFindHost).Request(&RequestDataFindHost{receiver.ID.Bytes()}).Build()

	serialized, _ := SerializePacket(msg)

	var buffer bytes.Buffer

	buffer.Write(serialized)

	deserialized, err := DeserializePacket(&buffer)

	assert.NoError(t, err)
	assert.Equal(t, deserialized, msg)
}

func TestDeserializeBigPacket(t *testing.T) {
	address, _ := host.NewAddress("127.0.0.1:31337")
	hostOne := host.NewHost(address)

	data := make([]byte, 1024*1024*10)
	rand.Read(data)

	builder := NewBuilder()
	msg := builder.Sender(hostOne).Receiver(hostOne).Type(TypeStore).Request(&RequestDataStore{data, true}).Build()
	assert.True(t, msg.IsValid())

	serialized, err := SerializePacket(msg)
	assert.NoError(t, err)

	var buffer bytes.Buffer
	buffer.Write(serialized)

	deserializedMsg, err := DeserializePacket(&buffer)
	assert.NoError(t, err)

	deserializedData := deserializedMsg.Data.(*RequestDataStore).Data
	assert.Equal(t, data, deserializedData)
}
