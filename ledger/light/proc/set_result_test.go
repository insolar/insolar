// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc_test

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/testutils"
)

func TestSetResult_Proceed(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	defer mc.Finish()

	flowPulse := insolar.GenesisPulse.PulseNumber + 2
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		flowPulse,
	)

	writeAccessor := executor.NewWriteAccessorMock(mc)
	writeAccessor.BeginMock.Return(func() {}, nil)
	pcs := testutils.NewPlatformCryptographyScheme()

	sender := bus.NewSenderMock(mc)
	sender.ReplyMock.Return()
	detachedNotifier := executor.NewDetachedNotifierMock(mc)

	jetID := gen.JetID()
	objectID := gen.ID()
	requestID := gen.ID()

	resultRecord := &record.Result{
		Request: *insolar.NewReference(requestID),
		Object:  objectID,
	}
	virtual := record.Virtual{
		Union: &record.Virtual_Result{
			Result: resultRecord,
		},
	}
	hash := record.HashVirtual(pcs.ReferenceHasher(), virtual)
	resultID := *insolar.NewID(flow.Pulse(ctx), hash)
	virtualBuf, err := virtual.Marshal()
	require.NoError(t, err)

	result := payload.SetResult{
		Result: virtualBuf,
	}
	resultBuf, err := result.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: resultBuf,
	}
	LatestRequest := gen.IDWithPulse(flowPulse)
	expectedFilament := record.PendingFilament{
		RecordID:       resultID,
		PreviousRecord: &LatestRequest,
	}
	hash = record.HashVirtual(pcs.ReferenceHasher(), record.Wrap(&expectedFilament))
	expectedFilamentID := *insolar.NewID(resultID.Pulse(), hash)

	indexes := object.NewMemoryIndexStorageMock(mc)
	indexes.ForIDMock.Set(func(_ context.Context, pn insolar.PulseNumber, id insolar.ID) (record.Index, error) {
		require.Equal(t, flow.Pulse(ctx), pn)
		require.Equal(t, objectID, id)
		earliestPN := requestID.Pulse()
		return record.Index{
			Lifeline: record.Lifeline{
				LatestRequest:       &LatestRequest,
				EarliestOpenRequest: &earliestPN,
				OpenRequestsCount:   1,
			},
		}, nil
	})

	parent := gen.Reference()
	sideEffects := record.Activate{
		Request: gen.Reference(),
		Parent:  parent,
	}
	hash = record.HashVirtual(pcs.ReferenceHasher(), record.Wrap(&sideEffects))
	expectedSideEffectID := *insolar.NewID(resultID.Pulse(), hash)
	earliestID := gen.ID()
	earliestPulse := earliestID.Pulse()

	indexes.SetMock.Set(func(_ context.Context, pn insolar.PulseNumber, idx record.Index) {
		require.Equal(t, resultID.Pulse(), pn)
		expectedIndex := record.Index{
			LifelineLastUsed: resultID.Pulse(),
			Lifeline: record.Lifeline{
				LatestRequest:       &expectedFilamentID,
				LatestState:         &expectedSideEffectID,
				StateID:             record.StateActivation,
				Parent:              parent,
				EarliestOpenRequest: &earliestPulse,
			},
		}
		require.Equal(t, expectedIndex, idx)
	})

	records := object.NewAtomicRecordModifierMock(mc)
	records.SetAtomicMock.Set(func(_ context.Context, recs ...record.Material) (r error) {
		require.Equal(t, 3, len(recs))

		result := recs[0]
		filament := recs[1]
		sideEffect := recs[2]
		require.Equal(t, resultID, result.ID)
		require.Equal(t, resultRecord, record.Unwrap(&result.Virtual))

		require.Equal(t, expectedFilamentID, filament.ID)
		require.Equal(t, &expectedFilament, record.Unwrap(&filament.Virtual))

		require.Equal(t, &sideEffects, record.Unwrap(&sideEffect.Virtual))
		return nil
	})

	vImmut := record.Wrap(&record.IncomingRequest{Immutable: true})
	vMut := record.Wrap(&record.IncomingRequest{})
	opened := []record.CompositeFilamentRecord{
		{
			RecordID: earliestID,
			Record:   record.Material{Virtual: vImmut},
		},
		{
			RecordID: requestID,
			Record:   record.Material{Virtual: vMut},
		},
		{
			RecordID: gen.ID(),
			Record: record.Material{
				Virtual: record.Wrap(&record.OutgoingRequest{
					Reason:     *insolar.NewReference(requestID),
					ReturnMode: record.ReturnSaga,
				}),
			},
		},
	}

	filaments := executor.NewFilamentCalculatorMock(mc)
	filaments.ResultDuplicateMock.Set(func(_ context.Context, objID insolar.ID, resID insolar.ID, r record.Result) (*record.CompositeFilamentRecord, error) {
		require.Equal(t, objectID, objID)
		require.Equal(t, *resultRecord, r)
		return nil, nil
	})
	filaments.OpenedRequestsMock.Inspect(func(_ context.Context, pn insolar.PulseNumber, objID insolar.ID, pendingOnly bool) {
		require.Equal(t, objectID, objID)
		require.Equal(t, flow.Pulse(ctx), pn)
		require.False(t, pendingOnly)
	}).Return(opened, nil)

	detachedNotifier.NotifyMock.Inspect(func(ctx context.Context, openedRequests []record.CompositeFilamentRecord, objID insolar.ID, closedRequestID insolar.ID) {
		require.Equal(t, objectID, objID)
		require.Equal(t, requestID, closedRequestID)
		require.Equal(t, opened, openedRequests)
	}).Return()

	setResultProc := proc.NewSetResult(msg, jetID, *resultRecord, &sideEffects)
	setResultProc.Dep(writeAccessor, sender, object.NewIndexLocker(), filaments, records, indexes, pcs, detachedNotifier)

	err = setResultProc.Proceed(ctx)
	require.NoError(t, err)
}

