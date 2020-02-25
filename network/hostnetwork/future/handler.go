// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package future

import (
	"context"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/hostnetwork/packet"
)

type packetHandler struct {
	futureManager Manager
}

func NewPacketHandler(futureManager Manager) PacketHandler {
	return &packetHandler{
		futureManager: futureManager,
	}
}

func (ph *packetHandler) Handle(ctx context.Context, response *packet.ReceivedPacket) {
	metrics.NetworkPacketReceivedTotal.WithLabelValues(response.GetType().String()).Inc()
	if !response.IsResponse() {
		return
	}

	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"type":       response.Type,
		"request_id": response.RequestID,
	})
	logger.Debugf("[ processResponse ] Processing %s response from host %s; RequestID = %d",
		response.GetType(), response.Sender, response.RequestID)

	future := ph.futureManager.Get(response.Packet)
	if future != nil {
		if shouldProcessPacket(future, response) {
			start := time.Now()
			future.SetResponse(response)
			logger.Debugf("[ processResponse ] Finished processing future RequestID = %d in %s", future.ID(), time.Since(start).String())
		} else {
			logger.Debugf("[ processResponse ] Canceling future RequestID = %d", future.ID())
			future.Cancel()
		}
	}
}

func shouldProcessPacket(future Future, p *packet.ReceivedPacket) bool {
	typesShouldBeEqual := p.GetType() == future.Request().GetType()
	responseIsForRightSender := future.Receiver().Equal(*p.Sender)

	return typesShouldBeEqual && responseIsForRightSender
}
