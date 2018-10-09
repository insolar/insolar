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
	"encoding/binary"
	"fmt"
	"hash"
	"sort"
	"sync"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
)

type NodeKeeper interface {
	// GetSelf get active node for the current insolard. Returns nil if the current insolard is not an active node
	GetSelf() *core.ActiveNode
	// GetActiveNode get active node by its reference. Returns nil if node is not found.
	GetActiveNode(ref core.RecordRef) *core.ActiveNode
	// GetActiveNodes get active nodes.
	GetActiveNodes() []*core.ActiveNode
	// AddActiveNodes set active nodes.
	AddActiveNodes([]*core.ActiveNode)
	// GetUnsyncHash get hash computed based on the list of unsync nodes, and the size of this list.
	GetUnsyncHash() (hash []byte, unsyncCount int, err error)
	// GetUnsync gets the local unsync list (excluding other nodes unsync lists).
	GetUnsync() []*core.ActiveNode
	// SetPulse sets internal PulseNumber to number. Returns true if set was successful, false if number is less
	// or equal to internal PulseNumber
	SetPulse(number core.PulseNumber) bool
	// Sync initiate transferring unsync -> sync, sync -> active. If approved is false, unsync is not transferred to sync.
	// If number is less than internal PulseNumber then ignore Sync.
	Sync(approved bool, number core.PulseNumber)
	// AddUnsync add unsync node to the local unsync list.
	// Returns error if node's PulseNumber is not equal to the NodeKeeper internal PulseNumber.
	AddUnsync(*core.ActiveNode) error
	// AddUnsyncGossip merge unsync list from another node to the local unsync list.
	// Returns error if:
	// 1. One of the nodes' PulseNumber is not equal to the NodeKeeper internal PulseNumber;
	// 2. One of the nodes' reference is equal to one of the local unsync nodes' reference.
	AddUnsyncGossip([]*core.ActiveNode) error
}

// NewNodeKeeper create new NodeKeeper. unsyncDiscardAfter = timeout after which each unsync node is discarded.
func NewNodeKeeper(nodeID core.RecordRef, unsyncDiscardAfter time.Duration) NodeKeeper {
	return &nodekeeper{
		nodeID:       nodeID,
		state:        undefined,
		timeout:      unsyncDiscardAfter,
		active:       make(map[core.RecordRef]*core.ActiveNode),
		sync:         make([]*core.ActiveNode, 0),
		unsync:       make([]*core.ActiveNode, 0),
		unsyncGossip: make(map[core.RecordRef]*core.ActiveNode),
	}
}

type nodekeeperState uint8

const (
	undefined = nodekeeperState(iota + 1)
	awaitUnsync
	hashCalculated
	synced
)