func TestSetResult_Proceed_ResultDuplicated(t *testing.T) {
	t.Parallel()

	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)
	mc := minimock.NewController(t)
	defer mc.Finish()

	detachedNotifier := executor.NewDetachedNotifierMock(mc)
	writeAccessor := executor.NewWriteAccessorMock(mc)
	records := object.NewAtomicRecordModifierMock(mc)
	indexes := object.NewMemoryIndexStorageMock(mc)
	indexes.ForIDMock.Return(record.Index{}, nil)
	pcs := testutils.NewPlatformCryptographyScheme()

	sender := bus.NewSenderMock(mc)

	jetID := gen.JetID()
	objectID := gen.ID()
	resultID := gen.ID()

	res := &record.Result{
		Object: objectID,
	}
	virtual := record.Virtual{
		Union: &record.Virtual_Result{
			Result: res,
		},
	}
	m := record.Material{Virtual: virtual}
	duplicateBuf, err := m.Marshal()
	require.NoError(t, err)

	result := payload.SetResult{
		Result: duplicateBuf,
	}
	resultBuf, err := result.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: resultBuf,
	}

	filaments := executor.NewFilamentCalculatorMock(mc)
	filaments.ResultDuplicateMock.Set(func(_ context.Context, objID insolar.ID, resID insolar.ID, r record.Result) (*record.CompositeFilamentRecord, error) {
		require.Equal(t, objectID, objID)
		require.Equal(t, *res, r)

		return &record.CompositeFilamentRecord{
			Record:   record.Material{Virtual: virtual},
			RecordID: resultID,
		}, nil
	})
	sender.ReplyMock.Set(func(_ context.Context, receivedMeta payload.Meta, resMsg *message.Message) {
		require.Equal(t, msg, receivedMeta)

		resp, err := payload.Unmarshal(resMsg.Payload)
		require.NoError(t, err)

		res, ok := resp.(*payload.ErrorResultExists)
		require.True(t, ok)
		require.Equal(t, resultID, res.ResultID)
		receivedResult := record.Material{}
		err = receivedResult.Unmarshal(res.Result)
		require.NoError(t, err)
		require.Equal(t, virtual, receivedResult.Virtual)
	})

	setResultProc := proc.NewSetResult(msg, jetID, *res, nil)
	setResultProc.Dep(writeAccessor, sender, object.NewIndexLocker(), filaments, records, indexes, pcs, detachedNotifier)
	err = setResultProc.Proceed(ctx)
	require.NoError(t, err)
}

