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
 *
 */

package nodenetwork

import (
	"sort"

	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/pkg/errors"
)

func copyActiveNodes(m map[core.RecordRef]core.Node) map[core.RecordRef]core.Node {
	result := make(map[core.RecordRef]core.Node, len(m))
	for k, v := range m {
		v.(MutableNode).ChangeState()
		result[k] = v
	}
	return result
}

type unsyncList struct {
	length      int
	origin      core.Node
	activeNodes map[core.RecordRef]core.Node
	claims      map[core.RecordRef][]consensus.ReferendumClaim
	refToIndex  map[core.RecordRef]int
	proofs      map[core.RecordRef]*consensus.NodePulseProof
	ghs         map[core.RecordRef]consensus.GlobuleHashSignature
	indexToRef  map[int]core.RecordRef
}

func (ul *unsyncList) GetGlobuleHashSignature(ref core.RecordRef) (consensus.GlobuleHashSignature, bool) {
	ghs, ok := ul.ghs[ref]
	return ghs, ok
}

func (ul *unsyncList) SetGlobuleHashSignature(ref core.RecordRef, ghs consensus.GlobuleHashSignature) {
	ul.ghs[ref] = ghs
}

func (ul *unsyncList) RemoveNode(nodeID core.RecordRef) {
	delete(ul.activeNodes, nodeID)
	delete(ul.claims, nodeID)
	delete(ul.proofs, nodeID)
	delete(ul.ghs, nodeID)
}

func (ul *unsyncList) ApproveSync(sync []core.RecordRef) {
	prevActive := make([]core.RecordRef, 0, len(ul.activeNodes))
	for nodeID := range ul.activeNodes {
		prevActive = append(prevActive, nodeID)
	}
	diff := removeFromList(prevActive, sync)
	for _, node := range diff {
		ul.removeNode(node)
	}
}

func (ul *unsyncList) AddNode(node core.Node, bitsetIndex uint16) {
	ul.addNode(node, int(bitsetIndex))
}

func (ul *unsyncList) GetClaims(nodeID core.RecordRef) []consensus.ReferendumClaim {
	return ul.claims[nodeID]
}

func (ul *unsyncList) AddProof(nodeID core.RecordRef, proof *consensus.NodePulseProof) {
	ul.proofs[nodeID] = proof
}

func (ul *unsyncList) GetProof(nodeID core.RecordRef) *consensus.NodePulseProof {
	return ul.proofs[nodeID]
}

func (ul *unsyncList) InsertClaims(ref core.RecordRef, claims []consensus.ReferendumClaim) {
	ul.claims[ref] = claims
}

func newUnsyncList(origin core.Node, activeNodesSorted []core.Node, length int) *unsyncList {
	result := &unsyncList{
		length:      length,
		origin:      origin,
		indexToRef:  make(map[int]core.RecordRef, len(activeNodesSorted)),
		refToIndex:  make(map[core.RecordRef]int, len(activeNodesSorted)),
		activeNodes: make(map[core.RecordRef]core.Node, len(activeNodesSorted)),
	}
	for i, node := range activeNodesSorted {
		result.addNode(node, i)
	}
	result.proofs = make(map[core.RecordRef]*consensus.NodePulseProof)
	result.claims = make(map[core.RecordRef][]consensus.ReferendumClaim)
	result.ghs = make(map[core.RecordRef]consensus.GlobuleHashSignature)

	return result
}

func (ul *unsyncList) addNodes(nodes []core.Node) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].ID().Compare(nodes[j].ID()) < 0
	})

	for index, node := range nodes {
		ul.addNode(node, index)
	}
}

func (ul *unsyncList) addNode(node core.Node, index int) {
	ul.indexToRef[index] = node.ID()
	ul.refToIndex[node.ID()] = index
	ul.activeNodes[node.ID()] = node
}

func (ul *unsyncList) removeNode(nodeID core.RecordRef) {
	delete(ul.activeNodes, nodeID)
	delete(ul.claims, nodeID)
	delete(ul.proofs, nodeID)
	delete(ul.ghs, nodeID)
	i, ok := ul.refToIndex[nodeID]
	if ok {
		delete(ul.indexToRef, i)
	}
	delete(ul.refToIndex, nodeID)
}

