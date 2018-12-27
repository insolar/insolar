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
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/require"
)

var testType = core.MessageType(123)

var replyOk = &reply.OK{}
var replyWrong = &reply.WrongPulseNumber{
	CurrentPulse: &core.Pulse{
		PrevPulseNumber: core.PulseNumber(100),
		PulseNumber:     core.PulseNumber(101),
		NextPulseNumber: core.PulseNumber(102),
	},
}
var ref = testutils.RandomRef()
var opts = &core.MessageSendOptions{
	Receiver: &ref,
}

type tester struct {
	t            *testing.T
	ctx          context.Context
	mb           core.MessageBus
	currentPulse core.Pulse
	messagePulse core.PulseNumber
	reply        core.Reply
}

func (t *tester) updatePulse(pulse uint32) {
	t.currentPulse = core.Pulse{
		PrevPulseNumber: t.currentPulse.PulseNumber,
		PulseNumber:     core.PulseNumber(pulse),
		NextPulseNumber: core.PulseNumber(pulse + 1),
	}
	t.mb.OnPulse(t.ctx, t.currentPulse)
}

func (t *tester) setReply(rep core.Reply) {
	t.reply = rep
}

func prepare(t *testing.T, ctx context.Context, currentPulse uint32, messagePulse uint32) (*MessageBus, *tester, core.Parcel) {
	mb, err := NewMessageBus(configuration.Configuration{})
	require.NoError(t, err)
	zz := &tester{
		t:   t,
		ctx: ctx,
		mb:  mb,
		currentPulse: core.Pulse{
			PrevPulseNumber: core.PulseNumber(currentPulse - 1),
			PulseNumber:     core.PulseNumber(currentPulse),
			NextPulseNumber: core.PulseNumber(currentPulse + 1),
		},
		messagePulse: core.PulseNumber(messagePulse),
		reply:        replyOk,
	}

	net := testutils.NewNetworkMock(t)
	jc := testutils.NewJetCoordinatorMock(t)
	ls := testutils.NewLocalStorageMock(t)
	nn := network.NewNodeNetworkMock(t)
	pcs := testutils.NewPlatformCryptographyScheme()
	cs := testutils.NewCryptographyServiceMock(t)
	dtf := testutils.NewDelegationTokenFactoryMock(t)
	pf := NewParcelFactory()
	ps := testutils.NewPulseStorageMock(t)

	(&component.Manager{}).Inject(net, jc, ls, nn, pcs, cs, dtf, pf, ps, mb)

	ps.CurrentFunc = func(context.Context) (*core.Pulse, error) {
		return &zz.currentPulse, nil
	}

	net.SendMessageFunc = func(core.RecordRef, string, core.Parcel) ([]byte, error) {
		rd, err := reply.Serialize(zz.reply)
		require.NoError(t, err)
		return rd.(*bytes.Buffer).Bytes(), nil
	}

	nn.GetOriginFunc = func() (r core.Node) {
		node := network.NewNodeMock(t)
		node.IDFunc = func() (r core.RecordRef) {
			return testutils.RandomRef()
		}
		return node
	}

	err = mb.Register(testType, func(context.Context, core.Parcel) (core.Reply, error) {
		return replyOk, nil
	})
	require.NoError(t, err)

	parcel := testutils.NewParcelMock(t)
	parcel.PulseFunc = func() core.PulseNumber {
		return zz.messagePulse
	}
	parcel.TypeFunc = func() core.MessageType {
		return testType
	}
	parcel.UpdatePulseFunc = func(newPulse core.PulseNumber) {
		zz.messagePulse = newPulse
	}

	mb.Unlock(ctx)

	return mb, zz, parcel
}

func TestMessageBus_doDeliverSamePulse(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mb, _, parcel := prepare(t, ctx, 100, 100)

	result, err := mb.doDeliver(ctx, parcel)
	require.NoError(t, err)
	require.Equal(t, replyOk, result)
}

func TestMessageBus_doDeliverNextPulse(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mb, zz, parcel := prepare(t, ctx, 100, 101)

	pulseUpdated := false

	go func() {
		time.Sleep(time.Second)
		pulseUpdated = true
		zz.updatePulse(101)
	}()
	result, err := mb.doDeliver(ctx, parcel)
	require.NoError(t, err)
	require.Equal(t, replyOk, result)
	require.True(t, pulseUpdated)
}

func TestMessageBus_doDeliverWrongPulse(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mb, _, parcel := prepare(t, ctx, 100, 200)

	_, err := mb.doDeliver(ctx, parcel)
	require.EqualError(t, err, "[ MessageBus ] Incorrect message pulse 200 100")
}

func TestMessageBus_SendWithSamePulse(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mb, zz, parcel := prepare(t, ctx, 100, 100)

	pulse := zz.currentPulse
	res, err := mb.SendParcel(ctx, parcel, pulse, opts)
	require.NoError(t, err)
	require.Equal(t, replyOk, res)
}

func TestMessageBus_SendWithPrevPulse(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mb, zz, parcel := prepare(t, ctx, 100, 101)

	pulseUpdated := false

	go func() {
		time.Sleep(time.Second)
		zz.setReply(replyOk)
		pulseUpdated = true
		zz.updatePulse(101)
	}()

	zz.setReply(replyWrong)
	pulse := zz.currentPulse
	res, err := mb.SendParcel(ctx, parcel, pulse, opts)
	require.NoError(t, err)
	require.True(t, pulseUpdated)
	require.Equal(t, replyOk, res)
}
