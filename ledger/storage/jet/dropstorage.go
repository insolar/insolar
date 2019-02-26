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
	"github.com/insolar/insolar/ledger/storage/jet"
)

type DropSaver interface {
	Set(ctx context.Context, jetID core.JetID, drop JetDrop, pulse core.PulseNumber) error
}

type DropFetcher interface {
	ForPulse(ctx context.Context, jetID core.JetID, pulse core.PulseNumber) (JetDrop, error)
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

func (m *dropForPulseManager) forPulse(jetID core.JetID, pulse core.PulseNumber) (JetDrop, error) {
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
	jets map[core.JetID]*dropForPulseManager
}

func (m *dropStorageMemory) fetchStorage(jetID core.JetID) (ds *dropForPulseManager) {
	m.lock.Lock()
	defer m.lock.Unlock()

	ds, ok := m.jets[jetID]
	if !ok {
		m.jets[jetID] = new(dropForPulseManager)
		ds = m.jets[jetID]
	}
	return
}

func (m *dropStorageMemory) ForPulse(ctx context.Context, jetID core.JetID, pulse core.PulseNumber) (JetDrop, error) {
	ds := m.fetchStorage(jetID)
	return ds.forPulse(jetID, pulse)
}

func (m *dropStorageMemory) Set(ctx context.Context, jetID core.JetID, drop JetDrop, pulse core.PulseNumber) error {
	ds := m.fetchStorage(jetID)
	return ds.set(drop, pulse)
}

type dropStorageDB struct {
	DB storage.DBContext
}

func (*dropStorageDB) ForPulse(ctx context.Context, jetID core.JetID, pulse core.PulseNumber) (JetDrop, error) {
	panic("implement me")
}

func (ds *dropStorageDB) Set(ctx context.Context, jetID core.JetID, drop JetDrop, pulse core.PulseNumber) error {
	_, prefix := Jet(jetID)
	k := storage.Prefixkey(scopeIDJetDrop, prefix, drop.Pulse.Bytes())
	_, err := ds.DB.get(ctx, k)
	if err == nil {
		return ErrOverride
	}

	encoded, err := jet.Encode(drop)
	if err != nil {
		return err
	}
	return ds.DB.set(ctx, k, encoded)
}
