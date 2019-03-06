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
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
)

func TestJoinClaimHandler_HandleClaim(t *testing.T) {
	activeNodesCount := 5
	claimsCount := 10
	entropy := core.Entropy{}
	_, err := rand.Read(entropy[:])
	assert.NoError(t, err)

	claims := make([]packets.ReferendumClaim, claimsCount)
	priorityMap := make(map[core.RecordRef][]byte, claimsCount)
	for i := 0; i < claimsCount; i++ {
		claim := getJoinClaim(t)
		claims[i] = claim
		priorityMap[claim.NodeRef] = getPriority(claim.NodeRef, entropy)
	}

	handler := NewJoinHandler(activeNodesCount)
	handler.AddClaims(claims, entropy)
	res := handler.HandleAndReturnClaims()
	assert.Len(t, res, int(float64(activeNodesCount)*NodesToJoinPercent))

	for _, claim := range res {
		delete(priorityMap, claim.NodeRef)
	}

	for i := len(res) - 1; i >= 0; i-- {
		highPriority := getPriority(res[i].NodeRef, entropy)
		for _, claim := range priorityMap {
			assert.True(t, bytes.Compare(highPriority, claim) >= 0)
		}
	}
}
