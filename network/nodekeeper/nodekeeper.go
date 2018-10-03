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
	"sync/atomic"
	"time"

	"github.com/anacrolix/sync"
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
)

type NodeKeeper interface {
	// GetActiveNode get active node by its reference. Returns nil if node is not found.
	GetActiveNode(ref core.RecordRef) *core.ActiveNode
	// GetActiveNodes get active nodes.
	GetActiveNodes() []*core.ActiveNode
	// GetUnsyncHash get hash computed based on the list of unsync nodes, and the size of this list.
	GetUnsyncHash() (hash []byte, unsyncCount int, err error)
	// GetUnsync gets the local unsync list (excluding other nodes unsync lists).
	GetUnsync() []*core.ActiveNode
	// SetPulse sets internal PulseNumber to number.
	SetPulse(number core.PulseNumber)
	// Sync initiate transferring unsync -> sync, sync -> active. If approved is false, unsync is not transferred to sync.
	Sync(approved bool)
	// AddUnsync add unsync node to the local unsync list.
	// Returns error if node's PulseNumber is not equal to the NodeKeeper internal PulseNumber.
	AddUnsync(*core.ActiveNode) error
	// AddUnsyncGossip merge unsync list from another node to the local unsync list.
	// Returns error if:
	// 1. One of the nodes' PulseNumber is not equal to the NodeKeeper internal PulseNumber;
	// 2. One of the nodes' reference is equal to one of the local unsync nodes' reference.
	AddUnsyncGossip([]*core.ActiveNode) error
}

// NewNodeKeeper create new NodeKeeper. unsyncDiscardAfter = timeout after which each unsync node is discarded
func NewNodeKeeper(unsyncDiscardAfter time.Duration) NodeKeeper {
	return &nodekeeper{
		timeout:      unsyncDiscardAfter,
		active:       make(map[core.RecordRef]*core.ActiveNode),
		sync:         make([]*core.ActiveNode, 0),
		unsync:       make([]*core.ActiveNode, 0),
		unsyncGossip: make(map[core.RecordRef]*core.ActiveNode),
	}
}

type nodekeeper struct {
	pulse   uint32
	timeout time.Duration

	activeLock sync.RWMutex
	active     map[core.RecordRef]*core.ActiveNode
	sync       []*core.ActiveNode

	unsyncLock    sync.Mutex
	unsync        []*core.ActiveNode
	unsyncTimeout []time.Time
	unsyncGossip  map[core.RecordRef]*core.ActiveNode
}

func (nk *nodekeeper) GetActiveNodes() []*core.ActiveNode {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	result := make([]*core.ActiveNode, len(nk.active))
	index := 0
	for _, node := range nk.active {
		result[index] = node
		index++
	}
	return result
}

func (nk *nodekeeper) GetActiveNode(ref core.RecordRef) *core.ActiveNode {
	nk.activeLock.RLock()
	defer nk.activeLock.RUnlock()

	return nk.active[ref]
}

func (nk *nodekeeper) GetUnsyncHash() ([]byte, int, error) {
	nk.unsyncLock.Lock()
	unsync := nk.collectUnsync()
	nk.unsyncLock.Unlock()
	hash, err := calculateHash(unsync)
	if err != nil {
		return nil, 0, err
	}
	return hash, len(unsync), nil
}

func (nk *nodekeeper) GetUnsync() []*core.ActiveNode {
	nk.unsyncLock.Lock()
	defer nk.unsyncLock.Unlock()

	result := make([]*core.ActiveNode, len(nk.unsync))
	copy(result, nk.unsync)
	return result
}

func (nk *nodekeeper) SetPulse(number core.PulseNumber) {
	atomic.StoreUint32(&nk.pulse, uint32(number))
}

func (nk *nodekeeper) Sync(approved bool) {
	nk.unsyncLock.Lock()
	nk.activeLock.Lock()

	defer func() {
		nk.activeLock.Unlock()
		nk.unsyncLock.Unlock()
	}()

	// sync -> active
	for _, node := range nk.sync {
		nk.active[node.NodeID] = node
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
}

func (nk *nodekeeper) AddUnsync(node *core.ActiveNode) error {
	nk.unsyncLock.Lock()
	defer nk.unsyncLock.Unlock()

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
}

func (nk *nodekeeper) checkPulse(nodes []*core.ActiveNode) error {
	pulse := core.PulseNumber(atomic.LoadUint32(&nk.pulse))
	for _, node := range nodes {
		if node.PulseNum != pulse {
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

func calculateHash(list []*core.ActiveNode) (result []byte, err error) {
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
