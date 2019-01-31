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
	"fmt"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/sequence"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/insolar/insolar/network/transport/relay"
	"github.com/pkg/errors"
)

type transportConsensus struct {
	transportBase
	resolver network.RoutingTable
	handlers map[types.PacketType]network.ConsensusRequestHandler
}

// RegisterRequestHandler register a handler function to process incoming requests of a specific type.
func (tc *transportConsensus) RegisterRequestHandler(t types.PacketType, handler network.ConsensusRequestHandler) {
	_, exists := tc.handlers[t]
	if exists {
		panic(fmt.Sprintf("multiple handlers for packet type %s are not supported!", t.String()))
	}
	tc.handlers[t] = handler
}

func (tc *transportConsensus) SendRequest(request network.Request, receiver core.RecordRef) error {
	log.Debugf("Send %s request to host %s", request.GetType().String(), receiver.String())
	receiverHost, err := tc.resolver.Resolve(receiver)
	if err != nil {
		return errors.Wrapf(err, "Failed to send %s request to node %s",
			request.GetType().String(), receiver.String())
	}
	ctx := context.Background()
	p := tc.buildRequest(ctx, request, receiverHost)
	return tc.transport.SendPacket(ctx, p)
}

func (tc *transportConsensus) processMessage(msg *packet.Packet) {
	log.Debugf("Got %s request from host, shortID: %d", msg.Type.String(), msg.Sender.ShortID)
	sender, err := tc.resolver.ResolveS(msg.Sender.ShortID)
	if err != nil {
		log.Errorf("Error processing incoming message: failed to resolve ShortID (%d) -> NodeID", msg.Sender.ShortID)
		return
	}
	msg.Sender = sender
	handler, exist := tc.handlers[msg.Type]
	if !exist {
		log.Errorf("No handler set for packet type %s from node %s",
			msg.Type.String(), msg.Sender.NodeID.String())
		return
	}
	handler((*packetWrapper)(msg))
}

func NewConsensusNetwork(address, nodeID string, shortID core.ShortNodeID,
	resolver network.RoutingTable) (network.ConsensusNetwork, error) {

	conf := configuration.Transport{}
	conf.Address = address
	conf.Protocol = "PURE_UDP"
	conf.BehindNAT = false

	tp, err := transport.NewTransport(conf, relay.NewProxy())
	if err != nil {
		return nil, errors.Wrap(err, "error creating transport")
	}
	origin, err := getOrigin(tp, nodeID)
	if err != nil {
		go tp.Stop()
		<-tp.Stopped()
		tp.Close()
		return nil, errors.Wrap(err, "error getting origin")
	}
	origin.ShortID = shortID
	result := &transportConsensus{handlers: make(map[types.PacketType]network.ConsensusRequestHandler)}

	result.transport = tp
	result.sequenceGenerator = sequence.NewGeneratorImpl()
	result.resolver = resolver
	result.origin = origin
	result.messageProcessor = result.processMessage
	return result, nil
}
