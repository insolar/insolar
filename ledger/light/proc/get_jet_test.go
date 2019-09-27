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

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/pulse"
)

func TestGetJet_Proceed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	var (
		jetAccessor *jet.AccessorMock
		sender      *bus.SenderMock
	)

	setup := func() {
		jetAccessor = jet.NewAccessorMock(mc)
		sender = bus.NewSenderMock(mc)
	}

	t.Run("basic ok", func(t *testing.T) {
		setup()
		defer mc.Finish()

		jetID := gen.JetID()
		jetAccessor.ForIDMock.Return(jetID, true)

		expectedMsg, _ := payload.NewMessage(&payload.Jet{
			JetID:  jetID,
			Actual: true,
		})

		sender.ReplyMock.Inspect(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
			assert.Equal(t, expectedMsg.Payload, reply.Payload)
		}).Return()

		p := proc.NewGetJet(payload.Meta{}, gen.ID(), pulse.MinTimePulse)
		p.Dep(jetAccessor, sender)
		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})

}
