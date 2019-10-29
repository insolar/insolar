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
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/writecontroller"
)

func genAPIRequestID() string {
	APIRequestID := utils.RandTraceID()
	if strings.Contains(APIRequestID, "createRandomTraceIDFailed") {
		panic("Failed to generate uuid: " + APIRequestID)
	}
	return APIRequestID
}

func genIncomingRequest() *record.IncomingRequest {
	baseRef := gen.Reference()
	objectRef := gen.Reference()
	prototypeRef := gen.Reference()

	return &record.IncomingRequest{
		Polymorph:       rand.Int31(),
		CallType:        record.CTMethod,
		Caller:          gen.Reference(),
		CallerPrototype: gen.Reference(),
		Nonce:           0,
		ReturnMode:      record.ReturnSaga,
		Immutable:       false,
		Base:            &baseRef,
		Object:          &objectRef,
		Prototype:       &prototypeRef,
		Method:          "Call",
		Arguments:       []byte{},
		APIRequestID:    genAPIRequestID(),
		Reason:          gen.RecordReference(),
		APINode:         insolar.Reference{},
	}
}

func TestHandleAdditionalCallFromPreviousExecutor_Present(t *testing.T) {
	defer leaktest.Check(t)()

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
				incoming := genIncomingRequest()

				receivedPayload := &payload.AdditionalCallFromPreviousExecutor{
					RequestRef:      gen.RecordReference(),
					Pending:         insolar.NotPending,
					ObjectReference: *incoming.Object,
					Request:         incoming,
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
		if useLeakTest {
			defer leaktest.Check(t)()
		} else {
			t.Parallel()
		}

		mc := minimock.NewController(t)

		ctx := inslogger.TestContext(t)
		incoming := genIncomingRequest()

		msg := &payload.AdditionalCallFromPreviousExecutor{
			ObjectReference: *incoming.Object,
			RequestRef:      gen.RecordReference(),
			Request:         incoming,
			Pending:         insolar.NotPending,
		}

		stateStorage := NewStateStorageMock(t).
			UpsertExecutionStateMock.Expect(*incoming.Object).Return(
			NewExecutionBrokerIMock(t).SetNotPendingMock.Return().
				HasMoreRequestsMock.Return(),
		)

		proc := AdditionalCallFromPreviousExecutor{stateStorage: stateStorage, message: msg}
		err := proc.Proceed(ctx)

		require.NoError(t, err)

		mc.Wait(2 * time.Minute)
		mc.Finish()
	})

	t.Run("Proceed with pending", func(t *testing.T) {
		if useLeakTest {
			defer leaktest.Check(t)()
		} else {
			t.Parallel()
		}

		mc := minimock.NewController(t)

		ctx := inslogger.TestContext(t)
		incoming := genIncomingRequest()

		msg := &payload.AdditionalCallFromPreviousExecutor{
			ObjectReference: *incoming.Object,
			RequestRef:      gen.RecordReference(),
			Request:         incoming,
			Pending:         insolar.InPending,
		}

		stateStorage := NewStateStorageMock(t).
			UpsertExecutionStateMock.Expect(*incoming.Object).Return(
			NewExecutionBrokerIMock(t).
				HasMoreRequestsMock.Return(),
		)

		proc := AdditionalCallFromPreviousExecutor{stateStorage: stateStorage, message: msg}
		err := proc.Proceed(ctx)

		require.NoError(t, err)

		mc.Wait(2 * time.Minute)
		mc.Finish()
	})
}
