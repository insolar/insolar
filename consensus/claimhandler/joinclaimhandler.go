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
	"math"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/utils"
)

// NodesToJoinPercent how many nodes from active list can connect to the network.
const NodesToJoinPercent = 1.0 / 3.0

type JoinClaimHandler struct {
	queue       Queue
	knownClaims map[core.RecordRef]bool
	activeCount int
}

func NewJoinHandler(activeNodesCount int) *JoinClaimHandler {
	handler := &JoinClaimHandler{
		queue:       Queue{},
		activeCount: activeNodesCount,
	}
	handler.knownClaims = make(map[core.RecordRef]bool)
	return handler
}

func (jch *JoinClaimHandler) AddClaims(claims []packets.ReferendumClaim, entropy core.Entropy) {
	for _, claim := range claims {
		join, ok := claim.(*packets.NodeJoinClaim)
		if !ok || jch.isKnownClaim(join) {
			continue
		}
		priority := getPriority(join.NodeRef, entropy)
		jch.queue.PushClaim(claim, priority)
		jch.knownClaims[join.NodeRef] = true
	}
}

func (jch *JoinClaimHandler) HandleAndReturnClaims() []*packets.NodeJoinClaim {
	return jch.getClaimsByPriority()
}

func (jch *JoinClaimHandler) getClaimsByPriority() []*packets.NodeJoinClaim {
	res := make([]*packets.NodeJoinClaim, 0)
	nodesToJoin := float64(jch.activeCount) * NodesToJoinPercent

	if nodesToJoin == 0 {
		nodesToJoin++
	}
	queueLen := float64(jch.queue.Len())
	for i := 0; i < int(math.Min(nodesToJoin, queueLen)); i++ {
		res = append(res, jch.queue.PopClaim().(*packets.NodeJoinClaim))
	}

	return res
}

func getPriority(ref core.RecordRef, entropy core.Entropy) []byte {
	return utils.CircleXOR(ref[:], entropy[:])
}

func (jch *JoinClaimHandler) isKnownClaim(claim *packets.NodeJoinClaim) bool {
	_, ok := jch.knownClaims[claim.NodeRef]
	return ok
}
