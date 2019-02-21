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

package bootstrap

import (
	"crypto/rand"
	"fmt"
	"math"
	"sort"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/utils"
	"github.com/pkg/errors"
)

const nonceSize int = 128

// CheckShortIDCollision returns true if NodeKeeper already contains node with such ShortID
func CheckShortIDCollision(keeper network.NodeKeeper, id core.ShortNodeID) bool {
	return keeper.GetActiveNodeByShortID(id) != nil
}

// GenerateShortID correct ShortID of the node so it does not conflict with existing active node list
func GenerateShortID(keeper network.NodeKeeper, nodeID core.RecordRef) core.ShortNodeID {
	shortID := utils.GenerateShortID(nodeID)
	if !CheckShortIDCollision(keeper, shortID) {
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
