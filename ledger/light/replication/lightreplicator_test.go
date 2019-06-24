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

package replication

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestLightReplicator_sendToHeavy(t *testing.T) {
	t.Parallel()
	mb := testutils.NewMessageBusMock(t)
	mb.SendMock.Return(&reply.OK{}, nil)
	r := &LightReplicatorDefault{
		msgBus: mb,
	}

	res := r.sendToHeavy(inslogger.TestContext(t), nil)
	require.Nil(t, res)
}

func TestLightReplicator_sendToHeavy_ErrReturned(t *testing.T) {
	t.Parallel()
	mb := testutils.NewMessageBusMock(t)
	mb.SendMock.Return(nil, errors.New("expected"))
	r := LightReplicatorDefault{
		msgBus: mb,
	}

	res := r.sendToHeavy(inslogger.TestContext(t), nil)

	require.Equal(t, res, errors.New("expected"))
}

func TestLightReplicator_sendToHeavy_HeavyErr(t *testing.T) {
	t.Parallel()
	mb := testutils.NewMessageBusMock(t)
	heavyErr := reply.HeavyError{JetID: gen.JetID(), PulseNum: gen.PulseNumber()}
	mb.SendMock.Return(&heavyErr, nil)
	r := LightReplicatorDefault{
		msgBus: mb,
	}

	res := r.sendToHeavy(inslogger.TestContext(t), nil)

	require.Equal(t, &heavyErr, res)
}

func Test_NotifyAboutPulse(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)
	ctx := inslogger.TestContext(t)

	jetID := jet.NewIDFromString("1010")
	expectPN := insolar.PulseNumber(2835341939)
	expectDrop := drop.Drop{
		JetID: jetID,
		Pulse: expectPN,
		Hash:  []byte{4, 2, 3},
	}
	expectBlob := blob.Blob{
		JetID: gen.JetID(),
		Value: []byte{1, 2, 3, 4, 5, 6, 7},
	}
	expectIndexes := []object.FilamentIndex{
		{ObjID: gen.ID()},
		{ObjID: gen.ID()},
	}
	expectRecords := []record.Material{
		{Signature: gen.Signature(256)},
		{Signature: gen.Signature(256)},
	}

	expectMsg := &message.HeavyPayload{
		JetID:        jetID,
		PulseNum:     expectPN,
		IndexBuckets: convertIndexBuckets(ctx, expectIndexes),
		Drop:         drop.MustEncode(&expectDrop),
		Blobs:        convertBlobs([]blob.Blob{expectBlob}),
		Records:      convertRecords(ctx, expectRecords),
	}

	mb := testutils.NewMessageBusMock(ctrl)
	mb.SendFunc = func(_ context.Context, msg insolar.Message, opts *insolar.MessageSendOptions) (insolar.Reply, error) {
		require.IsType(t, &message.HeavyPayload{}, msg, "got heavy payload message")
		hMsg := msg.(*message.HeavyPayload)
		require.Equal(t, expectMsg, hMsg, "heavy message payload")
		return &reply.OK{}, nil
	}

	jetCalc := executor.NewJetCalculatorMock(ctrl)
	jetCalc.MineForPulseFunc = func(_ context.Context, _ insolar.PulseNumber) []insolar.JetID {
		return []insolar.JetID{jetID}
	}

	cleaner := NewCleanerMock(ctrl)
	cleaner.NotifyAboutPulseMock.Expect(ctx, expectPN)

	pulseCalc := pulse.NewCalculatorMock(ctrl)
	pulseCalc.BackwardsMock.Expect(ctx, expectPN+1, 1).Return(
		insolar.Pulse{PulseNumber: expectPN}, nil)

	dropAccessor := drop.NewAccessorMock(ctrl)
	dropAccessor.ForPulseMock.Expect(ctx, jetID, expectPN).Return(expectDrop, nil)

	blobAccessor := blob.NewCollectionAccessorMock(ctrl)
	blobAccessor.ForPulseMock.Expect(ctx, jetID, expectPN).Return([]blob.Blob{expectBlob})

	recordAccessor := object.NewRecordCollectionAccessorMock(ctrl)
	recordAccessor.ForPulseMock.Expect(ctx, jetID, expectPN).Return(expectRecords)

	indexAccessor := object.NewIndexBucketAccessorMock(ctrl)
	indexAccessor.ForPulseFunc = func(_ context.Context, _ insolar.PulseNumber) []object.FilamentIndex {
		return expectIndexes
	}

	jetAccessor := jet.NewAccessorMock(ctrl)
	jetAccessor.ForIDFunc = func(_ context.Context, _ insolar.PulseNumber, _ insolar.ID) (insolar.JetID, bool) {
		return jetID, false
	}

	r := NewReplicatorDefault(
		jetCalc,
		cleaner,
		mb,
		pulseCalc,
		dropAccessor,
		blobAccessor,
		recordAccessor,
		indexAccessor,
		jetAccessor,
	)
	defer close(r.syncWaitingPulses)

	r.NotifyAboutPulse(ctx, expectPN+1)
	ctrl.Wait(time.Minute)
	ctrl.Finish()
}
