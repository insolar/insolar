// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package logicrunner

import (
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/writecontroller"
	"github.com/insolar/insolar/testutils"
)

func TestHandleStillExecuting_Present(t *testing.T) {
	defer testutils.LeakTester(t)

	tests := []struct {
		name  string
		mocks func(t minimock.Tester) (*HandleStillExecuting, flow.Flow)
	}{
		{
			name: "success",
			mocks: func(t minimock.Tester) (*HandleStillExecuting, flow.Flow) {
				obj := gen.Reference()
				receivedPayload := &payload.StillExecuting{
					ObjectRef:   obj,
					Executor:    gen.Reference(),
					RequestRefs: []insolar.Reference{gen.RecordReference()},
				}

				buf, err := payload.Marshal(receivedPayload)
				require.NoError(t, err, "marshal")

				h := &HandleStillExecuting{
					dep: &Dependencies{
						StateStorage: NewStateStorageMock(t).
							UpsertExecutionStateMock.Expect(obj).
							Return(
								NewExecutionBrokerIMock(t).
									PrevExecutorStillExecutingMock.Return(nil),
							),
						ResultsMatcher: NewResultMatcherMock(t).
							AddStillExecutionMock.Return(),
						WriteAccessor: writecontroller.NewWriteControllerMock(t).BeginMock.Return(func() {}, nil),
					},
					Message: payload.Meta{Payload: buf},
				}
				return h, flow.NewFlowMock(t)
			},
		},
		{
			name: "not in pending",
			mocks: func(t minimock.Tester) (*HandleStillExecuting, flow.Flow) {
				obj := gen.Reference()
				receivedPayload := &payload.StillExecuting{
					ObjectRef:   obj,
					Executor:    gen.Reference(),
					RequestRefs: []insolar.Reference{gen.RecordReference()},
				}

				buf, err := payload.Marshal(receivedPayload)
				require.NoError(t, err, "marshal")

				h := &HandleStillExecuting{
					dep: &Dependencies{
						StateStorage: NewStateStorageMock(t).
							UpsertExecutionStateMock.Expect(obj).
							Return(
								NewExecutionBrokerIMock(t).
									PrevExecutorStillExecutingMock.Return(ErrNotInPending),
							),
						ResultsMatcher: NewResultMatcherMock(t).
							AddStillExecutionMock.Return(),
						WriteAccessor: writecontroller.NewWriteControllerMock(t).BeginMock.Return(func() {}, nil),
					},
					Message: payload.Meta{Payload: buf},
				}
				return h, flow.NewFlowMock(t)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
			mc := minimock.NewController(t)

			h, f := test.mocks(mc)
			err := h.Present(ctx, f)
			require.NoError(t, err)

			mc.Wait(1 * time.Minute)
			mc.Finish()
		})
	}
}
