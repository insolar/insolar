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
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/proc"
)

func TestGetPendings_Proceed(t *testing.T) {
	ctx := flow.TestContextWithPulse(inslogger.TestContext(t), insolar.FirstPulseNumber+10)
	mc := minimock.NewController(t)

	var (
		filaments *executor.FilamentCalculatorMock
		sender    *bus.SenderMock
	)

	setup := func() {
		filaments = executor.NewFilamentCalculatorMock(mc)
		sender = bus.NewSenderMock(mc)
	}

	t.Run("ok, pendings is empty", func(t *testing.T) {
		setup()
		defer mc.Finish()

		filaments.OpenedRequestsMock.Return([]record.CompositeFilamentRecord{}, nil)

		expectedMsg, _ := payload.NewMessage(&payload.Error{
			Code: payload.CodeNoPendings,
			Text: insolar.ErrNoPendingRequest.Error(),
		})

		sender.ReplyMock.Inspect(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
			assert.Equal(t, expectedMsg.Payload, reply.Payload)
		}).Return()

		p := proc.NewGetPendings(payload.Meta{}, gen.ID())
		p.Dep(filaments, sender)

		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})

	t.Run("ok, pendings found", func(t *testing.T) {
		setup()
		defer mc.Finish()
		pendings := []record.CompositeFilamentRecord{
			{RecordID: gen.ID()},
			{RecordID: gen.ID()},
			{RecordID: gen.ID()},
			{RecordID: gen.ID()},
		}

		ids := make([]insolar.ID, len(pendings))
		for i, pend := range pendings {
			ids[i] = pend.RecordID
		}

		expectedMsg, _ := payload.NewMessage(&payload.IDs{
			IDs: ids,
		})

		filaments.OpenedRequestsMock.Return(pendings, nil)

		sender.ReplyMock.Inspect(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
			assert.Equal(t, expectedMsg.Payload, reply.Payload)
		}).Return()

		p := proc.NewGetPendings(payload.Meta{}, gen.ID())
		p.Dep(filaments, sender)

		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})
}
