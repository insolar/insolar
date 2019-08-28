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

package object

import (
	"fmt"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/pulse"
)

type mucount struct {
	*sync.RWMutex
	count int32
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexLocker -o ./ -s _mock.go -g

// IndexLocker provides Lock/Unlock methods per record ID.
type IndexLocker interface {
	Lock(id insolar.ID)
	Unlock(id insolar.ID)
}

// IndexLocker provides Lock/Unlock methods per record ID.
type idLocker struct {
	mu   sync.Mutex
	muxs map[insolar.ID]*mucount
}

// NewIndexLocker creates new initialized IndexLocker.
func NewIndexLocker() IndexLocker {
	return &idLocker{
		muxs: make(map[insolar.ID]*mucount),
	}
}

// Lock locks mutex belonged to record ID.
// If mutex does not exist, it will be created in concurrent safe fashion.
func (l *idLocker) Lock(id insolar.ID) {
	// Reset pulse. It should not be considered when locking.
	normalizedID := *insolar.NewID(pulse.LocalRelative, id.Hash())

	l.mu.Lock()
	mc, ok := l.muxs[normalizedID]
	if !ok {
		mc = &mucount{RWMutex: &sync.RWMutex{}}
		l.muxs[normalizedID] = mc
	}
	mc.count++
	l.mu.Unlock()

	mc.Lock()
}

// Unlock unlocks mutex belonged to record ID.
func (l *idLocker) Unlock(id insolar.ID) {
	// Reset pulse. It should not be considered when locking.
	zeroID := *insolar.NewID(pulse.LocalRelative, id.Hash())

	l.mu.Lock()
	defer l.mu.Unlock()

	mc, ok := l.muxs[zeroID]
	if !ok {
		panic(fmt.Sprintf("try to unlock not initialized mutex for ID %+v", zeroID))
	}
	mc.count--
	mc.Unlock()
	if mc.count == 0 {
		delete(l.muxs, zeroID)
	}
}
