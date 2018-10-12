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

package storage

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/record"
)

type mucount struct {
	*sync.Mutex
	count int32
}

// IDLocker provides Lock/Unlock methods per record ID.
//
// TODO: for further optimization we could use sync.Pool for mutexes.
type IDLocker struct {
	m    map[core.RecordID]*mucount
	rwmu sync.RWMutex
}

// NewIDLocker creates new initialized IDLocker.
func NewIDLocker() *IDLocker {
	return &IDLocker{
		m: make(map[core.RecordID]*mucount),
	}
}

func (l *IDLocker) getmu(cid *core.RecordID) *mucount {
	var (
		mc *mucount
		ok bool
	)
	l.rwmu.RLock()
	mc, ok = l.m[*cid]
	l.rwmu.RUnlock()
	if ok {
		return mc
	}
	// initialize mutex for recordID
	l.rwmu.Lock()
	// check if not already initialized before Lock has been acquired
	if mc, ok = l.m[*cid]; !ok {
		mc = &mucount{Mutex: &sync.Mutex{}}
		l.m[*cid] = mc
	}
	l.rwmu.Unlock()
	return mc
}

// Lock locks mutex belonged to record ID.
// If mutex does not exist, it will be created in concurrent safe fashion.
func (l *IDLocker) Lock(id *record.ID) {
	cid := id.CoreID()
	mc := l.getmu(cid)
	atomic.AddInt32(&mc.count, 1)
	mc.Lock()
}

// Unlock unlocks mutex belonged to record ID.
func (l *IDLocker) Unlock(id *record.ID) {
	cid := id.CoreID()

	l.rwmu.Lock()
	defer l.rwmu.Unlock()
	mc, ok := l.m[*cid]
	if !ok {
		panic(fmt.Sprintf("try to unlock not initialized mutex for ID %+v", cid))
	}

	cnt := atomic.AddInt32(&mc.count, -1)
	mc.Unlock()
	if cnt == 0 {
		delete(l.m, *cid)
	}
}
