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

func TestBuilder_Build_EmptyMessage(t *testing.T) {
	builder := NewBuilder()

	assert.Equal(t, builder.Build(), &Message{})
}

func TestBuilder_Build_RequestMessage(t *testing.T) {
	builder := NewBuilder()
	senderAddress, _ := node.NewAddress("127.0.0.1:31337")
	sender := node.NewNode(senderAddress)
	sender.ID, _ = id.NewID([]byte{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106})
	receiverAddress, _ := node.NewAddress("127.0.0.2:31338")
	receiver := node.NewNode(receiverAddress)
	receiver.ID, _ = id.NewID([]byte{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106})

	m := builder.Sender(sender).Receiver(receiver).Type(TypeRPC).Request(&RequestDataRPC{"test", [][]byte{}}).Build()

	expectedMessage := &Message{
		Sender:     sender,
		Receiver:   receiver,
		Type:       TypeRPC,
		Data:       &RequestDataRPC{"test", [][]byte{}},
		IsResponse: false,
		Error:      nil,
	}
	assert.Equal(t, expectedMessage, m)
}

func TestBuilder_Build_ResponseMessage(t *testing.T) {
	builder := NewBuilder()
	senderAddress, _ := node.NewAddress("127.0.0.1:31337")
	sender := node.NewNode(senderAddress)
	sender.ID, _ = id.NewID([]byte{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106})
	receiverAddress, _ := node.NewAddress("127.0.0.2:31338")
	receiver := node.NewNode(receiverAddress)
	receiver.ID, _ = id.NewID([]byte{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106})

	m := builder.Sender(sender).Receiver(receiver).Type(TypeRPC).Response(&ResponseDataRPC{true, []byte("ok"), ""}).Build()

	expectedMessage := &Message{
		Sender:     sender,
		Receiver:   receiver,
		Type:       TypeRPC,
		Data:       &ResponseDataRPC{true, []byte("ok"), ""},
		IsResponse: true,
		Error:      nil,
	}
	assert.Equal(t, expectedMessage, m)
}

func TestBuilder_Build_ErrorMessage(t *testing.T) {
	builder := NewBuilder()
	senderAddress, _ := node.NewAddress("127.0.0.1:31337")
	sender := node.NewNode(senderAddress)
	sender.ID, _ = id.NewID([]byte{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106})
	receiverAddress, _ := node.NewAddress("127.0.0.2:31338")
	receiver := node.NewNode(receiverAddress)
	receiver.ID, _ = id.NewID([]byte{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106})

	m := builder.Sender(sender).Receiver(receiver).Type(TypeRPC).Response(&ResponseDataRPC{}).Error(errors.New("test error")).Build()

	expectedMessage := &Message{
		Sender:     sender,
		Receiver:   receiver,
		Type:       TypeRPC,
		Data:       &ResponseDataRPC{},
		IsResponse: true,
		Error:      errors.New("test error"),
	}
	assert.Equal(t, expectedMessage, m)
}
