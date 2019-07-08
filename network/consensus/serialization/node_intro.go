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
	"math/bits"

	"github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
	"github.com/pkg/errors"
)

const (
	primaryRoleMask    = 63 // 0b00111111
	primaryRoleBitSize = 6

	addrModeShift   = 6
	addrModeBitSize = 2
)

type NodeBriefIntro struct {
	// ByteSize= 135, 137, 147
	// ByteSize= 135 + (0, 2, 12)

	/*
		This field MUST be excluded from the packet, but considered for signature calculation.
		Value of this field equals SourceID or AnnounceID.
	*/
	ShortID common.ShortNodeID `insolar-transport:"ignore=send"` // ByteSize = 0

	PrimaryRoleAndFlags uint8 `insolar-transport:"[0:5]=header:NodePrimaryRole;[6:7]=header:AddrMode"` //AddrMode =0 reserved, =1 Relay, =2 IPv4 =3 IPv6
	SpecialRoles        common2.NodeSpecialRole
	StartPower          common2.MemberPower

	// 4 | 6 | 18 bytes
	// InboundRelayID common.ShortNodeID `insolar-transport:"AddrMode=2"`
	BasePort    uint16 `insolar-transport:"AddrMode=0,1"`
	PrimaryIPv4 uint32 `insolar-transport:"AddrMode=0"`
	// PrimaryIPv6    [4]uint32          `insolar-transport:"AddrMode=1"`

	// 128 bytes
	NodePK          common.Bits512 // works as a unique node identity
	JoinerSignature common.Bits512 // ByteSize=64
}

func (bi *NodeBriefIntro) GetPrimaryRole() common2.NodePrimaryRole {
	return common2.NodePrimaryRole(bi.PrimaryRoleAndFlags & primaryRoleMask)
}

func (bi *NodeBriefIntro) SetPrimaryRole(primaryRole common2.NodePrimaryRole) {
	if bits.Len(uint(primaryRole)) > primaryRoleBitSize {
		panic("invalid primary role")
	}

	bi.PrimaryRoleAndFlags |= uint8(primaryRole)
}
func (bi *NodeBriefIntro) GetAddrMode() common.NodeEndpointType {
	return common.NodeEndpointType(bi.PrimaryRoleAndFlags >> addrModeShift)
}

