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

	"github.com/insolar/insolar/core"
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
	rwmu sync.Mutex
}

// NewIDLocker creates new initialized IDLocker.
func NewIDLocker() *IDLocker {
	return &IDLocker{
		m: make(map[core.RecordID]*mucount),
	}
}

// Lock locks mutex belonged to record ID.
// If mutex does not exist, it will be created in concurrent safe fashion.
func (l *IDLocker) Lock(id *core.RecordID) {
	l.rwmu.Lock()
	mc, ok := l.m[*id]
	if !ok {
		mc = &mucount{Mutex: &sync.Mutex{}}
		l.m[*id] = mc
	}
	mc.count++
	l.rwmu.Unlock()

	mc.Lock()
}

// Unlock unlocks mutex belonged to record ID.
func (l *IDLocker) Unlock(id *core.RecordID) {
	l.rwmu.Lock()
	defer l.rwmu.Unlock()

	mc, ok := l.m[*id]
	if !ok {
		panic(fmt.Sprintf("try to unlock not initialized mutex for ID %+v", id))
	}
	mc.count--
	mc.Unlock()
	if mc.count == 0 {
		delete(l.m, *id)
	}
}
