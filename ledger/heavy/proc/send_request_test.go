// Copyright 2020 Insolar Network Ltd.
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

package proc_test

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock/v3"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/proc"
	"github.com/insolar/insolar/ledger/object"
	"github.com/stretchr/testify/require"
)

func TestGetRequest_Proceed(t *testing.T) {
	mc := minimock.NewController(t)
	ctx := inslogger.TestContext(t)

	var (
		sender  *bus.SenderMock
		records *object.RecordAccessorMock
	)

	resetComponents := func() {
		sender = bus.NewSenderMock(mc)
		records = object.NewRecordAccessorMock(t)
	}

	newProc := func(msg payload.Meta) *proc.SendRequest {
		p := proc.NewSendRequest(msg)
		p.Dep(records, sender)
		return p
	}

	resetComponents()
	t.Run("request does not exist", func(t *testing.T) {
		sender.ReplyMock.Set(func(_ context.Context, _ payload.Meta, msg *message.Message) {
			rep := payload.Error{}
			err := rep.Unmarshal(msg.Payload)
			require.NoError(t, err)
			require.Equal(t, rep.Code, payload.CodeNotFound)
			require.Equal(t, rep.Text, object.ErrNotFound.Error())
		})
		p := newProc(payload.Meta{})
		records.ForIDMock.Return(record.Material{}, object.ErrNotFound)

		err := p.Proceed(ctx)
		require.NoError(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("happy basic", func(t *testing.T) {
		reqID := gen.ID()
		msg := payload.GetRequest{
			RequestID: reqID,
		}
		buf, err := msg.Marshal()
		require.NoError(t, err)
		receivedMeta := payload.Meta{Payload: buf}
		p := newProc(receivedMeta)

		ref := gen.Reference()
		req := record.Virtual{
			Union: &record.Virtual_IncomingRequest{
				IncomingRequest: &record.IncomingRequest{
					Object: &ref,
				},
			},
		}

		records.ForIDMock.Return(record.Material{
			Virtual: req,
		}, nil)
		sender.ReplyMock.Set(func(_ context.Context, origin payload.Meta, rep *message.Message) {
			require.Equal(t, receivedMeta, origin)

			resp, err := payload.Unmarshal(rep.Payload)
			require.NoError(t, err)

			res, ok := resp.(*payload.Request)
			require.True(t, ok)
			require.Equal(t, msg.RequestID, res.RequestID)
			require.Equal(t, req, res.Request)
		})

		err = p.Proceed(ctx)
		require.NoError(t, err)

		mc.Finish()
	})
}
