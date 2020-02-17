// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/object"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.FilamentCalculator -o ./ -s _mock.go -g

type FilamentCalculator interface {
	// Requests returns request records for objectID's chain, starts from provided id until provided pulse.
	// TODO: remove calcPulse param
	Requests(
		ctx context.Context,
		objectID, from insolar.ID,
		readUntil insolar.PulseNumber,
	) ([]record.CompositeFilamentRecord, error)

	// OpenedRequests returns all opened requests of object for provided pulse.
	OpenedRequests(
		ctx context.Context,
		pulse insolar.PulseNumber,
		objectID insolar.ID,
		pendingOnly bool,
	) ([]record.CompositeFilamentRecord, error)

	// RequestDuplicate searches two records on objectID chain:
	// First one with same ID as requestID param.
	// Second is the Result record Request field of which equals requestID param.
	// Uses request parameter to check if Reason is not empty and to set pulse for scan limit.
	RequestDuplicate(
		ctx context.Context,
		objectID, requestID insolar.ID,
		request record.Request,
	) (
		foundRequest *record.CompositeFilamentRecord,
		foundResult *record.CompositeFilamentRecord,
		err error,
	)

	ResultDuplicate(
		ctx context.Context,
		objectID, resultID insolar.ID,
		result record.Result,
	) (
		foundResult *record.CompositeFilamentRecord,
		err error,
	)

	// RequestInfo is searching for request and result by objectID, requestID and pulse number
	RequestInfo(
		ctx context.Context,
		objectID insolar.ID,
		requestID insolar.ID,
		pulse insolar.PulseNumber,
	) (
		requestInfo FilamentsRequestInfo,
		err error,
	)
}

type FilamentsRequestInfo struct {
	Request       *record.CompositeFilamentRecord
	Result        *record.CompositeFilamentRecord
	OldestMutable bool
}

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.FilamentCleaner -o ./ -s _mock.go -g

type FilamentCleaner interface {
	ClearIfLonger(limit int)
	ClearAllExcept(ids []insolar.ID)
}

type FilamentCalculatorDefault struct {
	cache       *cacheStore
	indexes     object.MemoryIndexAccessor
	coordinator jet.Coordinator
	jetFetcher  JetFetcher
	sender      bus.Sender
	pulses      pulse.Calculator
}

func NewFilamentCalculator(
	indexes object.MemoryIndexAccessor,
	records object.RecordAccessor,
	coordinator jet.Coordinator,
	jetFetcher JetFetcher,
	sender bus.Sender,
	pulses pulse.Calculator,
) *FilamentCalculatorDefault {
	return &FilamentCalculatorDefault{
		cache:       newCacheStore(records),
		indexes:     indexes,
		coordinator: coordinator,
		jetFetcher:  jetFetcher,
		sender:      sender,
		pulses:      pulses,
	}
}

func (c *FilamentCalculatorDefault) Requests(
	ctx context.Context,
	objectID,
	from insolar.ID,
	readUntil insolar.PulseNumber,
) ([]record.CompositeFilamentRecord, error) {
	_, err := c.indexes.ForID(ctx, from.Pulse(), objectID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch index")
	}

	cache := c.cache.Get(objectID)
	cache.Lock()
	defer cache.Unlock()

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

func (c *FilamentCalculatorDefault) OpenedRequests(ctx context.Context, pulse insolar.PulseNumber, objectID insolar.ID, pendingOnly bool) ([]record.CompositeFilamentRecord, error) {
	idx, err := c.indexes.ForID(ctx, pulse, objectID)
	if err != nil {
		return nil, err
	}

	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"object_id":           objectID.DebugString(),
		"pending_filament_id": idx.Lifeline.LatestRequest.DebugString(),
	})
	logger.Debug("started collecting opened requests")
	defer logger.Debug("finished collecting opened requests")

	cache := c.cache.Get(objectID)
	cache.Lock()
	defer cache.Unlock()

	if idx.Lifeline.LatestRequest == nil {
		return []record.CompositeFilamentRecord{}, nil
	}
	if idx.Lifeline.EarliestOpenRequest == nil {
		return []record.CompositeFilamentRecord{}, nil
	}

	iter := newFetchingIterator(
		ctx,
		cache,
		objectID,
		*idx.Lifeline.LatestRequest,
		*idx.Lifeline.EarliestOpenRequest,
		c.jetFetcher,
		c.coordinator,
		c.sender,
	)

	var opened []record.CompositeFilamentRecord
	hasResult := map[insolar.ID]struct{}{}
	for iter.HasPrev() {
		rec, err := iter.Prev(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate opened")
		}

		// Skip closed requests.
		if _, ok := hasResult[rec.RecordID]; ok {
			continue
		}

		virtual := record.Unwrap(&rec.Record.Virtual)
		switch r := virtual.(type) {
		// result should always go first, before initial request
		case *record.Result:
			hasResult[*r.Request.GetLocal()] = struct{}{}

		case *record.IncomingRequest:
			opened = append(opened, rec)

		case *record.OutgoingRequest:
			_, reasonClosed := hasResult[*r.Reason.GetLocal()]
			isReadyDetached := r.IsDetached() && reasonClosed
			if pendingOnly && !isReadyDetached {
				break
			}

			opened = append(opened, rec)
		}
	}

	// We need to reverse opened to time-ascending because we iterated from the end when selecting them.
	ordered := make([]record.CompositeFilamentRecord, len(opened))
	count := len(opened)
	for i, pend := range opened {
		ordered[count-i-1] = pend
	}

	return ordered, nil
}

