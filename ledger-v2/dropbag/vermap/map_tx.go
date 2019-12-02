//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package vermap

import (
	"github.com/insolar/insolar/longbits"
	"math"
	"sync"
	"sync/atomic"
)

var _ TxMap = &Tx{}

const pendingTxVersion = uint64(math.MaxUint64)
const discardedTxVersion = uint64(0)
const firstTxVersion = uint64(1)

type Tx struct {
	container   *IncrementalMap
	viewVersion uint64
	txMark      *txMark
	options     Options

	mutex   sync.RWMutex
	pending map[longbits.ByteString]txEntry
	//deleted map[longbits.ByteString]struct{}
}

func (t *Tx) canReplace() bool {
	return t.options&CanReplace != 0
}

func (t *Tx) canDelete() bool {
	return t.options&CanDelete != 0
}

func (t *Tx) canUpdate() (canUpdate bool, isDiscarded bool) {
	v := t.txMark.getVersion()
	return v == pendingTxVersion && t.options&CanUpdate != 0, v == discardedTxVersion
}

func (t *Tx) isPending() bool {
	v := t.txMark.getVersion()
	return v == pendingTxVersion
}

func (t *Tx) ensureUpdate() error {
	switch canUpdate, isDiscarded := t.canUpdate(); {
	case isDiscarded:
		return ErrDiscardedTxn
	case !canUpdate:
		return ErrReadOnlyTxn
	}
	return nil
}

func (t *Tx) ensureView() error {
	if _, isDiscarded := t.canUpdate(); isDiscarded {
		return ErrDiscardedTxn
	}
	return nil
}

func (t *Tx) Set(k Key, v Value) error {
	return t.SetEntry(NewEntry(k, v))
}

func (t *Tx) SetEntry(kv Entry) error {
	if err := t.container.validateEntry(kv); err != nil {
		return err
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	if err := t.ensureUpdate(); err != nil {
		return err
	}

	if !t.canReplace() && t.contains(kv.Key) {
		return ErrExistingKey
	}

	if t.pending == nil {
		t.pending = make(map[longbits.ByteString]txEntry)
	}
	t.pending[kv.Key] = txEntry{kv, t.txMark}
	//delete(t.deleted, kv.Key)
	return nil
}

//func (t *Tx) Delete(k Key) error {
//	if err := t.container.validateKey(k); err != nil {
//		return err
//	}
//	if !t.canDelete() {
//		return ErrExistingKey
//	}
//
//	t.mutex.Lock()
//	defer t.mutex.Unlock()
//
//	if err := t.ensureUpdate(); err != nil {
//		return err
//	}
//
//	delete(t.pending, k)
//
//	if _, ok := t.txMarkOf(k); ok {
//		if t.deleted == nil {
//			t.deleted = make(map[longbits.ByteString]struct{})
//		}
//		t.deleted[k] = struct{}{}
//	}
//	return nil
//}

func (t *Tx) GetEntry(k Key) (Entry, error) {
	if err := t.container.validateKey(k); err != nil {
		return Entry{}, err
	}

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if err := t.ensureView(); err != nil {
		return Entry{}, err
	}

	if v, ok := t.pending[k]; ok {
		return v.Entry, nil
	}

	return t.container.getAsOf(k, t.viewVersion)
}

func (t *Tx) Get(k Key) (Value, error) {
	if kv, err := t.GetEntry(k); err != nil {
		return nil, err
	} else {
		return kv.Value, nil
	}
}

func (t *Tx) Contains(k Key) bool {
	if err := t.container.validateKey(k); err != nil {
		return false
	}

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.contains(k)
}

func (t *Tx) contains(k Key) bool {
	if _, ok := t.pending[k]; ok {
		return true
	}
	_, ok := t.container.markAsOf(k, t.viewVersion)
	return ok
}

func (t *Tx) txMarkOf(k Key) (*txMark, bool) {
	return t.container.markAsOf(k, t.viewVersion)
}

func (t *Tx) Discard() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	_ = t.txMark.discard()
}

func (t *Tx) Commit() error {
	return t.container.commitTx(t)
}

// a minimal tx portion to avoid retention of whole Tx in memory after commit
type txMark struct {
	commitVersion uint64 //atomic
}

func (t *txMark) getVersion() uint64 {
	return atomic.LoadUint64(&t.commitVersion)
}

func (t *txMark) getVersionNil() uint64 {
	if t == nil {
		return firstTxVersion
	}
	return t.getVersion()
}

func (t *txMark) discard() error {
	for {
		v := t.getVersion()
		if v != pendingTxVersion {
			return ErrDiscardedTxn
		}
		if atomic.CompareAndSwapUint64(&t.commitVersion, v, discardedTxVersion) {
			return nil
		}
	}
}

func (t *txMark) setVersion(version uint64) bool {
	switch version {
	case pendingTxVersion, discardedTxVersion:
		panic("illegal value")
	}
	return atomic.CompareAndSwapUint64(&t.commitVersion, pendingTxVersion, version)
}
