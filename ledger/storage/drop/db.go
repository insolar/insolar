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
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/jet"
)

// Modifier provides interface for modifying jetdrops
//go:generate minimock -i github.com/insolar/insolar/ledger/storage/drop.Modifier -o ./ -s _mock.go
type Modifier interface {
	Set(ctx context.Context, jetID core.JetID, drop jet.Drop) error
}

// Accessor provides interface for accessing jetdrops
//go:generate minimock -i github.com/insolar/insolar/ledger/storage/drop.Accessor -o ./ -s _mock.go
type Accessor interface {
	ForPulse(ctx context.Context, jetID core.JetID, pulse core.PulseNumber) (jet.Drop, error)
}

type dropKey struct {
	pulse core.PulseNumber
	jetID core.JetID
}

type dropStorageMemory struct {
	lock sync.RWMutex
	jets map[dropKey]jet.Drop
}

// NewStorageMemory creates new storage, that holds data in-memory
func NewStorageMemory() *dropStorageMemory { // nolint: golint
	return &dropStorageMemory{
		jets: map[dropKey]jet.Drop{},
	}
}

// ForPulse returns a jet.Drop for a provided pulse, that is stored in memory
func (m *dropStorageMemory) ForPulse(ctx context.Context, jetID core.JetID, pulse core.PulseNumber) (jet.Drop, error) {
	m.lock.RLock()
	m.lock.RUnlock()

	key := dropKey{jetID: jetID, pulse: pulse}
	d, ok := m.jets[key]
	if !ok {
		return jet.Drop{}, core.ErrNotFound
	}

	return d, nil
}

// Set saves a provided jet.Drop to memory
func (m *dropStorageMemory) Set(ctx context.Context, jetID core.JetID, drop jet.Drop) error {
	m.lock.RLock()
	m.lock.RUnlock()

	key := dropKey{jetID: jetID, pulse: drop.Pulse}
	m.jets[key] = drop

	return nil
}

type dropStorageDB struct {
	DB storage.DBContext `inject:""`
}

// NewStorageDB creates new storage, that holds data in db
func NewStorageDB() *dropStorageDB { // nolint: golint
	return &dropStorageDB{}
}

// ForPulse returns a jet.Drop for a provided pulse, that is stored in db
func (ds *dropStorageDB) ForPulse(ctx context.Context, jetID core.JetID, pulse core.PulseNumber) (jet.Drop, error) {
	_, prefix := jetID.Jet()
	k := storage.JetDropPrefixKey(prefix, pulse)

	// buf, err := db.get(ctx, k)
	buf, err := ds.DB.Get(ctx, k)
	if err != nil {
		return jet.Drop{}, err
	}
	drop, err := jet.Decode(buf)
	if err != nil {
		return jet.Drop{}, err
	}
	return *drop, nil
}

// Set saves a provided jet.Drop to db
func (ds *dropStorageDB) Set(ctx context.Context, jetID core.JetID, drop jet.Drop) error {
	_, prefix := jetID.Jet()
	k := storage.JetDropPrefixKey(prefix, drop.Pulse)
	_, err := ds.DB.Get(ctx, k)
	if err == nil {
		return storage.ErrOverride
	}

	encoded, err := jet.Encode(&drop)
	if err != nil {
		return err
	}
	return ds.DB.Set(ctx, k, encoded)
}
