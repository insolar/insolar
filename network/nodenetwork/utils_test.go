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

package nodenetwork

import (
	"testing"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/node"
	"github.com/stretchr/testify/assert"
)

func newTestNode(reference insolar.Reference, state insolar.NodeState) insolar.NetworkNode {
	return newTestNodeWithRole(reference, state, insolar.StaticRoleUnknown)
}

func newTestNodeWithRole(reference insolar.Reference, state insolar.NodeState, role insolar.StaticRole) insolar.NetworkNode {
	result := node.NewNode(reference, role, nil, "127.0.0.1:5432", "")
	result.(node.MutableNode).SetState(state)
	return result
}

func newTestJoinClaim(reference insolar.Reference) packets.ReferendumClaim {
	return &packets.NodeJoinClaim{NodeRef: reference}
}

func newTestLeaveClaim(reference insolar.Reference, ETA insolar.PulseNumber) packets.ReferendumClaim {
	return &packets.NodeLeaveClaim{NodeID: reference, ETA: ETA}
}

func Test_copyActiveNodes(t *testing.T) {
	nodes := []insolar.NetworkNode{
		newTestNode(insolar.Reference{0}, insolar.NodeUndefined),
		newTestNode(insolar.Reference{1}, insolar.NodePending),
		newTestNode(insolar.Reference{2}, insolar.NodeReady),
		newTestNode(insolar.Reference{3}, insolar.NodeLeaving),
	}

	copy := copyActiveNodes(nodes)
	assert.NotNil(t, copy[insolar.Reference{0}])
	// state changed during copy NodeUndefined -> NodePending
	assert.Equal(t, insolar.NodePending, copy[insolar.Reference{0}].GetState())
	assert.NotNil(t, copy[insolar.Reference{1}])
	// state changed during copy NodePending -> NodeReady
	assert.Equal(t, insolar.NodeReady, copy[insolar.Reference{1}].GetState())
	assert.NotNil(t, copy[insolar.Reference{2}])
	// state is not changed if node state is NodeReady or above
	assert.Equal(t, insolar.NodeReady, copy[insolar.Reference{2}].GetState())
	assert.NotNil(t, copy[insolar.Reference{3}])
	// state is not changed if node state is NodeReady or above
	assert.Equal(t, insolar.NodeLeaving, copy[insolar.Reference{3}].GetState())
	assert.Nil(t, copy[insolar.Reference{4}])
}

func TestGetMergedCopy_NilClaims(t *testing.T) {
	nodes := []insolar.NetworkNode{
		newTestNode(insolar.Reference{1}, insolar.NodePending),
		newTestNode(insolar.Reference{2}, insolar.NodeReady),
		newTestNode(insolar.Reference{3}, insolar.NodeLeaving),
	}

	result, err := GetMergedCopy(nodes, nil)
	assert.NoError(t, err)
	assert.False(t, result.NodesJoinedDuringPrevPulse)
	assert.Equal(t, 3, len(result.ActiveList))
}

func TestGetMergedCopy_JoinClaims(t *testing.T) {
	nodes := []insolar.NetworkNode{
		newTestNode(insolar.Reference{1}, insolar.NodePending),
		newTestNode(insolar.Reference{2}, insolar.NodeReady),
		newTestNode(insolar.Reference{3}, insolar.NodeLeaving),
	}
	claims := []packets.ReferendumClaim{
		newTestJoinClaim(insolar.Reference{4}),
		newTestJoinClaim(insolar.Reference{5}),
		&packets.NodeBroadcast{},
	}
	result, err := GetMergedCopy(nodes, claims)
	assert.NoError(t, err)
	assert.True(t, result.NodesJoinedDuringPrevPulse)
	assert.Equal(t, 5, len(result.ActiveList))
	assert.Equal(t, insolar.NodePending, result.ActiveList[insolar.Reference{4}].GetState())
	assert.Equal(t, insolar.NodePending, result.ActiveList[insolar.Reference{5}].GetState())
	assert.Equal(t, insolar.NodeReady, result.ActiveList[insolar.Reference{1}].GetState())
	assert.Equal(t, insolar.NodeReady, result.ActiveList[insolar.Reference{2}].GetState())
	assert.Equal(t, insolar.NodeLeaving, result.ActiveList[insolar.Reference{3}].GetState())
}

func TestGetMergedCopy_LeaveClaims(t *testing.T) {
	nodes := []insolar.NetworkNode{
		newTestNode(insolar.Reference{1}, insolar.NodePending),
		newTestNode(insolar.Reference{2}, insolar.NodeReady),
		newTestNode(insolar.Reference{3}, insolar.NodeLeaving),
	}
	claims := []packets.ReferendumClaim{
		newTestJoinClaim(insolar.Reference{4}),
		newTestLeaveClaim(insolar.Reference{2}, insolar.PulseNumber(50)),
		&packets.NodeBroadcast{},
		&packets.NodeBroadcast{},
	}
	result, err := GetMergedCopy(nodes, claims)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(result.ActiveList))
	assert.Equal(t, insolar.NodePending, result.ActiveList[insolar.Reference{4}].GetState())
	assert.Equal(t, insolar.NodeReady, result.ActiveList[insolar.Reference{1}].GetState())
	assert.Equal(t, insolar.NodeLeaving, result.ActiveList[insolar.Reference{2}].GetState())
	assert.Equal(t, insolar.PulseNumber(50), result.ActiveList[insolar.Reference{2}].LeavingETA())
	assert.Equal(t, insolar.NodeLeaving, result.ActiveList[insolar.Reference{3}].GetState())
}

func TestGetMergedCopy_InvalidLeaveClaim(t *testing.T) {
	nodes := []insolar.NetworkNode{
		newTestNode(insolar.Reference{1}, insolar.NodePending),
		newTestNode(insolar.Reference{2}, insolar.NodeReady),
		newTestNode(insolar.Reference{3}, insolar.NodeLeaving),
	}
	claims := []packets.ReferendumClaim{
		newTestLeaveClaim(insolar.Reference{5}, insolar.PulseNumber(50)),
	}
	result, err := GetMergedCopy(nodes, claims)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(result.ActiveList))
}
