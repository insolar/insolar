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
	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/relay"
	"github.com/pkg/errors"
)

type transportConsensus struct {
	transportBase
	resolver network.RoutingTable
	handlers map[consensus.PacketType]network.ConsensusPacketHandler
}

// RegisterPacketHandler register a handler function to process incoming requests of a specific type.
func (tc *transportConsensus) RegisterPacketHandler(t consensus.PacketType, handler network.ConsensusPacketHandler) {
	_, exists := tc.handlers[t]
	if exists {
		panic(fmt.Sprintf("multiple handlers for packet type %s are not supported!", t.String()))
	}
	tc.handlers[t] = handler
}

func (tc *transportConsensus) SignAndSendPacket(packet consensus.ConsensusPacket,
	receiver core.RecordRef, service core.CryptographyService) error {

	log.Debugf("Send %s request to host %s", packet.GetType(), receiver.String())
	receiverHost, err := tc.resolver.ResolveConsensusRef(receiver)
	if err != nil {
		return errors.Wrapf(err, "Failed to resolve %s request to node %s", packet.GetType(), receiver.String())
	}
	packet.SetRouting(tc.origin.ShortID, receiverHost.ShortID)
	err = packet.Sign(service)
	if err != nil {
		return errors.Wrapf(err, "Failed to sign %s request to node %s", packet.GetType(), receiver.String())
	}
	p := tc.buildPacket(packet, receiverHost)
	return tc.transport.SendPacket(p)
}

func (tc *transportConsensus) buildPacket(p consensus.ConsensusPacket, receiver *host.Host) *packet.Packet {
	return packet.NewBuilder(tc.origin).Receiver(receiver).Request(p).Build()
}

func (tc *transportConsensus) processMessage(ctx context.Context, msg *packet.Packet) {
	p, ok := msg.Data.(consensus.ConsensusPacket)
	if !ok {
		log.Error("Error processing incoming message: failed to convert to ConsensusPacket")
		return
	}
	log.Debugf("Got %s request from host, shortID: %d", p.GetType(), p.GetOrigin())
	if p.GetTarget() != tc.origin.ShortID {
		log.Errorf("Error processing incoming message: target ID %d differs from origin %d", p.GetTarget(), tc.origin.ShortID)
		return
	}
	sender, err := tc.resolver.ResolveConsensus(p.GetOrigin())
	// TODO: NETD18-79
	// special case for Phase1 because we can get a valid packet from a node we don't know yet (first consensus case)
	if err != nil && p.GetType() != consensus.Phase1 {
		log.Errorf("Error processing incoming message: failed to resolve ShortID (%d) -> NodeID", msg.Sender.ShortID)
		return
	}
	handler, exist := tc.handlers[p.GetType()]
	if !exist {
		log.Errorf("No handler set for packet type %s from node %d, %s", p.GetType(), sender.ShortID, sender.NodeID)
		return
	}
	handler(p, sender.NodeID)
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
	result := &transportConsensus{handlers: make(map[consensus.PacketType]network.ConsensusPacketHandler)}

	result.transport = tp
	result.resolver = resolver
	result.origin = origin
	result.messageProcessor = result.processMessage
	return result, nil
}
