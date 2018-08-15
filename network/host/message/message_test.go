/*
 *    Copyright 2018 INS Ecosystem
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

package message

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/insolar/insolar/network/host/node"
	"github.com/stretchr/testify/assert"
)

func TestNewPingMessage(t *testing.T) {
	senderAddress, _ := node.NewAddress("127.0.0.1:31337")
	sender := node.NewNode(senderAddress)
	sender.ID, _ = node.NewID()
	receiverAddress, _ := node.NewAddress("127.0.0.2:31338")
	receiver := node.NewNode(receiverAddress)
	receiver.ID, _ = node.NewID()

	m := NewPingMessage(sender, receiver)

	expectedMessage := &Message{
		Sender:   sender,
		Receiver: receiver,
		Type:     TypePing,
	}
	assert.Equal(t, expectedMessage, m)
}

func TestMessage_IsValid(t *testing.T) {
	builder := NewBuilder()

	correctMessage := builder.Type(TypeRPC).Request(&RequestDataRPC{"test", [][]byte{}}).Build()
	assert.True(t, correctMessage.IsValid())

	badtMessage := builder.Type(TypeStore).Request(&RequestDataRPC{"test", [][]byte{}}).Build()
	assert.False(t, badtMessage.IsValid())
}

func TestMessage_IsValid_Ok(t *testing.T) {
	tests := []struct {
		name        string
		messageType messageType
		data        interface{}
	}{
		{"TypePing", TypePing, nil},
		{"TypeFindNode", TypeFindNode, &RequestDataFindNode{}},
		{"TypeFindValue", TypeFindValue, &RequestDataFindValue{}},
		{"TypeStore", TypeStore, &RequestDataStore{}},
		{"TypeRPC", TypeRPC, &RequestDataRPC{"test", [][]byte{}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := NewBuilder()
			message := builder.Type(test.messageType).Request(test.data).Build()
			assert.True(t, message.IsValid())
		})
	}
}

func TestMessage_IsValid_Fail(t *testing.T) {
	tests := []struct {
		name        string
		messageType messageType
		data        interface{}
	}{
		{"incorrect request", TypeStore, &RequestDataRPC{"test", [][]byte{}}},
		{"incorrect type", messageType(1337), &RequestDataFindNode{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			builder := NewBuilder()
			message := builder.Type(test.messageType).Request(test.data).Build()
			assert.False(t, message.IsValid())
		})
	}
}

func TestMessage_IsForMe(t *testing.T) {
	senderAddress, _ := node.NewAddress("127.0.0.1:31337")
	sender := node.NewNode(senderAddress)
	sender.ID, _ = node.NewID()
	receiverAddress, _ := node.NewAddress("127.0.0.2:31338")
	receiver := node.NewNode(receiverAddress)
	receiver.ID, _ = node.NewID()
	builder := NewBuilder()
	origin, _ := node.NewOrigin([]node.ID{receiver.ID}, receiver.Address)

	myMessage := builder.Receiver(receiver).Build()
	notMyMessage := builder.Receiver(sender).Build()

	assert.True(t, myMessage.IsForMe(*origin))
	assert.False(t, notMyMessage.IsForMe(*origin))
}

func TestSerializeMessage(t *testing.T) {
	senderAddress, _ := node.NewAddress("127.0.0.1:31337")
	sender := node.NewNode(senderAddress)
	sender.ID, _ = node.NewID()
	receiverAddress, _ := node.NewAddress("127.0.0.2:31338")
	receiver := node.NewNode(receiverAddress)
	receiver.ID, _ = node.NewID()
	builder := NewBuilder()
	msg := builder.Sender(sender).Receiver(receiver).Type(TypeFindNode).Request(&RequestDataFindNode{receiver.ID}).Build()

	_, err := SerializeMessage(msg)

	assert.NoError(t, err)
}

func TestDeserializeMessage(t *testing.T) {
	senderAddress, _ := node.NewAddress("127.0.0.1:31337")
	sender := node.NewNode(senderAddress)
	sender.ID, _ = node.NewID()
	receiverAddress, _ := node.NewAddress("127.0.0.2:31338")
	receiver := node.NewNode(receiverAddress)
	receiver.ID, _ = node.NewID()
	builder := NewBuilder()
	msg := builder.Sender(sender).Receiver(receiver).Type(TypeFindNode).Request(&RequestDataFindNode{receiver.ID}).Build()

	serialized, _ := SerializeMessage(msg)

	var buffer bytes.Buffer

	buffer.Write(serialized)

	deserialized, err := DeserializeMessage(&buffer)

	assert.NoError(t, err)
	assert.Equal(t, deserialized, msg)
}

func TestDeserializeBigMessage(t *testing.T) {
	address, _ := node.NewAddress("127.0.0.1:31337")
	nodeOne := node.NewNode(address)

	data := make([]byte, 1024*1024*10)
	rand.Read(data)

	builder := NewBuilder()
	msg := builder.Sender(nodeOne).Receiver(nodeOne).Type(TypeStore).Request(&RequestDataStore{data, true}).Build()
	assert.True(t, msg.IsValid())

	serialized, err := SerializeMessage(msg)
	assert.NoError(t, err)

	var buffer bytes.Buffer
	buffer.Write(serialized)

	deserializedMsg, err := DeserializeMessage(&buffer)
	assert.NoError(t, err)

	deserializedData := deserializedMsg.Data.(*RequestDataStore).Data
	assert.Equal(t, data, deserializedData)
}
