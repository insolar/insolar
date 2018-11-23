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
	claims      map[core.RecordRef][]consensus.ReferendumClaim
	refToIndex  map[core.RecordRef]int
	indexToRef  []core.RecordRef
	cache       []byte
}

func newUnsyncList(activeNodesSorted []core.Node) *unsyncList {
	indexToRef := make([]core.RecordRef, len(activeNodesSorted))
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

func (ul *unsyncList) AddClaims(from core.RecordRef, claims []consensus.ReferendumClaim) {
	ul.claims[from] = claims
	ul.cache = nil
}

func (ul *unsyncList) CalculateHash() ([]byte, error) {
	if ul.cache != nil {
		return ul.cache, nil
	}
	m := copyMap(ul.activeNodes)
	merge(m, ul.claims)
	sorted := sortedNodeList(m)
	var err error
	ul.cache, err = CalculateHash(nil, sorted)
	return ul.cache, err
}

func merge(nodes map[core.RecordRef]core.Node, claims map[core.RecordRef][]consensus.ReferendumClaim) {
	// TODO: implement
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
	return ul.indexToRef[index], nil
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
