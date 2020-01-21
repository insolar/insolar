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

package executor_test

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
)

func TestDetachedNotifier_Notify(t *testing.T) {
	t.Parallel()

	t.Run("Simple success", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()

		ctx := flow.TestContextWithPulse(
			inslogger.TestContext(t),
			gen.PulseNumber(),
		)

		objectID := gen.ID()
		closedReqID := gen.ID()
		detachedReqID := gen.ID()
		detachedReq := record.Wrap(&record.OutgoingRequest{
			Reason:     *insolar.NewReference(closedReqID),
			ReturnMode: record.ReturnSaga,
		})
		detachedReqBuf, _ := detachedReq.Marshal()

		opened := []record.CompositeFilamentRecord{
			// wrong
			{
				RecordID: closedReqID,
				Record:   record.Material{Virtual: record.Wrap(&record.IncomingRequest{})},
			},
			// wrong
			{
				RecordID: gen.ID(),
				Record:   record.Material{Virtual: record.Wrap(&record.OutgoingRequest{})},
			},
			// right
			{
				RecordID: detachedReqID,
				Record:   record.Material{Virtual: detachedReq},
			},
			// wrong
			{
				RecordID: gen.ID(),
				Record: record.Material{Virtual: record.Wrap(&record.OutgoingRequest{
					Reason:     *insolar.NewReference(gen.ID()),
					ReturnMode: record.ReturnSaga,
				})},
			},
		}

		sender := bus.NewSenderMock(mc)
		expectedToVirtualMsg, _ := payload.NewMessage(&payload.SagaCallAcceptNotification{
			ObjectID:          objectID,
			DetachedRequestID: detachedReqID,
			Request:           detachedReqBuf,
		})

		sender.SendRoleMock.Inspect(func(ctx context.Context, msg *message.Message, role insolar.DynamicRole, object insolar.Reference) {
			require.Equal(t, expectedToVirtualMsg.Payload, msg.Payload)
			require.Equal(t, insolar.DynamicRoleVirtualExecutor, role)
			require.Equal(t, *insolar.NewReference(objectID), object)
		}).Return(make(chan *message.Message), func() {})

		dn := executor.NewDetachedNotifierDefault(sender)
		dn.Notify(ctx, opened, objectID, closedReqID)
	})
}
