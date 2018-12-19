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

package hostnetwork

import (
	"context"
	"sync/atomic"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/pkg/errors"
)

type transportBase struct {
	started          uint32
	transport        transport.Transport
	origin           *host.Host
	messageProcessor func(ctx context.Context, msg *packet.Packet)
}

// Listen start listening to network requests, should be started in goroutine.
func (h *transportBase) Start(ctx context.Context) {
	if !atomic.CompareAndSwapUint32(&h.started, 0, 1) {
		inslogger.FromContext(ctx).Warn("double listen initiated")
		return
	}
	go h.listen(ctx)
	go func(ctx context.Context) {
		err := h.transport.Listen(ctx)
		if err != nil {
			inslogger.FromContext(ctx).Error(err)
		}
	}(ctx)
}

func (h *transportBase) listen(ctx context.Context) {
	for {
		select {
		case msg := <-h.transport.Packets():
			if msg == nil {
				log.Error("HostNetwork receiving channel is closed")
				break
			}
			if msg.Error != nil {
				log.Warnf("Received error response: %s", msg.Error.Error())
			}
			go h.messageProcessor(ctx, msg)
		case <-h.transport.Stopped():
			if atomic.CompareAndSwapUint32(&h.started, 1, 0) {
				h.transport.Close()
			}
			return
		}
	}
}

// Disconnect stop listening to network requests.
func (h *transportBase) Stop() {
	go h.transport.Stop()
	if atomic.CompareAndSwapUint32(&h.started, 1, 0) {
		<-h.transport.Stopped()
		h.transport.Close()
	}
}

func (h *transportBase) buildRequest(request network.Request, receiver *host.Host) *packet.Packet {
	return packet.NewBuilder(h.origin).Receiver(receiver).
		Type(request.GetType()).Request(request.GetData()).Build()
}

// PublicAddress returns public address that can be published for all nodes.
func (h *transportBase) PublicAddress() string {
	return h.origin.Address.String()
}

// GetNodeID get current node ID.
func (h *transportBase) GetNodeID() core.RecordRef {
	return h.origin.NodeID
}

// NewRequestBuilder create packet Builder for an outgoing request with sender set to current node.
func (h *transportBase) NewRequestBuilder() network.RequestBuilder {
	return &Builder{sender: h.origin}
}

func getOrigin(tp transport.Transport, id string) (*host.Host, error) {
	address, err := host.NewAddress(tp.PublicAddress())
	if err != nil {
		return nil, errors.Wrap(err, "error resolving address")
	}
	nodeID, err := core.NewRefFromBase58(id)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing NodeID from string")
	}
	origin := &host.Host{NodeID: *nodeID, Address: address}
	return origin, nil
}
