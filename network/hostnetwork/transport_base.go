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

package hostnetwork

import (
	"context"
	"sync/atomic"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/sequence"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/pkg/errors"
)

type transportBase struct {
	started           uint32
	transport         transport.Transport
	origin            *host.Host
	messageProcessor  func(msg *packet.Packet)
	sequenceGenerator sequence.Generator
}

// Listen start listening to network requests, should be started in goroutine.
func (h *transportBase) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapUint32(&h.started, 0, 1) {
		return errors.New("Failed to start transport: double listen initiated")
	}
	if err := transport.ListenAndWaitUntilReady(ctx, h.transport); err != nil {
		return errors.Wrap(err, "Failed to start transport: listen syscall failed")
	}

	go h.listen(ctx)
	return nil
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
			go h.messageProcessor(msg)
		case <-h.transport.Stopped():
			return
		}
	}
}

// Disconnect stop listening to network requests.
func (h *transportBase) Stop(ctx context.Context) error {
	if atomic.CompareAndSwapUint32(&h.started, 1, 0) {
		go h.transport.Stop()
		<-h.transport.Stopped()
		h.transport.Close()
	}
	return nil
}

func (h *transportBase) buildRequest(ctx context.Context, request network.Request, receiver *host.Host) *packet.Packet {
	return packet.NewBuilder(h.origin).Receiver(receiver).Type(request.GetType()).RequestID(request.GetRequestID()).
		Request(request.GetData()).TraceID(inslogger.TraceID(ctx)).Build()
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
	return &Builder{sender: h.origin, id: network.RequestID(h.sequenceGenerator.Generate())}
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
