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

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/consensus"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/sequence"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/relay"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

type transportConsensus struct {
	transportBase
	resolver network.RoutingTable
	handlers map[packets.PacketType]network.ConsensusPacketHandler
}

func (tc *transportConsensus) Start(ctx context.Context) error {
	return tc.transportBase.Start(ctx)
}

func (tc *transportConsensus) Stop(ctx context.Context) error {
	return tc.transportBase.Stop(ctx)
}

// RegisterPacketHandler register a handler function to process incoming requests of a specific type.
func (tc *transportConsensus) RegisterPacketHandler(t packets.PacketType, handler network.ConsensusPacketHandler) {
	_, exists := tc.handlers[t]
	if exists {
		log.Warnf("Multiple handlers for packet type %s are not supported! New handler will replace the old one!", t)
	}
	tc.handlers[t] = handler
}

func (tc *transportConsensus) SignAndSendPacket(packet packets.ConsensusPacket,
	receiver core.RecordRef, service core.CryptographyService) error {

	receiverHost, err := tc.resolver.ResolveConsensusRef(receiver)
	if err != nil {
		return errors.Wrapf(err, "Failed to resolve %s request to node %s", packet.GetType(), receiver.String())
	}
	log.Debugf("Send %s request to host %s", packet.GetType(), receiverHost)
	packet.SetRouting(tc.origin.ShortID, receiverHost.ShortID)
	err = packet.Sign(service)
	if err != nil {
		return errors.Wrapf(err, "Failed to sign %s request to node %s", packet.GetType(), receiver.String())
	}
	ctx := context.Background()
	p := tc.buildPacket(packet, receiverHost)
	err = tc.transport.SendPacket(ctx, p)
	if err == nil {
		statsErr := stats.RecordWithTags(ctx, []tag.Mutator{
			tag.Upsert(consensus.TagPhase, packet.GetType().String()),
		}, consensus.PacketsSent.M(1))
		if statsErr != nil {
			log.Warn(" [ transportConsensus ] Failed to record sent packets metric: " + statsErr.Error())
		}
	}
	return err
}

func (tc *transportConsensus) buildPacket(p packets.ConsensusPacket, receiver *host.Host) *packet.Packet {
	return packet.NewBuilder(tc.origin).Receiver(receiver).Request(p).Build()
}

func (tc *transportConsensus) processMessage(msg *packet.Packet) {
	p, ok := msg.Data.(packets.ConsensusPacket)
	if !ok {
		log.Error("Error processing incoming message: failed to convert to ConsensusPacket")
		return
	}
	log.Debugf("Got %s request from host, shortID: %d", p.GetType(), p.GetOrigin())
	if p.GetTarget() != tc.origin.ShortID {
		log.Errorf("Error processing incoming message: target ID %d differs from origin %d", p.GetTarget(), tc.origin.ShortID)
		return
	}
	if p.GetOrigin() == tc.origin.ShortID {
		log.Errorf("Error processing incoming message: sender ID %d equals to origin %d", p.GetTarget(), tc.origin.ShortID)
		return
	}
	sender, err := tc.resolver.ResolveConsensus(p.GetOrigin())
	// TODO: NETD18-79
	// special case for Phase1 because we can get a valid packet from a node we don't know yet (first consensus case)
	if err != nil && p.GetType() != packets.Phase1 {
		log.Errorf("Error processing incoming message: failed to resolve ShortID (%d) -> NodeID", p.GetOrigin())
		return
	}
	if sender == nil {
		sender = &host.Host{}
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
	result := &transportConsensus{handlers: make(map[packets.PacketType]network.ConsensusPacketHandler)}

	result.transport = tp
	result.sequenceGenerator = sequence.NewGeneratorImpl()
	result.resolver = resolver
	result.origin = origin
	result.messageProcessor = result.processMessage
	return result, nil
}
