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
	"github.com/insolar/insolar/instrumentation/inslogger"
	"go.opencensus.io/stats"
)

type IndexStorageMemory struct {
	bucketsLock sync.RWMutex
	buckets     map[insolar.PulseNumber]map[insolar.ID]*FilamentIndex
}

func NewIndexStorageMemory() *IndexStorageMemory {
	return &IndexStorageMemory{buckets: map[insolar.PulseNumber]map[insolar.ID]*FilamentIndex{}}
}

func (i *IndexStorageMemory) Index(pn insolar.PulseNumber, objID insolar.ID) *FilamentIndex {
	i.bucketsLock.RLock()
	defer i.bucketsLock.RUnlock()

	objsByPn, ok := i.buckets[pn]
	if !ok {
		return nil
	}

	return objsByPn[objID]
}

// ForPNAndJet returns a collection of buckets for a provided pn and jetID
func (i *IndexStorageMemory) ForPNAndJet(ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID) []FilamentIndex {
	i.bucketsLock.RLock()
	defer i.bucketsLock.RUnlock()

	bucks, ok := i.buckets[pn]
	if !ok {
		return nil
	}

	var res []FilamentIndex

	for _, b := range bucks {
		if b.Lifeline.JetID != jetID {
			continue
		}

		clonedLfl := CloneLifeline(b.Lifeline)
		var clonedRecords []insolar.ID

		clonedRecords = append(clonedRecords, b.PendingRecords...)

		res = append(res, FilamentIndex{
			ObjID:            b.ObjID,
			Lifeline:         clonedLfl,
			LifelineLastUsed: b.LifelineLastUsed,
			PendingRecords:   clonedRecords,
		})
	}

	return res
}

func (i *IndexStorageMemory) CreateIndex(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) *FilamentIndex {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucket := &FilamentIndex{
		ObjID:          objID,
		PendingRecords: []insolar.ID{},
	}

	objsByPn, ok := i.buckets[pn]
	if !ok {
		objsByPn = map[insolar.ID]*FilamentIndex{}
		i.buckets[pn] = objsByPn
	}

	_, ok = objsByPn[objID]
	if !ok {
		objsByPn[objID] = bucket
	}

	inslogger.FromContext(ctx).Debugf("[createPendingBucket] create bucket for obj - %v was created successfully", objID.DebugString())
	return bucket
}

// SetIndex adds a bucket with provided pulseNumber and ID
func (i *IndexStorageMemory) SetIndex(ctx context.Context, pn insolar.PulseNumber, bucket FilamentIndex) error {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucks, ok := i.buckets[pn]
	if !ok {
		bucks = map[insolar.ID]*FilamentIndex{}
		i.buckets[pn] = bucks
	}

	bucks[bucket.ObjID] = &bucket

	stats.Record(ctx,
		statBucketAddedCount.M(1),
	)

	return nil
}

// DeleteForPN deletes all buckets for a provided pulse number
func (i *IndexStorageMemory) DeleteForPN(ctx context.Context, pn insolar.PulseNumber) {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	delete(i.buckets, pn)
}
