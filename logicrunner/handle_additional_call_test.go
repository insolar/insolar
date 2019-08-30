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

package logicrunner

import (
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/writecontroller"
)

func TestHandleAdditionalCallFromPreviousExecutor_Present(t *testing.T) {
	table := []struct {
		name                      string
		clarifyPendingStateResult error
		startQueueProcessorResult error
		mocks                     func(t minimock.Tester) (*HandleAdditionalCallFromPreviousExecutor, flow.Flow)
		error                     bool
	}{
		{
			name: "Happy path",
			mocks: func(t minimock.Tester) (*HandleAdditionalCallFromPreviousExecutor, flow.Flow) {
				obj := gen.Reference()
				reqRef := gen.Reference()

				receivedPayload := &payload.AdditionalCallFromPreviousExecutor{
					RequestRef:      reqRef,
					Pending:         insolar.NotPending,
					ObjectReference: obj,
					Request:         &record.IncomingRequest{},
					ServiceData:     &payload.ServiceData{},
				}

				buf, err := payload.Marshal(receivedPayload)
				require.NoError(t, err, "marshal")

				h := &HandleAdditionalCallFromPreviousExecutor{
					dep: &Dependencies{
						Sender: bus.NewSenderMock(t).ReplyMock.Return(),
						WriteAccessor: writecontroller.NewAccessorMock(t).
							BeginMock.Return(func() {}, nil),
					},
					Message: payload.Meta{Payload: buf},
				}
				f := flow.NewFlowMock(t).ProcedureMock.Return(nil)
				return h, f
			},
		},
	}

	for _, test := range table {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
			mc := minimock.NewController(t)

			h, f := test.mocks(mc)
			err := h.Present(ctx, f)
			if test.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mc.Wait(2 * time.Minute)
			mc.Finish()
		})
	}
}

func TestAdditionalCallFromPreviousExecutor_Proceed(t *testing.T) {

	t.Run("Proceed without pending", func(t *testing.T) {
		t.Parallel()

		mc := minimock.NewController(t)

		ctx := inslogger.TestContext(t)
		obj := gen.Reference()
		reqRef := gen.Reference()

		msg := &payload.AdditionalCallFromPreviousExecutor{
			ObjectReference: obj,
			RequestRef:      reqRef,
			Request: &record.IncomingRequest{
				Object: &obj,
			},
			Pending: insolar.NotPending,
		}

		stateStorage := NewStateStorageMock(t).
			UpsertExecutionStateMock.Expect(obj).Return(
			NewExecutionBrokerIMock(t).
				AddAdditionalRequestFromPrevExecutorMock.Return().
				SetNotPendingMock.Return(),
		)

		proc := AdditionalCallFromPreviousExecutor{stateStorage: stateStorage, message: msg}
		err := proc.Proceed(ctx)

		require.NoError(t, err)

		mc.Wait(2 * time.Minute)
		mc.Finish()
	})

	t.Run("Proceed with pending", func(t *testing.T) {
		t.Parallel()

		mc := minimock.NewController(t)

		ctx := inslogger.TestContext(t)
		obj := gen.Reference()
		reqRef := gen.Reference()

		msg := &payload.AdditionalCallFromPreviousExecutor{
			ObjectReference: obj,
			RequestRef:      reqRef,
			Request: &record.IncomingRequest{
				Object: &obj,
			},
			Pending: insolar.InPending,
		}

		stateStorage := NewStateStorageMock(t).
			UpsertExecutionStateMock.Expect(obj).Return(
			NewExecutionBrokerIMock(t).
				AddAdditionalRequestFromPrevExecutorMock.Return(),
		)

		proc := AdditionalCallFromPreviousExecutor{stateStorage: stateStorage, message: msg}
		err := proc.Proceed(ctx)

		require.NoError(t, err)

		mc.Wait(2 * time.Minute)
		mc.Finish()
	})
}
