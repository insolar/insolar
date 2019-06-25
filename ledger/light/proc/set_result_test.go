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

package proc

import (
	"context"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/hot"
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

	records := object.NewRecordModifierMock(t)
	records.SetMock.Return(nil)

	jetID := gen.JetID()
	id := gen.ID()

	res := &record.Result{
		Object: id,
	}

	filamentMod := executor.NewFilamentModifierMock(t)
	filamentMod.SetResultFunc = func(p context.Context, p1 insolar.ID, p2 insolar.JetID, p3 record.Result) (r error) {
		require.Equal(t, jetID, p2)
		require.Equal(t, id, p1)
		return nil
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

	// Pendings limit not reached.
	setResultProc := NewSetResult(msg, *res, id, jetID)
	setResultProc.dep.writer = writeAccessor
	setResultProc.dep.sender = sender
	setResultProc.dep.records = records
	setResultProc.dep.locker = object.NewIndexLocker()
	setResultProc.dep.filament = filamentMod

	err = setResultProc.Proceed(ctx)
	require.NoError(t, err)
}
