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

package claimhandler

import (
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/utils"
)

// NodesToJoinPercent how many nodes from active list can connect to the network.
const NodesToJoinPercent = 1.0 / 3.0

func maxJoinersForPulse(activeNodesCount int) int {
	nodesToJoin := int(float64(activeNodesCount) * NodesToJoinPercent)
	if nodesToJoin == 0 {
		nodesToJoin++
	}
	return nodesToJoin
}

func min(first, second int) int {
	if first < second {
		return first
	}
	return second
}

func ApprovedJoinersCount(requestedJoinersCount, activeNodesCount int) int {
	return min(requestedJoinersCount, maxJoinersForPulse(activeNodesCount))
}

type ClaimHandler struct {
	claims      map[insolar.RecordRef][]packets.ReferendumClaim
	activeCount int
}

func NewClaimHandler(activeNodesCount int, claims map[insolar.RecordRef][]packets.ReferendumClaim) *ClaimHandler {
	return &ClaimHandler{
		activeCount: activeNodesCount,
		claims:      claims,
	}
}

func (ch *ClaimHandler) SetClaimsFromNode(node insolar.RecordRef, claims []packets.ReferendumClaim) {
	ch.claims[node] = claims
}

func (ch *ClaimHandler) GetClaimsFromNode(node insolar.RecordRef) []packets.ReferendumClaim {
	return ch.claims[node]
}

func (ch *ClaimHandler) GetClaims() []packets.ReferendumClaim {
	result := make([]packets.ReferendumClaim, 0)
	for _, claims := range ch.claims {
		result = append(result, claims...)
	}
	return result
}

type ClaimSplit struct {
	ApprovedClaims []packets.ReferendumClaim
	// TODO: add logic to return unallowed local claims back to ClaimQueue
	LeftForNextPulse []packets.ReferendumClaim
}

type none struct{}
type recordRefSet map[insolar.RecordRef]none

func (ch *ClaimHandler) FilterClaims(approvedNodes []insolar.RecordRef, entropy insolar.Entropy) ClaimSplit {
	knownClaims := make(recordRefSet)
	queue := Queue{}

	for _, node := range approvedNodes {
		addKnownClaimsToQueue(&queue, knownClaims, ch.GetClaimsFromNode(node), entropy)
	}

	joinClaims := ch.getApprovedJoinClaims(&queue)
	joinClaimsSet := make(recordRefSet)
	for _, joinClaim := range joinClaims {
		joinClaimsSet[joinClaim.NodeRef] = none{}
	}

	approvedClaims := make([]packets.ReferendumClaim, 0)
	for _, node := range approvedNodes {
		approvedClaimsForNode := getApprovedClaimsForNode(joinClaimsSet, ch.GetClaimsFromNode(node))
		approvedClaims = append(approvedClaims, approvedClaimsForNode...)
	}

	return ClaimSplit{
		ApprovedClaims: approvedClaims,
	}
}

func getApprovedClaimsForNode(approvedJoinClaims recordRefSet, claimsForNode []packets.ReferendumClaim) []packets.ReferendumClaim {
	result := make([]packets.ReferendumClaim, 0)
	for _, claim := range claimsForNode {
		joinClaim, ok := claim.(*packets.NodeJoinClaim)
		if !ok {
			result = append(result, claim)
			continue
		}
		_, ok = approvedJoinClaims[joinClaim.NodeRef]
		if ok {
			result = append(result, claim)
		}
	}
	return result
}

func addKnownClaimsToQueue(queue *Queue, knownClaims recordRefSet, claims []packets.ReferendumClaim, entropy insolar.Entropy) {
	for _, claim := range claims {
		join, ok := claim.(*packets.NodeJoinClaim)
		if !ok {
			continue
		}
		_, ok = knownClaims[join.NodeRef]
		if ok {
			continue
		}
		priority := getPriority(join.NodeRef, entropy)
		queue.PushClaim(claim, priority)
		knownClaims[join.NodeRef] = none{}
	}
}

func (ch *ClaimHandler) getApprovedJoinClaims(queue *Queue) []*packets.NodeJoinClaim {
	res := make([]*packets.NodeJoinClaim, 0)
	nodesToJoin := ApprovedJoinersCount(queue.Len(), ch.activeCount)
	for i := 0; i < nodesToJoin; i++ {
		res = append(res, queue.PopClaim().(*packets.NodeJoinClaim))
	}
	return res
}

func getPriority(ref insolar.RecordRef, entropy insolar.Entropy) []byte {
	return utils.CircleXOR(ref[:], entropy[:])
}
