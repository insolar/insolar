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

package executor

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/network/storage"
	"github.com/pkg/errors"
)

type FilamentModifier interface {
	SetRequest(ctx context.Context, reqID insolar.ID, jetID insolar.JetID, request record.Request) error
	SetResult(ctx context.Context, resID insolar.ID, jetID insolar.JetID, result record.Result) error
}

func NewFilamentManager(r object.RecordAccessor) *FilamentManager {
	return &FilamentManager{
		cache: newCacheStore(r),
	}
}

type FilamentManager struct {
	cache *cacheStore

	idxAccessor     object.IndexAccessor
	idxModifier     object.IndexModifier
	idLocker        object.IDLocker
	recordStorage   object.RecordStorage
	coordinator     jet.Coordinator
	pcs             insolar.PlatformCryptographyScheme
	pulseCalculator storage.PulseCalculator
	busWM           bus.Sender
	jetFetcher      jet.Fetcher
}

func (m *FilamentManager) SetRequest(ctx context.Context, requestID insolar.ID, jetID insolar.JetID, request record.Request) error {
	if request.Object == nil {
		return errors.New("object is empty")
	}

	objectID := *request.Object.Record()

	m.idLocker.Lock(&objectID)
	defer m.idLocker.Unlock(&objectID)

	idx := m.idxAccessor.Index(requestID.Pulse(), objectID)
	if idx == nil {
		return object.ErrLifelineNotFound
	}

	if idx.Lifeline.PendingPointer != nil && requestID.Pulse() < idx.Lifeline.PendingPointer.Pulse() {
		return errors.New("request from the past")
	}

	var composite record.CompositeFilamentRecord

	// Save request record to storage.
	{
		virtual := record.Wrap(request)
		hash := record.HashVirtual(m.pcs.ReferenceHasher(), virtual)
		id := *insolar.NewID(requestID.Pulse(), hash)
		material := record.Material{Virtual: &virtual, JetID: jetID}
		err := m.recordStorage.Set(ctx, id, material)
		if err != nil {
			return errors.Wrap(err, "failed to save request record")
		}
		composite.RecordID = id
		composite.Record = material
	}

	// Save filament record to storage.
	{
		rec := record.PendingFilament{
			RecordID:       requestID,
			PreviousRecord: idx.Lifeline.PendingPointer,
		}
		virtual := record.Wrap(rec)
		hash := record.HashVirtual(m.pcs.ReferenceHasher(), virtual)
		id := *insolar.NewID(requestID.Pulse(), hash)
		material := record.Material{Virtual: &virtual, JetID: jetID}
		err := m.recordStorage.Set(ctx, id, material)
		if err != nil {
			return errors.Wrap(err, "failed to save filament record")
		}
		composite.MetaID = id
		composite.Meta = material
	}

	idx.Lifeline.PendingPointer = &composite.MetaID
	if idx.Lifeline.EarliestOpenRequest == nil {
		pn := requestID.Pulse()
		idx.Lifeline.EarliestOpenRequest = &pn
	}

	err := m.idxModifier.SetIndex(ctx, requestID.Pulse(), *idx)
	if err != nil {
		return errors.Wrap(err, "failed to update index")
	}

	return nil
}

func (m *FilamentManager) SetResult(ctx context.Context, resultID insolar.ID, jetID insolar.JetID, result record.Result) error {
	if result.Object.IsEmpty() {
		return errors.New("object is empty")
	}

	objectID := result.Object

	m.idLocker.Lock(&objectID)
	defer m.idLocker.Unlock(&objectID)

	idx := m.idxAccessor.Index(resultID.Pulse(), objectID)
	if idx == nil {
		return object.ErrLifelineNotFound
	}

	var filamentRecord record.CompositeFilamentRecord

	// Save request record to storage.
	{
		virtual := record.Wrap(result)
		hash := record.HashVirtual(m.pcs.ReferenceHasher(), virtual)
		id := *insolar.NewID(resultID.Pulse(), hash)
		material := record.Material{Virtual: &virtual, JetID: jetID}
		err := m.recordStorage.Set(ctx, id, material)
		if err != nil {
			return errors.Wrap(err, "failed to save request record")
		}
		filamentRecord.RecordID = id
		filamentRecord.Record = material
	}

	// Save filament record to storage.
	{
		rec := record.PendingFilament{
			RecordID:       resultID,
			PreviousRecord: idx.Lifeline.PendingPointer,
		}
		virtual := record.Wrap(rec)
		hash := record.HashVirtual(m.pcs.ReferenceHasher(), virtual)
		id := *insolar.NewID(resultID.Pulse(), hash)
		material := record.Material{Virtual: &virtual, JetID: jetID}
		err := m.recordStorage.Set(ctx, id, material)
		if err != nil {
			return errors.Wrap(err, "failed to save filament record")
		}
		filamentRecord.MetaID = id
		filamentRecord.Meta = material
	}

	pending, err := m.calculatePending(ctx, resultID.Pulse(), objectID, *idx)
	if err != nil {
		return errors.Wrap(err, "failed to calculate pending requests")
	}
	if len(pending) > 0 {
		calculatedEarliest := pending[0].Pulse()
		idx.Lifeline.EarliestOpenRequest = &calculatedEarliest
		err = m.idxModifier.SetIndex(ctx, resultID.Pulse(), *idx)
		if err != nil {
			return errors.Wrap(err, "failed to create a meta-record about pending request")
		}

	}

	return nil
}

