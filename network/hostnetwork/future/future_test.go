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
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
)

func TestNewFuture(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(types.RequestID(1), n, m, cb)

	require.Implements(t, (*Future)(nil), f)
}

func TestFuture_ID(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(types.RequestID(1), n, m, cb)

	require.Equal(t, f.ID(), types.RequestID(1))
}

func TestFuture_Actor(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(types.RequestID(1), n, m, cb)

	require.Equal(t, f.Receiver(), n)
}

func TestFuture_Result(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(types.RequestID(1), n, m, cb)

	require.Empty(t, f.Response())
}

func TestFuture_Request(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(types.RequestID(1), n, m, cb)

	require.Equal(t, f.Request(), m)
}

func TestFuture_SetResponse(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	cb := func(f Future) {}
	m := &packet.Packet{}
	f := NewFuture(types.RequestID(1), n, m, cb)

	require.Empty(t, f.Response())

	go f.SetResponse(m)

	m2 := <-f.Response() // Response() call closes channel

	require.Equal(t, m, m2)

	m3, err := f.WaitResponse(10 * time.Millisecond)
	// legal behavior, the channel is closed because of the previous f.Response() call finished the Future
	require.EqualError(t, err, "channel closed")
	require.Nil(t, m3)
}

func TestFuture_Cancel(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")

	cbCalled := false

	cb := func(f Future) { cbCalled = true }

	m := &packet.Packet{}
	f := NewFuture(types.RequestID(1), n, m, cb)

	f.Cancel()

	_, closed := <-f.Response()

	require.False(t, closed)
	require.True(t, cbCalled)
}

func TestFuture_WaitResponse_Cancel(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	c := make(chan network.Packet)
	var f Future = &future{
		response:       c,
		receiver:       n,
		request:        &packet.Packet{},
		requestID:      types.RequestID(1),
		cancelCallback: func(f Future) {},
	}
	go func() {
		time.Sleep(time.Millisecond)
		f.Cancel()
	}()
	_, err := f.WaitResponse(1000 * time.Millisecond)
	require.Error(t, err)
}

func TestFuture_WaitResponse_Timeout(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	c := make(chan network.Packet)
	cancelled := false
	var f Future = &future{
		response:       c,
		receiver:       n,
		request:        &packet.Packet{},
		requestID:      types.RequestID(1),
		cancelCallback: func(f Future) { cancelled = true },
	}
	_, err := f.WaitResponse(time.Millisecond)
	require.Error(t, err)
	require.Equal(t, err, ErrTimeout)
	require.True(t, cancelled)
}

func TestFuture_WaitResponse_Success(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")
	c := make(chan network.Packet, 1)
	var f Future = &future{
		response:       c,
		receiver:       n,
		request:        &packet.Packet{},
		requestID:      types.RequestID(1),
		cancelCallback: func(f Future) {},
	}

	p := &packet.Packet{}
	c <- p

	res, err := f.WaitResponse(time.Millisecond)
	require.NoError(t, err)
	require.Equal(t, res, p)
}

func TestFuture_SetResponse_Cancel_Concurrency(t *testing.T) {
	n, _ := host.NewHost("127.0.0.1:8080")

	cbCalled := false

	cb := func(f Future) { cbCalled = true }

	m := &packet.Packet{}
	f := NewFuture(types.RequestID(1), n, m, cb)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		f.Cancel()
		wg.Done()
	}()
	go func() {
		f.SetResponse(&packet.Packet{})
		wg.Done()
	}()

	wg.Wait()
	res, ok := <-f.Response()

	cancelDone := res == nil && !ok
	resultDone := res != nil && ok

	require.True(t, cancelDone || resultDone)
	require.True(t, cbCalled)
}
