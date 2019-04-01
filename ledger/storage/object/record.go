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
	"context"
	"sync"

	"github.com/insolar/insolar/insolar/record"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/ledger/storage/db"
	"go.opencensus.io/stats"
)

//go:generate go run gen/type.go

// TypeID encodes a record object type.
type TypeID uint32

// TypeIDSize is a size of TypeID type.
const TypeIDSize = 4

func init() {
	// ID can be any unique int value.
	// Never change id constants. They are used for serialization.
	register(100, new(GenesisRecord))
	register(101, new(ChildRecord))

	register(200, new(RequestRecord))

	register(300, new(ResultRecord))
	register(301, new(TypeRecord))
	register(302, new(CodeRecord))
	register(303, new(ActivateRecord))
	register(304, new(AmendRecord))
	register(305, new(DeactivationRecord))
}

//go:generate minimock -i github.com/insolar/insolar/ledger/storage/object.RecordAccessor -o ./ -s _mock.go

// RecordAccessor provides info about record-values from storage.
type RecordAccessor interface {
	// ForID returns record for provided id.
	ForID(ctx context.Context, id insolar.ID) (record.MaterialRecord, error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/storage/object.RecordModifier -o ./ -s _mock.go

// RecordModifier provides methods for setting record-values to storage.
type RecordModifier interface {
	// Set saves new record-value in storage.
	Set(ctx context.Context, id insolar.ID, rec record.MaterialRecord) error
}

// RecordMemory is an in-memory struct for record-storage.
type RecordMemory struct {
	jetIndex db.JetIndexModifier

	lock   sync.RWMutex
	memory map[insolar.ID]record.MaterialRecord
}

// NewRecordMemory creates a new instance of RecordMemory storage.
func NewRecordMemory() *RecordMemory {
	return &RecordMemory{
		memory:   map[insolar.ID]record.MaterialRecord{},
		jetIndex: db.NewJetIndex(),
	}
}

// Set saves new record-value in storage.
func (m *RecordMemory) Set(ctx context.Context, id insolar.ID, rec record.MaterialRecord) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	_, ok := m.memory[id]
	if ok {
		return ErrOverride
	}

	m.memory[id] = rec
	m.jetIndex.Add(id, rec.JetID)

	stats.Record(ctx,
		statIndexInMemoryCount.M(1),
	)

	return nil
}

// ForID returns record for provided id.
func (m *RecordMemory) ForID(ctx context.Context, id insolar.ID) (rec record.MaterialRecord, err error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	rec, ok := m.memory[id]
	if !ok {
		err = ErrNotFound
		return
	}

	return
}

// RecordDB is a DB storage implementation. It saves records to disk and does not allow removal.
type RecordDB struct {
	DB   db.DB
	lock sync.RWMutex
}

type recordKey insolar.ID

func (k recordKey) Scope() db.Scope {
	return db.ScopeRecord
}

func (k recordKey) ID() []byte {
	res := insolar.ID(k)
	return (&res).Bytes()
}

// NewRecordDB creates new DB storage instance.
func NewRecordDB(db db.DB) *RecordDB {
	return &RecordDB{
		DB: db,
	}
}

// Set saves new record-value in storage.
func (r *RecordDB) Set(ctx context.Context, id insolar.ID, rec record.MaterialRecord) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.set(id, rec)
}

// ForID returns record for provided id.
func (r *RecordDB) ForID(ctx context.Context, id insolar.ID) (record.MaterialRecord, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.get(id)
}

func (r *RecordDB) set(id insolar.ID, rec record.MaterialRecord) error {
	key := recordKey(id)

	_, err := r.DB.Get(key)
	if err == nil {
		return ErrOverride
	}

	return r.DB.Set(key, EncodeRecord(rec))
}

func (r *RecordDB) get(id insolar.ID) (rec record.MaterialRecord, err error) {
	buff, err := r.DB.Get(recordKey(id))
	if err == db.ErrNotFound {
		err = ErrNotFound
		return
	}
	if err != nil {
		return
	}
	rec = DecodeRecord(buff)
	return
}

func EncodeRecord(rec record.MaterialRecord) []byte {
	buff := SerializeRecord(rec.Record)
	result := append(buff[:], rec.JetID[:]...)

	return result
}

func DecodeRecord(buff []byte) record.MaterialRecord {
	recBuff := buff[:len(buff)-insolar.RecordIDSize]
	jetIDBuff := buff[len(buff)-insolar.RecordIDSize:]

	rec := DeserializeRecord(recBuff)

	var jetID insolar.JetID
	copy(jetID[:], jetIDBuff)

	return record.MaterialRecord{Record: rec, JetID: jetID}
}
