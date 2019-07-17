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
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.FilamentModifier -o ./ -s _mock.go

type FilamentModifier interface {
	SetRequest(ctx context.Context, reqID insolar.ID, jetID insolar.JetID, request record.Request) (foundRequest *record.CompositeFilamentRecord, foundResult *record.CompositeFilamentRecord, err error)
	SetResult(ctx context.Context, resID insolar.ID, jetID insolar.JetID, result record.Result) (foundResult *record.CompositeFilamentRecord, err error)
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

	RequestDuplicate(
		ctx context.Context,
		startFrom insolar.PulseNumber,
		objectID, requestID insolar.ID,
		request record.Request,
	) (
		foundRequest *record.CompositeFilamentRecord,
		foundResult *record.CompositeFilamentRecord,
		err error,
	)

	ResultDuplicate(ctx context.Context, startFrom insolar.PulseNumber, objectID, resultID insolar.ID, result record.Result) (foundResult *record.CompositeFilamentRecord, err error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.FilamentCleaner -o ./ -s _mock.go

type FilamentCleaner interface {
	Clear(objID insolar.ID)
}

func NewFilamentModifier(
	indexes object.IndexStorage,
	recordStorage object.RecordModifier,
	pcs insolar.PlatformCryptographyScheme,
	calculator FilamentCalculator,
	pulses pulse.Calculator,
) *FilamentModifierDefault {
	return &FilamentModifierDefault{
		calculator: calculator,
		indexes:    indexes,
		records:    recordStorage,
		pcs:        pcs,
		pulses:     pulses,
	}
}

type FilamentModifierDefault struct {
	calculator FilamentCalculator
	indexes    object.IndexStorage
	records    object.RecordModifier
	pcs        insolar.PlatformCryptographyScheme
	pulses     pulse.Calculator
}

func (m *FilamentModifierDefault) checkObject(ctx context.Context, currentPN insolar.PulseNumber, untilPN insolar.PulseNumber, requestID insolar.ID) (record.Index, error) {
	for {
		idx, err := m.indexes.ForID(ctx, currentPN, requestID)
		if err != nil && err != object.ErrIndexNotFound {
			return idx, errors.Wrap(err, "failed to fetch index")
		}
		if err == nil {
			return idx, nil
		}

		tmpPN, err := m.pulses.Backwards(ctx, currentPN, 1)
		if err != nil {
			return record.Index{}, object.ErrIndexNotFound
		}

		currentPN = tmpPN.PulseNumber
		if currentPN > untilPN {
			return record.Index{}, object.ErrIndexNotFound
		}
	}
}

func (m *FilamentModifierDefault) prepareCreationRequest(ctx context.Context, requestID insolar.ID, request record.Request) error {
	currentPN := requestID.Pulse()
	reason := request.ReasonRef()
	untilPN := reason.Record().Pulse()

	_, err := m.checkObject(ctx, currentPN, untilPN, requestID)
	if err == object.ErrIndexNotFound {
		idx := record.Index{
			ObjID:            requestID,
			PendingRecords:   []insolar.ID{},
			LifelineLastUsed: requestID.Pulse(),
		}
		err := m.indexes.SetIndex(ctx, requestID.Pulse(), idx)
		if err != nil {
			return errors.Wrap(err, "failed to create an object")
		}
		return nil
	}

	return err
}

func (m *FilamentModifierDefault) SetRequest(
	ctx context.Context,
	requestID insolar.ID,
	jetID insolar.JetID,
	request record.Request,
) (*record.CompositeFilamentRecord, *record.CompositeFilamentRecord, error) {
	if requestID.IsEmpty() {
		return nil, nil, errors.New("request id is empty")
	}
	if !jetID.IsValid() {
		return nil, nil, errors.New("jet is not valid")
	}
	if request.ReasonRef().IsEmpty() {
		return nil, nil, ErrEmptyReason
	}

	var objectID insolar.ID

	if request.IsCreationRequest() {
		err := m.prepareCreationRequest(ctx, requestID, request)
		if err != nil {
			return nil, nil, err
		}
		objectID = requestID
	} else {
		if request.AffinityRef() == nil && request.AffinityRef().Record().IsEmpty() {
			return nil, nil, errors.New("request object id is empty")
		}
		objectID = *request.AffinityRef().Record()
	}

	idx, err := m.indexes.ForID(ctx, requestID.Pulse(), objectID)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to fetch index")
	}

	if idx.Lifeline.PendingPointer != nil && requestID.Pulse() < idx.Lifeline.PendingPointer.Pulse() {
		return nil, nil, errors.New("request from the past")
	}

	foundRequest, foundResult, err := m.calculator.RequestDuplicate(ctx, requestID.Pulse(), objectID, requestID, request)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to set request")
	}
	if foundRequest != nil || foundResult != nil {
		return foundRequest, foundResult, err
	}

	// Save request record to storage.
	{
		virtual := record.Wrap(request)
		material := record.Material{Virtual: &virtual, JetID: jetID}
		err := m.records.Set(ctx, requestID, material)
		if err != nil && err != object.ErrOverride {
			return nil, nil, errors.Wrap(err, "failed to save a request record")
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
			return nil, nil, errors.Wrap(err, "failed to save filament record")
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
		return nil, nil, errors.Wrap(err, "failed to update index")
	}

	inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"object_id":  objectID.DebugString(),
		"request_id": requestID.DebugString(),
	}).Debug("set request")

	return nil, nil, nil
}

