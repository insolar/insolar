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
	"errors"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/writecontroller"
)

func TestHandleExecutorResults_Present(t *testing.T) {
	tests := []struct {
		name  string
		mocks func(t minimock.Tester) (*HandleExecutorResults, flow.Flow)
		error bool
	}{
		{
			name: "success, every call to broker",
			mocks: func(t minimock.Tester) (*HandleExecutorResults, flow.Flow) {
				incoming1 := genIncomingRequest()
				incoming2 := genIncomingRequest()
				incoming2.Object = incoming1.Object

				receivedPayload := &payload.ExecutorResults{
					RecordRef:             *incoming1.Object,
					Pending:               insolar.NotPending,
					LedgerHasMoreRequests: true,
					Queue: []*payload.ExecutionQueueElement{
						{
							RequestRef:  gen.RecordReference(),
							Incoming:    incoming1,
							ServiceData: &payload.ServiceData{},
						},
						{
							RequestRef:  gen.RecordReference(),
							Incoming:    incoming2,
							ServiceData: &payload.ServiceData{},
						},
					},
				}

				buf, err := payload.Marshal(receivedPayload)
				require.NoError(t, err, "marshal")

				h := &HandleExecutorResults{
					dep: &Dependencies{
						WriteAccessor: writecontroller.NewWriteControllerMock(t).BeginMock.Return(func() {}, nil),
						StateStorage: NewStateStorageMock(t).
							UpsertExecutionStateMock.Expect(*incoming1.Object).
							Return(
								NewExecutionBrokerIMock(t).
									PrevExecutorPendingResultMock.Return().
									MoreRequestsOnLedgerMock.Return(),
							),
					},
					meta: payload.Meta{Payload: buf},
				}
				f := flow.NewFlowMock(t)
				return h, f
			},
		},
		{
			name: "success, minimum calls to broker",
			mocks: func(t minimock.Tester) (*HandleExecutorResults, flow.Flow) {
				obj := gen.Reference()

				receivedPayload := &payload.ExecutorResults{
					RecordRef: obj,
					Pending:   insolar.NotPending,
				}

				buf, err := payload.Marshal(receivedPayload)
				require.NoError(t, err, "marshal")

				h := &HandleExecutorResults{
					dep: &Dependencies{
						WriteAccessor: writecontroller.NewWriteControllerMock(t).BeginMock.Return(func() {}, nil),
						StateStorage: NewStateStorageMock(t).
							UpsertExecutionStateMock.Expect(obj).
							Return(
								NewExecutionBrokerIMock(t).
									PrevExecutorPendingResultMock.Return(),
							),
					},
					meta: payload.Meta{Payload: buf},
				}
				f := flow.NewFlowMock(t)
				return h, f
			},
		},
		{
			name: "write controller is closed",
			mocks: func(t minimock.Tester) (*HandleExecutorResults, flow.Flow) {
				obj := gen.Reference()
				receivedPayload := &payload.ExecutorResults{
					RecordRef: obj,
					Pending:   insolar.NotPending,
				}

				buf, err := payload.Marshal(receivedPayload)
				require.NoError(t, err, "marshal")

				h := &HandleExecutorResults{
					dep: &Dependencies{
						WriteAccessor: writecontroller.NewWriteControllerMock(t).
							BeginMock.Return(nil, writecontroller.ErrWriteClosed),
					},
					meta: payload.Meta{Payload: buf},
				}
				f := flow.NewFlowMock(t)
				return h, f
			},
			error: true,
		},
		{
			name: "write controller error",
			mocks: func(t minimock.Tester) (*HandleExecutorResults, flow.Flow) {
				obj := gen.Reference()
				receivedPayload := &payload.ExecutorResults{
					RecordRef: obj,
					Pending:   insolar.NotPending,
				}

				buf, err := payload.Marshal(receivedPayload)
				require.NoError(t, err, "marshal")

				h := &HandleExecutorResults{
					dep: &Dependencies{
						WriteAccessor: writecontroller.NewWriteControllerMock(t).
							BeginMock.Return(nil, errors.New("some")),
					},
					meta: payload.Meta{Payload: buf},
				}
				f := flow.NewFlowMock(t)
				return h, f
			},
			error: true,
		},
		{
			name: "error, bad data",
			mocks: func(t minimock.Tester) (*HandleExecutorResults, flow.Flow) {
				h := &HandleExecutorResults{
					dep:  &Dependencies{},
					meta: payload.Meta{Payload: []byte{3, 2, 1}},
				}
				f := flow.NewFlowMock(t)
				return h, f
			},
			error: true,
		},
	}
	for _, test := range tests {
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

			mc.Wait(1 * time.Minute)
			mc.Finish()
		})
	}
}
