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

package auth

import (
	"sort"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/nodenetwork"
)

// MajorityRuleCheck
func MajorityRuleCheck(activeNodesLists [][]core.Node, majorityRule int) (activeNodesList []core.Node, success bool) {
	// TODO: fair majorityRule check keeping in mind possible discovery redirects
	return activeNodesLists[0], true
}

// CheckShortIDCollision returns true if NodeKeeper already contains node with such ShortID
func CheckShortIDCollision(keeper network.NodeKeeper, id core.ShortNodeID) bool {
	return keeper.GetActiveNodeByShortID(id) != nil
}

// CorrectShortIDCollision correct ShortID of the node so it does not conflict with existing active node list
func CorrectShortIDCollision(keeper network.NodeKeeper, node core.Node) {
	activeNodes := keeper.GetActiveNodes()
	shortIDs := make([]core.ShortNodeID, len(activeNodes))
	for i, activeNode := range activeNodes {
		shortIDs[i] = activeNode.ShortID()
	}
	sort.Slice(shortIDs, func(i, j int) bool {
		return shortIDs[i] < shortIDs[j]
	})
	shortID := generateNonConflictingID(shortIDs, node.ShortID())
	mutable := node.(nodenetwork.MutableNode)
	mutable.SetShortID(shortID)
}

func generateNonConflictingID(sortedSlice []core.ShortNodeID, conflictingID core.ShortNodeID) core.ShortNodeID {
	index := sort.Search(len(sortedSlice), func(i int) bool {
		return sortedSlice[i] >= conflictingID
	})
	result := conflictingID
	for {
		index++
		result++
		if index >= len(sortedSlice) || result != sortedSlice[index] {
			return result
		}
	}
	// TODO: handle uint32 overflow
}
