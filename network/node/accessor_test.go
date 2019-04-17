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

package node

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessor(t *testing.T) {
	m := make(map[insolar.Reference]insolar.NetworkNode)

	node := newMutableNode(testutils.RandomRef(), insolar.StaticRoleVirtual, nil, insolar.NodeReady, "127.0.0.1:0", "")
	m[node.ID()] = node

	node2 := newMutableNode(testutils.RandomRef(), insolar.StaticRoleVirtual, nil, insolar.NodePending, "127.0.0.1:0", "")
	node2.SetShortID(11)
	m[node2.ID()] = node2

	node3 := newMutableNode(testutils.RandomRef(), insolar.StaticRoleVirtual, nil, insolar.NodeLeaving, "127.0.0.1:0", "")
	node3.SetShortID(10)
	m[node3.ID()] = node3

	node4 := newMutableNode(testutils.RandomRef(), insolar.StaticRoleVirtual, nil, insolar.NodeUndefined, "127.0.0.1:0", "")
	m[node4.ID()] = node4

	snapshot := NewSnapshot(insolar.FirstPulseNumber, m)
	accessor := NewAccessor(snapshot)
	assert.Equal(t, 4, len(accessor.GetActiveNodes()))
	assert.Equal(t, 1, len(accessor.GetWorkingNodes()))
	assert.NotNil(t, accessor.GetWorkingNode(node.ID()))
	assert.Nil(t, accessor.GetWorkingNode(node2.ID()))
	assert.NotNil(t, accessor.GetActiveNode(node2.ID()))
	assert.NotNil(t, accessor.GetActiveNode(node3.ID()))
	assert.NotNil(t, accessor.GetActiveNode(node4.ID()))

	assert.Empty(t, accessor.GetWorkingNodesByRole(insolar.DynamicRoleLightValidator))
	assert.Equal(t, 1, len(accessor.GetWorkingNodesByRole(insolar.DynamicRoleVirtualValidator)))

	assert.NotNil(t, accessor.GetActiveNodeByShortID(10))
	assert.NotNil(t, accessor.GetActiveNodeByShortID(11))
	assert.Nil(t, accessor.GetActiveNodeByShortID(12))
}

func Test_dynamicToStaticRole(t *testing.T) {
	assert.Equal(t, insolar.StaticRoleVirtual, dynamicToStaticRole(insolar.DynamicRoleVirtualExecutor))
	assert.Equal(t, insolar.StaticRoleLightMaterial, dynamicToStaticRole(insolar.DynamicRoleLightExecutor))
	assert.Equal(t, insolar.StaticRoleHeavyMaterial, dynamicToStaticRole(insolar.DynamicRoleHeavyExecutor))
	assert.Equal(t, insolar.StaticRoleUnknown, dynamicToStaticRole(150))
}

func TestMutator(t *testing.T) {
	snapshot := NewSnapshot(insolar.FirstPulseNumber, nil)
	mutator := NewMutator(snapshot)
	node1 := newMutableNode(insolar.Reference{22}, insolar.StaticRoleVirtual, nil, insolar.NodeReady, "127.0.0.1:0", "")
	mutator.AddWorkingNode(node1)
	node2 := newMutableNode(insolar.Reference{11}, insolar.StaticRoleLightMaterial, nil, insolar.NodeReady, "127.0.0.1:0", "")
	mutator.AddWorkingNode(node2)

	assert.Equal(t, 2, len(mutator.GetActiveNodes()))
	// duplicating key
	node3 := newMutableNode(insolar.Reference{11}, insolar.StaticRoleHeavyMaterial, nil, insolar.NodeReady, "127.0.0.1:0", "")
	mutator.AddWorkingNode(node3)
	assert.Equal(t, 2, len(mutator.GetActiveNodes()))

	nodes := mutator.GetWorkingNodes()
	require.Equal(t, 2, len(mutator.GetWorkingNodes()))
	assert.Equal(t, insolar.Reference{11}, nodes[0].ID())
	assert.Equal(t, insolar.Reference{22}, nodes[1].ID())
}
