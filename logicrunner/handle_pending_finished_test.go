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

	"github.com/fortytw2/leaktest"
	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/writecontroller"
)

func TestHandlePendingFinished_Present(t *testing.T) {
	defer leaktest.Check(t)()

	tests := []struct {
		name  string
		mocks func(t minimock.Tester) (*HandlePendingFinished, flow.Flow)
		error bool
	}{
		{
			name: "success",
			mocks: func(t minimock.Tester) (*HandlePendingFinished, flow.Flow) {
				obj := gen.Reference()
				receivedPayload := payload.PendingFinished{
					ObjectRef: obj,
				}

				buf, err := payload.Marshal(&receivedPayload)
				require.NoError(t, err, "marshal")

				h := &HandlePendingFinished{
					dep: &Dependencies{
						Sender: bus.NewSenderMock(t).ReplyMock.Return(),
						StateStorage: NewStateStorageMock(t).
							UpsertExecutionStateMock.Expect(obj).
							Return(
								NewExecutionBrokerIMock(t).
									PrevExecutorSentPendingFinishedMock.Return(nil),
							),
						WriteAccessor: writecontroller.NewWriteControllerMock(t).BeginMock.Return(func() {}, nil),
					},
					Message: payload.Meta{
						Payload: buf,
					},
				}
				return h, flow.NewFlowMock(t)
			},
		},
		{
			name: "error",
			mocks: func(t minimock.Tester) (*HandlePendingFinished, flow.Flow) {
				obj := gen.Reference()
				receivedPayload := payload.PendingFinished{
					ObjectRef: obj,
				}

				buf, err := payload.Marshal(&receivedPayload)
				require.NoError(t, err, "marshal")

				h := &HandlePendingFinished{
					dep: &Dependencies{
						StateStorage: NewStateStorageMock(t).
							UpsertExecutionStateMock.Expect(obj).
							Return(
								NewExecutionBrokerIMock(t).
									PrevExecutorSentPendingFinishedMock.Return(errors.New("some")),
							),
						WriteAccessor: writecontroller.NewWriteControllerMock(t).BeginMock.Return(func() {}, nil),
					},
					Message: payload.Meta{
						Payload: buf,
					},
				}
				return h, flow.NewFlowMock(t)
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
