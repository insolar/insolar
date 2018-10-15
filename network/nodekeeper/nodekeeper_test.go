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
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
)

func newActiveNode(ref byte) *core.ActiveNode {
	return &core.ActiveNode{
		NodeID:    core.RecordRef{ref},
		PulseNum:  core.PulseNumber(0),
		State:     core.NodeActive,
		Role:      core.RoleUnknown,
		PublicKey: []byte{0, 0, 0},
	}
}

func newSelfNode(ref core.RecordRef) *core.ActiveNode {
	return &core.ActiveNode{
		NodeID:    ref,
		PulseNum:  core.PulseNumber(0),
		State:     core.NodeActive,
		Role:      core.RoleUnknown,
		PublicKey: []byte{0, 0, 0},
	}
}

func newNodeKeeper() NodeKeeper {
	id := core.RecordRef{255}
	keeper := NewNodeKeeper(id)
	keeper.AddActiveNodes([]*core.ActiveNode{newSelfNode(id)})
	return keeper
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
	assert.Equal(t, 1, len(list.GetUnsync()))
}

func TestNodekeeper_AddUnsync2(t *testing.T) {
	keeper := newNodeKeeper()
	success, list := keeper.SetPulse(core.PulseNumber(0))
	err := keeper.AddUnsync(newActiveNode(0))
	assert.NoError(t, err)
	assert.True(t, success)
	assert.Equal(t, 0, len(list.GetUnsync()))
}

func TestNodekeeper_AddUnsync3(t *testing.T) {
	keeper := newNodeKeeper()
	err := keeper.AddUnsync(newActiveNode(0))
	success, list := keeper.SetPulse(core.PulseNumber(0))
	err = keeper.AddUnsync(newActiveNode(1))
	assert.NoError(t, err)
	assert.True(t, success)
	assert.Equal(t, 1, len(list.GetUnsync()))
}

func TestNodekeeper_pipeline(t *testing.T) {
	keeper := newNodeKeeper()
	for i := 0; i < 4; i++ {
		err := keeper.AddUnsync(newActiveNode(byte(2 * i)))
		assert.NoError(t, err)
		pulse := core.PulseNumber(i)
		success, list := keeper.SetPulse(pulse)
		assert.True(t, success)
		err = keeper.AddUnsync(newActiveNode(byte(2*i + 1)))
		assert.NoError(t, err)
		keeper.Sync(list.GetUnsync(), pulse)
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
	assert.Equal(t, 1, len(list.GetUnsync()))
	keeper.Sync(list.GetUnsync(), pulse)
	// second sync should be ignored because pulse has not changed
	keeper.Sync(list.GetUnsync(), pulse)
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
	keeper.Sync(list.GetUnsync(), pulse)
	_, _ = keeper.SetPulse(core.PulseNumber(1))
	_, _ = keeper.SetPulse(core.PulseNumber(2))
	// node with ref 0 advanced to active list
	assert.Equal(t, 2, len(keeper.GetActiveNodes()))
	assert.NotNil(t, keeper.GetActiveNode(core.RecordRef{0}))
}

func TestNodekeeper_outdatedSync(t *testing.T) {
	keeper := newNodeKeeper()
	num := 4
	wg := sync.WaitGroup{}
	wg.Add(num)
	for i := 0; i < num; i++ {
		time.Sleep(100 * time.Millisecond)
		go func(k NodeKeeper, i int) {
			_ = k.AddUnsync(newActiveNode(byte(2 * i)))
			_ = k.AddUnsync(newActiveNode(byte(2*i + 1)))
			pulse := core.PulseNumber(i)
			success, list := k.SetPulse(pulse)
			assert.True(t, success)
			// imitate long consensus process
			time.Sleep(200 * time.Millisecond)
			k.Sync(list.GetUnsync(), pulse)
			wg.Done()
		}(keeper, i)
	}
	wg.Wait()
	// All Syncs calls are executed out of date
	// So, no nodes should advance to active list (we should have only 1 self node in active)
	assert.Equal(t, 1, len(keeper.GetActiveNodes()))
}

func TestNodekeeper_SetPulse(t *testing.T) {
	keeper := newNodeKeeper()
	success, _ := keeper.SetPulse(core.PulseNumber(10))
	assert.True(t, success)
	// Pulses should pass in ascending order
	success, _ = keeper.SetPulse(core.PulseNumber(9))
	assert.False(t, success)
}

func TestNodekeeper_notifyWaiters(t *testing.T) {
	keeper := newNodeKeeper()
	success, _ := keeper.SetPulse(core.PulseNumber(10))
	assert.True(t, success)
}
