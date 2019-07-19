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

package executor_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilamentModifierDefault_SetRequest(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	ctx := inslogger.TestContext(t)

	var (
		pcs        insolar.PlatformCryptographyScheme
		indexes    object.IndexStorage
		records    object.RecordStorage
		manager    *executor.FilamentModifierDefault
		calculator *executor.FilamentCalculatorMock
	)
	resetComponents := func() {
		pcs = testutils.NewPlatformCryptographyScheme()
		indexes = object.NewIndexStorageMemory()
		records = object.NewRecordMemory()
		calculator = executor.NewFilamentCalculatorMock(t)
		manager = executor.NewFilamentModifier(indexes, records, pcs, calculator, nil)
	}

	objRef := gen.Reference()
	validRequest := record.IncomingRequest{Object: &objRef, CallType: record.CTMethod, Reason: *insolar.NewReference(*objRef.Record())}

	resetComponents()
	t.Run("object id is empty", func(t *testing.T) {
		_, _, err := manager.SetRequest(ctx, insolar.ID{}, gen.JetID(), &validRequest)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("jet is not valid", func(t *testing.T) {
		_, _, err := manager.SetRequest(ctx, gen.ID(), insolar.JetID{}, &validRequest)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("index does not exist", func(t *testing.T) {
		_, _, err := manager.SetRequest(ctx, gen.ID(), gen.JetID(), &validRequest)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("request from the past", func(t *testing.T) {
		reqID := gen.ID()
		reqID.SetPulse(insolar.FirstPulseNumber + 1)
		latestPendingID := gen.ID()
		latestPendingID.SetPulse(insolar.FirstPulseNumber + 2)

		err := indexes.SetIndex(ctx, reqID.Pulse(), record.Index{
			Lifeline: record.Lifeline{
				PendingPointer: &latestPendingID,
			},
		})
		require.NoError(t, err)

		_, _, err = manager.SetRequest(ctx, reqID, gen.JetID(), &validRequest)
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

		err := indexes.SetIndex(ctx, requestID.Pulse(), record.Index{
			ObjID: *validRequest.Object.Record(),
			Lifeline: record.Lifeline{
				PendingPointer: &latestPendingID,
			},
		})
		require.NoError(t, err)

		calculator.RequestDuplicateFunc = func(_ context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 insolar.ID, p4 record.Request) (r *record.CompositeFilamentRecord, r1 *record.CompositeFilamentRecord, r2 error) {
			require.Equal(t, p1, requestID.Pulse())
			require.Equal(t, p2, *validRequest.Object.Record())
			require.Equal(t, p3, requestID)
			return nil, nil, nil
		}

		_, _, err = manager.SetRequest(ctx, requestID, jetID, &validRequest)
		require.NoError(t, err)

		idx, err := indexes.ForID(ctx, requestID.Pulse(), *validRequest.Object.Record())
		require.NoError(t, err)

		expectedFilamentRecord := record.PendingFilament{
			RecordID:       requestID,
			PreviousRecord: &latestPendingID,
		}
		virtual := record.Wrap(expectedFilamentRecord)
		hash := record.HashVirtual(pcs.ReferenceHasher(), virtual)
		expectedFilamentRecordID := *insolar.NewID(requestID.Pulse(), hash)

		require.Equal(t, expectedFilamentRecordID, *idx.Lifeline.PendingPointer)
		require.Equal(t, requestID.Pulse(), *idx.Lifeline.EarliestOpenRequest)

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

func TestFilamentModifierDefault_SetRequest_NewObject(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	ctx := inslogger.TestContext(t)

	var (
		pcs        insolar.PlatformCryptographyScheme
		indexes    object.IndexStorage
		records    object.RecordStorage
		manager    *executor.FilamentModifierDefault
		calculator *executor.FilamentCalculatorMock
		pulses     *pulse.CalculatorMock
	)
	resetComponents := func() {
		pcs = testutils.NewPlatformCryptographyScheme()
		indexes = object.NewIndexStorageMemory()
		records = object.NewRecordMemory()
		calculator = executor.NewFilamentCalculatorMock(t)
		pulses = pulse.NewCalculatorMock(t)
		manager = executor.NewFilamentModifier(indexes, records, pcs, calculator, pulses)
	}

	resetComponents()
	t.Run("request from the past", func(t *testing.T) {
		reqID := gen.ID()
		reqID.SetPulse(insolar.FirstPulseNumber + 1)

		objRef := gen.Reference()
		validRequest := record.IncomingRequest{
			Object:   &objRef,
			Reason:   *insolar.NewReference(reqID),
			CallType: record.CTSaveAsChild,
		}

		latestPendingID := gen.ID()
		latestPendingID.SetPulse(insolar.FirstPulseNumber + 2)

		err := indexes.SetIndex(ctx, reqID.Pulse(), record.Index{
			ObjID: reqID,
			Lifeline: record.Lifeline{
				PendingPointer: &latestPendingID,
			},
		})
		require.NoError(t, err)

		_, _, err = manager.SetRequest(ctx, reqID, gen.JetID(), &validRequest)
		require.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("happy basic. new object", func(t *testing.T) {
		requestID := gen.ID()
		requestID.SetPulse(insolar.FirstPulseNumber + 2)
		jetID := gen.JetID()

		objRef := gen.Reference()
		validRequest := record.IncomingRequest{
			Object:   &objRef,
			Reason:   *insolar.NewReference(requestID),
			CallType: record.CTSaveAsChild,
		}
		pulses.BackwardsMock.Return(*insolar.GenesisPulse, errors.New("stub error"))
		calculator.RequestDuplicateMock.Return(nil, nil, nil)

		_, _, err := manager.SetRequest(ctx, requestID, jetID, &validRequest)
		assert.NoError(t, err)

		idx, err := indexes.ForID(ctx, requestID.Pulse(), requestID)
		require.NoError(t, err)

		expectedFilamentRecord := record.PendingFilament{
			RecordID: requestID,
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

	resetComponents()
	t.Run("happy basic. existed object", func(t *testing.T) {
		requestID := gen.ID()
		requestID.SetPulse(insolar.FirstPulseNumber + 2)
		jetID := gen.JetID()

		objRef := gen.Reference()
		validRequest := record.IncomingRequest{
			Object:   &objRef,
			Reason:   *insolar.NewReference(requestID),
			CallType: record.CTSaveAsChild,
		}
		pulses.BackwardsMock.Return(*insolar.GenesisPulse, errors.New("stub error"))
		calculator.RequestDuplicateFunc = func(_ context.Context, _ insolar.PulseNumber, _ insolar.ID, _ insolar.ID, _ record.Request) (*record.CompositeFilamentRecord, *record.CompositeFilamentRecord, error) {
			return nil, nil, nil
		}

		req, res, err := manager.SetRequest(ctx, requestID, jetID, &validRequest)
		require.NoError(t, err)
		require.Nil(t, req)
		require.Nil(t, res)
		_, err = indexes.ForID(ctx, requestID.Pulse(), requestID)
		require.NoError(t, err)

		calculator.RequestDuplicateFunc = func(_ context.Context, _ insolar.PulseNumber, objID insolar.ID, reqID insolar.ID, _ record.Request) (r *record.CompositeFilamentRecord, r1 *record.CompositeFilamentRecord, r2 error) {
			require.Equal(t, requestID, objID)
			require.Equal(t, requestID, reqID)

			return &record.CompositeFilamentRecord{
				RecordID: requestID,
			}, nil, nil
		}

		req, res, err = manager.SetRequest(ctx, requestID, jetID, &validRequest)
		require.NoError(t, err)
		require.Nil(t, res)
		require.NotNil(t, req)
		require.Equal(t, requestID, req.RecordID)

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
		manager = executor.NewFilamentModifier(indexes, records, pcs, calculator, nil)
	}

	validResult := record.Result{Object: gen.ID()}

	resetComponents()
	t.Run("object id is empty", func(t *testing.T) {
		_, err := manager.SetResult(ctx, insolar.ID{}, gen.JetID(), validResult)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("jet is not valid", func(t *testing.T) {
		_, err := manager.SetResult(ctx, gen.ID(), insolar.JetID{}, validResult)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("index does not exist", func(t *testing.T) {
		_, err := manager.SetResult(ctx, gen.ID(), gen.JetID(), validResult)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("problems during duplicates searching", func(t *testing.T) {
		calculator.ResultDuplicateMock.Return(nil, errors.New("smth terrible happened"))
		resultID := gen.ID()
		err := indexes.SetIndex(ctx, resultID.Pulse(), record.Index{
			ObjID:    validResult.Object,
			Lifeline: record.Lifeline{},
		})
		require.NoError(t, err)

		_, err = manager.SetResult(ctx, resultID, gen.JetID(), validResult)
		require.Error(t, err)
		require.Contains(t, err.Error(), "smth terrible happened")

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
		calculator.ResultDuplicateFunc = func(_ context.Context, inPN insolar.PulseNumber, inObjID insolar.ID, inResID insolar.ID, inRes record.Result) (*record.CompositeFilamentRecord, error) {
			require.Equal(t, inPN, resultID.Pulse())
			require.Equal(t, validResult.Object, inObjID)
			require.Equal(t, resultID, inResID)
			return nil, nil
		}

		latestPendingPulse := latestPendingID.Pulse()
		err := indexes.SetIndex(ctx, resultID.Pulse(), record.Index{
			ObjID: validResult.Object,
			Lifeline: record.Lifeline{
				PendingPointer:      &latestPendingID,
				EarliestOpenRequest: &latestPendingPulse,
			},
		})
		require.NoError(t, err)

		_, err = manager.SetResult(ctx, resultID, jetID, validResult)
		assert.NoError(t, err)

		idx, err := indexes.ForID(ctx, resultID.Pulse(), validResult.Object)
		require.NoError(t, err)

		require.Equal(t, expectedFilamentRecordID, *idx.Lifeline.PendingPointer)
		require.Equal(t, resultID.Pulse(), *idx.Lifeline.EarliestOpenRequest)

		rec, err := records.ForID(ctx, expectedFilamentRecordID)
		require.NoError(t, err)
		virtual = record.Wrap(expectedFilamentRecord)
		require.Equal(t, record.Material{Virtual: &virtual, JetID: jetID}, rec)

		rec, err = records.ForID(ctx, resultID)
		require.NoError(t, err)
		virtual = record.Wrap(validResult)
		require.Equal(t, record.Material{Virtual: &virtual, JetID: jetID}, rec)

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
		calculator.ResultDuplicateFunc = func(_ context.Context, pn insolar.PulseNumber, objID insolar.ID, inResID insolar.ID, _ record.Result) (*record.CompositeFilamentRecord, error) {
			require.Equal(t, resultID, inResID)
			require.Equal(t, validResult.Object, objID)
			require.Equal(t, resultID.Pulse(), pn)

			return nil, nil
		}

		latestPendingPulse := latestPendingID.Pulse()
		err := indexes.SetIndex(ctx, resultID.Pulse(), record.Index{
			ObjID: validResult.Object,
			Lifeline: record.Lifeline{
				PendingPointer:      &latestPendingID,
				EarliestOpenRequest: &latestPendingPulse,
			},
		})
		require.NoError(t, err)

		_, err = manager.SetResult(ctx, resultID, jetID, validResult)
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
		records    *object.RecordMemory
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
		err := indexes.SetIndex(ctx, fromID.Pulse(), record.Index{
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
		b := newFilamentBuilder(ctx, pcs, records, indexes)
		storageRecs := make([]record.CompositeFilamentRecord, 5)
		storageRecs[0] = b.Append(insolar.FirstPulseNumber+1, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})
		storageRecs[1] = b.Append(insolar.FirstPulseNumber+2, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})
		storageRecs[2] = b.Append(insolar.FirstPulseNumber+2, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})
		storageRecs[3] = b.Append(insolar.FirstPulseNumber+3, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})
		storageRecs[4] = b.Append(insolar.FirstPulseNumber+4, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})

		objectID := gen.ID()
		fromID := storageRecs[3].MetaID
		earliestPending := storageRecs[0].MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromID.Pulse(), record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer:      &storageRecs[3].MetaID,
				EarliestOpenRequest: &earliestPending,
			},
		})
		require.NoError(t, err)

		// First time, records accessed from storage.
		recs, err := calculator.Requests(ctx, objectID, fromID, storageRecs[1].MetaID.Pulse(), storageRecs[3].MetaID.Pulse())
		assert.NoError(t, err)
		require.Equal(t, 3, len(recs))
		assert.Equal(t, []record.CompositeFilamentRecord{storageRecs[3], storageRecs[2], storageRecs[1]}, recs)

		// Second time storage is cleared. Records are accessed from cache.
		for _, rec := range storageRecs {
			records.DeleteForPN(ctx, rec.MetaID.Pulse())
		}
		recs, err = calculator.Requests(ctx, objectID, fromID, storageRecs[1].MetaID.Pulse(), storageRecs[3].MetaID.Pulse())
		assert.NoError(t, err)
		require.Equal(t, 3, len(recs))
		assert.Equal(t, []record.CompositeFilamentRecord{storageRecs[3], storageRecs[2], storageRecs[1]}, recs)

		mc.Finish()
	})
}

func TestFilamentCalculatorDefault_PendingRequests_RequestOnly(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	c := newComponents(mc)
	b := newFilamentBuilder(ctx, c.pcs, c.records, c.indexes)

	inRequest1 := b.appendInRequest(pulseNum(1))
	inRequestID1 := inRequest1.RecordID

	objectID, err := b.setIndex(ctx)
	require.NoError(t, err)

	recs, err := c.calculator.PendingRequests(ctx, *b.earliestOpenRequest, objectID)
	require.NoError(t, err)
	require.Equal(t, []insolar.ID{inRequestID1}, recs, "request in pendings")

	mc.Finish()
}

func TestFilamentCalculatorDefault_PendingRequests_RequestWithDetached(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	c := newComponents(mc)
	b := newFilamentBuilder(ctx, c.pcs, c.records, c.indexes)

	inRequest1 := b.appendInRequest(pulseNum(1))
	inRequestID1 := inRequest1.RecordID

	outRequestDetached1 := b.appendOutRequest(pulseNum(2), inRequestID1, record.ReturnSaga)
	outRequestDetachedID1 := outRequestDetached1.RecordID

	objectID, err := b.setIndex(ctx)
	require.NoError(t, err)

	recs, err := c.calculator.PendingRequests(ctx, *b.earliestOpenRequest, objectID)
	require.NoError(t, err)
	require.Equal(t, []insolar.ID{inRequestID1, outRequestDetachedID1}, recs,
		"single unclosed request with detached outgoing (both are pendings)")

	mc.Finish()
}

func TestFilamentCalculatorDefault_PendingRequests_RequestWithNotDetached(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	c := newComponents(mc)
	b := newFilamentBuilder(ctx, c.pcs, c.records, c.indexes)

	inRequest1 := b.appendInRequest(pulseNum(1))
	inRequestID1 := inRequest1.RecordID

	_ = b.appendOutRequest(pulseNum(2), inRequestID1, record.ReturnResult)

	objectID, err := b.setIndex(ctx)
	require.NoError(t, err)

	recs, err := c.calculator.PendingRequests(ctx, *b.earliestOpenRequest, objectID)
	require.NoError(t, err)
	require.Equal(t, []insolar.ID{inRequestID1}, recs,
		"single unclosed request with not detached outgoing (request is only pending)")

	mc.Finish()
}

func TestFilamentCalculatorDefault_RequestWithDetachedAndResult(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	c := newComponents(mc)
	b := newFilamentBuilder(ctx, c.pcs, c.records, c.indexes)

	inRequest1 := b.appendInRequest(pulseNum(1))
	inRequestID1 := inRequest1.RecordID
	_ = b.appendOutRequest(pulseNum(2), inRequestID1, record.ReturnSaga)
	_ = b.appendResult(pulseNum(2), inRequestID1)

	objectID, err := b.setIndex(ctx)
	require.NoError(t, err)

	recs, err := c.calculator.PendingRequests(ctx, *b.earliestOpenRequest, objectID)
	require.NoError(t, err)
	require.Equal(t, []insolar.ID{}, recs, "in-request with result detached out-request, as result no pendings")

	mc.Finish()
}

func TestFilamentCalculatorDefault_RequestWithResult(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	c := newComponents(mc)
	b := newFilamentBuilder(ctx, c.pcs, c.records, c.indexes)

	inRequest1 := b.appendInRequest(pulseNum(1))
	_ = b.appendResult(pulseNum(2), inRequest1.RecordID)

	objectID, err := b.setIndex(ctx)
	require.NoError(t, err)

	recs, err := c.calculator.PendingRequests(ctx, *b.earliestOpenRequest, objectID)
	require.NoError(t, err)
	require.Equal(t, 0, len(recs), "no pendings if request with result")

	mc.Finish()
}

func TestFilamentCalculatorDefault2_RequestWithResultAndNotDetachedOutgoing(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	c := newComponents(mc)
	b := newFilamentBuilder(ctx, c.pcs, c.records, c.indexes)

	inRequest1 := b.appendInRequest(pulseNum(1))
	inRequestID1 := inRequest1.RecordID
	_ = b.appendOutRequest(pulseNum(2), inRequestID1, record.ReturnResult)
	_ = b.appendResult(pulseNum(2), inRequestID1)

	objectID, err := b.setIndex(ctx)
	require.NoError(t, err)

	recs, err := c.calculator.PendingRequests(ctx, *b.earliestOpenRequest, objectID)
	require.NoError(t, err)
	require.Equal(t, 0, len(recs), "no pendings")

	mc.Finish()
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
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
		})
		require.NoError(t, err)

		recs, err := calculator.PendingRequests(ctx, fromPulse, objectID)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(recs))

		mc.Finish()
	})

	resetComponents()
	t.Run("happy basic", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records, indexes)
		rec1 := b.Append(insolar.FirstPulseNumber+1, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})
		rec2 := b.Append(insolar.FirstPulseNumber+2, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})
		b.Append(insolar.FirstPulseNumber+3, record.Result{Request: *insolar.NewReference(rec1.RecordID)})
		rec4 := b.Append(insolar.FirstPulseNumber+3, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})
		b.Append(insolar.FirstPulseNumber+4, record.IncomingRequest{Nonce: rand.Uint64(), CallType: record.CTMethod})

		objectID := gen.ID()
		fromPulse := rec4.MetaID.Pulse()
		earliestPending := rec1.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer:      &rec4.MetaID,
				EarliestOpenRequest: &earliestPending,
			},
		})
		require.NoError(t, err)

		recs, err := calculator.PendingRequests(ctx, fromPulse, objectID)
		assert.NoError(t, err)
		require.Equal(t, 2, len(recs))
		assert.Equal(t, []insolar.ID{rec2.RecordID, rec4.RecordID}, recs)

		mc.Finish()
	})

	resetComponents()
	t.Run("happy fetches from light", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records, indexes)
		rec1 := b.Append(insolar.FirstPulseNumber+1, record.IncomingRequest{Nonce: rand.Uint64()})
		rec2 := b.Append(insolar.FirstPulseNumber+2, record.IncomingRequest{Nonce: rand.Uint64()})
		// This result is not in the storage.
		missingRec := b.AppendNoPersist(insolar.FirstPulseNumber+3, record.Result{Request: *insolar.NewReference(rec1.RecordID)})
		rec4 := b.Append(insolar.FirstPulseNumber+4, record.IncomingRequest{Nonce: rand.Uint64()})
		b.Append(insolar.FirstPulseNumber+5, record.IncomingRequest{Nonce: rand.Uint64()})

		objectID := gen.ID()
		fromPulse := rec4.MetaID.Pulse()
		earliestPending := rec1.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer:      &rec4.MetaID,
				EarliestOpenRequest: &earliestPending,
			},
		})
		require.NoError(t, err)

		coordinator.IsBeyondLimitFunc = func(_ context.Context, target insolar.PulseNumber) (bool, error) {
			require.Equal(t, missingRec.MetaID.Pulse(), target)
			return false, nil
		}

		jetID := gen.JetID()
		fetcher.FetchFunc = func(_ context.Context, targetID insolar.ID, pn insolar.PulseNumber) (*insolar.ID, error) {
			require.Equal(t, objectID, targetID)
			require.Equal(t, missingRec.MetaID.Pulse(), pn)
			id := insolar.ID(jetID)
			return &id, nil
		}

		node := gen.Reference()
		coordinator.NodeForJetFunc = func(_ context.Context, jet insolar.ID, target insolar.PulseNumber) (*insolar.Reference, error) {
			require.Equal(t, insolar.ID(jetID), jet)
			require.Equal(t, missingRec.MetaID.Pulse(), target)
			return &node, nil
		}

		coordinator.MeMock.Return(node)

		recs, err := calculator.PendingRequests(ctx, fromPulse, objectID)
		assert.Error(t, err, "returns error if trying to fetch from self")

		coordinator.MeMock.Return(gen.Reference())

		sender.SendTargetFunc = func(_ context.Context, msg *message.Message, target insolar.Reference) (<-chan *message.Message, func()) {
			pl, err := payload.Unmarshal(msg.Payload)
			require.NoError(t, err)

			getFilament, ok := pl.(*payload.GetFilament)
			require.True(t, ok)

			require.Equal(t, objectID, getFilament.ObjectID)
			require.Equal(t, missingRec.MetaID, getFilament.StartFrom)
			require.Equal(t, earliestPending, getFilament.ReadUntil)

			require.NoError(t, err)
			respMsg, err := payload.NewMessage(&payload.FilamentSegment{
				ObjectID: objectID,
				Records:  []record.CompositeFilamentRecord{missingRec},
			})
			require.NoError(t, err)
			meta := payload.Meta{Payload: respMsg.Payload}
			buf, err := meta.Marshal()
			require.NoError(t, err)
			respMsg.Payload = buf
			ch := make(chan *message.Message, 1)
			ch <- respMsg
			return ch, func() {}
		}

		recs, err = calculator.PendingRequests(ctx, fromPulse, objectID)
		assert.NoError(t, err)
		require.Equal(t, 2, len(recs))
		assert.Equal(t, []insolar.ID{rec2.RecordID, rec4.RecordID}, recs)

		mc.Finish()
	})

	resetComponents()
	t.Run("happy fetches from heavy", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records, indexes)
		rec1 := b.Append(insolar.FirstPulseNumber+1, record.IncomingRequest{Nonce: rand.Uint64()})
		rec2 := b.Append(insolar.FirstPulseNumber+2, record.IncomingRequest{Nonce: rand.Uint64()})
		// This result is not in the storage.
		missingRec := b.AppendNoPersist(insolar.FirstPulseNumber+3, record.Result{Request: *insolar.NewReference(rec1.RecordID)})
		rec4 := b.Append(insolar.FirstPulseNumber+4, record.IncomingRequest{Nonce: rand.Uint64()})
		b.Append(insolar.FirstPulseNumber+5, record.IncomingRequest{Nonce: rand.Uint64()})

		objectID := gen.ID()
		fromPulse := rec4.MetaID.Pulse()
		earliestPending := rec1.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer:      &rec4.MetaID,
				EarliestOpenRequest: &earliestPending,
			},
		})
		require.NoError(t, err)

		coordinator.IsBeyondLimitFunc = func(_ context.Context, target insolar.PulseNumber) (bool, error) {
			require.Equal(t, missingRec.MetaID.Pulse(), target)
			return true, nil
		}

		node := gen.Reference()
		coordinator.HeavyFunc = func(_ context.Context) (*insolar.Reference, error) {
			return &node, nil
		}
		coordinator.MeMock.Return(node)

		recs, err := calculator.PendingRequests(ctx, fromPulse, objectID)
		assert.Error(t, err, "returns error if trying to fetch from self")

		coordinator.MeMock.Return(gen.Reference())

		sender.SendTargetFunc = func(_ context.Context, msg *message.Message, target insolar.Reference) (<-chan *message.Message, func()) {
			pl, err := payload.Unmarshal(msg.Payload)
			require.NoError(t, err)

			getFilament, ok := pl.(*payload.GetFilament)
			require.True(t, ok)

			require.Equal(t, objectID, getFilament.ObjectID)
			require.Equal(t, missingRec.MetaID, getFilament.StartFrom)
			require.Equal(t, earliestPending, getFilament.ReadUntil)

			require.NoError(t, err)
			respMsg, err := payload.NewMessage(&payload.FilamentSegment{
				ObjectID: objectID,
				Records:  []record.CompositeFilamentRecord{missingRec},
			})
			require.NoError(t, err)
			meta := payload.Meta{Payload: respMsg.Payload}
			buf, err := meta.Marshal()
			require.NoError(t, err)
			respMsg.Payload = buf
			ch := make(chan *message.Message, 1)
			ch <- respMsg
			return ch, func() {}
		}

		recs, err = calculator.PendingRequests(ctx, fromPulse, objectID)
		assert.NoError(t, err)
		require.Equal(t, 2, len(recs))
		assert.Equal(t, []insolar.ID{rec2.RecordID, rec4.RecordID}, recs)

		mc.Finish()
	})
}