type nodekeeper struct {
	nodeID          core.RecordRef
	self            *core.ActiveNode
	state           nodekeeperState
	pulse           core.PulseNumber
	timeout         time.Duration
	cacheUnsyncCalc []byte
	cacheUnsyncSize int

	activeLock sync.RWMutex
	active     map[core.RecordRef]*core.ActiveNode
	sync       []*core.ActiveNode

	unsyncLock    sync.Mutex
	unsync        []*core.ActiveNode
	unsyncTimeout []time.Time
	unsyncGossip  map[core.RecordRef]*core.ActiveNode
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
			log.Warnf("AddActiveNodes: trying to add self ID: %s. Typically it must happen via Sync", nk.nodeID)
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

func (nk *nodekeeper) GetUnsyncHash() ([]byte, int, error) {
	nk.unsyncLock.Lock()
	defer nk.unsyncLock.Unlock()

	if nk.state != awaitUnsync {
		log.Warn("NodeKeeper: GetUnsyncHash called more than once during one pulse")
		return nk.cacheUnsyncCalc, nk.cacheUnsyncSize, nil
	}
	unsync := nk.collectUnsync()
	hash, err := CalculateHash(unsync)
	if err != nil {
		return nil, 0, err
	}
	nk.cacheUnsyncCalc, nk.cacheUnsyncSize = hash, len(unsync)
	nk.state = hashCalculated
	return nk.cacheUnsyncCalc, nk.cacheUnsyncSize, nil
}

func (nk *nodekeeper) GetUnsync() []*core.ActiveNode {
	nk.unsyncLock.Lock()
	defer nk.unsyncLock.Unlock()

	result := make([]*core.ActiveNode, len(nk.unsync))
	copy(result, nk.unsync)
	return result
}

func (nk *nodekeeper) SetPulse(number core.PulseNumber) bool {
	nk.unsyncLock.Lock()
	defer nk.unsyncLock.Unlock()

	if nk.state == undefined {
		nk.pulse = number
		nk.state = awaitUnsync
		return true
	}

	if number <= nk.pulse {
		log.Warnf("NodeKeeper: ignored SetPulse call with number=%d while current=%d", uint32(number), uint32(nk.pulse))
		return false
	}

	if nk.state == hashCalculated || nk.state == awaitUnsync {
		log.Warn("NodeKeeper: SetPulse called not from `undefined` or `synced` state")
		nk.activeLock.Lock()
		nk.syncUnsafe(false)
		nk.activeLock.Unlock()
	}

	nk.pulse = number
	nk.state = awaitUnsync
	nk.invalidateCache()
	nk.updateUnsyncPulse()
	return true
}

func (nk *nodekeeper) Sync(approved bool, number core.PulseNumber) {
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

	nk.syncUnsafe(approved)
}

func (nk *nodekeeper) AddUnsync(node *core.ActiveNode) error {
	nk.unsyncLock.Lock()
	defer nk.unsyncLock.Unlock()

	if nk.state != awaitUnsync {
		return errors.New("Cannot add node to unsync list: try again in next pulse slot")
	}

	checkedList := []*core.ActiveNode{node}
	if err := nk.checkPulse(checkedList); err != nil {
		return errors.Wrap(err, "Error adding local unsync node")
	}

	nk.unsync = append(nk.unsync, node)
	tm := time.Now().Add(nk.timeout)
	nk.unsyncTimeout = append(nk.unsyncTimeout, tm)
	return nil
}

func (nk *nodekeeper) AddUnsyncGossip(nodes []*core.ActiveNode) error {
	nk.unsyncLock.Lock()
	defer nk.unsyncLock.Unlock()

	if nk.state != awaitUnsync {
		return errors.New("Cannot add node to unsync list: try again in next pulse slot")
	}

	if err := nk.checkPulse(nodes); err != nil {
		return errors.Wrap(err, "Error adding unsync gossip nodes")
	}

	if err := nk.checkReference(nodes); err != nil {
		return errors.Wrap(err, "Error adding unsync gossip nodes")
	}

	for _, node := range nodes {
		nk.unsyncGossip[node.NodeID] = node
	}
	return nil
}

func (nk *nodekeeper) syncUnsafe(approved bool) {
	// sync -> active
	for _, node := range nk.sync {
		nk.active[node.NodeID] = node
		if node.NodeID.Equal(nk.nodeID) {
			log.Infof("Sync: current node %s reached the active node list", nk.nodeID)
			nk.self = node
		}
	}

	if approved {
		// unsync -> sync
		unsync := nk.collectUnsync()
		nk.sync = unsync
		// clear unsync
		nk.unsync = make([]*core.ActiveNode, 0)
	} else {
		// clear sync
		nk.sync = make([]*core.ActiveNode, 0)
		nk.discardTimedOutUnsync()
	}
	// clear unsyncGossip
	nk.unsyncGossip = make(map[core.RecordRef]*core.ActiveNode)
	nk.state = synced
}

func (nk *nodekeeper) discardTimedOutUnsync() {
	index := 0
	for _, tm := range nk.unsyncTimeout {
		if tm.After(time.Now()) {
			break
		}
		index++
	}
	if index == 0 {
		return
	}
	// discard all unsync nodes before index
	nk.unsyncTimeout = nk.unsyncTimeout[index:]
	nk.unsync = nk.unsync[index:]
	log.Infof("NodeKeeper: discarded %d unsync nodes due to timeout", index)
}

func (nk *nodekeeper) checkPulse(nodes []*core.ActiveNode) error {
	for _, node := range nodes {
		if node.PulseNum != nk.pulse {
			return errors.Errorf("Node ID:%s pulse:%d is not equal to NodeKeeper current pulse:%d",
				node.NodeID.String(), node.PulseNum, nk.pulse)
		}
	}
	return nil
}

func (nk *nodekeeper) checkReference(nodes []*core.ActiveNode) error {
	// quadratic, should not be a problem because unsync lists are usually empty or have few elements
	for _, localNode := range nk.unsync {
		for _, node := range nodes {
			if node.NodeID.Equal(localNode.NodeID) {
				return errors.Errorf("Node %s cannot be added to gossip unsync list "+
					"because it is in local unsync list", node.NodeID.String())
			}
		}
	}
	return nil
}

func (nk *nodekeeper) collectUnsync() []*core.ActiveNode {
	unsync := make([]*core.ActiveNode, len(nk.unsyncGossip)+len(nk.unsync))
	index := 0
	for _, node := range nk.unsyncGossip {
		unsync[index] = node
		index++
	}
	copy(unsync[index:], nk.unsync)
	return unsync
}

func (nk *nodekeeper) invalidateCache() {
	nk.cacheUnsyncCalc = nil
	nk.cacheUnsyncSize = 0
}

func (nk *nodekeeper) updateUnsyncPulse() {
	for _, node := range nk.unsync {
		node.PulseNum = nk.pulse
	}
	count := len(nk.unsync)
	if count != 0 {
		log.Infof("NodeKeeper: updated pulse for %d stored unsync nodes", count)
	}
}

func hashWriteChecked(hash hash.Hash, data []byte) {
	n, err := hash.Write(data)
	if n != len(data) {
		panic(fmt.Sprintf("Error writing hash. Bytes expected: %d; bytes actual: %d", len(data), n))
	}
	if err != nil {
		panic(err.Error())
	}
}

func calculateNodeHash(node *core.ActiveNode) []byte {
	hash := sha3.New224()
	hashWriteChecked(hash, node.NodeID[:])
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(node.JetRoles))
	hashWriteChecked(hash, b[:])
	binary.LittleEndian.PutUint32(b, uint32(node.PulseNum))
	hashWriteChecked(hash, b[:4])
	b[0] = byte(node.State)
	hashWriteChecked(hash, b[:1])
	hashWriteChecked(hash, node.PublicKey)
	return hash.Sum(nil)
}

func CalculateHash(list []*core.ActiveNode) (result []byte, err error) {
	sort.Slice(list[:], func(i, j int) bool {
		return bytes.Compare(list[i].NodeID[:], list[j].NodeID[:]) < 0
	})

	// catch possible panic from hashWriteChecked in this function and in all calculateNodeHash funcs
	defer func() {
		if r := recover(); r != nil {
			result, err = nil, fmt.Errorf("error calculating hash: %s", r)
		}
	}()

	hash := sha3.New224()
	for _, node := range list {
		nodeHash := calculateNodeHash(node)
		hashWriteChecked(hash, nodeHash)
	}
	return hash.Sum(nil), nil
}