func TestSetResult_Proceed_ImmutableRequest_Error(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	defer mc.Finish()

	flowPulse := insolar.GenesisPulse.PulseNumber + 2
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		flowPulse,
	)

	writeAccessor := executor.NewWriteAccessorMock(mc)
	pcs := testutils.NewPlatformCryptographyScheme()

	sender := bus.NewSenderMock(mc)
	detachedNotifier := executor.NewDetachedNotifierMock(mc)

	jetID := gen.JetID()
	objectID := gen.ID()
	requestID := gen.ID()

	resultRecord := &record.Result{
		Request: *insolar.NewReference(requestID),
		Object:  objectID,
	}
	virtual := record.Virtual{
		Union: &record.Virtual_Result{
			Result: resultRecord,
		},
	}
	hash := record.HashVirtual(pcs.ReferenceHasher(), virtual)
	resultID := *insolar.NewID(flow.Pulse(ctx), hash)
	virtualBuf, err := virtual.Marshal()
	require.NoError(t, err)

	result := payload.SetResult{
		Result: virtualBuf,
	}
	resultBuf, err := result.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: resultBuf,
	}
	LatestRequest := gen.IDWithPulse(flowPulse)
	expectedFilament := record.PendingFilament{
		RecordID:       resultID,
		PreviousRecord: &LatestRequest,
	}
	hash = record.HashVirtual(pcs.ReferenceHasher(), record.Wrap(&expectedFilament))

	indexes := object.NewMemoryIndexStorageMock(mc)
	indexes.ForIDMock.Set(func(_ context.Context, pn insolar.PulseNumber, id insolar.ID) (record.Index, error) {
		require.Equal(t, flow.Pulse(ctx), pn)
		require.Equal(t, objectID, id)
		earliestPN := requestID.Pulse()
		return record.Index{
			Lifeline: record.Lifeline{
				LatestRequest:       &LatestRequest,
				EarliestOpenRequest: &earliestPN,
				OpenRequestsCount:   1,
			},
		}, nil
	})

	parent := gen.Reference()
	sideEffects := record.Activate{
		Request: gen.Reference(),
		Parent:  parent,
	}
	hash = record.HashVirtual(pcs.ReferenceHasher(), record.Wrap(&sideEffects))

	records := object.NewAtomicRecordModifierMock(mc)

	filaments := executor.NewFilamentCalculatorMock(mc)
	filaments.ResultDuplicateMock.Set(func(_ context.Context, objID insolar.ID, resID insolar.ID, r record.Result) (*record.CompositeFilamentRecord, error) {
		require.Equal(t, objectID, objID)
		require.Equal(t, *resultRecord, r)
		return nil, nil
	})
	filaments.OpenedRequestsMock.Set(func(_ context.Context, pn insolar.PulseNumber, objID insolar.ID, pendingOnly bool) ([]record.CompositeFilamentRecord, error) {
		require.Equal(t, objectID, objID)
		require.Equal(t, flow.Pulse(ctx), pn)
		require.False(t, pendingOnly)

		opened := []record.CompositeFilamentRecord{
			// req that we closing
			{
				RecordID: requestID,
				Record: record.Material{
					Virtual: record.Wrap(&record.IncomingRequest{
						Immutable: true,
					}),
				},
			},
		}
		return opened, nil
	})

	setResultProc := proc.NewSetResult(msg, jetID, *resultRecord, &sideEffects)
	setResultProc.Dep(writeAccessor, sender, object.NewIndexLocker(), filaments, records, indexes, pcs, detachedNotifier)

	err = setResultProc.Proceed(ctx)
	require.Error(t, err)
	insError, ok := errors.Cause(err).(*payload.CodedError)
	require.True(t, ok)
	require.Equal(t, payload.CodeRequestInvalid, insError.GetCode())
}

