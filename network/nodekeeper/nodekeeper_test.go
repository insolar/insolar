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

package nodekeeper

import (
	"encoding/hex"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
)

func TestNodekeeper_calculateNodeHash(t *testing.T) {
	hash1, _ := calculateHash(nil)
	hash2, _ := calculateHash([]*core.ActiveNode{})

	assert.Equal(t, "6b4e03423667dbb73b6e15454f0eb1abd4597f9a1b078e3f5b5a6bc7", hex.EncodeToString(hash1))
	assert.Equal(t, hash1, hash2)

	var mask1 core.JetRoleMask
	mask1.Set(core.RoleVirtualExecutor)

	activeNode1 := core.ActiveNode{
		NodeID:    core.RecordRef{0},
		PulseNum:  core.PulseNumber(0),
		State:     core.NodeActive,
		JetRoles:  mask1,
		PublicKey: []byte{0, 0, 0},
	}
	activeNode2 := core.ActiveNode{
		NodeID:    core.RecordRef{0},
		PulseNum:  core.PulseNumber(0),
		State:     core.NodeActive,
		JetRoles:  mask1,
		PublicKey: []byte{0, 0, 0},
	}

	activeNode1Slice := []*core.ActiveNode{&activeNode1}
	activeNode2Slice := []*core.ActiveNode{&activeNode2}

	hash1, _ = calculateHash(activeNode1Slice)
	hash2, _ = calculateHash(activeNode2Slice)
	assert.Equal(t, hash1, hash2)
	activeNode2.NodeID = core.RecordRef{1}
	hash2, _ = calculateHash(activeNode2Slice)
	assert.NotEqual(t, hash1, hash2)

	// nodes order in slice should not affect hash calculating
	slice1 := []*core.ActiveNode{&activeNode1, &activeNode2}
	slice2 := []*core.ActiveNode{&activeNode2, &activeNode1}
	hash1, _ = calculateHash(slice1)
	hash2, _ = calculateHash(slice2)
	assert.Equal(t, hash1, hash2)
}
