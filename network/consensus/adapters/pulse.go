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
	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/common/long_bits"
	"github.com/insolar/insolar/network/consensus/common/pulse_data"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"io"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

const nanosecondsInSecond = int64(time.Second / time.Nanosecond)

func NewPulse(pulseData pulse_data.PulseData) insolar.Pulse {
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

func NewPulseData(pulse insolar.Pulse) pulse_data.PulseData {
	data := pulse_data.NewPulsarData(
		pulse_data.PulseNumber(pulse.PulseNumber),
		uint16(pulse.NextPulseNumber-pulse.PulseNumber),
		uint16(pulse.PulseNumber-pulse.PrevPulseNumber),
		long_bits.NewBits512FromBytes(pulse.Entropy[:]).FoldToBits256(),
	)
	data.Timestamp = uint32(pulse.PulseTimestamp / nanosecondsInSecond)
	return *data
}

type PulsePacketReader struct {
	pulse insolar.Pulse
}

func (p *PulsePacketReader) OriginalPulsarPacket() {}

func (p *PulsePacketReader) GetPulseData() pulse_data.PulseData {
	return NewPulseData(p.pulse)
}

func (p *PulsePacketReader) GetPulseDataEvidence() packets.OriginalPulsarPacket {
	return p
}

func (p *PulsePacketReader) WriteTo(w io.Writer) (n int64, err error) {
	panic("implement me")
}

func (p *PulsePacketReader) Read(b []byte) (n int, err error) {
	panic("implement me")
}

func (p *PulsePacketReader) AsBytes() []byte {
	panic("implement me")
}

func (p *PulsePacketReader) AsByteString() string {
	panic("implement me")
}

func (p *PulsePacketReader) FixedByteSize() int {
	panic("implement me")
}

func NewPulsePacketReader(pulse insolar.Pulse) *PulsePacketReader {
	return &PulsePacketReader{pulse}
}

type PulsePacketParser struct {
	pulse insolar.Pulse
}

func (p *PulsePacketParser) IsRelayForbidden() bool {
	return false
}

func NewPulsePacketParser(pulse insolar.Pulse) *PulsePacketParser {
	return &PulsePacketParser{pulse}
}

func (p *PulsePacketParser) GetSourceID() common.ShortNodeID {
	panic("implement me")
}

func (p *PulsePacketParser) GetReceiverID() common.ShortNodeID {
	panic("implement me")
}

func (p *PulsePacketParser) GetTargetID() common.ShortNodeID {
	panic("implement me")
}

func (p *PulsePacketParser) GetPacketType() api.PacketType {
	return api.PacketPulse
}

func (p *PulsePacketParser) GetPulseNumber() pulse_data.PulseNumber {
	return pulse_data.PulseNumber(p.pulse.PulseNumber)
}

func (p *PulsePacketParser) GetPulsePacket() packets.PulsePacketReader {
	return NewPulsePacketReader(p.pulse)
}

func (p *PulsePacketParser) GetMemberPacket() packets.MemberPacketReader {
	return nil
}

func (p *PulsePacketParser) GetPacketSignature() cryptography_containers.SignedDigest {
	return cryptography_containers.SignedDigest{}
}
