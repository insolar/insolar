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
	"bytes"
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/ugorji/go/codec"
	"go.opencensus.io/stats"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexAccessor -o ./ -s _mock.go

// IndexAccessor provides info about Index-values from storage.
type IndexAccessor interface {
	// ForID returns Index for provided id.
	ForID(ctx context.Context, id insolar.ID) (Lifeline, error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexCollectionAccessor -o ./ -s _mock.go

// IndexCollectionAccessor provides methods for querying a collection of blobs with specific search conditions.
type IndexCollectionAccessor interface {
	// ForJet returns a collection of lifelines for a provided jetID
	ForJet(ctx context.Context, jetID insolar.JetID) map[insolar.ID]LifelineMeta
	// ForPulseAndJet returns a collection of lifelines for a provided jetID and a pulse number
	ForPulseAndJet(ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID) map[insolar.ID]Lifeline
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexModifier -o ./ -s _mock.go

// IndexModifier provides methods for setting Index-values to storage.
type IndexModifier interface {
	// Set saves new Index-value in storage.
	Set(ctx context.Context, id insolar.ID, index Lifeline) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.ExtendedIndexModifier -o ./ -s _mock.go

// ExtendedIndexModifier provides methods for setting Index-values to storage.
// The main difference with IndexModifier is an opportunity to modify a state of an internal pulse-index
type ExtendedIndexModifier interface {
	// SetWithMeta saves index to the storage and sets its index and pulse number in internal indexes
	SetWithMeta(ctx context.Context, id insolar.ID, pn insolar.PulseNumber, index Lifeline) error
	// SetUsageForPulse updates an internal state of an internal pulse-index
	// Calling this method guaranties that provied pn will be used as a LastUsagePulse for an id
	SetUsageForPulse(ctx context.Context, id insolar.ID, pn insolar.PulseNumber)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexStorage -o ./ -s _mock.go

// IndexStorage is an union of IndexAccessor and IndexModifier.
type IndexStorage interface {
	IndexAccessor
	IndexModifier
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexCleaner -o ./ -s _mock.go

// IndexCleaner provides an interface for removing interfaces from a storage.
type IndexCleaner interface {
	// DeleteForPN method removes indexes from a storage for a provided
	DeleteForPN(ctx context.Context, pn insolar.PulseNumber)
}

// Lifeline represents meta information for record object.
type Lifeline struct {
	LatestState         *insolar.ID // Amend or activate record.
	LatestStateApproved *insolar.ID // State approved by VM.
	ChildPointer        *insolar.ID // Meta record about child activation.
	Parent              insolar.Reference
	Delegates           map[insolar.Reference]insolar.Reference
	State               StateID
	LatestUpdate        insolar.PulseNumber
	JetID               insolar.JetID
}

// LifelineMeta holds additional info about Lifeline
// It provides LastUsed pulse number
// That can be used for placed in a special bucket in processing structs
type LifelineMeta struct {
	Index    Lifeline
	LastUsed insolar.PulseNumber
}

// EncodeIndex converts lifeline index into binary format.
func EncodeIndex(index Lifeline) []byte {
	buff := bytes.NewBuffer(nil)
	enc := codec.NewEncoder(buff, &codec.CborHandle{})
	enc.MustEncode(index)

	return buff.Bytes()
}

// MustDecodeIndex converts byte array into lifeline index struct.
func MustDecodeIndex(buff []byte) (index Lifeline) {
	dec := codec.NewDecoderBytes(buff, &codec.CborHandle{})
	dec.MustDecode(&index)

	return
}

// DecodeIndex converts byte array into lifeline index struct.
func DecodeIndex(buff []byte) (index Lifeline, err error) {
	dec := codec.NewDecoderBytes(buff, &codec.CborHandle{})
	err = dec.Decode(&index)

	return
}

// CloneIndex returns copy of argument idx value.
func CloneIndex(idx Lifeline) Lifeline {
	if idx.LatestState != nil {
		tmp := *idx.LatestState
		idx.LatestState = &tmp
	}

	if idx.LatestStateApproved != nil {
		tmp := *idx.LatestStateApproved
		idx.LatestStateApproved = &tmp
	}

	if idx.ChildPointer != nil {
		tmp := *idx.ChildPointer
		idx.ChildPointer = &tmp
	}

	if idx.Delegates != nil {
		cp := make(map[insolar.Reference]insolar.Reference)
		for k, v := range idx.Delegates {
			cp[k] = v
		}
		idx.Delegates = cp
	} else {
		idx.Delegates = map[insolar.Reference]insolar.Reference{}
	}

	return idx
}

// IndexMemory is an in-indexStorage struct for index-storage.
type IndexMemory struct {
	jetIndexModifier store.JetIndexModifier
	jetIndexAccessor store.JetIndexAccessor
	pulseIndex       PulseIndex

	storageLock  sync.RWMutex
	indexStorage map[insolar.ID]Lifeline
}

// NewIndexMemory creates a new instance of IndexMemory storage.
func NewIndexMemory() *IndexMemory {
	idx := store.NewJetIndex()
	return &IndexMemory{
		indexStorage:     map[insolar.ID]Lifeline{},
		jetIndexModifier: idx,
		jetIndexAccessor: idx,
		pulseIndex:       NewPulseIndex(),
	}
}

// Set saves new Index-value in storage.
func (m *IndexMemory) Set(ctx context.Context, id insolar.ID, index Lifeline) error {
	m.storageLock.Lock()
	defer m.storageLock.Unlock()

	idx := CloneIndex(index)

	m.indexStorage[id] = idx
	m.jetIndexModifier.Add(id, idx.JetID)

	stats.Record(ctx,
		statIndexInMemoryAddedCount.M(1),
	)

	return nil
}

// SetWithMeta saves index to the storage and sets its index and pulse number in internal indexes
func (m *IndexMemory) SetWithMeta(ctx context.Context, id insolar.ID, pn insolar.PulseNumber, index Lifeline) error {
	m.storageLock.Lock()
	defer m.storageLock.Unlock()

	idx := CloneIndex(index)

	m.indexStorage[id] = idx
	m.jetIndexModifier.Add(id, idx.JetID)
	m.pulseIndex.Add(id, pn)

	stats.Record(ctx,
		statIndexInMemoryAddedCount.M(1),
	)

	return nil
}

// SetUsageForPulse updates an internal state of an internal pulse-index
// Calling this method guaranties that provied pn will be used as a LastUsagePulse for an id
func (m *IndexMemory) SetUsageForPulse(ctx context.Context, id insolar.ID, pn insolar.PulseNumber) {
	m.pulseIndex.Add(id, pn)
}

// ForID returns Index for provided id.
func (m *IndexMemory) ForID(ctx context.Context, id insolar.ID) (Lifeline, error) {
	m.storageLock.RLock()
	defer m.storageLock.RUnlock()
	var index Lifeline

	idx, ok := m.indexStorage[id]
	if !ok {
		return index, ErrIndexNotFound

	}

	index = CloneIndex(idx)

	return index, nil
}

// ForJet returns a collection of lifelines for a provided jetID
func (m *IndexMemory) ForJet(ctx context.Context, jetID insolar.JetID) map[insolar.ID]LifelineMeta {
	m.storageLock.RLock()
	defer m.storageLock.RUnlock()

	idxByJet := m.jetIndexAccessor.For(jetID)

	res := map[insolar.ID]LifelineMeta{}

	for id := range idxByJet {
		idx, ok := m.indexStorage[id]
		if ok {
			lstPN, lstOk := m.pulseIndex.LastUsage(id)
			if !lstOk {
				panic("index isn't in a consistent state")
			}

			res[id] = LifelineMeta{
				Index:    CloneIndex(idx),
				LastUsed: lstPN,
			}
		}
	}

	return res
}

// ForPulseAndJet returns a collection of lifelines for a provided jetID and a pulse number
func (m *IndexMemory) ForPulseAndJet(
	ctx context.Context,
	pn insolar.PulseNumber,
	jetID insolar.JetID,
) map[insolar.ID]Lifeline {
	m.storageLock.RLock()
	defer m.storageLock.RUnlock()

	idxByJet := m.jetIndexAccessor.For(jetID)
	idxByPN := m.pulseIndex.ForPN(pn)

	res := map[insolar.ID]Lifeline{}

	for id := range idxByJet {
		_, existInPn := idxByPN[id]
		if existInPn {
			res[id] = m.indexStorage[id]
		}

	}

	return res
}

// DeleteForPN method removes indexes from a indexByPulseStor for a provided pulse
func (m *IndexMemory) DeleteForPN(ctx context.Context, pn insolar.PulseNumber) {
	m.storageLock.Lock()
	defer m.storageLock.Unlock()

	rmIDs := m.pulseIndex.ForPN(pn)
	m.pulseIndex.DeleteForPulse(pn)

	for id := range rmIDs {
		idx, ok := m.indexStorage[id]
		if ok {
			m.jetIndexModifier.Delete(id, idx.JetID)
			delete(m.indexStorage, id)
			stats.Record(ctx,
				statIndexInMemoryRemovedCount.M(1),
			)
		}
	}

}

type IndexDB struct {
	lock sync.RWMutex
	db   store.DB
}

type indexKey insolar.ID

func (k indexKey) Scope() store.Scope {
	return store.ScopeIndex
}

func (k indexKey) ID() []byte {
	res := insolar.ID(k)
	return (&res).Bytes()
}

// NewIndexDB creates new DB storage instance.
func NewIndexDB(db store.DB) *IndexDB {
	return &IndexDB{db: db}
}

// Set saves new index-value in storage.
func (i *IndexDB) Set(ctx context.Context, id insolar.ID, index Lifeline) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	if index.Delegates == nil {
		index.Delegates = map[insolar.Reference]insolar.Reference{}
	}

	return i.set(id, index)
}

// ForID returns index for provided id.
func (i *IndexDB) ForID(ctx context.Context, id insolar.ID) (index Lifeline, err error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	return i.get(id)
}

func (i *IndexDB) set(id insolar.ID, index Lifeline) error {
	key := indexKey(id)

	return i.db.Set(key, EncodeIndex(index))
}

func (i *IndexDB) get(id insolar.ID) (index Lifeline, err error) {
	buff, err := i.db.Get(indexKey(id))
	if err == store.ErrNotFound {
		err = ErrIndexNotFound
		return
	}
	if err != nil {
		return
	}
	index = MustDecodeIndex(buff)
	return
}
