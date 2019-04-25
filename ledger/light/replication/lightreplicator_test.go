/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package replication

import (
	"errors"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestNewReplicatorDefault(t *testing.T) {
	t.Parallel()
	r := NewReplicatorDefault(
		jet.NewCalculatorMock(t),
		NewDataGathererMock(t),
		NewCleanerMock(t),
		testutils.NewMessageBusMock(t),
		pulse.NewCalculatorMock(t),
	)
	defer close(r.syncWaitingPulses)

	require.NotNil(t, r.jetCalculator)
	require.NotNil(t, r.dataGatherer)
	require.NotNil(t, r.cleaner)
	require.NotNil(t, r.msgBus)
	require.NotNil(t, r.pulseCalculator)
	require.NotNil(t, r.syncWaitingPulses)
}

func TestLightReplicator_sendToHeavy(t *testing.T) {
	t.Parallel()
	mb := testutils.NewMessageBusMock(t)
	mb.SendMock.Return(&reply.OK{}, nil)
	r := LightReplicatorDefault{
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

func TestLightReplicatorDefault_sync(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)
	ctx := inslogger.TestContext(t)
	jc := jet.NewCalculatorMock(ctrl)
	c := NewCleanerMock(ctrl)
	mb := testutils.NewMessageBusMock(ctrl)
	pc := pulse.NewCalculatorMock(ctrl)
	dg := NewDataGathererMock(ctrl)
	r := NewReplicatorDefault(
		jc,
		dg,
		c,
		mb,
		pc,
	)
	defer close(r.syncWaitingPulses)

	pn := gen.PulseNumber()
	msg := message.HeavyPayload{
		JetID:    gen.JetID(),
		PulseNum: gen.PulseNumber(),
	}

	jetID := gen.JetID()
	jc.MineForPulseMock.Expect(ctx, pn).Return([]insolar.JetID{jetID})
	dg.ForPulseAndJetMock.Return(&msg, nil)
	mb.SendMock.Return(&reply.OK{}, nil)
	c.NotifyAboutPulseMock.Expect(ctx, pn)

	go r.sync(ctx)
	r.syncWaitingPulses <- pn

	ctrl.Wait(time.Minute)
	ctrl.Finish()
}

func TestLightReplicatorDefault_NotifyAboutPulse(t *testing.T) {
	t.Parallel()
	ctrl := minimock.NewController(t)
	ctx := inslogger.TestContext(t)
	jc := jet.NewCalculatorMock(ctrl)
	c := NewCleanerMock(ctrl)
	mb := testutils.NewMessageBusMock(ctrl)

	inputPN := gen.PulseNumber()
	expectedPN := gen.PulseNumber()

	pc := pulse.NewCalculatorMock(ctrl)
	pc.BackwardsMock.Expect(ctx, inputPN, 1).Return(insolar.Pulse{PulseNumber: expectedPN}, nil)
	dg := NewDataGathererMock(ctrl)
	r := NewReplicatorDefault(
		jc,
		dg,
		c,
		mb,
		pc,
	)
	defer close(r.syncWaitingPulses)

	msg := message.HeavyPayload{
		JetID:    gen.JetID(),
		PulseNum: gen.PulseNumber(),
	}

	jetID := gen.JetID()
	jc.MineForPulseMock.Expect(ctx, expectedPN).Return([]insolar.JetID{jetID})
	dg.ForPulseAndJetMock.Return(&msg, nil)
	mb.SendMock.Return(&reply.OK{}, nil)
	c.NotifyAboutPulseMock.Expect(ctx, expectedPN)

	go r.NotifyAboutPulse(ctx, inputPN)

	ctrl.Wait(time.Minute)
	ctrl.Finish()
}
