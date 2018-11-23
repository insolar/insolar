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
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

const (
	nullHash = "a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26"
)

func TestNodeConsensus_calculateNodeHash(t *testing.T) {
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	nodeID := testutils.RandomRef()
	hash1, _ := CalculateNodeUnsyncHash(scheme, nodeID, nil)
	hash2, _ := CalculateNodeUnsyncHash(scheme, nodeID, []core.Node{})

	require.Equal(t, nullHash, hex.EncodeToString(hash1.Hash))
	require.Equal(t, hash1, hash2)
	require.Equal(t, nodeID, hash1.NodeID)

	activeNode1 := newActiveNode(0, 0)
	activeNode2 := newActiveNode(0, 0)

	activeNode1Slice := []core.Node{activeNode1}
	activeNode2Slice := []core.Node{activeNode2}

	hash1, _ = CalculateNodeUnsyncHash(scheme, nodeID, activeNode1Slice)
	hash2, _ = CalculateNodeUnsyncHash(scheme, nodeID, activeNode2Slice)
	require.Equal(t, hash1, hash2)

	activeNode3 := newActiveNode(1, 0)
	activeNode3Slice := []core.Node{activeNode3}
	hash3, _ := CalculateNodeUnsyncHash(scheme, nodeID, activeNode3Slice)
	require.NotEqual(t, hash1, hash3)

	// nodes order in slice should not affect hash calculating
	slice1 := []core.Node{activeNode1, activeNode2}
	slice2 := []core.Node{activeNode2, activeNode1}
	hash1, _ = CalculateNodeUnsyncHash(scheme, nodeID, slice1)
	hash2, _ = CalculateNodeUnsyncHash(scheme, nodeID, slice2)
	require.Equal(t, hash1, hash2)
	require.Equal(t, nodeID, hash1.NodeID)
}
