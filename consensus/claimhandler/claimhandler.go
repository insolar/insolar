/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package claimhandler

import (
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/utils"
)

// NodesToJoinPercent how many nodes from active list can connect to the network.
const NodesToJoinPercent = 1.0 / 3.0

type ClaimHandler struct {
	queue       Queue
	claims      map[core.RecordRef][]packets.ReferendumClaim
	knownClaims map[core.RecordRef]bool
	activeCount int
}

func NewClaimHandler(activeNodesCount int, claims map[core.RecordRef][]packets.ReferendumClaim) *ClaimHandler {
	return &ClaimHandler{
		queue:       Queue{},
		activeCount: activeNodesCount,
		knownClaims: make(map[core.RecordRef]bool),
		claims:      claims,
	}
}

func (ch *ClaimHandler) SetClaimsFromNode(node core.RecordRef, claims []packets.ReferendumClaim) {
	ch.claims[node] = claims
}

func (ch *ClaimHandler) GetClaimsFromNode(node core.RecordRef) []packets.ReferendumClaim {
	return ch.claims[node]
}

func (ch *ClaimHandler) AddKnownClaims(claims []packets.ReferendumClaim, entropy core.Entropy) {
	for _, claim := range claims {
		join, ok := claim.(*packets.NodeJoinClaim)
		if !ok || ch.isKnownClaim(join) {
			continue
		}
		priority := getPriority(join.NodeRef, entropy)
		ch.queue.PushClaim(claim, priority)
		ch.knownClaims[join.NodeRef] = true
	}
}

func (ch *ClaimHandler) HandleAndReturnClaims() []*packets.NodeJoinClaim {
	return ch.getClaimsByPriority()
}

func (ch *ClaimHandler) getClaimsByPriority() []*packets.NodeJoinClaim {
	res := make([]*packets.NodeJoinClaim, 0)
	nodesToJoin := int(float64(ch.activeCount) * NodesToJoinPercent)

	if nodesToJoin == 0 {
		nodesToJoin++
	}

	min := func(first, second int) int {
		if first < second {
			return first
		}
		return second
	}

	for i := 0; i < min(nodesToJoin, ch.queue.Len()); i++ {
		res = append(res, ch.queue.PopClaim().(*packets.NodeJoinClaim))
	}

	return res
}

func getPriority(ref core.RecordRef, entropy core.Entropy) []byte {
	return utils.CircleXOR(ref[:], entropy[:])
}

func (ch *ClaimHandler) isKnownClaim(claim *packets.NodeJoinClaim) bool {
	_, ok := ch.knownClaims[claim.NodeRef]
	return ok
}
