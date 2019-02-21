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
	"context"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

const NodesToJoinPercent = 0.3

type JoinHandler struct {
	queue       Queue
	ref         core.RecordRef
	activeCount int
}

func NewJoinHandler(activeNodesCount int) *JoinHandler {
	handler := &JoinHandler{
		queue:       Queue{},
		activeCount: activeNodesCount,
	}
	return handler
}

func (jch *JoinHandler) HandleClaims(claims []*packets.NodeJoinClaim, entropy core.Entropy) []*packets.NodeJoinClaim {
	for _, claim := range claims {
		priority := getPriority(claim.NodeRef, entropy)
		jch.queue.PushClaim(claim, priority)
	}
	return jch.getClaimsByPriority()
}

func (jch *JoinHandler) getClaimsByPriority() []*packets.NodeJoinClaim {
	res := make([]*packets.NodeJoinClaim, 0)
	nodesToJoin := int(float64(jch.activeCount) * NodesToJoinPercent)

	if nodesToJoin == 0 {
		nodesToJoin++
	}
	logger := inslogger.FromContext(context.Background())
	for i := 0; i < nodesToJoin; i++ {
		if i >= jch.queue.Len() {
			break
		}
		res = append(res, jch.queue.PopClaim().(*packets.NodeJoinClaim))
	}
	logger.Debugf("[ getClaimsByPriority ] handle join claims. max nodes to join: %d, join nodes count: %d", nodesToJoin, len(res))

	return res
}

func getPriority(ref core.RecordRef, entropy core.Entropy) []byte {
	// TODO: try to delete this if block, but be careful
	if len(ref) != len(entropy) {
		logger := inslogger.FromContext(context.Background())
		logger.Errorf("[ joinClaimHandler ] getPriority: length not match! reference: %d, entropy: %d", len(ref), len(entropy))
	}
	res := make([]byte, len(ref))
	for i := 0; i < len(ref); i++ {
		res[i] = ref[i] ^ entropy[i]
	}
	return res
}
