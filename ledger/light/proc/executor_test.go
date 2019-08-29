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
	"errors"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/pulse"
)

func TestFetchJet_Proceed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	var (
		jetAccessor *jet.AccessorMock
		jetFetcher  *executor.JetFetcherMock
		coordinator *jet.CoordinatorMock
		sender      *bus.SenderMock
	)
	setup := func() {
		jetAccessor = jet.NewAccessorMock(mc)
		jetFetcher = executor.NewJetFetcherMock(mc)
		coordinator = jet.NewCoordinatorMock(mc)
		sender = bus.NewSenderMock(mc)

	}

	t.Run("error wrong jetID", func(t *testing.T) {
		setup()
		defer mc.Finish()

		jetFetcher.FetchMock.Return(&insolar.ID{}, errors.New("test"))

		p := proc.NewFetchJet(gen.ID(), pulse.MinTimePulse, payload.Meta{}, false)
		p.Dep(jetAccessor, jetFetcher, coordinator, sender)
		err := p.Proceed(ctx)
		assert.Error(t, err, "expected error 'failed to fetch jet'")
	})

	t.Run("virtual is executor", func(t *testing.T) {
		setup()
		defer mc.Finish()

		id := gen.ID()
		jetFetcher.FetchMock.Return(&id, nil)
		me := insolar.NewReference(gen.ID())
		coordinator.LightExecutorForJetMock.Return(me, nil)
		coordinator.MeMock.Return(*me)

		p := proc.NewFetchJet(gen.ID(), pulse.MinTimePulse, payload.Meta{}, false)
		p.Dep(jetAccessor, jetFetcher, coordinator, sender)
		err := p.Proceed(ctx)
		assert.NoError(t, err)
		assert.Equal(t, p.Result.Jet, insolar.JetID(id))
	})

	t.Run("virtual is not executor", func(t *testing.T) {
		setup()
		defer mc.Finish()

		id := gen.ID()
		jetFetcher.FetchMock.Return(&id, nil)
		me := insolar.NewReference(gen.ID())
		notMe := insolar.NewReference(gen.ID())
		coordinator.LightExecutorForJetMock.Return(me, nil)
		coordinator.MeMock.Return(*notMe)

		p := proc.NewFetchJet(gen.ID(), pulse.MinTimePulse, payload.Meta{}, false)
		p.Dep(jetAccessor, jetFetcher, coordinator, sender)
		err := p.Proceed(ctx)
		assert.Error(t, err)
	})

	t.Run("virtual passing", func(t *testing.T) {
		setup()
		defer mc.Finish()

		id := gen.ID()
		jetFetcher.FetchMock.Return(&id, nil)
		msgSender := insolar.NewReference(gen.ID())
		me := insolar.NewReference(gen.ID())
		notMe := insolar.NewReference(gen.ID())
		coordinator.LightExecutorForJetMock.Return(notMe, nil)
		coordinator.MeMock.Return(*me)

		reps := make(chan *message.Message, 1)

		sender.SendTargetMock.Inspect(func(ctx context.Context, msg *message.Message, target insolar.Reference) {
			assert.True(t, target == *notMe || target == *msgSender, "expected messages to msgSender and to right executor")
		}).Return(reps, func() {})

		p := proc.NewFetchJet(gen.ID(), pulse.MinTimePulse, payload.Meta{Sender: *msgSender}, true)
		p.Dep(jetAccessor, jetFetcher, coordinator, sender)
		err := p.Proceed(ctx)
		assert.Error(t, err)
		assert.Equal(t, err, proc.ErrNotExecutor)
	})

}