func TestSetResult_Proceed_OutgoingRequest_Error(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	defer mc.Finish()

	flowPulse := insolar.GenesisPulse.PulseNumber + 2
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		flowPulse,
	)

	writeAccessor := executor.NewWriteAccessorMock(mc)
	pcs := testutils.NewPlatformCryptographyScheme()

	sender := bus.NewSenderMock(mc)
	detachedNotifier := executor.NewDetachedNotifierMock(mc)

	jetID := gen.JetID()
	objectID := gen.ID()
	requestID := gen.ID()

	resultRecord := &record.Result{
		Request: *insolar.NewReference(requestID),
		Object:  objectID,
	}
	virtual := record.Virtual{
		Union: &record.Virtual_Result{
			Result: resultRecord,
		},
	}
	hash := record.HashVirtual(pcs.ReferenceHasher(), virtual)
	resultID := *insolar.NewID(flow.Pulse(ctx), hash)
	virtualBuf, err := virtual.Marshal()
	require.NoError(t, err)

	result := payload.SetResult{
		Result: virtualBuf,
	}
	resultBuf, err := result.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: resultBuf,
	}
	LatestRequest := gen.IDWithPulse(flowPulse)
	expectedFilament := record.PendingFilament{
		RecordID:       resultID,
		PreviousRecord: &LatestRequest,
	}
	hash = record.HashVirtual(pcs.ReferenceHasher(), record.Wrap(&expectedFilament))

	indexes := object.NewMemoryIndexStorageMock(mc)
	indexes.ForIDMock.Set(func(_ context.Context, pn insolar.PulseNumber, id insolar.ID) (record.Index, error) {
		require.Equal(t, flow.Pulse(ctx), pn)
		require.Equal(t, objectID, id)
		earliestPN := requestID.Pulse()
		return record.Index{
			Lifeline: record.Lifeline{
				LatestRequest:       &LatestRequest,
				EarliestOpenRequest: &earliestPN,
				OpenRequestsCount:   1,
			},
		}, nil
	})

	parent := gen.Reference()
	sideEffects := record.Activate{
		Request: gen.Reference(),
		Parent:  parent,
	}
	hash = record.HashVirtual(pcs.ReferenceHasher(), record.Wrap(&sideEffects))

	records := object.NewAtomicRecordModifierMock(mc)

	filaments := executor.NewFilamentCalculatorMock(mc)
	filaments.ResultDuplicateMock.Set(func(_ context.Context, objID insolar.ID, resID insolar.ID, r record.Result) (*record.CompositeFilamentRecord, error) {
		require.Equal(t, objectID, objID)
		require.Equal(t, *resultRecord, r)
		return nil, nil
	})
	filaments.OpenedRequestsMock.Set(func(_ context.Context, pn insolar.PulseNumber, objID insolar.ID, pendingOnly bool) ([]record.CompositeFilamentRecord, error) {
		require.Equal(t, objectID, objID)
		require.Equal(t, flow.Pulse(ctx), pn)
		require.False(t, pendingOnly)

		opened := []record.CompositeFilamentRecord{
			// req that we closing
			{
				RecordID: requestID,
				Record: record.Material{
					Virtual: record.Wrap(&record.OutgoingRequest{}),
				},
			},
		}
		return opened, nil
	})

	setResultProc := proc.NewSetResult(msg, jetID, *resultRecord, &sideEffects)
	setResultProc.Dep(writeAccessor, sender, object.NewIndexLocker(), filaments, records, indexes, pcs, detachedNotifier)

	err = setResultProc.Proceed(ctx)
	require.Error(t, err)
	insError, ok := errors.Cause(err).(*payload.CodedError)
	require.True(t, ok)
	require.Equal(t, payload.CodeRequestInvalid, insError.GetCode())
}

