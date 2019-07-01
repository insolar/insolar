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

func TestSetRequest_Proceed(t *testing.T) {
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

	filaments := executor.NewFilamentModifierMock(t)
	filaments.SetRequestMock.Return(nil)

	ref := gen.Reference()
	jetID := gen.JetID()
	id := gen.ID()

	request := record.IncomingRequest{
		Object:   &ref,
		CallType: record.CTMethod,
	}
	virtual := record.Virtual{
		Union: &record.Virtual_Request{
			Request: &request,
		},
	}
	virtualBuf, err := virtual.Marshal()
	require.NoError(t, err)

	pl := payload.SetRequest{
		Request: virtualBuf,
	}
	requestBuf, err := pl.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: requestBuf,
	}

	pmm := object.NewPendingModifierMock(t)
	pmm.SetRequestMock.Return(nil)

	// Pendings limit not reached.
	p := NewSetRequest(msg, request, id, jetID)
	p.Dep(writeAccessor, records, filaments, sender, object.NewIndexLocker())

	err = p.Proceed(ctx)
	require.NoError(t, err)
}
