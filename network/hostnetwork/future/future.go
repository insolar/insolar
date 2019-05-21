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
	"sync/atomic"
	"time"

	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
)

type future struct {
	response       chan network.Packet
	receiver       *host.Host
	request        *packet.PacketBackend
	requestID      types.RequestID
	cancelCallback CancelCallback
	finished       uint32
}

// NewFuture creates a new Future.
func NewFuture(requestID types.RequestID, receiver *host.Host, packet *packet.PacketBackend, cancelCallback CancelCallback) Future {
	metrics.NetworkFutures.WithLabelValues(packet.GetType().String()).Inc()
	return &future{
		response:       make(chan network.Packet, 1),
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
func (f *future) Response() <-chan network.Packet {
	return f.response
}

// SetResponse write packet to the response channel.
func (f *future) SetResponse(response network.Packet) {
	if atomic.CompareAndSwapUint32(&f.finished, 0, 1) {
		f.response <- response
		f.finish()
	}
}

// WaitResponse gets the future response from Response() channel with a timeout set to `duration`.
func (f *future) WaitResponse(duration time.Duration) (network.Packet, error) {
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
		metrics.NetworkFutures.WithLabelValues(f.request.GetType().String()).Dec()
	}
}

func (f *future) finish() {
	close(f.response)
	f.cancelCallback(f)
}
