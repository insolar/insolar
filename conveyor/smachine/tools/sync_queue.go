///
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
///

package tools

import (
	"github.com/insolar/insolar/network/consensus/common/rwlock"
	"sync"
)

func NewSyncQueue(locker sync.Locker) SyncQueue {
	if locker == nil {
		panic("illegal value")
	}
	return SyncQueue{locker: locker}
}

func NewSignalCondQueue(signal *sync.Cond) SyncQueue {
	if signal == nil {
		panic("illegal value")
	}
	return SyncQueue{locker: signal.L, signalFn: signal.Broadcast}
}

func NewSignalFuncQueue(locker sync.Locker, signalFn func()) SyncQueue {
	if locker == nil {
		panic("illegal value")
	}
	return SyncQueue{locker: locker, signalFn: signalFn}
}

func NewNoSyncQueue() SyncQueue {
	return SyncQueue{locker: rwlock.DummyLocker()}
}

type SyncFunc func()
type SyncFuncList []SyncFunc

type SyncQueue struct {
	locker   sync.Locker
	signalFn func()
	queue    SyncFuncList
}

func (p *SyncQueue) IsZero() bool {
	return p.locker == nil
}

func (p *SyncQueue) Add(fn SyncFunc) {
	if fn == nil {
		panic("illegal value")
	}
	p.locker.Lock()
	defer p.locker.Unlock()

	p.queue = append(p.queue, fn)
	if p.signalFn != nil {
		p.signalFn()
	}
}

func (p *SyncQueue) Flush() []SyncFunc {
	p.locker.Lock()
	defer p.locker.Unlock()

	if len(p.queue) == 0 {
		return nil
	}

	nextCap := cap(p.queue)
	if nextCap > 128 && len(p.queue)<<1 < nextCap {
		nextCap >>= 1
	}
	queue := p.queue
	p.queue = make([]SyncFunc, 0, nextCap)

	return queue
}