func (c *FilamentCalculatorDefault) ResultDuplicate(
	ctx context.Context, objectID, resultID insolar.ID, result record.Result,
) (*record.CompositeFilamentRecord, error) {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"object_id":  objectID.DebugString(),
		"result_id":  resultID.DebugString(),
		"request_id": result.Request.GetLocal().DebugString(),
	})

	logger.Debug("started to search for duplicated results")
	defer logger.Debug("finished to search for duplicated results")

	if result.Request.IsEmpty() {
		return nil, errors.New("request is empty")
	}
	idx, err := c.indexes.ForID(ctx, resultID.Pulse(), objectID)
	if err != nil {
		return nil, err
	}
	if idx.Lifeline.LatestRequest == nil {
		return nil, nil
	}

	cache := c.cache.Get(objectID)
	cache.Lock()
	defer cache.Unlock()

	iter := newFetchingIterator(
		ctx,
		cache,
		objectID,
		*idx.Lifeline.LatestRequest,
		result.Request.GetLocal().Pulse(),
		c.jetFetcher,
		c.coordinator,
		c.sender,
	)

	for iter.HasPrev() {
		rec, err := iter.Prev(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate pending")
		}

		// Result already exists, return it. It should happen before request.
		if bytes.Equal(rec.RecordID.Hash(), resultID.Hash()) {
			logger.Debugf("found duplicate %s", rec.RecordID.DebugString())
			return &rec, nil
		}

		// Another result already exists, return it.
		if res, ok := record.Unwrap(&rec.Record.Virtual).(*record.Result); ok {
			if *res.Request.GetLocal() == *result.Request.GetLocal() {
				return &rec, nil
			}
		}

		// Request found, return nil. It means we didn't find the result since result goes before request on
		// iteration.
		if bytes.Equal(rec.RecordID.Hash(), result.Request.GetLocal().Hash()) {
			return nil, nil
		}
	}

	return nil, fmt.Errorf(
		"request %s for result %s is not found",
		result.Request.GetLocal().DebugString(),
		resultID.DebugString(),
	)
}

func (c *FilamentCalculatorDefault) RequestDuplicate(
	ctx context.Context, objectID, requestID insolar.ID, request record.Request,
) (*record.CompositeFilamentRecord, *record.CompositeFilamentRecord, error) {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"object_id":  objectID.DebugString(),
		"request_id": requestID.DebugString(),
	})

	logger.Debug("started to search for duplicated requests")
	defer logger.Debug("finished searching for duplicated requests")

	reasonRef := request.ReasonRef()
	reasonID := *reasonRef.GetLocal()
	var (
		index record.Index
		err   error
	)
	if request.IsCreationRequest() {
		index, err = c.findIndex(ctx, reasonID, requestID)
		if err == object.ErrIndexNotFound {
			return nil, nil, nil
		}
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to find index")
		}
	} else {
		index, err = c.indexes.ForID(ctx, requestID.Pulse(), objectID)
		if err != nil {
			return nil, nil, err
		}
	}

	if index.Lifeline.LatestRequest == nil {
		logger.Warn("request pointer is nil")
		return nil, nil, nil
	}

	cache := c.cache.Get(index.ObjID)
	cache.Lock()
	defer cache.Unlock()

	iter := newFetchingIterator(
		ctx,
		cache,
		index.ObjID,
		*index.Lifeline.LatestRequest,
		reasonID.Pulse(),
		c.jetFetcher,
		c.coordinator,
		c.sender,
	)

	var foundRequest *record.CompositeFilamentRecord
	var foundResult *record.CompositeFilamentRecord

	for iter.HasPrev() {
		rec, err := iter.Prev(ctx)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to calculate pending")
		}

		if bytes.Equal(rec.RecordID.Hash(), requestID.Hash()) {
			foundRequest = &rec
			logger.Debugf("found duplicate %s", rec.RecordID.DebugString())
		}

		virtual := record.Unwrap(&rec.Record.Virtual)
		if r, ok := virtual.(*record.Result); ok {
			if bytes.Equal(r.Request.GetLocal().Hash(), requestID.Hash()) {
				foundResult = &rec
				logger.Debugf("found result %s", rec.RecordID.DebugString())
			}
		}
	}

	return foundRequest, foundResult, nil
}

