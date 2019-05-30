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
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/consensus/packets"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/sequence"
	"github.com/insolar/insolar/network/transport"
)

type networkConsensus struct {
	Resolver network.RoutingTable `inject:""`
	Factory  transport.Factory    `inject:""`

	nodeID            insolar.Reference
	shortID           insolar.ShortNodeID
	transport         transport.DatagramTransport
	started           uint32
	sequenceGenerator sequence.Generator

	muHandlers sync.RWMutex
	handlers   map[packets.PacketType]network.ConsensusPacketHandler

	muOrigin sync.RWMutex
	origin   *host.Host
}

func (nc *networkConsensus) Init(ctx context.Context) error {
	var err error
	nc.transport, err = nc.Factory.CreateDatagramTransport(nc)
	if err != nil {
		return errors.Wrap(err, "Failed to create datagram transport")
	}

	return err
}

func (nc *networkConsensus) Start(ctx context.Context) error {
	if !atomic.CompareAndSwapUint32(&nc.started, 0, 1) {
		inslogger.FromContext(ctx).Warn("NetworkConsensus component already started")
		return nil
	}

	nc.muOrigin.Lock()
	defer nc.muOrigin.Unlock()

	if err := nc.transport.Start(ctx); err != nil {
		return errors.Wrap(err, "Failed to start datagram transport")
	}

	h, err := host.NewHostNS(nc.transport.Address(), nc.nodeID, nc.shortID)
	if err != nil {
		return errors.Wrap(err, "failed to create host")
	}

	nc.origin = h

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
	return nc.getOrigin().Address.String()
}

// RegisterPacketHandler register a handler function to process incoming requests of a specific type.
func (nc *networkConsensus) RegisterPacketHandler(t packets.PacketType, handler network.ConsensusPacketHandler) {
	nc.muHandlers.Lock()
	defer nc.muHandlers.Unlock()

	_, exists := nc.handlers[t]
	if exists {
		log.Warnf("Multiple handlers for packet type %s are not supported! New handler will replace the old one!", t)
	}
	nc.handlers[t] = handler
}

func (nc *networkConsensus) SignAndSendPacket(packet packets.ConsensusPacket,
	receiver insolar.Reference, service insolar.CryptographyService) error {

	if atomic.LoadUint32(&nc.started) == 0 {
		return errors.New("consensus network is not started")
	}

	receiverHost, err := nc.Resolver.ResolveConsensusRef(receiver)
	if err != nil {
		return errors.Wrapf(err, "Failed to resolve %s request to node %s", packet.GetType(), receiver.String())
	}
	log.Debugf("Send %s request to host %s", packet.GetType(), receiverHost)
	packet.SetRouting(nc.getOrigin().ShortID, receiverHost.ShortID)
	err = packet.Sign(service)
	if err != nil {
		return errors.Wrapf(err, "Failed to sign %s request to node %s", packet.GetType(), receiver.String())
	}
	ctx := context.Background()

	buf, err := packet.Serialize()
	if err != nil {
		return errors.Wrap(err, "Failed to serialize packet.")
	}

	err = nc.transport.SendDatagram(ctx, receiverHost.Address.String(), buf)
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

// NewConsensusNetwork constructor creates new ConsensusNetwork
func NewConsensusNetwork(nodeID string, shortID insolar.ShortNodeID) (network.ConsensusNetwork, error) {

	id, err := insolar.NewReferenceFromBase58(nodeID)
	if err != nil {
		return nil, errors.Wrap(err, "invalid nodeID")
	}

	result := &networkConsensus{
		handlers:          make(map[packets.PacketType]network.ConsensusPacketHandler),
		sequenceGenerator: sequence.NewGenerator(),
		nodeID:            *id,
		shortID:           shortID,
	}

	return result, nil
}

// HandleDatagram callback method handles udp datagram from transport
func (nc *networkConsensus) HandleDatagram(address string, buf []byte) {
	logger := inslogger.FromContext(context.Background())
	r := bytes.NewReader(buf)
	p, err := packets.ExtractPacket(r)
	if err != nil {
		logger.Error("[ HandleDatagram ] could not convert network datagram to ConsensusPacket")
		return
	}

	origin := nc.getOrigin()

	log.Debugf("Got %s request from host, shortID: %d", p.GetType(), p.GetOrigin())
	if p.GetTarget() != origin.ShortID {
		logger.Errorf("[ HandleDatagram ] target ID %d differs from origin %d", p.GetTarget(), origin.ShortID)
		return
	}
	if p.GetOrigin() == origin.ShortID {
		logger.Errorf("[ HandleDatagram ] sender ID %d equals to origin %d", p.GetTarget(), origin.ShortID)
		return
	}
	sender, err := nc.Resolver.ResolveConsensus(p.GetOrigin())
	// TODO: NETD18-79
	// special case for Phase1 because we can get a valid packet from a node we don't know yet (first consensus case)
	if err != nil && p.GetType() != packets.Phase1 {
		logger.Errorf("[ HandleDatagram ] failed to resolve ShortID (%d) -> NodeID", p.GetOrigin())
		return
	}
	if sender == nil {
		sender = &host.Host{}
	}

	nc.muHandlers.RLock()
	defer nc.muHandlers.RUnlock()

	handler, exist := nc.handlers[p.GetType()]
	if !exist {
		logger.Errorf("[ HandleDatagram ] No handler set for packet type %d from node %d, %s", p.GetType(), sender.ShortID, sender.NodeID)
		return
	}
	handler(p, sender.NodeID)
}

func (nc *networkConsensus) getOrigin() *host.Host {
	nc.muOrigin.RLock()
	defer nc.muOrigin.RUnlock()

	return nc.origin
}
