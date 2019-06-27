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
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.FilamentModifier -o ./ -s _mock.go

type FilamentModifier interface {
	SetRequest(ctx context.Context, reqID insolar.ID, jetID insolar.JetID, request record.Request) error
	SetResult(ctx context.Context, resID insolar.ID, jetID insolar.JetID, result record.Result) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.FilamentCalculator -o ./ -s _mock.go

type FilamentCalculator interface {
	// Requests goes to network.
	Requests(
		ctx context.Context,
		objectID, from insolar.ID,
		readUntil, calcPulse insolar.PulseNumber,
	) ([]record.CompositeFilamentRecord, error)

	// PendingRequests only looks locally.
	PendingRequests(ctx context.Context, pulse insolar.PulseNumber, objectID insolar.ID) ([]insolar.ID, error)
}

type FilamentCleaner interface {
	Clear(objID insolar.ID)
}

func NewFilamentModifier(
	indexes object.IndexStorage,
	recordStorage object.RecordModifier,
	pcs insolar.PlatformCryptographyScheme,
	calculator FilamentCalculator,
) *FilamentModifierDefault {
	return &FilamentModifierDefault{
		calculator: calculator,
		indexes:    indexes,
		records:    recordStorage,
		pcs:        pcs,
	}
}

type FilamentModifierDefault struct {
	cache *cacheStore

	calculator FilamentCalculator
	indexes    object.IndexStorage
	records    object.RecordModifier
	pcs        insolar.PlatformCryptographyScheme
}

func (m *FilamentModifierDefault) SetRequest(ctx context.Context, requestID insolar.ID, jetID insolar.JetID, request record.Request) error {
	if requestID.IsEmpty() {
		return errors.New("request id is empty")
	}
	if !jetID.IsValid() {
		return errors.New("jet is not valid")
	}
	if request.Object == nil && request.Object.Record().IsEmpty() {
		return errors.New("request object id is empty")
	}

	objectID := *request.Object.Record()

	idx, err := m.indexes.ForID(ctx, requestID.Pulse(), objectID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch index")
	}

	if idx.Lifeline.PendingPointer != nil && requestID.Pulse() < idx.Lifeline.PendingPointer.Pulse() {
		return errors.New("request from the past")
	}

	// Save request record to storage.
	{
		virtual := record.Wrap(request)
		material := record.Material{Virtual: &virtual, JetID: jetID}
		err := m.records.Set(ctx, requestID, material)
		if err != nil && err != object.ErrOverride {
			return errors.Wrap(err, "failed to save a request record")
		}
	}

	var filamentID insolar.ID
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
		err := m.records.Set(ctx, id, material)
		if err != nil {
			return errors.Wrap(err, "failed to save filament record")
		}
		filamentID = id
	}

	idx.Lifeline.PendingPointer = &filamentID
	if idx.Lifeline.EarliestOpenRequest == nil {
		pn := requestID.Pulse()
		idx.Lifeline.EarliestOpenRequest = &pn
	}

	err = m.indexes.SetIndex(ctx, requestID.Pulse(), idx)
	if err != nil {
		return errors.Wrap(err, "failed to update index")
	}

	return nil
}

func (m *FilamentModifierDefault) SetResult(ctx context.Context, resultID insolar.ID, jetID insolar.JetID, result record.Result) error {
	if resultID.IsEmpty() {
		return errors.New("request id is empty")
	}
	if !jetID.IsValid() {
		return errors.New("jet is not valid")
	}
	if result.Object.IsEmpty() {
		return errors.New("object is empty")
	}

	objectID := result.Object

	idx, err := m.indexes.ForID(ctx, resultID.Pulse(), objectID)
	if err != nil {
		return errors.Wrap(err, "failed to update a result's filament")
	}

	// Save request record to storage.
	{
		virtual := record.Wrap(result)
		material := record.Material{Virtual: &virtual, JetID: jetID}
		err := m.records.Set(ctx, resultID, material)
		if err != nil && err != object.ErrOverride {
			return errors.Wrap(err, "failed to save a result record")
		}
	}

	var filamentID insolar.ID
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
		err := m.records.Set(ctx, id, material)
		if err != nil {
			return errors.Wrap(err, "failed to save filament record")
		}
		filamentID = id
	}

	pending, err := m.calculator.PendingRequests(ctx, resultID.Pulse(), objectID)
	if err != nil {
		return errors.Wrap(err, "failed to calculate pending requests")
	}
	if len(pending) > 0 {
		calculatedEarliest := pending[0].Pulse()
		idx.Lifeline.EarliestOpenRequest = &calculatedEarliest
	} else {
		idx.Lifeline.EarliestOpenRequest = nil
	}

	idx.Lifeline.PendingPointer = &filamentID
	err = m.indexes.SetIndex(ctx, resultID.Pulse(), idx)
	if err != nil {
		return errors.Wrap(err, "failed to create a meta-record about pending request")
	}

	return nil
}

