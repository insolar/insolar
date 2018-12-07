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
	"fmt"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
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
	p := tc.buildRequest(request, receiverHost)
	return tc.transport.SendPacket(p)
}

func (tc *transportConsensus) processMessage(ctx context.Context, msg *packet.Packet) {
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
	result.resolver = resolver
	result.origin = origin
	result.messageProcessor = result.processMessage
	return result, nil
}
