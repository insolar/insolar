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
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/longbits"
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
)

const nanosecondsInSecond = int64(time.Second / time.Nanosecond)

func NewPulse(pulseData pulse.Data) insolar.Pulse {
	var prev insolar.PulseNumber
	if !pulseData.IsFirstPulse() {
		prev = insolar.PulseNumber(pulseData.GetPrevPulseNumber())
	} else {
		prev = insolar.PulseNumber(pulseData.PulseNumber)
	}

	entropy := insolar.Entropy{}
	bytes := pulseData.PulseEntropy.AsBytes()
	copy(entropy[:], bytes)
	copy(entropy[pulseData.PulseEntropy.FixedByteSize():], bytes)

	return insolar.Pulse{
		PulseNumber:      insolar.PulseNumber(pulseData.PulseNumber),
		NextPulseNumber:  insolar.PulseNumber(pulseData.GetNextPulseNumber()),
		PrevPulseNumber:  prev,
		PulseTimestamp:   int64(pulseData.Timestamp) * nanosecondsInSecond,
		EpochPulseNumber: int(pulseData.PulseEpoch),
		Entropy:          entropy,
	}
}

func NewPulseData(p insolar.Pulse) pulse.Data {
	data := pulse.NewPulsarData(
		pulse.Number(p.PulseNumber),
		uint16(p.NextPulseNumber-p.PulseNumber),
		uint16(p.PulseNumber-p.PrevPulseNumber),
		longbits.NewBits512FromBytes(p.Entropy[:]).FoldToBits256(),
	)
	data.Timestamp = uint32(p.PulseTimestamp / nanosecondsInSecond)
	return *data
}

type PulsePacketReader struct {
	longbits.FixedReader
	data pulse.Data
}

func (p *PulsePacketReader) GetPulseNumber() pulse.Number {
	return p.data.PulseNumber
}

func (p *PulsePacketReader) GetPulseDataDigest() cryptkit.DigestHolder {
	return nil
}

func (p *PulsePacketReader) OriginalPulsarPacket() {}

func (p *PulsePacketReader) GetPulseData() pulse.Data {
	return p.data
}

func (p *PulsePacketReader) GetPulseDataEvidence() proofs.OriginalPulsarPacket {
	return p
}

func NewPulsePacketReader(pulse insolar.Pulse, data []byte) *PulsePacketReader {
	return &PulsePacketReader{
		FixedReader: longbits.NewFixedReader(data),
		data:        NewPulseData(pulse),
	}
}

type PulsePacketParser struct {
	pulse       insolar.Pulse
	pulsePacket *PulsePacketReader
}

func (p *PulsePacketParser) ParsePacketBody() (transport.PacketParser, error) {
	return nil, nil
}

func NewPulsePacketParser(pulse insolar.Pulse, data []byte) *PulsePacketParser {
	return &PulsePacketParser{
		pulse:       pulse,
		pulsePacket: NewPulsePacketReader(pulse, data),
	}
}

func (p *PulsePacketParser) IsRelayForbidden() bool {
	return true
}

func (p *PulsePacketParser) GetSourceID() insolar.ShortNodeID {
	return insolar.AbsentShortNodeID
}

func (p *PulsePacketParser) GetReceiverID() insolar.ShortNodeID {
	return insolar.AbsentShortNodeID
}

func (p *PulsePacketParser) GetTargetID() insolar.ShortNodeID {
	return insolar.AbsentShortNodeID
}

func (p *PulsePacketParser) GetPacketType() phases.PacketType {
	return phases.PacketPulse
}

func (p *PulsePacketParser) GetPulseNumber() pulse.Number {
	return pulse.Number(p.pulse.PulseNumber)
}

func (p *PulsePacketParser) GetPulsePacket() transport.PulsePacketReader {
	return p.pulsePacket
}

func (p *PulsePacketParser) GetMemberPacket() transport.MemberPacketReader {
	return nil
}

func (p *PulsePacketParser) GetPacketSignature() cryptkit.SignedDigest {
	return cryptkit.SignedDigest{}
}
