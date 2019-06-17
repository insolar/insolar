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

type IndexStorageConcrete struct {
	bucketsLock sync.RWMutex
	buckets     map[insolar.PulseNumber]map[insolar.ID]*LockedIndex
}

func NewIndexStorageConcrete() *IndexStorageConcrete {
	return &IndexStorageConcrete{buckets: map[insolar.PulseNumber]map[insolar.ID]*LockedIndex{}}
}

type LockedIndex struct {
	sync.RWMutex

	objectMeta FilamentIndex
}

func (i *IndexStorageConcrete) Index(pn insolar.PulseNumber, objID insolar.ID) *LockedIndex {
	i.bucketsLock.RLock()
	defer i.bucketsLock.RUnlock()

	objsByPn, ok := i.buckets[pn]
	if !ok {
		return nil
	}

	return objsByPn[objID]
}

// ForPNAndJet returns a collection of buckets for a provided pn and jetID
func (i *IndexStorageConcrete) ForPNAndJet(ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID) []FilamentIndex {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucks, ok := i.buckets[pn]
	if !ok {
		return nil
	}

	var res []FilamentIndex

	for _, b := range bucks {
		if b.objectMeta.Lifeline.JetID != jetID {
			continue
		}

		clonedLfl := CloneLifeline(b.objectMeta.Lifeline)
		var clonedRecords []insolar.ID

		clonedRecords = append(clonedRecords, b.objectMeta.PendingRecords...)

		res = append(res, FilamentIndex{
			ObjID:            b.objectMeta.ObjID,
			Lifeline:         clonedLfl,
			LifelineLastUsed: b.objectMeta.LifelineLastUsed,
			PendingRecords:   clonedRecords,
		})
	}

	return res
}

func (i *IndexStorageConcrete) CreateIndex(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) *LockedIndex {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucket := &LockedIndex{
		objectMeta: FilamentIndex{
			ObjID:          objID,
			PendingRecords: []insolar.ID{},
		},
	}

	objsByPn, ok := i.buckets[pn]
	if !ok {
		objsByPn = map[insolar.ID]*LockedIndex{}
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
func (i *IndexStorageConcrete) SetIndex(ctx context.Context, pn insolar.PulseNumber, bucket FilamentIndex) error {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucks, ok := i.buckets[pn]
	if !ok {
		bucks = map[insolar.ID]*LockedIndex{}
		i.buckets[pn] = bucks
	}

	bucks[bucket.ObjID] = &LockedIndex{
		objectMeta: bucket,
	}

	stats.Record(ctx,
		statBucketAddedCount.M(1),
	)

	return nil
}

// DeleteForPN deletes all buckets for a provided pulse number
func (i *IndexStorageConcrete) DeleteForPN(ctx context.Context, pn insolar.PulseNumber) {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucks, ok := i.buckets[pn]
	if !ok {
		return
	}

	delete(i.buckets, pn)

	stats.Record(ctx,
		statBucketRemovedCount.M(int64(len(bucks))),
	)

	for _, buck := range bucks {
		stats.Record(ctx,
			statObjectPendingRecordsInMemoryRemovedCount.M(int64(len(buck.objectMeta.PendingRecords))),
		)
	}
}
