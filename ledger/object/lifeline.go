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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/internal/ledger/store"
	"go.opencensus.io/stats"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineAccessor -o ./ -s _mock.go

// LifelineAccessor provides info about Index-values from storage.
type LifelineAccessor interface {
	// ForID returns Index for provided id.
	ForID(ctx context.Context, id insolar.ID) (Lifeline, error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineCollectionAccessor -o ./ -s _mock.go

// LifelineCollectionAccessor provides methods for querying a collection of blobs with specific search conditions.
type LifelineCollectionAccessor interface {
	// ForJet returns a collection of lifelines for a provided jetID
	ForJet(ctx context.Context, jetID insolar.JetID) map[insolar.ID]LifelineMeta
	// ForPulseAndJet returns a collection of lifelines for a provided jetID and a pulse number
	ForPulseAndJet(ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID) map[insolar.ID]Lifeline
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineModifier -o ./ -s _mock.go

// LifelineModifier provides methods for setting Index-values to storage.
type LifelineModifier interface {
	// Set saves new Index-value in storage.
	Set(ctx context.Context, id insolar.ID, index Lifeline) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.ExtendedLifelineModifier -o ./ -s _mock.go

// ExtendedLifelineModifier provides methods for setting Index-values to storage.
// The main difference with LifelineModifier is an opportunity to modify a state of an internal pulse-index
type ExtendedLifelineModifier interface {
	// SetWithMeta saves index to the storage and sets its index and pulse number in internal indexes
	SetWithMeta(ctx context.Context, id insolar.ID, pn insolar.PulseNumber, index Lifeline) error
	// SetUsageForPulse updates an internal state of an internal pulse-index
	// Calling this method guaranties that provied pn will be used as a LastUsagePulse for an id
	SetUsageForPulse(ctx context.Context, id insolar.ID, pn insolar.PulseNumber)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineStorage -o ./ -s _mock.go

// LifelineStorage is an union of LifelineAccessor and LifelineModifier.
type LifelineStorage interface {
	LifelineAccessor
	LifelineModifier
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineCleaner -o ./ -s _mock.go

// LifelineCleaner provides an interface for removing interfaces from a storage.
type LifelineCleaner interface {
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

	LatestRequest *insolar.ID
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
	rawIdx := LifelineRaw{
		LatestState:         index.LatestState,
		LatestStateApproved: index.LatestStateApproved,
		ChildPointer:        index.ChildPointer,
		Parent:              index.Parent,
		State:               index.State,
		Delegates:           []DelegateKeyValue{},
		LatestUpdate:        index.LatestUpdate,
		JetID:               index.JetID,
		LatestRequest:       index.LatestRequest,
	}
	for k, d := range index.Delegates {
		rawIdx.Delegates = append(rawIdx.Delegates, DelegateKeyValue{
			Key:   k,
			Value: d,
		})
	}
	data, err := rawIdx.Marshal()
	if err != nil {
		panic("can't marshal lifeline")
	}

	return data
}

// MustDecodeIndex converts byte array into lifeline index struct.
func MustDecodeIndex(buff []byte) (index Lifeline) {
	idx, err := DecodeIndex(buff)
	if err != nil {
		panic(err)
	}

	return idx
}

// DecodeIndex converts byte array into lifeline index struct.
func DecodeIndex(buff []byte) (Lifeline, error) {
	rawIdx := LifelineRaw{}
	err := rawIdx.Unmarshal(buff)
	if err != nil {
		return Lifeline{}, nil
	}

	idx := Lifeline{
		LatestState:         rawIdx.LatestState,
		LatestStateApproved: rawIdx.LatestStateApproved,
		ChildPointer:        rawIdx.ChildPointer,
		Parent:              rawIdx.Parent,
		State:               rawIdx.State,
		Delegates:           map[insolar.Reference]insolar.Reference{},
		LatestUpdate:        rawIdx.LatestUpdate,
		JetID:               rawIdx.JetID,
		LatestRequest:       rawIdx.LatestRequest,
	}
	for _, v := range rawIdx.Delegates {
		idx.Delegates[v.Key] = v.Value
	}

	return idx, nil
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

// LifelineStorageMemory is an in-indexStorage struct for index-storage.
type LifelineStorageMemory struct {
	jetIndexModifier store.JetIndexModifier
	jetIndexAccessor store.JetIndexAccessor
	pulseIndex       PulseIndex

	storageLock  sync.RWMutex
	indexStorage map[insolar.ID]Lifeline
}

// NewIndexMemory creates a new instance of LifelineStorageMemory storage.
func NewIndexMemory() *LifelineStorageMemory {
	idx := store.NewJetIndex()
	return &LifelineStorageMemory{
		indexStorage:     map[insolar.ID]Lifeline{},
		jetIndexModifier: idx,
		jetIndexAccessor: idx,
		pulseIndex:       NewPulseIndex(),
	}
}

// Set saves new Index-value in storage.
func (m *LifelineStorageMemory) Set(ctx context.Context, id insolar.ID, index Lifeline) error {
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
func (m *LifelineStorageMemory) SetWithMeta(ctx context.Context, id insolar.ID, pn insolar.PulseNumber, index Lifeline) error {
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
func (m *LifelineStorageMemory) SetUsageForPulse(ctx context.Context, id insolar.ID, pn insolar.PulseNumber) {
	m.pulseIndex.Add(id, pn)
}

// ForID returns Index for provided id.
func (m *LifelineStorageMemory) ForID(ctx context.Context, id insolar.ID) (Lifeline, error) {
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
func (m *LifelineStorageMemory) ForJet(ctx context.Context, jetID insolar.JetID) map[insolar.ID]LifelineMeta {
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
func (m *LifelineStorageMemory) ForPulseAndJet(
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
func (m *LifelineStorageMemory) DeleteForPN(ctx context.Context, pn insolar.PulseNumber) {
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
