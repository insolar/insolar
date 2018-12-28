/*
 *    Copyright 2018 Insolar
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

package messagebus

import (
	"context"
	"testing"
	"time"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/require"
)

var testType = core.MessageType(123)

type replyMock int

func (replyMock) Type() core.ReplyType {
	return core.ReplyType(124)
}

var testReply replyMock = 124

func testHandler(ctx context.Context, msg core.Parcel) (core.Reply, error) {
	return testReply, nil
}

func prepare(t *testing.T, ctx context.Context, currentPulse int, msgPulse int) (*MessageBus, *testutils.PulseStorageMock, core.Parcel) {
	mb, err := NewMessageBus(configuration.Configuration{})
	require.NoError(t, err)

	net := network.GetTestNetwork()
	jc := testutils.NewJetCoordinatorMock(t)
	ls := testutils.NewLocalStorageMock(t)
	nn := network.NewNodeNetworkMock(t)
	pcs := testutils.NewPlatformCryptographyScheme()
	cs := testutils.NewCryptographyServiceMock(t)
	dtf := testutils.NewDelegationTokenFactoryMock(t)
	pf := NewParcelFactory()
	ps := testutils.NewPulseStorageMock(t)

	(&component.Manager{}).Inject(net, jc, ls, nn, pcs, cs, dtf, pf, ps, mb)

	ps.CurrentFunc = func(ctx context.Context) (*core.Pulse, error) {
		return &core.Pulse{
			PulseNumber:     core.PulseNumber(currentPulse),
			NextPulseNumber: core.PulseNumber(currentPulse + 1),
		}, nil
	}

	err = mb.Register(testType, testHandler)
	require.NoError(t, err)

	parcel := testutils.NewParcelMock(t)

	parcel.PulseFunc = func() core.PulseNumber {
		return core.PulseNumber(msgPulse)
	}
	parcel.TypeFunc = func() core.MessageType {
		return testType
	}

	mb.Unlock(ctx)

	return mb, ps, parcel
}

func TestMessageBus_doDeliverSamePulse(t *testing.T) {
	ctx := context.Background()
	mb, _, parcel := prepare(t, ctx, 100, 100)

	result, err := mb.doDeliver(ctx, parcel)
	require.NoError(t, err)
	require.Equal(t, testReply, result)
}

func TestMessageBus_doDeliverNextPulse(t *testing.T) {
	ctx := context.Background()
	mb, ps, parcel := prepare(t, ctx, 100, 101)

	pulseUpdated := false

	go func() {
		time.Sleep(time.Second)
		newPulse := &core.Pulse{
			PulseNumber:     101,
			NextPulseNumber: 102,
		}
		ps.CurrentFunc = func(ctx context.Context) (*core.Pulse, error) {
			return newPulse, nil
		}
		pulseUpdated = true
		mb.OnPulse(ctx, *newPulse)
	}()
	result, err := mb.doDeliver(ctx, parcel)
	require.NoError(t, err)
	require.Equal(t, testReply, result)
	require.True(t, pulseUpdated)
}

func TestMessageBus_doDeliverWrongPulse(t *testing.T) {
	ctx := context.Background()
	mb, _, parcel := prepare(t, ctx, 100, 200)

	_, err := mb.doDeliver(ctx, parcel)
	require.EqualError(t, err, "[ MessageBus ] Incorrect message pulse 200 100")
}
