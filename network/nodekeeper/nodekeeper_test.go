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

const (
	nullHash = "6b4e03423667dbb73b6e15454f0eb1abd4597f9a1b078e3f5b5a6bc7"
)

func newActiveNode(ref byte) *core.ActiveNode {
	var mask core.JetRoleMask
	mask.Set(core.RoleVirtualExecutor)

	return &core.ActiveNode{
		NodeID:    core.RecordRef{ref},
		PulseNum:  core.PulseNumber(0),
		State:     core.NodeActive,
		JetRoles:  mask,
		PublicKey: []byte{0, 0, 0},
	}
}

func newSelfNode(ref core.RecordRef) *core.ActiveNode {
	var mask core.JetRoleMask
	mask.Set(core.RoleVirtualExecutor)

	return &core.ActiveNode{
		NodeID:    ref,
		PulseNum:  core.PulseNumber(0),
		State:     core.NodeActive,
		JetRoles:  mask,
		PublicKey: []byte{0, 0, 0},
	}
}

func newNodeKeeper() NodeKeeper {
	id := core.RecordRef{255}
	keeper := NewNodeKeeper(id)
	keeper.AddActiveNodes([]*core.ActiveNode{newSelfNode(id)})
	return keeper
}

func TestNodekeeper_calculateNodeHash(t *testing.T) {
	hash1, _ := CalculateHash(nil)
	hash2, _ := CalculateHash([]*core.ActiveNode{})

	assert.Equal(t, nullHash, hex.EncodeToString(hash1))
	assert.Equal(t, hash1, hash2)

	activeNode1 := newActiveNode(0)
	activeNode2 := newActiveNode(0)

	activeNode1Slice := []*core.ActiveNode{activeNode1}
	activeNode2Slice := []*core.ActiveNode{activeNode2}

	hash1, _ = CalculateHash(activeNode1Slice)
	hash2, _ = CalculateHash(activeNode2Slice)
	assert.Equal(t, hash1, hash2)
	activeNode2.NodeID = core.RecordRef{1}
	hash2, _ = CalculateHash(activeNode2Slice)
	assert.NotEqual(t, hash1, hash2)

	// nodes order in slice should not affect hash calculating
	slice1 := []*core.ActiveNode{activeNode1, activeNode2}
	slice2 := []*core.ActiveNode{activeNode2, activeNode1}
	hash1, _ = CalculateHash(slice1)
	hash2, _ = CalculateHash(slice2)
	assert.Equal(t, hash1, hash2)
}

func TestNodekeeper_AddUnsync(t *testing.T) {
	id := core.RecordRef{}
	keeper := NewNodeKeeper(id)
	// AddUnsync should return error if we are not an active node
	err := keeper.AddUnsync(newActiveNode(0))
	assert.Error(t, err)
	// Add active node with NodeKeeper id, so we are now active and can add unsyncs
	keeper.AddActiveNodes([]*core.ActiveNode{newSelfNode(id)})
	err = keeper.AddUnsync(newActiveNode(0))
	assert.NoError(t, err)
	success, list := keeper.SetPulse(core.PulseNumber(0))
	assert.True(t, success)
	assert.Equal(t, 1, len(list))
}

func TestNodekeeper_AddUnsync2(t *testing.T) {
	keeper := newNodeKeeper()
	success, list := keeper.SetPulse(core.PulseNumber(0))
	err := keeper.AddUnsync(newActiveNode(0))
	assert.NoError(t, err)
	assert.True(t, success)
	assert.Equal(t, 0, len(list))
}

func TestNodekeeper_AddUnsync3(t *testing.T) {
	keeper := newNodeKeeper()
	err := keeper.AddUnsync(newActiveNode(0))
	success, list := keeper.SetPulse(core.PulseNumber(0))
	err = keeper.AddUnsync(newActiveNode(1))
	assert.NoError(t, err)
	assert.True(t, success)
	assert.Equal(t, 1, len(list))
}

func TestNodekeeper_pipeline(t *testing.T) {
	keeper := newNodeKeeper()
	for i := 0; i < 4; i++ {
		err := keeper.AddUnsync(newActiveNode(byte(2 * i)))
		assert.NoError(t, err)
		success, list := keeper.SetPulse(core.PulseNumber(i))
		assert.True(t, success)
		err = keeper.AddUnsync(newActiveNode(byte(2*i + 1)))
		assert.NoError(t, err)
		keeper.Sync(list, keeper.GetPulse())
	}
	// 3 nodes should not advance to join active list
	// 5 nodes should advance + 1 self node
	assert.Equal(t, 6, len(keeper.GetActiveNodes()))
	for i := 0; i < 5; i++ {
		assert.NotNil(t, keeper.GetActiveNode(core.RecordRef{byte(i)}))
	}
}

func TestNodekeeper_doubleSync(t *testing.T) {
	keeper := newNodeKeeper()
	err := keeper.AddUnsync(newActiveNode(0))
	assert.NoError(t, err)
	pulse := core.PulseNumber(0)
	success, list := keeper.SetPulse(pulse)
	assert.True(t, success)
	assert.Equal(t, 1, len(list))
	keeper.Sync(list, pulse)
	// second sync should be ignored because pulse has not changed
	keeper.Sync(list, pulse)
	// and added unsync node should not advance to active list (only one self node would be in the list)
	assert.Equal(t, 1, len(keeper.GetActiveNodes()))
	assert.Equal(t, keeper.GetSelf().NodeID, keeper.GetActiveNodes()[0].NodeID)
}

