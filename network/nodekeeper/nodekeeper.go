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
	"bytes"
	"sort"
	"sync"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

// NodeKeeper manages unsync, sync and active lists
type NodeKeeper interface {
	// GetID get current node ID
	GetID() core.RecordRef
	// GetSelf get active node for the current insolard. Returns nil if the current insolard is not an active node.
	GetSelf() *core.ActiveNode
	// GetActiveNode get active node by its reference. Returns nil if node is not found.
	GetActiveNode(ref core.RecordRef) *core.ActiveNode
	// GetActiveNodes get active nodes.
	GetActiveNodes() []*core.ActiveNode
	// AddActiveNodes add active nodes.
	AddActiveNodes([]*core.ActiveNode)
	// SetPulse sets internal PulseNumber to number. Returns true if set was successful, false if number is less
	// or equal to internal PulseNumber. If set is successful, returns collected unsync list and starts collecting new unsync list
	SetPulse(number core.PulseNumber) (bool, *UnsyncList)
	// Sync initiates transferring syncCandidates -> sync, sync -> active.
	// If number is less than internal PulseNumber then ignore Sync.
	Sync(syncCandidates []*core.ActiveNode, number core.PulseNumber)
	// AddUnsync add unsync node to the unsync list. Returns error if current node is not active and cannot participate in consensus
	AddUnsync(*core.ActiveNode) error
	// GetUnsyncHolder get unsync list executed in consensus for specific pulse.
	// 1. If pulse is less than internal NodeKeeper pulse, returns error.
	// 2. If pulse is equal to internal NodeKeeper pulse, returns unsync list holder for currently executed consensus.
	// 3. If pulse is more than internal NodeKeeper pulse, blocks till next SetPulse or duration timeout and then acts like in par. 2
	GetUnsyncHolder(pulse core.PulseNumber, duration time.Duration) (*UnsyncList, error)
}

// NewNodeKeeper create new NodeKeeper
func NewNodeKeeper(nodeID core.RecordRef) NodeKeeper {
	return &nodekeeper{
		nodeID:        nodeID,
		state:         undefined,
		active:        make(map[core.RecordRef]*core.ActiveNode),
		sync:          make([]*core.ActiveNode, 0),
		unsync:        make([]*core.ActiveNode, 0),
		unsyncWaiters: make([]chan *UnsyncList, 0),
	}
}

type nodekeeperState uint8

const (
	undefined = nodekeeperState(iota + 1)
	pulseSet
	synced
)

type nodekeeper struct {
	nodeID core.RecordRef
	self   *core.ActiveNode
	state  nodekeeperState
	pulse  core.PulseNumber

	activeLock sync.RWMutex
	active     map[core.RecordRef]*core.ActiveNode
	sync       []*core.ActiveNode

	unsyncLock    sync.Mutex
	unsync        []*core.ActiveNode
	unsyncList    *UnsyncList
	unsyncWaiters []chan *UnsyncList
}

func (nk *nodekeeper) GetID() core.RecordRef {
	return nk.nodeID
}

func (nk *nodekeeper) GetSelf() *core.ActiveNode {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	return nk.self
}

func (nk *nodekeeper) GetActiveNodes() []*core.ActiveNode {
	nk.activeLock.RLock()
	result := make([]*core.ActiveNode, len(nk.active))
	index := 0
	for _, node := range nk.active {
		result[index] = node
		index++
	}
	nk.activeLock.RUnlock()
	// Sort active nodes to return list with determinate order on every node.
	// If we have more than 10k nodes, we need to optimize this
	sort.Slice(result, func(i, j int) bool {
		return bytes.Compare(result[i].NodeID[:], result[j].NodeID[:]) < 0
	})
	return result
}

func (nk *nodekeeper) AddActiveNodes(nodes []*core.ActiveNode) {
	nk.activeLock.Lock()
	defer nk.activeLock.Unlock()

	for _, node := range nodes {
		if node.NodeID.Equal(nk.nodeID) {
			nk.self = node
		}
		nk.active[node.NodeID] = node
	}
}

