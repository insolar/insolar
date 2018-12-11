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

// claims auxiliar constants
const (
	headerTypeShift = 10
	headerTypeMask  = 0xfc00
	// headerLengthMask = 0x3ff
)

const HeaderSize = 2

func extractTypeFromHeader(claimHeader uint16) uint8 {
	return uint8((claimHeader & headerTypeMask) >> headerTypeShift)
}

// func extractLengthFromHeader(header uint16) uint16 {
// 	return header & headerLengthMask
// }

func makeClaimHeader(claim ReferendumClaim) uint16 {
	if claim == nil {
		panic("invalid claim")
	}
	var result = getClaimSize(claim)
	result |= uint16(claim.Type()) << headerTypeShift
	return result
}

func getClaimSize(claim ReferendumClaim) uint16 {
	return claimSizeMap[claim.Type()]
}

func getClaimWithHeaderSize(claim ReferendumClaim) uint16 {
	return getClaimSize(claim) + HeaderSize
}

func makeVoteHeader(vote ReferendumVote) uint16 {
	if vote == nil {
		panic("invalid vote")
	}
	var result = getVoteSize(vote)
	result |= uint16(vote.Type()) << headerTypeShift
	return result
}

func getVoteSize(vote ReferendumVote) uint16 {
	return voteSizeMap[vote.Type()]
}
