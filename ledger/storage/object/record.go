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

package object

import (
	"context"
	"io"
	"sync"

	"github.com/insolar/insolar"
	"github.com/insolar/insolar/ledger/storage/db"
	"go.opencensus.io/stats"
)

//go:generate go run gen/type.go

// TypeID encodes a record object type.
type TypeID uint32

// TypeIDSize is a size of TypeID type.
const TypeIDSize = 4

// VirtualRecord is base interface for all records.
type VirtualRecord interface {
	// WriteHashData writes record data to provided writer. This data is used to calculate record's hash.
	WriteHashData(w io.Writer) (int, error)
}

type MaterialRecord struct {
	Record VirtualRecord

	JetID insolar.JetID
}

func init() {
	// ID can be any unique int value.
	// Never change id constants. They are used for serialization.
	register(100, new(GenesisRecord))
	register(101, new(ChildRecord))
	register(102, new(JetRecord))

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
	ForID(ctx context.Context, id insolar.ID) (MaterialRecord, error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/storage/object.RecordModifier -o ./ -s _mock.go

// RecordModifier provides methods for setting record-values to storage.
type RecordModifier interface {
	// Set saves new record-value in storage.
	Set(ctx context.Context, id insolar.ID, rec MaterialRecord) error
}

// RecordMemory is an in-memory struct for record-storage.
type RecordMemory struct {
	jetIndex db.JetIndexModifier

	lock   sync.RWMutex
	memory map[insolar.ID]MaterialRecord
}

// NewRecordMemory creates a new instance of RecordMemory storage.
func NewRecordMemory() *RecordMemory {
	return &RecordMemory{
		memory:   map[insolar.ID]MaterialRecord{},
		jetIndex: db.NewJetIndex(),
	}
}

// Set saves new record-value in storage.
func (m *RecordMemory) Set(ctx context.Context, id insolar.ID, rec MaterialRecord) error {
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
func (m *RecordMemory) ForID(ctx context.Context, id insolar.ID) (rec MaterialRecord, err error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	rec, ok := m.memory[id]
	if !ok {
		err = RecNotFound
		return
	}

	return
}
