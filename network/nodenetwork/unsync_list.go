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
	"sort"

	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
	"github.com/pkg/errors"
)

func copyActiveNodes(nodes []insolar.NetworkNode) map[insolar.RecordRef]insolar.NetworkNode {
	result := make(map[insolar.RecordRef]insolar.NetworkNode, len(nodes))
	for _, node := range nodes {
		node.(MutableNode).ChangeState()
		result[node.ID()] = node
	}
	return result
}

type unsyncList struct {
	length      int
	origin      insolar.NetworkNode
	activeNodes map[insolar.RecordRef]insolar.NetworkNode
	refToIndex  map[insolar.RecordRef]int
	proofs      map[insolar.RecordRef]*consensus.NodePulseProof
	ghs         map[insolar.RecordRef]consensus.GlobuleHashSignature
	indexToRef  map[int]insolar.RecordRef
}

func (ul *unsyncList) GetOrigin() insolar.NetworkNode {
	return ul.origin
}

func (ul *unsyncList) GetGlobuleHashSignature(ref insolar.RecordRef) (consensus.GlobuleHashSignature, bool) {
	ghs, ok := ul.ghs[ref]
	return ghs, ok
}

func (ul *unsyncList) SetGlobuleHashSignature(ref insolar.RecordRef, ghs consensus.GlobuleHashSignature) {
	ul.ghs[ref] = ghs
}

func (ul *unsyncList) RemoveNode(nodeID insolar.RecordRef) {
	delete(ul.activeNodes, nodeID)
	delete(ul.proofs, nodeID)
	delete(ul.ghs, nodeID)
}

func (ul *unsyncList) AddNode(node insolar.NetworkNode, bitsetIndex uint16) {
	ul.addNode(node, int(bitsetIndex))
}

func (ul *unsyncList) AddProof(nodeID insolar.RecordRef, proof *consensus.NodePulseProof) {
	ul.proofs[nodeID] = proof
}

func (ul *unsyncList) GetProof(nodeID insolar.RecordRef) *consensus.NodePulseProof {
	return ul.proofs[nodeID]
}

func newUnsyncList(origin insolar.NetworkNode, activeNodesSorted []insolar.NetworkNode, length int) *unsyncList {
	result := &unsyncList{
		length:      length,
		origin:      origin,
		indexToRef:  make(map[int]insolar.RecordRef, len(activeNodesSorted)),
		refToIndex:  make(map[insolar.RecordRef]int, len(activeNodesSorted)),
		activeNodes: make(map[insolar.RecordRef]insolar.NetworkNode, len(activeNodesSorted)),
	}
	for i, node := range activeNodesSorted {
		result.addNode(node, i)
	}
	result.proofs = make(map[insolar.RecordRef]*consensus.NodePulseProof)
	result.ghs = make(map[insolar.RecordRef]consensus.GlobuleHashSignature)

	return result
}

func (ul *unsyncList) addNodes(nodes []insolar.NetworkNode) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].ID().Compare(nodes[j].ID()) < 0
	})

	for index, node := range nodes {
		ul.addNode(node, index)
	}
}

func (ul *unsyncList) addNode(node insolar.NetworkNode, index int) {
	ul.indexToRef[index] = node.ID()
	ul.refToIndex[node.ID()] = index
	ul.activeNodes[node.ID()] = node
}

func (ul *unsyncList) GetActiveNode(ref insolar.RecordRef) insolar.NetworkNode {
	return ul.activeNodes[ref]
}

func (ul *unsyncList) GetActiveNodes() []insolar.NetworkNode {
	return sortedNodeList(ul.activeNodes)
}

func sortedNodeList(nodes map[insolar.RecordRef]insolar.NetworkNode) []insolar.NetworkNode {
	result := make([]insolar.NetworkNode, len(nodes))
	i := 0
	for _, node := range nodes {
		result[i] = node
		i++
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID().Compare(result[j].ID()) < 0
	})
	return result
}

func (ul *unsyncList) IndexToRef(index int) (insolar.RecordRef, error) {
	if index < 0 || index >= ul.length {
		return insolar.RecordRef{}, consensus.ErrBitSetOutOfRange
	}
	result, ok := ul.indexToRef[index]
	if !ok {
		return insolar.RecordRef{}, consensus.ErrBitSetNodeIsMissing
	}
	return result, nil
}

func (ul *unsyncList) RefToIndex(nodeID insolar.RecordRef) (int, error) {
	index, ok := ul.refToIndex[nodeID]
	if !ok {
		return 0, consensus.ErrBitSetIncorrectNode
	}
	return index, nil
}

func (ul *unsyncList) Length() int {
	return ul.length
}

func ApplyClaims(ul network.UnsyncList, claims []consensus.ReferendumClaim) error {
	for _, claim := range claims {
		c, ok := claim.(*consensus.NodeAnnounceClaim)
		if !ok {
			continue
		}

		// TODO: fix version
		node, err := ClaimToNode("", &c.NodeJoinClaim)
		if err != nil {
			return errors.Wrap(err, "[ AddClaims ] failed to convert Claim -> NetworkNode")
		}
		// TODO: check these two
		ul.AddNode(node, c.NodeAnnouncerIndex)
		ul.AddNode(ul.GetOrigin(), c.NodeJoinerIndex)
	}
	return nil
}