func TestSetResult_Proceed_NotFoundInOpened_Error(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	defer mc.Finish()

	flowPulse := insolar.GenesisPulse.PulseNumber + 2
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		flowPulse,
	)

	writeAccessor := executor.NewWriteAccessorMock(mc)
	pcs := testutils.NewPlatformCryptographyScheme()

	sender := bus.NewSenderMock(mc)
	detachedNotifier := executor.NewDetachedNotifierMock(mc)

	jetID := gen.JetID()
	objectID := gen.ID()
	requestID := gen.ID()

	resultRecord := &record.Result{
		Request: *insolar.NewReference(requestID),
		Object:  objectID,
	}
	virtual := record.Virtual{
		Union: &record.Virtual_Result{
			Result: resultRecord,
		},
	}
	hash := record.HashVirtual(pcs.ReferenceHasher(), virtual)
	resultID := *insolar.NewID(flow.Pulse(ctx), hash)
	virtualBuf, err := virtual.Marshal()
	require.NoError(t, err)

	result := payload.SetResult{
		Result: virtualBuf,
	}
	resultBuf, err := result.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: resultBuf,
	}
	LatestRequest := gen.IDWithPulse(flowPulse)
	expectedFilament := record.PendingFilament{
		RecordID:       resultID,
		PreviousRecord: &LatestRequest,
	}
	hash = record.HashVirtual(pcs.ReferenceHasher(), record.Wrap(&expectedFilament))

	indexes := object.NewMemoryIndexStorageMock(mc)
	indexes.ForIDMock.Set(func(_ context.Context, pn insolar.PulseNumber, id insolar.ID) (record.Index, error) {
		require.Equal(t, flow.Pulse(ctx), pn)
		require.Equal(t, objectID, id)
		earliestPN := requestID.Pulse()
		return record.Index{
			Lifeline: record.Lifeline{
				LatestRequest:       &LatestRequest,
				EarliestOpenRequest: &earliestPN,
				OpenRequestsCount:   1,
			},
		}, nil
	})

	parent := gen.Reference()
	sideEffects := record.Activate{
		Request: gen.Reference(),
		Parent:  parent,
	}
	hash = record.HashVirtual(pcs.ReferenceHasher(), record.Wrap(&sideEffects))

	records := object.NewAtomicRecordModifierMock(mc)

	filaments := executor.NewFilamentCalculatorMock(mc)
	filaments.ResultDuplicateMock.Set(func(_ context.Context, objID insolar.ID, resID insolar.ID, r record.Result) (*record.CompositeFilamentRecord, error) {
		require.Equal(t, objectID, objID)
		require.Equal(t, *resultRecord, r)
		return nil, nil
	})
	filaments.OpenedRequestsMock.Set(func(_ context.Context, pn insolar.PulseNumber, objID insolar.ID, pendingOnly bool) ([]record.CompositeFilamentRecord, error) {
		require.Equal(t, objectID, objID)
		require.Equal(t, flow.Pulse(ctx), pn)
		require.False(t, pendingOnly)

		opened := []record.CompositeFilamentRecord{
			{
				RecordID: gen.ID(),
				Record: record.Material{
					Virtual: record.Wrap(&record.IncomingRequest{
						Immutable: true,
					}),
				},
			},
			{
				RecordID: gen.ID(),
				Record: record.Material{
					Virtual: record.Wrap(&record.IncomingRequest{
						Immutable: true,
					}),
				},
			},
		}
		return opened, nil
	})

	setResultProc := proc.NewSetResult(msg, jetID, *resultRecord, &sideEffects)
	setResultProc.Dep(writeAccessor, sender, object.NewIndexLocker(), filaments, records, indexes, pcs, detachedNotifier)

	err = setResultProc.Proceed(ctx)
	require.Error(t, err)
	insError, ok := errors.Cause(err).(*payload.CodedError)
	require.True(t, ok)
	require.Equal(t, payload.CodeRequestNotFound, insError.GetCode())
}

