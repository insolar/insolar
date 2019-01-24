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