type FilamentCalculatorDefault struct {
	cache       *cacheStore
	indexes     object.IndexAccessor
	coordinator jet.Coordinator
	jetFetcher  jet.Fetcher
	sender      bus.Sender
}

func NewFilamentCalculator(
	indexes object.IndexAccessor,
	records object.RecordAccessor,
	coordinator jet.Coordinator,
	jetFetcher jet.Fetcher,
	sender bus.Sender,
) *FilamentCalculatorDefault {
	return &FilamentCalculatorDefault{
		cache:       newCacheStore(records),
		indexes:     indexes,
		coordinator: coordinator,
		jetFetcher:  jetFetcher,
		sender:      sender,
	}
}

func (c *FilamentCalculatorDefault) Requests(
	ctx context.Context, objectID insolar.ID, from insolar.ID, readUntil, calcPulse insolar.PulseNumber,
) ([]record.CompositeFilamentRecord, error) {
	_, err := c.indexes.ForID(ctx, from.Pulse(), objectID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch index")
	}

	cache := c.cache.Get(objectID)
	cache.RLock()
	defer cache.RUnlock()

	iter := cache.NewIterator(ctx, from)
	var segment []record.CompositeFilamentRecord
	for iter.HasPrev() {
		rec, err := iter.Prev(ctx)
		if err == object.ErrNotFound {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to get filament")
		}
		if rec.MetaID.Pulse() < readUntil {
			break
		}

		segment = append(segment, rec)
	}

	return segment, nil
}

func (c *FilamentCalculatorDefault) PendingRequests(
	ctx context.Context, pulse insolar.PulseNumber, objectID insolar.ID,
) ([]insolar.ID, error) {
	idx, err := c.indexes.ForID(ctx, pulse, objectID)
	if err != nil {
		return nil, err
	}

	cache := c.cache.Get(objectID)
	cache.Lock()
	defer cache.Unlock()

	if idx.Lifeline.EarliestOpenRequest == nil {
		return []insolar.ID{}, nil
	}

	iter := newFetchingIterator(
		ctx,
		cache,
		objectID,
		*idx.Lifeline.PendingPointer,
		*idx.Lifeline.EarliestOpenRequest,
		pulse,
		c.jetFetcher,
		c.coordinator,
		c.sender,
	)

	var pending []insolar.ID
	hasResult := map[insolar.ID]struct{}{}
	for iter.HasPrev() {
		rec, err := iter.Prev(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate pending")
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
		ordered[count-i-1] = id
	}

	return ordered, nil
}

func (c *FilamentCalculatorDefault) Clear(objID insolar.ID) {
	cache := c.cache.Get(objID)
	cache.Lock()
	cache.Clear()
	cache.Unlock()
}

type fetchingIterator struct {
	iter  filamentIterator
	cache *filamentCache

	objectID             insolar.ID
	readUntil, calcPulse insolar.PulseNumber

	jetFetcher  jet.Fetcher
	coordinator jet.Coordinator
	sender      bus.Sender
}

type cacheStore struct {
	records object.RecordAccessor

	lock   sync.Mutex
	caches map[insolar.ID]*filamentCache
}

func newCacheStore(r object.RecordAccessor) *cacheStore {
	return &cacheStore{
		caches:  map[insolar.ID]*filamentCache{},
		records: r,
	}
}

func (c *cacheStore) Get(id insolar.ID) *filamentCache {
	c.lock.Lock()
	defer c.lock.Unlock()

	obj, ok := c.caches[id]
	if !ok {
		obj = newFilamentCache(c.records)
		c.caches[id] = obj
	}

	return obj
}

func (c *cacheStore) Delete(id insolar.ID) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.caches, id)
}

type filamentCache struct {
	sync.RWMutex
	cache map[insolar.ID]record.CompositeFilamentRecord

	records object.RecordAccessor
}

func newFilamentCache(r object.RecordAccessor) *filamentCache {
	return &filamentCache{
		cache:   map[insolar.ID]record.CompositeFilamentRecord{},
		records: r,
	}
}

func (c *filamentCache) Update(recs []record.CompositeFilamentRecord) {
	for _, rec := range recs {
		c.cache[rec.MetaID] = rec
	}
}

