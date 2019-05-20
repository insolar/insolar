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
	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/flow/internal/thread"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestNewDispatcher(t *testing.T) {
	t.Parallel()
	ok := false
	var f flow.MakeHandle = func(bus.Message) flow.Handle {
		ok = true
		return nil
	}
	require.False(t, ok)

	d := NewDispatcher(f, f)
	require.NotNil(t, d.controller)
	handle := d.handles.present(bus.Message{})
	require.Nil(t, handle)
	require.True(t, ok)
}

type replyMock int

func (replyMock) Type() insolar.ReplyType {
	return insolar.ReplyType(42)
}

func TestDispatcher_WrapBusHandle(t *testing.T) {
	t.Parallel()

	d := &Dispatcher{
		controller:         thread.NewController(),
		currentPulseNumber: insolar.FirstPulseNumber,
	}
	reply := bus.Reply{
		Reply: replyMock(42),
	}
	d.handles.present = func(msg bus.Message) flow.Handle {
		return func(ctx context.Context, f flow.Flow) error {
			msg.ReplyTo <- reply
			return nil
		}
	}
	parcel := &testutils.ParcelMock{}
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	result, err := d.WrapBusHandle(context.Background(), parcel)
	require.NoError(t, err)
	require.Equal(t, reply.Reply, result)
}

func TestDispatcher_WrapBusHandle_Error(t *testing.T) {
	t.Parallel()

	d := &Dispatcher{
		controller:         thread.NewController(),
		currentPulseNumber: insolar.FirstPulseNumber,
	}
	d.handles.present = func(msg bus.Message) flow.Handle {
		return func(ctx context.Context, f flow.Flow) error {
			return errors.New("test error")
		}
	}
	parcel := &testutils.ParcelMock{}
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	result, err := d.WrapBusHandle(context.Background(), parcel)
	require.EqualError(t, err, "test error")
	require.Nil(t, result)
}

func TestDispatcher_WrapBusHandle_ReplyError(t *testing.T) {
	t.Parallel()

	d := &Dispatcher{
		controller:         thread.NewController(),
		currentPulseNumber: insolar.FirstPulseNumber,
	}
	d.handles.present = func(msg bus.Message) flow.Handle {
		return func(ctx context.Context, f flow.Flow) error {
			msg.ReplyTo <- bus.Reply{
				Err: errors.New("reply error"),
			}
			return errors.New("test error")
		}
	}
	parcel := &testutils.ParcelMock{}
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	result, err := d.WrapBusHandle(context.Background(), parcel)
	require.EqualError(t, err, "reply error")
	require.Nil(t, result)
}

func TestDispatcher_WrapBusHandle_NoReply(t *testing.T) {
	t.Parallel()

	d := &Dispatcher{
		controller:         thread.NewController(),
		currentPulseNumber: insolar.FirstPulseNumber,
	}
	d.handles.present = func(msg bus.Message) flow.Handle {
		return func(ctx context.Context, f flow.Flow) error {
			return nil
		}
	}
	parcel := &testutils.ParcelMock{}
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	result, err := d.WrapBusHandle(context.Background(), parcel)
	require.EqualError(t, err, "no reply from handler")
	require.Nil(t, result)
}

func TestDispatcher_WrapBusHandle_ReplyWithError(t *testing.T) {
	t.Parallel()

	d := &Dispatcher{
		controller:         thread.NewController(),
		currentPulseNumber: insolar.FirstPulseNumber,
	}
	reply := bus.Reply{
		Reply: replyMock(42),
	}
	d.handles.present = func(msg bus.Message) flow.Handle {
		return func(ctx context.Context, f flow.Flow) error {
			msg.ReplyTo <- reply
			return errors.New("test error")
		}
	}
	parcel := &testutils.ParcelMock{}
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	result, err := d.WrapBusHandle(context.Background(), parcel)
	require.NoError(t, err)
	require.Equal(t, reply.Reply, result)
}

func TestDispatcher_WrapBusHandle_CallFutureDispatcher(t *testing.T) {
	t.Parallel()
	d := &Dispatcher{
		controller:         thread.NewController(),
		currentPulseNumber: 0,
	}

	reply := bus.Reply{
		Reply: replyMock(42),
	}

	d.handles.future = func(msg bus.Message) flow.Handle {
		return func(ctx context.Context, f flow.Flow) error {
			msg.ReplyTo <- reply
			return nil
		}
	}
	parcel := &testutils.ParcelMock{}
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	result, err := d.WrapBusHandle(context.Background(), parcel)
	require.NoError(t, err)
	require.Equal(t, reply.Reply, result)
}

func makeWMMessage(ctx context.Context, payLoad watermillMsg.Payload) *watermillMsg.Message {
	wmMsg := watermillMsg.NewMessage(watermill.NewUUID(), payLoad)
	wmMsg.Metadata.Set("TraceID", inslogger.TraceID(ctx))

	return wmMsg
}

func TestDispatcher_InnerSubscriber(t *testing.T) {
	t.Parallel()
	d := &Dispatcher{
		controller:         thread.NewController(),
		currentPulseNumber: insolar.FirstPulseNumber,
	}

	testResult := 77
	result := make(chan int)

	d.handles.present = func(msg bus.Message) flow.Handle {
		return func(ctx context.Context, f flow.Flow) error {
			result <- testResult
			return nil
		}
	}

	_, err := d.InnerSubscriber(makeWMMessage(context.Background(), nil))
	require.NoError(t, err)
	require.Equal(t, testResult, <-result)
}

func TestDispatcher_InnerSubscriber_Error(t *testing.T) {
	t.Parallel()
	d := &Dispatcher{
		controller:         thread.NewController(),
		currentPulseNumber: insolar.FirstPulseNumber,
	}
	testResult := 77
	result := make(chan int)

	d.handles.present = func(msg bus.Message) flow.Handle {
		return func(ctx context.Context, f flow.Flow) error {
			result <- testResult
			return errors.New("some error.")
		}
	}
	_, err := d.InnerSubscriber(makeWMMessage(context.Background(), nil))
	require.NoError(t, err)
	require.Equal(t, testResult, <-result)
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
