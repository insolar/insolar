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

package messagebus

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/bus/meta"
	busMeta "github.com/insolar/insolar/insolar/bus/meta"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/reply"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
)

var testType = insolar.MessageType(123)

type replyMock int

func (replyMock) Type() insolar.ReplyType {
	return insolar.ReplyType(124)
}

var testReply replyMock = 124

func testHandler(_ context.Context, _ insolar.Parcel) (insolar.Reply, error) {
	return testReply, nil
}

func prepare(t *testing.T, ctx context.Context, currentPulse int, msgPulse int) (*MessageBus, *pulse.AccessorMock, insolar.Parcel, insolar.Reference) {
	mb, err := NewMessageBus(configuration.Configuration{})
	require.NoError(t, err)

	net := testutils.GetTestNetwork(t)
	jc := jet.NewCoordinatorMock(t)
	expectedRef := gen.Reference()
	jc.QueryRoleMock.Set(func(p context.Context, p1 insolar.DynamicRole, p2 insolar.ID, p3 insolar.PulseNumber) (r []insolar.Reference, r1 error) {
		return []insolar.Reference{expectedRef}, nil
	})

	nn := network.NewNodeNetworkMock(t)
	nn.GetOriginMock.Set(func() (r insolar.NetworkNode) {
		n := network.NewNetworkNodeMock(t)
		n.IDMock.Return(*insolar.NewEmptyReference())
		return n
	})

	pcs := testutils.NewPlatformCryptographyScheme()
	cs := testutils.NewCryptographyServiceMock(t)
	cs.SignMock.Set(func(p []byte) (r *insolar.Signature, r1 error) {
		return &insolar.Signature{}, nil
	})

	dtf := testutils.NewDelegationTokenFactoryMock(t)
	pf := NewParcelFactory()
	ps := pulse.NewAccessorMock(t)

	b := bus.NewSenderMock(t)

	(&component.Manager{}).Inject(net, jc, nn, pcs, cs, dtf, pf, ps, mb, b)

	ps.LatestMock.Set(func(ctx context.Context) (insolar.Pulse, error) {
		return insolar.Pulse{
			PulseNumber:     insolar.PulseNumber(currentPulse),
			NextPulseNumber: insolar.PulseNumber(currentPulse + 1),
		}, nil
	})

	err = mb.Register(testType, testHandler)
	require.NoError(t, err)

	parcel := testutils.NewParcelMock(t)

	parcel.PulseMock.Set(func() insolar.PulseNumber {
		return insolar.PulseNumber(msgPulse)
	})
	parcel.TypeMock.Set(func() insolar.MessageType {
		return testType
	})
	parcel.GetSenderMock.Set(func() (r insolar.Reference) {
		return gen.Reference()
	})
	parcel.MessageMock.Return(&message.GenesisRequest{})

	return mb, ps, parcel, expectedRef
}

func TestMessageBus_doDeliver_SamePulse(t *testing.T) {
	ctx := context.Background()
	mb, _, parcel, _ := prepare(t, ctx, 100, 100)

	result, err := mb.doDeliver(ctx, parcel)
	require.NoError(t, err)
	require.Equal(t, testReply, result)
}

func TestMessageBus_doDeliver_NextPulse(t *testing.T) {
	ctx := context.Background()
	mb, ps, parcel, _ := prepare(t, ctx, 100, 101)

	pulseUpdated := false

	var triggerUnlock int32
	newPulse := insolar.Pulse{
		PulseNumber:     101,
		NextPulseNumber: 102,
	}
	fn := func(ctx context.Context) (insolar.Pulse, error) {
		return insolar.Pulse{
			PulseNumber:     insolar.PulseNumber(100),
			NextPulseNumber: insolar.PulseNumber(101),
		}, nil
	}
	ps.LatestMock.Set(func(ctx context.Context) (insolar.Pulse, error) {
		if atomic.LoadInt32(&triggerUnlock) > 0 {
			return newPulse, nil
		}
		return fn(ctx)
	})
	go func() {
		// should unlock
		time.Sleep(time.Second)
		atomic.AddInt32(&triggerUnlock, 1)

		pulseUpdated = true
		err := mb.OnPulse(ctx, newPulse)
		require.NoError(t, err)
	}()
	// blocks until newPulse returns
	result, err := mb.doDeliver(ctx, parcel)
	require.NoError(t, err)
	require.Equal(t, testReply, result)
	require.True(t, pulseUpdated)
}

