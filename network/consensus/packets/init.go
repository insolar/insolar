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

package packets

import (
	"github.com/insolar/insolar/log"
)

// packetMaxSize should be less then MTU
const packetMaxSize = 1400

var (
	phase1PacketSizeForClaims int
	// claimSizeMap contains serialized size of each claim type without header(2 bytes)
	claimSizeMap map[ClaimType]uint16
	// claimSizeMap contains sizes of serialized votes for each type without header (2 bytes)
	voteSizeMap map[VoteType]uint16
)

// init packets and claims size variables
func init() {
	sizeOf := func(s Serializer) uint16 {
		data, err := s.Serialize()
		if err != nil {
			log.Fatal("Failed to init packets package: ", err.Error())
		}
		return uint16(len(data))
	}

	phase1PacketSizeForClaims = packetMaxSize - int(sizeOf(&Phase1Packet{}))

	claimSizeMap = make(map[ClaimType]uint16)
	claimSizeMap[TypeNodeJoinClaim] = sizeOf(&NodeJoinClaim{})
	claimSizeMap[TypeNodeAnnounceClaim] = sizeOf(&NodeAnnounceClaim{})
	claimSizeMap[TypeNodeLeaveClaim] = sizeOf(&NodeLeaveClaim{})

	voteSizeMap = make(map[VoteType]uint16)
	voteSizeMap[TypeMissingNodeRespVote] = sizeOf(&MissingNodeRespVote{})
	voteSizeMap[TypeMissingNodeReqVote] = sizeOf(&MissingNodeReqVote{})
	// MissingNodeClaimsVote has variable-sized length
	voteSizeMap[TypeMissingNodeClaimsVote] = 0
}
