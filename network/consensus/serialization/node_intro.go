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
	"github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
)

type NodeBriefIntro struct {
	// ByteSize= 74, 76, 88
	ShortID common.ShortNodeID
	NodeBriefIntroExt
}

type NodeBriefIntroExt struct {
	// ByteSize= 3 + (4, 6, 18) + 64 = 70, 72, 84

	//ShortID             common.ShortNodeID

	PrimaryRoleAndFlags uint8 `insolar-transport:"[0:5]=header:NodePrimaryRole;[6:7]=header:AddrMode"` //AddrMode =0 reserved, =1 Relay, =2 IPv4 =3 IPv6
	SpecialRoles        common2.NodeSpecialRole
	StartPower          common2.MemberPower

	// 4 | 6 | 18 bytes
	InboundRelayID common.ShortNodeID `insolar-transport:"AddrMode=2"`
	BasePort       uint16             `insolar-transport:"AddrMode=0,1"`
	PrimaryIPv4    uint32             `insolar-transport:"AddrMode=0"`
	PrimaryIPv6    [4]uint32          `insolar-transport:"AddrMode=1"`

	// 64 bytes
	NodePK common.Bits512 // works as a unique node identity
}

type NodeFullIntro struct {
	NodeBriefIntro
	NodeFullIntroExt
}

type NodeFullIntroExt struct {
	// ByteSize>=86
	IssuedAtPulse common.PulseNumber // =0 when a node was connected during zeronet
	IssuedAtTime  uint64

	PowerLevels common2.MemberPowerSet // ByteSize=4

	EndpointLen    uint8
	ExtraEndpoints []uint16

	ProofLen     uint8
	NodeRefProof []common.Bits512

	DiscoveryIssuerNodeId         common.ShortNodeID
	FullIntroSignatureByDiscovery common.Bits512
}