func TestMessageBus_doDeliver_TwoAheadPulses(t *testing.T) {
	ctx := context.Background()
	mb, ps, parcel, _ := prepare(t, ctx, 100, 102)

	var mu sync.Mutex
	pulse := &insolar.Pulse{
		PulseNumber:     100,
		NextPulseNumber: 101,
	}
	ps.LatestMock.Set(func(ctx context.Context) (insolar.Pulse, error) {
		mu.Lock()
		defer mu.Unlock()

		return *pulse, nil
	})
	go func() {
		for i := 1; i <= 2; i++ {
			mu.Lock()
			pulse = &insolar.Pulse{
				PulseNumber:     insolar.PulseNumber(100 + i),
				NextPulseNumber: insolar.PulseNumber(100 + i + 1),
			}
			mu.Unlock()
			err := mb.OnPulse(ctx, *pulse)
			require.NoError(t, err)
		}
	}()

	_, err := mb.doDeliver(ctx, parcel)
	require.NoError(t, err)
	require.Equal(t, insolar.PulseNumber(102), pulse.PulseNumber)
}

func TestMessageBus_createWatermillMessage(t *testing.T) {
	ctx := context.Background()
	mb, _, _, _ := prepare(t, ctx, 100, 100)

	pulse := insolar.Pulse{
		PulseNumber: insolar.PulseNumber(100),
	}
	parcel := &message.Parcel{
		Msg: &message.GenesisRequest{},
	}

	msg := mb.createWatermillMessage(ctx, parcel, pulse)

	require.NotNil(t, msg)
	require.NotNil(t, msg.Payload)
	require.Equal(t, parcel.Msg.Type().String(), msg.Metadata.Get(meta.Type))
	require.Equal(t, insolar.NewEmptyReference().String(), msg.Metadata.Get(meta.Sender))
}

func TestMessageBus_deserializePayload_GetReply(t *testing.T) {
	rep := &reply.OK{}
	pl := reply.ToBytes(rep)
	meta := payload.Meta{
		Payload: pl,
	}
	buf, err := meta.Marshal()
	require.NoError(t, err)
	msg := watermillMsg.NewMessage(watermill.NewUUID(), buf)
	msg.Metadata.Set(busMeta.Type, busMeta.TypeReply)

	r, err := deserializePayload(msg)

	require.NoError(t, err)
	require.Equal(t, rep, r)
}

func TestMessageBus_deserializePayload_GetError(t *testing.T) {
	rep := errors.New("test error for deserializePayload")

	msg, err := payload.NewMessage(&payload.Error{Text: rep.Error()})
	require.NoError(t, err)
	meta := payload.Meta{
		Payload: msg.Payload,
	}
	buf, err := meta.Marshal()
	require.NoError(t, err)
	msg.Payload = buf

	r, err := deserializePayload(msg)

	require.Equal(t, rep.Error(), err.Error())
	require.Nil(t, r)
}

func TestMessageBus_deserializePayload_GetReply_WrongType(t *testing.T) {
	meta := payload.Meta{
		Payload: []byte{1, 2, 3, 4},
	}
	buf, err := meta.Marshal()
	require.NoError(t, err)
	msg := watermillMsg.NewMessage(watermill.NewUUID(), buf)
	msg.Metadata.Set(busMeta.Type, busMeta.TypeReply)

	r, err := deserializePayload(msg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "can't deserialize payload to reply")
	require.Nil(t, r)
}

func TestMessageBus_deserializePayload_GetError_WrongBytes(t *testing.T) {
	msg, err := payload.NewMessage(&payload.GetCode{})
	require.NoError(t, err)

	r, err := deserializePayload(msg)

	require.Error(t, err)
	require.Contains(t, err.Error(), "wrong polymorph field number")
	require.Nil(t, r)
}
func TestMessageBus_deserializePayload_Nil(t *testing.T) {
	r, err := deserializePayload(nil)

	require.Error(t, err)
	require.Contains(t, err.Error(), "can't deserialize payload of nil message")
	require.Nil(t, r)
}
