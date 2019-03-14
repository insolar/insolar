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

package blob

import (
	"context"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage/db"
	"go.opencensus.io/stats"
)

// StorageMemory is an in-memory struct for blob-storage.
type StorageMemory struct {
	jetIndex db.JetIndexModifier

	lock   sync.RWMutex
	memory map[core.RecordID]Blob
}

// NewStorageMemory creates a new instance of Storage.
func NewStorageMemory() *StorageMemory {
	return &StorageMemory{
		memory:   map[core.RecordID]Blob{},
		jetIndex: db.NewJetIndex(),
	}
}

// ForID returns Blob for provided id.
func (s *StorageMemory) ForID(ctx context.Context, id core.RecordID) (blob Blob, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	b, ok := s.memory[id]
	if !ok {
		err = ErrNotFound
		return
	}

	blob = Clone(b)

	return
}

// Set saves new Blob-value in storage.
func (s *StorageMemory) Set(ctx context.Context, id core.RecordID, blob Blob) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, ok := s.memory[id]
	if ok {
		return ErrOverride
	}

	b := Clone(blob)

	s.memory[id] = b
	s.jetIndex.Add(id, b.JetID)

	blobSize := int64(len(b.Value))

	stats.Record(ctx,
		statBlobInMemorySize.M(blobSize),
		statBlobInMemoryCount.M(1),
	)

	return nil
}