func TestSetResult_Proceed_HasOpenedOutgoing_Error(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	defer mc.Finish()

	flowPulse := insolar.GenesisPulse.PulseNumber + 2
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		flowPulse,
	)

	writeAccessor := executor.NewWriteAccessorMock(mc)
	pcs := testutils.NewPlatformCryptographyScheme()

	sender := bus.NewSenderMock(mc)
	detachedNotifier := executor.NewDetachedNotifierMock(mc)

	jetID := gen.JetID()
	objectID := gen.ID()
	requestID := gen.ID()

	resultRecord := &record.Result{
		Request: *insolar.NewReference(requestID),
		Object:  objectID,
	}
	virtual := record.Virtual{
		Union: &record.Virtual_Result{
			Result: resultRecord,
		},
	}
	hash := record.HashVirtual(pcs.ReferenceHasher(), virtual)
	resultID := *insolar.NewID(flow.Pulse(ctx), hash)
	virtualBuf, err := virtual.Marshal()
	require.NoError(t, err)

	result := payload.SetResult{
		Result: virtualBuf,
	}
	resultBuf, err := result.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: resultBuf,
	}
	LatestRequest := gen.IDWithPulse(flowPulse)
	expectedFilament := record.PendingFilament{
		RecordID:       resultID,
		PreviousRecord: &LatestRequest,
	}
	hash = record.HashVirtual(pcs.ReferenceHasher(), record.Wrap(&expectedFilament))

	indexes := object.NewMemoryIndexStorageMock(mc)
	indexes.ForIDMock.Set(func(_ context.Context, pn insolar.PulseNumber, id insolar.ID) (record.Index, error) {
		require.Equal(t, flow.Pulse(ctx), pn)
		require.Equal(t, objectID, id)
		earliestPN := requestID.Pulse()
		return record.Index{
			Lifeline: record.Lifeline{
				LatestRequest:       &LatestRequest,
				EarliestOpenRequest: &earliestPN,
				OpenRequestsCount:   1,
			},
		}, nil
	})

	parent := gen.Reference()
	sideEffects := record.Activate{
		Request: gen.Reference(),
		Parent:  parent,
	}
	hash = record.HashVirtual(pcs.ReferenceHasher(), record.Wrap(&sideEffects))
	records := object.NewAtomicRecordModifierMock(mc)

	filaments := executor.NewFilamentCalculatorMock(mc)
	filaments.ResultDuplicateMock.Set(func(_ context.Context, objID insolar.ID, resID insolar.ID, r record.Result) (*record.CompositeFilamentRecord, error) {
		require.Equal(t, objectID, objID)
		require.Equal(t, *resultRecord, r)
		return nil, nil
	})
	filaments.OpenedRequestsMock.Set(func(_ context.Context, pn insolar.PulseNumber, objID insolar.ID, pendingOnly bool) ([]record.CompositeFilamentRecord, error) {
		require.Equal(t, objectID, objID)
		require.Equal(t, flow.Pulse(ctx), pn)
		require.False(t, pendingOnly)

		opened := []record.CompositeFilamentRecord{
			// req that we closing
			{
				RecordID: requestID,
				Record: record.Material{
					Virtual: record.Wrap(&record.IncomingRequest{}),
				},
			},
			// outgoing where reason is closing incoming
			{
				RecordID: gen.ID(),
				Record: record.Material{
					Virtual: record.Wrap(&record.OutgoingRequest{
						Reason: *insolar.NewReference(requestID),
					}),
				},
			},
		}
		return opened, nil
	})

	setResultProc := proc.NewSetResult(msg, jetID, *resultRecord, &sideEffects)
	setResultProc.Dep(writeAccessor, sender, object.NewIndexLocker(), filaments, records, indexes, pcs, detachedNotifier)

	err = setResultProc.Proceed(ctx)
	require.Error(t, err)
	insError, ok := errors.Cause(err).(*payload.CodedError)
	require.True(t, ok)
	require.Equal(t, payload.CodeRequestNonClosedOutgoing, insError.GetCode())
}

