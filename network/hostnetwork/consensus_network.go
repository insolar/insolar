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

package hostnetwork

import (
	"bytes"
	"context"
	"sync/atomic"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/consensus"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/sequence"
	"github.com/insolar/insolar/network/transport"
)

type networkConsensus struct {
	Resolver network.RoutingTable `inject:""`

	transport         transport.Transport2
	origin            *host.Host
	started           uint32
	sequenceGenerator sequence.Generator
	messageProcessor  func(p packets.ConsensusPacket)
	handlers          map[packets.PacketType]network.ConsensusPacketHandler
}

func (nc *networkConsensus) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapUint32(&nc.started, 0, 1) {
		return errors.New("Failed to start transport: double listen initiated")
	}
	if err := nc.transport.Start(ctx); err != nil {
		return errors.Wrap(err, "Failed to start transport: listen syscall failed")
	}

	return nil
}

func (nc *networkConsensus) Stop(ctx context.Context) error {
	if atomic.CompareAndSwapUint32(&nc.started, 1, 0) {
		err := nc.transport.Stop(ctx)
		if err != nil {
			return errors.Wrap(err, "Failed to stop transport.")
		}
	}
	return nil
}

// PublicAddress returns public address that can be published for all nodes.
func (nc *networkConsensus) PublicAddress() string {
	return nc.origin.Address.String()
}

// GetNodeID get current node ID.
func (nc *networkConsensus) GetNodeID() insolar.Reference {
	return nc.origin.NodeID
}

// RegisterPacketHandler register a handler function to process incoming requests of a specific type.
func (nc *networkConsensus) RegisterPacketHandler(t packets.PacketType, handler network.ConsensusPacketHandler) {
	_, exists := nc.handlers[t]
	if exists {
		log.Warnf("Multiple handlers for packet type %s are not supported! New handler will replace the old one!", t)
	}
	nc.handlers[t] = handler
}

func (nc *networkConsensus) SignAndSendPacket(packet packets.ConsensusPacket,
	receiver insolar.Reference, service insolar.CryptographyService) error {

	receiverHost, err := nc.Resolver.ResolveConsensusRef(receiver)
	if err != nil {
		return errors.Wrapf(err, "Failed to resolve %s request to node %s", packet.GetType(), receiver.String())
	}
	log.Debugf("Send %s request to host %s", packet.GetType(), receiverHost)
	packet.SetRouting(nc.origin.ShortID, receiverHost.ShortID)
	err = packet.Sign(service)
	if err != nil {
		return errors.Wrapf(err, "Failed to sign %s request to node %s", packet.GetType(), receiver.String())
	}
	ctx := context.Background()

	buf, err := packet.Serialize()
	if err != nil {
		return errors.Wrap(err, "Failed to serialize packet.")
	}

	err = nc.transport.SendDgram(ctx, receiverHost.Address.String(), buf)
	if err == nil {
		statsErr := stats.RecordWithTags(ctx, []tag.Mutator{
			tag.Upsert(consensus.TagPhase, packet.GetType().String()),
		}, consensus.PacketsSent.M(1))
		if statsErr != nil {
			log.Warn(" [ networkConsensus ] Failed to record sent packets metric: " + statsErr.Error())
		}
	}
	return err
}

func (nc *networkConsensus) processMessage(p packets.ConsensusPacket) {
	// p, ok := msg.Data.(packets.ConsensusPacket)
	// if !ok {
	// 	log.Error("Error processing incoming message: failed to convert to ConsensusPacket")
	// 	return
	// }
	log.Debugf("Got %s request from host, shortID: %d", p.GetType(), p.GetOrigin())
	if p.GetTarget() != nc.origin.ShortID {
		log.Errorf("Error processing incoming message: target ID %d differs from origin %d", p.GetTarget(), nc.origin.ShortID)
		return
	}
	if p.GetOrigin() == nc.origin.ShortID {
		log.Errorf("Error processing incoming message: sender ID %d equals to origin %d", p.GetTarget(), nc.origin.ShortID)
		return
	}
	sender, err := nc.Resolver.ResolveConsensus(p.GetOrigin())
	// TODO: NETD18-79
	// special case for Phase1 because we can get a valid packet from a node we don't know yet (first consensus case)
	if err != nil && p.GetType() != packets.Phase1 {
		log.Errorf("Error processing incoming message: failed to resolve ShortID (%d) -> NodeID", p.GetOrigin())
		return
	}
	if sender == nil {
		sender = &host.Host{}
	}
	handler, exist := nc.handlers[p.GetType()]
	if !exist {
		log.Errorf("No handler set for packet type %s from node %d, %s", p.GetType(), sender.ShortID, sender.NodeID)
		return
	}
	handler(p, sender.NodeID)
}

func NewConsensusNetwork(address, nodeID string, shortID insolar.ShortNodeID) (network.ConsensusNetwork, error) {
	conf := configuration.Transport{}
	conf.Address = address
	conf.Protocol = "PURE_UDP"

	tp, publicAddress, err := transport.NewTransport2(conf)
	if err != nil {
		return nil, errors.Wrap(err, "error creating transport")
	}
	id, err := insolar.NewReferenceFromBase58(nodeID)
	if err != nil {
		return nil, errors.Wrap(err, "invalid nodeID")
	}

	origin, err := host.NewHostNS(publicAddress, *id, shortID)
	if err != nil {
		return nil, errors.Wrap(err, "error getting origin")
	}
	result := &networkConsensus{handlers: make(map[packets.PacketType]network.ConsensusPacketHandler)}

	result.transport = tp
	tp.SetDgramProcessor(result)

	result.sequenceGenerator = sequence.NewGeneratorImpl()
	result.origin = origin
	result.messageProcessor = result.processMessage
	return result, nil
}

func (nc *networkConsensus) ProcessDgram(address string, buf []byte) error {
	r := bytes.NewReader(buf)
	p, err := packets.ExtractPacket(r)
	if err != nil {
		return errors.Wrap(err, "could not convert network datagram to ConsensusPacket")
	}

	nc.processMessage(p)
	return nil
}
