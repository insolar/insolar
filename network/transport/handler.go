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
	"context"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/packet/types"
)

type packetHandlerImpl struct {
	futureManager futureManager

	received chan *packet.Packet
}

func newPacketHandlerImpl(futureManager futureManager) *packetHandlerImpl {
	return &packetHandlerImpl{
		futureManager: futureManager,
		received:      make(chan *packet.Packet),
	}
}

func (ph *packetHandlerImpl) Handle(ctx context.Context, msg *packet.Packet) {
	metrics.NetworkPacketReceivedTotal.WithLabelValues(msg.Type.String()).Inc()
	if msg.IsResponse {
		ph.processResponse(ctx, msg)
		return
	}

	ph.processRequest(ctx, msg)
}

func (ph *packetHandlerImpl) Received() <-chan *packet.Packet {
	return ph.received
}

func (ph *packetHandlerImpl) processResponse(ctx context.Context, msg *packet.Packet) {
	logger := inslogger.FromContext(ctx)

	logger.Debugf("[ processResponse ] Process response %s from %s with RequestID = %d", msg.Type, msg.RemoteAddress, msg.RequestID)

	future := ph.futureManager.Get(msg)
	if future != nil {
		if shouldProcessPacket(future, msg) {
			logger.Debugf("[ processResponse ] Processing future with RequestID = %v", msg.RequestID)
			future.SetResult(msg)
		} else {
			logger.Debugf("[ processResponse ] Canceling future with RequestID = %v", msg.RequestID)
		}
		future.Cancel()
	}
}

func (ph *packetHandlerImpl) processRequest(ctx context.Context, msg *packet.Packet) {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("[ processRequest ] Process request %s from %s with RequestID = %d", msg.Type, msg.RemoteAddress, msg.RequestID)

	ph.received <- msg
}

func shouldProcessPacket(future Future, msg *packet.Packet) bool {
	typesShouldBeEqual := msg.Type == future.Request().Type
	responseIsForRightSender := future.Actor().Equal(*msg.Sender)

	return typesShouldBeEqual && (responseIsForRightSender || msg.Type == types.Ping)
}