func TestSetResult_Proceed_OldestMutableRequest(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	defer mc.Finish()

	flowPulse := insolar.GenesisPulse.PulseNumber + 2
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		flowPulse,
	)

	writeAccessor := executor.NewWriteAccessorMock(mc)
	pcs := testutils.NewPlatformCryptographyScheme()

	sender := bus.NewSenderMock(mc)
	detachedNotifier := executor.NewDetachedNotifierMock(mc)

	jetID := gen.JetID()
	objectID := gen.ID()
	requestID := gen.ID()

	resultRecord := &record.Result{
		Request: *insolar.NewReference(requestID),
		Object:  objectID,
	}
	virtual := record.Virtual{
		Union: &record.Virtual_Result{
			Result: resultRecord,
		},
	}
	hash := record.HashVirtual(pcs.ReferenceHasher(), virtual)
	resultID := *insolar.NewID(flow.Pulse(ctx), hash)
	virtualBuf, err := virtual.Marshal()
	require.NoError(t, err)

	result := payload.SetResult{
		Result: virtualBuf,
	}
	resultBuf, err := result.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: resultBuf,
	}
	LatestRequest := gen.IDWithPulse(flowPulse)
	expectedFilament := record.PendingFilament{
		RecordID:       resultID,
		PreviousRecord: &LatestRequest,
	}
	hash = record.HashVirtual(pcs.ReferenceHasher(), record.Wrap(&expectedFilament))

	indexes := object.NewMemoryIndexStorageMock(mc)
	indexes.ForIDMock.Set(func(_ context.Context, pn insolar.PulseNumber, id insolar.ID) (record.Index, error) {
		require.Equal(t, flow.Pulse(ctx), pn)
		require.Equal(t, objectID, id)
		earliestPN := requestID.Pulse()
		return record.Index{
			Lifeline: record.Lifeline{
				LatestRequest:       &LatestRequest,
				EarliestOpenRequest: &earliestPN,
				OpenRequestsCount:   1,
			},
		}, nil
	})

	records := object.NewAtomicRecordModifierMock(mc)

	filaments := executor.NewFilamentCalculatorMock(mc)
	filaments.ResultDuplicateMock.Set(func(_ context.Context, objID insolar.ID, resID insolar.ID, r record.Result) (*record.CompositeFilamentRecord, error) {
		require.Equal(t, objectID, objID)
		require.Equal(t, *resultRecord, r)
		return nil, nil
	})
	filaments.OpenedRequestsMock.Set(func(_ context.Context, pn insolar.PulseNumber, objID insolar.ID, pendingOnly bool) ([]record.CompositeFilamentRecord, error) {
		require.Equal(t, objectID, objID)
		require.Equal(t, flow.Pulse(ctx), pn)
		require.False(t, pendingOnly)

		opened := []record.CompositeFilamentRecord{
			// Here we have unclosed Mutable request (Immutable == false below)
			// and other IDs for closing attempt. We should get an error from check this record.
			{
				RecordID: gen.ID(),
				Record: record.Material{
					Virtual: record.Wrap(&record.IncomingRequest{
						Immutable: false,
					}),
				},
			},
			// We shouldn't process this record,
			// because we have unclosed oldest mutable request (see above).
			{
				RecordID: requestID,
				Record: record.Material{
					Virtual: record.Wrap(&record.IncomingRequest{
						Immutable: false,
					}),
				},
			},
		}
		return opened, nil
	})

	setResultProc := proc.NewSetResult(msg, jetID, *resultRecord, nil)
	setResultProc.Dep(writeAccessor, sender, object.NewIndexLocker(), filaments, records, indexes, pcs, detachedNotifier)

	err = setResultProc.Proceed(ctx)
	require.Error(t, err)
	insError, ok := errors.Cause(err).(*payload.CodedError)
	require.True(t, ok)
	require.Equal(t, payload.CodeRequestNonOldestMutable, insError.GetCode())
}

