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
	"sort"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/packets"
	"github.com/insolar/insolar/network/node"
	"github.com/pkg/errors"
)

type BitsetMapper struct {
	length     int
	refToIndex map[insolar.Reference]int
	indexToRef map[int]insolar.Reference
}

func NewBitsetMapper(activeNodes []insolar.NetworkNode) *BitsetMapper {
	bm := NewSparseBitsetMapper(len(activeNodes))
	sort.Slice(activeNodes, func(i, j int) bool {
		return activeNodes[i].ID().Compare(activeNodes[j].ID()) < 0
	})
	for i, node := range activeNodes {
		bm.AddNode(node, uint16(i))
	}
	return bm
}

func NewSparseBitsetMapper(length int) *BitsetMapper {
	return &BitsetMapper{
		length:     length,
		refToIndex: make(map[insolar.Reference]int),
		indexToRef: make(map[int]insolar.Reference),
	}
}

func (bm *BitsetMapper) AddNode(node insolar.NetworkNode, bitsetIndex uint16) {
	index := int(bitsetIndex)
	bm.indexToRef[index] = node.ID()
	bm.refToIndex[node.ID()] = index
}

func (bm *BitsetMapper) IndexToRef(index int) (insolar.Reference, error) {
	if index < 0 || index >= bm.length {
		return insolar.Reference{}, packets.ErrBitSetOutOfRange
	}
	result, ok := bm.indexToRef[index]
	if !ok {
		return insolar.Reference{}, packets.ErrBitSetNodeIsMissing
	}
	return result, nil
}

func (bm *BitsetMapper) RefToIndex(nodeID insolar.Reference) (int, error) {
	index, ok := bm.refToIndex[nodeID]
	if !ok {
		return 0, packets.ErrBitSetIncorrectNode
	}
	return index, nil
}

func (bm *BitsetMapper) Length() int {
	return bm.length
}

// ApplyClaims processes Announce Claims and modifies  ConsensusState bitset and Mutator
func ApplyClaims(state *ConsensusState, origin insolar.NetworkNode, claims []packets.ReferendumClaim) error {
	var NodeJoinerIndex uint16
	for _, claim := range claims {
		c, ok := claim.(*packets.NodeAnnounceClaim)
		if !ok {
			continue
		}

		// TODO: fix version
		n, err := node.ClaimToNode("", c)
		if err != nil {
			return errors.Wrap(err, "[ AddClaims ] failed to convert Claim -> Node")
		}
		// TODO: check bitset indexes from every announce claim for fraud
		NodeJoinerIndex = c.NodeJoinerIndex
		state.BitsetMapper.AddNode(n, c.NodeAnnouncerIndex)
		if c.LeavingETA > 0 {
			state.NodesMutator.AddNode(n, node.ListLeaving)
		} else {
			state.NodesMutator.AddNode(n, node.ListWorking)
		}
	}
	state.BitsetMapper.AddNode(origin, NodeJoinerIndex)
	state.NodesMutator.AddNode(origin, node.ListWorking)

	return nil
}