func (m *FilamentManager) calculatePending(
	ctx context.Context,
	pulse insolar.PulseNumber,
	objectID insolar.ID,
	idx object.FilamentIndex,
) ([]insolar.ID, error) {
	cache := m.cache.Get(objectID)
	iter := m.newFetchingIterator(
		ctx,
		cache,
		objectID,
		*idx.Lifeline.PendingPointer,
		*idx.Lifeline.EarliestOpenRequest,
		pulse,
	)

	var pending []insolar.ID
	hasResult := map[insolar.ID]struct{}{}
	for iter.HasPrev() {
		rec, err := iter.Prev(ctx)
		if err != nil {
			break
		}

		virtual := record.Unwrap(rec.Record.Virtual)
		switch r := virtual.(type) {
		case *record.Request:
			if _, ok := hasResult[rec.RecordID]; !ok {
				pending = append(pending, rec.RecordID)
			}
		case *record.Result:
			hasResult[*r.Request.Record()] = struct{}{}
		}
	}

	// We need to reverse pending because we iterated from the end when selecting them.
	ordered := make([]insolar.ID, len(pending))
	count := len(pending)
	for i, id := range pending {
		ordered[count-i] = id
	}

	return ordered, nil
}

type fetchingIterator struct {
	iter  filamentIterator
	cache *filamentCache

	objectID             insolar.ID
	readUntil, calcPulse insolar.PulseNumber

	records     object.RecordAccessor
	jetFetcher  jet.Fetcher
	coordinator jet.Coordinator
	sender      bus.Sender
}

func (m *FilamentManager) newFetchingIterator(
	ctx context.Context,
	cache *filamentCache,
	objectID, from insolar.ID,
	readUntil, calcPulse insolar.PulseNumber,
) *fetchingIterator {
	return &fetchingIterator{
		iter:        cache.NewIterator(ctx, from),
		cache:       cache,
		objectID:    objectID,
		readUntil:   readUntil,
		calcPulse:   calcPulse,
		records:     m.recordStorage,
		jetFetcher:  m.jetFetcher,
		coordinator: m.coordinator,
		sender:      m.busWM,
	}
}

func (i *fetchingIterator) PrevID() *insolar.ID {
	return i.iter.PrevID()
}

func (i *fetchingIterator) HasPrev() bool {
	return i.iter.HasPrev()
}

func (i *fetchingIterator) Prev(ctx context.Context) (record.CompositeFilamentRecord, error) {
	rec, err := i.iter.Prev(ctx)
	switch err {
	case nil:
		return rec, nil

	case object.ErrNotFound:
		// Update cache from network.
		recs, err := i.fetchFromNetwork(ctx, *i.PrevID())
		if err != nil {
			return record.CompositeFilamentRecord{}, errors.Wrap(err, "failed to fetch filament")
		}
		i.cache.Update(recs)

		// Try to iterate again.
		rec, err = i.iter.Prev(ctx)
		if err != nil {
			return record.CompositeFilamentRecord{}, errors.Wrap(err, "failed to update filament")
		}
		return rec, nil

	default:
		return record.CompositeFilamentRecord{}, errors.Wrap(err, "failed to fetch filament")
	}
}

func (i *fetchingIterator) fetchFromNetwork(ctx context.Context, forID insolar.ID) ([]record.CompositeFilamentRecord, error) {
	jetID, err := i.jetFetcher.Fetch(ctx, i.objectID, forID.Pulse())
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch jet")
	}
	node, err := i.coordinator.NodeForJet(ctx, *jetID, i.calcPulse, forID.Pulse())
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate node")
	}
	if *node == i.coordinator.Me() {
		return nil, errors.New("tried to send message to self")
	}

	msg, err := payload.NewMessage(&payload.GetFilament{
		ObjectID:  i.objectID,
		StartFrom: forID,
		ReadUntil: i.readUntil,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create fetching message")
	}
	reps, done := i.sender.SendTarget(ctx, msg, *node)
	defer done()
	res, ok := <-reps
	if !ok {
		return nil, errors.New("no reply for filament fetch")
	}

	pl, err := payload.UnmarshalFromMeta(res.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal reply")
	}
	filaments, ok := pl.(*payload.FilamentSegment)
	if !ok {
		return nil, fmt.Errorf("unexpected reply %T", pl)
	}
	return filaments.Records, nil
}
