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

	var (
		writeAccessor *hot.WriteAccessorMock
		sender        *bus.SenderMock
		filaments     *executor.FilamentModifierMock
		idxStorage    *object.IndexStorageMock
		coordinator   *jet.CoordinatorMock
	)

	resetComponents := func() {
		writeAccessor = hot.NewWriteAccessorMock(mc)
		sender = bus.NewSenderMock(mc)
		filaments = executor.NewFilamentModifierMock(t)
		idxStorage = object.NewIndexStorageMock(t)
		coordinator = jet.NewCoordinatorMock(t)
	}

	ref := gen.Reference()
	jetID := gen.JetID()
	id := gen.ID()

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

		writeAccessor.BeginMock.Return(func() {}, nil)
		sender.ReplyMock.Return()
		filaments.SetRequestMock.Return(nil, nil, nil)
		coordinator.VirtualExecutorForObjectFunc = func(_ context.Context, objID insolar.ID, pn insolar.PulseNumber) (r *insolar.Reference, r1 error) {
			require.Equal(t, flowPN, pn)
			require.Equal(t, *ref.Record(), objID)

			return &virtualRef, nil
		}

		p := proc.NewSetRequest(msg, &request, id, jetID)
		p.Dep(writeAccessor, filaments, sender, object.NewIndexLocker(), idxStorage, coordinator)

		err = p.Proceed(ctx)
		require.NoError(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("duplicate returns correct id", func(t *testing.T) {
		idxStorage.ForIDMock.Return(record.Index{
			Lifeline: record.Lifeline{
				StateID: record.StateActivation,
			},
		}, nil)

		writeAccessor.BeginMock.Return(func() {}, nil)
		reqID := gen.ID()
		resID := gen.ID()
		filaments.SetRequestMock.Return(
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

		p := proc.NewSetRequest(msg, &request, id, jetID)
		p.Dep(writeAccessor, filaments, sender, object.NewIndexLocker(), idxStorage, coordinator)

		err = p.Proceed(ctx)
		require.NoError(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("wrong sender", func(t *testing.T) {
		idxStorage.ForIDMock.Return(record.Index{
			Lifeline: record.Lifeline{
				StateID: record.StateActivation,
			},
		}, nil)

		writeAccessor.BeginMock.Return(func() {}, nil)
		coordinator.VirtualExecutorForObjectFunc = func(_ context.Context, objID insolar.ID, pn insolar.PulseNumber) (r *insolar.Reference, r1 error) {
			require.Equal(t, flowPN, pn)
			require.Equal(t, *ref.Record(), objID)

			virtualRef := gen.Reference()
			return &virtualRef, nil
		}

		p := proc.NewSetRequest(msg, &request, id, jetID)
		p.Dep(writeAccessor, filaments, sender, object.NewIndexLocker(), idxStorage, coordinator)

		err = p.Proceed(ctx)
		require.Error(t, err)
		require.Equal(t, err.Error(), proc.ErrExecutorMismatch.Error())

		mc.Finish()
	})
}

func TestDeactivateObject_ObjectIsDeactivated(t *testing.T) {
	t.Parallel()

	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	writeAccessor := hot.NewWriteAccessorMock(t)
	writeAccessor.BeginMock.Return(func() {}, nil)

	idxLockMock := object.NewIndexLockerMock(t)
	idxLockMock.LockMock.Return()
	idxLockMock.UnlockMock.Return()

	idxStorageMock := object.NewIndexStorageMock(t)
	idxStorageMock.ForIDMock.Return(record.Index{
		Lifeline: record.Lifeline{
			StateID: record.StateDeactivation,
		},
	}, nil)
	idxStorageMock.SetIndexMock.Return(nil)

	sender := bus.NewSenderMock(t)
	sender.ReplyFunc = func(_ context.Context, _ payload.Meta, inMsg *message.Message) {
		resp, err := payload.Unmarshal(inMsg.Payload)
		require.NoError(t, err)

		res, ok := resp.(*payload.Error)
		require.True(t, ok)
		require.Equal(t, payload.CodeDeactivated, int(res.Code))
	}

	p := proc.NewDeactivateObject(
		payload.Meta{},
		record.Deactivate{},
		gen.ID(),
		record.Result{},
		gen.ID(),
		gen.JetID(),
	)
	p.Dep(
		writeAccessor,
		idxLockMock,
		nil,
		idxStorageMock,
		nil,
		sender,
	)

	err := p.Proceed(ctx)
	require.NoError(t, err)
}
