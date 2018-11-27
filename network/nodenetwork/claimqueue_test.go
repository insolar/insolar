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