func (c *FilamentCalculatorDefault) RequestInfo(
	ctx context.Context,
	objectID insolar.ID,
	requestID insolar.ID,
	pulse insolar.PulseNumber,
) (
	FilamentsRequestInfo,
	error,
) {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"object_id":  objectID.DebugString(),
		"request_id": requestID.DebugString(),
	})

	logger.Debug("start searching request info")
	defer logger.Debug("finished searching request info")

	idx, err := c.indexes.ForID(ctx, pulse, objectID)
	if err != nil {
		return FilamentsRequestInfo{}, errors.Wrap(err, fmt.Sprintf("object: %s", objectID.DebugString()))
	}

	if idx.Lifeline.LatestRequest == nil {
		return FilamentsRequestInfo{}, errors.New("latest request in lifeline is empty")
	}

	logger.Debugf("latest request from index %s", idx.Lifeline.LatestRequest.DebugString())

	cache := c.cache.Get(objectID)
	cache.Lock()
	defer cache.Unlock()

	// Select the min pulse to we can search for closed requests
	var minPulse insolar.PulseNumber
	if idx.Lifeline.EarliestOpenRequest != nil && *idx.Lifeline.EarliestOpenRequest <= requestID.GetPulseNumber() {
		minPulse = *idx.Lifeline.EarliestOpenRequest
	} else {
		minPulse = requestID.GetPulseNumber()
	}

	iter := newFetchingIterator(
		ctx,
		cache,
		objectID,
		*idx.Lifeline.LatestRequest,
		minPulse,
		c.jetFetcher,
		c.coordinator,
		c.sender,
	)

	var foundRequestInfo FilamentsRequestInfo
	foundRequestInfo.OldestMutable = true
	closedRequests := map[insolar.ID]struct{}{}

	isMutableRequest := func(virtual record.Record) bool {
		if in, ok := virtual.(*record.IncomingRequest); ok && !in.Immutable {
			return true
		}
		return false
	}

	for iter.HasPrev() {
		rec, err := iter.Prev(ctx)
		if err != nil {
			return FilamentsRequestInfo{}, errors.Wrap(err, "failed to calculate filament")
		}
		virtual := record.Unwrap(&rec.Record.Virtual)

		if rec.RecordID == requestID {
			foundRequestInfo.Request = &rec
			logger.Debugf("found request %s", rec.RecordID.DebugString())
			// if immutable found, we can don't need to search deeper
			if !isMutableRequest(virtual) {
				foundRequestInfo.OldestMutable = false
				return foundRequestInfo, nil
			}
			continue
		}

		if r, ok := virtual.(*record.Result); ok {
			closedRequests[*r.Request.GetLocal()] = struct{}{}

			if *r.Request.GetLocal() == requestID {
				foundRequestInfo.Result = &rec
				logger.Debugf("found result %s", rec.RecordID.DebugString())
			}
		}

		// request found, check if we have another opened mutable incoming older than that
		if foundRequestInfo.Request != nil {
			if _, ok := closedRequests[rec.RecordID]; !ok {
				if isMutableRequest(virtual) {
					foundRequestInfo.OldestMutable = false
					logger.Debugf("found oldest %s", rec.RecordID.DebugString())
				}
			}
		}
	}

	if foundRequestInfo.Request == nil {
		return FilamentsRequestInfo{},
			&payload.CodedError{
				Text: fmt.Sprintf("requestInfo not found request %s", requestID.DebugString()),
				Code: payload.CodeRequestNotFound,
			}
	}

	return foundRequestInfo, nil
}

