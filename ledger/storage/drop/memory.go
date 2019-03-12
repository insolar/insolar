/*
 *    Copyright 2019 Insolar Technologies
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

package drop

import (
	"context"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage/db"
	"github.com/insolar/insolar/ledger/storage/jet"
)

type dropKey struct {
	pulse core.PulseNumber
	jetID core.JetID
}

type dropStorageMemory struct {
	lock sync.RWMutex
	jets map[dropKey]jet.Drop
}

// NewStorageMemory creates a new storage, that holds data in a memory.
func NewStorageMemory() *dropStorageMemory { // nolint: golint
	return &dropStorageMemory{
		jets: map[dropKey]jet.Drop{},
	}
}

// ForPulse returns a jet.Drop for a provided pulse, that is stored in a memory
func (m *dropStorageMemory) ForPulse(ctx context.Context, jetID core.JetID, pulse core.PulseNumber) (jet.Drop, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	key := dropKey{jetID: jetID, pulse: pulse}
	d, ok := m.jets[key]
	if !ok {
		return jet.Drop{}, db.ErrNotFound
	}

	return d, nil
}

// Set saves a provided jet.Drop to a memory
func (m *dropStorageMemory) Set(ctx context.Context, jetID core.JetID, drop jet.Drop) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	key := dropKey{jetID: jetID, pulse: drop.Pulse}
	m.jets[key] = drop

	return nil
}

func (m *dropStorageMemory) Delete(pulse core.PulseNumber) {
	m.lock.Lock()
	for key := range m.jets {
		if key.pulse == pulse {
			delete(m.jets, key)
		}
	}
	m.lock.Unlock()
}
