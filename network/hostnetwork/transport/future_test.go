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

package transport

import (
	"testing"

	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/stretchr/testify/assert"
)

func TestNewFuture(t *testing.T) {
	addr, _ := host.NewAddress("127.0.0.1:8080")
	n := host.NewHost(addr)
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(packet.RequestID(1), n, m, cb)

	assert.Implements(t, (*Future)(nil), f)
}

func TestFuture_ID(t *testing.T) {
	addr, _ := host.NewAddress("127.0.0.1:8080")
	n := host.NewHost(addr)
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(packet.RequestID(1), n, m, cb)

	assert.Equal(t, f.ID(), packet.RequestID(1))
}

func TestFuture_Actor(t *testing.T) {
	addr, _ := host.NewAddress("127.0.0.1:8080")
	n := host.NewHost(addr)
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(packet.RequestID(1), n, m, cb)

	assert.Equal(t, f.Actor(), n)
}

func TestFuture_Result(t *testing.T) {
	addr, _ := host.NewAddress("127.0.0.1:8080")
	n := host.NewHost(addr)
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(packet.RequestID(1), n, m, cb)

	assert.Empty(t, f.Result())
}

func TestFuture_Request(t *testing.T) {
	addr, _ := host.NewAddress("127.0.0.1:8080")
	n := host.NewHost(addr)
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(packet.RequestID(1), n, m, cb)

	assert.Equal(t, f.Request(), m)
}

func TestFuture_SetResult(t *testing.T) {
	addr, _ := host.NewAddress("127.0.0.1:8080")
	n := host.NewHost(addr)
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(packet.RequestID(1), n, m, cb)

	assert.Empty(t, f.Result())

	go f.SetResult(m)

	m2 := <-f.Result()

	assert.Equal(t, m, m2)
}

func TestFuture_Cancel(t *testing.T) {
	addr, _ := host.NewAddress("127.0.0.1:8080")
	n := host.NewHost(addr)

	cbCalled := false

	cb := func(f Future) { cbCalled = true }

	m := &packet.Packet{}
	f := NewFuture(packet.RequestID(1), n, m, cb)

	f.Cancel()

	_, closed := <-f.Result()

	assert.False(t, closed)
	assert.True(t, cbCalled)
}
