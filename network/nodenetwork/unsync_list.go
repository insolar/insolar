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
	"sort"

	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
)

func copyMap(m map[core.RecordRef]core.Node) map[core.RecordRef]core.Node {
	result := make(map[core.RecordRef]core.Node, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

type unsyncList struct {
	origin      core.Node
	activeNodes map[core.RecordRef]core.Node
	claims      map[core.RecordRef][]consensus.ReferendumClaim
	refToIndex  map[core.RecordRef]int
	proofs      map[core.RecordRef]*consensus.NodePulseProof
	ghs         map[core.RecordRef]consensus.GlobuleHashSignature
	indexToRef  map[int]core.RecordRef
	cache       []byte
}

func (ul *unsyncList) GlobuleHashSignatures() map[core.RecordRef]consensus.GlobuleHashSignature {
	return ul.ghs
}

func (ul *unsyncList) ApproveSync(sync []core.RecordRef) {
	prevActive := make([]core.RecordRef, 0, len(ul.activeNodes))
	for nodeID := range ul.activeNodes {
		prevActive = append(prevActive, nodeID)
	}
	diff := diffList(prevActive, sync)
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

func newUnsyncList(origin core.Node, activeNodesSorted []core.Node) *unsyncList {
	result := &unsyncList{
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
	ul.cache = nil
}

func (ul *unsyncList) AddClaims(claims map[core.RecordRef][]consensus.ReferendumClaim) {
	ul.claims = claims
	ul.cache = nil
}

func (ul *unsyncList) CalculateHash(scheme core.PlatformCryptographyScheme) ([]byte, error) {
	if ul.cache != nil {
		return ul.cache, nil
	}
	m := ul.getMergedNodeMap()
	sorted := sortedNodeList(m)
	var err error
	ul.cache, err = CalculateHash(scheme, sorted)
	return ul.cache, err
}

func (ul *unsyncList) GetActiveNode(ref core.RecordRef) core.Node {
	return ul.activeNodes[ref]
}

func (ul *unsyncList) GetActiveNodes() []core.Node {
	return sortedNodeList(ul.activeNodes)
}

func (ul *unsyncList) getMergedNodeMap() map[core.RecordRef]core.Node {
	nodes := copyMap(ul.activeNodes)

	for _, claimList := range ul.claims {
		for _, claim := range claimList {
			ul.mergeClaim(nodes, claim)
		}
	}

	return nodes
}

func (ul *unsyncList) mergeClaim(nodes map[core.RecordRef]core.Node, claim consensus.ReferendumClaim) {
	switch t := claim.(type) {
	case *consensus.NodeJoinClaim:
		// TODO: fix version
		node, err := ClaimToNode("", t)
		if err != nil {
			log.Error("[ mergeClaim ] failed to convert Claim -> Node")
		}
		nodes[node.ID()] = node
	case *consensus.NodeLeaveClaim:
		// TODO: add node ID to node leave claim (only to struct, not packet)
		// delete(nodes, ref)
		break
	}
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
	if index < 0 || index >= len(ul.indexToRef) {
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
	return len(ul.activeNodes)
}

type sparseUnsyncList struct {
	unsyncList
	capacity int
}

func newSparseUnsyncList(origin core.Node, capacity int) *sparseUnsyncList {
	return &sparseUnsyncList{unsyncList: *newUnsyncList(origin, nil), capacity: capacity}
}

func (ul *sparseUnsyncList) Length() int {
	return ul.capacity
}

func (ul *sparseUnsyncList) AddClaims(claims map[core.RecordRef][]consensus.ReferendumClaim) {
	ul.unsyncList.AddClaims(claims)

	for _, claimList := range claims {
		for _, claim := range claimList {
			t := claim.Type()
			if t != consensus.TypeNodeAnnounceClaim {
				continue
			}
			c, ok := claim.(*consensus.NodeAnnounceClaim)
			if !ok {
				log.Error("[ AddClaims ] Could not convert claim with type TypeNodeAnnounceClaim to NodeAnnounceClaim")
			}

			// TODO: fix version
			node, err := ClaimToNode("", &c.NodeJoinClaim)
			if err != nil {
				log.Error("[ AddClaims ] failed to convert Claim -> Node")
			}
			// TODO: check these two
			ul.addNode(node, int(c.NodeAnnouncerIndex))
			ul.addNode(ul.origin, int(c.NodeJoinerIndex))
		}
	}
}