func TestNodekeeper_doubleSetPulse(t *testing.T) {
	keeper := newNodeKeeper()
	err := keeper.AddUnsync(newActiveNode(0))
	assert.NoError(t, err)
	pulse := core.PulseNumber(0)
	_, list := keeper.SetPulse(pulse)
	keeper.Sync(list, pulse)
	_, _ = keeper.SetPulse(core.PulseNumber(1))
	_, _ = keeper.SetPulse(core.PulseNumber(2))
	// node with ref 0 advanced to active list
	assert.Equal(t, 2, len(keeper.GetActiveNodes()))
	assert.NotNil(t, keeper.GetActiveNode(core.RecordRef{0}))
}

// func TestNodekeeper_AddUnsync3(t *testing.T) {
// 	_ = keeper.AddUnsync(newActiveNode(0, 0))
// 	_ = keeper.AddUnsync(newActiveNode(1, 0))
// 	gossip := []*core.ActiveNode{newActiveNode(2, 0), newActiveNode(3, 0)}
// 	_ = keeper.AddUnsyncGossip(gossip)
// 	assert.Equal(t, 2, len(keeper.GetUnsync()))
// 	keeper.Sync(true, core.PulseNumber(0))
// 	assert.Equal(t, 0, len(keeper.GetUnsync()))
// 	_ = keeper.SetPulse(core.PulseNumber(1))
// 	keeper.Sync(true, core.PulseNumber(1))
// 	assert.Equal(t, 4, len(keeper.GetActiveNodes()))
// 	for i := 0; i < 4; i++ {
// 		assert.NotNil(t, keeper.GetActiveNode(core.RecordRef{byte(i)}))
// 	}
// }
//
// func TestNodekeeper_AddUnsync_checks(t *testing.T) {
// 	keeper := NewNodeKeeper(core.RecordRef{}, time.Hour)
// 	_ = keeper.SetPulse(core.PulseNumber(0))
//
// 	// Unsync node pulse number should be equal to the NodeKeeper pulse number
// 	err := keeper.AddUnsync(newActiveNode(0, 1))
// 	assert.Error(t, err)
// 	err = keeper.AddUnsync(newActiveNode(0, 0))
// 	assert.NoError(t, err)
//
// 	// Gossip unsync node should not have reference id equal to one of the local unsync nodes
// 	err = keeper.AddUnsyncGossip([]*core.ActiveNode{newActiveNode(0, 0)})
// 	assert.Error(t, err)
// 	// Gossip unsync node pulse number should be equal to the NodeKeeper pulse number
// 	err = keeper.AddUnsyncGossip([]*core.ActiveNode{newActiveNode(1, 1)})
// 	assert.Error(t, err)
// 	err = keeper.AddUnsyncGossip([]*core.ActiveNode{newActiveNode(1, 0)})
// 	assert.NoError(t, err)
// }
//
// func TestNodekeeper_AddActiveNodes(t *testing.T) {
// 	keeper := NewNodeKeeper(core.RecordRef{}, time.Hour)
// 	_ = keeper.SetPulse(core.PulseNumber(0))
//
// 	node2 := newActiveNode(0, 0)
// 	node1 := newActiveNode(1, 0)
// 	nodes := []*core.ActiveNode{node1, node2}
// 	keeper.AddActiveNodes(nodes)
//
// 	assert.Equal(t, 2, len(keeper.GetActiveNodes()))
// 	assert.NotNil(t, keeper.GetActiveNode(core.RecordRef{0}))
// 	assert.NotNil(t, keeper.GetActiveNode(core.RecordRef{1}))
// }
//
// func TestNodekeeper_transitions1(t *testing.T) {
// 	keeper := NewNodeKeeper(core.RecordRef{}, time.Hour)
// 	_ = keeper.SetPulse(core.PulseNumber(0))
//
// 	keeper.AddUnsync(newActiveNode(0, 0))
// 	keeper.Sync(true, core.PulseNumber(0))
// 	// check that Sync is not called and the transition sync -> active is not performed
// 	_ = keeper.SetPulse(core.PulseNumber(0))
// 	_ = keeper.SetPulse(core.PulseNumber(0))
// 	assert.Equal(t, 0, len(keeper.GetActiveNodes()))
// }
//
// func TestNodekeeper_transitions2(t *testing.T) {
// 	keeper := NewNodeKeeper(core.RecordRef{}, time.Hour)
// 	_ = keeper.SetPulse(core.PulseNumber(0))
//
// 	keeper.AddUnsync(newActiveNode(0, 0))
// 	// check that Sync is called correctly every time and the transition unsync -> sync -> active is performed
// 	_ = keeper.SetPulse(core.PulseNumber(1))
// 	keeper.Sync(true, core.PulseNumber(1))
// 	_ = keeper.SetPulse(core.PulseNumber(2))
// 	keeper.Sync(true, core.PulseNumber(2))
// 	assert.Equal(t, 1, len(keeper.GetActiveNodes()))
// }
//
// func TestNodekeeper_unsyncUpdatePulse(t *testing.T) {
// 	keeper := NewNodeKeeper(core.RecordRef{}, time.Hour)
// 	_ = keeper.SetPulse(core.PulseNumber(0))
//
// 	keeper.AddUnsync(newActiveNode(0, 0))
// 	_ = keeper.SetPulse(core.PulseNumber(1))
// 	nodePulse := keeper.GetUnsync()[0].PulseNum
// 	assert.Equal(t, uint32(1), uint32(nodePulse))
// }
