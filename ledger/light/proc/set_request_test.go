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
	"github.com/stretchr/testify/require"
)

func TestSetRequest_Proceed(t *testing.T) {
	t.Parallel()

	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)
	mc := minimock.NewController(t)

	var (
		writeAccessor *hot.WriteAccessorMock
		sender        *bus.SenderMock
		filaments     *executor.FilamentModifierMock
	)

	resetComponents := func() {
		writeAccessor = hot.NewWriteAccessorMock(mc)
		sender = bus.NewSenderMock(mc)
		filaments = executor.NewFilamentModifierMock(t)
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

	msg := payload.Meta{
		Payload: requestBuf,
	}

	resetComponents()
	t.Run("happy basic", func(t *testing.T) {
		writeAccessor.BeginMock.Return(func() {}, nil)
		sender.ReplyMock.Return()
		filaments.SetRequestMock.Return(nil, nil, nil)

		p := proc.NewSetRequest(msg, &request, id, jetID)
		p.Dep(writeAccessor, filaments, sender, object.NewIndexLocker())

		err = p.Proceed(ctx)
		require.NoError(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("duplicate returns correct id", func(t *testing.T) {
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

		p := proc.NewSetRequest(msg, &request, id, jetID)
		p.Dep(writeAccessor, filaments, sender, object.NewIndexLocker())

		err = p.Proceed(ctx)
		require.NoError(t, err)

		mc.Finish()
	})
}
