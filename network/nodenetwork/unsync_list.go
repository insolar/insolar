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

func (ul *unsyncList) AddClaims(from core.RecordRef, claims []consensus.ReferendumClaim, addressMap map[core.RecordRef]string) {
	ul.addressMap = addressMap
	ul.claims[from] = claims
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
	case *consensus.NodeAnnounceClaim:
		// TODO: fix version
		node, err := claimToNode(ul.addressMap[t.NodeRef], "", t)
		if err != nil {
			log.Error("[ mergeClaim ] failed to get a Node")
		}
		addFunc(node)
	case *consensus.NodeJoinClaim:
		// TODO: fix version
		node, err := claimToNode(ul.addressMap[t.NodeRef], "", t)
		if err != nil {
			log.Error("[ mergeClaim ] failed to get a Node")
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

func claimToNode(address, version string, claim consensus.ReferendumClaim) (core.Node, error) {
	njc, ok := claim.(*consensus.NodeJoinClaim)
	if !ok {
		return nil, errors.New("[ ClaimToNode ] failed to convert a claim to node koin claim")
	}
	keyProc := platformpolicy.NewKeyProcessor()
	key, err := keyProc.ImportPublicKey(njc.NodePK[:])
	if err != nil {
		return nil, errors.Wrap(err, "[ Node ] failed to import a public key")
	}
	node := NewNode(
		njc.NodeRef,
		[]core.StaticRole{core.StaticRole(int(njc.NodeRoleRecID))},
		key,
		address,
		version)
	return node, nil
}

func nodeToClaim(node core.Node) (consensus.ReferendumClaim, error) {
	// TODO: do this if u know how
	return nil, nil
}
