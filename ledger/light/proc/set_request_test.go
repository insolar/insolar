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

package proc_test

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
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
		writeAccessor *hot.WriteAccessorMock
		sender        *bus.SenderMock
		filaments     *executor.FilamentCalculatorMock
		idxStorage    *object.IndexStorageMock
		records       *object.RecordModifierMock
		checker       *executor.RequestCheckerMock
		coordinator   *jet.CoordinatorMock
	)

	resetComponents := func() {
		writeAccessor = hot.NewWriteAccessorMock(mc)
		sender = bus.NewSenderMock(mc)
		filaments = executor.NewFilamentCalculatorMock(mc)
		idxStorage = object.NewIndexStorageMock(mc)
		records = object.NewRecordModifierMock(mc)
		checker = executor.NewRequestCheckerMock(mc)
		coordinator = jet.NewCoordinatorMock(t)
	}

	ref := gen.Reference()
	jetID := gen.JetID()
	requestID := gen.ID()

	request := record.IncomingRequest{
		Object:   &ref,
		CallType: record.CTMethod,
		// APINode:  gen.Reference(),
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
		idxStorage.SetIndexFunc = func(_ context.Context, pn insolar.PulseNumber, idx record.Index) (r error) {
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
					PendingPointer:      pendingID,
				},
			}
			require.Equal(t, expectedIndex, idx)
			return nil
		}

		writeAccessor.BeginMock.Return(func() {}, nil)
		filaments.RequestDuplicateMock.Return(nil, nil, nil)
		sender.ReplyMock.Return()
		coordinator.VirtualExecutorForObjectFunc = func(_ context.Context, objID insolar.ID, pn insolar.PulseNumber) (r *insolar.Reference, r1 error) {
			require.Equal(t, flowPN, pn)
			require.Equal(t, *ref.Record(), objID)

			return &virtualRef, nil
		}
		records.SetFunc = func(_ context.Context, id insolar.ID, rec record.Material) (r error) {
			switch record.Unwrap(&rec.Virtual).(type) {
			case *record.IncomingRequest:
				require.Equal(t, requestID, id)
			case *record.PendingFilament:
				hash := record.HashVirtual(pcs.ReferenceHasher(), rec.Virtual)
				calcID := *insolar.NewID(requestID.Pulse(), hash)
				require.Equal(t, calcID, id)
			default:
				t.Fatal("unknown record saved")
			}

			return nil
		}
		checker.CheckRequestFunc = func(_ context.Context, id insolar.ID, req record.Request) (r error) {
			require.Equal(t, requestID, id)
			require.Equal(t, &request, req)
			return nil
		}

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
		filaments.RequestDuplicateMock.Return(
			&record.CompositeFilamentRecord{RecordID: reqID},
			&record.CompositeFilamentRecord{RecordID: resID},
			nil,
		)

		sender.ReplyFunc = func(_ context.Context, meta payload.Meta, msg *message.Message) {
			pl, err := payload.Unmarshal(msg.Payload)
			require.NoError(t, err)
			rep, ok := pl.(*payload.RequestInfo)
			require.True(t, ok)
			require.Equal(t, reqID, rep.RequestID)
		}
		coordinator.VirtualExecutorForObjectFunc = func(_ context.Context, objID insolar.ID, pn insolar.PulseNumber) (r *insolar.Reference, r1 error) {
			require.Equal(t, flowPN, pn)
			require.Equal(t, *ref.Record(), objID)

			return &virtualRef, nil
		}

		p := proc.NewSetRequest(msg, &request, requestID, jetID)
		p.Dep(writeAccessor, filaments, sender, object.NewIndexLocker(), idxStorage, records, pcs, checker, coordinator)

		err = p.Proceed(ctx)
		require.NoError(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("wrong sender", func(t *testing.T) {
		t.Skip("virtual doesn't pass this check")
		coordinator.VirtualExecutorForObjectFunc = func(_ context.Context, objID insolar.ID, pn insolar.PulseNumber) (r *insolar.Reference, r1 error) {
			require.Equal(t, flowPN, pn)
			require.Equal(t, *ref.Record(), objID)

			virtualRef := gen.Reference()
			return &virtualRef, nil
		}

		p := proc.NewSetRequest(msg, &request, requestID, jetID)
		p.Dep(writeAccessor, filaments, sender, object.NewIndexLocker(), idxStorage, records, pcs, checker, coordinator)

		err = p.Proceed(ctx)
		require.Error(t, err)
		require.Equal(t, err.Error(), proc.ErrExecutorMismatch.Error())

		mc.Finish()
	})
}
