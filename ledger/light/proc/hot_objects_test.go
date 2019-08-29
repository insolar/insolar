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
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/pulse"
)

func TestHotObjects_Proceed(t *testing.T) {
	ctx := flow.TestContextWithPulse(inslogger.TestContext(t), pulse.MinTimePulse+10)
	mc := minimock.NewController(t)

	var (
		drops       *drop.ModifierMock
		indexes     *object.MemoryIndexModifierMock
		jetStorage  *jet.StorageMock
		jetFetcher  *executor.JetFetcherMock
		jetReleaser *executor.JetReleaserMock
		coordinator *jet.CoordinatorMock
		calculator  *insolarPulse.CalculatorMock
		sender      *bus.SenderMock
	)

	setup := func(mc minimock.MockController) {
		drops = drop.NewModifierMock(mc)
		indexes = object.NewMemoryIndexModifierMock(mc)
		jetStorage = jet.NewStorageMock(mc)
		jetFetcher = executor.NewJetFetcherMock(mc)
		jetReleaser = executor.NewJetReleaserMock(mc)
		coordinator = jet.NewCoordinatorMock(mc)
		calculator = insolarPulse.NewCalculatorMock(mc)
		sender = bus.NewSenderMock(mc)
	}

	t.Run("basic ok", func(t *testing.T) {
		setup(mc)
		defer mc.Finish()

		expectedPulse := insolar.Pulse{
			PulseNumber: pulse.MinTimePulse + 10,
		}
		expectedJetID := gen.JetID()
		expectedObjJetID := expectedJetID
		meta := payload.Meta{}
		expectedDrop := drop.Drop{
			Pulse: expectedPulse.PulseNumber,
			JetID: expectedJetID,
			Split: false,
			Hash:  []byte{1, 2, 3},
		}
		idxs := []record.Index{
			{
				ObjID: gen.ID(),
				// this is hack, PendingRecords in record.Index should be always empty
				PendingRecords: []insolar.ID{},
			},
		}

		drops.SetMock.Inspect(func(ctx context.Context, drop drop.Drop) {
			assert.Equal(t, expectedDrop, drop, "didn't set drop")
		}).Return(nil)

		jetStorage.UpdateMock.Inspect(func(ctx context.Context, pulse insolar.PulseNumber, actual bool, ids ...insolar.JetID) {
			assert.Equal(t, expectedPulse.PulseNumber, pulse, "wrong pulse received")
			assert.Equal(t, expectedJetID, ids[0], "wrong jetID received")
		}).Return(nil)

		calculator.BackwardsMock.Return(insolar.Pulse{}, insolarPulse.ErrNotFound)

		jetStorage.ForIDMock.Return(expectedObjJetID, false)

		indexes.SetMock.Inspect(func(ctx context.Context, pn insolar.PulseNumber, index record.Index) {
			assert.Equal(t, idxs[0], index)
		}).Return()

		jetFetcher.ReleaseMock.Inspect(func(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) {
			assert.Equal(t, expectedJetID, jetID)
		}).Return()
		jetReleaser.UnlockMock.Inspect(func(ctx context.Context, pulse insolar.PulseNumber, jetID insolar.JetID) {
			assert.Equal(t, expectedJetID, jetID)
		}).Return(nil)

		expectedToHeavyMsg, _ := payload.NewMessage(&payload.GotHotConfirmation{
			JetID: expectedJetID,
			Pulse: expectedPulse.PulseNumber,
			Split: expectedDrop.Split,
		})

		sender.SendRoleMock.Inspect(func(ctx context.Context, msg *message.Message, role insolar.DynamicRole, object insolar.Reference) {
			assert.Equal(t, expectedToHeavyMsg.Payload, msg.Payload)
		}).Return(make(chan *message.Message), func() {})

		// start test
		p := proc.NewHotObjects(meta, expectedPulse.PulseNumber, expectedJetID, expectedDrop, idxs)
		p.Dep(drops, indexes, jetStorage, jetFetcher, jetReleaser, coordinator, calculator, sender)

		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})

	t.Run("ok, with pendings", func(t *testing.T) {
		setup(mc)
		defer mc.Finish()

		currentPulse := insolar.Pulse{
			PulseNumber: pulse.MinTimePulse + 100,
		}
		abandonedRequestPulse := insolar.Pulse{
			PulseNumber: pulse.MinTimePulse,
		}
		thresholdAbandonedRequestPulse := insolar.Pulse{
			PulseNumber: pulse.MinTimePulse + 80,
		}

		expectedJetID := gen.JetID()
		expectedObjJetID := expectedJetID
		expectedObjectID := gen.ID()
		meta := payload.Meta{}
		expectedDrop := drop.Drop{
			Pulse: currentPulse.PulseNumber,
			JetID: expectedJetID,
			Split: false,
			Hash:  []byte{1, 2, 3},
		}
		idxs := []record.Index{
			{
				ObjID: expectedObjectID,
				Lifeline: record.Lifeline{
					EarliestOpenRequest: &abandonedRequestPulse.PulseNumber,
				},
				// this is hack, PendingRecords in record.Index should be always empty
				PendingRecords: []insolar.ID{},
			},
		}

		drops.SetMock.Inspect(func(ctx context.Context, drop drop.Drop) {
			assert.Equal(t, expectedDrop, drop, "didn't set drop")
		}).Return(nil)

		jetStorage.UpdateMock.Inspect(func(ctx context.Context, pulse insolar.PulseNumber, actual bool, ids ...insolar.JetID) {
			assert.Equal(t, currentPulse.PulseNumber, pulse, "wrong pulse received")
			assert.Equal(t, expectedJetID, ids[0], "wrong jetID received")
		}).Return(nil)

		calculator.BackwardsMock.Return(thresholdAbandonedRequestPulse, nil)
		jetStorage.ForIDMock.Return(expectedObjJetID, false)

		indexes.SetMock.Inspect(func(ctx context.Context, pn insolar.PulseNumber, index record.Index) {
			assert.Equal(t, idxs[0], index)
		}).Return()

		jetFetcher.ReleaseMock.Inspect(func(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) {
			assert.Equal(t, expectedJetID, jetID)
		}).Return()
		jetReleaser.UnlockMock.Inspect(func(ctx context.Context, pulse insolar.PulseNumber, jetID insolar.JetID) {
			assert.Equal(t, expectedJetID, jetID)
		}).Return(nil)

		expectedToHeavyMsg, _ := payload.NewMessage(&payload.GotHotConfirmation{
			JetID: expectedJetID,
			Pulse: currentPulse.PulseNumber,
			Split: expectedDrop.Split,
		})

		expectedToVirtualMsg, _ := payload.NewMessage(&payload.AbandonedRequestsNotification{
			ObjectID: expectedObjectID,
		})

		sender.SendRoleMock.Inspect(func(ctx context.Context, msg *message.Message, role insolar.DynamicRole, object insolar.Reference) {
			if role == insolar.DynamicRoleHeavyExecutor {
				assert.Equal(t, expectedToHeavyMsg.Payload, msg.Payload)
				return
			} else if role == insolar.DynamicRoleVirtualExecutor {
				assert.Equal(t, expectedToVirtualMsg.Payload, msg.Payload)
				return
			}
			assert.True(t, false, "didn't receive at least 2 messages")
		}).Return(make(chan *message.Message), func() {})

		// start test
		p := proc.NewHotObjects(meta, currentPulse.PulseNumber, expectedJetID, expectedDrop, idxs)
		p.Dep(drops, indexes, jetStorage, jetFetcher, jetReleaser, coordinator, calculator, sender)

		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})
}
