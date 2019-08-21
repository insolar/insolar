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

package serialization

import (
	"io"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/pulse"
)

const (
	primaryRoleBitSize = 6
	primaryRoleMask    = 1<<primaryRoleBitSize - 1 // 0b00111111
	primaryRoleMax     = primaryRoleMask

	addrModeBitSize = 2
	addrModeShift   = primaryRoleBitSize
	addrModeMax     = 1<<addrModeBitSize - 1
)

type NodeBriefIntro struct {
	// ByteSize= 135, 137, 147
	// ByteSize= 135 + (0, 2, 12)

	/*
		This field MUST be excluded from the packet, but considered for signature calculation.
		Value of this field equals SourceID or AnnounceID.
	*/
	ShortID insolar.ShortNodeID `insolar-transport:"ignore=send"` // ByteSize = 0

	PrimaryRoleAndFlags uint8 `insolar-transport:"[0:5]=header:PrimaryRole;[6:7]=header:AddrMode"` // AddrMode =0 reserved, =1 Relay, =2 IPv4 =3 IPv6
	SpecialRoles        member.SpecialRole
	StartPower          member.Power

	// 4 | 6 | 18 bytes
	// InboundRelayID common.ShortNodeID `insolar-transport:"AddrMode=2"`
	// BasePort    uint16 `insolar-transport:"AddrMode=0,1"`
	// PrimaryIPv4 uint32 `insolar-transport:"AddrMode=0"`
	// PrimaryIPv6    [4]uint32          `insolar-transport:"AddrMode=1"`
	Endpoint [18]byte

	// 128 bytes
	NodePK longbits.Bits512 // works as a unique node identity

	JoinerData      []byte           `insolar-transport:"ignore=send"` // ByteSize = 0
	JoinerSignature longbits.Bits512 // ByteSize=64
}

func (bi *NodeBriefIntro) GetPrimaryRole() member.PrimaryRole {
	return member.PrimaryRole(bi.PrimaryRoleAndFlags & primaryRoleMask)
}

func (bi *NodeBriefIntro) SetPrimaryRole(primaryRole member.PrimaryRole) {
	if primaryRole > primaryRoleMax {
		panic("invalid primary role")
	}

	bi.PrimaryRoleAndFlags |= uint8(primaryRole)
}
func (bi *NodeBriefIntro) GetAddrMode() endpoints.NodeEndpointType {
	return endpoints.NodeEndpointType(bi.PrimaryRoleAndFlags >> addrModeShift)
}

func (bi *NodeBriefIntro) SetAddrMode(addrMode endpoints.NodeEndpointType) {
	if addrMode > addrModeMax {
		panic("invalid addr mode")
	}

	bi.PrimaryRoleAndFlags |= uint8(addrMode) << addrModeShift
}

func (bi *NodeBriefIntro) SerializeTo(ctx SerializeContext, writer io.Writer) error {
	if err := write(writer, bi.PrimaryRoleAndFlags); err != nil {
		return errors.Wrap(err, "failed to serialize PrimaryRoleAndFlags")
	}

	if err := write(writer, bi.SpecialRoles); err != nil {
		return errors.Wrap(err, "failed to serialize SpecialRoles")
	}

	if err := write(writer, bi.StartPower); err != nil {
		return errors.Wrap(err, "failed to serialize StartPower")
	}

	if err := write(writer, bi.Endpoint); err != nil {
		return errors.Wrap(err, "failed to serialize Endpoint")
	}

	if err := write(writer, bi.NodePK); err != nil {
		return errors.Wrap(err, "failed to serialize NodePK")
	}

	if err := write(writer, bi.JoinerSignature); err != nil {
		return errors.Wrap(err, "failed to serialize JoinerSignature")
	}

	return nil
}

func (bi *NodeBriefIntro) DeserializeFrom(ctx DeserializeContext, reader io.Reader) error {
	capture := network.NewCapturingReader(reader)

	if err := read(capture, &bi.PrimaryRoleAndFlags); err != nil {
		return errors.Wrap(err, "failed to deserialize PrimaryRoleAndFlags")
	}

	if err := read(capture, &bi.SpecialRoles); err != nil {
		return errors.Wrap(err, "failed to deserialize SpecialRoles")
	}

	if err := read(capture, &bi.StartPower); err != nil {
		return errors.Wrap(err, "failed to deserialize StartPower")
	}

	if err := read(capture, &bi.Endpoint); err != nil {
		return errors.Wrap(err, "failed to deserialize BasePort")
	}

	if err := read(capture, &bi.NodePK); err != nil {
		return errors.Wrap(err, "failed to deserialize NodePK")
	}

	bi.JoinerData = capture.Captured()

	if err := read(reader, &bi.JoinerSignature); err != nil {
		return errors.Wrap(err, "failed to deserialize JoinerSignature")
	}

	return nil
}

type NodeExtendedIntro struct {
	// ByteSize>=86
	IssuedAtPulse pulse.Number // =0 when a node was connected during zeronet
	IssuedAtTime  uint64

	PowerLevels member.PowerSet // ByteSize=4

	EndpointLen    uint8
	ExtraEndpoints []uint16

	ProofLen     uint8
	NodeRefProof []longbits.Bits512

	DiscoveryIssuerNodeID insolar.ShortNodeID
	IssuerSignature       longbits.Bits512
}

