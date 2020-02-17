// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package object

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/store"
)

// TypeID encodes a record object type.
type TypeID uint32

// TypeIDSize is a size of TypeID type.
const TypeIDSize = 4

//go:generate minimock -i github.com/insolar/insolar/ledger/object.RecordStorage -o ./ -s _mock.go -g

// RecordStorage is an union of RecordAccessor and RecordModifier
type RecordStorage interface {
	RecordAccessor
	RecordModifier
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.AtomicRecordStorage -o ./ -s _mock.go -g

// AtomicRecordStorage is an union of RecordAccessor and AtomicRecordModifier
type AtomicRecordStorage interface {
	RecordAccessor
	AtomicRecordModifier
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.RecordAccessor -o ./ -s _mock.go -g

// RecordAccessor provides info about record-values from storage.
type RecordAccessor interface {
	// ForID returns record for provided id.
	ForID(ctx context.Context, id insolar.ID) (record.Material, error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.RecordCollectionAccessor -o ./ -s _mock.go -g

// RecordCollectionAccessor provides methods for querying records with specific search conditions.
type RecordCollectionAccessor interface {
	// ForPulse returns []MaterialRecord for a provided jetID and a pulse number.
	ForPulse(ctx context.Context, jetID insolar.JetID, pn insolar.PulseNumber) []record.Material
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.RecordModifier -o ./ -s _mock.go -g

// RecordModifier provides methods for setting record-values to storage.
type RecordModifier interface {
	// Set saves new record-value in storage with order-processing.
	Set(ctx context.Context, rec record.Material) error
	// BatchSet saves a batch of records to storage with order-processing.
	BatchSet(ctx context.Context, recs []record.Material) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.AtomicRecordModifier -o ./ -s _mock.go

// AtomicRecordModifier allows to modify multiple record atomically.
type AtomicRecordModifier interface {
	// SetAtomic atomically stores records to storage. Guarantees to either store all records or none.
	SetAtomic(ctx context.Context, records ...record.Material) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.RecordCleaner -o ./ -s _mock.go -g

// RecordCleaner provides an interface for removing records from a storage.
type RecordCleaner interface {
	// DeleteForPN method removes records from a storage for a pulse
	DeleteForPN(ctx context.Context, pulse insolar.PulseNumber)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.RecordPositionAccessor -o ./ -s _mock.go

// RecordPositionAccessor provides an interface for fetcing position of records.
type RecordPositionAccessor interface {
	// LastKnownPosition returns last known position of record in Pulse.
	LastKnownPosition(pn insolar.PulseNumber) (uint32, error)
	// AtPosition returns record ID for a specific pulse and a position
	AtPosition(pn insolar.PulseNumber, position uint32) (insolar.ID, error)
}

// RecordMemory is an in-indexStorage struct for record-storage.
type RecordMemory struct {
	jetIndex         store.JetIndexModifier
	jetIndexAccessor store.JetIndexAccessor

	lock     sync.RWMutex
	recsStor map[insolar.ID]record.Material
}

// NewRecordMemory creates a new instance of RecordMemory storage.
func NewRecordMemory() *RecordMemory {
	ji := store.NewJetIndex()
	return &RecordMemory{
		recsStor:         map[insolar.ID]record.Material{},
		jetIndex:         ji,
		jetIndexAccessor: ji,
	}
}

// SetAtomic atomically stores records to storage. Guarantees to either store all records or none.
func (m *RecordMemory) SetAtomic(ctx context.Context, recs ...record.Material) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, r := range recs {
		if r.ID.IsEmpty() {
			return errors.New("id is empty")
		}
		_, ok := m.recsStor[r.ID]
		if ok {
			return ErrOverride
		}
	}

	for _, r := range recs {
		m.recsStor[r.ID] = r
		m.jetIndex.Add(r.ID, r.JetID)
	}

	stats.Record(ctx,
		statRecordInMemoryAddedCount.M(int64(len(recs))),
	)
	return nil
}

// ForID returns record for provided id.
func (m *RecordMemory) ForID(ctx context.Context, id insolar.ID) (rec record.Material, err error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	rec, ok := m.recsStor[id]
	if !ok {
		err = ErrNotFound
		return
	}

	return
}

// ForPulse returns []MaterialRecord for a provided jetID and a pulse number.
func (m *RecordMemory) ForPulse(
	ctx context.Context, jetID insolar.JetID, pn insolar.PulseNumber,
) []record.Material {
	m.lock.RLock()
	defer m.lock.RUnlock()

	ids := m.jetIndexAccessor.For(jetID)
	var res []record.Material
	for id := range ids {
		if id.Pulse() == pn {
			rec := m.recsStor[id]
			res = append(res, rec)
		}
	}

	return res
}

// DeleteForPN method removes records from a storage for all pulses until pulse (pulse included)
func (m *RecordMemory) DeleteForPN(ctx context.Context, pulse insolar.PulseNumber) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for id, rec := range m.recsStor {
		if id.Pulse() != pulse {
			continue
		}

		m.jetIndex.Delete(id, rec.JetID)
		delete(m.recsStor, id)

		stats.Record(ctx,
			statRecordInMemoryRemovedCount.M(1),
		)
	}
}
