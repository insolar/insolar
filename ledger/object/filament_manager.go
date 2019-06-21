package object

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/insolar"
	buswm "github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/network/storage"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

type FilamentModifier interface {
	SetRequest(ctx context.Context, reqID insolar.ID, jetID insolar.JetID, request record.Request) error
	SetResult(ctx context.Context, resID insolar.ID, jetID insolar.JetID, result record.Result) error
}

func NewFilamentManager() *FilamentManager {
	return &FilamentManager{
		cache: newCacheStore(),
	}
}

type FilamentManager struct {
	cache *cacheStore

	idxAccessor     IndexAccessor
	idxModifier     IndexModifier
	idLocker        IDLocker
	recordStorage   RecordStorage
	coordinator     jet.Coordinator
	pcs             insolar.PlatformCryptographyScheme
	pulseCalculator storage.PulseCalculator
	bus             insolar.MessageBus
	busWM           buswm.Sender
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
		return ErrLifelineNotFound
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

	stats.Record(ctx,
		statObjectPendingRequestsInMemoryAddedCount.M(int64(1)),
	)

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
		return ErrLifelineNotFound
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

	stats.Record(ctx,
		statObjectPendingResultsInMemoryAddedCount.M(int64(1)),
	)

	return nil
}

func (m *FilamentManager) calculatePending(
	ctx context.Context,
	pulse insolar.PulseNumber,
	objectID insolar.ID,
	idx FilamentIndex,
) ([]insolar.ID, error) {
	incoming := m.cache.Get(objectID).Incoming
	err := m.fetchFilament(
		ctx,
		incoming,
		pulse,
		objectID,
		*idx.Lifeline.PendingPointer,
		*idx.Lifeline.EarliestOpenRequest,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch pending")
	}

	var pending []insolar.ID
	hasResult := map[insolar.ID]struct{}{}
	iter := idx.Lifeline.PendingPointer
	for iter != nil && iter.Pulse() >= *idx.Lifeline.EarliestOpenRequest {
		recs, err := incoming.Get(ctx, *idx.Lifeline.PendingPointer)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch pending")
		}

		if len(recs) == 0 {
			return nil, nil
		}

		// Going in reverse to find result for request before the request. This way we preserve pending order.
		for i := len(recs) - 1; i >= 0; i-- {
			virtual := record.Unwrap(recs[i].Record.Virtual)
			switch r := virtual.(type) {
			case *record.Request:
				if _, ok := hasResult[recs[i].RecordID]; !ok {
					pending = append(pending, recs[i].RecordID)
				}
			case *record.Result:
				hasResult[*r.Request.Record()] = struct{}{}
			}
		}

		virtual := record.Unwrap(recs[len(recs)-1].Meta.Virtual)
		filament, ok := virtual.(*record.PendingFilament)
		if !ok {
			return nil, errors.New("failed to convert filament")
		}

		// Jumping to the next pulse.
		iter = filament.PreviousRecord
	}

	// We need to reverse pending because we iterated from the end when selected them.
	ordered := make([]insolar.ID, len(pending))
	count := len(pending)
	for i, id := range pending {
		ordered[count-i] = id
	}

	return ordered, nil
}

func (m *FilamentManager) fetchFilament(
	ctx context.Context,
	cache *filamentCache,
	pulse insolar.PulseNumber,
	objectID, from insolar.ID,
	until insolar.PulseNumber,
) error {
	fetchFromNetwork := func(forID insolar.ID) ([]record.CompositeFilamentRecord, error) {
		jetID, err := m.jetFetcher.Fetch(ctx, objectID, forID.Pulse())
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch jet")
		}
		node, err := m.coordinator.NodeForJet(ctx, *jetID, pulse, forID.Pulse())
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate node")
		}
		if *node == m.coordinator.Me() {
			return nil, errors.New("tried to send message to self")
		}

		msg, err := payload.NewMessage(&payload.GetFilament{
			ObjectID:  objectID,
			StartFrom: forID,
			ReadUntil: until,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create fetching message")
		}
		reps, done := m.busWM.SendTarget(ctx, msg, *node)
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

	iter := &from
	for iter != nil && iter.Pulse() <= until {
		recs, err := cache.Get(ctx, from)
		if err != nil {
			if err == ErrNotFound {
				recs, err := fetchFromNetwork(from)
				if err != nil {
					return errors.Wrap(err, "failed to fetch filament")
				}
				saveToCache(cache, recs)
			}
			return errors.Wrap(err, "failed to fetch filament")
		}
		if len(recs) == 0 {
			return nil
		}

		virtual := record.Unwrap(recs[len(recs)-1].Meta.Virtual)
		filament, ok := virtual.(*record.PendingFilament)
		if !ok {
			return errors.New("failed to convert filament record")
		}

		iter = filament.PreviousRecord
	}

	return nil
}

func saveToCache(cache *filamentCache, recs []record.CompositeFilamentRecord) {
	if len(recs) == 0 {
		return
	}

	tail := 0
	iterPN := recs[0].MetaID.Pulse()
	for i, rec := range recs {
		if rec.MetaID.Pulse() != iterPN {
			cache.Set(iterPN, recs[tail:i-tail])
			tail = i
			iterPN = rec.MetaID.Pulse()
		}
	}
}
