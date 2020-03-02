// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package object

import (
	"context"
	"runtime/debug"
	"sync"

	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/log"
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
func (i *IndexStorageMemory) ForPulse(ctx context.Context, pn insolar.PulseNumber) ([]record.Index, error) {
	i.bucketsLock.RLock()
	defer i.bucketsLock.RUnlock()

	bucks, ok := i.buckets[pn]
	if !ok {
		return nil, ErrIndexNotFound
	}

	res := make([]record.Index, 0, len(bucks))
	for _, b := range bucks {
		res = append(res, clone(b))
	}
	return res, nil
}

func (i *IndexStorageMemory) Set(ctx context.Context, pn insolar.PulseNumber, bucket record.Index) {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	_, ok := i.buckets[pn]
	if !ok {
		i.buckets[pn] = map[insolar.ID]*record.Index{}
	}
	i.set(ctx, pn, bucket)
}

func (i *IndexStorageMemory) SetIfNone(ctx context.Context, pn insolar.PulseNumber, bucket record.Index) {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	_, ok := i.buckets[pn]
	if !ok {
		i.buckets[pn] = map[insolar.ID]*record.Index{}
	}
	if _, ok := i.buckets[pn][bucket.ObjID]; ok {
		return
	}
	i.set(ctx, pn, bucket)
}

func (i *IndexStorageMemory) set(ctx context.Context, pn insolar.PulseNumber, bucket record.Index) {
	if i.buckets[pn][bucket.ObjID] != nil {
		savedBuck := i.buckets[pn][bucket.ObjID]
		if savedBuck.LifelineLastUsed > bucket.LifelineLastUsed {
			debug.PrintStack()
			log.Fatal("savedBuck.LifelineLastUsed > bucket.LifelineLastUsed")
		}
		if len(savedBuck.PendingRecords) > len(bucket.PendingRecords) {
			debug.PrintStack()
			log.Fatal("len(savedBuck.PendingRecords) > len(bucket.PendingRecords)")
		}
		if savedBuck.Lifeline.EarliestOpenRequest != nil && bucket.Lifeline.EarliestOpenRequest != nil &&
			*savedBuck.Lifeline.EarliestOpenRequest > *bucket.Lifeline.EarliestOpenRequest {
			debug.PrintStack()
			log.Fatal("*savedBuck.Lifeline.EarliestOpenRequest > *bucket.Lifeline.EarliestOpenRequest")
		}
		if !savedBuck.Lifeline.Parent.IsEmpty() && savedBuck.Lifeline.Parent != bucket.Lifeline.Parent {
			debug.PrintStack()
			log.Fatalf("savedBuck.Lifeline.Parent:%v != bucket.Lifeline.Parent:%v", savedBuck.Lifeline.Parent, bucket.Lifeline.Parent)
		}

		if savedBuck.Lifeline.LatestRequest != nil && bucket.Lifeline.LatestRequest == nil {
			debug.PrintStack()
			log.Fatal("savedBuck.Lifeline.EarliestRequest != nil && bucket.Lifeline.EarliestRequest == nil")
		}
		if savedBuck.Lifeline.LatestRequest != nil && savedBuck.Lifeline.LatestRequest.Pulse() > bucket.Lifeline.LatestRequest.Pulse() {
			debug.PrintStack()
			log.Fatal("savedBuck.Lifeline.EarliestRequest.Pulse() < bucket.Lifeline.EarliestRequest.Pulse()")
		}

		if savedBuck.Lifeline.LatestState != nil && bucket.Lifeline.LatestState == nil {
			debug.PrintStack()
			log.Fatal("savedBuck.Lifeline.LatestState != nil && bucket.Lifeline.LatestState == nil")
		}
		if savedBuck.Lifeline.LatestState != nil && savedBuck.Lifeline.LatestState.Pulse() > bucket.Lifeline.LatestState.Pulse() {
			debug.PrintStack()
			log.Fatal("savedBuck.Lifeline.LatestState.Pulse() < bucket.Lifeline.LatestState.Pulse()")
		}
	}

	if _, ok := i.buckets[pn][bucket.ObjID]; !ok {
		stats.Record(ctx, statIndexesAddedCount.M(1))
	}
	i.buckets[pn][bucket.ObjID] = &bucket
}

// DeleteForPN deletes all buckets for a provided pulse number
func (i *IndexStorageMemory) DeleteForPN(ctx context.Context, pn insolar.PulseNumber) {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	stats.Record(ctx, statIndexesRemovedCount.M(int64(len(i.buckets[pn]))))
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
