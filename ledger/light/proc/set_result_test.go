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

func TestSetResult_Proceed(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	flowPulse := insolar.GenesisPulse.PulseNumber + 2
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		flowPulse,
	)

	writeAccessor := hot.NewWriteAccessorMock(mc)
	writeAccessor.BeginMock.Return(func() {}, nil)
	pcs := testutils.NewPlatformCryptographyScheme()

	sender := bus.NewSenderMock(t)
	sender.ReplyMock.Return()

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
	pendingPointer := gen.IDWithPulse(flowPulse)
	expectedFilament := record.PendingFilament{
		RecordID:       resultID,
		PreviousRecord: &pendingPointer,
	}
	hash = record.HashVirtual(pcs.ReferenceHasher(), record.Wrap(&expectedFilament))
	expectedFilamentID := *insolar.NewID(resultID.Pulse(), hash)

	indexes := object.NewIndexStorageMock(mc)
	indexes.ForIDFunc = func(_ context.Context, pn insolar.PulseNumber, id insolar.ID) (record.Index, error) {
		require.Equal(t, flow.Pulse(ctx), pn)
		require.Equal(t, objectID, id)
		earliestPN := requestID.Pulse()
		return record.Index{
			Lifeline: record.Lifeline{
				PendingPointer:      &pendingPointer,
				EarliestOpenRequest: &earliestPN,
			},
		}, nil
	}
	indexes.SetIndexFunc = func(_ context.Context, pn insolar.PulseNumber, idx record.Index) (r error) {
		require.Equal(t, resultID.Pulse(), pn)
		expectedIndex := record.Index{
			LifelineLastUsed: resultID.Pulse(),
			Lifeline: record.Lifeline{
				PendingPointer:      &expectedFilamentID,
				EarliestOpenRequest: nil,
			},
		}
		require.Equal(t, expectedIndex, idx)
		return nil
	}
	records := object.NewAtomicRecordModifierMock(mc)
	records.SetAtomicFunc = func(_ context.Context, recs ...record.Material) (r error) {
		require.Equal(t, 1, len(recs))
		rec := recs[0]

		switch r := record.Unwrap(&rec.Virtual).(type) {
		case *record.Result:
			require.Equal(t, resultID, rec.ID)
			require.Equal(t, resultRecord, r)
		case *record.PendingFilament:
			require.Equal(t, expectedFilamentID, rec.ID)
			require.Equal(t, &expectedFilament, record.Unwrap(&rec.Virtual))
		}

		return nil
	}

	filaments := executor.NewFilamentCalculatorMock(mc)
	filaments.ResultDuplicateFunc = func(_ context.Context, objID insolar.ID, resID insolar.ID, r record.Result) (*record.CompositeFilamentRecord, error) {
		require.Equal(t, objectID, objID)
		require.Equal(t, *resultRecord, r)
		return nil, nil
	}
	filaments.OpenedRequestsFunc = func(_ context.Context, pn insolar.PulseNumber, objID insolar.ID, pendingOnly bool) ([]record.CompositeFilamentRecord, error) {
		require.Equal(t, objectID, objID)
		require.Equal(t, flow.Pulse(ctx), pn)
		require.False(t, pendingOnly)

		v := record.Wrap(&record.IncomingRequest{})
		return []record.CompositeFilamentRecord{
			{
				RecordID: requestID,
				Record:   record.Material{Virtual: v},
			},
		}, nil
	}

	setResultProc := proc.NewSetResult(msg, jetID, *resultRecord, nil)
	setResultProc.Dep(writeAccessor, sender, object.NewIndexLocker(), filaments, records, indexes, pcs)

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

	writeAccessor := hot.NewWriteAccessorMock(mc)
	writeAccessor.BeginMock.Return(func() {}, nil)
	records := object.NewAtomicRecordModifierMock(mc)
	indexes := object.NewIndexStorageMock(mc)
	indexes.ForIDMock.Return(record.Index{}, nil)
	pcs := testutils.NewPlatformCryptographyScheme()

	sender := bus.NewSenderMock(t)

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
	filaments.ResultDuplicateFunc = func(_ context.Context, objID insolar.ID, resID insolar.ID, r record.Result) (*record.CompositeFilamentRecord, error) {
		require.Equal(t, objectID, objID)
		require.Equal(t, *res, r)

		return &record.CompositeFilamentRecord{
			Record:   record.Material{Virtual: virtual},
			RecordID: resultID,
		}, nil
	}
	sender.ReplyFunc = func(_ context.Context, receivedMeta payload.Meta, resMsg *message.Message) {
		require.Equal(t, msg, receivedMeta)

		resp, err := payload.Unmarshal(resMsg.Payload)
		require.NoError(t, err)

		res, ok := resp.(*payload.ResultInfo)
		require.True(t, ok)
		require.Equal(t, duplicateBuf, res.Result)
		require.Equal(t, resultID, res.ResultID)
	}

	setResultProc := proc.NewSetResult(msg, jetID, *res, nil)
	setResultProc.Dep(writeAccessor, sender, object.NewIndexLocker(), filaments, records, indexes, pcs)
	err = setResultProc.Proceed(ctx)
	require.NoError(t, err)
}
