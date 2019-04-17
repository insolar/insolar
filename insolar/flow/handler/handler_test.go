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

package handler

import (
	"context"
	"runtime"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/flow/internal/thread"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	t.Parallel()
	ok := false
	var f flow.MakeHandle = func(bus.Message) flow.Handle {
		ok = true
		return nil
	}
	require.False(t, ok)
	h := NewHandler(f)
	require.NotNil(t, h.controller)
	handle := h.handles.present(bus.Message{})
	require.Nil(t, handle)
	require.True(t, ok)
}

type replyMock int

func (replyMock) Type() insolar.ReplyType {
	return insolar.ReplyType(42)
}

func TestHandler_WrapBusHandle(t *testing.T) {
	t.Parallel()
	h := &Handler{
		controller: thread.NewController(),
	}
	reply := bus.Reply{
		Reply: replyMock(42),
	}
	h.handles.present = func(msg bus.Message) flow.Handle {
		msg.ReplyTo <- reply
		runtime.Gosched()
		return func(ctx context.Context, f flow.Flow) error {
			return nil
		}
	}
	parcel := &testutils.ParcelMock{}
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	result, err := h.WrapBusHandle(context.Background(), parcel)
	require.NoError(t, err)
	require.Equal(t, reply.Reply, result)
}

func TestHandler_WrapBusHandle_Error(t *testing.T) {
	t.Parallel()
	h := &Handler{
		controller: thread.NewController(),
	}
	h.handles.present = func(msg bus.Message) flow.Handle {
		return func(ctx context.Context, f flow.Flow) error {
			return errors.New("test error")
		}
	}
	parcel := &testutils.ParcelMock{}
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	result, err := h.WrapBusHandle(context.Background(), parcel)
	require.EqualError(t, err, "test error")
	require.Nil(t, result)
}

func TestHandler_WrapBusHandle_ReplyError(t *testing.T) {
	t.Parallel()
	h := &Handler{
		controller: thread.NewController(),
	}
	h.handles.present = func(msg bus.Message) flow.Handle {
		msg.ReplyTo <- bus.Reply{
			Err: errors.New("reply error"),
		}
		runtime.Gosched()
		return func(ctx context.Context, f flow.Flow) error {
			return errors.New("test error")
		}
	}
	parcel := &testutils.ParcelMock{}
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	result, err := h.WrapBusHandle(context.Background(), parcel)
	require.EqualError(t, err, "reply error")
	require.Nil(t, result)
}

func TestHandler_WrapBusHandle_NoReply(t *testing.T) {
	t.Parallel()
	h := &Handler{
		controller: thread.NewController(),
	}
	h.handles.present = func(msg bus.Message) flow.Handle {
		return func(ctx context.Context, f flow.Flow) error {
			return nil
		}
	}
	parcel := &testutils.ParcelMock{}
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	result, err := h.WrapBusHandle(context.Background(), parcel)
	require.EqualError(t, err, "no reply from handler")
	require.Nil(t, result)
}

func TestHandler_WrapBusHandle_ReplyWithError(t *testing.T) {
	t.Parallel()
	h := &Handler{
		controller: thread.NewController(),
	}
	reply := bus.Reply{
		Reply: replyMock(42),
	}
	h.handles.present = func(msg bus.Message) flow.Handle {
		msg.ReplyTo <- reply
		runtime.Gosched()
		return func(ctx context.Context, f flow.Flow) error {
			return errors.New("test error")
		}
	}
	parcel := &testutils.ParcelMock{}
	parcel.PulseFunc = func() insolar.PulseNumber {
		return 42
	}
	result, err := h.WrapBusHandle(context.Background(), parcel)
	require.NoError(t, err)
	require.Equal(t, reply.Reply, result)
}
