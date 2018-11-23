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
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestBuilder_Build_EmptyPacket(t *testing.T) {
	builder := NewBuilder(nil)

	require.Equal(t, builder.Build(), &Packet{})
}

func TestBuilder_Build_RequestPacket(t *testing.T) {
	sender, _ := host.NewHostN("127.0.0.1:31337", testutils.RandomRef())
	receiver, _ := host.NewHostN("127.0.0.2:31338", testutils.RandomRef())
	builder := NewBuilder(sender)
	m := builder.Receiver(receiver).Type(TestPacket).Request(&RequestTest{[]byte{0, 1, 2, 3}}).Build()

	expectedPacket := &Packet{
		Sender:     sender,
		Receiver:   receiver,
		Type:       TestPacket,
		Data:       &RequestTest{[]byte{0, 1, 2, 3}},
		IsResponse: false,
		Error:      nil,
	}
	require.Equal(t, expectedPacket, m)
}

func TestBuilder_Build_ResponsePacket(t *testing.T) {
	sender, _ := host.NewHostN("127.0.0.1:31337", testutils.RandomRef())
	receiver, _ := host.NewHostN("127.0.0.2:31338", testutils.RandomRef())
	builder := NewBuilder(sender)
	m := builder.Receiver(receiver).Type(TestPacket).Response(&ResponseTest{42}).Build()

	expectedPacket := &Packet{
		Sender:     sender,
		Receiver:   receiver,
		Type:       TestPacket,
		Data:       &ResponseTest{42},
		IsResponse: true,
		Error:      nil,
	}
	require.Equal(t, expectedPacket, m)
}

func TestBuilder_Build_ErrorPacket(t *testing.T) {
	sender, _ := host.NewHostN("127.0.0.1:31337", testutils.RandomRef())
	receiver, _ := host.NewHostN("127.0.0.2:31338", testutils.RandomRef())
	builder := NewBuilder(sender)
	m := builder.Receiver(receiver).Type(TestPacket).Response(&ResponseTest{}).Error(errors.New("test error")).Build()

	expectedPacket := &Packet{
		Sender:     sender,
		Receiver:   receiver,
		Type:       TestPacket,
		Data:       &ResponseTest{},
		IsResponse: true,
		Error:      errors.New("test error"),
	}
	require.Equal(t, expectedPacket, m)
}
