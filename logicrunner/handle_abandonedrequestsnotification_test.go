package logicrunner

import (
	"errors"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/writecontroller"
)

func TestHandleAbandonedRequestsNotification_Present(t *testing.T) {
	tests := []struct {
		name  string
		mocks func(t minimock.Tester) (*HandleAbandonedRequestsNotification, flow.Flow)
		error bool
	}{
		{
			name: "success",
			mocks: func(t minimock.Tester) (*HandleAbandonedRequestsNotification, flow.Flow) {
				obj := gen.Reference()
				receivedPayload := payload.AbandonedRequestsNotification{
					ObjectID: *obj.GetLocal(),
				}

				buf, err := payload.Marshal(&receivedPayload)
				require.NoError(t, err, "marshal")

				h := &HandleAbandonedRequestsNotification{
					dep: &Dependencies{
						StateStorage: NewStateStorageMock(t).
							UpsertExecutionStateMock.Expect(obj).
							Return(
								NewExecutionBrokerIMock(t).
									AbandonedRequestsOnLedgerMock.Return(),
							),
						WriteAccessor: writecontroller.NewWriteControllerMock(t).
							BeginMock.Return(func() {}, nil),
					},
					meta: payload.Meta{Payload: buf},
				}
				return h, flow.NewFlowMock(t)
			},
		},
		{
			name: "write controller is closed",
			mocks: func(t minimock.Tester) (*HandleAbandonedRequestsNotification, flow.Flow) {
				obj := gen.Reference()
				receivedPayload := payload.AbandonedRequestsNotification{
					ObjectID: *obj.GetLocal(),
				}

				buf, err := payload.Marshal(&receivedPayload)
				require.NoError(t, err, "marshal")

				h := &HandleAbandonedRequestsNotification{
					dep: &Dependencies{
						WriteAccessor: writecontroller.NewWriteControllerMock(t).
							BeginMock.Return(nil, errors.New("some")),
					},
					meta: payload.Meta{Payload: buf},
				}
				return h, flow.NewFlowMock(t)
			},
		},
		{
			name: "error unmarshaling",
			mocks: func(t minimock.Tester) (*HandleAbandonedRequestsNotification, flow.Flow) {
				h := &HandleAbandonedRequestsNotification{
					dep: &Dependencies{},
					meta: payload.Meta{Payload: []byte{3,2,1}},
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
