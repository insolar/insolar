package executor_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilamentModifierDefault_SetRequest(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	ctx := inslogger.TestContext(t)

	var (
		pcs     insolar.PlatformCryptographyScheme
		indexes object.IndexStorage
		records object.RecordStorage
		manager *executor.FilamentModifierDefault
	)
	resetComponents := func() {
		pcs = testutils.NewPlatformCryptographyScheme()
		indexes = object.NewIndexStorageMemory()
		records = object.NewRecordMemory()
		manager = executor.NewFilamentModifier(indexes, records, pcs, nil)
	}

	objRef := gen.Reference()
	validRequest := record.Request{Object: &objRef}

	resetComponents()
	t.Run("object id is empty", func(t *testing.T) {
		err := manager.SetRequest(ctx, insolar.ID{}, gen.JetID(), validRequest)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("jet is not valid", func(t *testing.T) {
		err := manager.SetRequest(ctx, gen.ID(), insolar.JetID{}, validRequest)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("index does not exist", func(t *testing.T) {
		err := manager.SetRequest(ctx, gen.ID(), gen.JetID(), validRequest)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("request from the past", func(t *testing.T) {
		reqID := gen.ID()
		reqID.SetPulse(insolar.FirstPulseNumber + 1)
		latestPendingID := gen.ID()
		latestPendingID.SetPulse(insolar.FirstPulseNumber + 2)

		err := indexes.SetIndex(ctx, reqID.Pulse(), object.FilamentIndex{
			Lifeline: object.Lifeline{
				PendingPointer: &latestPendingID,
			},
		})
		require.NoError(t, err)

		err = manager.SetRequest(ctx, reqID, gen.JetID(), validRequest)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("happy basic", func(t *testing.T) {
		requestID := gen.ID()
		requestID.SetPulse(insolar.FirstPulseNumber + 2)
		latestPendingID := gen.ID()
		latestPendingID.SetPulse(insolar.FirstPulseNumber + 1)
		jetID := gen.JetID()

		err := indexes.SetIndex(ctx, requestID.Pulse(), object.FilamentIndex{
			ObjID: *validRequest.Object.Record(),
			Lifeline: object.Lifeline{
				PendingPointer: &latestPendingID,
			},
		})
		require.NoError(t, err)

		err = manager.SetRequest(ctx, requestID, jetID, validRequest)
		assert.NoError(t, err)

		idx, err := indexes.ForID(ctx, requestID.Pulse(), *validRequest.Object.Record())
		require.NoError(t, err)

		expectedFilamentRecord := record.PendingFilament{
			RecordID:       requestID,
			PreviousRecord: &latestPendingID,
		}
		virtual := record.Wrap(expectedFilamentRecord)
		hash := record.HashVirtual(pcs.ReferenceHasher(), virtual)
		expectedFilamentRecordID := *insolar.NewID(requestID.Pulse(), hash)

		assert.Equal(t, expectedFilamentRecordID, *idx.Lifeline.PendingPointer)
		assert.Equal(t, requestID.Pulse(), *idx.Lifeline.EarliestOpenRequest)

		rec, err := records.ForID(ctx, expectedFilamentRecordID)
		require.NoError(t, err)
		virtual = record.Wrap(expectedFilamentRecord)
		assert.Equal(t, record.Material{Virtual: &virtual, JetID: jetID}, rec)

		rec, err = records.ForID(ctx, requestID)
		require.NoError(t, err)
		virtual = record.Wrap(validRequest)
		assert.Equal(t, record.Material{Virtual: &virtual, JetID: jetID}, rec)

		mc.Finish()
	})
}

func TestFilamentModifierDefault_SetResult(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	ctx := inslogger.TestContext(t)

	var (
		pcs        insolar.PlatformCryptographyScheme
		indexes    object.IndexStorage
		records    object.RecordStorage
		calculator *executor.FilamentCalculatorMock
		manager    *executor.FilamentModifierDefault
	)
	resetComponents := func() {
		pcs = testutils.NewPlatformCryptographyScheme()
		indexes = object.NewIndexStorageMemory()
		records = object.NewRecordMemory()
		calculator = executor.NewFilamentCalculatorMock(mc)
		manager = executor.NewFilamentModifier(indexes, records, pcs, calculator)
	}

	validResult := record.Result{Object: gen.ID()}

	resetComponents()
	t.Run("object id is empty", func(t *testing.T) {
		err := manager.SetResult(ctx, insolar.ID{}, gen.JetID(), validResult)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("jet is not valid", func(t *testing.T) {
		err := manager.SetResult(ctx, gen.ID(), insolar.JetID{}, validResult)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("index does not exist", func(t *testing.T) {
		err := manager.SetResult(ctx, gen.ID(), gen.JetID(), validResult)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("happy basic", func(t *testing.T) {
		resultID := gen.ID()
		resultID.SetPulse(insolar.FirstPulseNumber + 2)
		latestPendingID := gen.ID()
		latestPendingID.SetPulse(insolar.FirstPulseNumber + 1)
		jetID := gen.JetID()

		expectedFilamentRecord := record.PendingFilament{
			RecordID:       resultID,
			PreviousRecord: &latestPendingID,
		}
		virtual := record.Wrap(expectedFilamentRecord)
		hash := record.HashVirtual(pcs.ReferenceHasher(), virtual)
		expectedFilamentRecordID := *insolar.NewID(resultID.Pulse(), hash)

		calculator.PendingRequestsFunc = func(_ context.Context, pn insolar.PulseNumber, id insolar.ID) ([]insolar.ID, error) {
			require.Equal(t, resultID.Pulse(), pn)
			require.Equal(t, validResult.Object, id)

			return []insolar.ID{expectedFilamentRecordID}, nil
		}

		latestPendingPulse := latestPendingID.Pulse()
		err := indexes.SetIndex(ctx, resultID.Pulse(), object.FilamentIndex{
			ObjID: validResult.Object,
			Lifeline: object.Lifeline{
				PendingPointer:      &latestPendingID,
				EarliestOpenRequest: &latestPendingPulse,
			},
		})
		require.NoError(t, err)

		err = manager.SetResult(ctx, resultID, jetID, validResult)
		assert.NoError(t, err)

		idx, err := indexes.ForID(ctx, resultID.Pulse(), validResult.Object)
		require.NoError(t, err)

		assert.Equal(t, expectedFilamentRecordID, *idx.Lifeline.PendingPointer)
		assert.Equal(t, resultID.Pulse(), *idx.Lifeline.EarliestOpenRequest)

		rec, err := records.ForID(ctx, expectedFilamentRecordID)
		require.NoError(t, err)
		virtual = record.Wrap(expectedFilamentRecord)
		assert.Equal(t, record.Material{Virtual: &virtual, JetID: jetID}, rec)

		rec, err = records.ForID(ctx, resultID)
		require.NoError(t, err)
		virtual = record.Wrap(validResult)
		assert.Equal(t, record.Material{Virtual: &virtual, JetID: jetID}, rec)

		mc.Finish()
	})

	resetComponents()
	t.Run("happy resets pending pointers in index", func(t *testing.T) {
		resultID := gen.ID()
		resultID.SetPulse(insolar.FirstPulseNumber + 2)
		latestPendingID := gen.ID()
		latestPendingID.SetPulse(insolar.FirstPulseNumber + 1)
		jetID := gen.JetID()

		calculator.PendingRequestsFunc = func(_ context.Context, pn insolar.PulseNumber, id insolar.ID) ([]insolar.ID, error) {
			require.Equal(t, resultID.Pulse(), pn)
			require.Equal(t, validResult.Object, id)

			return []insolar.ID{}, nil
		}

		latestPendingPulse := latestPendingID.Pulse()
		err := indexes.SetIndex(ctx, resultID.Pulse(), object.FilamentIndex{
			ObjID: validResult.Object,
			Lifeline: object.Lifeline{
				PendingPointer:      &latestPendingID,
				EarliestOpenRequest: &latestPendingPulse,
			},
		})
		require.NoError(t, err)

		err = manager.SetResult(ctx, resultID, jetID, validResult)
		assert.NoError(t, err)

		idx, err := indexes.ForID(ctx, resultID.Pulse(), validResult.Object)
		require.NoError(t, err)

		assert.Nil(t, idx.Lifeline.EarliestOpenRequest)

		mc.Finish()
	})
}

func TestFilamentCalculatorDefault_Requests(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	ctx := inslogger.TestContext(t)

	var (
		indexes    object.IndexStorage
		records    object.RecordStorage
		pcs        insolar.PlatformCryptographyScheme
		calculator *executor.FilamentCalculatorDefault
	)
	resetComponents := func() {
		indexes = object.NewIndexStorageMemory()
		records = object.NewRecordMemory()
		pcs = testutils.NewPlatformCryptographyScheme()
		calculator = executor.NewFilamentCalculator(indexes, records, nil, nil, nil)
	}

	resetComponents()
	t.Run("returns error if object does not exist", func(t *testing.T) {
		_, err := calculator.Requests(ctx, gen.ID(), gen.ID(), gen.PulseNumber(), gen.PulseNumber())
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("empty response", func(t *testing.T) {
		objectID := gen.ID()
		fromID := gen.ID()
		err := indexes.SetIndex(ctx, fromID.Pulse(), object.FilamentIndex{
			ObjID: objectID,
		})
		require.NoError(t, err)

		recs, err := calculator.Requests(ctx, objectID, fromID, gen.PulseNumber(), gen.PulseNumber())
		assert.NoError(t, err)
		assert.Equal(t, 0, len(recs))

		mc.Finish()
	})

	resetComponents()
	t.Run("happy basic", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records)
		rec1 := b.AppendRequest(insolar.FirstPulseNumber+1, record.Request{Nonce: rand.Uint64()})
		rec2 := b.AppendRequest(insolar.FirstPulseNumber+2, record.Request{Nonce: rand.Uint64()})
		rec3 := b.AppendRequest(insolar.FirstPulseNumber+2, record.Request{Nonce: rand.Uint64()})
		rec4 := b.AppendRequest(insolar.FirstPulseNumber+3, record.Request{Nonce: rand.Uint64()})
		b.AppendRequest(insolar.FirstPulseNumber+4, record.Request{Nonce: rand.Uint64()})

		objectID := gen.ID()
		fromID := rec4.MetaID
		earliestPending := rec1.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromID.Pulse(), object.FilamentIndex{
			ObjID: objectID,
			Lifeline: object.Lifeline{
				PendingPointer:      &rec4.MetaID,
				EarliestOpenRequest: &earliestPending,
			},
		})
		require.NoError(t, err)

		recs, err := calculator.Requests(ctx, objectID, fromID, rec2.MetaID.Pulse(), rec4.MetaID.Pulse())
		assert.NoError(t, err)
		require.Equal(t, 3, len(recs))
		assert.Equal(t, []record.CompositeFilamentRecord{rec4, rec3, rec2}, recs)

		mc.Finish()
	})
}

func TestFilamentCalculatorDefault_PendingRequests(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	ctx := inslogger.TestContext(t)

	var (
		indexes     object.IndexStorage
		records     object.RecordStorage
		coordinator *jet.CoordinatorMock
		fetcher     *jet.FetcherMock
		sender      *bus.SenderMock
		pcs         insolar.PlatformCryptographyScheme
		calculator  *executor.FilamentCalculatorDefault
	)
	resetComponents := func() {
		indexes = object.NewIndexStorageMemory()
		records = object.NewRecordMemory()
		coordinator = jet.NewCoordinatorMock(mc)
		fetcher = jet.NewFetcherMock(mc)
		sender = bus.NewSenderMock(mc)
		pcs = testutils.NewPlatformCryptographyScheme()
		calculator = executor.NewFilamentCalculator(indexes, records, coordinator, fetcher, sender)
	}

	resetComponents()
	t.Run("returns error if object does not exist", func(t *testing.T) {
		_, err := calculator.PendingRequests(ctx, gen.PulseNumber(), gen.ID())
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("empty response", func(t *testing.T) {
		objectID := gen.ID()
		fromPulse := gen.PulseNumber()
		err := indexes.SetIndex(ctx, fromPulse, object.FilamentIndex{
			ObjID: objectID,
		})
		require.NoError(t, err)

		recs, err := calculator.PendingRequests(ctx, fromPulse, objectID)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(recs))

		mc.Finish()
	})
}

type filamentBuilder struct {
	records   object.RecordModifier
	currentID insolar.ID
	ctx       context.Context
	pcs       insolar.PlatformCryptographyScheme
}

func newFilamentBuilder(
	ctx context.Context,
	pcs insolar.PlatformCryptographyScheme,
	records object.RecordModifier,
) *filamentBuilder {
	return &filamentBuilder{
		ctx:     ctx,
		records: records,
		pcs:     pcs,
	}
}

func (b *filamentBuilder) AppendRequest(pn insolar.PulseNumber, request record.Request) record.CompositeFilamentRecord {
	var composite record.CompositeFilamentRecord
	{
		virtual := record.Wrap(request)
		hash := record.HashVirtual(b.pcs.ReferenceHasher(), virtual)
		id := *insolar.NewID(pn, hash)
		material := record.Material{Virtual: &virtual, JetID: insolar.ZeroJetID}
		err := b.records.Set(b.ctx, id, material)
		if err != nil {
			panic(err)
		}
		composite.RecordID = id
		composite.Record = material
	}

	{
		rec := record.PendingFilament{RecordID: composite.RecordID}
		if !b.currentID.IsEmpty() {
			curr := b.currentID
			rec.PreviousRecord = &curr
		}
		virtual := record.Wrap(rec)
		hash := record.HashVirtual(b.pcs.ReferenceHasher(), virtual)
		id := *insolar.NewID(pn, hash)
		material := record.Material{Virtual: &virtual, JetID: insolar.ZeroJetID}
		err := b.records.Set(b.ctx, id, material)
		if err != nil {
			panic(err)
		}
		composite.MetaID = id
		composite.Meta = material
	}

	b.currentID = composite.MetaID

	return composite
}
