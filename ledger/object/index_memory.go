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
	"github.com/insolar/insolar/insolar/record"
	"go.opencensus.io/stats"
)

type IndexStorageMemory struct {
	bucketsLock sync.RWMutex
	buckets     map[insolar.PulseNumber]map[insolar.ID]*record.Index
}

func NewIndexStorageMemory() *IndexStorageMemory {
	return &IndexStorageMemory{
		buckets: map[insolar.PulseNumber]map[insolar.ID]*record.Index{},
	}
}

func (i *IndexStorageMemory) ForID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (record.Index, error) {
	i.bucketsLock.RLock()
	defer i.bucketsLock.RUnlock()

	objsByPn, ok := i.buckets[pn]
	if !ok {
		return record.Index{}, ErrIndexNotFound
	}

	idx, ok := objsByPn[objID]
	if !ok {
		return record.Index{}, ErrIndexNotFound
	}

	return clone(idx), nil
}

// ForPulse returns a collection of buckets for a provided pulse number.
func (i *IndexStorageMemory) ForPulse(ctx context.Context, pn insolar.PulseNumber) []record.Index {
	i.bucketsLock.RLock()
	defer i.bucketsLock.RUnlock()

	bucks, ok := i.buckets[pn]
	if !ok {
		return nil
	}

	res := make([]record.Index, 0, len(bucks))
	for _, b := range bucks {
		res = append(res, clone(b))
	}
	return res
}

func (i *IndexStorageMemory) Set(ctx context.Context, pn insolar.PulseNumber, bucket record.Index) {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	_, ok := i.buckets[pn]
	if !ok {
		i.buckets[pn] = map[insolar.ID]*record.Index{}
	}

	i.buckets[pn][bucket.ObjID] = &bucket

	stats.Record(ctx,
		statBucketAddedCount.M(1),
	)
}

// DeleteForPN deletes all buckets for a provided pulse number
func (i *IndexStorageMemory) DeleteForPN(ctx context.Context, pn insolar.PulseNumber) {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	delete(i.buckets, pn)
}

func clone(index *record.Index) record.Index {
	var clonedRecords []insolar.ID

	clonedRecords = append(clonedRecords, index.PendingRecords...)
	return record.Index{
		Polymorph:        index.Polymorph,
		ObjID:            index.ObjID,
		Lifeline:         CloneLifeline(index.Lifeline),
		LifelineLastUsed: index.LifelineLastUsed,
		PendingRecords:   clonedRecords,
	}
}
