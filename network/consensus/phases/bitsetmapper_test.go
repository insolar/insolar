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

package phases

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/packets"
	"github.com/insolar/insolar/network/node"
	"github.com/stretchr/testify/assert"
)

func TestNewBitsetMapper(t *testing.T) {
	nodes := []insolar.NetworkNode{
		node.NewNode(insolar.Reference{0}, insolar.StaticRoleVirtual, nil, "127.0.0.1:0", ""),
		node.NewNode(insolar.Reference{5}, insolar.StaticRoleLightMaterial, nil, "127.0.0.1:0", ""),
		node.NewNode(insolar.Reference{1}, insolar.StaticRoleHeavyMaterial, nil, "127.0.0.1:0", ""),
		node.NewNode(insolar.Reference{3}, insolar.StaticRoleVirtual, nil, "127.0.0.1:0", ""),
	}

	bm := NewBitsetMapper(nodes)
	assert.Equal(t, 4, bm.Length())
	_, err := bm.IndexToRef(-1)
	assert.Equal(t, packets.ErrBitSetOutOfRange, err)
	_, err = bm.IndexToRef(5)
	assert.Equal(t, packets.ErrBitSetOutOfRange, err)

	// keep in mind that bitset holds nodes sorted by their references
	ref, err := bm.IndexToRef(2)
	assert.NoError(t, err)
	assert.Equal(t, insolar.Reference{3}, ref)
	ref, err = bm.IndexToRef(3)
	assert.NoError(t, err)
	assert.Equal(t, insolar.Reference{5}, ref)

	_, err = bm.RefToIndex(insolar.Reference{2})
	assert.Equal(t, packets.ErrBitSetIncorrectNode, err)
	index, err := bm.RefToIndex(insolar.Reference{5})
	assert.NoError(t, err)
	assert.Equal(t, 3, index)
	index, err = bm.RefToIndex(insolar.Reference{1})
	assert.NoError(t, err)
	assert.Equal(t, 1, index)
}

func TestNewSparseBitsetMapper(t *testing.T) {
	bm := NewSparseBitsetMapper(5)
	node1 := node.NewNode(insolar.Reference{11}, insolar.StaticRoleHeavyMaterial, nil, "127.0.0.1:0", "")
	bm.AddNode(node1, 1)
	node3 := node.NewNode(insolar.Reference{13}, insolar.StaticRoleVirtual, nil, "127.0.0.1:0", "")
	bm.AddNode(node3, 3)
	node4 := node.NewNode(insolar.Reference{14}, insolar.StaticRoleLightMaterial, nil, "127.0.0.1:0", "")
	bm.AddNode(node4, 4)

	assert.Equal(t, 5, bm.Length())
	_, err := bm.IndexToRef(-1)
	assert.Equal(t, packets.ErrBitSetOutOfRange, err)
	_, err = bm.IndexToRef(6)
	assert.Equal(t, packets.ErrBitSetOutOfRange, err)
	_, err = bm.IndexToRef(0)
	assert.Equal(t, packets.ErrBitSetNodeIsMissing, err)

	ref, err := bm.IndexToRef(3)
	assert.NoError(t, err)
	assert.Equal(t, node3.ID(), ref)
	ref, err = bm.IndexToRef(4)
	assert.NoError(t, err)
	assert.Equal(t, node4.ID(), ref)

	_, err = bm.RefToIndex(insolar.Reference{22})
	assert.Equal(t, packets.ErrBitSetIncorrectNode, err)
	index, err := bm.RefToIndex(node3.ID())
	assert.NoError(t, err)
	assert.Equal(t, 3, index)
	index, err = bm.RefToIndex(node1.ID())
	assert.NoError(t, err)
	assert.Equal(t, 1, index)
}

func newTestAnnounceClaim(announcerIndex, joinerIndex, count uint16, announcer insolar.Reference) *packets.NodeAnnounceClaim {
	result := &packets.NodeAnnounceClaim{}
	result.NodeJoinClaim.NodeRef = announcer
	result.NodeAnnouncerIndex = announcerIndex
	result.NodeJoinerIndex = joinerIndex
	result.NodeCount = count
	result.NodeAddress = packets.NewNodeAddress("127.0.0.1:0")
	return result
}

func TestApplyClaims(t *testing.T) {
	mutator := node.NewMutator(node.NewSnapshot(insolar.FirstPulseNumber, nil))
	count := 10
	bm := NewSparseBitsetMapper(count)
	announce1 := newTestAnnounceClaim(1, 5, uint16(count), insolar.Reference{11})
	announce2 := newTestAnnounceClaim(2, 5, uint16(count), insolar.Reference{22})
	announce3 := newTestAnnounceClaim(9, 5, uint16(count), insolar.Reference{99})
	cs := ConsensusState{NodesMutator: mutator, BitsetMapper: bm}
	origin := node.NewNode(insolar.Reference{55}, insolar.StaticRoleLightMaterial, nil, "127.0.0.1:0", "")
	claims := []packets.ReferendumClaim{announce1, announce2, announce3}
	err := ApplyClaims(&cs, origin, claims)
	assert.NoError(t, err)

	assert.NotNil(t, mutator.GetActiveNode(insolar.Reference{11}))
	assert.NotNil(t, mutator.GetActiveNode(insolar.Reference{22}))
	assert.NotNil(t, mutator.GetActiveNode(insolar.Reference{55}))
	assert.Nil(t, mutator.GetActiveNode(insolar.Reference{33}))

	ref, err := bm.IndexToRef(2)
	assert.NoError(t, err)
	assert.Equal(t, insolar.Reference{22}, ref)
	ref, err = bm.IndexToRef(9)
	assert.NoError(t, err)
	assert.Equal(t, insolar.Reference{99}, ref)
	_, err = bm.IndexToRef(3)
	assert.Equal(t, packets.ErrBitSetNodeIsMissing, err)

	_, err = bm.RefToIndex(insolar.Reference{41})
	assert.Equal(t, packets.ErrBitSetIncorrectNode, err)
	index, err := bm.RefToIndex(insolar.Reference{11})
	assert.NoError(t, err)
	assert.Equal(t, 1, index)
	index, err = bm.RefToIndex(insolar.Reference{99})
	assert.NoError(t, err)
	assert.Equal(t, 9, index)
}
