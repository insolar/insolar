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

	"github.com/anacrolix/sync"
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
)

type NodeKeeper interface {
	// GetActiveNode get active node by its reference. Returns nil if node is not found
	GetActiveNode(ref core.RecordRef) *core.ActiveNode
	// GetActiveNodes get active nodes.
	GetActiveNodes() []*core.ActiveNode
	// GetUnsyncHash get hash computed based on the list of unsync nodes, and the size of this list.
	GetUnsyncHash() (hash []byte, unsyncCount int, err error)
	// GetUnsync gets the local unsync list (excluding other nodes unsync lists)
	GetUnsync() []*core.ActiveNode
	// Sync initiate transferring unsync -> sync, sync -> active. If approved is false, unsync is not transferred to sync.
	// This function also sets internal PulseNumber to `number`.
	Sync(number core.PulseNumber, approved bool)
	// AddUnsync add unsync node to the local unsync list.
	// Returns error if node's PulseNumber is not equal to the NodeKeeper internal PulseNumber
	AddUnsync(*core.ActiveNode) error
	// AddUnsyncGossip merge unsync list from another node to the local unsync list.
	// Returns error if one of the nodes' PulseNumber is not equal to the NodeKeeper internal PulseNumber
	AddUnsyncGossip([]*core.ActiveNode) error
}

func NewNodeKeeper() NodeKeeper {
	return &nodekeeper{}
}

type nodekeeper struct {
	pulse core.PulseNumber

	activeLock sync.RWMutex
	active     map[core.RecordRef]*core.ActiveNode

	syncLock sync.Mutex
	sync     []*core.ActiveNode

	unsyncLock   sync.Mutex
	unsync       []*core.ActiveNode
	unsyncGossip map[core.RecordRef]*core.ActiveNode
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

func (nk *nodekeeper) Sync(number core.PulseNumber, approved bool) {
	nk.pulse = number
	// TODO: unsync discard logic due to timeout
	// TODO: sync logic
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

func (nk *nodekeeper) AddUnsync(node *core.ActiveNode) error {
	nk.unsyncLock.Lock()
	defer nk.unsyncLock.Unlock()

	checkedList := []*core.ActiveNode{node}
	if err := nk.checkPulse(checkedList); err != nil {
		return err
	}

	nk.unsync = append(nk.unsync, node)
	return nil
}

func (nk *nodekeeper) AddUnsyncGossip(nodes []*core.ActiveNode) error {
	nk.unsyncLock.Lock()
	defer nk.unsyncLock.Unlock()

	if err := nk.checkPulse(nodes); err != nil {
		return err
	}

	for _, node := range nodes {
		nk.unsyncGossip[node.NodeID] = node
	}
	return nil
}
