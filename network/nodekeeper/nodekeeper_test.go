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
	"time"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
)

const (
	nullHash = "6b4e03423667dbb73b6e15454f0eb1abd4597f9a1b078e3f5b5a6bc7"
)

func newActiveNode(ref byte, pulse int) *core.ActiveNode {
	var mask core.JetRoleMask
	mask.Set(core.RoleVirtualExecutor)

	return &core.ActiveNode{
		NodeID:    core.RecordRef{ref},
		PulseNum:  core.PulseNumber(pulse),
		State:     core.NodeActive,
		JetRoles:  mask,
		PublicKey: []byte{0, 0, 0},
	}
}

func TestNodekeeper_calculateNodeHash(t *testing.T) {
	hash1, _ := calculateHash(nil)
	hash2, _ := calculateHash([]*core.ActiveNode{})

	assert.Equal(t, nullHash, hex.EncodeToString(hash1))
	assert.Equal(t, hash1, hash2)

	activeNode1 := newActiveNode(0, 0)
	activeNode2 := newActiveNode(0, 0)

	activeNode1Slice := []*core.ActiveNode{activeNode1}
	activeNode2Slice := []*core.ActiveNode{activeNode2}

	hash1, _ = calculateHash(activeNode1Slice)
	hash2, _ = calculateHash(activeNode2Slice)
	assert.Equal(t, hash1, hash2)
	activeNode2.NodeID = core.RecordRef{1}
	hash2, _ = calculateHash(activeNode2Slice)
	assert.NotEqual(t, hash1, hash2)

	// nodes order in slice should not affect hash calculating
	slice1 := []*core.ActiveNode{activeNode1, activeNode2}
	slice2 := []*core.ActiveNode{activeNode2, activeNode1}
	hash1, _ = calculateHash(slice1)
	hash2, _ = calculateHash(slice2)
	assert.Equal(t, hash1, hash2)
}

func TestNodekeeper_AddUnsync(t *testing.T) {
	keeper := NewNodeKeeper(time.Hour)
	keeper.SetPulse(core.PulseNumber(0))
	_ = keeper.AddUnsync(newActiveNode(0, 0))
	_ = keeper.AddUnsync(newActiveNode(1, 0))
	gossip := []*core.ActiveNode{newActiveNode(2, 0), newActiveNode(3, 0)}
	_ = keeper.AddUnsyncGossip(gossip)
	keeper.Sync(true)
	keeper.SetPulse(core.PulseNumber(1))
	keeper.Sync(true)
	assert.Equal(t, 4, len(keeper.GetActiveNodes()))
	for i := 0; i < 4; i++ {
		assert.NotNil(t, keeper.GetActiveNode(core.RecordRef{byte(i)}))
	}
}

func TestNodekeeper_GetUnsyncHash(t *testing.T) {
	keeper := NewNodeKeeper(time.Hour)
	keeper.SetPulse(core.PulseNumber(0))
	hash, count, _ := keeper.GetUnsyncHash()
	assert.Equal(t, nullHash, hex.EncodeToString(hash))
	assert.Equal(t, 0, count)

	keeper.SetPulse(core.PulseNumber(1))
	_ = keeper.AddUnsync(newActiveNode(0, 1))
	_ = keeper.AddUnsyncGossip([]*core.ActiveNode{newActiveNode(1, 1)})

	keeper2 := NewNodeKeeper(time.Hour)
	keeper2.SetPulse(core.PulseNumber(1))
	_ = keeper2.AddUnsync(newActiveNode(1, 1))
	_ = keeper2.AddUnsyncGossip([]*core.ActiveNode{newActiveNode(0, 1)})

	hash, count, _ = keeper.GetUnsyncHash()
	hash2, count2, _ := keeper2.GetUnsyncHash()
	assert.Equal(t, hash, hash2)
	assert.Equal(t, count, count2)
}

