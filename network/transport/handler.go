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
	"context"

	"github.com/insolar/insolar/instrumentation/inslogger"
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

	logger.Debugf("[ processResponse ] Process response %s with RequestID = %d", msg.RemoteAddress, msg.RequestID)

	future := ph.futureManager.Get(msg)
	if future != nil {
		if shouldProcessPacket(future, msg) {
			logger.Debugf("[ processResponse ] Processing future with RequestID = %s", msg.RequestID)
			future.SetResult(msg)
		} else {
			logger.Debugf("[ processResponse ] Canceling future with RequestID = %s", msg.RequestID)
		}
		future.Cancel()
	}
}

func (ph *packetHandlerImpl) processRequest(ctx context.Context, msg *packet.Packet) {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("[ processRequest ] Process request %s with RequestID = %d", msg.RemoteAddress, msg.RequestID)

	ph.received <- msg
}

func shouldProcessPacket(future Future, msg *packet.Packet) bool {
	typesShouldBeEqual := msg.Type == future.Request().Type
	responseIsForRightSender := future.Actor().Equal(*msg.Sender)

	return typesShouldBeEqual && (responseIsForRightSender || msg.Type == types.Ping)
}
