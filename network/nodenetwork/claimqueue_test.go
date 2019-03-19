/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted (subject to the limitations in the disclaimer below) provided that
 * the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *  * Neither the name of Insolar Technologies nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
 * BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
 * CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING,
 * BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package nodenetwork

import (
	"testing"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/stretchr/testify/assert"
)

func newTestClaim(claimType packets.ClaimType) packets.ReferendumClaim {
	switch claimType {
	case packets.TypeNodeJoinClaim:
		return &packets.NodeJoinClaim{}
	case packets.TypeCapabilityPollingAndActivation:
		return &packets.CapabilityPoolingAndActivation{}
	case packets.TypeNodeViolationBlame:
		return &packets.NodeViolationBlame{}
	case packets.TypeNodeBroadcast:
		return &packets.NodeBroadcast{}
	case packets.TypeNodeLeaveClaim:
		return &packets.NodeLeaveClaim{}
	}
	return nil
}

func TestClaimQueue_Pop(t *testing.T) {
	cq := newClaimQueue()
	assert.Equal(t, 0, cq.Length())
	assert.Nil(t, cq.Front())
	assert.Nil(t, cq.Pop())

	cq.Push(newTestClaim(packets.TypeNodeJoinClaim))
	cq.Push(newTestClaim(packets.TypeNodeBroadcast))
	assert.Equal(t, 2, cq.Length())

	assert.NotNil(t, cq.Front())
	assert.Equal(t, packets.TypeNodeJoinClaim, cq.Front().Type())

	assert.Equal(t, packets.TypeNodeJoinClaim, cq.Pop().Type())
	assert.Equal(t, packets.TypeNodeBroadcast, cq.Pop().Type())
}