func (m *FilamentModifierDefault) SetResult(ctx context.Context, resultID insolar.ID, jetID insolar.JetID, result record.Result) (*record.CompositeFilamentRecord, error) {
	if resultID.IsEmpty() {
		return nil, errors.New("request id is empty")
	}
	if !jetID.IsValid() {
		return nil, errors.New("jet is not valid")
	}
	if result.Object.IsEmpty() {
		return nil, errors.New("object is empty")
	}

	objectID := result.Object

	idx, err := m.indexes.ForID(ctx, resultID.Pulse(), objectID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update a result's filament")
	}

	foundRes, err := m.calculator.ResultDuplicate(ctx, resultID.Pulse(), objectID, resultID, result)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save a result record")
	}
	if foundRes != nil {
		return foundRes, nil
	}

	// Save request record to storage.
	{
		virtual := record.Wrap(result)
		material := record.Material{Virtual: &virtual, JetID: jetID}
		err := m.records.Set(ctx, resultID, material)
		if err != nil && err != object.ErrOverride {
			return nil, errors.Wrap(err, "failed to save a result record")
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
			return nil, errors.Wrap(err, "failed to save filament record")
		}
		filamentID = id
	}

	pending, err := m.calculator.PendingRequests(ctx, resultID.Pulse(), objectID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate pending requests")
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
		return nil, errors.Wrap(err, "failed to create a meta-record about pending request")
	}

	inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"object_id":  objectID.DebugString(),
		"request_id": result.Request.Record().DebugString(),
		"result_id":  resultID.DebugString(),
	}).Debug("set result")

	return nil, nil
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
	logger := inslogger.FromContext(ctx).WithField("object_id", objectID.DebugString())

	logger.Debug("started collecting pending requests")
	defer logger.Debug("finished collecting pending requests")

	idx, err := c.indexes.ForID(ctx, pulse, objectID)
	if err != nil {
		return nil, err
	}

	cache := c.cache.Get(objectID)
	cache.Lock()
	defer cache.Unlock()

	if idx.Lifeline.PendingPointer == nil {
		return []insolar.ID{}, nil
	}
	if idx.Lifeline.EarliestOpenRequest == nil {
		return []insolar.ID{}, nil
	}

	iter := newFetchingIterator(
		ctx,
		cache,
		objectID,
		*idx.Lifeline.PendingPointer,
		*idx.Lifeline.EarliestOpenRequest,
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
		case *record.IncomingRequest:
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

func (c *FilamentCalculatorDefault) ResultDuplicate(
	ctx context.Context, startFrom insolar.PulseNumber, objectID, resultID insolar.ID, result record.Result,
) (*record.CompositeFilamentRecord, error) {
	logger := inslogger.FromContext(ctx).WithField("object_id", objectID.DebugString())

	logger.Debug("started to search duplicated requests")
	defer logger.Debug("finished to search duplicated requests")

	if result.Request.IsEmpty() {
		return nil, errors.New("request is empty")
	}
	idx, err := c.indexes.ForID(ctx, startFrom, objectID)
	if err != nil {
		return nil, err
	}
	if idx.Lifeline.PendingPointer == nil {
		return nil, nil
	}

	cache := c.cache.Get(objectID)
	cache.Lock()
	defer cache.Unlock()

	iter := newFetchingIterator(
		ctx,
		cache,
		objectID,
		*idx.Lifeline.PendingPointer,
		result.Request.Record().Pulse(),
		c.jetFetcher,
		c.coordinator,
		c.sender,
	)

	var foundResult *record.CompositeFilamentRecord

	for iter.HasPrev() {
		rec, err := iter.Prev(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate pending")
		}

		if bytes.Equal(rec.RecordID.Hash(), resultID.Hash()) {
			foundResult = &rec
		}

		if bytes.Equal(rec.RecordID.Hash(), result.Request.Record().Hash()) {
			return foundResult, nil
		}
	}

	return foundResult, errors.New(fmt.Sprintf("request for result is not found"))
}

func (c *FilamentCalculatorDefault) RequestDuplicate(
	ctx context.Context, startFrom insolar.PulseNumber, objectID, requestID insolar.ID, request record.Request,
) (*record.CompositeFilamentRecord, *record.CompositeFilamentRecord, error) {
	logger := inslogger.FromContext(ctx).WithField("object_id", objectID.DebugString())

	logger.Debug("started to search duplicated requests")
	defer logger.Debug("finished to search duplicated requests")

	if request.ReasonRef().IsEmpty() {
		return nil, nil, ErrEmptyReason
	}
	reason := request.ReasonRef()

	idx, err := c.indexes.ForID(ctx, startFrom, objectID)
	if err != nil {
		return nil, nil, err
	}
	if idx.Lifeline.PendingPointer == nil {
		return nil, nil, nil
	}

	cache := c.cache.Get(objectID)
	cache.Lock()
	defer cache.Unlock()

	iter := newFetchingIterator(
		ctx,
		cache,
		objectID,
		*idx.Lifeline.PendingPointer,
		reason.Record().Pulse(),
		c.jetFetcher,
		c.coordinator,
		c.sender,
	)

	_, isOutgoing := request.(*record.OutgoingRequest)
	if !isOutgoing && reason.Record().Pulse() != insolar.PulseNumberAPIRequest {
		exists, err := c.checkReason(ctx, reason)
		if err != nil {
			return nil, nil, err
		}
		if !exists {
			return nil, nil, errors.New("request reason is not found")
		}
	}

	isReasonFound := false
	var foundRequest *record.CompositeFilamentRecord
	var foundResult *record.CompositeFilamentRecord

	for iter.HasPrev() {
		rec, err := iter.Prev(ctx)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to calculate pending")
		}

		if bytes.Equal(rec.RecordID.Hash(), requestID.Hash()) {
			foundRequest = &rec
		}
		if rec.RecordID == *reason.Record() {
			isReasonFound = true
		}

		virtual := record.Unwrap(rec.Record.Virtual)
		if r, ok := virtual.(*record.Result); ok {
			if bytes.Equal(r.Request.Record().Hash(), requestID.Hash()) {
				foundResult = &rec
			}
		}
	}

	if isOutgoing && !isReasonFound {
		return nil, nil, errors.New("request reason is not found")
	}

	return foundRequest, foundResult, nil
}

