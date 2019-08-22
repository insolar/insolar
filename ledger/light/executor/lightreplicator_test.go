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

package executor

import (
	"context"
	"testing"
	"time"

	message2 "github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/stretchr/testify/require"
)

func Test_NotifyAboutPulse(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	ctx := context.Background()

	jetID := jet.NewIDFromString("1010")
	expectPN := insolar.PulseNumber(2835341939)
	expectDrop := drop.Drop{
		JetID: jetID,
		Pulse: expectPN,
		Hash:  []byte{4, 2, 3},
	}
	expectIndexes := []record.Index{
		{ObjID: gen.ID()},
		{ObjID: gen.ID()},
	}
	expectRecords := []record.Material{
		{Signature: gen.Signature(256)},
		{Signature: gen.Signature(256)},
	}

	expectPL := payload.Replication{
		Polymorph: uint32(payload.TypeReplication),
		JetID:     jetID,
		Pulse:     expectPN,
		Indexes:   expectIndexes,
		Drop:      drop.MustEncode(&expectDrop),
		Records:   expectRecords,
	}

	sender := bus.NewSenderMock(mc)
	sender.SendRoleMock.Set(func(_ context.Context, msg *message2.Message, role insolar.DynamicRole, _ insolar.Reference) (r <-chan *message2.Message, r1 func()) {
		pl, err := payload.Unmarshal(msg.Payload)
		require.NoError(t, err)
		require.Equal(t, &expectPL, pl, "heavy message payload")
		return nil, func() {}
	})

	jetCalc := NewJetCalculatorMock(mc)
	jetCalc.MineForPulseMock.Set(func(_ context.Context, _ insolar.PulseNumber) ([]insolar.JetID, error) {
		return []insolar.JetID{jetID}, nil
	})

	cleaner := NewCleanerMock(mc)
	cleaner.NotifyAboutPulseMock.Set(func(_ context.Context, _ insolar.PulseNumber) {})

	pulseCalc := pulse.NewCalculatorMock(mc)
	pulseCalc.BackwardsMock.Expect(ctx, expectPN+1, 1).Return(
		insolar.Pulse{PulseNumber: expectPN}, nil)

	dropAccessor := drop.NewAccessorMock(mc)
	dropAccessor.ForPulseMock.Set(func(_ context.Context, _ insolar.JetID, _ insolar.PulseNumber) (r drop.Drop, r1 error) {
		return expectDrop, nil
	})

	recordAccessor := object.NewRecordCollectionAccessorMock(mc)
	recordAccessor.ForPulseMock.Set(func(_ context.Context, _ insolar.JetID, _ insolar.PulseNumber) (r []record.Material) {
		return expectRecords
	})

	indexAccessor := object.NewIndexAccessorMock(mc)
	indexAccessor.ForPulseMock.Set(func(_ context.Context, _ insolar.PulseNumber) ([]record.Index, error) {
		return expectIndexes, nil
	})

	jetAccessor := jet.NewAccessorMock(mc)
	jetAccessor.ForIDMock.Set(func(_ context.Context, _ insolar.PulseNumber, _ insolar.ID) (insolar.JetID, bool) {
		return jetID, false
	})

	r := NewReplicatorDefault(
		jetCalc,
		cleaner,
		sender,
		pulseCalc,
		dropAccessor,
		recordAccessor,
		indexAccessor,
		jetAccessor,
	)
	defer close(r.syncWaitingPulses)

	r.NotifyAboutPulse(ctx, expectPN+1)
	mc.Wait(time.Minute)
	mc.Finish()
}
