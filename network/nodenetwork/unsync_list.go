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
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

func copyMap(m map[core.RecordRef]core.Node) map[core.RecordRef]core.Node {
	result := make(map[core.RecordRef]core.Node, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

type unsyncList struct {
	activeNodes map[core.RecordRef]core.Node
	addressMap  map[core.RecordRef]string
	claims      map[core.RecordRef][]consensus.ReferendumClaim
	refToIndex  map[core.RecordRef]int
	indexToRef  map[int]core.RecordRef
	cache       []byte
}

func newUnsyncList(activeNodesSorted []core.Node) *unsyncList {
	indexToRef := make(map[int]core.RecordRef, len(activeNodesSorted))
	refToIndex := make(map[core.RecordRef]int, len(activeNodesSorted))
	activeNodes := make(map[core.RecordRef]core.Node, len(activeNodesSorted))
	for i, node := range activeNodesSorted {
		indexToRef[i] = node.ID()
		refToIndex[node.ID()] = i
		activeNodes[node.ID()] = node
	}
	claims := make(map[core.RecordRef][]consensus.ReferendumClaim)

	return &unsyncList{activeNodes: activeNodes, claims: claims, refToIndex: refToIndex, indexToRef: indexToRef}
}

func (ul *unsyncList) RemoveClaims(from core.RecordRef) {
	delete(ul.claims, from)
	ul.cache = nil
}

func (ul *unsyncList) AddClaims(claims map[core.RecordRef][]consensus.ReferendumClaim, addressMap map[core.RecordRef]string) {
	ul.addressMap = addressMap
	ul.claims = claims
	ul.cache = nil
}

func (ul *unsyncList) CalculateHash() ([]byte, error) {
	if ul.cache != nil {
		return ul.cache, nil
	}
	m := copyMap(ul.activeNodes)
	ul.merge(m, ul.claims)
	sorted := sortedNodeList(m)
	var err error
	ul.cache, err = CalculateHash(nil, sorted)
	return ul.cache, err
}

func (ul *unsyncList) GetActiveNode(ref core.RecordRef) core.Node {
	return ul.activeNodes[ref]
}

func (ul *unsyncList) GetActiveNodes() []core.Node {
	return sortedNodeList(ul.activeNodes)
}

type adder func(core.Node)
type deleter func(core.RecordRef)

func (ul *unsyncList) merge(nodes map[core.RecordRef]core.Node, claims map[core.RecordRef][]consensus.ReferendumClaim) {
	addNode := func(node core.Node) {
		nodes[node.ID()] = node
	}
	delNode := func(ref core.RecordRef) {
		delete(nodes, ref)
	}
	ul.mergeWith(claims, addNode, delNode)
}

func (ul *unsyncList) mergeWith(claims map[core.RecordRef][]consensus.ReferendumClaim, addFunc adder, delFunc deleter) {
	for _, claimList := range claims {
		for _, claim := range claimList {
			ul.mergeClaim(claim, addFunc, delFunc)
		}
	}
}

func (ul *unsyncList) mergeClaim(claim consensus.ReferendumClaim, addFunc adder, delFunc deleter) {
	switch t := claim.(type) {
	case *consensus.NodeJoinClaim:
		// TODO: fix version
		node, err := claimToNode(ul.addressMap[t.NodeRef], "", t)
		if err != nil {
			log.Error("[ mergeClaim ] failed to convert Claim -> Node")
		}
		addFunc(node)
	case *consensus.NodeLeaveClaim:
		// TODO: add node ID to node leave claim (only to struct, not packet)
		// delFunc()
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

func newSparseUnsyncList(capacity int) *sparseUnsyncList {
	return &sparseUnsyncList{unsyncList: *newUnsyncList(nil), capacity: capacity}
}

func (ul *sparseUnsyncList) Length() int {
	return ul.capacity
}

func (ul *sparseUnsyncList) AddClaims(claims map[core.RecordRef][]consensus.ReferendumClaim, addressMap map[core.RecordRef]string) {
	ul.unsyncList.AddClaims(claims, addressMap)

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
			node, err := claimToNode(ul.addressMap[c.NodeRef], "", &c.NodeJoinClaim)
			if err != nil {
				log.Error("[ AddClaims ] failed to convert Claim -> Node")
			}
			ul.activeNodes[node.ID()] = node
		}
	}
}

func claimToNode(address, version string, claim *consensus.NodeJoinClaim) (core.Node, error) {
	keyProc := platformpolicy.NewKeyProcessor()
	key, err := keyProc.ImportPublicKeyPEM(claim.NodePK[:])
	if err != nil {
		return nil, errors.Wrap(err, "[ ClaimToNode ] failed to import a public key")
	}
	node := NewNode(
		claim.NodeRef,
		core.StaticRole(int(claim.NodeRoleRecID)),
		key,
		address,
		version)
	return node, nil
}
