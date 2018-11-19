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

package transport

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/stretchr/testify/require"
)

func TestNewFuture(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(packet.RequestID(1), n, m, cb)

	require.Implements(t, (*Future)(nil), f)
}

func TestFuture_ID(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(packet.RequestID(1), n, m, cb)

	require.Equal(t, f.ID(), packet.RequestID(1))
}

func TestFuture_Actor(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(packet.RequestID(1), n, m, cb)

	require.Equal(t, f.Actor(), n)
}

func TestFuture_Result(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(packet.RequestID(1), n, m, cb)

	require.Empty(t, f.Result())
}

func TestFuture_Request(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(packet.RequestID(1), n, m, cb)

	require.Equal(t, f.Request(), m)
}

func TestFuture_SetResult(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(packet.RequestID(1), n, m, cb)

	require.Empty(t, f.Result())

	go f.SetResult(m)

	m2 := <-f.Result()

	require.Equal(t, m, m2)

	go f.SetResult(m)

	m3, err := f.GetResult(10 * time.Millisecond)
	require.NoError(t, err)
	require.Equal(t, m, m3)

	// no result, timeout
	_, err = f.GetResult(10 * time.Millisecond)
	require.Error(t, err)
}

func TestFuture_Cancel(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")

	cbCalled := false

	cb := func(f Future) { cbCalled = true }

	m := &packet.Packet{}
	f := NewFuture(packet.RequestID(1), n, m, cb)

	f.Cancel()

	_, closed := <-f.Result()

	require.False(t, closed)
	require.True(t, cbCalled)
}

func TestFuture_GetResult(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	m := &packet.Packet{}
	var cancelled uint32 = 0
	cancelCallback := func(f Future) {
		atomic.StoreUint32(&cancelled, 1)
	}
	f := NewFuture(packet.RequestID(1), n, m, cancelCallback)
	go func() {
		time.Sleep(time.Millisecond)
		f.Cancel()
	}()

	_, err := f.GetResult(10 * time.Millisecond)
	require.Error(t, err)
	require.Equal(t, uint32(1), atomic.LoadUint32(&cancelled))
}

func TestFuture_GetResult2(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	c := make(chan *packet.Packet)
	var f Future = &future{
		result:         c,
		actor:          n,
		request:        &packet.Packet{},
		requestID:      packet.RequestID(1),
		cancelCallback: func(f Future) {},
	}
	go func() {
		time.Sleep(time.Millisecond)
		close(c)
	}()
	_, err := f.GetResult(10 * time.Millisecond)
	require.Error(t, err)
}
