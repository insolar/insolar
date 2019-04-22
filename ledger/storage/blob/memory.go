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

package blob

import (
	"context"
	"sync"

	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/internal/ledger/store"
)

// StorageMemory is an in-blobsStor struct for blob-storage.
type StorageMemory struct {
	jetIndexModifier store.JetIndexModifier
	jetIndexAccessor store.JetIndexAccessor

	lock      sync.RWMutex
	blobsStor map[insolar.ID]Blob
}

// NewStorageMemory creates a new instance of Storage.
func NewStorageMemory() *StorageMemory {
	ji := store.NewJetIndex()
	return &StorageMemory{
		blobsStor:        map[insolar.ID]Blob{},
		jetIndexModifier: ji,
		jetIndexAccessor: ji,
	}
}

// ForID returns Blob for provided id.
func (s *StorageMemory) ForID(ctx context.Context, id insolar.ID) (blob Blob, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	b, ok := s.blobsStor[id]
	if !ok {
		err = ErrNotFound
		return
	}

	blob = Clone(b)

	return
}

// Set saves new Blob-value in storage.
func (s *StorageMemory) Set(ctx context.Context, id insolar.ID, blob Blob) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, ok := s.blobsStor[id]
	if ok {
		return ErrOverride
	}

	b := Clone(blob)

	s.blobsStor[id] = b
	s.jetIndexModifier.Add(id, b.JetID)

	blobSize := int64(len(b.Value))

	stats.Record(ctx,
		statBlobInMemorySize.M(blobSize),
		statBlobInMemoryCount.M(1),
	)

	return nil
}

// ForPulse returns []Blob for a provided jetID and a pulse number.
func (s *StorageMemory) ForPulse(ctx context.Context, jetID insolar.JetID, pn insolar.PulseNumber) []Blob {
	s.lock.RLock()
	defer s.lock.RUnlock()

	ids := s.jetIndexAccessor.For(jetID)
	var res []Blob
	for id := range ids {
		if id.Pulse() == pn {
			b := s.blobsStor[id]
			res = append(res, b)
		}
	}

	return res
}

// DeleteForPN cleans blobs for a provided pulse from blobsStor
func (s *StorageMemory) DeleteForPN(ctx context.Context, pulse insolar.PulseNumber) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for id, blob := range s.blobsStor {
		if id.Pulse() != pulse {
			continue
		}

		s.jetIndexModifier.Delete(id, blob.JetID)
		delete(s.blobsStor, id)

		stats.Record(ctx,
			statBlobInMemoryRemovedCount.M(1),
		)
	}
}