func TestSetResult_Proceed_CloseImmutOutWhenOpenedMutableInFilaments(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	defer mc.Finish()

	flowPulse := insolar.GenesisPulse.PulseNumber + 2
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		flowPulse,
	)

	writeAccessor := executor.NewWriteAccessorMock(mc)
	writeAccessor.BeginMock.Return(func() {}, nil)
	pcs := testutils.NewPlatformCryptographyScheme()

	sender := bus.NewSenderMock(mc)
	sender.ReplyMock.Return()
	detachedNotifier := executor.NewDetachedNotifierMock(mc)

	jetID := gen.JetID()
	objectID := gen.ID()
	requestID := gen.ID()

	resultRecord := &record.Result{
		Request: *insolar.NewReference(requestID),
		Object:  objectID,
	}
	virtual := record.Virtual{
		Union: &record.Virtual_Result{
			Result: resultRecord,
		},
	}
	hash := record.HashVirtual(pcs.ReferenceHasher(), virtual)
	resultID := *insolar.NewID(flow.Pulse(ctx), hash)
	virtualBuf, err := virtual.Marshal()
	require.NoError(t, err)

	result := payload.SetResult{
		Result: virtualBuf,
	}
	resultBuf, err := result.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: resultBuf,
	}
	LatestRequest := gen.IDWithPulse(flowPulse)
	expectedFilament := record.PendingFilament{
		RecordID:       resultID,
		PreviousRecord: &LatestRequest,
	}
	hash = record.HashVirtual(pcs.ReferenceHasher(), record.Wrap(&expectedFilament))
	expectedFilamentID := *insolar.NewID(resultID.Pulse(), hash)

	indexes := object.NewMemoryIndexStorageMock(mc)
	indexes.ForIDMock.Set(func(_ context.Context, pn insolar.PulseNumber, id insolar.ID) (record.Index, error) {
		require.Equal(t, flow.Pulse(ctx), pn)
		require.Equal(t, objectID, id)
		earliestPN := requestID.Pulse()
		return record.Index{
			Lifeline: record.Lifeline{
				LatestRequest:       &LatestRequest,
				EarliestOpenRequest: &earliestPN,
				OpenRequestsCount:   1,
			},
		}, nil
	})

	earliestID := gen.ID()
	earliestPulse := earliestID.Pulse()

	indexes.SetMock.Set(func(_ context.Context, pn insolar.PulseNumber, idx record.Index) {
		require.Equal(t, resultID.Pulse(), pn)
		expectedIndex := record.Index{
			LifelineLastUsed: resultID.Pulse(),
			Lifeline: record.Lifeline{
				LatestRequest:       &expectedFilamentID,
				LatestState:         nil,
				StateID:             record.StateUndefined,
				EarliestOpenRequest: &earliestPulse,
			},
		}
		require.Equal(t, expectedIndex, idx)
	})

	records := object.NewAtomicRecordModifierMock(mc)
	records.SetAtomicMock.Set(func(_ context.Context, recs ...record.Material) (r error) {
		require.Equal(t, 2, len(recs))

		result := recs[0]
		filament := recs[1]
		require.Equal(t, resultID, result.ID)
		require.Equal(t, resultRecord, record.Unwrap(&result.Virtual))

		require.Equal(t, expectedFilamentID, filament.ID)
		require.Equal(t, &expectedFilament, record.Unwrap(&filament.Virtual))

		return nil
	})

	filaments := executor.NewFilamentCalculatorMock(mc)
	filaments.ResultDuplicateMock.Set(func(_ context.Context, objID insolar.ID, resID insolar.ID, r record.Result) (*record.CompositeFilamentRecord, error) {
		require.Equal(t, objectID, objID)
		require.Equal(t, *resultRecord, r)
		return nil, nil
	})
	filaments.OpenedRequestsMock.Set(func(_ context.Context, pn insolar.PulseNumber, objID insolar.ID, pendingOnly bool) ([]record.CompositeFilamentRecord, error) {
		require.Equal(t, objectID, objID)
		require.Equal(t, flow.Pulse(ctx), pn)
		require.False(t, pendingOnly)

		vImmut := record.Wrap(&record.IncomingRequest{Immutable: true})
		vMut := record.Wrap(&record.IncomingRequest{})
		opened := []record.CompositeFilamentRecord{
			{
				RecordID: earliestID,
				Record:   record.Material{Virtual: vImmut},
			},
			{
				RecordID: gen.ID(),
				Record:   record.Material{Virtual: vMut},
			},
			{
				RecordID: requestID,
				Record: record.Material{
					Virtual: record.Wrap(&record.OutgoingRequest{
						Reason:     *insolar.NewReference(requestID),
						ReturnMode: record.ReturnSaga,
						Immutable:  false,
					}),
				},
			},
		}
		return opened, nil
	})

	setResultProc := proc.NewSetResult(msg, jetID, *resultRecord, nil)
	setResultProc.Dep(writeAccessor, sender, object.NewIndexLocker(), filaments, records, indexes, pcs, detachedNotifier)

	err = setResultProc.Proceed(ctx)
	require.NoError(t, err)
}
