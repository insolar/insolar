//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package future

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
)

func TestNewFuture(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(network.RequestID(1), n, m, cb)

	require.Implements(t, (*Future)(nil), f)
}

func TestFuture_ID(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(network.RequestID(1), n, m, cb)

	require.Equal(t, f.ID(), network.RequestID(1))
}

func TestFuture_Actor(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(network.RequestID(1), n, m, cb)

	require.Equal(t, f.Actor(), n)
}

func TestFuture_Result(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(network.RequestID(1), n, m, cb)

	require.Empty(t, f.Result())
}

func TestFuture_Request(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(network.RequestID(1), n, m, cb)

	require.Equal(t, f.Request(), m)
}

func TestFuture_SetResult(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(network.RequestID(1), n, m, cb)

	require.Empty(t, f.Result())

	go f.SetResult(m)

	m2 := <-f.Result() // Result() call closes channel

	require.Equal(t, m, m2)

	m3, err := f.GetResult(10 * time.Millisecond)
	// legal behavior, the channel is closed because of the previous f.Result() call finished the Future
	require.EqualError(t, err, "channel closed")
	require.Nil(t, m3)
}

func TestFuture_Cancel(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")

	cbCalled := false

	cb := func(f Future) { cbCalled = true }

	m := &packet.Packet{}
	f := NewFuture(network.RequestID(1), n, m, cb)

	f.Cancel()

	_, closed := <-f.Result()

	require.False(t, closed)
	require.True(t, cbCalled)
}

func TestFuture_GetResult(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	m := &packet.Packet{}
	canceled := make(chan bool, 1)
	cancelCallback := func(f Future) {
		canceled <- true
	}
	f := NewFuture(network.RequestID(1), n, m, cancelCallback)
	_, err := f.GetResult(time.Millisecond)
	require.Error(t, err)
	require.Equal(t, ErrTimeout, err)
	tmp := <-canceled
	require.Equal(t, true, tmp)
}

func TestFuture_GetResult2(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	c := make(chan *packet.Packet)
	var f Future = &future{
		result:         c,
		actor:          n,
		request:        &packet.Packet{},
		requestID:      network.RequestID(1),
		cancelCallback: func(f Future) {},
	}
	go func() {
		time.Sleep(time.Millisecond)
		f.Cancel()
	}()
	_, err := f.GetResult(1000 * time.Millisecond)
	require.Error(t, err)
}

func TestFuture_SetResult_Cancel_Concurrency(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")

	cbCalled := false

	cb := func(f Future) { cbCalled = true }

	m := &packet.Packet{}
	f := NewFuture(network.RequestID(1), n, m, cb)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		f.Cancel()
		wg.Done()
	}()
	go func() {
		f.SetResult(&packet.Packet{})
		wg.Done()
	}()

	wg.Wait()
	res, ok := <-f.Result()

	cancelDone := res == nil && !ok
	resultDone := res != nil && ok

	require.True(t, cancelDone || resultDone)
	require.True(t, cbCalled)
}
