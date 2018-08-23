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

package packet

import (
	"errors"
	"testing"

	"github.com/insolar/insolar/network/host/id"
	"github.com/insolar/insolar/network/host/node"
	"github.com/stretchr/testify/assert"
)

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder()

	assert.Equal(t, builder, Builder{})
	assert.Empty(t, builder.actions)
}

func TestBuilder_Build_EmptyPacket(t *testing.T) {
	builder := NewBuilder()

	assert.Equal(t, builder.Build(), &Packet{})
}

func TestBuilder_Build_RequestPacket(t *testing.T) {
	builder := NewBuilder()
	senderAddress, _ := node.NewAddress("127.0.0.1:31337")
	sender := node.NewNode(senderAddress)
	sender.ID, _ = id.NewID(id.GetRandomKey())
	receiverAddress, _ := node.NewAddress("127.0.0.2:31338")
	receiver := node.NewNode(receiverAddress)
	receiver.ID, _ = id.NewID(id.GetRandomKey())

	m := builder.Sender(sender).Receiver(receiver).Type(TypeRPC).Request(&RequestDataRPC{"test", [][]byte{}}).Build()

	expectedPacket := &Packet{
		Sender:     sender,
		Receiver:   receiver,
		Type:       TypeRPC,
		Data:       &RequestDataRPC{"test", [][]byte{}},
		IsResponse: false,
		Error:      nil,
	}
	assert.Equal(t, expectedPacket, m)
}

func TestBuilder_Build_ResponsePacket(t *testing.T) {
	builder := NewBuilder()
	senderAddress, _ := node.NewAddress("127.0.0.1:31337")
	sender := node.NewNode(senderAddress)
	sender.ID, _ = id.NewID(id.GetRandomKey())
	receiverAddress, _ := node.NewAddress("127.0.0.2:31338")
	receiver := node.NewNode(receiverAddress)
	receiver.ID, _ = id.NewID(id.GetRandomKey())

	m := builder.Sender(sender).Receiver(receiver).Type(TypeRPC).Response(&ResponseDataRPC{true, []byte("ok"), ""}).Build()

	expectedPacket := &Packet{
		Sender:     sender,
		Receiver:   receiver,
		Type:       TypeRPC,
		Data:       &ResponseDataRPC{true, []byte("ok"), ""},
		IsResponse: true,
		Error:      nil,
	}
	assert.Equal(t, expectedPacket, m)
}

func TestBuilder_Build_ErrorPacket(t *testing.T) {
	builder := NewBuilder()
	senderAddress, _ := node.NewAddress("127.0.0.1:31337")
	sender := node.NewNode(senderAddress)
	sender.ID, _ = id.NewID(id.GetRandomKey())
	receiverAddress, _ := node.NewAddress("127.0.0.2:31338")
	receiver := node.NewNode(receiverAddress)
	receiver.ID, _ = id.NewID(id.GetRandomKey())

	m := builder.Sender(sender).Receiver(receiver).Type(TypeRPC).Response(&ResponseDataRPC{}).Error(errors.New("test error")).Build()

	expectedPacket := &Packet{
		Sender:     sender,
		Receiver:   receiver,
		Type:       TypeRPC,
		Data:       &ResponseDataRPC{},
		IsResponse: true,
		Error:      errors.New("test error"),
	}
	assert.Equal(t, expectedPacket, m)
}
