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

package nodenetwork

import (
	"errors"
	"sync"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
)

type UnsyncList struct {
	unsync []core.Node
	pulse  core.PulseNumber
	hash   []*network.NodeUnsyncHash

	waiters     []chan []*network.NodeUnsyncHash
	waitersLock sync.Mutex

	unsyncListCache map[core.RecordRef][]core.Node
	unsyncListLock  sync.Mutex
	unsyncHashCache map[core.RecordRef][]*network.NodeUnsyncHash
	unsyncHashLock  sync.Mutex
}

// NewUnsyncHolder create new object to hold data for consensus
func NewUnsyncHolder(pulse core.PulseNumber, unsync []core.Node) *UnsyncList {
	return &UnsyncList{
		pulse:           pulse,
		unsync:          unsync,
		waiters:         make([]chan []*network.NodeUnsyncHash, 0),
		unsyncListCache: make(map[core.RecordRef][]core.Node),
		unsyncHashCache: make(map[core.RecordRef][]*network.NodeUnsyncHash),
	}
}

// GetUnsync returns list of local unsync nodes. This list is created
func (u *UnsyncList) GetUnsync() []core.Node {
	return u.unsync
}

// GetPulse returns actual pulse for current consensus process.
func (u *UnsyncList) GetPulse() core.PulseNumber {
	return u.pulse
}

// SetHash sets hash of unsync lists for each node of consensus.
func (u *UnsyncList) SetHash(hash []*network.NodeUnsyncHash) {
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
	u.waiters = make([]chan []*network.NodeUnsyncHash, 0)
}

// GetHash get hash of unsync lists for each node of consensus. If hash is not calculated yet, then this call blocks
// until the hash is calculated with SetHash() call
func (u *UnsyncList) GetHash(blockTimeout time.Duration) ([]*network.NodeUnsyncHash, error) {
	u.waitersLock.Lock()
	if u.hash != nil {
		result := u.hash
		u.waitersLock.Unlock()
		return result, nil
	}
	ch := make(chan []*network.NodeUnsyncHash, 1)
	u.waiters = append(u.waiters, ch)
	u.waitersLock.Unlock()
	var result []*network.NodeUnsyncHash
	select {
	case data := <-ch:
		if data == nil {
			return nil, errors.New("GetHash: channel closed")
		}
		result = data
	case <-time.After(blockTimeout):
		return nil, errors.New("GetHash: timeout")
	}
	return result, nil
}

// AddUnsyncList add unsync list for remote ref
func (u *UnsyncList) AddUnsyncList(ref core.RecordRef, unsync []core.Node) {
	u.unsyncListLock.Lock()
	defer u.unsyncListLock.Unlock()

	u.unsyncListCache[ref] = unsync
}

// AddUnsyncHash add unsync hash for remote ref
func (u *UnsyncList) AddUnsyncHash(ref core.RecordRef, hash []*network.NodeUnsyncHash) {
	u.unsyncHashLock.Lock()
	defer u.unsyncHashLock.Unlock()

	u.unsyncHashCache[ref] = hash
}

// GetUnsyncList get unsync list for remote ref
func (u *UnsyncList) GetUnsyncList(ref core.RecordRef) ([]core.Node, bool) {
	u.unsyncListLock.Lock()
	defer u.unsyncListLock.Unlock()

	result, ok := u.unsyncListCache[ref]
	return result, ok
}

// GetUnsyncHash get unsync hash for remote ref
func (u *UnsyncList) GetUnsyncHash(ref core.RecordRef) ([]*network.NodeUnsyncHash, bool) {
	u.unsyncHashLock.Lock()
	defer u.unsyncHashLock.Unlock()

	result, ok := u.unsyncHashCache[ref]
	return result, ok
}
