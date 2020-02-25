// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package adapters

import (
	"fmt"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/serialization/pulseserialization"
	"github.com/insolar/insolar/pulse"
)

const nanosecondsInSecond = int64(time.Second / time.Nanosecond)

func NewPulse(pulseData pulse.Data) insolar.Pulse {
	var prev insolar.PulseNumber
	if !pulseData.IsFirstPulse() {
		prev = pulseData.PrevPulseNumber()
	} else {
		prev = pulseData.PulseNumber
	}

	entropy := insolar.Entropy{}
	bs := pulseData.PulseEntropy.AsBytes()
	copy(entropy[:], bs)
	copy(entropy[pulseData.PulseEntropy.FixedByteSize():], bs)

	return insolar.Pulse{
		PulseNumber:      pulseData.PulseNumber,
		NextPulseNumber:  pulseData.NextPulseNumber(),
		PrevPulseNumber:  prev,
		PulseTimestamp:   int64(pulseData.Timestamp) * nanosecondsInSecond,
		EpochPulseNumber: pulseData.PulseEpoch,
		Entropy:          entropy,
	}
}

func NewPulseData(p insolar.Pulse) pulse.Data {
	data := pulse.NewPulsarData(
		p.PulseNumber,
		uint16(p.NextPulseNumber-p.PulseNumber),
		uint16(p.PulseNumber-p.PrevPulseNumber),
		longbits.NewBits512FromBytes(p.Entropy[:]).FoldToBits256(),
	)
	data.Timestamp = uint32(p.PulseTimestamp / nanosecondsInSecond)
	data.PulseEpoch = p.EpochPulseNumber
	return data
}

func NewPulseDigest(data pulse.Data) cryptkit.Digest {
	entropySize := data.PulseEntropy.FixedByteSize()

	bits := longbits.Bits512{}
	copy(bits[:entropySize], data.PulseEntropy[:])
	copy(bits[entropySize:], data.PulseEntropy[:])

	// It's not digest actually :)
	return cryptkit.NewDigest(&bits, SHA3512Digest)
}

type PulsePacketParser struct {
	longbits.FixedReader
	digest cryptkit.DigestHolder
	pulse  pulse.Data
}

func NewPulsePacketParser(pulse pulse.Data) *PulsePacketParser {
	data, err := pulseserialization.Serialize(pulse)
	if err != nil {
		panic(err.Error())
	}

	return &PulsePacketParser{
		FixedReader: longbits.NewFixedReader(data),
		digest:      NewPulseDigest(pulse).AsDigestHolder(),
		pulse:       pulse,
	}
}

func (p PulsePacketParser) String() string {
	return fmt.Sprintf("<pt=pulse body=<%s>>", p.pulse.String())
}

func (p *PulsePacketParser) ParsePacketBody() (transport.PacketParser, error) {
	return nil, nil
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
	return phases.PacketPulsarPulse
}

func (p *PulsePacketParser) GetPulseNumber() pulse.Number {
	return p.pulse.PulseNumber
}

func (p *PulsePacketParser) GetPulsePacket() transport.PulsePacketReader {
	return p
}

func (p *PulsePacketParser) GetMemberPacket() transport.MemberPacketReader {
	return nil
}

func (p *PulsePacketParser) GetPacketSignature() cryptkit.SignedDigest {
	return cryptkit.SignedDigest{}
}

func (p *PulsePacketParser) GetPulseDataDigest() cryptkit.DigestHolder {
	return p.digest
}

func (p *PulsePacketParser) OriginalPulsarPacket() {}

func (p *PulsePacketParser) GetPulseData() pulse.Data {
	return p.pulse
}

func (p *PulsePacketParser) GetPulseDataEvidence() proofs.OriginalPulsarPacket {
	return p
}
