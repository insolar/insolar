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

package jet

import (
	"context"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage"
)

// DropModifier provides interface for modifying jetdrops
//go:generate minimock -i github.com/insolar/insolar/ledger/storage/jet.DropModifier -o ./ -s _mock.go
type DropModifier interface {
	Set(ctx context.Context, jetID storage.JetID, drop JetDrop) error
}

// DropAccessor provides interface for accessing jetdrops
//go:generate minimock -i github.com/insolar/insolar/ledger/storage/jet.DropAccessor -o ./ -s _mock.go
type DropAccessor interface {
	ForPulse(ctx context.Context, jetID storage.JetID, pulse core.PulseNumber) (JetDrop, error)
}

type dropForPulseManager struct {
	lock  sync.RWMutex
	drops map[core.PulseNumber]JetDrop
}

func (m *dropForPulseManager) set(drop JetDrop, pulse core.PulseNumber) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.drops[pulse]; ok {
		return storage.ErrOverride
	}

	m.drops[pulse] = drop

	return nil
}

func (m *dropForPulseManager) forPulse(jetID storage.JetID, pulse core.PulseNumber) (JetDrop, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	drop, ok := m.drops[pulse]
	if !ok {
		return JetDrop{}, core.ErrNotFound
	}

	return drop, nil
}

type dropStorageMemory struct {
	lock sync.Mutex
	jets map[storage.JetID]*dropForPulseManager
}

func NewDropStorageMemory() *dropStorageMemory {
	return &dropStorageMemory{
		jets: map[storage.JetID]*dropForPulseManager{},
	}
}

func (m *dropStorageMemory) fetchStorage(jetID storage.JetID) (ds *dropForPulseManager) {
	m.lock.Lock()
	defer m.lock.Unlock()

	ds, ok := m.jets[jetID]
	if !ok {
		m.jets[jetID] = new(dropForPulseManager)
		ds = m.jets[jetID]
	}
	return
}

func (m *dropStorageMemory) ForPulse(ctx context.Context, jetID storage.JetID, pulse core.PulseNumber) (JetDrop, error) {
	ds := m.fetchStorage(jetID)
	return ds.forPulse(jetID, pulse)
}

func (m *dropStorageMemory) Set(ctx context.Context, jetID storage.JetID, drop JetDrop) error {
	ds := m.fetchStorage(jetID)
	return ds.set(drop, drop.Pulse)
}

type dropStorageDB struct {
	DB storage.DBContext `inject:""`
}

func NewDropStorageDB() *dropStorageDB {
	return &dropStorageDB{}
}

func (ds *dropStorageDB) ForPulse(ctx context.Context, jetID storage.JetID, pulse core.PulseNumber) (JetDrop, error) {
	_, prefix := jetID.Jet()
	k := storage.JetDropPrefixKey(prefix, pulse)

	// buf, err := db.get(ctx, k)
	buf, err := ds.DB.Get(ctx, k)
	if err != nil {
		return JetDrop{}, err
	}
	drop, err := Decode(buf)
	if err != nil {
		return JetDrop{}, err
	}
	return *drop, nil
}

func (ds *dropStorageDB) Set(ctx context.Context, jetID storage.JetID, drop JetDrop) error {
	_, prefix := jetID.Jet()
	k := storage.JetDropPrefixKey(prefix, drop.Pulse)
	_, err := ds.DB.Get(ctx, k)
	if err == nil {
		return storage.ErrOverride
	}

	encoded, err := Encode(&drop)
	if err != nil {
		return err
	}
	return ds.DB.Set(ctx, k, encoded)
}