func (c *FilamentCalculatorDefault) ClearIfLonger(limit int) {
	c.cache.DeleteIfLonger(limit)
}

func (c *FilamentCalculatorDefault) ClearAllExcept(ids []insolar.ID) {
	c.cache.DeleteAllExcept(ids)
}

func (c FilamentCalculatorDefault) findIndex(ctx context.Context, reasonID, requestID insolar.ID) (record.Index, error) {
	logger := inslogger.FromContext(ctx)

	logger.Debug("looking for index locally")
	idx, err := c.findIndexLocal(ctx, reasonID.Pulse(), requestID)
	if err == nil {
		return idx, nil
	}
	if err != object.ErrIndexNotFound {
		return record.Index{}, errors.Wrap(err, "failed to find index")
	}

	// Searching for the requests in the network
	// We need to be sure, that there is no duplicate of creationg request
	// INS-3607
	logger.Debug("looking for index on heavy")
	return c.findIndexHeavy(ctx, requestID, reasonID.Pulse())
}

func (c FilamentCalculatorDefault) findIndexHeavy(
	ctx context.Context, objID insolar.ID, readUntil insolar.PulseNumber,
) (record.Index, error) {
	node, err := c.coordinator.Heavy(ctx)
	if err != nil {
		return record.Index{}, errors.Wrap(err, "failed to check index origin")
	}

	msg, err := payload.NewMessage(&payload.SearchIndex{
		ObjectID: objID,
		Until:    readUntil,
	})
	if err != nil {
		return record.Index{}, errors.Wrap(err, "failed to create message")
	}

	reps, done := c.sender.SendTarget(ctx, msg, *node)
	defer done()

	res, ok := <-reps
	if !ok {
		return record.Index{}, errors.New("no reply")
	}

	pl, err := payload.UnmarshalFromMeta(res.Payload)
	if err != nil {
		return record.Index{}, errors.Wrap(err, "failed to unmarshal reply")
	}

	switch rep := pl.(type) {
	case *payload.SearchIndexInfo:
		if rep.Index == nil {
			return record.Index{}, object.ErrIndexNotFound
		}

		return *rep.Index, nil
	case *payload.Error:
		return record.Index{}, &payload.CodedError{
			Text: fmt.Sprint("failed to fetch index from heavy: ", rep.Text),
			Code: rep.Code,
		}
	default:
		return record.Index{}, fmt.Errorf("unexpected reply %T", pl)
	}
}

func (c *FilamentCalculatorDefault) findIndexLocal(
	ctx context.Context, until insolar.PulseNumber, requestID insolar.ID,
) (record.Index, error) {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"object_id": requestID.DebugString(),
		"until":     until,
	})
	iter := requestID.Pulse()
	logger.Debug("findIndex. start executing")
	for iter >= until {
		// We search for combination of id in a latest bucket
		// requestID.Pulse() is a latest, because findIndex is called only for Creationg requests
		idx, err := c.indexes.ForID(ctx, requestID.Pulse(), *insolar.NewID(iter, requestID.Hash()))
		if err != nil && err != object.ErrIndexNotFound {
			return record.Index{}, errors.Wrap(err, "failed to fetch index")
		}
		if err == nil {
			logger.Debug("findIndex. found:", idx.ObjID)
			return idx, nil
		}
		logger.Debug("findIndex. didn't find for:", iter)

		prev, err := c.pulses.Backwards(ctx, iter, 1)
		if err != nil {
			return record.Index{}, object.ErrIndexNotFound
		}

		iter = prev.PulseNumber
		logger.Debug("findIndex. next iter:", iter)
	}

	return record.Index{}, object.ErrIndexNotFound
}

type fetchingIterator struct {
	iter  filamentIterator
	cache *filamentCache

	objectID  insolar.ID
	readUntil insolar.PulseNumber

	jetFetcher  JetFetcher
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

func (c *cacheStore) DeleteIfLonger(limit int) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for id, cache := range c.caches {
		cache.RLock()
		cacheLen := len(cache.cache)
		cache.RUnlock()
		if cacheLen > limit {
			delete(c.caches, id)
		}
	}
}