func (nk *nodekeeper) GetActiveNode(ref core.RecordRef) *core.ActiveNode {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	return nk.active[ref]
}

func (nk *nodekeeper) SetPulse(number core.PulseNumber) (bool, *UnsyncList) {
	nk.unsyncLock.Lock()
	defer nk.unsyncLock.Unlock()

	if nk.state == undefined {
		return true, nk.collectUnsync(number)
	}

	if number <= nk.pulse {
		log.Warnf("NodeKeeper: ignored SetPulse call with number=%d while current=%d", uint32(number), uint32(nk.pulse))
		return false, nil
	}

	if nk.state == pulseSet {
		log.Warn("NodeKeeper: SetPulse called from `pulseSet` state")
		nk.activeLock.Lock()
		nk.syncUnsafe(nil)
		nk.activeLock.Unlock()
	}

	return true, nk.collectUnsync(number)
}

func (nk *nodekeeper) Sync(syncCandidates []*core.ActiveNode, number core.PulseNumber) {
	nk.unsyncLock.Lock()
	nk.activeLock.Lock()

	defer func() {
		nk.activeLock.Unlock()
		nk.unsyncLock.Unlock()
	}()

	if nk.state == synced || nk.state == undefined {
		log.Warn("NodeKeeper: ignored Sync call from `synced` or `undefined` state")
		return
	}

	if nk.pulse > number {
		log.Warnf("NodeKeeper: ignored Sync call because passed number %d is less than internal number %d",
			number, nk.pulse)
		return
	}

	nk.syncUnsafe(syncCandidates)
}

func (nk *nodekeeper) AddUnsync(node *core.ActiveNode) error {
	nk.unsyncLock.Lock()
	defer nk.unsyncLock.Unlock()

	if nk.self == nil {
		return errors.New("cannot add node to unsync list: current node is not active")
	}

	nk.unsync = append(nk.unsync, node)
	return nil
}

func (nk *nodekeeper) GetUnsyncHolder(pulse core.PulseNumber, duration time.Duration) (*UnsyncList, error) {
	nk.unsyncLock.Lock()
	currentPulse := nk.pulse

	if currentPulse == pulse {
		result := nk.unsyncList
		nk.unsyncLock.Unlock()
		return result, nil
	}
	if currentPulse > pulse || duration < 0 {
		nk.unsyncLock.Unlock()
		return nil, errors.Errorf("GetUnsyncHolder called with pulse %d, but current NodeKeeper pulse is %d",
			pulse, currentPulse)
	}
	ch := make(chan *UnsyncList, 1)
	nk.unsyncWaiters = append(nk.unsyncWaiters, ch)
	nk.unsyncLock.Unlock()
	// TODO: timeout
	result := <-ch
	if result.GetPulse() != pulse {
		return nil, errors.Errorf("GetUnsyncHolder called with pulse %d, but current UnsyncHolder pulse is %d",
			pulse, result.GetPulse())
	}
	return result, nil
}

func (nk *nodekeeper) syncUnsafe(syncCandidates []*core.ActiveNode) {
	// sync -> active
	for _, node := range nk.sync {
		nk.active[node.NodeID] = node
	}
	// unsync -> sync
	nk.sync = syncCandidates
	nk.state = synced
}

func (nk *nodekeeper) collectUnsync(number core.PulseNumber) *UnsyncList {
	nk.pulse = number
	nk.state = pulseSet

	for _, node := range nk.unsync {
		node.PulseNum = nk.pulse
	}
	tmp := nk.unsync
	nk.unsync = make([]*core.ActiveNode, 0)
	nk.unsyncList = NewUnsyncHolder(nk.pulse, tmp)
	if len(nk.unsyncWaiters) == 0 {
		return nk.unsyncList
	}
	// notify waiters that new unsync holder is available for read
	for _, ch := range nk.unsyncWaiters {
		ch <- nk.unsyncList
		close(ch)
	}
	nk.unsyncWaiters = make([]chan *UnsyncList, 0)
	return nk.unsyncList
}
