/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package transport

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet"
)

var (
	// ErrTimeout is returned when the operation timeout is exceeded.
	ErrTimeout = errors.New("timeout")
	// ErrChannelClosed is returned when the input channel is closed.
	ErrChannelClosed = errors.New("channel closed")
)

// Future is network response future.
type Future interface {

	// ID returns packet sequence number.
	ID() network.RequestID

	// Actor returns the initiator of the packet.
	Actor() *host.Host

	// Request returns origin request.
	Request() *packet.Packet

	// Result is a channel to listen for future result.
	Result() <-chan *packet.Packet

	// SetResult makes packet to appear in result channel.
	SetResult(*packet.Packet)

	// GetResult gets the future result from Result() channel with a timeout set to `duration`.
	GetResult(duration time.Duration) (*packet.Packet, error)

	// Cancel closes all channels and cleans up underlying structures.
	Cancel()
}

// CancelCallback is a callback function executed when cancelling Future.
type CancelCallback func(Future)

type future struct {
	result         chan *packet.Packet
	actor          *host.Host
	request        *packet.Packet
	requestID      network.RequestID
	cancelCallback CancelCallback
	finished       uint32
}

// NewFuture creates new Future.
func NewFuture(requestID network.RequestID, actor *host.Host, msg *packet.Packet, cancelCallback CancelCallback) Future {
	metrics.NetworkFutures.WithLabelValues(msg.Type.String()).Inc()
	return &future{
		result:         make(chan *packet.Packet, 1),
		actor:          actor,
		request:        msg,
		requestID:      requestID,
		cancelCallback: cancelCallback,
	}
}

// ID returns RequestID of packet.
func (future *future) ID() network.RequestID {
	return future.requestID
}

// Actor returns Host address that was used to create packet.
func (future *future) Actor() *host.Host {
	return future.actor
}

// Request returns original request packet.
func (future *future) Request() *packet.Packet {
	return future.request
}

// Result returns result packet channel.
func (future *future) Result() <-chan *packet.Packet {
	return future.result
}

// SetResult write packet to the result channel.
func (future *future) SetResult(msg *packet.Packet) {
	if atomic.CompareAndSwapUint32(&future.finished, 0, 1) {
		future.result <- msg
		future.finish()
	}
}

// GetResult gets the future result from Result() channel with a timeout set to `duration`.
func (future *future) GetResult(duration time.Duration) (*packet.Packet, error) {
	select {
	case result, ok := <-future.Result():
		if !ok {
			return nil, ErrChannelClosed
		}
		return result, nil
	case <-time.After(duration):
		future.Cancel()
		metrics.NetworkPacketTimeoutTotal.WithLabelValues(future.request.Type.String()).Inc()
		return nil, ErrTimeout
	}
}

// Cancel allows to cancel Future processing.
func (future *future) Cancel() {
	if atomic.CompareAndSwapUint32(&future.finished, 0, 1) {
		future.finish()
		metrics.NetworkFutures.WithLabelValues(future.request.Type.String()).Dec()
	}
}

func (future *future) finish() {
	close(future.result)
	future.cancelCallback(future)
}