func (bi *NodeBriefIntro) SetAddrMode(addrMode common.NodeEndpointType) {
	if bits.Len(uint(addrMode)) > addrModeBitSize {
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

	if err := write(writer, bi.BasePort); err != nil {
		return errors.Wrap(err, "failed to serialize BasePort")
	}

	if err := write(writer, bi.PrimaryIPv4); err != nil {
		return errors.Wrap(err, "failed to serialize PrimaryIPv4")
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
	if err := read(reader, &bi.PrimaryRoleAndFlags); err != nil {
		return errors.Wrap(err, "failed to deserialize PrimaryRoleAndFlags")
	}

	if err := read(reader, &bi.SpecialRoles); err != nil {
		return errors.Wrap(err, "failed to deserialize SpecialRoles")
	}

	if err := read(reader, &bi.StartPower); err != nil {
		return errors.Wrap(err, "failed to deserialize StartPower")
	}

	if err := read(reader, &bi.BasePort); err != nil {
		return errors.Wrap(err, "failed to deserialize BasePort")
	}

	if err := read(reader, &bi.PrimaryIPv4); err != nil {
		return errors.Wrap(err, "failed to deserialize PrimaryIPv4")
	}

	if err := read(reader, &bi.NodePK); err != nil {
		return errors.Wrap(err, "failed to deserialize NodePK")
	}

	if err := read(reader, &bi.JoinerSignature); err != nil {
		return errors.Wrap(err, "failed to deserialize JoinerSignature")
	}

	return nil
}

type NodeFullIntro struct {
	// ByteSize= >=86 + (135, 137, 147) = >(221, 223, 233)

	NodeBriefIntro // ByteSize= 135, 137, 147

	// ByteSize>=86
	IssuedAtPulse common.PulseNumber // =0 when a node was connected during zeronet
	IssuedAtTime  uint64

	PowerLevels common2.MemberPowerSet // ByteSize=4

	EndpointLen    uint8
	ExtraEndpoints []uint16

	ProofLen     uint8
	NodeRefProof []common.Bits512

	DiscoveryIssuerNodeId common.ShortNodeID
	IssuerSignature       common.Bits512
}

func (fi *NodeFullIntro) SerializeTo(ctx SerializeContext, writer io.Writer) error {
	if err := fi.NodeBriefIntro.SerializeTo(ctx, writer); err != nil {
		return errors.Wrap(err, "failed to serialize NodeBriefIntro")
	}

	if err := write(writer, fi.IssuedAtPulse); err != nil {
		return errors.Wrap(err, "failed to serialize IssuedAtPulse")
	}

	if err := write(writer, fi.IssuedAtTime); err != nil {
		return errors.Wrap(err, "failed to serialize IssuedAtTime")
	}

	if err := write(writer, fi.PowerLevels); err != nil {
		return errors.Wrap(err, "failed to serialize PowerLevels")
	}

	if err := write(writer, fi.EndpointLen); err != nil {
		return errors.Wrap(err, "failed to serialize EndpointLen")
	}

	for i := 0; i < int(fi.EndpointLen); i++ {
		if err := write(writer, fi.ExtraEndpoints[i]); err != nil {
			return errors.Wrapf(err, "failed to serialize ExtraEndpoints[%d]", i)
		}
	}

	if err := write(writer, fi.ProofLen); err != nil {
		return errors.Wrap(err, "failed to serialize ProofLen")
	}

	for i := 0; i < int(fi.ProofLen); i++ {
		if err := write(writer, fi.NodeRefProof[i]); err != nil {
			return errors.Wrapf(err, "failed to serialize NodeRefProof[%d]", i)
		}
	}

	if err := write(writer, fi.DiscoveryIssuerNodeId); err != nil {
		return errors.Wrap(err, "failed to serialize DiscoveryIssuerNodeId")
	}

	if err := write(writer, fi.IssuerSignature); err != nil {
		return errors.Wrap(err, "failed to serialize IssuerSignature")
	}

	return nil
}

func (fi *NodeFullIntro) DeserializeFrom(ctx DeserializeContext, reader io.Reader) error {
	if err := fi.NodeBriefIntro.DeserializeFrom(ctx, reader); err != nil {
		return errors.Wrap(err, "failed to deserialize NodeBriefIntro")
	}

	if err := read(reader, &fi.IssuedAtPulse); err != nil {
		return errors.Wrap(err, "failed to deserialize IssuedAtPulse")
	}

	if err := read(reader, &fi.IssuedAtTime); err != nil {
		return errors.Wrap(err, "failed to deserialize IssuedAtTime")
	}

	if err := read(reader, &fi.PowerLevels); err != nil {
		return errors.Wrap(err, "failed to deserialize PowerLevels")
	}

	if err := read(reader, &fi.EndpointLen); err != nil {
		return errors.Wrap(err, "failed to deserialize EndpointLen")
	}

	fi.ExtraEndpoints = make([]uint16, fi.EndpointLen)
	for i := 0; i < int(fi.EndpointLen); i++ {
		if err := read(reader, &fi.ExtraEndpoints[i]); err != nil {
			return errors.Wrapf(err, "failed to deserialize ExtraEndpoints[%d]", i)
		}
	}

	if err := read(reader, &fi.ProofLen); err != nil {
		return errors.Wrap(err, "failed to deserialize ProofLen")
	}

	fi.NodeRefProof = make([]common.Bits512, fi.ProofLen)
	for i := 0; i < int(fi.EndpointLen); i++ {
		if err := read(reader, &fi.NodeRefProof[i]); err != nil {
			return errors.Wrapf(err, "failed to deserialize NodeRefProof[%d]", i)
		}
	}

	if err := read(reader, &fi.DiscoveryIssuerNodeId); err != nil {
		return errors.Wrap(err, "failed to deserialize DiscoveryIssuerNodeId")
	}

	if err := read(reader, &fi.IssuerSignature); err != nil {
		return errors.Wrap(err, "failed to deserialize IssuerSignature")
	}

	return nil
}