func (c *FilamentCalculatorDefault) Clear(objID insolar.ID) {
	c.cache.Delete(objID)
}

func (c *FilamentCalculatorDefault) checkReason(ctx context.Context, reason insolar.Reference) (bool, error) {
	isBeyond, err := c.coordinator.IsBeyondLimit(ctx, reason.Record().Pulse())
	if err != nil {
		return false, errors.Wrap(err, "failed to calculate limit")
	}
	var node *insolar.Reference
	if isBeyond {
		node, err = c.coordinator.Heavy(ctx)
		if err != nil {
			return false, errors.Wrap(err, "failed to calculate node")
		}
	} else {
		jetID, err := c.jetFetcher.Fetch(ctx, *reason.Record(), reason.Record().Pulse())
		if err != nil {
			return false, errors.Wrap(err, "failed to fetch jet")
		}
		node, err = c.coordinator.NodeForJet(ctx, *jetID, reason.Record().Pulse())
		if err != nil {
			return false, errors.Wrap(err, "failed to calculate node")
		}
	}
	msg, err := payload.NewMessage(&payload.GetRequest{
		RequestID: *reason.Record(),
	})
	if err != nil {
		return false, errors.Wrap(err, "failed to check an object existence")
	}

	reps, done := c.sender.SendTarget(ctx, msg, *node)
	defer done()
	res, ok := <-reps
	if !ok {
		return false, errors.New("no reply for reason check")
	}

	pl, err := payload.UnmarshalFromMeta(res.Payload)
	if err != nil {
		return false, errors.Wrap(err, "failed to unmarshal reply")
	}

	switch concrete := pl.(type) {
	case *payload.Request:
		return true, nil
	case *payload.Error:
		if concrete.Code == payload.CodeObjectNotFound {
			inslogger.FromContext(ctx).Errorf("reason is wrong. %v", concrete.Text)
			return true, nil
		}
		return false, errors.New(concrete.Text)
	default:
		return false, fmt.Errorf("unexpected reply %T", pl)
	}
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

func (c *cacheStore) getIDs() []insolar.ID {
	c.lock.Lock()
	defer c.lock.Unlock()

	res := make([]insolar.ID, len(c.caches))
	i := 0
	for id := range c.caches {
		res[i] = id
		i++
	}

	return res
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
	readUntil insolar.PulseNumber,
	fetcher jet.Fetcher,
	coordinator jet.Coordinator,
	sender bus.Sender,
) *fetchingIterator {
	return &fetchingIterator{
		iter:        cache.NewIterator(ctx, from),
		cache:       cache,
		objectID:    objectID,
		readUntil:   readUntil,
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
	logger := inslogger.FromContext(ctx)

	rec, err := i.iter.Prev(ctx)
	switch err {
	case nil:
		return rec, nil

	case object.ErrNotFound:
		// Update cache from network.
		logger.Debug("fetching requests from network")
		defer logger.Debug("received requests from network")
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

func (i *fetchingIterator) fetchFromNetwork(
	ctx context.Context, forID insolar.ID,
) ([]record.CompositeFilamentRecord, error) {
	ctx, span := instracer.StartSpan(ctx, "fetchingIterator.fetchFromNetwork")
	defer span.End()

	isBeyond, err := i.coordinator.IsBeyondLimit(ctx, forID.Pulse())
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate limit")
	}
	var node *insolar.Reference
	if isBeyond {
		node, err = i.coordinator.Heavy(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate node")
		}
	} else {
		jetID, err := i.jetFetcher.Fetch(ctx, i.objectID, forID.Pulse())
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch jet")
		}
		node, err = i.coordinator.NodeForJet(ctx, *jetID, forID.Pulse())
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
