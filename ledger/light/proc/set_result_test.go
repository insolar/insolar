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
	"github.com/stretchr/testify/require"
)

func TestSetResult_Proceed(t *testing.T) {
	t.Parallel()

	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	writeAccessor := hot.NewWriteAccessorMock(t)
	writeAccessor.BeginMock.Return(func() {}, nil)

	sender := bus.NewSenderMock(t)
	sender.ReplyMock.Return()

	jetID := gen.JetID()
	id := gen.ID()

	res := &record.Result{
		Object: id,
	}
	virtual := record.Virtual{
		Union: &record.Virtual_Result{
			Result: res,
		},
	}
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

	filamentModifier := executor.NewFilamentManagerMock(t)
	filamentModifier.SetResultFunc = func(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Result) (fRes *record.CompositeFilamentRecord, r error) {
		require.Equal(t, id, p1)
		require.Equal(t, jetID, p2)
		require.Equal(t, *res, p3)

		return nil, nil
	}

	// Pendings limit not reached.
	setResultProc := proc.NewSetResult(msg, *res, id, jetID)
	setResultProc.Dep(writeAccessor, sender, object.NewIndexLocker(), filamentModifier)

	err = setResultProc.Proceed(ctx)
	require.NoError(t, err)
}

func TestSetResult_Proceed_ResultDuplicated(t *testing.T) {
	t.Parallel()

	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	writeAccessor := hot.NewWriteAccessorMock(t)
	writeAccessor.BeginMock.Return(func() {}, nil)

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

	filamentModifier := executor.NewFilamentManagerMock(t)
	filamentModifier.SetResultFunc = func(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Result) (fRes *record.CompositeFilamentRecord, r error) {
		require.Equal(t, objectID, p1)
		require.Equal(t, jetID, p2)
		require.Equal(t, *res, p3)

		return nil, nil
	}

	// Pendings limit not reached.
	setResultProc := proc.NewSetResult(msg, *res, objectID, jetID)
	setResultProc.Dep(writeAccessor, sender, object.NewIndexLocker(), filamentModifier)
	sender.ReplyFunc = func(_ context.Context, receivedMeta payload.Meta, resMsg *message.Message) {
		require.Equal(t, msg, receivedMeta)

		resp, err := payload.Unmarshal(resMsg.Payload)
		require.NoError(t, err)

		res, ok := resp.(*payload.ResultInfo)
		require.True(t, ok)
		require.Nil(t, res.Result)
		require.Equal(t, objectID, res.ResultID)
	}

	err = setResultProc.Proceed(ctx)
	require.NoError(t, err)

	filamentModifier.SetResultFunc = func(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Result) (fRes *record.CompositeFilamentRecord, r error) {
		require.Equal(t, objectID, p1)
		require.Equal(t, jetID, p2)
		require.Equal(t, *res, p3)

		return &record.CompositeFilamentRecord{
			Record:   record.Material{Virtual: &virtual},
			RecordID: resultID,
		}, nil
	}
	sender.ReplyFunc = func(_ context.Context, receivedMeta payload.Meta, resMsg *message.Message) {
		require.Equal(t, msg, receivedMeta)

		resp, err := payload.Unmarshal(resMsg.Payload)
		require.NoError(t, err)

		res, ok := resp.(*payload.ResultInfo)
		require.True(t, ok)
		require.Equal(t, virtualBuf, res.Result)
		require.Equal(t, resultID, res.ResultID)
	}

	// CheckDuplication
	err = setResultProc.Proceed(ctx)
	require.NoError(t, err)
}
