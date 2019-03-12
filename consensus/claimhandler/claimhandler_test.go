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
	"testing"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

func TestApprovedJoinersCount(t *testing.T) {
	assert.Equal(t, 1, ApprovedJoinersCount(1, 1))
	assert.Equal(t, 1, ApprovedJoinersCount(2, 3))
	assert.Equal(t, 3, ApprovedJoinersCount(5, 10))
	assert.Equal(t, 2, ApprovedJoinersCount(2, 10))
}

func TestClaimHandler_FilterClaims(t *testing.T) {
	// announcers references do not affect joiner claims filter logic, so choose random
	ref1 := testutils.RandomRef()
	ref2 := testutils.RandomRef()
	ref3 := testutils.RandomRef()

	claims := make(map[core.RecordRef][]packets.ReferendumClaim)
	claims[ref1] = []packets.ReferendumClaim{&packets.NodeBroadcast{}, &packets.NodeBroadcast{}, getJoinClaim(t, core.RecordRef{152})}
	claims[ref2] = []packets.ReferendumClaim{getJoinClaim(t, core.RecordRef{0}), getJoinClaim(t, core.RecordRef{154})}
	claims[ref3] = []packets.ReferendumClaim{getJoinClaim(t, core.RecordRef{1}), getJoinClaim(t, core.RecordRef{153})}

	containsJoinClaim := func(claims []packets.ReferendumClaim, ref core.RecordRef) bool {
		for _, claim := range claims {
			joinClaim, ok := claim.(*packets.NodeJoinClaim)
			if !ok {
				continue
			}
			if joinClaim.NodeRef.Equal(ref) {
				return true
			}
		}
		return false
	}

	handler := NewClaimHandler(6, claims)
	result := handler.FilterClaims([]core.RecordRef{ref1, ref2, ref3}, core.Entropy{0})
	// 2 NodeBroadcast + 2 JoinClaims
	assert.Equal(t, 4, len(result.ApprovedClaims))
	assert.True(t, containsJoinClaim(result.ApprovedClaims, core.RecordRef{154}))
	assert.True(t, containsJoinClaim(result.ApprovedClaims, core.RecordRef{153}))
	assert.False(t, containsJoinClaim(result.ApprovedClaims, core.RecordRef{0}))
	assert.False(t, containsJoinClaim(result.ApprovedClaims, core.RecordRef{1}))

	// 2 NodeBroadcast + 2 JoinClaims
	result = handler.FilterClaims([]core.RecordRef{ref1, ref2}, core.Entropy{0})
	assert.Equal(t, 4, len(result.ApprovedClaims))
	assert.True(t, containsJoinClaim(result.ApprovedClaims, core.RecordRef{154}))
	assert.True(t, containsJoinClaim(result.ApprovedClaims, core.RecordRef{152}))
	assert.False(t, containsJoinClaim(result.ApprovedClaims, core.RecordRef{0}))
	assert.False(t, containsJoinClaim(result.ApprovedClaims, core.RecordRef{1}))

	// only 2 JoinClaims
	result = handler.FilterClaims([]core.RecordRef{ref2, ref3}, core.Entropy{0})
	assert.Equal(t, 2, len(result.ApprovedClaims))
	assert.True(t, containsJoinClaim(result.ApprovedClaims, core.RecordRef{154}))
	assert.True(t, containsJoinClaim(result.ApprovedClaims, core.RecordRef{153}))
	assert.False(t, containsJoinClaim(result.ApprovedClaims, core.RecordRef{0}))
	assert.False(t, containsJoinClaim(result.ApprovedClaims, core.RecordRef{1}))
}

func TestClaimHandler_GetClaims(t *testing.T) {
	// announcers references do not affect joiner claims filter logic, so choose random
	ref1 := testutils.RandomRef()
	ref2 := testutils.RandomRef()
	ref3 := testutils.RandomRef()

	claims := make(map[core.RecordRef][]packets.ReferendumClaim)
	claims[ref1] = []packets.ReferendumClaim{&packets.NodeBroadcast{}, getJoinClaim(t, core.RecordRef{0})}
	claims[ref2] = []packets.ReferendumClaim{getJoinClaim(t, core.RecordRef{1}), getJoinClaim(t, core.RecordRef{2})}

	handler := NewClaimHandler(6, claims)
	assert.Equal(t, 4, len(handler.GetClaims()))

	handler.SetClaimsFromNode(ref3, []packets.ReferendumClaim{&packets.NodeBroadcast{}})
	assert.Equal(t, 5, len(handler.GetClaims()))
}