func (c *cacheStore) DeleteAllExcept(ids []insolar.ID) {
	c.lock.Lock()
	defer c.lock.Unlock()

	newCaches := map[insolar.ID]*filamentCache{}
	for _, id := range ids {
		if cache, ok := c.caches[id]; ok {
			newCaches[id] = cache
		}
	}
	c.caches = newCaches
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

	defer stats.Record(ctx, statFilamentLength.M(1))

	composite, ok := i.cache.cache[*i.currentID]
	if ok {
		virtual := record.Unwrap(&composite.Meta.Virtual)
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
	virtual := record.Unwrap(&filamentRecord.Virtual)
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

type fetchIterator interface {
	PrevID() *insolar.ID
	HasPrev() bool
	Prev(ctx context.Context) (record.CompositeFilamentRecord, error)
}

func newFetchingIterator(
	ctx context.Context,
	cache *filamentCache,
	objectID, from insolar.ID,
	readUntil insolar.PulseNumber,
	fetcher JetFetcher,
	coordinator jet.Coordinator,
	sender bus.Sender,
) fetchIterator {
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

	if i.readUntil == 0 {
		return record.CompositeFilamentRecord{}, errors.New("invalid fetching parameters")
	}

	rec, err := i.iter.Prev(ctx)
	if err == nil {
		return rec, nil
	}

	if err != object.ErrNotFound {
		return record.CompositeFilamentRecord{}, errors.Wrap(err, "failed to fetch filament")
	}

	// Update cache from network.
	logger.Debug("fetching requests from network")
	recs, err := i.fetchFromNetwork(ctx, *i.PrevID())
	logger.Debug("received requests from network")
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

}

func (i *fetchingIterator) fetchFromNetwork(
	ctx context.Context, forID insolar.ID,
) ([]record.CompositeFilamentRecord, error) {
	ctx, span := instracer.StartSpan(ctx, "fetchingIterator.fetchFromNetwork")
	defer span.Finish()

	isBeyond, err := i.coordinator.IsBeyondLimit(ctx, forID.Pulse())
	if err != nil {
		instracer.AddError(span, err)
		return nil, errors.Wrap(err, "failed to calculate limit")
	}
	var node *insolar.Reference
	if isBeyond {
		node, err = i.coordinator.Heavy(ctx)
		if err != nil {
			instracer.AddError(span, err)
			return nil, errors.Wrap(err, "failed to calculate node")
		}
	} else {
		jetID, err := i.jetFetcher.Fetch(ctx, i.objectID, forID.Pulse())
		if err != nil {
			instracer.AddError(span, err)
			return nil, errors.Wrap(err, "failed to fetch jet")
		}
		node, err = i.coordinator.NodeForJet(ctx, *jetID, forID.Pulse())
		if err != nil {
			instracer.AddError(span, err)
			return nil, errors.Wrap(err, "failed to calculate node")
		}
	}
	if *node == i.coordinator.Me() {
		instracer.AddError(span, errors.New("tried to send message to self"))
		return nil, errors.New("tried to send message to self")
	}

	span.SetTag("objID", i.objectID.DebugString()).
		SetTag("startFrom", forID.DebugString()).
		SetTag("readUntil", i.readUntil.String())

	msg, err := payload.NewMessage(&payload.GetFilament{
		ObjectID:  i.objectID,
		StartFrom: forID,
		ReadUntil: i.readUntil,
	})
	if err != nil {
		instracer.AddError(span, err)
		return nil, errors.Wrap(err, "failed to create fetching message")
	}
	reps, done := i.sender.SendTarget(ctx, msg, *node)
	defer done()
	res, ok := <-reps
	if !ok {
		instracer.AddError(span, errors.New("no reply for filament fetch"))
		return nil, errors.New("no reply for filament fetch")
	}

	pl, err := payload.UnmarshalFromMeta(res.Payload)
	if err != nil {
		instracer.AddError(span, err)
		return nil, errors.Wrap(err, "failed to unmarshal reply")
	}
	switch p := pl.(type) {
	case *payload.FilamentSegment:
		stats.Record(ctx, statFilamentFetchedCount.M(int64(len(p.Records))))
		return p.Records, nil
	case *payload.Error:
		return nil, errors.New(p.Text)
	}
	instracer.AddError(span, fmt.Errorf("unexpected reply %T", pl))
	return nil, fmt.Errorf("unexpected reply %T", pl)
}
