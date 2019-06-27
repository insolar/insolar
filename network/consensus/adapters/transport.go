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

package adapters

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/consensus/gcpv2/nodeset"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
	"github.com/insolar/insolar/network/transport"
)

type PacketBuilder struct{}

func NewPacketBuilder() *PacketBuilder {
	return &PacketBuilder{}
}

func (pb *PacketBuilder) GetNeighbourhoodSize(populationCount int) common2.NeighbourhoodSizes {
	panic("implement me")
}

func (pb *PacketBuilder) PreparePhase0Packet(sender common2.NodeProfile, pulsarPacket common2.OriginalPulsarPacket, options core.PacketSendOptions) core.PreparedPacketSender {
	panic("implement me")
}

func (pb *PacketBuilder) PreparePhase1Packet(sender common2.NodeProfile, pulsarPacket common2.OriginalPulsarPacket, nsh common2.NodeStateHashEvidence,
	options core.PacketSendOptions) core.PreparedPacketSender {
	panic("implement me")
}

func (pb *PacketBuilder) PreparePhase2Packet(sender common2.NodeProfile, pd common.PulseData, neighbourhood []packets.NodeStateHashReportReader,
	intros []common2.NodeIntroduction, options core.PacketSendOptions) core.PreparedPacketSender {
	panic("implement me")
}

func (pb *PacketBuilder) PreparePhase3Packet(sender common2.NodeProfile, pd common.PulseData, bitset nodeset.NodeBitset,
	gshTrusted common2.GlobulaStateHash, gshDoubted common2.GlobulaStateHash,
	options core.PacketSendOptions) core.PreparedPacketSender {
	panic("implement me")
}

type PacketSender struct {
	datagramTransport transport.DatagramTransport
}

func NewPacketSender(datagramTransport transport.DatagramTransport) *PacketSender {
	return &PacketSender{
		datagramTransport: datagramTransport,
	}
}

type payloadWrapper struct {
	Payload interface{}
}

// TODO: signature seems to be wrong :( context missed
func (ps *PacketSender) SendPacketToTransport(to common2.NodeProfile, sendOptions core.PacketSendOptions, payload interface{}) {
	ctx := context.TODO()
	addr := to.GetDefaultEndpoint().String()

	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"receiver_addr":    addr,
		"receiver_node_id": to.GetShortNodeID(),
		"options":          sendOptions,
	})

	bs, err := insolar.Serialize(payload)
	if err != nil {
		logger.Error("Failed to serialize payload")
	}

	err = ps.datagramTransport.SendDatagram(ctx, addr, bs)
	if err != nil {
		logger.Error("Failed to send datagram")
	}
}

type DatagramHandler struct {
	consensusController core.ConsensusController
}

func NewDatagramHandler() *DatagramHandler {
	return &DatagramHandler{}
}

func (dh *DatagramHandler) SetConsensusController(consensusController core.ConsensusController) {
	dh.consensusController = consensusController
}

func (dh *DatagramHandler) HandleDatagram(ctx context.Context, address string, buf []byte) {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"address": address,
	})

	p := payloadWrapper{}
	err := insolar.Deserialize(buf, p)
	if err != nil {
		logger.Error(err)
		return
	}

	packetParser, ok := p.Payload.(packets.PacketParser)
	if !ok {
		logger.Error("Failed to get PacketParser")
		return
	}

	if packetParser == nil {
		logger.Error("PacketParser is nil")
		return
	}

	hostIdentity := common.HostIdentity{
		Addr: common.HostAddress(address),
	}
	err = dh.consensusController.ProcessPacket( /*TODO: ctx, */ packetParser, &hostIdentity)
	if err != nil {
		logger.Error("Failed to process p")
		return
	}
}
