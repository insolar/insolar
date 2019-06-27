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

package consensusadapters

import (
	"context"

	"github.com/insolar/insolar/insolar"
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
func (ps *PacketSender) SendPacketToTransport(t common2.NodeProfile, sendOptions core.PacketSendOptions, payload interface{}) {
	ctx := context.TODO()
	addr := t.GetDefaultEndpoint()

	bs := insolar.MustSerialize(payload)

	err := ps.datagramTransport.SendDatagram(ctx, addr.String(), bs)
	if err != nil {
		panic(err)
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

func (dh *DatagramHandler) HandleDatagram(address string, buf []byte) {
	packet := payloadWrapper{}
	insolar.MustDeserialize(buf, packet)

	packetParser, ok := packet.Payload.(packets.PacketParser)
	if !ok {
		panic("Failed to cast PacketParser")
	}

	if packetParser == nil {
		panic("PacketParser is nil")
	}

	hostIdentity := common.HostIdentity{
		Addr: common.HostAddress(address),
	}
	err := dh.consensusController.ProcessPacket(packetParser, &hostIdentity)
	if err != nil {
		panic(err)
	}
}