func (ei *NodeExtendedIntro) SerializeTo(ctx SerializeContext, writer io.Writer) error {
	if err := write(writer, ei.IssuedAtPulse); err != nil {
		return errors.Wrap(err, "failed to serialize IssuedAtPulse")
	}

	if err := write(writer, ei.IssuedAtTime); err != nil {
		return errors.Wrap(err, "failed to serialize IssuedAtTime")
	}

	if err := write(writer, ei.PowerLevels); err != nil {
		return errors.Wrap(err, "failed to serialize PowerLevels")
	}

	if err := write(writer, ei.EndpointLen); err != nil {
		return errors.Wrap(err, "failed to serialize EndpointLen")
	}

	for i := 0; i < int(ei.EndpointLen); i++ {
		if err := write(writer, ei.ExtraEndpoints[i]); err != nil {
			return errors.Wrapf(err, "failed to serialize ExtraEndpoints[%d]", i)
		}
	}

	if err := write(writer, ei.ProofLen); err != nil {
		return errors.Wrap(err, "failed to serialize ProofLen")
	}

	for i := 0; i < int(ei.ProofLen); i++ {
		if err := write(writer, ei.NodeRefProof[i]); err != nil {
			return errors.Wrapf(err, "failed to serialize NodeRefProof[%d]", i)
		}
	}

	if err := write(writer, ei.DiscoveryIssuerNodeID); err != nil {
		return errors.Wrap(err, "failed to serialize DiscoveryIssuerNodeID")
	}

	if err := write(writer, ei.IssuerSignature); err != nil {
		return errors.Wrap(err, "failed to serialize IssuerSignature")
	}

	return nil
}

func (ei *NodeExtendedIntro) DeserializeFrom(ctx DeserializeContext, reader io.Reader) error {
	if err := read(reader, &ei.IssuedAtPulse); err != nil {
		return errors.Wrap(err, "failed to deserialize IssuedAtPulse")
	}

	if err := read(reader, &ei.IssuedAtTime); err != nil {
		return errors.Wrap(err, "failed to deserialize IssuedAtTime")
	}

	if err := read(reader, &ei.PowerLevels); err != nil {
		return errors.Wrap(err, "failed to deserialize PowerLevels")
	}

	if err := read(reader, &ei.EndpointLen); err != nil {
		return errors.Wrap(err, "failed to deserialize EndpointLen")
	}

	if ei.EndpointLen > 0 {
		ei.ExtraEndpoints = make([]uint16, ei.EndpointLen)
		for i := 0; i < int(ei.EndpointLen); i++ {
			if err := read(reader, &ei.ExtraEndpoints[i]); err != nil {
				return errors.Wrapf(err, "failed to deserialize ExtraEndpoints[%d]", i)
			}
		}
	}

	if err := read(reader, &ei.ProofLen); err != nil {
		return errors.Wrap(err, "failed to deserialize ProofLen")
	}

	if ei.ProofLen > 0 {
		ei.NodeRefProof = make([]longbits.Bits512, ei.ProofLen)
		for i := 0; i < int(ei.ProofLen); i++ {
			if err := read(reader, &ei.NodeRefProof[i]); err != nil {
				return errors.Wrapf(err, "failed to deserialize NodeRefProof[%d]", i)
			}
		}
	}

	if err := read(reader, &ei.DiscoveryIssuerNodeID); err != nil {
		return errors.Wrap(err, "failed to deserialize DiscoveryIssuerNodeID")
	}

	if err := read(reader, &ei.IssuerSignature); err != nil {
		return errors.Wrap(err, "failed to deserialize IssuerSignature")
	}

	return nil
}

type NodeFullIntro struct {
	// ByteSize= >=86 + (135, 137, 147) = >(221, 223, 233)

	NodeBriefIntro    // ByteSize= 135, 137, 147
	NodeExtendedIntro // ByteSize>=86
}

func (fi *NodeFullIntro) SerializeTo(ctx SerializeContext, writer io.Writer) error {
	if err := fi.NodeBriefIntro.SerializeTo(ctx, writer); err != nil {
		return errors.Wrap(err, "failed to serialize NodeBriefIntro")
	}

	if err := fi.NodeExtendedIntro.SerializeTo(ctx, writer); err != nil {
		return errors.Wrap(err, "failed to serialize NodeExtendedIntro")
	}

	return nil
}

func (fi *NodeFullIntro) DeserializeFrom(ctx DeserializeContext, reader io.Reader) error {
	if err := fi.NodeBriefIntro.DeserializeFrom(ctx, reader); err != nil {
		return errors.Wrap(err, "failed to deserialize NodeBriefIntro")
	}

	if err := fi.NodeExtendedIntro.DeserializeFrom(ctx, reader); err != nil {
		return errors.Wrap(err, "failed to deserialize NodeExtendedIntro")
	}

	return nil
}