func TestNodekeeper_AddUnsync_checks(t *testing.T) {
	keeper := NewNodeKeeper(time.Hour)
	keeper.SetPulse(core.PulseNumber(0))

	// Unsync node pulse number should be equal to the NodeKeeper pulse number
	err := keeper.AddUnsync(newActiveNode(0, 1))
	assert.Error(t, err)
	err = keeper.AddUnsync(newActiveNode(0, 0))
	assert.NoError(t, err)

	// Gossip unsync node should not have reference id equal to one of the local unsync nodes
	err = keeper.AddUnsyncGossip([]*core.ActiveNode{newActiveNode(0, 0)})
	assert.Error(t, err)
	// Gossip unsync node pulse number should be equal to the NodeKeeper pulse number
	err = keeper.AddUnsyncGossip([]*core.ActiveNode{newActiveNode(1, 1)})
	assert.Error(t, err)
	err = keeper.AddUnsyncGossip([]*core.ActiveNode{newActiveNode(1, 0)})
	assert.NoError(t, err)
}

func TestNodekeeper_discardTimedOutUnsync(t *testing.T) {
	keeper := NewNodeKeeper(250 * time.Millisecond)
	for i := 0; i < 4; i++ {
		keeper.SetPulse(core.PulseNumber(i))
		_ = keeper.AddUnsync(newActiveNode(byte(i), i))
		time.Sleep(100 * time.Millisecond)
		keeper.Sync(false)
	}
	assert.Equal(t, 2, len(keeper.GetUnsync()))
}

func TestNodekeeper_cache(t *testing.T) {
	keeper := &nodekeeper{
		state:        undefined,
		timeout:      time.Hour,
		active:       make(map[core.RecordRef]*core.ActiveNode),
		sync:         make([]*core.ActiveNode, 0),
		unsync:       make([]*core.ActiveNode, 0),
		unsyncGossip: make(map[core.RecordRef]*core.ActiveNode),
	}
	keeper.SetPulse(core.PulseNumber(0))
	assert.Equal(t, awaitUnsync, keeper.state)
	err := keeper.AddUnsync(newActiveNode(0, 0))
	assert.NoError(t, err)
	keeper.AddUnsyncGossip([]*core.ActiveNode{newActiveNode(1, 0)})
	assert.NoError(t, err)
	hash1, _, _ := keeper.GetUnsyncHash()
	hash2, _, _ := keeper.GetUnsyncHash()
	assert.Equal(t, hash1, hash2)
	assert.Equal(t, hashCalculated, keeper.state)
	err = keeper.AddUnsync(newActiveNode(2, 0))
	assert.Error(t, err)
	keeper.AddUnsyncGossip([]*core.ActiveNode{newActiveNode(3, 0)})
	assert.Error(t, err)
	keeper.Sync(true)
	assert.Equal(t, synced, keeper.state)

	keeper.SetPulse(core.PulseNumber(1))
	assert.Equal(t, awaitUnsync, keeper.state)
}

func TestNodekeeper_AddActiveNodes(t *testing.T) {
	keeper := NewNodeKeeper(time.Hour)
	keeper.SetPulse(core.PulseNumber(0))

	node2 := newActiveNode(0, 0)
	node1 := newActiveNode(1, 0)
	nodes := []*core.ActiveNode{node1, node2}
	keeper.AddActiveNodes(nodes)

	assert.Equal(t, 2, len(keeper.GetActiveNodes()))
	assert.NotNil(t, keeper.GetActiveNode(core.RecordRef{0}))
	assert.NotNil(t, keeper.GetActiveNode(core.RecordRef{1}))
}

func TestNodekeeper_transitions1(t *testing.T) {
	keeper := NewNodeKeeper(time.Hour)
	keeper.SetPulse(core.PulseNumber(0))

	keeper.AddUnsync(newActiveNode(0, 0))
	// check that Sync is not called and the transition unsync -> sync -> active is not performed
	keeper.SetPulse(core.PulseNumber(0))
	keeper.SetPulse(core.PulseNumber(0))
	assert.Equal(t, 0, len(keeper.GetActiveNodes()))
}

func TestNodekeeper_transitions2(t *testing.T) {
	keeper := NewNodeKeeper(time.Hour)
	keeper.SetPulse(core.PulseNumber(0))

	keeper.AddUnsync(newActiveNode(0, 0))
	// check that Sync is called correctly every time and the transition unsync -> sync -> active is performed
	keeper.SetPulse(core.PulseNumber(1))
	keeper.Sync(true)
	keeper.SetPulse(core.PulseNumber(2))
	keeper.Sync(true)
	assert.Equal(t, 1, len(keeper.GetActiveNodes()))
}
