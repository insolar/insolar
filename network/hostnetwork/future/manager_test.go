/*
 *    Copyright 2019 Insolar Technologies
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

package future

import (
	"testing"

	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	m := NewManager()

	require.IsType(t, m, &futureManager{})
}

func TestFutureManager_Create(t *testing.T) {
	m := NewManager()

	sender, _ := host.NewHostN("127.0.0.1:31337", testutils.RandomRef())
	receiver, _ := host.NewHostN("127.0.0.2:31338", testutils.RandomRef())
	builder := packet.NewBuilder(sender)
	p := builder.
		Receiver(receiver).
		Type(packet.TestPacket).
		Request(&packet.RequestTest{[]byte{0, 1, 2, 3}}).
		RequestID(network.RequestID(123)).
		Build()

	future := m.Create(p)

	require.Equal(t, future.ID(), p.RequestID)
	require.Equal(t, future.Request(), p)
	require.Equal(t, future.Receiver(), receiver)
	require.Equal(t, future.Request(), p)
}

func TestFutureManager_Get(t *testing.T) {
	m := NewManager()

	sender, _ := host.NewHostN("127.0.0.1:31337", testutils.RandomRef())
	receiver, _ := host.NewHostN("127.0.0.2:31338", testutils.RandomRef())
	builder := packet.NewBuilder(sender)
	p := builder.
		Receiver(receiver).
		Type(packet.TestPacket).
		Request(&packet.RequestTest{[]byte{0, 1, 2, 3}}).
		RequestID(network.RequestID(123)).
		Build()

	require.Nil(t, m.Get(p))

	expectedFuture := m.Create(p)
	actualFuture := m.Get(p)

	require.Equal(t, expectedFuture, actualFuture)
}

func TestFutureManager_Canceler(t *testing.T) {
	m := NewManager()

	sender, _ := host.NewHostN("127.0.0.1:31337", testutils.RandomRef())
	receiver, _ := host.NewHostN("127.0.0.2:31338", testutils.RandomRef())
	builder := packet.NewBuilder(sender)
	p := builder.
		Receiver(receiver).
		Type(packet.TestPacket).
		Request(&packet.RequestTest{[]byte{0, 1, 2, 3}}).
		RequestID(network.RequestID(123)).
		Build()

	future := m.Create(p)
	require.NotNil(t, future)

	future.Cancel()

	require.Nil(t, m.Get(p))
}
