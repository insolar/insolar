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
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/pulse"
	"github.com/insolar/insolar/testutils"
)

func TestSetRequest_Proceed(t *testing.T) {
	t.Parallel()
	flowPN := insolar.GenesisPulse.PulseNumber + 10

	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		flowPN,
	)
	mc := minimock.NewController(t)
	pcs := testutils.NewPlatformCryptographyScheme()

	var (
		writeAccessor *executor.WriteAccessorMock
		sender        *bus.SenderMock
		filaments     *executor.FilamentCalculatorMock
		idxStorage    *object.MemoryIndexStorageMock
		records       *object.AtomicRecordModifierMock
		checker       *executor.RequestCheckerMock
		coordinator   *jet.CoordinatorMock
	)

	resetComponents := func() {
		writeAccessor = executor.NewWriteAccessorMock(mc)
		sender = bus.NewSenderMock(mc)
		filaments = executor.NewFilamentCalculatorMock(mc)
		idxStorage = object.NewMemoryIndexStorageMock(mc)
		records = object.NewAtomicRecordModifierMock(mc)
		checker = executor.NewRequestCheckerMock(mc)
		coordinator = jet.NewCoordinatorMock(t)
	}

	ref := gen.Reference()
	jetID := gen.JetID()
	requestID := gen.ID()

	request := record.IncomingRequest{
		Object:   &ref,
		CallType: record.CTMethod,
	}
	virtual := record.Virtual{
		Union: &record.Virtual_IncomingRequest{
			IncomingRequest: &request,
		},
	}

	pl := payload.SetIncomingRequest{
		Request: virtual,
	}
	requestBuf, err := pl.Marshal()
	require.NoError(t, err)

	virtualRef := gen.Reference()
	msg := payload.Meta{
		Payload: requestBuf,
		Sender:  virtualRef,
	}

	resetComponents()
	t.Run("happy basic", func(t *testing.T) {
		idxStorage.ForIDMock.Return(record.Index{
			Lifeline: record.Lifeline{
				StateID: record.StateActivation,
			},
		}, nil)
		idxStorage.SetMock.Set(func(_ context.Context, pn insolar.PulseNumber, idx record.Index) {
			require.Equal(t, requestID.Pulse(), pn)

			virtual = record.Wrap(&record.PendingFilament{
				RecordID:       requestID,
				PreviousRecord: nil,
			})
			hash := record.HashVirtual(pcs.ReferenceHasher(), virtual)
			pendingID := insolar.NewID(requestID.Pulse(), hash)
			expectedIndex := record.Index{
				LifelineLastUsed: pn,
				Lifeline: record.Lifeline{
					StateID:             record.StateActivation,
					EarliestOpenRequest: &pn,
					LatestRequest:       pendingID,
					OpenRequestsCount:   1,
				},
			}
			require.Equal(t, expectedIndex, idx)
		})

		writeAccessor.BeginMock.Return(func() {}, nil)
		filaments.RequestDuplicateMock.Return(nil, nil, nil)
		sender.ReplyMock.Return()
		coordinator.VirtualExecutorForObjectMock.Set(func(_ context.Context, objID insolar.ID, pn insolar.PulseNumber) (r *insolar.Reference, r1 error) {
			require.Equal(t, flowPN, pn)
			require.Equal(t, *ref.GetLocal(), objID)

			return &virtualRef, nil
		})
		records.SetAtomicMock.Set(func(_ context.Context, recs ...record.Material) (r error) {
			require.Equal(t, len(recs), 2)
			req := recs[0]
			filament := recs[1]

			require.Equal(t, requestID, req.ID)
			require.Equal(t, record.Unwrap(&req.Virtual), &request)
			hash := record.HashVirtual(pcs.ReferenceHasher(), filament.Virtual)
			calcID := *insolar.NewID(requestID.Pulse(), hash)
			require.Equal(t, calcID, filament.ID)
			return nil
		})
		checker.ValidateRequestMock.Set(func(_ context.Context, id insolar.ID, req record.Request) (r error) {
			require.Equal(t, requestID, id)
			require.Equal(t, &request, req)
			return nil
		})
		checker.CheckRequestMock.Set(func(_ context.Context, id insolar.ID, req record.Request) (r error) {
			require.Equal(t, requestID, id)
			require.Equal(t, &request, req)
			return nil
		})

		p := proc.NewSetRequest(msg, &request, requestID, jetID)
		p.Dep(writeAccessor, filaments, sender, object.NewIndexLocker(), idxStorage, records, pcs, checker, coordinator)

		err = p.Proceed(ctx)
		require.NoError(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("duplicate returns correct requestID", func(t *testing.T) {
		reqID := gen.ID()
		resID := gen.ID()
		idxStorage.ForIDMock.Return(record.Index{
			Lifeline: record.Lifeline{
				StateID: record.StateActivation,
			},
		}, nil)
		checker.ValidateRequestMock.Set(func(_ context.Context, id insolar.ID, req record.Request) (r error) {
			require.Equal(t, requestID, id)
			require.Equal(t, &request, req)
			return nil
		})
		filaments.RequestDuplicateMock.Return(
			&record.CompositeFilamentRecord{RecordID: reqID},
			&record.CompositeFilamentRecord{RecordID: resID},
			nil,
		)

		sender.ReplyMock.Set(func(_ context.Context, meta payload.Meta, msg *message.Message) {
			pl, err := payload.Unmarshal(msg.Payload)
			require.NoError(t, err)
			rep, ok := pl.(*payload.RequestInfo)
			require.True(t, ok)
			require.Equal(t, reqID, rep.RequestID)
		})
		coordinator.VirtualExecutorForObjectMock.Set(func(_ context.Context, objID insolar.ID, pn insolar.PulseNumber) (r *insolar.Reference, r1 error) {
			require.Equal(t, flowPN, pn)
			require.Equal(t, *ref.GetLocal(), objID)

			return &virtualRef, nil
		})

		p := proc.NewSetRequest(msg, &request, requestID, jetID)
		p.Dep(writeAccessor, filaments, sender, object.NewIndexLocker(), idxStorage, records, pcs, checker, coordinator)

		err = p.Proceed(ctx)
		require.NoError(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("wrong sender", func(t *testing.T) {
		coordinator.VirtualExecutorForObjectMock.Set(func(_ context.Context, objID insolar.ID, pn insolar.PulseNumber) (r *insolar.Reference, r1 error) {
			require.Equal(t, flowPN, pn)
			require.Equal(t, *ref.GetLocal(), objID)

			virtualRef := gen.Reference()
			return &virtualRef, nil
		})

		p := proc.NewSetRequest(msg, &request, requestID, jetID)
		p.Dep(writeAccessor, filaments, sender, object.NewIndexLocker(), idxStorage, records, pcs, checker, coordinator)

		err = p.Proceed(ctx)
		require.Error(t, err)

		mc.Finish()
	})

	t.Run("object is not activated error", func(t *testing.T) {
		resetComponents()
		defer mc.Finish()

		request := record.IncomingRequest{
			Object:   &ref,
			CallType: record.CTMethod,
		}
		coordinator.VirtualExecutorForObjectMock.Set(func(_ context.Context, objID insolar.ID, pn insolar.PulseNumber) (r *insolar.Reference, r1 error) {
			require.Equal(t, flowPN, pn)
			require.Equal(t, *ref.GetLocal(), objID)

			return &virtualRef, nil
		})
		idxStorage.ForIDMock.Return(record.Index{
			Lifeline: record.Lifeline{
				StateID: record.StateUndefined,
			},
		}, nil)

		p := proc.NewSetRequest(msg, &request, requestID, jetID)
		p.Dep(writeAccessor, filaments, sender, object.NewIndexLocker(), idxStorage, records, pcs, checker, coordinator)

		err = p.Proceed(ctx)
		require.Error(t, err)
		insError, ok := errors.Cause(err).(*payload.CodedError)
		require.True(t, ok)
		require.Equal(t, payload.CodeNonActivated, insError.GetCode())
	})

	t.Run("object is deactivated error", func(t *testing.T) {
		resetComponents()
		defer mc.Finish()

		request := record.IncomingRequest{
			Object:   &ref,
			CallType: record.CTMethod,
		}
		coordinator.VirtualExecutorForObjectMock.Set(func(_ context.Context, objID insolar.ID, pn insolar.PulseNumber) (r *insolar.Reference, r1 error) {
			require.Equal(t, flowPN, pn)
			require.Equal(t, *ref.GetLocal(), objID)

			return &virtualRef, nil
		})
		idxStorage.ForIDMock.Return(record.Index{
			Lifeline: record.Lifeline{
				StateID: record.StateDeactivation,
			},
		}, nil)

		p := proc.NewSetRequest(msg, &request, requestID, jetID)
		p.Dep(writeAccessor, filaments, sender, object.NewIndexLocker(), idxStorage, records, pcs, checker, coordinator)

		err = p.Proceed(ctx)
		require.Error(t, err)
		insError, ok := errors.Cause(err).(*payload.CodedError)
		require.True(t, ok)
		require.Equal(t, payload.CodeDeactivated, insError.GetCode())
	})

	t.Run("request from past error", func(t *testing.T) {
		resetComponents()
		defer mc.Finish()

		last := gen.IDWithPulse(pulse.MinTimePulse + 100)
		requestID := gen.IDWithPulse(pulse.MinTimePulse + 1)

		request := record.IncomingRequest{
			Object:   &ref,
			CallType: record.CTMethod,
		}
		coordinator.VirtualExecutorForObjectMock.Set(func(_ context.Context, objID insolar.ID, pn insolar.PulseNumber) (r *insolar.Reference, r1 error) {
			require.Equal(t, flowPN, pn)
			require.Equal(t, *ref.GetLocal(), objID)

			return &virtualRef, nil
		})
		idxStorage.ForIDMock.Return(record.Index{
			Lifeline: record.Lifeline{
				StateID:       record.StateActivation,
				LatestRequest: &last,
			},
		}, nil)

		p := proc.NewSetRequest(msg, &request, requestID, jetID)
		p.Dep(writeAccessor, filaments, sender, object.NewIndexLocker(), idxStorage, records, pcs, checker, coordinator)

		err = p.Proceed(ctx)
		require.Error(t, err)
	})
}
