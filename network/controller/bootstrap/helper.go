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

package bootstrap

import (
	"math"
	"sort"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/utils"
)

// CheckShortIDCollision returns true if NodeKeeper already contains node with such ShortID
func CheckShortIDCollision(keeper network.NodeKeeper, id insolar.ShortNodeID) bool {
	return keeper.GetAccessor().GetActiveNodeByShortID(id) != nil
}

// GenerateShortID correct ShortID of the node so it does not conflict with existing active node list
func GenerateShortID(keeper network.NodeKeeper, nodeID insolar.Reference) insolar.ShortNodeID {
	shortID := utils.GenerateShortID(nodeID)
	if !CheckShortIDCollision(keeper, shortID) {
		return shortID
	}
	return regenerateShortID(keeper, shortID)
}

func regenerateShortID(keeper network.NodeKeeper, shortID insolar.ShortNodeID) insolar.ShortNodeID {
	activeNodes := keeper.GetAccessor().GetActiveNodes()
	shortIDs := make([]insolar.ShortNodeID, len(activeNodes))
	for i, activeNode := range activeNodes {
		shortIDs[i] = activeNode.ShortID()
	}
	sort.Slice(shortIDs, func(i, j int) bool {
		return shortIDs[i] < shortIDs[j]
	})
	return generateNonConflictingID(shortIDs, shortID)
}

func generateNonConflictingID(sortedSlice []insolar.ShortNodeID, conflictingID insolar.ShortNodeID) insolar.ShortNodeID {
	index := sort.Search(len(sortedSlice), func(i int) bool {
		return sortedSlice[i] >= conflictingID
	})
	result := conflictingID
	repeated := false
	for {
		if result == math.MaxUint32 {
			if !repeated {
				repeated = true
				result = 0
				index = 0
			} else {
				panic("[ generateNonConflictingID ] shortID overflow twice")
			}
		}
		index++
		result++
		if index >= len(sortedSlice) || result != sortedSlice[index] {
			return result
		}
	}
}

func RemoveOrigin(discoveryNodes []insolar.DiscoveryNode, origin insolar.Reference) []insolar.DiscoveryNode {
	for i, discoveryNode := range discoveryNodes {
		if origin.Equal(*discoveryNode.GetNodeRef()) {
			return append(discoveryNodes[:i], discoveryNodes[i+1:]...)
		}
	}
	return discoveryNodes
}

func FindDiscovery(cert insolar.Certificate, ref insolar.Reference) insolar.DiscoveryNode {
	bNodes := cert.GetDiscoveryNodes()
	for _, discoveryNode := range bNodes {
		if ref.Equal(*discoveryNode.GetNodeRef()) {
			return discoveryNode
		}
	}
	return nil
}
