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

package bootstrap

import (
	"crypto/rand"
	"fmt"
	"sort"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/utils"
	"github.com/pkg/errors"
)

const nonceSize int = 128

// checkShortIDCollision returns true if NodeKeeper already contains node with such ShortID
func checkShortIDCollision(keeper network.NodeKeeper, id core.ShortNodeID) bool {
	return keeper.GetActiveNodeByShortID(id) != nil
}

// GenerateShortID correct ShortID of the node so it does not conflict with existing active node list
func GenerateShortID(keeper network.NodeKeeper, nodeID core.RecordRef) core.ShortNodeID {
	shortID := utils.GenerateShortID(nodeID)
	if !checkShortIDCollision(keeper, shortID) {
		return shortID
	}
	return regenerateShortID(keeper, shortID)
}

func regenerateShortID(keeper network.NodeKeeper, shortID core.ShortNodeID) core.ShortNodeID {
	activeNodes := keeper.GetActiveNodes()
	shortIDs := make([]core.ShortNodeID, len(activeNodes))
	for i, activeNode := range activeNodes {
		shortIDs[i] = activeNode.ShortID()
	}
	sort.Slice(shortIDs, func(i, j int) bool {
		return shortIDs[i] < shortIDs[j]
	})
	return generateNonConflictingID(shortIDs, shortID)
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

func RemoveOrigin(discoveryNodes []core.DiscoveryNode, origin core.RecordRef) ([]core.DiscoveryNode, error) {
	for i, discoveryNode := range discoveryNodes {
		if origin.Equal(*discoveryNode.GetNodeRef()) {
			return append(discoveryNodes[:i], discoveryNodes[i+1:]...), nil
		}
	}
	return nil, errors.New("Origin not found in discovery nodes list")
}

func FindDiscovery(cert core.Certificate, ref core.RecordRef) core.DiscoveryNode {
	bNodes := cert.GetDiscoveryNodes()
	for _, discoveryNode := range bNodes {
		if ref.Equal(*discoveryNode.GetNodeRef()) {
			return discoveryNode
		}
	}
	return nil
}

func Xor(first, second []byte) []byte {
	if len(second) < len(first) {
		temp := second
		second = first
		first = temp
	}
	result := make([]byte, len(second))
	for i, d := range second {
		result[i] = first[i%len(first)] ^ d
	}
	return result
}

func GenerateNonce() (Nonce, error) {
	buffer := [nonceSize]byte{}
	l, err := rand.Read(buffer[:])
	if err != nil {
		return nil, errors.Wrapf(err, "error generating nonce")
	}
	if l != nonceSize {
		return nil, errors.New(fmt.Sprintf("GenerateNonce: generated size %d does equal to required size %d", l, nonceSize))
	}
	return buffer[:], nil
}
