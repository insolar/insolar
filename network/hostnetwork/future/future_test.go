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

	receivedPacket := packet.NewReceivedPacket(m, nil)
	go f.SetResponse(receivedPacket)

	m2 := <-f.Response() // Response() call closes channel

	require.Equal(t, receivedPacket, m2)

	m3, err := f.WaitResponse(time.Minute)
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
	c := make(chan network.ReceivedPacket)
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
	c := make(chan network.ReceivedPacket)
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
	c := make(chan network.ReceivedPacket, 1)
	var f Future = &future{
		response:       c,
		receiver:       n,
		request:        &packet.Packet{},
		requestID:      types.RequestID(1),
		cancelCallback: func(f Future) {},
	}

	p := packet.NewReceivedPacket(&packet.Packet{}, nil)
	c <- p

	res, err := f.WaitResponse(time.Minute)
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
		f.SetResponse(packet.NewReceivedPacket(&packet.Packet{}, nil))
		wg.Done()
	}()

	wg.Wait()
	res, ok := <-f.Response()

	cancelDone := res == nil && !ok
	resultDone := res != nil && ok

	require.True(t, cancelDone || resultDone)
	require.True(t, cbCalled)
}
