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
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/consensus"
)

type UnsyncList struct {
	unsync []*core.ActiveNode
	pulse  core.PulseNumber
	hash   []*consensus.NodeUnsyncHash

	waiters     []chan []*consensus.NodeUnsyncHash
	waitersLock sync.Mutex

	unsyncListCache map[core.RecordRef][]*core.ActiveNode
	unsyncListLock  sync.Mutex
	unsyncHashCache map[core.RecordRef][]*consensus.NodeUnsyncHash
	unsyncHashLock  sync.Mutex
}

// NewUnsyncHolder create new object to hold data for consensus
func NewUnsyncHolder(pulse core.PulseNumber, unsync []*core.ActiveNode) *UnsyncList {
	return &UnsyncList{pulse: pulse, unsync: unsync}
}

// GetUnsync returns list of local unsync nodes. This list is created
func (u *UnsyncList) GetUnsync() []*core.ActiveNode {
	return u.unsync
}

// GetPulse returns actual pulse for current consensus process.
func (u *UnsyncList) GetPulse() core.PulseNumber {
	return u.pulse
}

// SetHash sets hash of unsync lists for each node of consensus.
func (u *UnsyncList) SetHash(hash []*consensus.NodeUnsyncHash) {
	u.waitersLock.Lock()
	defer u.waitersLock.Unlock()

	u.hash = hash
	if len(u.waiters) == 0 {
		return
	}
	for _, ch := range u.waiters {
		ch <- u.hash
		close(ch)
	}
	u.waiters = make([]chan []*consensus.NodeUnsyncHash, 0)
}

// GetHash get hash of unsync lists for each node of consensus. If hash is not calculated yet, then this call blocks
// until the hash is calculated with SetHash() call
func (u *UnsyncList) GetHash(blockTimeout time.Duration) ([]*consensus.NodeUnsyncHash, error) {
	u.waitersLock.Lock()
	if u.hash != nil {
		result := u.hash
		u.waitersLock.Unlock()
		return result, nil
	}
	ch := make(chan []*consensus.NodeUnsyncHash)
	u.waiters = append(u.waiters, ch)
	u.waitersLock.Unlock()
	// TODO: timeout
	result := <-ch
	return result, nil
}

// AddUnsyncList add unsync list for remote ref
func (u *UnsyncList) AddUnsyncList(ref core.RecordRef, unsync []*core.ActiveNode) {
	u.unsyncListLock.Lock()
	defer u.unsyncListLock.Unlock()

	u.unsyncListCache[ref] = unsync
}

// AddUnsyncHash add unsync hash for remote ref
func (u *UnsyncList) AddUnsyncHash(ref core.RecordRef, hash []*consensus.NodeUnsyncHash) {
	u.unsyncHashLock.Lock()
	defer u.unsyncHashLock.Unlock()

	u.unsyncHashCache[ref] = hash
}

// GetUnsyncList get unsync list for remote ref
func (u *UnsyncList) GetUnsyncList(ref core.RecordRef) ([]*core.ActiveNode, bool) {
	u.unsyncListLock.Lock()
	defer u.unsyncListLock.Unlock()

	result, ok := u.unsyncListCache[ref]
	return result, ok
}

// GetUnsyncHash get unsync hash for remote ref
func (u *UnsyncList) GetUnsyncHash(ref core.RecordRef) ([]*consensus.NodeUnsyncHash, bool) {
	u.unsyncHashLock.Lock()
	defer u.unsyncHashLock.Unlock()

	result, ok := u.unsyncHashCache[ref]
	return result, ok
}
