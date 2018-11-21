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
	"errors"
	"testing"

	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/id"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Build_EmptyPacket(t *testing.T) {
	builder := NewBuilder(nil)

	assert.Equal(t, builder.Build(), &Packet{})
}

func TestBuilder_Build_RequestPacket(t *testing.T) {
	senderAddress, _ := host.NewAddress("127.0.0.1:31337")
	sender := host.NewHost(senderAddress)
	sender.ID, _ = id.NewID()
	builder := NewBuilder(sender)
	receiverAddress, _ := host.NewAddress("127.0.0.2:31338")
	receiver := host.NewHost(receiverAddress)
	receiver.ID, _ = id.NewID()
	ref := testutils.RandomRef()

	m := builder.Receiver(receiver).Type(types.TypeRPC).Request(&RequestDataRPC{ref, "test", [][]byte{}}).Build()

	expectedPacket := &Packet{
		Sender:     sender,
		Receiver:   receiver,
		Type:       types.TypeRPC,
		Data:       &RequestDataRPC{ref, "test", [][]byte{}},
		IsResponse: false,
		Error:      nil,
	}
	assert.Equal(t, expectedPacket, m)
}

func TestBuilder_Build_ResponsePacket(t *testing.T) {
	senderAddress, _ := host.NewAddress("127.0.0.1:31337")
	sender := host.NewHost(senderAddress)
	sender.ID, _ = id.NewID()
	builder := NewBuilder(sender)
	receiverAddress, _ := host.NewAddress("127.0.0.2:31338")
	receiver := host.NewHost(receiverAddress)
	receiver.ID, _ = id.NewID()

	m := builder.Receiver(receiver).Type(types.TypeRPC).Response(&ResponseDataRPC{true, []byte("ok"), ""}).Build()

	expectedPacket := &Packet{
		Sender:     sender,
		Receiver:   receiver,
		Type:       types.TypeRPC,
		Data:       &ResponseDataRPC{true, []byte("ok"), ""},
		IsResponse: true,
		Error:      nil,
	}
	assert.Equal(t, expectedPacket, m)
}

func TestBuilder_Build_ErrorPacket(t *testing.T) {
	senderAddress, _ := host.NewAddress("127.0.0.1:31337")
	sender := host.NewHost(senderAddress)
	sender.ID, _ = id.NewID()
	builder := NewBuilder(sender)
	receiverAddress, _ := host.NewAddress("127.0.0.2:31338")
	receiver := host.NewHost(receiverAddress)
	receiver.ID, _ = id.NewID()

	m := builder.Receiver(receiver).Type(types.TypeRPC).Response(&ResponseDataRPC{}).Error(errors.New("test error")).Build()

	expectedPacket := &Packet{
		Sender:     sender,
		Receiver:   receiver,
		Type:       types.TypeRPC,
		Data:       &ResponseDataRPC{},
		IsResponse: true,
		Error:      errors.New("test error"),
	}
	assert.Equal(t, expectedPacket, m)
}
