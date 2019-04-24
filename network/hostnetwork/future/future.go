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
	"errors"
	"sync/atomic"
	"time"

	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
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

	GetRequest() network.Request
	Response() <-chan network.Response
	GetResponse(duration time.Duration) (network.Response, error)
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

// NewFuture creates a new Future.
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
func (f *future) ID() network.RequestID {
	return f.requestID
}

// Actor returns Host address that was used to create packet.
func (f *future) Actor() *host.Host {
	return f.actor
}

// Request returns original request packet.
func (f *future) Request() *packet.Packet {
	return f.request
}

// Result returns result packet channel.
func (f *future) Result() <-chan *packet.Packet {
	return f.result
}

// SetResult write packet to the result channel.
func (f *future) SetResult(msg *packet.Packet) {
	if atomic.CompareAndSwapUint32(&f.finished, 0, 1) {
		f.result <- msg
		f.finish()
	}
}

// GetResult gets the future result from Result() channel with a timeout set to `duration`.
func (f *future) GetResult(duration time.Duration) (*packet.Packet, error) {
	select {
	case result, ok := <-f.Result():
		if !ok {
			return nil, ErrChannelClosed
		}
		return result, nil
	case <-time.After(duration):
		f.Cancel()
		metrics.NetworkPacketTimeoutTotal.WithLabelValues(f.request.Type.String()).Inc()
		return nil, ErrTimeout
	}
}

// Cancel cancels Future processing.
// Please note that cancelCallback is called asynchronously. In other words it's not guaranteed
// it was called and finished when GetResult() returns ErrChannelClosed.
func (f *future) Cancel() {
	if atomic.CompareAndSwapUint32(&f.finished, 0, 1) {
		f.finish()
		metrics.NetworkFutures.WithLabelValues(f.request.Type.String()).Dec()
	}
}

func (f *future) finish() {
	close(f.result)
	f.cancelCallback(f)
}

// Response get channel that receives response to sent request
func (f *future) Response() <-chan network.Response {
	in := f.Result()
	out := make(chan network.Response, cap(in))
	go func(in <-chan *packet.Packet, out chan<- network.Response) {
		for packet := range in {
			out <- packet
		}
		close(out)
	}(in, out)
	return out
}

// GetResponse get response to sent request with `duration` timeout
func (f *future) GetResponse(duration time.Duration) (network.Response, error) {
	result, err := f.GetResult(duration)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetRequest get initiating request.
func (f *future) GetRequest() network.Request {
	request := f.Request()
	return request
}
