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

type NodeVectors struct {
	// ByteSize=133..599
	/*
		GlobulaNodeBitset is a 5-state bitset, each node has a state at the same index as it was given in the current rank.
		Node have following states:
		0 - z-value (same as missing value) Trusted node
		1 - Doubted node
		2 -
		3 - Fraud node
		4 - Missing node
	*/
	StateVectorMask    *NodeAppearanceBitset `insolar-transport:"Packet=3"`                                                          // ByteSize=1..335
	TrustedStateVector *GlobulaStateVector   `insolar-transport:"Packet=3;TrustedStateVector.ExpectedRank[30-31]=flags:phase3Flags"` // ByteSize=132
	DoubtedStateVector *GlobulaStateVector   `insolar-transport:"optional=phase3Flags[0];Packet=3"`                                  // ByteSize=132
	// FraudStateVector *GlobulaStateVector `insolar-transport:"optional=phase3Flags[1];Packet=3"` //ByteSize=132
}

type NodeAppearanceBitset struct {
	// ByteSize=1..335
	FlagsAndLoLength uint8 // [00-05] LoByteLength, [06] Compressed, [07] HasHiLength (to be compatible with Protobuf VarInt)
	HiLength         uint8 // [00-06] HiByteLength, [07] MUST = 0 (to be compatible with Protobuf VarInt)
	Bytes            []byte
}

type GlobulaStateVector struct {
	// ByteSize=132
	ExpectedRank common2.MembershipRank // ByteSize=4
	/*
		GlobulaVectorHash = merkle(GlobulaNodeStateSignature of all nodes of this vector)
		SignedVectorHash = sign(GlobulaVectorHash, SK(sending node))
	*/
	SignedVectorHash common.Bits512 // ByteSize=64
	/*
		Hash(all MembershipUpdate.Signature of this vector)
	*/
	AnnouncementVectorHash common.Bits512
}
