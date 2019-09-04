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
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/proc"
)

func TestSendFilament_Proceed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	var (
		sender    *bus.SenderMock
		filaments *executor.FilamentCalculatorMock
	)

	setup := func() {
		sender = bus.NewSenderMock(mc)
		filaments = executor.NewFilamentCalculatorMock(mc)
	}

	t.Run("simple success", func(t *testing.T) {
		setup()
		defer mc.Finish()

		obj := gen.ID()
		pl, _ := (&payload.GetFilament{}).Marshal()

		msg := payload.Meta{
			Payload: pl,
		}

		storageRecs := make([]record.CompositeFilamentRecord, 5)
		filaments.RequestsMock.Return(storageRecs, nil)
		expectedMsg, _ := payload.NewMessage(&payload.FilamentSegment{
			ObjectID: obj,
			Records:  storageRecs,
		})

		sender.ReplyMock.Inspect(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
			assert.Equal(t, expectedMsg.Payload, reply.Payload)
			assert.Equal(t, msg, origin)
		}).Return()

		p := proc.NewSendFilament(msg, obj, gen.ID(), gen.PulseNumber())
		p.Dep(sender, filaments)

		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})

	t.Run("requests not found sends error", func(t *testing.T) {
		setup()
		defer mc.Finish()

		obj := gen.ID()
		pl, _ := (&payload.GetFilament{}).Marshal()

		msg := payload.Meta{
			Payload: pl,
		}

		var storageRecs []record.CompositeFilamentRecord
		filaments.RequestsMock.Return(storageRecs, nil)

		p := proc.NewSendFilament(msg, obj, gen.ID(), gen.PulseNumber())
		p.Dep(sender, filaments)

		err := p.Proceed(ctx)
		assert.Error(t, err)
		insError, ok := errors.Cause(err).(*payload.CodedError)
		assert.True(t, ok)
		assert.Equal(t, uint32(payload.CodeNotFound), insError.GetCode())

	})
}