func TestFilamentCalculatorDefault_ResultDuplicate(t *testing.T) {
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
	t.Run("returns error if reason is empty", func(t *testing.T) {
		_, err := calculator.ResultDuplicate(ctx, gen.PulseNumber(), gen.ID(), gen.ID(), record.Result{})
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("no records", func(t *testing.T) {
		objectID := gen.ID()
		fromPulse := gen.PulseNumber()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
		})
		require.NoError(t, err)

		res, err := calculator.ResultDuplicate(ctx, fromPulse, objectID, gen.ID(), record.Result{Request: gen.Reference()})

		assert.NoError(t, err)
		assert.Nil(t, res)

		mc.Finish()
	})

	resetComponents()
	t.Run("returns result. result duplicate is found", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records, indexes)
		req := record.IncomingRequest{Nonce: rand.Uint64(), Reason: *insolar.NewReference(*insolar.NewID(insolar.FirstPulseNumber, nil))}
		req1 := b.Append(insolar.FirstPulseNumber+1, req)
		res := record.Result{Request: *insolar.NewReference(req1.RecordID)}
		res1 := b.Append(insolar.FirstPulseNumber+2, res)

		objectID := gen.ID()
		fromPulse := res1.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer: &res1.MetaID,
			},
		})
		require.NoError(t, err)

		fRes, err := calculator.ResultDuplicate(ctx, fromPulse, objectID, res1.RecordID, res)
		require.NoError(t, err)
		require.Equal(t, *fRes, res1)

		mc.Finish()
	})

	resetComponents()
	t.Run("returns result. request not found", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records, indexes)
		req := b.Append(
			insolar.FirstPulseNumber+1,
			record.IncomingRequest{Nonce: rand.Uint64(), Reason: *insolar.NewReference(*insolar.NewID(insolar.FirstPulseNumber, nil))},
		)

		objectID := gen.ID()
		fromPulse := req.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer: &req.MetaID,
			},
		})
		require.NoError(t, err)

		_, err = calculator.ResultDuplicate(ctx, fromPulse, objectID, req.RecordID, record.Result{Request: gen.Reference()})
		require.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("returns no result. request found", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records, indexes)
		req := record.IncomingRequest{Nonce: rand.Uint64(), Reason: *insolar.NewReference(*insolar.NewID(insolar.FirstPulseNumber, nil))}
		req1 := b.Append(insolar.FirstPulseNumber+1, req)
		res := record.Result{Request: *insolar.NewReference(req1.RecordID)}
		resID := insolar.NewID(insolar.FirstPulseNumber+1, []byte{1})

		objectID := gen.ID()
		fromPulse := req1.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer: &req1.MetaID,
			},
		})
		require.NoError(t, err)

		fRes, err := calculator.ResultDuplicate(ctx, fromPulse, objectID, *resID, res)
		require.NoError(t, err)
		require.Nil(t, fRes)

		mc.Finish()
	})
}

