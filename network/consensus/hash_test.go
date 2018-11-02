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

package consensus

import (
	"encoding/hex"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

const (
	nullHash = "6b4e03423667dbb73b6e15454f0eb1abd4597f9a1b078e3f5b5a6bc7"
)

func TestNodeConsensus_calculateNodeHash(t *testing.T) {
	nodeID := testutils.RandomRef()
	hash1, _ := CalculateNodeUnsyncHash(nodeID, nil)
	hash2, _ := CalculateNodeUnsyncHash(nodeID, []core.Node{})

	assert.Equal(t, nullHash, hex.EncodeToString(hash1.Hash))
	assert.Equal(t, hash1, hash2)
	assert.Equal(t, nodeID, hash1.NodeID)

	activeNode1 := newActiveNode(0, 0)
	activeNode2 := newActiveNode(0, 0)

	activeNode1Slice := []core.Node{activeNode1}
	activeNode2Slice := []core.Node{activeNode2}

	hash1, _ = CalculateNodeUnsyncHash(nodeID, activeNode1Slice)
	hash2, _ = CalculateNodeUnsyncHash(nodeID, activeNode2Slice)
	assert.Equal(t, hash1, hash2)

	activeNode3 := newActiveNode(1, 0)
	activeNode3Slice := []core.Node{activeNode3}
	hash3, _ := CalculateNodeUnsyncHash(nodeID, activeNode3Slice)
	assert.NotEqual(t, hash1, hash3)

	// nodes order in slice should not affect hash calculating
	slice1 := []core.Node{activeNode1, activeNode2}
	slice2 := []core.Node{activeNode2, activeNode1}
	hash1, _ = CalculateNodeUnsyncHash(nodeID, slice1)
	hash2, _ = CalculateNodeUnsyncHash(nodeID, slice2)
	assert.Equal(t, hash1, hash2)
	assert.Equal(t, nodeID, hash1.NodeID)
}