func (ul *unsyncList) AddClaims(claims map[core.RecordRef][]consensus.ReferendumClaim) error {
	ul.claims = claims
	return nil
}

func (ul *unsyncList) GetActiveNode(ref core.RecordRef) core.Node {
	return ul.activeNodes[ref]
}

func (ul *unsyncList) GetActiveNodes() []core.Node {
	return sortedNodeList(ul.activeNodes)
}

func (ul *unsyncList) GetMergedCopy() (*network.MergedListCopy, error) {
	nodes := copyActiveNodes(ul.activeNodes)

	var nodesJoinedDuringPrevPulse bool
	for _, claimList := range ul.claims {
		for _, claim := range claimList {
			isJoin, err := mergeClaim(nodes, claim)
			if err != nil {
				return nil, errors.Wrap(err, "[ GetMergedCopy ] failed to merge a claim")
			}

			nodesJoinedDuringPrevPulse = nodesJoinedDuringPrevPulse || isJoin
		}
	}

	return &network.MergedListCopy{
		ActiveList:                 nodes,
		NodesJoinedDuringPrevPulse: nodesJoinedDuringPrevPulse,
	}, nil
}

func mergeClaim(nodes map[core.RecordRef]core.Node, claim consensus.ReferendumClaim) (bool, error) {
	isJoinClaim := false
	switch t := claim.(type) {
	case *consensus.NodeJoinClaim:
		isJoinClaim = true
		// TODO: fix version
		node, err := ClaimToNode("", t)
		if err != nil {
			return isJoinClaim, errors.Wrap(err, "[ mergeClaim ] failed to convert Claim -> Node")
		}
		node.(MutableNode).SetState(core.NodeJoining)
		nodes[node.ID()] = node
	case *consensus.NodeLeaveClaim:
		if nodes[t.NodeID] == nil {
			break
		}

		node := nodes[t.NodeID].(MutableNode)
		if t.ETA == 0 || !node.Leaving() {
			node.SetLeavingETA(t.ETA)
		}
	}

	return isJoinClaim, nil
}

func sortedNodeList(nodes map[core.RecordRef]core.Node) []core.Node {
	result := make([]core.Node, len(nodes))
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

func (ul *unsyncList) IndexToRef(index int) (core.RecordRef, error) {
	if index < 0 || index >= ul.length {
		return core.RecordRef{}, consensus.ErrBitSetOutOfRange
	}
	result, ok := ul.indexToRef[index]
	if !ok {
		return core.RecordRef{}, consensus.ErrBitSetNodeIsMissing
	}
	return result, nil
}

func (ul *unsyncList) RefToIndex(nodeID core.RecordRef) (int, error) {
	index, ok := ul.refToIndex[nodeID]
	if !ok {
		return 0, consensus.ErrBitSetIncorrectNode
	}
	return index, nil
}

func (ul *unsyncList) Length() int {
	return ul.length
}

type sparseUnsyncList struct {
	unsyncList
}

func newSparseUnsyncList(origin core.Node, capacity int) *sparseUnsyncList {
	return &sparseUnsyncList{unsyncList: *newUnsyncList(origin, nil, capacity)}
}

func (ul *sparseUnsyncList) AddClaims(claims map[core.RecordRef][]consensus.ReferendumClaim) error {
	err := ul.unsyncList.AddClaims(claims)
	if err != nil {
		return errors.Wrap(err, "[ AddClaims ] failed to add a claims")
	}

	for _, claimList := range claims {
		for _, claim := range claimList {
			c, ok := claim.(*consensus.NodeAnnounceClaim)
			if !ok {
				log.Error("[ AddClaims ] Could not convert claim with type TypeNodeAnnounceClaim to NodeAnnounceClaim")
				continue
			}

			// TODO: fix version
			node, err := ClaimToNode("", &c.NodeJoinClaim)
			if err != nil {
				return errors.Wrap(err, "[ AddClaims ] failed to convert Claim -> Node")
			}
			// TODO: check these two
			ul.addNode(node, int(c.NodeAnnouncerIndex))
			ul.addNode(ul.origin, int(c.NodeJoinerIndex))
		}
	}
	return nil
}

func (ul *sparseUnsyncList) UpdateClaims(ref core.RecordRef, claims []consensus.ReferendumClaim) {
	ul.unsyncList.claims[ref] = claims
}