func TestFilamentCalculatorDefault_RequestDuplicate(t *testing.T) {
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
	t.Run("returns error if reason is empty", func(t *testing.T) {
		_, _, err := calculator.RequestDuplicate(ctx, gen.PulseNumber(), gen.ID(), gen.ID(), &record.IncomingRequest{})
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("no records", func(t *testing.T) {
		objectID := gen.ID()
		fromPulse := gen.PulseNumber()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
		})
		require.NoError(t, err)

		req, res, err := calculator.RequestDuplicate(ctx, fromPulse, objectID, gen.ID(), &record.IncomingRequest{
			Reason: gen.Reference(),
		})

		assert.NoError(t, err)
		assert.Nil(t, req)
		assert.Nil(t, res)

		mc.Finish()
	})

	resetComponents()
	t.Run("returns request and result", func(t *testing.T) {
		coordinator.IsBeyondLimitMock.Return(true, nil)
		hNode := gen.Reference()
		coordinator.HeavyMock.Return(&hNode, nil)

		b := newFilamentBuilder(ctx, pcs, records, indexes)
		reason := *insolar.NewReference(*insolar.NewID(insolar.FirstPulseNumber, nil))
		req := record.IncomingRequest{Nonce: rand.Uint64(), Reason: reason}
		req1 := b.Append(insolar.FirstPulseNumber+1, req)
		res1 := b.Append(insolar.FirstPulseNumber+2, record.Result{Request: *insolar.NewReference(req1.RecordID)})

		sender.SendTargetFunc = func(_ context.Context, inMsg *message.Message, inNode insolar.Reference) (<-chan *message.Message, func()) {
			getReq := payload.GetRequest{}
			err := getReq.Unmarshal(inMsg.Payload)
			require.NoError(t, err)

			require.Equal(t, *reason.Record(), getReq.RequestID)
			require.Equal(t, hNode, inNode)

			reqMsg, err := payload.NewMessage(&payload.Request{})
			require.NoError(t, err)

			meta := payload.Meta{Payload: reqMsg.Payload}
			buf, err := meta.Marshal()
			require.NoError(t, err)
			reqMsg.Payload = buf
			ch := make(chan *message.Message, 1)
			ch <- reqMsg
			return ch, func() {}
		}

		objectID := gen.ID()
		fromPulse := res1.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer: &res1.MetaID,
			},
		})
		require.NoError(t, err)

		fReq, fRes, err := calculator.RequestDuplicate(ctx, fromPulse, objectID, req1.RecordID, &req)
		assert.NoError(t, err)
		require.Equal(t, *fReq, req1)
		assert.Equal(t, *fRes, res1)

		mc.Finish()
	})

	resetComponents()
	t.Run("returns only request", func(t *testing.T) {
		coordinator.IsBeyondLimitMock.Return(true, nil)
		hNode := gen.Reference()
		coordinator.HeavyMock.Return(&hNode, nil)

		b := newFilamentBuilder(ctx, pcs, records, indexes)
		reason := *insolar.NewReference(*insolar.NewID(insolar.FirstPulseNumber, nil))
		reqR := record.IncomingRequest{Nonce: rand.Uint64(), Reason: reason}
		req1 := b.Append(insolar.FirstPulseNumber+1, reqR)
		reqR2 := record.IncomingRequest{Nonce: rand.Uint64(), Reason: *insolar.NewReference(*insolar.NewID(insolar.FirstPulseNumber, nil))}
		req2 := b.Append(insolar.FirstPulseNumber+2, reqR2)

		sender.SendTargetFunc = func(_ context.Context, inMsg *message.Message, inNode insolar.Reference) (<-chan *message.Message, func()) {
			getReq := payload.GetRequest{}
			err := getReq.Unmarshal(inMsg.Payload)
			require.NoError(t, err)

			require.Equal(t, *reason.Record(), getReq.RequestID)
			require.Equal(t, hNode, inNode)

			reqMsg, err := payload.NewMessage(&payload.Request{})
			require.NoError(t, err)

			meta := payload.Meta{Payload: reqMsg.Payload}
			buf, err := meta.Marshal()
			require.NoError(t, err)
			reqMsg.Payload = buf
			ch := make(chan *message.Message, 1)
			ch <- reqMsg
			return ch, func() {}
		}

		objectID := gen.ID()
		fromPulse := req1.MetaID.Pulse()
		err := indexes.SetIndex(ctx, fromPulse, record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				PendingPointer: &req2.MetaID,
			},
		})
		require.NoError(t, err)

		fReq, fRes, err := calculator.RequestDuplicate(ctx, fromPulse, objectID, req1.RecordID, &reqR)
		require.NoError(t, err)
		require.Equal(t, *fReq, req1)
		require.Nil(t, fRes)

		mc.Finish()
	})

}