func (c *filamentCache) NewIterator(ctx context.Context, from insolar.ID) filamentIterator {
	return filamentIterator{
		currentID: &from,
		cache:     c,
	}
}

func (c *filamentCache) Clear() {
	c.cache = map[insolar.ID]record.CompositeFilamentRecord{}
}

type filamentIterator struct {
	currentID *insolar.ID
	cache     *filamentCache
}

func (i *filamentIterator) PrevID() *insolar.ID {
	return i.currentID
}

func (i *filamentIterator) HasPrev() bool {
	return i.currentID != nil
}

func (i *filamentIterator) Prev(ctx context.Context) (record.CompositeFilamentRecord, error) {
	if i.currentID == nil {
		return record.CompositeFilamentRecord{}, object.ErrNotFound
	}

	composite, ok := i.cache.cache[*i.currentID]
	if ok {
		virtual := record.Unwrap(composite.Meta.Virtual)
		filament, ok := virtual.(*record.PendingFilament)
		if !ok {
			return record.CompositeFilamentRecord{}, fmt.Errorf("unexpected filament record %T", virtual)
		}
		i.currentID = filament.PreviousRecord
		return composite, nil
	}

	// Fetching filament record.
	filamentRecord, err := i.cache.records.ForID(ctx, *i.currentID)
	if err != nil {
		return record.CompositeFilamentRecord{}, err
	}
	virtual := record.Unwrap(filamentRecord.Virtual)
	filament, ok := virtual.(*record.PendingFilament)
	if !ok {
		return record.CompositeFilamentRecord{}, fmt.Errorf("unexpected filament record %T", virtual)
	}
	composite.MetaID = *i.currentID
	composite.Meta = filamentRecord

	// Fetching primary record.
	rec, err := i.cache.records.ForID(ctx, filament.RecordID)
	if err != nil {
		return record.CompositeFilamentRecord{}, err
	}
	composite.RecordID = filament.RecordID
	composite.Record = rec

	// Adding to cache.
	i.cache.cache[*i.currentID] = composite
	i.currentID = filament.PreviousRecord

	return composite, nil
}

func newFetchingIterator(
	ctx context.Context,
	cache *filamentCache,
	objectID, from insolar.ID,
	readUntil, calcPulse insolar.PulseNumber,
	fetcher jet.Fetcher,
	coordinator jet.Coordinator,
	sender bus.Sender,
) *fetchingIterator {
	return &fetchingIterator{
		iter:        cache.NewIterator(ctx, from),
		cache:       cache,
		objectID:    objectID,
		readUntil:   readUntil,
		calcPulse:   calcPulse,
		jetFetcher:  fetcher,
		coordinator: coordinator,
		sender:      sender,
	}
}

func (i *fetchingIterator) PrevID() *insolar.ID {
	return i.iter.PrevID()
}

func (i *fetchingIterator) HasPrev() bool {
	return i.iter.HasPrev() && i.iter.PrevID().Pulse() >= i.readUntil
}

func (i *fetchingIterator) Prev(ctx context.Context) (record.CompositeFilamentRecord, error) {
	rec, err := i.iter.Prev(ctx)
	switch err {
	case nil:
		return rec, nil

	case object.ErrNotFound:
		// Update cache from network.
		recs, err := i.fetchFromNetwork(ctx, *i.PrevID(), i.calcPulse)
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

func (i *fetchingIterator) fetchFromNetwork(
	ctx context.Context, forID insolar.ID, calcPulse insolar.PulseNumber,
) ([]record.CompositeFilamentRecord, error) {
	ctx, span := instracer.StartSpan(ctx, "fetchingIterator.fetchFromNetwork")
	defer span.End()

	isBeyond, err := i.coordinator.IsBeyondLimit(ctx, i.calcPulse, forID.Pulse())
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate limit")
	}
	var node *insolar.Reference
	if isBeyond {
		node, err = i.coordinator.Heavy(ctx, calcPulse)
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate node")
		}
	} else {
		jetID, err := i.jetFetcher.Fetch(ctx, i.objectID, forID.Pulse())
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch jet")
		}
		node, err = i.coordinator.NodeForJet(ctx, *jetID, i.calcPulse, forID.Pulse())
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate node")
		}
	}
	if *node == i.coordinator.Me() {
		return nil, errors.New("tried to send message to self")
	}

	span.AddAttributes(
		trace.StringAttribute("objID", i.objectID.DebugString()),
		trace.StringAttribute("startFrom", forID.DebugString()),
		trace.StringAttribute("readUntil", i.readUntil.String()),
	)

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
