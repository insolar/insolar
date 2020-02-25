// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package future

import (
	"sync/atomic"
	"time"

	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
)

type future struct {
	response       chan network.ReceivedPacket
	receiver       *host.Host
	request        *packet.Packet
	requestID      types.RequestID
	cancelCallback CancelCallback
	finished       uint32
}

// NewFuture creates a new Future.
func NewFuture(requestID types.RequestID, receiver *host.Host, packet *packet.Packet, cancelCallback CancelCallback) Future {
	metrics.NetworkFutures.WithLabelValues(packet.GetType().String()).Inc()
	return &future{
		response:       make(chan network.ReceivedPacket, 1),
		receiver:       receiver,
		request:        packet,
		requestID:      requestID,
		cancelCallback: cancelCallback,
	}
}

// ID returns RequestID of packet.
func (f *future) ID() types.RequestID {
	return f.requestID
}

// Receiver returns Host address that was used to create packet.
func (f *future) Receiver() *host.Host {
	return f.receiver
}

// Request returns original request packet.
func (f *future) Request() network.Packet {
	return f.request
}

// Response returns response packet channel.
func (f *future) Response() <-chan network.ReceivedPacket {
	return f.response
}

// SetResponse write packet to the response channel.
func (f *future) SetResponse(response network.ReceivedPacket) {
	if atomic.CompareAndSwapUint32(&f.finished, 0, 1) {
		f.response <- response
		f.finish()
	}
}

// WaitResponse gets the future response from Response() channel with a timeout set to `duration`.
func (f *future) WaitResponse(duration time.Duration) (network.ReceivedPacket, error) {
	select {
	case response, ok := <-f.Response():
		if !ok {
			return nil, ErrChannelClosed
		}
		return response, nil
	case <-time.After(duration):
		f.Cancel()
		metrics.NetworkPacketTimeoutTotal.WithLabelValues(f.request.GetType().String()).Inc()
		return nil, ErrTimeout
	}
}

// Cancel cancels Future processing.
// Please note that cancelCallback is called asynchronously. In other words it's not guaranteed
// it was called and finished when WaitResponse() returns ErrChannelClosed.
func (f *future) Cancel() {
	if atomic.CompareAndSwapUint32(&f.finished, 0, 1) {
		f.finish()
	}
}

func (f *future) finish() {
	close(f.response)
	f.cancelCallback(f)
}
