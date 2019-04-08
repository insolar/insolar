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

//go:generate minimock -i github.com/insolar/insolar/ledger/storage/object.IndexAccessor -o ./ -s _mock.go

// IndexAccessor provides info about Index-values from storage.
type IndexAccessor interface {
	// ForID returns Index for provided id.
	ForID(ctx context.Context, id insolar.ID) (Lifeline, error)
}

// CollectionIndexAccessor provides methods for querying a collection of blobs with specific search conditions.
type CollectionIndexAccessor interface {
	// ForPulseAndJet returns []Blob for a provided jetID and a pulse number.
	ForPulseAndJet(ctx context.Context, jetID insolar.JetID, pn insolar.PulseNumber) map[insolar.ID]Lifeline
}

//go:generate minimock -i github.com/insolar/insolar/ledger/storage/object.IndexModifier -o ./ -s _mock.go

// IndexModifier provides methods for setting Index-values to storage.
type IndexModifier interface {
	// Set saves new Index-value in storage.
	Set(ctx context.Context, id insolar.ID, index Lifeline) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/storage/object.IndexCleaner -o ./ -s _mock.go

// IndexCleaner provides an interface for removing interfaces from a storage.
type IndexCleaner interface {
	// RemoveWithIDs method removes interfaces from a storage for a provided map of ids
	RemoveWithIDs(ctx context.Context, ids map[insolar.ID]struct{})
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

// IndexMemory is an in-memory struct for index-storage.
type IndexMemory struct {
	jetIndex store.JetIndexModifier

	lock   sync.RWMutex
	memory map[insolar.ID]Lifeline
}

// NewIndexMemory creates a new instance of IndexMemory storage.
func NewIndexMemory() *IndexMemory {
	return &IndexMemory{
		memory:   map[insolar.ID]Lifeline{},
		jetIndex: store.NewJetIndex(),
	}
}

// Set saves new Index-value in storage.
func (m *IndexMemory) Set(ctx context.Context, id insolar.ID, index Lifeline) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	idx := CloneIndex(index)

	m.memory[id] = idx
	m.jetIndex.Add(id, idx.JetID)

	stats.Record(ctx,
		statIndexInMemoryCount.M(1),
	)

	return nil
}

// ForID returns Index for provided id.
func (m *IndexMemory) ForID(ctx context.Context, id insolar.ID) (index Lifeline, err error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	idx, ok := m.memory[id]
	if !ok {
		err = ErrIndexNotFound
		return
	}

	index = CloneIndex(idx)

	return
}

// ForPulse returns Index for provided id.
func (m *IndexMemory) ForPulseAndJet(ctx context.Context, jetID insolar.JetID, pn insolar.PulseNumber) map[insolar.ID]Lifeline {
	m.lock.RLock()
	defer m.lock.RUnlock()

	res := map[insolar.ID]Lifeline{}
	for id, idx := range m.memory {
		if id.Pulse() != pn || idx.JetID != jetID {
			continue
		}

		res[id] = CloneIndex(idx)
	}

	return res
}

// RemoveWithIDs method removes interfaces from a storage for a provided map of ids
func (m *IndexMemory) RemoveWithIDs(ctx context.Context, ids map[insolar.ID]struct{}) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for id, idx := range m.memory {
		_, ok := ids[id]
		if ok {
			delete(m.memory, id)
			m.jetIndex.Delete(id, idx.JetID)
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
