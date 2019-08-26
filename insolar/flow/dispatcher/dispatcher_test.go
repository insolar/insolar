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

package dispatcher

import (
	"context"
	"fmt"
	"testing"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus/meta"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/internal/thread"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/instracer"
)

func makeMessage(t *testing.T, ctx context.Context, pn insolar.PulseNumber) *message.Message {
	payload := []byte{1, 2, 3, 4, 5}
	msg := message.NewMessage(watermill.NewUUID(), payload)
	msg.Metadata.Set(meta.Pulse, pn.String())
	sp, err := instracer.Serialize(ctx)
	require.NoError(t, err)
	msg.Metadata.Set(meta.SpanData, string(sp))

	return msg
}

func TestNewDispatcher(t *testing.T) {
	t.Parallel()
	ok := false
	var f flow.MakeHandle = func(*message.Message) flow.Handle {
		ok = true
		return nil
	}
	require.False(t, ok)

	dInterface := NewDispatcher(nil, f, f, f)
	d := dInterface.(*dispatcher)
	require.NotNil(t, d.controller)

	ctx := context.Background()
	currentPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(100)}
	d.pulses = pulse.NewAccessorMock(t).LatestMock.Return(currentPulse, nil)

	msg := makeMessage(t, ctx, currentPulse.PulseNumber)

	handle := d.handles.present(msg)
	require.Nil(t, handle)
	require.True(t, ok)
}

type replyMock int

func (replyMock) Type() insolar.ReplyType {
	return insolar.ReplyType(42)
}

func TestDispatcher_Process(t *testing.T) {
	t.Parallel()

	d := &dispatcher{
		controller: thread.NewController(),
	}
	reply := replyMock(42)
	replyChan := make(chan insolar.Reply, 1)
	d.handles.present = func(msg *message.Message) flow.Handle {
		return func(ctx context.Context, f flow.Flow) error {
			replyChan <- reply
			return nil
		}
	}

	ctx := context.Background()
	currentPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(100)}
	d.pulses = pulse.NewAccessorMock(t).LatestMock.Return(currentPulse, nil)

	msg := makeMessage(t, ctx, currentPulse.PulseNumber)

	err := d.Process(msg)
	require.NoError(t, err)
	rep := <-replyChan
	require.Equal(t, reply, rep)
}

func TestDispatcher_Process_ReplyError(t *testing.T) {
	t.Parallel()

	d := &dispatcher{
		controller: thread.NewController(),
	}
	replyChan := make(chan error, 1)
	d.handles.present = func(msg *message.Message) flow.Handle {
		return func(ctx context.Context, f flow.Flow) error {
			replyChan <- errors.New("reply error")
			return errors.New("test error")
		}
	}

	ctx := context.Background()
	currentPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(100)}
	d.pulses = pulse.NewAccessorMock(t).LatestMock.Return(currentPulse, nil)

	msg := makeMessage(t, ctx, currentPulse.PulseNumber)

	err := d.Process(msg)
	require.NoError(t, err)
	rep := <-replyChan
	require.Error(t, rep)
	require.Contains(t, rep.Error(), "reply error")
}

func TestDispatcher_Process_CallFutureDispatcher(t *testing.T) {
	t.Parallel()
	d := &dispatcher{
		controller: thread.NewController(),
	}

	reply := replyMock(42)
	replyChan := make(chan insolar.Reply, 1)
	d.handles.future = func(msg *message.Message) flow.Handle {
		return func(ctx context.Context, f flow.Flow) error {
			replyChan <- reply
			return nil
		}
	}

	ctx := context.Background()
	currentPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(100)}
	d.pulses = pulse.NewAccessorMock(t).LatestMock.Return(currentPulse, nil)

	msg := makeMessage(t, ctx, currentPulse.PulseNumber+1)

	err := d.Process(msg)
	require.NoError(t, err)
	rep := <-replyChan
	require.Equal(t, reply, rep)
}

func TestDispatcher_pulseFromString(t *testing.T) {
	expectedPulse := insolar.PulseNumber(666)
	pulse, err := pulseFromString(fmt.Sprintf("%d", expectedPulse))
	require.NoError(t, err)
	require.Equal(t, expectedPulse, pulse)
}

func TestDispatcher_pulseFromString_Err(t *testing.T) {
	pulse, err := pulseFromString("test_string")
	require.Error(t, err)
	require.Contains(t, err.Error(), "can't convert string value to pulse")
	require.Equal(t, insolar.PulseNumber(0), pulse)
}
