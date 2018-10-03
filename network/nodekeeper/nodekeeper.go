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
	// Sync initiate transferring unsync -> sync, sync -> active. If approved is false, unsync is not transferred to sync
	Sync(approved bool)
	// AddUnsync add unsync node to the local unsync list
	AddUnsync(*core.ActiveNode)
	// AddUnsyncGossip merge unsync list from another node to the local unsync list
	AddUnsyncGossip([]*core.ActiveNode)
}

func NewNodeKeeper() NodeKeeper {
	return &nodekeeper{}
}

type nodekeeper struct {
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

func calculateNodeHash(node *core.ActiveNode) (result []byte, err error) {
	// TODO: check Write
	hash := sha3.New224()

	hash.Write(node.NodeID[:])
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(node.JetRoles))
	hash.Write(b[:])
	binary.LittleEndian.PutUint32(b, uint32(node.PulseNum))
	hash.Write(b[:4])
	b[0] = byte(node.State)
	hash.Write(b[:1])
	hash.Write(node.PublicKey)
	return hash.Sum(nil), nil
}

func calculateHash(list []*core.ActiveNode) ([]byte, error) {
	hash := sha3.New224()
	sort.Slice(list[:], func(i, j int) bool {
		return bytes.Compare(list[i].NodeID[:], list[j].NodeID[:]) < 0
	})
	for _, node := range list {
		nodeHash, err := calculateNodeHash(node)
		if err != nil {
			return nil, errors.Wrap(err, "error calculating hash")
		}
		hash.Write(nodeHash)
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

func (nk *nodekeeper) Sync(approved bool) {
	// TODO: unsync discard logic due to timeout
	// TODO: sync logic
	// TODO: pass and store pulsenumber
}

func (nk *nodekeeper) AddUnsync(node *core.ActiveNode) {
	nk.unsyncLock.Lock()
	defer nk.unsyncLock.Unlock()

	// TODO: check node pulsenumber
	nk.unsync = append(nk.unsync, node)
}

func (nk *nodekeeper) AddUnsyncGossip(nodes []*core.ActiveNode) {
	nk.unsyncLock.Lock()
	defer nk.unsyncLock.Unlock()

	// TODO: check nodes pulsenumber
	for _, node := range nodes {
		nk.unsyncGossip[node.NodeID] = node
	}
}
