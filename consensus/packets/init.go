/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

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
			log.Fatalln("Failed to init packets package: ", err.Error())
		}
		return uint16(len(data))
	}

	phase1PacketSizeForClaims = packetMaxSize - int(sizeOf(&Phase1Packet{}))

	claimSizeMap = make(map[ClaimType]uint16)
	claimSizeMap[TypeNodeJoinClaim] = sizeOf(&NodeJoinClaim{})
	claimSizeMap[TypeNodeAnnounceClaim] = sizeOf(&NodeJoinClaim{})
	claimSizeMap[TypeCapabilityPollingAndActivation] = sizeOf(&CapabilityPoolingAndActivation{})
	claimSizeMap[TypeNodeViolationBlame] = sizeOf(&NodeViolationBlame{})
	claimSizeMap[TypeNodeBroadcast] = sizeOf(&NodeBroadcast{})
	claimSizeMap[TypeNodeLeaveClaim] = sizeOf(&NodeLeaveClaim{})
	claimSizeMap[TypeChangeNetworkClaim] = sizeOf(&NodeLeaveClaim{})

	voteSizeMap = make(map[VoteType]uint16)
	voteSizeMap[TypeNodeJoinSupplementaryVote] = sizeOf(&NodeJoinSupplementaryVote{})
	voteSizeMap[TypeStateFraudNodeSupplementaryVote] = sizeOf(&StateFraudNodeSupplementaryVote{})
	voteSizeMap[TypeNodeListSupplementaryVote] = sizeOf(&NodeListSupplementaryVote{})
	voteSizeMap[TypeMissingNodeSupplementaryVote] = sizeOf(&MissingNodeSupplementaryVote{})
	voteSizeMap[TypeMissingNode] = sizeOf(&MissingNode{})
}
